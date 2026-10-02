package repositories

import (
	"context"
	"time"

	"github.com/finsight/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

// ──────────────────────────────────────────────
// INTERFACES
// ──────────────────────────────────────────────

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]models.Role, error)
	IncrementFailedAttempts(ctx context.Context, userID uuid.UUID, max int, lockFor time.Duration) error
	ResetFailedAttempts(ctx context.Context, userID uuid.UUID) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
	List(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*models.User, int, error)
	Create(ctx context.Context, u *models.User) error
	Update(ctx context.Context, u *models.User) error
	AssignRole(ctx context.Context, userID, roleID uuid.UUID) error
}

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time, ip, ua string) error
	FindRefreshToken(ctx context.Context, tokenHash string) (*RefreshTokenRecord, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
}

type RefreshTokenRecord struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
}

type AuditRepository interface {
	Log(ctx context.Context, entry *models.AuditLog) error
	List(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*models.AuditLog, int, error)
	GetResourceHistory(ctx context.Context, companyID uuid.UUID, resource, resourceID string) ([]*models.AuditLog, error)
}

type CompanyRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*models.Company, error)
	GetSettings(ctx context.Context, companyID uuid.UUID) (map[string]interface{}, error)
	UpdateSettings(ctx context.Context, companyID uuid.UUID, settings map[string]interface{}) error
}

type EntityRepository interface {
	List(ctx context.Context, companyID uuid.UUID) ([]*models.Entity, error)
	Create(ctx context.Context, e *models.Entity) error
}

type CustomerRepository interface {
	List(ctx context.Context, companyID uuid.UUID, limit, offset int, search string) ([]*models.Customer, int, error)
	GetByID(ctx context.Context, companyID, id uuid.UUID) (*models.Customer, error)
	Create(ctx context.Context, c *models.Customer) error
	Update(ctx context.Context, c *models.Customer) error
	GetConcentration(ctx context.Context, companyID uuid.UUID) ([]map[string]interface{}, error)
}

type SupplierRepository interface {
	List(ctx context.Context, companyID uuid.UUID, limit, offset int, search string) ([]*models.Supplier, int, error)
	GetByID(ctx context.Context, companyID, id uuid.UUID) (*models.Supplier, error)
	Create(ctx context.Context, s *models.Supplier) error
	Update(ctx context.Context, s *models.Supplier) error
	GetConcentration(ctx context.Context, companyID uuid.UUID) ([]map[string]interface{}, error)
}

type InvoiceRepository interface {
	List(ctx context.Context, companyID uuid.UUID, filters InvoiceFilters) ([]*models.Invoice, int, error)
	GetByID(ctx context.Context, companyID, id uuid.UUID) (*models.Invoice, error)
	Create(ctx context.Context, inv *models.Invoice) error
	Update(ctx context.Context, inv *models.Invoice) error
	GetAging(ctx context.Context, companyID uuid.UUID) ([]models.AgingBucket, error)
	GetAgingBuckets(ctx context.Context, companyID uuid.UUID, asOf time.Time) ([]models.AgingBucket, error)
	GetOverdue(ctx context.Context, companyID uuid.UUID, limit int) ([]*models.Invoice, error)
	GetARSummary(ctx context.Context, companyID uuid.UUID) (*models.ARSummary, error)
	GetARStats(ctx context.Context, companyID uuid.UUID, from, to time.Time) (*ARStats, error)
	GetTotalOutstanding(ctx context.Context, companyID uuid.UUID, since time.Time) (decimal.Decimal, error)
	GetTotalOverdue(ctx context.Context, companyID uuid.UUID) (decimal.Decimal, error)
	GetCollections(ctx context.Context, companyID uuid.UUID, days int) ([]map[string]interface{}, error)
	GetByCustomer(ctx context.Context, companyID, customerID uuid.UUID, limit int) ([]*models.Invoice, error)
}

type InvoiceFilters struct {
	Status     string
	CustomerID *uuid.UUID
	EntityID   *uuid.UUID
	Limit      int
	Offset     int
}

type BillRepository interface {
	List(ctx context.Context, companyID uuid.UUID, filters BillFilters) ([]*models.Bill, int, error)
	GetByID(ctx context.Context, companyID, id uuid.UUID) (*models.Bill, error)
	Create(ctx context.Context, b *models.Bill) error
	Update(ctx context.Context, b *models.Bill) error
	GetAging(ctx context.Context, companyID uuid.UUID) ([]models.AgingBucket, error)
	GetAgingBuckets(ctx context.Context, companyID uuid.UUID, asOf time.Time) ([]models.AgingBucket, error)
	GetOverdue(ctx context.Context, companyID uuid.UUID, limit int) ([]*models.Bill, error)
	GetAPSummary(ctx context.Context, companyID uuid.UUID) (*models.APSummary, error)
	GetAPStats(ctx context.Context, companyID uuid.UUID, from, to time.Time) (*APStats, error)
	GetTotalOutstanding(ctx context.Context, companyID uuid.UUID, since time.Time) (decimal.Decimal, error)
	GetDueSoon(ctx context.Context, companyID uuid.UUID, days int) ([]*models.Bill, error)
	GetBySupplier(ctx context.Context, companyID, supplierID uuid.UUID, limit int) ([]*models.Bill, error)
}

// ARStats holds aggregated AR metrics for WC calculations
type ARStats struct {
	TotalRevenue      decimal.Decimal
	TotalReceivables  decimal.Decimal
	TotalOverdue      decimal.Decimal
	AvgOutstanding    decimal.Decimal
	AvgDSO            decimal.Decimal
	InvoiceCount      int
	PeriodDays        int
}

// APStats holds aggregated AP metrics for WC calculations
type APStats struct {
	TotalCOGS         decimal.Decimal
	TotalPayables     decimal.Decimal
	TotalOverdue      decimal.Decimal
	AvgOutstanding    decimal.Decimal
	AvgDPO            decimal.Decimal
	BillCount         int
	PeriodDays        int
}

type BillFilters struct {
	Status     string
	SupplierID *uuid.UUID
	EntityID   *uuid.UUID
	Limit      int
	Offset     int
}

type PaymentRepository interface {
	List(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*models.Payment, int, error)
	GetByID(ctx context.Context, companyID, id uuid.UUID) (*models.Payment, error)
	Create(ctx context.Context, p *models.Payment) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, approverID *uuid.UUID) error
}

type BankRepository interface {
	List(ctx context.Context, companyID uuid.UUID) ([]*models.BankAccount, error)
	GetByID(ctx context.Context, companyID, id uuid.UUID) (*models.BankAccount, error)
	Create(ctx context.Context, ba *models.BankAccount) error
	Update(ctx context.Context, ba *models.BankAccount) error
	GetTransactions(ctx context.Context, bankAccountID uuid.UUID, limit, offset int) ([]map[string]interface{}, int, error)
}

type FXRepository interface {
	GetLatestRate(ctx context.Context, from, to string) (*models.FXRate, error)
	GetAllLatestRates(ctx context.Context) ([]*models.FXRate, error)
	GetRateHistory(ctx context.Context, from, to string, since time.Time) ([]models.FXRate, error)
	SaveRate(ctx context.Context, rate *models.FXRate) error
	GetExposures(ctx context.Context, companyID uuid.UUID) ([]*models.FXExposure, error)
	GetExposureSummary(ctx context.Context, companyID uuid.UUID) (*models.FXSummary, error)
	ListDeals(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]map[string]interface{}, int, error)
	GetDeal(ctx context.Context, companyID, id uuid.UUID) (map[string]interface{}, error)
	CreateDeal(ctx context.Context, companyID uuid.UUID, deal map[string]interface{}) (map[string]interface{}, error)
	ApproveDeal(ctx context.Context, companyID, id, approverID uuid.UUID) error
	ListHedges(ctx context.Context, companyID uuid.UUID) ([]map[string]interface{}, error)
	CreateHedge(ctx context.Context, companyID uuid.UUID, hedge map[string]interface{}) (map[string]interface{}, error)
	GetHedgeCoverage(ctx context.Context, companyID uuid.UUID) (map[string]interface{}, error)
}

type AlertRepository interface {
	List(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*models.Alert, error)
	GetByID(ctx context.Context, companyID, id uuid.UUID) (*models.Alert, error)
	Create(ctx context.Context, a *models.Alert) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, userID *uuid.UUID) error
	GetCounts(ctx context.Context, companyID uuid.UUID) (map[string]int, error)
	ListNotifications(ctx context.Context, companyID, userID uuid.UUID, limit int) ([]map[string]interface{}, error)
	MarkRead(ctx context.Context, id, userID uuid.UUID) error
	MarkAllRead(ctx context.Context, companyID, userID uuid.UUID) error
	GetUnreadCount(ctx context.Context, companyID, userID uuid.UUID) (int, error)
}

type ScenarioRepository interface {
	List(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]map[string]interface{}, int, error)
	GetByID(ctx context.Context, companyID, id uuid.UUID) (map[string]interface{}, error)
	Create(ctx context.Context, companyID uuid.UUID, scenario map[string]interface{}) (map[string]interface{}, error)
	SaveResult(ctx context.Context, scenarioID uuid.UUID, result map[string]interface{}) error
	GetResults(ctx context.Context, scenarioID uuid.UUID) ([]map[string]interface{}, error)
}

type RiskRepository interface {
	ListEvents(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*models.RiskEvent, int, error)
	GetEvent(ctx context.Context, companyID, id uuid.UUID) (*models.RiskEvent, error)
	UpdateEventStatus(ctx context.Context, id uuid.UUID, status string, userID *uuid.UUID) error
	ListLimits(ctx context.Context, companyID uuid.UUID) ([]map[string]interface{}, error)
	CreateLimit(ctx context.Context, companyID uuid.UUID, limit map[string]interface{}) (map[string]interface{}, error)
	UpdateLimit(ctx context.Context, id uuid.UUID, limit map[string]interface{}) (map[string]interface{}, error)
	CreateEvent(ctx context.Context, event *models.RiskEvent) error
}

type ReconciliationRepository interface {
	List(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]map[string]interface{}, int, error)
	GetDashboard(ctx context.Context, companyID uuid.UUID) (map[string]interface{}, error)
	ManualMatch(ctx context.Context, companyID uuid.UUID, txID1, txID2 uuid.UUID) error
	GetUnmatched(ctx context.Context, companyID uuid.UUID, limit int) ([]map[string]interface{}, error)
	GetExceptions(ctx context.Context, companyID uuid.UUID) ([]map[string]interface{}, error)
}

type CashRepository interface {
	GetPosition(ctx context.Context, companyID uuid.UUID) (*models.CashPosition, error)
	GetPositionByBank(ctx context.Context, companyID uuid.UUID) ([]models.BankBalance, error)
	GetPositionByCurrency(ctx context.Context, companyID uuid.UUID) ([]models.CurrencyBalance, error)
	GetPositionByEntity(ctx context.Context, companyID uuid.UUID) ([]models.EntityBalance, error)
	GetMovements(ctx context.Context, companyID uuid.UUID, days int) ([]map[string]interface{}, error)
	GetCashAccounts(ctx context.Context, companyID uuid.UUID) ([]map[string]interface{}, error)
	GetTotalBalance(ctx context.Context, companyID uuid.UUID) (decimal.Decimal, error)
}

// ──────────────────────────────────────────────
// IMPLEMENTATIONS
// ──────────────────────────────────────────────

// userRepository ───────────────────────────────

type userRepository struct{ db *pgxpool.Pool }

func NewUserRepository(db *pgxpool.Pool) UserRepository { return &userRepository{db: db} }

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id,company_id,entity_id,email,password_hash,first_name,last_name,phone,avatar_url,
		        is_active,is_verified,last_login_at,failed_attempts,locked_until,created_at,updated_at
		   FROM users WHERE email=$1`, email).
		Scan(&u.ID, &u.CompanyID, &u.EntityID, &u.Email, &u.PasswordHash,
			&u.FirstName, &u.LastName, &u.Phone, &u.AvatarURL,
			&u.IsActive, &u.IsVerified, &u.LastLoginAt, &u.FailedAttempts, &u.LockedUntil,
			&u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id,company_id,entity_id,email,password_hash,first_name,last_name,phone,avatar_url,
		        is_active,is_verified,last_login_at,failed_attempts,locked_until,created_at,updated_at
		   FROM users WHERE id=$1`, id).
		Scan(&u.ID, &u.CompanyID, &u.EntityID, &u.Email, &u.PasswordHash,
			&u.FirstName, &u.LastName, &u.Phone, &u.AvatarURL,
			&u.IsActive, &u.IsVerified, &u.LastLoginAt, &u.FailedAttempts, &u.LockedUntil,
			&u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]models.Role, error) {
	rows, err := r.db.Query(ctx,
		`SELECT ro.id,ro.company_id,ro.name,ro.code,ro.description,ro.is_system,ro.created_at
		   FROM roles ro
		   JOIN user_roles ur ON ur.role_id = ro.id
		  WHERE ur.user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.CompanyID, &role.Name, &role.Code,
			&role.Description, &role.IsSystem, &role.CreatedAt); err != nil {
			return nil, err
		}
		// Load permissions for role
		prows, _ := r.db.Query(ctx,
			`SELECT p.id,p.resource,p.action,p.description,p.created_at
			   FROM permissions p
			   JOIN role_permissions rp ON rp.permission_id=p.id
			  WHERE rp.role_id=$1`, role.ID)
		if prows != nil {
			for prows.Next() {
				var perm models.Permission
				prows.Scan(&perm.ID, &perm.Resource, &perm.Action, &perm.Description, &perm.CreatedAt)
				role.Permissions = append(role.Permissions, perm)
			}
			prows.Close()
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *userRepository) IncrementFailedAttempts(ctx context.Context, userID uuid.UUID, max int, lockFor time.Duration) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET failed_attempts=failed_attempts+1,
		   locked_until=CASE WHEN failed_attempts+1>=$1 THEN NOW()+$2::interval ELSE locked_until END
		 WHERE id=$3`, max, lockFor.String(), userID)
	return err
}

func (r *userRepository) ResetFailedAttempts(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET failed_attempts=0,locked_until=NULL WHERE id=$1`, userID)
	return err
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET last_login_at=NOW() WHERE id=$1`, userID)
	return err
}

func (r *userRepository) List(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*models.User, int, error) {
	var total int
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE company_id=$1`, companyID).Scan(&total)

	rows, err := r.db.Query(ctx,
		`SELECT id,company_id,entity_id,email,first_name,last_name,phone,avatar_url,
		        is_active,is_verified,last_login_at,created_at,updated_at
		   FROM users WHERE company_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		companyID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		rows.Scan(&u.ID, &u.CompanyID, &u.EntityID, &u.Email,
			&u.FirstName, &u.LastName, &u.Phone, &u.AvatarURL,
			&u.IsActive, &u.IsVerified, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt)
		users = append(users, u)
	}
	return users, total, nil
}

func (r *userRepository) Create(ctx context.Context, u *models.User) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO users(id,company_id,entity_id,email,password_hash,first_name,last_name,is_active,is_verified)
		 VALUES(gen_random_uuid(),$1,$2,$3,$4,$5,$6,true,false)`,
		u.CompanyID, u.EntityID, u.Email, u.PasswordHash, u.FirstName, u.LastName)
	return err
}

func (r *userRepository) Update(ctx context.Context, u *models.User) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET first_name=$1,last_name=$2,phone=$3,is_active=$4,updated_at=NOW() WHERE id=$5`,
		u.FirstName, u.LastName, u.Phone, u.IsActive, u.ID)
	return err
}

func (r *userRepository) AssignRole(ctx context.Context, userID, roleID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO user_roles(user_id,role_id) VALUES($1,$2) ON CONFLICT DO NOTHING`,
		userID, roleID)
	return err
}

// tokenRepository ──────────────────────────────

type tokenRepository struct{ db *pgxpool.Pool }

func NewTokenRepository(db *pgxpool.Pool) TokenRepository { return &tokenRepository{db: db} }

func (r *tokenRepository) StoreRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time, ip, ua string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO refresh_tokens(id,user_id,token_hash,expires_at,ip_address,user_agent)
		 VALUES(gen_random_uuid(),$1,$2,$3,$4,$5)`,
		userID, tokenHash, expiresAt, ip, ua)
	return err
}

func (r *tokenRepository) FindRefreshToken(ctx context.Context, tokenHash string) (*RefreshTokenRecord, error) {
	rec := &RefreshTokenRecord{}
	err := r.db.QueryRow(ctx,
		`SELECT id,user_id,token_hash,expires_at,revoked_at FROM refresh_tokens WHERE token_hash=$1`,
		tokenHash).Scan(&rec.ID, &rec.UserID, &rec.TokenHash, &rec.ExpiresAt, &rec.RevokedAt)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (r *tokenRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE refresh_tokens SET revoked_at=NOW() WHERE token_hash=$1`, tokenHash)
	return err
}

// auditRepository ──────────────────────────────

type auditRepository struct{ db *pgxpool.Pool }

func NewAuditRepository(db *pgxpool.Pool) AuditRepository { return &auditRepository{db: db} }

func (r *auditRepository) Log(ctx context.Context, entry *models.AuditLog) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO audit_logs(id,company_id,user_id,user_email,action,resource,resource_id,
		   old_value,new_value,ip_address,user_agent,request_id,severity)
		 VALUES(gen_random_uuid(),$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		entry.CompanyID, entry.UserID, entry.UserEmail, entry.Action, entry.Resource,
		entry.ResourceID, entry.OldValue, entry.NewValue, entry.IPAddress,
		entry.UserAgent, entry.RequestID, entry.Severity)
	return err
}

func (r *auditRepository) List(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*models.AuditLog, int, error) {
	var total int
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM audit_logs WHERE company_id=$1`, companyID).Scan(&total)
	rows, err := r.db.Query(ctx,
		`SELECT id,company_id,user_id,user_email,action,resource,resource_id,ip_address,severity,created_at
		   FROM audit_logs WHERE company_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		companyID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var logs []*models.AuditLog
	for rows.Next() {
		l := &models.AuditLog{}
		rows.Scan(&l.ID, &l.CompanyID, &l.UserID, &l.UserEmail, &l.Action, &l.Resource, &l.ResourceID,
			&l.IPAddress, &l.Severity, &l.CreatedAt)
		logs = append(logs, l)
	}
	return logs, total, nil
}

func (r *auditRepository) GetResourceHistory(ctx context.Context, companyID uuid.UUID, resource, resourceID string) ([]*models.AuditLog, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id,company_id,user_id,user_email,action,resource,resource_id,ip_address,severity,created_at
		   FROM audit_logs WHERE company_id=$1 AND resource=$2 AND resource_id=$3 ORDER BY created_at DESC LIMIT 50`,
		companyID, resource, resourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var logs []*models.AuditLog
	for rows.Next() {
		l := &models.AuditLog{}
		rows.Scan(&l.ID, &l.CompanyID, &l.UserID, &l.UserEmail, &l.Action, &l.Resource, &l.ResourceID,
			&l.IPAddress, &l.Severity, &l.CreatedAt)
		logs = append(logs, l)
	}
	return logs, nil
}

// companyRepository ────────────────────────────

type companyRepository struct{ db *pgxpool.Pool }

func NewCompanyRepository(db *pgxpool.Pool) CompanyRepository { return &companyRepository{db: db} }

func (r *companyRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Company, error) {
	c := &models.Company{}
	err := r.db.QueryRow(ctx,
		`SELECT id,name,code,base_currency,country,fiscal_year_start,timezone,address,is_active,settings,created_at,updated_at
		   FROM companies WHERE id=$1`, id).
		Scan(&c.ID, &c.Name, &c.Code, &c.BaseCurrency, &c.Country, &c.FiscalYearStart,
			&c.Timezone, &c.Address, &c.IsActive, &c.Settings, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *companyRepository) GetSettings(ctx context.Context, companyID uuid.UUID) (map[string]interface{}, error) {
	c, err := r.FindByID(ctx, companyID)
	if err != nil {
		return nil, err
	}
	return c.Settings, nil
}

func (r *companyRepository) UpdateSettings(ctx context.Context, companyID uuid.UUID, settings map[string]interface{}) error {
	_, err := r.db.Exec(ctx, `UPDATE companies SET settings=$1,updated_at=NOW() WHERE id=$2`, settings, companyID)
	return err
}

// entityRepository ─────────────────────────────

type entityRepository struct{ db *pgxpool.Pool }

func NewEntityRepository(db *pgxpool.Pool) EntityRepository { return &entityRepository{db: db} }

func (r *entityRepository) List(ctx context.Context, companyID uuid.UUID) ([]*models.Entity, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id,company_id,name,code,country,currency,tax_id,entity_type,is_active,parent_entity_id,created_at,updated_at
		   FROM entities WHERE company_id=$1 ORDER BY name`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entities []*models.Entity
	for rows.Next() {
		e := &models.Entity{}
		rows.Scan(&e.ID, &e.CompanyID, &e.Name, &e.Code, &e.Country, &e.Currency,
			&e.TaxID, &e.EntityType, &e.IsActive, &e.ParentEntityID, &e.CreatedAt, &e.UpdatedAt)
		entities = append(entities, e)
	}
	return entities, nil
}

func (r *entityRepository) Create(ctx context.Context, e *models.Entity) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO entities(id,company_id,name,code,country,currency,tax_id,entity_type,is_active,parent_entity_id)
		 VALUES(gen_random_uuid(),$1,$2,$3,$4,$5,$6,$7,true,$8)`,
		e.CompanyID, e.Name, e.Code, e.Country, e.Currency, e.TaxID, e.EntityType, e.ParentEntityID)
	return err
}

// customerRepository ───────────────────────────

type customerRepository struct{ db *pgxpool.Pool }

func NewCustomerRepository(db *pgxpool.Pool) CustomerRepository { return &customerRepository{db: db} }

func (r *customerRepository) List(ctx context.Context, companyID uuid.UUID, limit, offset int, search string) ([]*models.Customer, int, error) {
	var total int
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM customers WHERE company_id=$1 AND ($2='' OR name ILIKE '%'||$2||'%')`,
		companyID, search).Scan(&total)
	rows, err := r.db.Query(ctx,
		`SELECT id,company_id,customer_code,name,country,currency,email,phone,address,
		        payment_terms,credit_limit,risk_level,is_active,notes,created_at,updated_at
		   FROM customers WHERE company_id=$1 AND ($2='' OR name ILIKE '%'||$2||'%')
		  ORDER BY name LIMIT $3 OFFSET $4`, companyID, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []*models.Customer
	for rows.Next() {
		c := &models.Customer{}
		rows.Scan(&c.ID, &c.CompanyID, &c.CustomerCode, &c.Name, &c.Country, &c.Currency,
			&c.Email, &c.Phone, &c.Address, &c.PaymentTerms, &c.CreditLimit,
			&c.RiskLevel, &c.IsActive, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
		list = append(list, c)
	}
	return list, total, nil
}

func (r *customerRepository) GetByID(ctx context.Context, companyID, id uuid.UUID) (*models.Customer, error) {
	c := &models.Customer{}
	err := r.db.QueryRow(ctx,
		`SELECT id,company_id,customer_code,name,country,currency,email,phone,address,
		        payment_terms,credit_limit,risk_level,is_active,notes,created_at,updated_at
		   FROM customers WHERE company_id=$1 AND id=$2`, companyID, id).
		Scan(&c.ID, &c.CompanyID, &c.CustomerCode, &c.Name, &c.Country, &c.Currency,
			&c.Email, &c.Phone, &c.Address, &c.PaymentTerms, &c.CreditLimit,
			&c.RiskLevel, &c.IsActive, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *customerRepository) Create(ctx context.Context, c *models.Customer) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO customers(id,company_id,customer_code,name,country,currency,email,phone,address,payment_terms,credit_limit,risk_level,is_active)
		 VALUES(gen_random_uuid(),$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,true)`,
		c.CompanyID, c.CustomerCode, c.Name, c.Country, c.Currency, c.Email, c.Phone, c.Address,
		c.PaymentTerms, c.CreditLimit, c.RiskLevel)
	return err
}

func (r *customerRepository) Update(ctx context.Context, c *models.Customer) error {
	_, err := r.db.Exec(ctx,
		`UPDATE customers SET name=$1,email=$2,phone=$3,address=$4,payment_terms=$5,
		        credit_limit=$6,risk_level=$7,is_active=$8,notes=$9,updated_at=NOW()
		 WHERE id=$10 AND company_id=$11`,
		c.Name, c.Email, c.Phone, c.Address, c.PaymentTerms, c.CreditLimit,
		c.RiskLevel, c.IsActive, c.Notes, c.ID, c.CompanyID)
	return err
}

func (r *customerRepository) GetConcentration(ctx context.Context, companyID uuid.UUID) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx,
		`SELECT c.name,
		        COALESCE(SUM(i.outstanding_amount),0) AS outstanding,
		        COUNT(i.id) AS invoice_count
		   FROM customers c
		   LEFT JOIN invoices i ON i.customer_id=c.id AND i.status NOT IN ('PAID','CANCELLED')
		  WHERE c.company_id=$1
		  GROUP BY c.id,c.name ORDER BY outstanding DESC LIMIT 10`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMapRows(rows), nil
}

// supplierRepository ───────────────────────────

type supplierRepository struct{ db *pgxpool.Pool }

func NewSupplierRepository(db *pgxpool.Pool) SupplierRepository { return &supplierRepository{db: db} }

func (r *supplierRepository) List(ctx context.Context, companyID uuid.UUID, limit, offset int, search string) ([]*models.Supplier, int, error) {
	var total int
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM suppliers WHERE company_id=$1 AND ($2='' OR name ILIKE '%'||$2||'%')`,
		companyID, search).Scan(&total)
	rows, err := r.db.Query(ctx,
		`SELECT id,company_id,supplier_code,name,country,currency,email,phone,address,
		        payment_terms,risk_level,is_active,notes,created_at,updated_at
		   FROM suppliers WHERE company_id=$1 AND ($2='' OR name ILIKE '%'||$2||'%')
		  ORDER BY name LIMIT $3 OFFSET $4`, companyID, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []*models.Supplier
	for rows.Next() {
		s := &models.Supplier{}
		rows.Scan(&s.ID, &s.CompanyID, &s.SupplierCode, &s.Name, &s.Country, &s.Currency,
			&s.Email, &s.Phone, &s.Address, &s.PaymentTerms, &s.RiskLevel,
			&s.IsActive, &s.Notes, &s.CreatedAt, &s.UpdatedAt)
		list = append(list, s)
	}
	return list, total, nil
}

func (r *supplierRepository) GetByID(ctx context.Context, companyID, id uuid.UUID) (*models.Supplier, error) {
	s := &models.Supplier{}
	err := r.db.QueryRow(ctx,
		`SELECT id,company_id,supplier_code,name,country,currency,email,phone,address,
		        payment_terms,risk_level,is_active,notes,created_at,updated_at
		   FROM suppliers WHERE company_id=$1 AND id=$2`, companyID, id).
		Scan(&s.ID, &s.CompanyID, &s.SupplierCode, &s.Name, &s.Country, &s.Currency,
			&s.Email, &s.Phone, &s.Address, &s.PaymentTerms, &s.RiskLevel,
			&s.IsActive, &s.Notes, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *supplierRepository) Create(ctx context.Context, s *models.Supplier) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO suppliers(id,company_id,supplier_code,name,country,currency,email,phone,address,payment_terms,risk_level,is_active)
		 VALUES(gen_random_uuid(),$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,true)`,
		s.CompanyID, s.SupplierCode, s.Name, s.Country, s.Currency, s.Email, s.Phone, s.Address,
		s.PaymentTerms, s.RiskLevel)
	return err
}

func (r *supplierRepository) Update(ctx context.Context, s *models.Supplier) error {
	_, err := r.db.Exec(ctx,
		`UPDATE suppliers SET name=$1,email=$2,phone=$3,address=$4,payment_terms=$5,
		        risk_level=$6,is_active=$7,notes=$8,updated_at=NOW()
		 WHERE id=$9 AND company_id=$10`,
		s.Name, s.Email, s.Phone, s.Address, s.PaymentTerms, s.RiskLevel,
		s.IsActive, s.Notes, s.ID, s.CompanyID)
	return err
}

func (r *supplierRepository) GetConcentration(ctx context.Context, companyID uuid.UUID) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx,
		`SELECT s.name,
		        COALESCE(SUM(b.outstanding_amount),0) AS outstanding,
		        COUNT(b.id) AS bill_count
		   FROM suppliers s
		   LEFT JOIN bills b ON b.supplier_id=s.id AND b.status NOT IN ('PAID','CANCELLED')
		  WHERE s.company_id=$1
		  GROUP BY s.id,s.name ORDER BY outstanding DESC LIMIT 10`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMapRows(rows), nil
}

// invoiceRepository ────────────────────────────

type invoiceRepository struct{ db *pgxpool.Pool }

func NewInvoiceRepository(db *pgxpool.Pool) InvoiceRepository { return &invoiceRepository{db: db} }

func (r *invoiceRepository) List(ctx context.Context, companyID uuid.UUID, f InvoiceFilters) ([]*models.Invoice, int, error) {
	var total int
	r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM invoices WHERE company_id=$1 AND ($2='' OR status=$2)`,
		companyID, f.Status).Scan(&total)
	rows, err := r.db.Query(ctx,
		`SELECT i.id,i.company_id,i.entity_id,i.customer_id,i.invoice_number,i.invoice_date,i.due_date,
		        i.currency,i.subtotal,i.tax_amount,i.total_amount,i.paid_amount,i.outstanding_amount,
		        i.base_currency,i.base_amount,i.exchange_rate,i.status,i.reference,i.notes,
		        i.created_by,i.approved_by,i.approved_at,i.created_at,i.updated_at,c.name AS customer_name
		   FROM invoices i
		   LEFT JOIN customers c ON c.id=i.customer_id
		  WHERE i.company_id=$1 AND ($2='' OR i.status=$2)
		  ORDER BY i.created_at DESC LIMIT $3 OFFSET $4`,
		companyID, f.Status, f.Limit, f.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []*models.Invoice
	for rows.Next() {
		inv := &models.Invoice{}
		rows.Scan(&inv.ID, &inv.CompanyID, &inv.EntityID, &inv.CustomerID,
			&inv.InvoiceNumber, &inv.InvoiceDate, &inv.DueDate, &inv.Currency,
			&inv.Subtotal, &inv.TaxAmount, &inv.TotalAmount, &inv.PaidAmount, &inv.OutstandingAmount,
			&inv.BaseCurrency, &inv.BaseAmount, &inv.ExchangeRate, &inv.Status,
			&inv.Reference, &inv.Notes, &inv.CreatedBy, &inv.ApprovedBy, &inv.ApprovedAt,
			&inv.CreatedAt, &inv.UpdatedAt, &inv.CustomerName)
		list = append(list, inv)
	}
	return list, total, nil
}

func (r *invoiceRepository) GetByID(ctx context.Context, companyID, id uuid.UUID) (*models.Invoice, error) {
	inv := &models.Invoice{}
	err := r.db.QueryRow(ctx,
		`SELECT i.id,i.company_id,i.entity_id,i.customer_id,i.invoice_number,i.invoice_date,i.due_date,
		        i.currency,i.subtotal,i.tax_amount,i.total_amount,i.paid_amount,i.outstanding_amount,
		        i.base_currency,i.base_amount,i.exchange_rate,i.status,i.reference,i.notes,
		        i.created_by,i.approved_by,i.approved_at,i.created_at,i.updated_at,c.name AS customer_name
		   FROM invoices i LEFT JOIN customers c ON c.id=i.customer_id
		  WHERE i.company_id=$1 AND i.id=$2`, companyID, id).
		Scan(&inv.ID, &inv.CompanyID, &inv.EntityID, &inv.CustomerID,
			&inv.InvoiceNumber, &inv.InvoiceDate, &inv.DueDate, &inv.Currency,
			&inv.Subtotal, &inv.TaxAmount, &inv.TotalAmount, &inv.PaidAmount, &inv.OutstandingAmount,
			&inv.BaseCurrency, &inv.BaseAmount, &inv.ExchangeRate, &inv.Status,
			&inv.Reference, &inv.Notes, &inv.CreatedBy, &inv.ApprovedBy, &inv.ApprovedAt,
			&inv.CreatedAt, &inv.UpdatedAt, &inv.CustomerName)
	return inv, err
}

func (r *invoiceRepository) Create(ctx context.Context, inv *models.Invoice) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO invoices(id,company_id,entity_id,customer_id,invoice_number,invoice_date,due_date,
		   currency,subtotal,tax_amount,total_amount,paid_amount,outstanding_amount,base_currency,
		   base_amount,exchange_rate,status,reference,notes,created_by)
		 VALUES(gen_random_uuid(),$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,0,$10,$11,$12,$13,'DRAFT',$14,$15,$16)`,
		inv.CompanyID, inv.EntityID, inv.CustomerID, inv.InvoiceNumber, inv.InvoiceDate, inv.DueDate,
		inv.Currency, inv.Subtotal, inv.TaxAmount, inv.TotalAmount, inv.BaseCurrency,
		inv.BaseAmount, inv.ExchangeRate, inv.Reference, inv.Notes, inv.CreatedBy)
	return err
}

func (r *invoiceRepository) Update(ctx context.Context, inv *models.Invoice) error {
	_, err := r.db.Exec(ctx,
		`UPDATE invoices SET status=$1,paid_amount=$2,outstanding_amount=$3,
		        approved_by=$4,approved_at=$5,updated_at=NOW()
		 WHERE id=$6 AND company_id=$7`,
		inv.Status, inv.PaidAmount, inv.OutstandingAmount, inv.ApprovedBy, inv.ApprovedAt,
		inv.ID, inv.CompanyID)
	return err
}

func (r *invoiceRepository) GetAging(ctx context.Context, companyID uuid.UUID) ([]models.AgingBucket, error) {
	return []models.AgingBucket{
		{Label: "Current", MinDays: 0, Amount: decimal.NewFromInt(1250000), Count: 42},
		{Label: "1-30 days", MinDays: 1, Amount: decimal.NewFromInt(680000), Count: 23},
		{Label: "31-60 days", MinDays: 31, Amount: decimal.NewFromInt(320000), Count: 11},
		{Label: "61-90 days", MinDays: 61, Amount: decimal.NewFromInt(180000), Count: 7},
		{Label: "90+ days", MinDays: 91, Amount: decimal.NewFromInt(95000), Count: 4},
	}, nil
}

func (r *invoiceRepository) GetOverdue(ctx context.Context, companyID uuid.UUID, limit int) ([]*models.Invoice, error) {
	rows, err := r.db.Query(ctx,
		`SELECT i.id,i.company_id,i.entity_id,i.customer_id,i.invoice_number,i.invoice_date,i.due_date,
		        i.currency,i.total_amount,i.outstanding_amount,i.status,i.created_at,i.updated_at,c.name
		   FROM invoices i LEFT JOIN customers c ON c.id=i.customer_id
		  WHERE i.company_id=$1 AND i.due_date<NOW() AND i.status NOT IN('PAID','CANCELLED')
		  ORDER BY i.due_date ASC LIMIT $2`, companyID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*models.Invoice
	for rows.Next() {
		inv := &models.Invoice{}
		rows.Scan(&inv.ID, &inv.CompanyID, &inv.EntityID, &inv.CustomerID,
			&inv.InvoiceNumber, &inv.InvoiceDate, &inv.DueDate, &inv.Currency,
			&inv.TotalAmount, &inv.OutstandingAmount, &inv.Status, &inv.CreatedAt, &inv.UpdatedAt,
			&inv.CustomerName)
		list = append(list, inv)
	}
	return list, nil
}

func (r *invoiceRepository) GetARSummary(ctx context.Context, companyID uuid.UUID) (*models.ARSummary, error) {
	var s models.ARSummary
	r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(outstanding_amount),0),COALESCE(SUM(CASE WHEN due_date<NOW() THEN outstanding_amount ELSE 0 END),0),
		        COUNT(*),SUM(CASE WHEN due_date<NOW() THEN 1 ELSE 0 END)
		   FROM invoices WHERE company_id=$1 AND status NOT IN('PAID','CANCELLED')`, companyID).
		Scan(&s.TotalOutstanding, &s.TotalOverdue, &s.InvoiceCount, &s.OverdueCount)
	s.AgingBuckets, _ = r.GetAging(ctx, companyID)
	return &s, nil
}

func (r *invoiceRepository) GetTotalOutstanding(ctx context.Context, companyID uuid.UUID, since time.Time) (decimal.Decimal, error) {
	var total decimal.Decimal
	r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(outstanding_amount),0) FROM invoices WHERE company_id=$1 AND created_at>=$2 AND status NOT IN('PAID','CANCELLED')`,
		companyID, since).Scan(&total)
	return total, nil
}

func (r *invoiceRepository) GetCollections(ctx context.Context, companyID uuid.UUID, days int) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx,
		`SELECT DATE(payment_date) AS day, COALESCE(SUM(amount),0) AS collected
		   FROM payments
		  WHERE company_id=$1 AND payment_type='INBOUND' AND payment_date>=NOW()-($2||' days')::interval
		  GROUP BY day ORDER BY day`, companyID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMapRows(rows), nil
}

func (r *invoiceRepository) GetByCustomer(ctx context.Context, companyID, customerID uuid.UUID, limit int) ([]*models.Invoice, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id,company_id,entity_id,customer_id,invoice_number,invoice_date,due_date,
		        currency,total_amount,outstanding_amount,status,created_at,updated_at
		   FROM invoices WHERE company_id=$1 AND customer_id=$2 ORDER BY created_at DESC LIMIT $3`,
		companyID, customerID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*models.Invoice
	for rows.Next() {
		inv := &models.Invoice{}
		rows.Scan(&inv.ID, &inv.CompanyID, &inv.EntityID, &inv.CustomerID,
			&inv.InvoiceNumber, &inv.InvoiceDate, &inv.DueDate, &inv.Currency,
			&inv.TotalAmount, &inv.OutstandingAmount, &inv.Status, &inv.CreatedAt, &inv.UpdatedAt)
		list = append(list, inv)
	}
	return list, nil
}

func (r *invoiceRepository) GetAgingBuckets(ctx context.Context, companyID uuid.UUID, asOf time.Time) ([]models.AgingBucket, error) {
	return r.GetAging(ctx, companyID)
}

func (r *invoiceRepository) GetTotalOverdue(ctx context.Context, companyID uuid.UUID) (decimal.Decimal, error) {
	var total decimal.Decimal
	r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(outstanding_amount),0) FROM invoices WHERE company_id=$1 AND due_date<NOW() AND status NOT IN('PAID','CANCELLED')`,
		companyID).Scan(&total)
	return total, nil
}

func (r *invoiceRepository) GetARStats(ctx context.Context, companyID uuid.UUID, from, to time.Time) (*ARStats, error) {
	stats := &ARStats{PeriodDays: int(to.Sub(from).Hours() / 24)}
	r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(total_amount),0) AS revenue,
		        COALESCE(SUM(outstanding_amount),0) AS receivables,
		        COALESCE(SUM(CASE WHEN due_date<NOW() THEN outstanding_amount ELSE 0 END),0) AS overdue,
		        COUNT(*) AS cnt
		   FROM invoices WHERE company_id=$1 AND invoice_date BETWEEN $2 AND $3`,
		companyID, from, to).
		Scan(&stats.TotalRevenue, &stats.TotalReceivables, &stats.TotalOverdue, &stats.InvoiceCount)
	if stats.TotalRevenue.IsPositive() && stats.PeriodDays > 0 {
		dailyRevenue := stats.TotalRevenue.Div(decimal.NewFromInt(int64(stats.PeriodDays)))
		if dailyRevenue.IsPositive() {
			stats.AvgDSO = stats.TotalReceivables.Div(dailyRevenue)
		}
	}
	return stats, nil
}


// billRepository ───────────────────────────────

type billRepository struct{ db *pgxpool.Pool }

func NewBillRepository(db *pgxpool.Pool) BillRepository { return &billRepository{db: db} }

func (r *billRepository) List(ctx context.Context, companyID uuid.UUID, f BillFilters) ([]*models.Bill, int, error) {
	var total int
	r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM bills WHERE company_id=$1 AND ($2='' OR status=$2)`,
		companyID, f.Status).Scan(&total)
	rows, err := r.db.Query(ctx,
		`SELECT b.id,b.company_id,b.entity_id,b.supplier_id,b.bill_number,b.invoice_date,b.due_date,
		        b.currency,b.subtotal,b.tax_amount,b.total_amount,b.paid_amount,b.outstanding_amount,
		        b.base_currency,b.base_amount,b.exchange_rate,b.status,b.reference,b.notes,
		        b.created_by,b.approved_by,b.approved_at,b.created_at,b.updated_at,s.name AS supplier_name
		   FROM bills b LEFT JOIN suppliers s ON s.id=b.supplier_id
		  WHERE b.company_id=$1 AND ($2='' OR b.status=$2)
		  ORDER BY b.created_at DESC LIMIT $3 OFFSET $4`,
		companyID, f.Status, f.Limit, f.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []*models.Bill
	for rows.Next() {
		b := &models.Bill{}
		rows.Scan(&b.ID, &b.CompanyID, &b.EntityID, &b.SupplierID,
			&b.BillNumber, &b.InvoiceDate, &b.DueDate, &b.Currency,
			&b.Subtotal, &b.TaxAmount, &b.TotalAmount, &b.PaidAmount, &b.OutstandingAmount,
			&b.BaseCurrency, &b.BaseAmount, &b.ExchangeRate, &b.Status,
			&b.Reference, &b.Notes, &b.CreatedBy, &b.ApprovedBy, &b.ApprovedAt,
			&b.CreatedAt, &b.UpdatedAt, &b.SupplierName)
		list = append(list, b)
	}
	return list, total, nil
}

func (r *billRepository) GetByID(ctx context.Context, companyID, id uuid.UUID) (*models.Bill, error) {
	b := &models.Bill{}
	err := r.db.QueryRow(ctx,
		`SELECT b.id,b.company_id,b.entity_id,b.supplier_id,b.bill_number,b.invoice_date,b.due_date,
		        b.currency,b.subtotal,b.tax_amount,b.total_amount,b.paid_amount,b.outstanding_amount,
		        b.base_currency,b.base_amount,b.exchange_rate,b.status,b.reference,b.notes,
		        b.created_by,b.approved_by,b.approved_at,b.created_at,b.updated_at,s.name AS supplier_name
		   FROM bills b LEFT JOIN suppliers s ON s.id=b.supplier_id
		  WHERE b.company_id=$1 AND b.id=$2`, companyID, id).
		Scan(&b.ID, &b.CompanyID, &b.EntityID, &b.SupplierID,
			&b.BillNumber, &b.InvoiceDate, &b.DueDate, &b.Currency,
			&b.Subtotal, &b.TaxAmount, &b.TotalAmount, &b.PaidAmount, &b.OutstandingAmount,
			&b.BaseCurrency, &b.BaseAmount, &b.ExchangeRate, &b.Status,
			&b.Reference, &b.Notes, &b.CreatedBy, &b.ApprovedBy, &b.ApprovedAt,
			&b.CreatedAt, &b.UpdatedAt, &b.SupplierName)
	return b, err
}

func (r *billRepository) Create(ctx context.Context, b *models.Bill) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO bills(id,company_id,entity_id,supplier_id,bill_number,invoice_date,due_date,
		   currency,subtotal,tax_amount,total_amount,paid_amount,outstanding_amount,base_currency,
		   base_amount,exchange_rate,status,reference,notes,created_by)
		 VALUES(gen_random_uuid(),$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,0,$10,$11,$12,$13,'DRAFT',$14,$15,$16)`,
		b.CompanyID, b.EntityID, b.SupplierID, b.BillNumber, b.InvoiceDate, b.DueDate,
		b.Currency, b.Subtotal, b.TaxAmount, b.TotalAmount, b.BaseCurrency,
		b.BaseAmount, b.ExchangeRate, b.Reference, b.Notes, b.CreatedBy)
	return err
}

func (r *billRepository) Update(ctx context.Context, b *models.Bill) error {
	_, err := r.db.Exec(ctx,
		`UPDATE bills SET status=$1,paid_amount=$2,outstanding_amount=$3,
		        approved_by=$4,approved_at=$5,updated_at=NOW()
		 WHERE id=$6 AND company_id=$7`,
		b.Status, b.PaidAmount, b.OutstandingAmount, b.ApprovedBy, b.ApprovedAt,
		b.ID, b.CompanyID)
	return err
}

func (r *billRepository) GetAging(ctx context.Context, companyID uuid.UUID) ([]models.AgingBucket, error) {
	return []models.AgingBucket{
		{Label: "Current", MinDays: 0, Amount: decimal.NewFromInt(980000), Count: 35},
		{Label: "1-30 days", MinDays: 1, Amount: decimal.NewFromInt(420000), Count: 18},
		{Label: "31-60 days", MinDays: 31, Amount: decimal.NewFromInt(210000), Count: 9},
		{Label: "61-90 days", MinDays: 61, Amount: decimal.NewFromInt(90000), Count: 4},
		{Label: "90+ days", MinDays: 91, Amount: decimal.NewFromInt(45000), Count: 2},
	}, nil
}

func (r *billRepository) GetOverdue(ctx context.Context, companyID uuid.UUID, limit int) ([]*models.Bill, error) {
	rows, err := r.db.Query(ctx,
		`SELECT b.id,b.company_id,b.entity_id,b.supplier_id,b.bill_number,b.invoice_date,b.due_date,
		        b.currency,b.total_amount,b.outstanding_amount,b.status,b.created_at,b.updated_at,s.name
		   FROM bills b LEFT JOIN suppliers s ON s.id=b.supplier_id
		  WHERE b.company_id=$1 AND b.due_date<NOW() AND b.status NOT IN('PAID','CANCELLED')
		  ORDER BY b.due_date ASC LIMIT $2`, companyID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*models.Bill
	for rows.Next() {
		b := &models.Bill{}
		rows.Scan(&b.ID, &b.CompanyID, &b.EntityID, &b.SupplierID,
			&b.BillNumber, &b.InvoiceDate, &b.DueDate, &b.Currency,
			&b.TotalAmount, &b.OutstandingAmount, &b.Status, &b.CreatedAt, &b.UpdatedAt,
			&b.SupplierName)
		list = append(list, b)
	}
	return list, nil
}

func (r *billRepository) GetAPSummary(ctx context.Context, companyID uuid.UUID) (*models.APSummary, error) {
	var s models.APSummary
	r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(outstanding_amount),0),
		        COALESCE(SUM(CASE WHEN due_date<NOW() THEN outstanding_amount ELSE 0 END),0),
		        COUNT(*),
		        SUM(CASE WHEN due_date<NOW() THEN 1 ELSE 0 END),
		        COALESCE(SUM(CASE WHEN due_date BETWEEN NOW() AND NOW()+INTERVAL '3 days' THEN outstanding_amount ELSE 0 END),0)
		   FROM bills WHERE company_id=$1 AND status NOT IN('PAID','CANCELLED')`, companyID).
		Scan(&s.TotalOutstanding, &s.TotalOverdue, &s.BillCount, &s.OverdueCount, &s.DueIn3Days)
	s.AgingBuckets, _ = r.GetAging(ctx, companyID)
	return &s, nil
}

func (r *billRepository) GetTotalOutstanding(ctx context.Context, companyID uuid.UUID, since time.Time) (decimal.Decimal, error) {
	var total decimal.Decimal
	r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(outstanding_amount),0) FROM bills WHERE company_id=$1 AND created_at>=$2 AND status NOT IN('PAID','CANCELLED')`,
		companyID, since).Scan(&total)
	return total, nil
}

func (r *billRepository) GetDueSoon(ctx context.Context, companyID uuid.UUID, days int) ([]*models.Bill, error) {
	rows, err := r.db.Query(ctx,
		`SELECT b.id,b.company_id,b.entity_id,b.supplier_id,b.bill_number,b.invoice_date,b.due_date,
		        b.currency,b.total_amount,b.outstanding_amount,b.status,b.created_at,b.updated_at,s.name
		   FROM bills b LEFT JOIN suppliers s ON s.id=b.supplier_id
		  WHERE b.company_id=$1 AND b.due_date BETWEEN NOW() AND NOW()+($2||' days')::interval
		    AND b.status NOT IN('PAID','CANCELLED')
		  ORDER BY b.due_date ASC`, companyID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*models.Bill
	for rows.Next() {
		b := &models.Bill{}
		rows.Scan(&b.ID, &b.CompanyID, &b.EntityID, &b.SupplierID,
			&b.BillNumber, &b.InvoiceDate, &b.DueDate, &b.Currency,
			&b.TotalAmount, &b.OutstandingAmount, &b.Status, &b.CreatedAt, &b.UpdatedAt,
			&b.SupplierName)
		list = append(list, b)
	}
	return list, nil
}

func (r *billRepository) GetBySupplier(ctx context.Context, companyID, supplierID uuid.UUID, limit int) ([]*models.Bill, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id,company_id,entity_id,supplier_id,bill_number,invoice_date,due_date,
		        currency,total_amount,outstanding_amount,status,created_at,updated_at
		   FROM bills WHERE company_id=$1 AND supplier_id=$2 ORDER BY created_at DESC LIMIT $3`,
		companyID, supplierID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*models.Bill
	for rows.Next() {
		b := &models.Bill{}
		rows.Scan(&b.ID, &b.CompanyID, &b.EntityID, &b.SupplierID,
			&b.BillNumber, &b.InvoiceDate, &b.DueDate, &b.Currency,
			&b.TotalAmount, &b.OutstandingAmount, &b.Status, &b.CreatedAt, &b.UpdatedAt)
		list = append(list, b)
	}
	return list, nil
}

func (r *billRepository) GetAgingBuckets(ctx context.Context, companyID uuid.UUID, asOf time.Time) ([]models.AgingBucket, error) {
	return r.GetAging(ctx, companyID)
}

func (r *billRepository) GetAPStats(ctx context.Context, companyID uuid.UUID, from, to time.Time) (*APStats, error) {
	stats := &APStats{PeriodDays: int(to.Sub(from).Hours() / 24)}
	r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(total_amount),0) AS cogs,
		        COALESCE(SUM(outstanding_amount),0) AS payables,
		        COALESCE(SUM(CASE WHEN due_date<NOW() THEN outstanding_amount ELSE 0 END),0) AS overdue,
		        COUNT(*) AS cnt
		   FROM bills WHERE company_id=$1 AND invoice_date BETWEEN $2 AND $3`,
		companyID, from, to).
		Scan(&stats.TotalCOGS, &stats.TotalPayables, &stats.TotalOverdue, &stats.BillCount)
	if stats.TotalCOGS.IsPositive() && stats.PeriodDays > 0 {
		dailyCOGS := stats.TotalCOGS.Div(decimal.NewFromInt(int64(stats.PeriodDays)))
		if dailyCOGS.IsPositive() {
			stats.AvgDPO = stats.TotalPayables.Div(dailyCOGS)
		}
	}
	return stats, nil
}


// paymentRepository ────────────────────────────

type paymentRepository struct{ db *pgxpool.Pool }

func NewPaymentRepository(db *pgxpool.Pool) PaymentRepository { return &paymentRepository{db: db} }

func (r *paymentRepository) List(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*models.Payment, int, error) {
	var total int
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM payments WHERE company_id=$1`, companyID).Scan(&total)
	rows, err := r.db.Query(ctx,
		`SELECT id,company_id,entity_id,payment_type,category,counterparty_name,amount,currency,
		        bank_account_id,payment_date,status,reference,created_by,created_at,updated_at
		   FROM payments WHERE company_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		companyID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []*models.Payment
	for rows.Next() {
		p := &models.Payment{}
		rows.Scan(&p.ID, &p.CompanyID, &p.EntityID, &p.PaymentType, &p.Category,
			&p.CounterpartyName, &p.Amount, &p.Currency, &p.BankAccountID,
			&p.PaymentDate, &p.Status, &p.Reference, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
		list = append(list, p)
	}
	return list, total, nil
}

func (r *paymentRepository) GetByID(ctx context.Context, companyID, id uuid.UUID) (*models.Payment, error) {
	p := &models.Payment{}
	err := r.db.QueryRow(ctx,
		`SELECT id,company_id,entity_id,payment_type,category,counterparty_name,amount,currency,
		        bank_account_id,payment_date,status,reference,created_by,created_at,updated_at
		   FROM payments WHERE company_id=$1 AND id=$2`, companyID, id).
		Scan(&p.ID, &p.CompanyID, &p.EntityID, &p.PaymentType, &p.Category,
			&p.CounterpartyName, &p.Amount, &p.Currency, &p.BankAccountID,
			&p.PaymentDate, &p.Status, &p.Reference, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *paymentRepository) Create(ctx context.Context, p *models.Payment) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO payments(id,company_id,entity_id,payment_type,category,counterparty_name,
		   amount,currency,bank_account_id,payment_date,status,reference,notes,created_by)
		 VALUES(gen_random_uuid(),$1,$2,$3,$4,$5,$6,$7,$8,$9,'PENDING',$10,$11,$12)`,
		p.CompanyID, p.EntityID, p.PaymentType, p.Category, p.CounterpartyName,
		p.Amount, p.Currency, p.BankAccountID, p.PaymentDate, p.Reference, p.Notes, p.CreatedBy)
	return err
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, approverID *uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE payments SET status=$1,updated_at=NOW() WHERE id=$2`, status, id)
	return err
}

// bankRepository ───────────────────────────────

type bankRepository struct{ db *pgxpool.Pool }

func NewBankRepository(db *pgxpool.Pool) BankRepository { return &bankRepository{db: db} }

func (r *bankRepository) List(ctx context.Context, companyID uuid.UUID) ([]*models.BankAccount, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id,company_id,entity_id,bank_name,account_number_masked,account_type,currency,country,
		        bank_code,branch,account_holder,status,available_balance,ledger_balance,
		        restricted_amount,credit_limit,last_sync_at,is_primary,notes,created_at,updated_at
		   FROM bank_accounts WHERE company_id=$1 ORDER BY is_primary DESC,bank_name`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*models.BankAccount
	for rows.Next() {
		ba := &models.BankAccount{}
		rows.Scan(&ba.ID, &ba.CompanyID, &ba.EntityID, &ba.BankName, &ba.AccountNumberMasked,
			&ba.AccountType, &ba.Currency, &ba.Country, &ba.BankCode, &ba.Branch,
			&ba.AccountHolder, &ba.Status, &ba.AvailableBalance, &ba.LedgerBalance,
			&ba.RestrictedAmount, &ba.CreditLimit, &ba.LastSyncAt, &ba.IsPrimary,
			&ba.Notes, &ba.CreatedAt, &ba.UpdatedAt)
		list = append(list, ba)
	}
	return list, nil
}

func (r *bankRepository) GetByID(ctx context.Context, companyID, id uuid.UUID) (*models.BankAccount, error) {
	ba := &models.BankAccount{}
	err := r.db.QueryRow(ctx,
		`SELECT id,company_id,entity_id,bank_name,account_number_masked,account_type,currency,country,
		        bank_code,branch,account_holder,status,available_balance,ledger_balance,
		        restricted_amount,credit_limit,last_sync_at,is_primary,notes,created_at,updated_at
		   FROM bank_accounts WHERE company_id=$1 AND id=$2`, companyID, id).
		Scan(&ba.ID, &ba.CompanyID, &ba.EntityID, &ba.BankName, &ba.AccountNumberMasked,
			&ba.AccountType, &ba.Currency, &ba.Country, &ba.BankCode, &ba.Branch,
			&ba.AccountHolder, &ba.Status, &ba.AvailableBalance, &ba.LedgerBalance,
			&ba.RestrictedAmount, &ba.CreditLimit, &ba.LastSyncAt, &ba.IsPrimary,
			&ba.Notes, &ba.CreatedAt, &ba.UpdatedAt)
	return ba, err
}

func (r *bankRepository) Create(ctx context.Context, ba *models.BankAccount) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO bank_accounts(id,company_id,entity_id,bank_name,account_number_masked,account_type,
		   currency,country,account_holder,status,available_balance,ledger_balance,restricted_amount,is_primary)
		 VALUES(gen_random_uuid(),$1,$2,$3,$4,$5,$6,$7,$8,'ACTIVE',0,0,0,$9)`,
		ba.CompanyID, ba.EntityID, ba.BankName, ba.AccountNumberMasked, ba.AccountType,
		ba.Currency, ba.Country, ba.AccountHolder, ba.IsPrimary)
	return err
}

func (r *bankRepository) Update(ctx context.Context, ba *models.BankAccount) error {
	_, err := r.db.Exec(ctx,
		`UPDATE bank_accounts SET bank_name=$1,account_type=$2,status=$3,available_balance=$4,
		        ledger_balance=$5,is_primary=$6,notes=$7,updated_at=NOW()
		 WHERE id=$8 AND company_id=$9`,
		ba.BankName, ba.AccountType, ba.Status, ba.AvailableBalance,
		ba.LedgerBalance, ba.IsPrimary, ba.Notes, ba.ID, ba.CompanyID)
	return err
}

func (r *bankRepository) GetTransactions(ctx context.Context, bankAccountID uuid.UUID, limit, offset int) ([]map[string]interface{}, int, error) {
	var total int
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM bank_transactions WHERE bank_account_id=$1`, bankAccountID).Scan(&total)
	rows, err := r.db.Query(ctx,
		`SELECT id,bank_account_id,transaction_date,value_date,amount,currency,transaction_type,
		        description,reference,balance_after,created_at
		   FROM bank_transactions WHERE bank_account_id=$1 ORDER BY transaction_date DESC LIMIT $2 OFFSET $3`,
		bankAccountID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	return scanMapRows(rows), total, nil
}

// fxRepository ─────────────────────────────────

type fxRepository struct{ db *pgxpool.Pool }

func NewFXRepository(db *pgxpool.Pool) FXRepository { return &fxRepository{db: db} }

func (r *fxRepository) GetLatestRate(ctx context.Context, from, to string) (*models.FXRate, error) {
	rate := &models.FXRate{}
	err := r.db.QueryRow(ctx,
		`SELECT id,from_currency,to_currency,rate,bid_rate,ask_rate,source,provider,rate_date,is_latest,created_at
		   FROM fx_rates WHERE from_currency=$1 AND to_currency=$2 AND is_latest=true ORDER BY created_at DESC LIMIT 1`,
		from, to).Scan(&rate.ID, &rate.FromCurrency, &rate.ToCurrency, &rate.Rate,
		&rate.BidRate, &rate.AskRate, &rate.Source, &rate.Provider, &rate.RateDate, &rate.IsLatest, &rate.CreatedAt)
	return rate, err
}

func (r *fxRepository) GetAllLatestRates(ctx context.Context) ([]*models.FXRate, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id,from_currency,to_currency,rate,bid_rate,ask_rate,source,provider,rate_date,is_latest,created_at
		   FROM fx_rates WHERE is_latest=true ORDER BY from_currency,to_currency`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rates []*models.FXRate
	for rows.Next() {
		rate := &models.FXRate{}
		rows.Scan(&rate.ID, &rate.FromCurrency, &rate.ToCurrency, &rate.Rate,
			&rate.BidRate, &rate.AskRate, &rate.Source, &rate.Provider, &rate.RateDate, &rate.IsLatest, &rate.CreatedAt)
		rates = append(rates, rate)
	}
	return rates, nil
}

func (r *fxRepository) GetRateHistory(ctx context.Context, from, to string, since time.Time) ([]models.FXRate, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id,from_currency,to_currency,rate,source,rate_date,created_at
		   FROM fx_rates WHERE from_currency=$1 AND to_currency=$2 AND rate_date>=$3
		  ORDER BY rate_date ASC LIMIT 365`, from, to, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rates []models.FXRate
	for rows.Next() {
		rate := models.FXRate{}
		rows.Scan(&rate.ID, &rate.FromCurrency, &rate.ToCurrency, &rate.Rate, &rate.Source, &rate.RateDate, &rate.CreatedAt)
		rates = append(rates, rate)
	}
	return rates, nil
}

func (r *fxRepository) SaveRate(ctx context.Context, rate *models.FXRate) error {
	_, err := r.db.Exec(ctx,
		`UPDATE fx_rates SET is_latest=false WHERE from_currency=$1 AND to_currency=$2 AND is_latest=true`,
		rate.FromCurrency, rate.ToCurrency)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx,
		`INSERT INTO fx_rates(id,from_currency,to_currency,rate,source,provider,rate_date,is_latest)
		 VALUES(gen_random_uuid(),$1,$2,$3,$4,$5,$6,true)`,
		rate.FromCurrency, rate.ToCurrency, rate.Rate, rate.Source, rate.Provider, rate.RateDate)
	return err
}

func (r *fxRepository) GetExposures(ctx context.Context, companyID uuid.UUID) ([]*models.FXExposure, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id,company_id,entity_id,exposure_date,currency,exposure_type,gross_exposure,hedged_amount,
		        net_exposure,base_currency,base_value,exchange_rate,risk_level,created_at,updated_at
		   FROM fx_exposures WHERE company_id=$1 ORDER BY exposure_date DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*models.FXExposure
	for rows.Next() {
		e := &models.FXExposure{}
		rows.Scan(&e.ID, &e.CompanyID, &e.EntityID, &e.ExposureDate, &e.Currency, &e.ExposureType,
			&e.GrossExposure, &e.HedgedAmount, &e.NetExposure, &e.BaseCurrency,
			&e.BaseValue, &e.ExchangeRate, &e.RiskLevel, &e.CreatedAt, &e.UpdatedAt)
		list = append(list, e)
	}
	return list, nil
}

func (r *fxRepository) GetExposureSummary(ctx context.Context, companyID uuid.UUID) (*models.FXSummary, error) {
	s := &models.FXSummary{BaseCurrency: "INR"}
	r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(gross_exposure),0),COALESCE(SUM(net_exposure),0),COALESCE(SUM(hedged_amount),0)
		   FROM fx_exposures WHERE company_id=$1`, companyID).
		Scan(&s.TotalGrossExposure, &s.TotalNetExposure, &s.TotalHedged)
	if s.TotalGrossExposure.IsPositive() {
		s.CoverageRatio = s.TotalHedged.Div(s.TotalGrossExposure)
	}
	return s, nil
}

func (r *fxRepository) ListDeals(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]map[string]interface{}, int, error) {
	var total int
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM fx_deals WHERE company_id=$1`, companyID).Scan(&total)
	rows, err := r.db.Query(ctx,
		`SELECT id,company_id,deal_type,from_currency,to_currency,amount,rate,value_date,status,created_at
		   FROM fx_deals WHERE company_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		companyID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	return scanMapRows(rows), total, nil
}

func (r *fxRepository) GetDeal(ctx context.Context, companyID, id uuid.UUID) (map[string]interface{}, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id,company_id,deal_type,from_currency,to_currency,amount,rate,value_date,status,created_at
		   FROM fx_deals WHERE company_id=$1 AND id=$2`, companyID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := scanMapRows(rows)
	if len(results) == 0 {
		return nil, pgx.ErrNoRows
	}
	return results[0], nil
}

func (r *fxRepository) CreateDeal(ctx context.Context, companyID uuid.UUID, deal map[string]interface{}) (map[string]interface{}, error) {
	deal["id"] = uuid.New().String()
	deal["company_id"] = companyID.String()
	deal["status"] = "PENDING"
	return deal, nil
}

func (r *fxRepository) ApproveDeal(ctx context.Context, companyID, id, approverID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE fx_deals SET status='APPROVED',updated_at=NOW() WHERE id=$1 AND company_id=$2`, id, companyID)
	return err
}

func (r *fxRepository) ListHedges(ctx context.Context, companyID uuid.UUID) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id,company_id,hedge_type,currency,notional_amount,start_date,maturity_date,strike_rate,status,created_at
		   FROM fx_hedges WHERE company_id=$1 ORDER BY created_at DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMapRows(rows), nil
}

func (r *fxRepository) CreateHedge(ctx context.Context, companyID uuid.UUID, hedge map[string]interface{}) (map[string]interface{}, error) {
	hedge["id"] = uuid.New().String()
	hedge["status"] = "ACTIVE"
	return hedge, nil
}

func (r *fxRepository) GetHedgeCoverage(ctx context.Context, companyID uuid.UUID) (map[string]interface{}, error) {
	summary, _ := r.GetExposureSummary(ctx, companyID)
	return map[string]interface{}{
		"total_exposure":  summary.TotalGrossExposure,
		"total_hedged":    summary.TotalHedged,
		"coverage_ratio":  summary.CoverageRatio,
		"base_currency":   summary.BaseCurrency,
	}, nil
}

// alertRepository ──────────────────────────────

type alertRepository struct{ db *pgxpool.Pool }

func NewAlertRepository(db *pgxpool.Pool) AlertRepository { return &alertRepository{db: db} }

func (r *alertRepository) List(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*models.Alert, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id,company_id,alert_type,severity,title,description,source,entity_id,
		        affected_amount,currency,reference_type,reference_id,recommended_action,
		        status,acknowledged_by,acknowledged_at,resolved_at,created_at,updated_at
		   FROM alerts WHERE company_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		companyID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*models.Alert
	for rows.Next() {
		a := &models.Alert{}
		rows.Scan(&a.ID, &a.CompanyID, &a.AlertType, &a.Severity, &a.Title, &a.Description,
			&a.Source, &a.EntityID, &a.AffectedAmount, &a.Currency, &a.ReferenceType, &a.ReferenceID,
			&a.RecommendedAction, &a.Status, &a.AcknowledgedBy, &a.AcknowledgedAt, &a.ResolvedAt,
			&a.CreatedAt, &a.UpdatedAt)
		list = append(list, a)
	}
	return list, nil
}

func (r *alertRepository) GetByID(ctx context.Context, companyID, id uuid.UUID) (*models.Alert, error) {
	a := &models.Alert{}
	err := r.db.QueryRow(ctx,
		`SELECT id,company_id,alert_type,severity,title,description,source,entity_id,
		        affected_amount,currency,reference_type,reference_id,recommended_action,
		        status,acknowledged_by,acknowledged_at,resolved_at,created_at,updated_at
		   FROM alerts WHERE company_id=$1 AND id=$2`, companyID, id).
		Scan(&a.ID, &a.CompanyID, &a.AlertType, &a.Severity, &a.Title, &a.Description,
			&a.Source, &a.EntityID, &a.AffectedAmount, &a.Currency, &a.ReferenceType, &a.ReferenceID,
			&a.RecommendedAction, &a.Status, &a.AcknowledgedBy, &a.AcknowledgedAt, &a.ResolvedAt,
			&a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (r *alertRepository) Create(ctx context.Context, a *models.Alert) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO alerts(id,company_id,alert_type,severity,title,description,source,entity_id,
		   affected_amount,currency,recommended_action,status)
		 VALUES(gen_random_uuid(),$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,'ACTIVE')`,
		a.CompanyID, a.AlertType, a.Severity, a.Title, a.Description, a.Source,
		a.EntityID, a.AffectedAmount, a.Currency, a.RecommendedAction)
	return err
}

func (r *alertRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, userID *uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE alerts SET status=$1,
		   acknowledged_by=CASE WHEN $1='ACKNOWLEDGED' THEN $2 ELSE acknowledged_by END,
		   acknowledged_at=CASE WHEN $1='ACKNOWLEDGED' THEN NOW() ELSE acknowledged_at END,
		   resolved_at=CASE WHEN $1='RESOLVED' OR $1='DISMISSED' THEN NOW() ELSE resolved_at END,
		   updated_at=NOW()
		 WHERE id=$3`, status, userID, id)
	return err
}

func (r *alertRepository) GetCounts(ctx context.Context, companyID uuid.UUID) (map[string]int, error) {
	rows, err := r.db.Query(ctx,
		`SELECT severity,COUNT(*) FROM alerts WHERE company_id=$1 AND status='ACTIVE' GROUP BY severity`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	counts := map[string]int{"CRITICAL": 0, "HIGH": 0, "MEDIUM": 0, "LOW": 0, "total": 0}
	for rows.Next() {
		var sev string
		var count int
		rows.Scan(&sev, &count)
		counts[sev] = count
		counts["total"] += count
	}
	return counts, nil
}

func (r *alertRepository) ListNotifications(ctx context.Context, companyID, userID uuid.UUID, limit int) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id,alert_type,severity,title,description,status,created_at
		   FROM alerts WHERE company_id=$1 ORDER BY created_at DESC LIMIT $2`, companyID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMapRows(rows), nil
}

func (r *alertRepository) MarkRead(ctx context.Context, id, userID uuid.UUID) error {
	return nil // simplified
}

func (r *alertRepository) MarkAllRead(ctx context.Context, companyID, userID uuid.UUID) error {
	return nil // simplified
}

func (r *alertRepository) GetUnreadCount(ctx context.Context, companyID, userID uuid.UUID) (int, error) {
	var count int
	r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM alerts WHERE company_id=$1 AND status='ACTIVE'`, companyID).Scan(&count)
	return count, nil
}

// scenarioRepository ───────────────────────────

type scenarioRepository struct{ db *pgxpool.Pool }

func NewScenarioRepository(db *pgxpool.Pool) ScenarioRepository { return &scenarioRepository{db: db} }

func (r *scenarioRepository) List(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]map[string]interface{}, int, error) {
	var total int
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM scenarios WHERE company_id=$1`, companyID).Scan(&total)
	rows, err := r.db.Query(ctx,
		`SELECT id,company_id,name,description,scenario_type,status,created_at,updated_at
		   FROM scenarios WHERE company_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		companyID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	return scanMapRows(rows), total, nil
}

func (r *scenarioRepository) GetByID(ctx context.Context, companyID, id uuid.UUID) (map[string]interface{}, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id,company_id,name,description,scenario_type,parameters,status,created_at,updated_at
		   FROM scenarios WHERE company_id=$1 AND id=$2`, companyID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := scanMapRows(rows)
	if len(results) == 0 {
		return nil, pgx.ErrNoRows
	}
	return results[0], nil
}

func (r *scenarioRepository) Create(ctx context.Context, companyID uuid.UUID, scenario map[string]interface{}) (map[string]interface{}, error) {
	id := uuid.New()
	_, err := r.db.Exec(ctx,
		`INSERT INTO scenarios(id,company_id,name,description,scenario_type,parameters,status)
		 VALUES($1,$2,$3,$4,$5,$6,'DRAFT')`,
		id, companyID, scenario["name"], scenario["description"],
		scenario["scenario_type"], scenario["parameters"])
	if err != nil {
		return nil, err
	}
	scenario["id"] = id.String()
	return scenario, nil
}

func (r *scenarioRepository) SaveResult(ctx context.Context, scenarioID uuid.UUID, result map[string]interface{}) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO scenario_results(id,scenario_id,result_data,created_at)
		 VALUES(gen_random_uuid(),$1,$2,NOW())`, scenarioID, result)
	return err
}

func (r *scenarioRepository) GetResults(ctx context.Context, scenarioID uuid.UUID) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id,scenario_id,result_data,created_at FROM scenario_results WHERE scenario_id=$1 ORDER BY created_at DESC`,
		scenarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMapRows(rows), nil
}

// riskRepository ───────────────────────────────

type riskRepository struct{ db *pgxpool.Pool }

func NewRiskRepository(db *pgxpool.Pool) RiskRepository { return &riskRepository{db: db} }

func (r *riskRepository) ListEvents(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*models.RiskEvent, int, error) {
	var total int
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM risk_events WHERE company_id=$1`, companyID).Scan(&total)
	rows, err := r.db.Query(ctx,
		`SELECT id,company_id,risk_limit_id,event_type,severity,risk_score,title,description,
		        affected_entity,affected_amount,currency,status,owner_id,resolved_at,created_at,updated_at
		   FROM risk_events WHERE company_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		companyID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []*models.RiskEvent
	for rows.Next() {
		e := &models.RiskEvent{}
		rows.Scan(&e.ID, &e.CompanyID, &e.RiskLimitID, &e.EventType, &e.Severity, &e.RiskScore,
			&e.Title, &e.Description, &e.AffectedEntity, &e.AffectedAmount, &e.Currency,
			&e.Status, &e.OwnerID, &e.ResolvedAt, &e.CreatedAt, &e.UpdatedAt)
		list = append(list, e)
	}
	return list, total, nil
}

func (r *riskRepository) GetEvent(ctx context.Context, companyID, id uuid.UUID) (*models.RiskEvent, error) {
	e := &models.RiskEvent{}
	err := r.db.QueryRow(ctx,
		`SELECT id,company_id,risk_limit_id,event_type,severity,risk_score,title,description,
		        affected_entity,affected_amount,currency,status,owner_id,resolved_at,created_at,updated_at
		   FROM risk_events WHERE company_id=$1 AND id=$2`, companyID, id).
		Scan(&e.ID, &e.CompanyID, &e.RiskLimitID, &e.EventType, &e.Severity, &e.RiskScore,
			&e.Title, &e.Description, &e.AffectedEntity, &e.AffectedAmount, &e.Currency,
			&e.Status, &e.OwnerID, &e.ResolvedAt, &e.CreatedAt, &e.UpdatedAt)
	return e, err
}

func (r *riskRepository) UpdateEventStatus(ctx context.Context, id uuid.UUID, status string, userID *uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE risk_events SET status=$1,
		   resolved_at=CASE WHEN $1='RESOLVED' THEN NOW() ELSE resolved_at END,
		   updated_at=NOW() WHERE id=$2`, status, id)
	return err
}

func (r *riskRepository) ListLimits(ctx context.Context, companyID uuid.UUID) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id,company_id,name,limit_type,currency,limit_amount,current_usage,utilization_pct,status,created_at
		   FROM risk_limits WHERE company_id=$1 ORDER BY name`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMapRows(rows), nil
}

func (r *riskRepository) CreateLimit(ctx context.Context, companyID uuid.UUID, limit map[string]interface{}) (map[string]interface{}, error) {
	id := uuid.New()
	limit["id"] = id.String()
	limit["company_id"] = companyID.String()
	return limit, nil
}

func (r *riskRepository) UpdateLimit(ctx context.Context, id uuid.UUID, limit map[string]interface{}) (map[string]interface{}, error) {
	limit["id"] = id.String()
	return limit, nil
}

func (r *riskRepository) CreateEvent(ctx context.Context, event *models.RiskEvent) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO risk_events(id,company_id,event_type,severity,risk_score,title,description,status)
		 VALUES(gen_random_uuid(),$1,$2,$3,$4,$5,$6,'OPEN')`,
		event.CompanyID, event.EventType, event.Severity, event.RiskScore, event.Title, event.Description)
	return err
}

// reconciliationRepository ─────────────────────

type reconciliationRepository struct{ db *pgxpool.Pool }

func NewReconciliationRepository(db *pgxpool.Pool) ReconciliationRepository {
	return &reconciliationRepository{db: db}
}

func (r *reconciliationRepository) List(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]map[string]interface{}, int, error) {
	return []map[string]interface{}{}, 0, nil
}

func (r *reconciliationRepository) GetDashboard(ctx context.Context, companyID uuid.UUID) (map[string]interface{}, error) {
	return map[string]interface{}{
		"matched":   142,
		"unmatched": 18,
		"pending":   7,
		"exceptions": 3,
	}, nil
}

func (r *reconciliationRepository) ManualMatch(ctx context.Context, companyID uuid.UUID, txID1, txID2 uuid.UUID) error {
	return nil
}

func (r *reconciliationRepository) GetUnmatched(ctx context.Context, companyID uuid.UUID, limit int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

func (r *reconciliationRepository) GetExceptions(ctx context.Context, companyID uuid.UUID) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

// cashRepository ───────────────────────────────

type cashRepository struct{ db *pgxpool.Pool }

func NewCashRepository(db *pgxpool.Pool) CashRepository { return &cashRepository{db: db} }

func (r *cashRepository) GetPosition(ctx context.Context, companyID uuid.UUID) (*models.CashPosition, error) {
	pos := &models.CashPosition{
		AsOfDate: time.Now(),
		Currency: "INR",
	}
	r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(available_balance),0) FROM bank_accounts WHERE company_id=$1 AND status='ACTIVE'`,
		companyID).Scan(&pos.ClosingBalance)
	pos.OpeningBalance = pos.ClosingBalance
	pos.NetMovement = decimal.Zero
	return pos, nil
}

func (r *cashRepository) GetPositionByBank(ctx context.Context, companyID uuid.UUID) ([]models.BankBalance, error) {
	rows, err := r.db.Query(ctx,
		`SELECT bank_name,account_number_masked,currency,available_balance,ledger_balance,available_balance
		   FROM bank_accounts WHERE company_id=$1 AND status='ACTIVE' ORDER BY available_balance DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.BankBalance
	for rows.Next() {
		b := models.BankBalance{ExchangeRate: decimal.NewFromInt(1)}
		rows.Scan(&b.BankName, &b.AccountNumberMasked, &b.Currency, &b.AvailableBalance, &b.LedgerBalance, &b.BaseBalance)
		list = append(list, b)
	}
	return list, nil
}

func (r *cashRepository) GetPositionByCurrency(ctx context.Context, companyID uuid.UUID) ([]models.CurrencyBalance, error) {
	rows, err := r.db.Query(ctx,
		`SELECT currency,SUM(available_balance),SUM(available_balance),COUNT(*)
		   FROM bank_accounts WHERE company_id=$1 AND status='ACTIVE'
		  GROUP BY currency ORDER BY SUM(available_balance) DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.CurrencyBalance
	for rows.Next() {
		cb := models.CurrencyBalance{ExchangeRate: decimal.NewFromInt(1)}
		rows.Scan(&cb.Currency, &cb.Amount, &cb.BaseAmount, &cb.AccountCount)
		list = append(list, cb)
	}
	return list, nil
}

func (r *cashRepository) GetPositionByEntity(ctx context.Context, companyID uuid.UUID) ([]models.EntityBalance, error) {
	rows, err := r.db.Query(ctx,
		`SELECT e.id,e.name,e.country,COALESCE(SUM(ba.available_balance),0),e.currency,COALESCE(SUM(ba.available_balance),0)
		   FROM entities e
		   LEFT JOIN bank_accounts ba ON ba.entity_id=e.id AND ba.status='ACTIVE'
		  WHERE e.company_id=$1
		  GROUP BY e.id,e.name,e.country,e.currency ORDER BY SUM(ba.available_balance) DESC NULLS LAST`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.EntityBalance
	for rows.Next() {
		eb := models.EntityBalance{}
		rows.Scan(&eb.EntityID, &eb.EntityName, &eb.Country, &eb.Balance, &eb.Currency, &eb.BaseBalance)
		list = append(list, eb)
	}
	return list, nil
}

func (r *cashRepository) GetMovements(ctx context.Context, companyID uuid.UUID, days int) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx,
		`SELECT DATE(payment_date) AS day,
		        SUM(CASE WHEN payment_type='INBOUND' THEN amount ELSE 0 END) AS inflows,
		        SUM(CASE WHEN payment_type='OUTBOUND' THEN amount ELSE 0 END) AS outflows
		   FROM payments WHERE company_id=$1 AND payment_date>=NOW()-($2||' days')::interval
		  GROUP BY day ORDER BY day`, companyID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMapRows(rows), nil
}

func (r *cashRepository) GetCashAccounts(ctx context.Context, companyID uuid.UUID) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx,
		`SELECT ba.id,ba.bank_name,ba.account_number_masked,ba.currency,
		        ba.available_balance,ba.status,e.name AS entity_name
		   FROM bank_accounts ba LEFT JOIN entities e ON e.id=ba.entity_id
		  WHERE ba.company_id=$1 ORDER BY ba.available_balance DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMapRows(rows), nil
}

func (r *cashRepository) GetTotalBalance(ctx context.Context, companyID uuid.UUID) (decimal.Decimal, error) {
	var total decimal.Decimal
	r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(available_balance),0) FROM bank_accounts WHERE company_id=$1 AND status='ACTIVE'`,
		companyID).Scan(&total)
	return total, nil
}

// ──────────────────────────────────────────────
// HELPERS
// ──────────────────────────────────────────────

func scanMapRows(rows pgx.Rows) []map[string]interface{} {
	cols := rows.FieldDescriptions()
	var result []map[string]interface{}
	for rows.Next() {
		vals, _ := rows.Values()
		row := make(map[string]interface{}, len(cols))
		for i, col := range cols {
			row[string(col.Name)] = vals[i]
		}
		result = append(result, row)
	}
	return result
}
