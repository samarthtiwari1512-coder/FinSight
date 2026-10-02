package workers

import (
	"context"
	"time"

	"github.com/finsight/backend/internal/repositories"
	"github.com/finsight/backend/internal/services"
	"go.uber.org/zap"
)

// configProvider is the subset of config the workers need
type configProvider interface {
	GetFXUpdateIntervalMinutes() int
}

// StartFXRateSyncer syncs FX rates on a schedule
func StartFXRateSyncer(ctx context.Context, fxSvc *services.FXService, cfg configProvider, logger *zap.Logger) {
	interval := time.Duration(cfg.GetFXUpdateIntervalMinutes()) * time.Minute
	if interval <= 0 {
		interval = 60 * time.Minute
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	logger.Info("FX rate syncer started", zap.Duration("interval", interval))
	// Sync immediately on startup
	if err := fxSvc.SyncRates(ctx); err != nil {
		logger.Warn("Initial FX rate sync failed", zap.Error(err))
	}

	for {
		select {
		case <-ctx.Done():
			logger.Info("FX rate syncer stopped")
			return
		case <-ticker.C:
			if err := fxSvc.SyncRates(ctx); err != nil {
				logger.Error("FX rate sync failed", zap.Error(err))
			} else {
				logger.Info("FX rates synced successfully")
			}
		}
	}
}

// StartAlertEvaluator evaluates alerts periodically
func StartAlertEvaluator(
	ctx context.Context,
	alertSvc *services.AlertService,
	invoiceRepo repositories.InvoiceRepository,
	billRepo repositories.BillRepository,
	cashSvc *services.CashService,
	riskSvc *services.RiskService,
	cfg configProvider,
	logger *zap.Logger,
) {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	logger.Info("Alert evaluator started")
	for {
		select {
		case <-ctx.Done():
			logger.Info("Alert evaluator stopped")
			return
		case <-ticker.C:
			logger.Debug("Running alert evaluation cycle")
			// Alert evaluation logic would run here in a full implementation
			// For now we just log that the cycle ran
		}
	}
}

// StartForecastGenerator regenerates cash flow forecasts periodically
func StartForecastGenerator(ctx context.Context, forecastSvc *services.ForecastService, logger *zap.Logger) {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	logger.Info("Forecast generator started")
	for {
		select {
		case <-ctx.Done():
			logger.Info("Forecast generator stopped")
			return
		case <-ticker.C:
			logger.Debug("Running forecast generation cycle")
			// Forecast generation would run here for each company
		}
	}
}
