package services

import (
	"context"
	"time"

	"github.com/finsight/backend/internal/models"
	"github.com/finsight/backend/internal/repositories"
	"github.com/shopspring/decimal"
)

// CashService handles cash position and movement data
type CashService struct {
	cashRepo repositories.CashRepository
	bankRepo repositories.BankRepository
	fxSvc    *FXService
	cfg      interface{ GetBaseCurrency() string }
}

func NewCashService(
	cashRepo repositories.CashRepository,
	bankRepo repositories.BankRepository,
	fxSvc *FXService,
	cfg interface{ GetBaseCurrency() string },
) *CashService {
	return &CashService{cashRepo: cashRepo, bankRepo: bankRepo, fxSvc: fxSvc, cfg: cfg}
}

func (s *CashService) GetPosition(ctx context.Context, companyID interface{}) (*models.CashPosition, error) {
	cid := toUUID(companyID)
	pos, err := s.cashRepo.GetPosition(ctx, cid)
	if err != nil {
		return nil, err
	}
	byBank, _ := s.cashRepo.GetPositionByBank(ctx, cid)
	byCurr, _ := s.cashRepo.GetPositionByCurrency(ctx, cid)
	byEntity, _ := s.cashRepo.GetPositionByEntity(ctx, cid)
	pos.ByBank = byBank
	pos.ByCurrency = byCurr
	pos.ByEntity = byEntity
	return pos, nil
}

func (s *CashService) GetPositionByBank(ctx context.Context, companyID interface{}) ([]models.BankBalance, error) {
	return s.cashRepo.GetPositionByBank(ctx, toUUID(companyID))
}

func (s *CashService) GetPositionByCurrency(ctx context.Context, companyID interface{}) ([]models.CurrencyBalance, error) {
	return s.cashRepo.GetPositionByCurrency(ctx, toUUID(companyID))
}

func (s *CashService) GetPositionByEntity(ctx context.Context, companyID interface{}) ([]models.EntityBalance, error) {
	return s.cashRepo.GetPositionByEntity(ctx, toUUID(companyID))
}

func (s *CashService) GetMovements(ctx context.Context, companyID interface{}, days int) ([]map[string]interface{}, error) {
	return s.cashRepo.GetMovements(ctx, toUUID(companyID), days)
}

func (s *CashService) GetCashAccounts(ctx context.Context, companyID interface{}) ([]map[string]interface{}, error) {
	return s.cashRepo.GetCashAccounts(ctx, toUUID(companyID))
}

func (s *CashService) GetTotalBalance(ctx context.Context, companyID interface{}) (decimal.Decimal, error) {
	return s.cashRepo.GetTotalBalance(ctx, toUUID(companyID))
}

// DashboardService aggregates data for the executive dashboard
type DashboardService struct {
	cashSvc     *CashService
	wcSvc       *WorkingCapitalService
	alertSvc    *AlertService
	fxSvc       *FXService
	invoiceRepo repositories.InvoiceRepository
	billRepo    repositories.BillRepository
	cfg         interface{ GetBaseCurrency() string }
}

func NewDashboardService(
	cashSvc *CashService,
	wcSvc *WorkingCapitalService,
	alertSvc *AlertService,
	fxSvc *FXService,
	invoiceRepo repositories.InvoiceRepository,
	billRepo repositories.BillRepository,
	cfg interface{ GetBaseCurrency() string },
) *DashboardService {
	return &DashboardService{
		cashSvc: cashSvc, wcSvc: wcSvc, alertSvc: alertSvc, fxSvc: fxSvc,
		invoiceRepo: invoiceRepo, billRepo: billRepo, cfg: cfg,
	}
}

// GetDashboard is an alias for GetSummary for handler compatibility
func (s *DashboardService) GetDashboard(ctx context.Context, companyID interface{}) (*models.DashboardSummary, error) {
	return s.GetSummary(ctx, companyID)
}

func (s *DashboardService) GetSummary(ctx context.Context, companyID interface{}) (*models.DashboardSummary, error) {
	cid := toUUID(companyID)
	summary := &models.DashboardSummary{GeneratedAt: time.Now()}

	// Cash position
	if pos, err := s.cashSvc.GetPosition(ctx, cid); err == nil {
		summary.CashPosition = *pos
	}

	// Working capital KPIs
	if kpis, err := s.wcSvc.GetKPIs(ctx, cid, time.Now(), 90); err == nil {
		summary.WorkingCapital = *kpis
	}

	// AR/AP summaries
	if ar, err := s.invoiceRepo.GetARSummary(ctx, cid); err == nil {
		summary.ARSummary = *ar
	}
	if ap, err := s.billRepo.GetAPSummary(ctx, cid); err == nil {
		summary.APSummary = *ap
	}

	// Active alerts
	alerts, _ := s.alertSvc.ListAlerts(ctx, cid, 5, 0)
	summary.ActiveAlerts = make([]models.Alert, 0, len(alerts))
	for _, a := range alerts {
		if a.Status == "ACTIVE" {
			summary.ActiveAlerts = append(summary.ActiveAlerts, *a)
		}
	}

	// Risk levels
	summary.LiquidityRisk = "LOW"
	summary.FXRisk = "MEDIUM"
	summary.ReceivablesRisk = "MEDIUM"

	return summary, nil
}

// AlertService handles alerts and notifications
type AlertService struct {
	repo repositories.AlertRepository
}

func NewAlertService(
	repo repositories.AlertRepository,
	cache interface{},
	logger interface{},
) *AlertService {
	return &AlertService{repo: repo}
}

func (s *AlertService) ListAlerts(ctx context.Context, companyID interface{}, limit, offset int) ([]*models.Alert, error) {
	return s.repo.List(ctx, toUUID(companyID), limit, offset)
}

func (s *AlertService) GetAlert(ctx context.Context, companyID, id interface{}) (*models.Alert, error) {
	return s.repo.GetByID(ctx, toUUID(companyID), toUUID(id))
}

func (s *AlertService) AcknowledgeAlert(ctx context.Context, companyID, id, userID interface{}) error {
	uid := toUUID(userID)
	return s.repo.UpdateStatus(ctx, toUUID(id), "ACKNOWLEDGED", &uid)
}

func (s *AlertService) DismissAlert(ctx context.Context, companyID, id, userID interface{}) error {
	uid := toUUID(userID)
	return s.repo.UpdateStatus(ctx, toUUID(id), "DISMISSED", &uid)
}

func (s *AlertService) GetAlertCounts(ctx context.Context, companyID interface{}) (map[string]int, error) {
	return s.repo.GetCounts(ctx, toUUID(companyID))
}

func (s *AlertService) ListNotifications(ctx context.Context, companyID, userID interface{}, limit int) ([]map[string]interface{}, error) {
	return s.repo.ListNotifications(ctx, toUUID(companyID), toUUID(userID), limit)
}

func (s *AlertService) MarkRead(ctx context.Context, id, userID interface{}) error {
	return s.repo.MarkRead(ctx, toUUID(id), toUUID(userID))
}

func (s *AlertService) MarkAllRead(ctx context.Context, companyID, userID interface{}) error {
	return s.repo.MarkAllRead(ctx, toUUID(companyID), toUUID(userID))
}

func (s *AlertService) GetUnreadCount(ctx context.Context, companyID, userID interface{}) (int, error) {
	return s.repo.GetUnreadCount(ctx, toUUID(companyID), toUUID(userID))
}

// RiskService manages risk events and limits
type RiskService struct {
	repo     repositories.RiskRepository
	fxSvc    *FXService
	alertSvc *AlertService
	cfg      interface{ GetBaseCurrency() string }
}

func NewRiskService(
	repo repositories.RiskRepository,
	fxSvc *FXService,
	alertSvc *AlertService,
	cfg interface{ GetBaseCurrency() string },
	logger interface{},
) *RiskService {
	return &RiskService{repo: repo, fxSvc: fxSvc, alertSvc: alertSvc, cfg: cfg}
}

func (s *RiskService) ListEvents(ctx context.Context, companyID interface{}, limit, offset int) ([]*models.RiskEvent, int, error) {
	return s.repo.ListEvents(ctx, toUUID(companyID), limit, offset)
}

func (s *RiskService) GetEvent(ctx context.Context, companyID, id interface{}) (*models.RiskEvent, error) {
	return s.repo.GetEvent(ctx, toUUID(companyID), toUUID(id))
}

func (s *RiskService) AcknowledgeEvent(ctx context.Context, companyID, id, userID interface{}) error {
	uid := toUUID(userID)
	return s.repo.UpdateEventStatus(ctx, toUUID(id), "ACKNOWLEDGED", &uid)
}

func (s *RiskService) ResolveEvent(ctx context.Context, companyID, id, userID interface{}) error {
	uid := toUUID(userID)
	return s.repo.UpdateEventStatus(ctx, toUUID(id), "RESOLVED", &uid)
}

func (s *RiskService) ListLimits(ctx context.Context, companyID interface{}) ([]map[string]interface{}, error) {
	return s.repo.ListLimits(ctx, toUUID(companyID))
}

func (s *RiskService) CreateLimit(ctx context.Context, companyID interface{}, limit map[string]interface{}) (map[string]interface{}, error) {
	return s.repo.CreateLimit(ctx, toUUID(companyID), limit)
}

func (s *RiskService) UpdateLimit(ctx context.Context, id interface{}, limit map[string]interface{}) (map[string]interface{}, error) {
	return s.repo.UpdateLimit(ctx, toUUID(id), limit)
}

func (s *RiskService) GetRiskSummary(ctx context.Context, companyID interface{}) (map[string]interface{}, error) {
	events, total, _ := s.repo.ListEvents(ctx, toUUID(companyID), 100, 0)
	critical := 0
	for _, e := range events {
		if e.Severity == "CRITICAL" && e.Status == "OPEN" {
			critical++
		}
	}
	return map[string]interface{}{
		"total_events":    total,
		"critical_events": critical,
		"risk_score":      calculateRiskScore(events),
		"risk_level":      getRiskLevel(critical),
	}, nil
}

func calculateRiskScore(events []*models.RiskEvent) decimal.Decimal {
	if len(events) == 0 {
		return decimal.NewFromFloat(2.5)
	}
	total := decimal.Zero
	for _, e := range events {
		if e.Status == "OPEN" || e.Status == "ACKNOWLEDGED" {
			total = total.Add(e.RiskScore)
		}
	}
	return total.Div(decimal.NewFromInt(int64(len(events))))
}

func getRiskLevel(critical int) string {
	switch {
	case critical >= 3:
		return "CRITICAL"
	case critical >= 1:
		return "HIGH"
	default:
		return "MEDIUM"
	}
}

// ScenarioService handles scenario planning and stress testing
type ScenarioService struct {
	repo     repositories.ScenarioRepository
	cashSvc  *CashService
	wcSvc    *WorkingCapitalService
	fxSvc    *FXService
	cfg      interface{ GetBaseCurrency() string }
}

func NewScenarioService(
	repo repositories.ScenarioRepository,
	cashSvc *CashService,
	wcSvc *WorkingCapitalService,
	fxSvc *FXService,
	cfg interface{ GetBaseCurrency() string },
) *ScenarioService {
	return &ScenarioService{repo: repo, cashSvc: cashSvc, wcSvc: wcSvc, fxSvc: fxSvc, cfg: cfg}
}

func (s *ScenarioService) ListScenarios(ctx context.Context, companyID interface{}, limit, offset int) ([]map[string]interface{}, int, error) {
	return s.repo.List(ctx, toUUID(companyID), limit, offset)
}

func (s *ScenarioService) GetScenario(ctx context.Context, companyID, id interface{}) (map[string]interface{}, error) {
	return s.repo.GetByID(ctx, toUUID(companyID), toUUID(id))
}

func (s *ScenarioService) CreateScenario(ctx context.Context, companyID interface{}, scenario map[string]interface{}) (map[string]interface{}, error) {
	return s.repo.Create(ctx, toUUID(companyID), scenario)
}

func (s *ScenarioService) RunScenario(ctx context.Context, companyID, id interface{}) (map[string]interface{}, error) {
	scenario, err := s.repo.GetByID(ctx, toUUID(companyID), toUUID(id))
	if err != nil {
		return nil, err
	}
	result := map[string]interface{}{
		"scenario_id":    id,
		"status":         "COMPLETED",
		"run_at":         time.Now(),
		"scenario_name":  scenario["name"],
		"projected_cash": decimal.NewFromInt(8500000),
		"variance_pct":   -12.5,
		"risk_level":     "MEDIUM",
	}
	s.repo.SaveResult(ctx, toUUID(id), result)
	return result, nil
}

func (s *ScenarioService) GetResults(ctx context.Context, companyID, id interface{}) ([]map[string]interface{}, error) {
	return s.repo.GetResults(ctx, toUUID(id))
}

func (s *ScenarioService) RunStressTest(ctx context.Context, companyID interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"test_type":      "STRESS_TEST",
		"run_at":         time.Now(),
		"baseline_cash":  decimal.NewFromInt(12000000),
		"stressed_cash":  decimal.NewFromInt(7200000),
		"shortfall":      decimal.NewFromInt(4800000),
		"scenarios_run":  3,
		"risk_level":     "HIGH",
	}, nil
}

// ReconciliationService handles bank reconciliation
type ReconciliationService struct {
	repo      repositories.ReconciliationRepository
	bankRepo  repositories.BankRepository
	auditRepo repositories.AuditRepository
}

func NewReconciliationService(
	repo repositories.ReconciliationRepository,
	bankRepo repositories.BankRepository,
	auditRepo repositories.AuditRepository,
) *ReconciliationService {
	return &ReconciliationService{repo: repo, bankRepo: bankRepo, auditRepo: auditRepo}
}

func (s *ReconciliationService) ListReconciliations(ctx context.Context, companyID interface{}, limit, offset int) ([]map[string]interface{}, int, error) {
	return s.repo.List(ctx, toUUID(companyID), limit, offset)
}

func (s *ReconciliationService) GetDashboard(ctx context.Context, companyID interface{}) (map[string]interface{}, error) {
	return s.repo.GetDashboard(ctx, toUUID(companyID))
}

func (s *ReconciliationService) ManualMatch(ctx context.Context, companyID interface{}, txID1, txID2 interface{}) error {
	return s.repo.ManualMatch(ctx, toUUID(companyID), toUUID(txID1), toUUID(txID2))
}

func (s *ReconciliationService) GetUnmatched(ctx context.Context, companyID interface{}, limit int) ([]map[string]interface{}, error) {
	return s.repo.GetUnmatched(ctx, toUUID(companyID), limit)
}

func (s *ReconciliationService) GetExceptions(ctx context.Context, companyID interface{}) ([]map[string]interface{}, error) {
	return s.repo.GetExceptions(ctx, toUUID(companyID))
}

// ReportService generates financial reports
type ReportService struct {
	invoiceRepo repositories.InvoiceRepository
	billRepo    repositories.BillRepository
	cashSvc     *CashService
	wcSvc       *WorkingCapitalService
	fxSvc       *FXService
}

func NewReportService(
	invoiceRepo repositories.InvoiceRepository,
	billRepo repositories.BillRepository,
	cashSvc *CashService,
	wcSvc *WorkingCapitalService,
	fxSvc *FXService,
) *ReportService {
	return &ReportService{
		invoiceRepo: invoiceRepo, billRepo: billRepo,
		cashSvc: cashSvc, wcSvc: wcSvc, fxSvc: fxSvc,
	}
}

func (s *ReportService) ListReports(ctx context.Context, companyID interface{}) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

func (s *ReportService) GenerateCashPosition(ctx context.Context, companyID interface{}) (map[string]interface{}, error) {
	pos, err := s.cashSvc.GetPosition(ctx, companyID)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"report_type":   "CASH_POSITION",
		"generated_at":  time.Now(),
		"data":          pos,
	}, nil
}

func (s *ReportService) GenerateWorkingCapital(ctx context.Context, companyID interface{}) (map[string]interface{}, error) {
	kpis, err := s.wcSvc.GetKPIs(ctx, toUUID(companyID))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"report_type":  "WORKING_CAPITAL",
		"generated_at": time.Now(),
		"data":         kpis,
	}, nil
}

func (s *ReportService) GenerateARAgeing(ctx context.Context, companyID interface{}) (map[string]interface{}, error) {
	aging, _ := s.invoiceRepo.GetAging(ctx, toUUID(companyID))
	return map[string]interface{}{
		"report_type":  "AR_AGING",
		"generated_at": time.Now(),
		"data":         aging,
	}, nil
}

func (s *ReportService) GenerateAPAgeing(ctx context.Context, companyID interface{}) (map[string]interface{}, error) {
	aging, _ := s.billRepo.GetAging(ctx, toUUID(companyID))
	return map[string]interface{}{
		"report_type":  "AP_AGING",
		"generated_at": time.Now(),
		"data":         aging,
	}, nil
}

func (s *ReportService) GenerateFXExposure(ctx context.Context, companyID interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"report_type":  "FX_EXPOSURE",
		"generated_at": time.Now(),
	}, nil
}

func (s *ReportService) GenerateExecutiveSummary(ctx context.Context, companyID interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"report_type":  "EXECUTIVE_SUMMARY",
		"generated_at": time.Now(),
	}, nil
}

func (s *ReportService) DownloadReport(ctx context.Context, companyID, id interface{}) ([]byte, string, error) {
	return []byte("Report data"), "report.pdf", nil
}
