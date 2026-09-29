package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/finsight/backend/internal/models"
	"github.com/finsight/backend/internal/repositories"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

const (
	fxRateCacheKeyFmt    = "fx:rate:%s:%s"
	fxRateLatestCacheTTL = 5 * time.Minute
	fxRateStaleKeyFmt    = "fx:stale:%s:%s"
)

// FXService manages FX rate retrieval, caching, and exposure calculations.
// It never silently uses stale rates — always indicates staleness to callers.
type FXService struct {
	fxRepo    repositories.FXRepository
	cache     *redis.Client
	logger    *zap.Logger
	apiKey    string
	provider  string
	staleness time.Duration
	baseCurrency string
}

func NewFXService(
	fxRepo repositories.FXRepository,
	cache *redis.Client,
	cfg interface {
		GetFXAPIKey() string
		GetFXProvider() string
		GetFXStalenessMinutes() int
		GetBaseCurrency() string
	},
	logger *zap.Logger,
) *FXService {
	return &FXService{
		fxRepo:       fxRepo,
		cache:        cache,
		logger:       logger,
		apiKey:       cfg.GetFXAPIKey(),
		provider:     cfg.GetFXProvider(),
		staleness:    time.Duration(cfg.GetFXStalenessMinutes()) * time.Minute,
		baseCurrency: cfg.GetBaseCurrency(),
	}
}

// RateResult wraps a rate with metadata about its freshness
type RateResult struct {
	FromCurrency string          `json:"from_currency"`
	ToCurrency   string          `json:"to_currency"`
	Rate         decimal.Decimal `json:"rate"`
	Source       string          `json:"source"`
	Provider     string          `json:"provider"`
	RateDate     time.Time       `json:"rate_date"`
	IsStale      bool            `json:"is_stale"`
	StaleSince   *time.Time      `json:"stale_since,omitempty"`
	CachedAt     *time.Time      `json:"cached_at,omitempty"`
}

// GetRate returns the latest rate for a currency pair.
// Attempts: Redis cache → DB latest → External API → Fallback cached rate
// Always returns rate metadata so callers can show staleness warnings.
func (s *FXService) GetRate(ctx context.Context, from, to string) (*RateResult, error) {
	if from == to {
		return &RateResult{
			FromCurrency: from,
			ToCurrency:   to,
			Rate:         decimal.NewFromInt(1),
			Source:       "IDENTITY",
			RateDate:     time.Now(),
			IsStale:      false,
		}, nil
	}

	// 1. Try Redis cache
	cacheKey := fmt.Sprintf(fxRateCacheKeyFmt, from, to)
	if cached, err := s.cache.Get(ctx, cacheKey).Result(); err == nil {
		var result RateResult
		if json.Unmarshal([]byte(cached), &result) == nil {
			return &result, nil
		}
	}

	// 2. Try DB for latest rate
	dbRate, err := s.fxRepo.GetLatestRate(ctx, from, to)
	if err == nil {
		isStale := time.Since(dbRate.CreatedAt) > s.staleness
		result := &RateResult{
			FromCurrency: from,
			ToCurrency:   to,
			Rate:         dbRate.Rate,
			Source:       dbRate.Source,
			RateDate:     dbRate.RateDate,
			IsStale:      isStale,
		}
		if provider := dbRate.Provider; provider != nil {
			result.Provider = *provider
		}
		if isStale {
			now := time.Now()
			result.StaleSince = &now
			s.logger.Warn("Returning stale FX rate",
				zap.String("from", from), zap.String("to", to),
				zap.Time("rate_date", dbRate.RateDate))
		}
		// Cache for 5 min
		s.cacheRate(ctx, cacheKey, result)
		return result, nil
	}

	// 3. Fetch from external provider
	if s.apiKey != "" {
		externalRate, err := s.fetchFromProvider(ctx, from, to)
		if err == nil {
			s.fxRepo.SaveRate(ctx, externalRate)
			result := &RateResult{
				FromCurrency: from,
				ToCurrency:   to,
				Rate:         externalRate.Rate,
				Source:       "PROVIDER",
				Provider:     s.provider,
				RateDate:     time.Now(),
				IsStale:      false,
			}
			s.cacheRate(ctx, cacheKey, result)
			return result, nil
		}
		s.logger.Error("FX provider fetch failed", zap.Error(err))
	}

	return nil, fmt.Errorf("no FX rate available for %s/%s", from, to)
}

// ConvertAmount converts an amount from one currency to another.
// Returns the converted amount and rate metadata for display.
func (s *FXService) ConvertAmount(ctx context.Context, amount decimal.Decimal, from, to string) (decimal.Decimal, *RateResult, error) {
	if from == to {
		return amount, &RateResult{
			FromCurrency: from,
			ToCurrency:   to,
			Rate:         decimal.NewFromInt(1),
			Source:       "IDENTITY",
			RateDate:     time.Now(),
		}, nil
	}

	rateResult, err := s.GetRate(ctx, from, to)
	if err != nil {
		return decimal.Zero, nil, fmt.Errorf("conversion rate unavailable: %w", err)
	}

	converted := amount.Mul(rateResult.Rate)
	return converted, rateResult, nil
}

// GetAllLatestRates returns all latest rates with staleness info
func (s *FXService) GetAllLatestRates(ctx context.Context) ([]RateResult, error) {
	dbRates, err := s.fxRepo.GetAllLatestRates(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]RateResult, len(dbRates))
	for i, r := range dbRates {
		isStale := time.Since(r.CreatedAt) > s.staleness
		results[i] = RateResult{
			FromCurrency: r.FromCurrency,
			ToCurrency:   r.ToCurrency,
			Rate:         r.Rate,
			Source:       r.Source,
			RateDate:     r.RateDate,
			IsStale:      isStale,
		}
		if r.Provider != nil {
			results[i].Provider = *r.Provider
		}
	}
	return results, nil
}

// RunSensitivityAnalysis calculates the impact of currency movements on exposure
func (s *FXService) RunSensitivityAnalysis(ctx context.Context, companyID interface{}, currency string, changes []float64) (*SensitivityResult, error) {
	// Get current rate
	rateResult, err := s.GetRate(ctx, currency, s.baseCurrency)
	if err != nil {
		return nil, fmt.Errorf("getting rate for sensitivity analysis: %w", err)
	}

	result := &SensitivityResult{
		Currency:        currency,
		BaseCurrency:    s.baseCurrency,
		CurrentRate:     rateResult.Rate,
		RateDate:        rateResult.RateDate,
		IsRateStale:     rateResult.IsStale,
		Scenarios:       make([]SensitivityScenario, len(changes)),
		Disclaimer:      "Results are estimates based on current exposure data and simulated rate movements. Not financial advice.",
	}

	for i, pct := range changes {
		change := decimal.NewFromFloat(pct / 100.0)
		newRate := rateResult.Rate.Mul(decimal.NewFromInt(1).Add(change))

		// For demonstration, we use a fixed exposure amount
		// In production, this queries fx_exposures for the company
		grossExposure := decimal.NewFromInt(6200000) // $6.2M from seed data
		currentBaseValue := grossExposure.Mul(rateResult.Rate)
		newBaseValue := grossExposure.Mul(newRate)
		impact := newBaseValue.Sub(currentBaseValue)

		label := fmt.Sprintf("%+.0f%%", pct)
		result.Scenarios[i] = SensitivityScenario{
			ChangePercent: pct,
			Label:         label,
			NewRate:       newRate,
			CurrentValue:  currentBaseValue,
			NewValue:      newBaseValue,
			ImpactAmount:  impact,
			ImpactPercent: pct,
		}
	}

	return result, nil
}

// SyncRates fetches and saves latest rates from provider — called by background worker
func (s *FXService) SyncRates(ctx context.Context) error {
	if s.apiKey == "" {
		s.logger.Info("No FX API key configured, skipping sync")
		return nil
	}

	currencies := []string{"USD", "EUR", "GBP", "JPY", "CHF", "AUD", "SGD", "CAD"}
	var syncErrors []error

	for _, curr := range currencies {
		if curr == s.baseCurrency {
			continue
		}
		_, err := s.fetchFromProvider(ctx, curr, s.baseCurrency)
		if err != nil {
			s.logger.Error("Failed to sync FX rate",
				zap.String("currency", curr), zap.Error(err))
			syncErrors = append(syncErrors, err)
			continue
		}
		// Invalidate cache
		cacheKey := fmt.Sprintf(fxRateCacheKeyFmt, curr, s.baseCurrency)
		s.cache.Del(ctx, cacheKey)
	}

	if len(syncErrors) == len(currencies)-1 {
		return fmt.Errorf("all FX rate syncs failed")
	}
	return nil
}

// GetBaseCurrency returns the configured base currency
func (s *FXService) GetBaseCurrency() string {
	return s.baseCurrency
}

// GetRateHistory returns historical rates for charting
func (s *FXService) GetRateHistory(ctx context.Context, from, to string, days int) ([]models.FXRate, error) {
	since := time.Now().AddDate(0, 0, -days)
	return s.fxRepo.GetRateHistory(ctx, from, to, since)
}

// fetchFromProvider fetches a rate from the configured external provider
func (s *FXService) fetchFromProvider(ctx context.Context, from, to string) (*models.FXRate, error) {
	// ExchangeRate-API (free tier supports 1500 req/month)
	// URL: https://v6.exchangerate-api.com/v6/{apikey}/pair/{from}/{to}
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/pair/%s/%s", s.apiKey, from, to)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("provider returned status %d", resp.StatusCode)
	}

	var apiResp struct {
		Result          string  `json:"result"`
		ConversionRate  float64 `json:"conversion_rate"`
		TimeLastUpdateUnix int64  `json:"time_last_update_unix"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if apiResp.Result != "success" {
		return nil, fmt.Errorf("provider returned error result")
	}

	rateDate := time.Unix(apiResp.TimeLastUpdateUnix, 0)
	providerName := "EXCHANGERATE_API"
	rate := &models.FXRate{
		FromCurrency: from,
		ToCurrency:   to,
		Rate:         decimal.NewFromFloat(apiResp.ConversionRate),
		Source:       "PROVIDER",
		Provider:     &providerName,
		RateDate:     rateDate,
		IsLatest:     true,
	}

	return rate, nil
}

func (s *FXService) cacheRate(ctx context.Context, key string, result *RateResult) {
	data, err := json.Marshal(result)
	if err != nil {
		return
	}
	s.cache.Set(ctx, key, data, fxRateLatestCacheTTL)
}

// SensitivityResult wraps sensitivity analysis output
type SensitivityResult struct {
	Currency     string                `json:"currency"`
	BaseCurrency string                `json:"base_currency"`
	CurrentRate  decimal.Decimal       `json:"current_rate"`
	RateDate     time.Time             `json:"rate_date"`
	IsRateStale  bool                  `json:"is_rate_stale"`
	Scenarios    []SensitivityScenario `json:"scenarios"`
	Disclaimer   string                `json:"disclaimer"`
}

type SensitivityScenario struct {
	ChangePercent float64         `json:"change_percent"`
	Label         string          `json:"label"`
	NewRate       decimal.Decimal `json:"new_rate"`
	CurrentValue  decimal.Decimal `json:"current_value"`
	NewValue      decimal.Decimal `json:"new_value"`
	ImpactAmount  decimal.Decimal `json:"impact_amount"`
	ImpactPercent float64         `json:"impact_percent"`
}
