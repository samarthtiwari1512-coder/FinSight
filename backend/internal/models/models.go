package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ============================================================
// COMPANY & ENTITIES
// ============================================================

type Company struct {
	ID                UUID      `db:"id" json:"id"`
	Name              string    `db:"name" json:"name"`
	Code              string    `db:"code" json:"code"`
	BaseCurrency      string    `db:"base_currency" json:"base_currency"`
	Country           string    `db:"country" json:"country"`
	FiscalYearStart   int       `db:"fiscal_year_start" json:"fiscal_year_start"`
	Timezone          string    `db:"timezone" json:"timezone"`
	Address           *string   `db:"address" json:"address,omitempty"`
	IsActive          bool      `db:"is_active" json:"is_active"`
	Settings          JSONB     `db:"settings" json:"settings"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

type Entity struct {
	ID             UUID      `db:"id" json:"id"`
	CompanyID      UUID      `db:"company_id" json:"company_id"`
	Name           string    `db:"name" json:"name"`
	Code           string    `db:"code" json:"code"`
	Country        string    `db:"country" json:"country"`
	Currency       string    `db:"currency" json:"currency"`
	TaxID          *string   `db:"tax_id" json:"tax_id,omitempty"`
	EntityType     string    `db:"entity_type" json:"entity_type"`
	IsActive       bool      `db:"is_active" json:"is_active"`
	ParentEntityID *UUID     `db:"parent_entity_id" json:"parent_entity_id,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// ============================================================
// USERS & RBAC
// ============================================================

type User struct {
	ID             UUID       `db:"id" json:"id"`
	CompanyID      UUID       `db:"company_id" json:"company_id"`
	EntityID       *UUID      `db:"entity_id" json:"entity_id,omitempty"`
	Email          string     `db:"email" json:"email"`
	PasswordHash   string     `db:"password_hash" json:"-"`
	FirstName      string     `db:"first_name" json:"first_name"`
	LastName       string     `db:"last_name" json:"last_name"`
	Phone          *string    `db:"phone" json:"phone,omitempty"`
	AvatarURL      *string    `db:"avatar_url" json:"avatar_url,omitempty"`
	IsActive       bool       `db:"is_active" json:"is_active"`
	IsVerified     bool       `db:"is_verified" json:"is_verified"`
	LastLoginAt    *time.Time `db:"last_login_at" json:"last_login_at,omitempty"`
	FailedAttempts int        `db:"failed_attempts" json:"-"`
	LockedUntil    *time.Time `db:"locked_until" json:"-"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`

	// Populated via JOIN
	Roles []Role `db:"-" json:"roles,omitempty"`
}

func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

func (u *User) HasPermission(resource, action string) bool {
	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			if perm.Resource == resource && perm.Action == action {
				return true
			}
		}
	}
	return false
}

func (u *User) HasRole(code string) bool {
	for _, r := range u.Roles {
		if r.Code == code {
			return true
		}
	}
	return false
}

type Role struct {
	ID          UUID         `db:"id" json:"id"`
	CompanyID   UUID         `db:"company_id" json:"company_id"`
	Name        string       `db:"name" json:"name"`
	Code        string       `db:"code" json:"code"`
	Description *string      `db:"description" json:"description,omitempty"`
	IsSystem    bool         `db:"is_system" json:"is_system"`
	CreatedAt   time.Time    `db:"created_at" json:"created_at"`
	Permissions []Permission `db:"-" json:"permissions,omitempty"`
}

type Permission struct {
	ID          UUID      `db:"id" json:"id"`
	Resource    string    `db:"resource" json:"resource"`
	Action      string    `db:"action" json:"action"`
	Description *string   `db:"description" json:"description,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

// ============================================================
// CURRENCY & FX
// ============================================================

type Currency struct {
	Code          string    `db:"code" json:"code"`
	Name          string    `db:"name" json:"name"`
	Symbol        string    `db:"symbol" json:"symbol"`
	DecimalPlaces int16     `db:"decimal_places" json:"decimal_places"`
	IsActive      bool      `db:"is_active" json:"is_active"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

type FXRate struct {
	ID           UUID            `db:"id" json:"id"`
	FromCurrency string          `db:"from_currency" json:"from_currency"`
	ToCurrency   string          `db:"to_currency" json:"to_currency"`
	Rate         decimal.Decimal `db:"rate" json:"rate"`
	BidRate      *decimal.Decimal `db:"bid_rate" json:"bid_rate,omitempty"`
	AskRate      *decimal.Decimal `db:"ask_rate" json:"ask_rate,omitempty"`
	Source       string          `db:"source" json:"source"`
	Provider     *string         `db:"provider" json:"provider,omitempty"`
	RateDate     time.Time       `db:"rate_date" json:"rate_date"`
	IsLatest     bool            `db:"is_latest" json:"is_latest"`
	CreatedAt    time.Time       `db:"created_at" json:"created_at"`
}

// ============================================================
// BANK ACCOUNTS
// ============================================================

type BankAccount struct {
	ID                  UUID            `db:"id" json:"id"`
	CompanyID           UUID            `db:"company_id" json:"company_id"`
	EntityID            UUID            `db:"entity_id" json:"entity_id"`
	BankName            string          `db:"bank_name" json:"bank_name"`
	AccountNumberMasked string          `db:"account_number_masked" json:"account_number_masked"`
	AccountType         string          `db:"account_type" json:"account_type"`
	Currency            string          `db:"currency" json:"currency"`
	Country             string          `db:"country" json:"country"`
	BankCode            *string         `db:"bank_code" json:"bank_code,omitempty"`
	Branch              *string         `db:"branch" json:"branch,omitempty"`
	AccountHolder       string          `db:"account_holder" json:"account_holder"`
	Status              string          `db:"status" json:"status"`
	AvailableBalance    decimal.Decimal `db:"available_balance" json:"available_balance"`
	LedgerBalance       decimal.Decimal `db:"ledger_balance" json:"ledger_balance"`
	RestrictedAmount    decimal.Decimal `db:"restricted_amount" json:"restricted_amount"`
	CreditLimit         *decimal.Decimal `db:"credit_limit" json:"credit_limit,omitempty"`
	LastSyncAt          *time.Time      `db:"last_sync_at" json:"last_sync_at,omitempty"`
	IsPrimary           bool            `db:"is_primary" json:"is_primary"`
	Notes               *string         `db:"notes" json:"notes,omitempty"`
	CreatedAt           time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time       `db:"updated_at" json:"updated_at"`
}

// ============================================================
// CUSTOMERS
// ============================================================

type Customer struct {
	ID           UUID            `db:"id" json:"id"`
	CompanyID    UUID            `db:"company_id" json:"company_id"`
	CustomerCode string          `db:"customer_code" json:"customer_code"`
	Name         string          `db:"name" json:"name"`
	Country      string          `db:"country" json:"country"`
	Currency     string          `db:"currency" json:"currency"`
	Email        *string         `db:"email" json:"email,omitempty"`
	Phone        *string         `db:"phone" json:"phone,omitempty"`
	Address      *string         `db:"address" json:"address,omitempty"`
	PaymentTerms int             `db:"payment_terms" json:"payment_terms"`
	CreditLimit  decimal.Decimal `db:"credit_limit" json:"credit_limit"`
	RiskLevel    string          `db:"risk_level" json:"risk_level"`
	IsActive     bool            `db:"is_active" json:"is_active"`
	Notes        *string         `db:"notes" json:"notes,omitempty"`
	CreatedAt    time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time       `db:"updated_at" json:"updated_at"`

	// Calculated
	OutstandingAmount *decimal.Decimal `db:"outstanding_amount" json:"outstanding_amount,omitempty"`
	OverdueAmount     *decimal.Decimal `db:"overdue_amount" json:"overdue_amount,omitempty"`
	AvgPaymentDays    *float64         `db:"avg_payment_days" json:"avg_payment_days,omitempty"`
	InvoiceCount      *int             `db:"invoice_count" json:"invoice_count,omitempty"`
}

// ============================================================
// SUPPLIERS
// ============================================================

type Supplier struct {
	ID           UUID            `db:"id" json:"id"`
	CompanyID    UUID            `db:"company_id" json:"company_id"`
	SupplierCode string          `db:"supplier_code" json:"supplier_code"`
	Name         string          `db:"name" json:"name"`
	Country      string          `db:"country" json:"country"`
	Currency     string          `db:"currency" json:"currency"`
	Email        *string         `db:"email" json:"email,omitempty"`
	Phone        *string         `db:"phone" json:"phone,omitempty"`
	Address      *string         `db:"address" json:"address,omitempty"`
	PaymentTerms int             `db:"payment_terms" json:"payment_terms"`
	RiskLevel    string          `db:"risk_level" json:"risk_level"`
	IsActive     bool            `db:"is_active" json:"is_active"`
	Notes        *string         `db:"notes" json:"notes,omitempty"`
	CreatedAt    time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time       `db:"updated_at" json:"updated_at"`

	// Calculated
	OutstandingPayable *decimal.Decimal `db:"outstanding_payable" json:"outstanding_payable,omitempty"`
	OverdueAmount      *decimal.Decimal `db:"overdue_amount" json:"overdue_amount,omitempty"`
	BillCount          *int             `db:"bill_count" json:"bill_count,omitempty"`
}

// ============================================================
// INVOICES
// ============================================================

type Invoice struct {
	ID                UUID             `db:"id" json:"id"`
	CompanyID         UUID             `db:"company_id" json:"company_id"`
	EntityID          UUID             `db:"entity_id" json:"entity_id"`
	CustomerID        UUID             `db:"customer_id" json:"customer_id"`
	InvoiceNumber     string           `db:"invoice_number" json:"invoice_number"`
	InvoiceDate       time.Time        `db:"invoice_date" json:"invoice_date"`
	DueDate           time.Time        `db:"due_date" json:"due_date"`
	Currency          string           `db:"currency" json:"currency"`
	Subtotal          decimal.Decimal  `db:"subtotal" json:"subtotal"`
	TaxAmount         decimal.Decimal  `db:"tax_amount" json:"tax_amount"`
	TotalAmount       decimal.Decimal  `db:"total_amount" json:"total_amount"`
	PaidAmount        decimal.Decimal  `db:"paid_amount" json:"paid_amount"`
	OutstandingAmount decimal.Decimal  `db:"outstanding_amount" json:"outstanding_amount"`
	BaseCurrency      string           `db:"base_currency" json:"base_currency"`
	BaseAmount        *decimal.Decimal `db:"base_amount" json:"base_amount,omitempty"`
	ExchangeRate      *decimal.Decimal `db:"exchange_rate" json:"exchange_rate,omitempty"`
	RateDate          *time.Time       `db:"rate_date" json:"rate_date,omitempty"`
	Status            string           `db:"status" json:"status"`
	Reference         *string          `db:"reference" json:"reference,omitempty"`
	PurchaseOrder     *string          `db:"purchase_order" json:"purchase_order,omitempty"`
	Notes             *string          `db:"notes" json:"notes,omitempty"`
	CreatedBy         *UUID            `db:"created_by" json:"created_by,omitempty"`
	ApprovedBy        *UUID            `db:"approved_by" json:"approved_by,omitempty"`
	ApprovedAt        *time.Time       `db:"approved_at" json:"approved_at,omitempty"`
	CreatedAt         time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time        `db:"updated_at" json:"updated_at"`

	// Joined
	CustomerName *string `db:"customer_name" json:"customer_name,omitempty"`
	EntityName   *string `db:"entity_name" json:"entity_name,omitempty"`
	DaysOverdue  *int    `db:"days_overdue" json:"days_overdue,omitempty"`
	AgeBucket    *string `db:"age_bucket" json:"age_bucket,omitempty"`
}

// ============================================================
// BILLS (AP)
// ============================================================

type Bill struct {
	ID                  UUID             `db:"id" json:"id"`
	CompanyID           UUID             `db:"company_id" json:"company_id"`
	EntityID            UUID             `db:"entity_id" json:"entity_id"`
	SupplierID          UUID             `db:"supplier_id" json:"supplier_id"`
	BillNumber          string           `db:"bill_number" json:"bill_number"`
	SupplierInvoiceRef  *string          `db:"supplier_invoice_ref" json:"supplier_invoice_ref,omitempty"`
	InvoiceDate         time.Time        `db:"invoice_date" json:"invoice_date"`
	DueDate             time.Time        `db:"due_date" json:"due_date"`
	Currency            string           `db:"currency" json:"currency"`
	Subtotal            decimal.Decimal  `db:"subtotal" json:"subtotal"`
	TaxAmount           decimal.Decimal  `db:"tax_amount" json:"tax_amount"`
	TotalAmount         decimal.Decimal  `db:"total_amount" json:"total_amount"`
	PaidAmount          decimal.Decimal  `db:"paid_amount" json:"paid_amount"`
	OutstandingAmount   decimal.Decimal  `db:"outstanding_amount" json:"outstanding_amount"`
	BaseCurrency        string           `db:"base_currency" json:"base_currency"`
	BaseAmount          *decimal.Decimal `db:"base_amount" json:"base_amount,omitempty"`
	ExchangeRate        *decimal.Decimal `db:"exchange_rate" json:"exchange_rate,omitempty"`
	Status              string           `db:"status" json:"status"`
	Reference           *string          `db:"reference" json:"reference,omitempty"`
	Notes               *string          `db:"notes" json:"notes,omitempty"`
	CreatedBy           *UUID            `db:"created_by" json:"created_by,omitempty"`
	ApprovedBy          *UUID            `db:"approved_by" json:"approved_by,omitempty"`
	ApprovedAt          *time.Time       `db:"approved_at" json:"approved_at,omitempty"`
	CreatedAt           time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time        `db:"updated_at" json:"updated_at"`

	// Joined
	SupplierName *string `db:"supplier_name" json:"supplier_name,omitempty"`
	EntityName   *string `db:"entity_name" json:"entity_name,omitempty"`
	DaysOverdue  *int    `db:"days_overdue" json:"days_overdue,omitempty"`
}

// ============================================================
// PAYMENTS
// ============================================================

type Payment struct {
	ID              UUID             `db:"id" json:"id"`
	CompanyID       UUID             `db:"company_id" json:"company_id"`
	EntityID        UUID             `db:"entity_id" json:"entity_id"`
	PaymentType     string           `db:"payment_type" json:"payment_type"`
	Category        string           `db:"category" json:"category"`
	ReferenceType   *string          `db:"reference_type" json:"reference_type,omitempty"`
	ReferenceID     *UUID            `db:"reference_id" json:"reference_id,omitempty"`
	CounterpartyName string          `db:"counterparty_name" json:"counterparty_name"`
	Amount          decimal.Decimal  `db:"amount" json:"amount"`
	Currency        string           `db:"currency" json:"currency"`
	BaseAmount      *decimal.Decimal `db:"base_amount" json:"base_amount,omitempty"`
	ExchangeRate    *decimal.Decimal `db:"exchange_rate" json:"exchange_rate,omitempty"`
	BankAccountID   *UUID            `db:"bank_account_id" json:"bank_account_id,omitempty"`
	PaymentDate     *time.Time       `db:"payment_date" json:"payment_date,omitempty"`
	ValueDate       *time.Time       `db:"value_date" json:"value_date,omitempty"`
	Status          string           `db:"status" json:"status"`
	BatchID         *UUID            `db:"batch_id" json:"batch_id,omitempty"`
	Reference       *string          `db:"reference" json:"reference,omitempty"`
	Notes           *string          `db:"notes" json:"notes,omitempty"`
	CreatedBy       UUID             `db:"created_by" json:"created_by"`
	CreatedAt       time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time        `db:"updated_at" json:"updated_at"`

	// Joined
	CreatedByName  *string          `db:"created_by_name" json:"created_by_name,omitempty"`
	Approvals      []PaymentApproval `db:"-" json:"approvals,omitempty"`
}

type PaymentApproval struct {
	ID         UUID      `db:"id" json:"id"`
	PaymentID  UUID      `db:"payment_id" json:"payment_id"`
	ApproverID UUID      `db:"approver_id" json:"approver_id"`
	Level      int16     `db:"level" json:"level"`
	Decision   string    `db:"decision" json:"decision"`
	Reason     *string   `db:"reason" json:"reason,omitempty"`
	DecidedAt  time.Time `db:"decided_at" json:"decided_at"`

	// Joined
	ApproverName *string `db:"approver_name" json:"approver_name,omitempty"`
}

// ============================================================
// FX EXPOSURE
// ============================================================

type FXExposure struct {
	ID           UUID            `db:"id" json:"id"`
	CompanyID    UUID            `db:"company_id" json:"company_id"`
	EntityID     UUID            `db:"entity_id" json:"entity_id"`
	ExposureDate time.Time       `db:"exposure_date" json:"exposure_date"`
	Currency     string          `db:"currency" json:"currency"`
	ExposureType string          `db:"exposure_type" json:"exposure_type"`
	GrossExposure decimal.Decimal `db:"gross_exposure" json:"gross_exposure"`
	HedgedAmount decimal.Decimal `db:"hedged_amount" json:"hedged_amount"`
	NetExposure  decimal.Decimal `db:"net_exposure" json:"net_exposure"`
	BaseCurrency string          `db:"base_currency" json:"base_currency"`
	BaseValue    *decimal.Decimal `db:"base_value" json:"base_value,omitempty"`
	ExchangeRate *decimal.Decimal `db:"exchange_rate" json:"exchange_rate,omitempty"`
	RiskLevel    string          `db:"risk_level" json:"risk_level"`
	CreatedAt    time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time       `db:"updated_at" json:"updated_at"`
}

// ============================================================
// RISK EVENTS
// ============================================================

type RiskEvent struct {
	ID             UUID             `db:"id" json:"id"`
	CompanyID      UUID             `db:"company_id" json:"company_id"`
	RiskLimitID    *UUID            `db:"risk_limit_id" json:"risk_limit_id,omitempty"`
	EventType      string           `db:"event_type" json:"event_type"`
	Severity       string           `db:"severity" json:"severity"`
	RiskScore      decimal.Decimal  `db:"risk_score" json:"risk_score"`
	Title          string           `db:"title" json:"title"`
	Description    string           `db:"description" json:"description"`
	AffectedEntity *string          `db:"affected_entity" json:"affected_entity,omitempty"`
	AffectedAmount *decimal.Decimal `db:"affected_amount" json:"affected_amount,omitempty"`
	Currency       *string          `db:"currency" json:"currency,omitempty"`
	Status         string           `db:"status" json:"status"`
	OwnerID        *UUID            `db:"owner_id" json:"owner_id,omitempty"`
	ResolvedAt     *time.Time       `db:"resolved_at" json:"resolved_at,omitempty"`
	CreatedAt      time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time        `db:"updated_at" json:"updated_at"`
}

// ============================================================
// ALERTS
// ============================================================

type Alert struct {
	ID                UUID             `db:"id" json:"id"`
	CompanyID         UUID             `db:"company_id" json:"company_id"`
	AlertType         string           `db:"alert_type" json:"alert_type"`
	Severity          string           `db:"severity" json:"severity"`
	Title             string           `db:"title" json:"title"`
	Description       string           `db:"description" json:"description"`
	Source            string           `db:"source" json:"source"`
	EntityID          *UUID            `db:"entity_id" json:"entity_id,omitempty"`
	AffectedAmount    *decimal.Decimal `db:"affected_amount" json:"affected_amount,omitempty"`
	Currency          *string          `db:"currency" json:"currency,omitempty"`
	ReferenceType     *string          `db:"reference_type" json:"reference_type,omitempty"`
	ReferenceID       *UUID            `db:"reference_id" json:"reference_id,omitempty"`
	RecommendedAction *string          `db:"recommended_action" json:"recommended_action,omitempty"`
	Status            string           `db:"status" json:"status"`
	AcknowledgedBy    *UUID            `db:"acknowledged_by" json:"acknowledged_by,omitempty"`
	AcknowledgedAt    *time.Time       `db:"acknowledged_at" json:"acknowledged_at,omitempty"`
	ResolvedAt        *time.Time       `db:"resolved_at" json:"resolved_at,omitempty"`
	CreatedAt         time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time        `db:"updated_at" json:"updated_at"`
}

// ============================================================
// AUDIT LOG
// ============================================================

type AuditLog struct {
	ID         UUID      `db:"id" json:"id"`
	CompanyID  UUID      `db:"company_id" json:"company_id"`
	UserID     *UUID     `db:"user_id" json:"user_id,omitempty"`
	UserEmail  *string   `db:"user_email" json:"user_email,omitempty"`
	Action     string    `db:"action" json:"action"`
	Resource   string    `db:"resource" json:"resource"`
	ResourceID *string   `db:"resource_id" json:"resource_id,omitempty"`
	OldValue   JSONB     `db:"old_value" json:"old_value,omitempty"`
	NewValue   JSONB     `db:"new_value" json:"new_value,omitempty"`
	IPAddress  *string   `db:"ip_address" json:"ip_address,omitempty"`
	UserAgent  *string   `db:"user_agent" json:"user_agent,omitempty"`
	RequestID  *UUID     `db:"request_id" json:"request_id,omitempty"`
	Severity   string    `db:"severity" json:"severity"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

// ============================================================
// UTILITY TYPES
// ============================================================

type UUID = uuid.UUID
type JSONB = map[string]interface{}

// Money represents a monetary amount with currency
type Money struct {
	Amount   decimal.Decimal `json:"amount"`
	Currency string          `json:"currency"`
}

// MoneyWithBase adds base currency conversion info
type MoneyWithBase struct {
	Amount        decimal.Decimal `json:"amount"`
	Currency      string          `json:"currency"`
	BaseCurrency  string          `json:"base_currency"`
	BaseAmount    decimal.Decimal `json:"base_amount"`
	ExchangeRate  decimal.Decimal `json:"exchange_rate"`
	RateDate      *time.Time      `json:"rate_date,omitempty"`
	RateIsStale   bool            `json:"rate_is_stale"`
}

// WorkingCapitalKPIs holds calculated WC metrics
type WorkingCapitalKPIs struct {
	// Core KPIs
	DSO decimal.Decimal `json:"dso"`
	DPO decimal.Decimal `json:"dpo"`
	DIO decimal.Decimal `json:"dio"`
	CCC decimal.Decimal `json:"ccc"` // CCC = DSO + DIO - DPO

	// Balances
	TotalReceivables  decimal.Decimal `json:"total_receivables"`
	TotalPayables     decimal.Decimal `json:"total_payables"`
	TotalInventory    decimal.Decimal `json:"total_inventory"`
	CurrentAssets     decimal.Decimal `json:"current_assets"`
	CurrentLiabilities decimal.Decimal `json:"current_liabilities"`
	NetWorkingCapital  decimal.Decimal `json:"net_working_capital"`

	// Ratios
	CurrentRatio        decimal.Decimal `json:"current_ratio"`
	QuickRatio          decimal.Decimal `json:"quick_ratio"`
	ReceivablesTurnover decimal.Decimal `json:"receivables_turnover"`
	PayablesTurnover    decimal.Decimal `json:"payables_turnover"`

	// Revenue/Expense base
	TotalRevenue         decimal.Decimal `json:"total_revenue"`
	TotalCOGS            decimal.Decimal `json:"total_cogs"`
	Period               string          `json:"period"`
	PeriodDays           int             `json:"period_days"`
	AsOfDate             time.Time       `json:"as_of_date"`

	// Targets & variances
	DSOTarget  *decimal.Decimal `json:"dso_target,omitempty"`
	DPOTarget  *decimal.Decimal `json:"dpo_target,omitempty"`
	DSOVariance *decimal.Decimal `json:"dso_variance,omitempty"`
	DPOVariance *decimal.Decimal `json:"dpo_variance,omitempty"`
}

// CashPosition represents a point-in-time cash summary
type CashPosition struct {
	AsOfDate       time.Time            `json:"as_of_date"`
	OpeningBalance decimal.Decimal      `json:"opening_balance"`
	TotalInflows   decimal.Decimal      `json:"total_inflows"`
	TotalOutflows  decimal.Decimal      `json:"total_outflows"`
	ClosingBalance decimal.Decimal      `json:"closing_balance"`
	NetMovement    decimal.Decimal      `json:"net_movement"`
	Currency       string               `json:"currency"`
	ByBank         []BankBalance        `json:"by_bank,omitempty"`
	ByCurrency     []CurrencyBalance    `json:"by_currency,omitempty"`
	ByEntity       []EntityBalance      `json:"by_entity,omitempty"`
}

type BankBalance struct {
	BankName            string          `json:"bank_name"`
	AccountNumberMasked string          `json:"account_number_masked"`
	Currency            string          `json:"currency"`
	AvailableBalance    decimal.Decimal `json:"available_balance"`
	LedgerBalance       decimal.Decimal `json:"ledger_balance"`
	BaseBalance         decimal.Decimal `json:"base_balance"`
	ExchangeRate        decimal.Decimal `json:"exchange_rate"`
}

type CurrencyBalance struct {
	Currency        string          `json:"currency"`
	Amount          decimal.Decimal `json:"amount"`
	BaseAmount      decimal.Decimal `json:"base_amount"`
	ExchangeRate    decimal.Decimal `json:"exchange_rate"`
	AccountCount    int             `json:"account_count"`
}

type EntityBalance struct {
	EntityID   UUID            `json:"entity_id"`
	EntityName string          `json:"entity_name"`
	Country    string          `json:"country"`
	Balance    decimal.Decimal `json:"balance"`
	Currency   string          `json:"currency"`
	BaseBalance decimal.Decimal `json:"base_balance"`
}

// AgingBucket for AR/AP aging
type AgingBucket struct {
	Label    string          `json:"label"`
	MinDays  int             `json:"min_days"`
	MaxDays  *int            `json:"max_days,omitempty"`
	Amount   decimal.Decimal `json:"amount"`
	Count    int             `json:"count"`
	Percent  float64         `json:"percent"`
}

// DashboardSummary for executive dashboard
type DashboardSummary struct {
	CashPosition      CashPosition        `json:"cash_position"`
	WorkingCapital    WorkingCapitalKPIs  `json:"working_capital"`
	ARSummary         ARSummary           `json:"ar_summary"`
	APSummary         APSummary           `json:"ap_summary"`
	FXSummary         FXSummary           `json:"fx_summary"`
	ActiveAlerts      []Alert             `json:"active_alerts"`
	RecentActivity    []AuditLog          `json:"recent_activity,omitempty"`
	LiquidityRisk     string              `json:"liquidity_risk"`
	FXRisk            string              `json:"fx_risk"`
	ReceivablesRisk   string              `json:"receivables_risk"`
	GeneratedAt       time.Time           `json:"generated_at"`
}

type ARSummary struct {
	TotalOutstanding decimal.Decimal `json:"total_outstanding"`
	TotalOverdue     decimal.Decimal `json:"total_overdue"`
	InvoiceCount     int             `json:"invoice_count"`
	OverdueCount     int             `json:"overdue_count"`
	AgingBuckets     []AgingBucket   `json:"aging_buckets"`
	LargestCustomer  *string         `json:"largest_customer,omitempty"`
	ConcentrationPct float64         `json:"concentration_pct"`
}

type APSummary struct {
	TotalOutstanding  decimal.Decimal `json:"total_outstanding"`
	TotalOverdue      decimal.Decimal `json:"total_overdue"`
	BillCount         int             `json:"bill_count"`
	OverdueCount      int             `json:"overdue_count"`
	DueIn3Days        decimal.Decimal `json:"due_in_3_days"`
	AgingBuckets      []AgingBucket   `json:"aging_buckets"`
	LargestSupplier   *string         `json:"largest_supplier,omitempty"`
	ConcentrationPct  float64         `json:"concentration_pct"`
}

type FXSummary struct {
	TotalGrossExposure decimal.Decimal   `json:"total_gross_exposure"`
	TotalNetExposure   decimal.Decimal   `json:"total_net_exposure"`
	TotalHedged        decimal.Decimal   `json:"total_hedged"`
	CoverageRatio      decimal.Decimal   `json:"coverage_ratio"`
	BaseCurrency       string            `json:"base_currency"`
	ByExposure         []ExposureSummary `json:"by_currency"`
	HighestRisk        *string           `json:"highest_risk_currency,omitempty"`
}

type ExposureSummary struct {
	Currency      string          `json:"currency"`
	NetExposure   decimal.Decimal `json:"net_exposure"`
	BaseValue     decimal.Decimal `json:"base_value"`
	ExchangeRate  decimal.Decimal `json:"exchange_rate"`
	RiskLevel     string          `json:"risk_level"`
	Direction     string          `json:"direction"` // LONG, SHORT, NEUTRAL
}
