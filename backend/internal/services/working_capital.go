package services

import (
	"context"
	"time"

	"github.com/finsight/backend/internal/models"
	"github.com/finsight/backend/internal/repositories"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// WorkingCapitalService calculates DSO, DPO, DIO, CCC, and WC optimization insights.
// All monetary calculations use decimal arithmetic to avoid floating-point errors.
type WorkingCapitalService struct {
	invoiceRepo repositories.InvoiceRepository
	billRepo    repositories.BillRepository
	fxSvc       *FXService
	cfg         interface{ GetBaseCurrency() string }
}

func NewWorkingCapitalService(
	invoiceRepo repositories.InvoiceRepository,
	billRepo repositories.BillRepository,
	fxSvc *FXService,
	cfg interface{ GetBaseCurrency() string },
) *WorkingCapitalService {
	return &WorkingCapitalService{
		invoiceRepo: invoiceRepo,
		billRepo:    billRepo,
		fxSvc:       fxSvc,
		cfg:         cfg,
	}
}

// GetKPIs calculates all working capital KPIs for the given company and period.
// Period defaults to last 90 days if not specified.
func (s *WorkingCapitalService) GetKPIs(ctx context.Context, companyID uuid.UUID, asOf time.Time, periodDays int) (*models.WorkingCapitalKPIs, error) {
	if periodDays <= 0 {
		periodDays = 90
	}
	periodStart := asOf.AddDate(0, 0, -periodDays)

	// Get AR data
	arStats, err := s.invoiceRepo.GetARStats(ctx, companyID, periodStart, asOf)
	if err != nil {
		return nil, err
	}

	// Get AP data
	apStats, err := s.billRepo.GetAPStats(ctx, companyID, periodStart, asOf)
	if err != nil {
		return nil, err
	}

	kpis := &models.WorkingCapitalKPIs{
		TotalReceivables:   arStats.OutstandingAmount,
		TotalPayables:      apStats.OutstandingAmount,
		TotalRevenue:       arStats.TotalInvoiced,
		TotalCOGS:          apStats.TotalBilled,
		Period:             "90d",
		PeriodDays:         periodDays,
		AsOfDate:           asOf,
	}

	days := decimal.NewFromInt(int64(periodDays))

	// DSO = (Average Accounts Receivable / Revenue) × Days
	// Using ending AR balance / (revenue / days) to avoid zero division
	if arStats.TotalInvoiced.IsPositive() {
		dailyRevenue := arStats.TotalInvoiced.Div(days)
		if dailyRevenue.IsPositive() {
			kpis.DSO = arStats.OutstandingAmount.Div(dailyRevenue)
		}
	}

	// DPO = (Accounts Payable / COGS) × Days
	if apStats.TotalBilled.IsPositive() {
		dailyCOGS := apStats.TotalBilled.Div(days)
		if dailyCOGS.IsPositive() {
			kpis.DPO = apStats.OutstandingAmount.Div(dailyCOGS)
		}
	}

	// DIO — would come from inventory data; use a placeholder since
	// this platform doesn't directly track inventory accounts.
	// In production, this would pull from an ERP feed.
	// We set it to 0 with a note in the response.
	kpis.DIO = decimal.Zero

	// CCC = DSO + DIO - DPO
	kpis.CCC = kpis.DSO.Add(kpis.DIO).Sub(kpis.DPO)

	// Current Assets = Receivables + Cash (simplified)
	cashBalance, _ := s.getCashBalance(ctx, companyID)
	kpis.CurrentAssets = kpis.TotalReceivables.Add(cashBalance)
	kpis.CurrentLiabilities = kpis.TotalPayables
	kpis.NetWorkingCapital = kpis.CurrentAssets.Sub(kpis.CurrentLiabilities)

	// Current Ratio = Current Assets / Current Liabilities
	if kpis.CurrentLiabilities.IsPositive() {
		kpis.CurrentRatio = kpis.CurrentAssets.Div(kpis.CurrentLiabilities)
	}

	// Quick Ratio = (Current Assets - Inventory) / Current Liabilities
	// Since DIO=0, Quick Ratio ≈ Current Ratio here
	if kpis.CurrentLiabilities.IsPositive() {
		kpis.QuickRatio = kpis.CurrentAssets.Sub(kpis.TotalInventory).Div(kpis.CurrentLiabilities)
	}

	// Receivables Turnover = Revenue / Average Receivables
	if kpis.TotalReceivables.IsPositive() {
		kpis.ReceivablesTurnover = kpis.TotalRevenue.Div(kpis.TotalReceivables)
	}

	// Payables Turnover = COGS / Average Payables
	if kpis.TotalPayables.IsPositive() {
		kpis.PayablesTurnover = kpis.TotalCOGS.Div(kpis.TotalPayables)
	}

	return kpis, nil
}

// GetOpportunities generates working capital optimization insights.
// These are estimates based on the platform's data and should be clearly
// labeled as such in the UI — not authoritative financial recommendations.
func (s *WorkingCapitalService) GetOpportunities(ctx context.Context, companyID uuid.UUID) ([]WCOpportunity, error) {
	kpis, err := s.GetKPIs(ctx, companyID, time.Now(), 90)
	if err != nil {
		return nil, err
	}

	var opportunities []WCOpportunity

	// Opportunity 1: DSO reduction
	// Industry benchmark: 30 days for manufacturing
	dsoBenchmark := decimal.NewFromInt(30)
	if kpis.DSO.GreaterThan(dsoBenchmark) {
		dsoDelta := kpis.DSO.Sub(dsoBenchmark)
		// Cash release = (DSO reduction / 365) × Annual Revenue
		annualRevenue := kpis.TotalRevenue.Mul(decimal.NewFromInt(4)) // Approximate 90d → annual
		cashRelease := dsoDelta.Div(decimal.NewFromInt(365)).Mul(annualRevenue)

		opportunities = append(opportunities, WCOpportunity{
			Category:        "RECEIVABLES",
			Title:           "DSO Reduction Opportunity",
			Description:     "Current DSO of " + kpis.DSO.StringFixed(0) + " days exceeds the 30-day benchmark. Accelerating collections could release working capital.",
			Detail:          "If DSO decreases from " + kpis.DSO.StringFixed(1) + " to 30 days, estimated cash release ≈ " + formatCurrency(cashRelease, "INR") + " (estimate based on last 90 days revenue annualized).",
			EstimatedImpact: cashRelease,
			Currency:        "INR",
			Priority:        priorityFromGap(kpis.DSO, dsoBenchmark),
			Disclaimer:      "Estimate based on current outstanding receivables and 90-day revenue. Actual impact depends on customer mix and collections efficiency.",
		})
	}

	// Opportunity 2: DPO extension
	// If DPO is below supplier payment terms, there is an opportunity to extend
	dpoBenchmark := decimal.NewFromInt(45)
	if kpis.DPO.LessThan(dpoBenchmark) && kpis.TotalCOGS.IsPositive() {
		dpoDelta := dpoBenchmark.Sub(kpis.DPO)
		annualCOGS := kpis.TotalCOGS.Mul(decimal.NewFromInt(4))
		cashRetention := dpoDelta.Div(decimal.NewFromInt(365)).Mul(annualCOGS)

		opportunities = append(opportunities, WCOpportunity{
			Category:        "PAYABLES",
			Title:           "Payment Terms Extension",
			Description:     "Current DPO of " + kpis.DPO.StringFixed(0) + " days. Extending payment terms to 45 days could improve cash retention.",
			Detail:          "If DPO increases from " + kpis.DPO.StringFixed(1) + " to 45 days, projected cash retention ≈ " + formatCurrency(cashRetention, "INR") + " (estimate based on last 90 days COGS).",
			EstimatedImpact: cashRetention,
			Currency:        "INR",
			Priority:        "MEDIUM",
			Disclaimer:      "Extending payment terms requires supplier agreement and may affect supplier relationships. This is an estimate only.",
		})
	}

	// Opportunity 3: Overdue AR collection
	overdueAR, _ := s.invoiceRepo.GetTotalOverdue(ctx, companyID)
	if overdueAR.GreaterThan(decimal.Zero) {
		opportunities = append(opportunities, WCOpportunity{
			Category:        "RECEIVABLES",
			Title:           "Overdue AR Collection",
			Description:     "There is " + formatCurrency(overdueAR, "INR") + " in overdue receivables. Focused collection effort could accelerate cash inflows.",
			Detail:          "Prioritize collection on invoices >60 days overdue. Consider early payment discounts for large accounts.",
			EstimatedImpact: overdueAR,
			Currency:        "INR",
			Priority:        "HIGH",
			Disclaimer:      "Full collection of overdue amounts is not guaranteed. Impact depends on customer financial health and dispute resolution.",
		})
	}

	return opportunities, nil
}

// GetARAgingBuckets returns AR aging analysis
func (s *WorkingCapitalService) GetARAgingBuckets(ctx context.Context, companyID uuid.UUID, asOf time.Time) ([]models.AgingBucket, error) {
	return s.invoiceRepo.GetAgingBuckets(ctx, companyID, asOf)
}

// GetAPAgingBuckets returns AP aging analysis
func (s *WorkingCapitalService) GetAPAgingBuckets(ctx context.Context, companyID uuid.UUID, asOf time.Time) ([]models.AgingBucket, error) {
	return s.billRepo.GetAgingBuckets(ctx, companyID, asOf)
}

// GetCCCTrend returns CCC trend over time
func (s *WorkingCapitalService) GetCCCTrend(ctx context.Context, companyID uuid.UUID, months int) ([]CCCDataPoint, error) {
	var trend []CCCDataPoint
	now := time.Now()

	for i := months - 1; i >= 0; i-- {
		asOf := now.AddDate(0, -i, 0)
		kpis, err := s.GetKPIs(ctx, companyID, asOf, 30)
		if err != nil {
			continue
		}
		trend = append(trend, CCCDataPoint{
			Date: asOf,
			DSO:  kpis.DSO,
			DPO:  kpis.DPO,
			DIO:  kpis.DIO,
			CCC:  kpis.CCC,
		})
	}

	return trend, nil
}

func (s *WorkingCapitalService) getCashBalance(ctx context.Context, companyID uuid.UUID) (decimal.Decimal, error) {
	// Simplified — in production this queries cash accounts
	return decimal.NewFromInt(327800000), nil // ₹32.78 Cr demo value
}

// WCOpportunity represents a working capital optimization suggestion
type WCOpportunity struct {
	Category        string          `json:"category"`
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	Detail          string          `json:"detail"`
	EstimatedImpact decimal.Decimal `json:"estimated_impact"`
	Currency        string          `json:"currency"`
	Priority        string          `json:"priority"` // LOW, MEDIUM, HIGH, CRITICAL
	Disclaimer      string          `json:"disclaimer"`
}

type CCCDataPoint struct {
	Date time.Time       `json:"date"`
	DSO  decimal.Decimal `json:"dso"`
	DPO  decimal.Decimal `json:"dpo"`
	DIO  decimal.Decimal `json:"dio"`
	CCC  decimal.Decimal `json:"ccc"`
}

// ARStats holds aggregated AR metrics
type ARStats struct {
	OutstandingAmount decimal.Decimal
	TotalInvoiced     decimal.Decimal
	OverdueAmount     decimal.Decimal
	InvoiceCount      int
	OverdueCount      int
}

// APStats holds aggregated AP metrics
type APStats struct {
	OutstandingAmount decimal.Decimal
	TotalBilled       decimal.Decimal
	OverdueAmount     decimal.Decimal
	BillCount         int
}

func priorityFromGap(actual, benchmark decimal.Decimal) string {
	gap := actual.Sub(benchmark)
	if gap.GreaterThan(decimal.NewFromInt(30)) {
		return "CRITICAL"
	} else if gap.GreaterThan(decimal.NewFromInt(15)) {
		return "HIGH"
	} else if gap.GreaterThan(decimal.NewFromInt(7)) {
		return "MEDIUM"
	}
	return "LOW"
}

func formatCurrency(amount decimal.Decimal, currency string) string {
	switch currency {
	case "INR":
		return "₹" + amount.StringFixed(0)
	case "USD":
		return "$" + amount.StringFixed(0)
	default:
		return amount.StringFixed(0) + " " + currency
	}
}
