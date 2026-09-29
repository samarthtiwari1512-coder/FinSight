-- FinSight Initial Schema
-- Version: 001
-- Description: Complete normalized schema for Corporate Treasury Platform
-- Uses UUID PKs, DECIMAL for money, proper FKs and indexes

-- Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- ============================================================
-- COMPANIES & ENTITIES
-- ============================================================

CREATE TABLE companies (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            VARCHAR(200) NOT NULL,
    code            VARCHAR(20) UNIQUE NOT NULL,
    base_currency   CHAR(3) NOT NULL DEFAULT 'INR',
    country         VARCHAR(100) NOT NULL,
    fiscal_year_start INT NOT NULL DEFAULT 4, -- Month (1-12)
    timezone        VARCHAR(50) NOT NULL DEFAULT 'Asia/Kolkata',
    address         TEXT,
    logo_url        VARCHAR(500),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    settings        JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE entities (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID NOT NULL REFERENCES companies(id) ON DELETE RESTRICT,
    name            VARCHAR(200) NOT NULL,
    code            VARCHAR(20) NOT NULL,
    country         VARCHAR(100) NOT NULL,
    currency        CHAR(3) NOT NULL,
    tax_id          VARCHAR(100),
    entity_type     VARCHAR(50) NOT NULL DEFAULT 'SUBSIDIARY', -- HQ, SUBSIDIARY, BRANCH, JV
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    parent_entity_id UUID REFERENCES entities(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, code)
);

-- ============================================================
-- USERS, ROLES, PERMISSIONS (RBAC)
-- ============================================================

CREATE TABLE roles (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id  UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name        VARCHAR(100) NOT NULL,
    code        VARCHAR(50) NOT NULL,
    description TEXT,
    is_system   BOOLEAN NOT NULL DEFAULT FALSE, -- System roles cannot be deleted
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, code)
);

CREATE TABLE permissions (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    resource    VARCHAR(100) NOT NULL, -- e.g., invoices, payments, fx_rates
    action      VARCHAR(50) NOT NULL,  -- create, read, update, delete, approve, export
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(resource, action)
);

CREATE TABLE role_permissions (
    role_id       UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY(role_id, permission_id)
);

CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID NOT NULL REFERENCES companies(id) ON DELETE RESTRICT,
    entity_id       UUID REFERENCES entities(id) ON DELETE SET NULL,
    email           VARCHAR(254) NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    first_name      VARCHAR(100) NOT NULL,
    last_name       VARCHAR(100) NOT NULL,
    phone           VARCHAR(30),
    avatar_url      VARCHAR(500),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    is_verified     BOOLEAN NOT NULL DEFAULT FALSE,
    last_login_at   TIMESTAMPTZ,
    failed_attempts INT NOT NULL DEFAULT 0,
    locked_until    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, email)
);

CREATE TABLE user_roles (
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id    UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    granted_by UUID REFERENCES users(id) ON DELETE SET NULL,
    granted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(user_id, role_id)
);

CREATE TABLE refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  VARCHAR(255) NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at  TIMESTAMPTZ,
    ip_address  INET,
    user_agent  TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- CURRENCIES & FX RATES
-- ============================================================

CREATE TABLE currencies (
    code            CHAR(3) PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    symbol          VARCHAR(10) NOT NULL,
    decimal_places  SMALLINT NOT NULL DEFAULT 2,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE fx_rates (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    from_currency   CHAR(3) NOT NULL REFERENCES currencies(code),
    to_currency     CHAR(3) NOT NULL REFERENCES currencies(code),
    rate            DECIMAL(20, 8) NOT NULL CHECK (rate > 0),
    bid_rate        DECIMAL(20, 8),
    ask_rate        DECIMAL(20, 8),
    source          VARCHAR(100) NOT NULL DEFAULT 'PROVIDER', -- PROVIDER, MANUAL, CACHED
    provider        VARCHAR(100),
    rate_date       DATE NOT NULL,
    is_latest       BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fx_rates_pair_date ON fx_rates(from_currency, to_currency, rate_date DESC);
CREATE INDEX idx_fx_rates_latest ON fx_rates(from_currency, to_currency) WHERE is_latest = TRUE;

-- ============================================================
-- BANK ACCOUNTS
-- ============================================================

CREATE TABLE bank_accounts (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id          UUID NOT NULL REFERENCES companies(id) ON DELETE RESTRICT,
    entity_id           UUID NOT NULL REFERENCES entities(id) ON DELETE RESTRICT,
    bank_name           VARCHAR(200) NOT NULL,
    account_number_masked VARCHAR(30) NOT NULL, -- XXXX1234
    account_number_hash VARCHAR(255) NOT NULL, -- SHA256 for lookup
    account_type        VARCHAR(50) NOT NULL, -- CURRENT, SAVINGS, OVERDRAFT, DEPOSIT, NOSTRO
    currency            CHAR(3) NOT NULL REFERENCES currencies(code),
    country             VARCHAR(100) NOT NULL,
    bank_code           VARCHAR(20), -- SWIFT / BIC / IFSC
    branch              VARCHAR(200),
    account_holder      VARCHAR(200) NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, INACTIVE, FROZEN, CLOSED
    available_balance   DECIMAL(20, 4) NOT NULL DEFAULT 0,
    ledger_balance      DECIMAL(20, 4) NOT NULL DEFAULT 0,
    restricted_amount   DECIMAL(20, 4) NOT NULL DEFAULT 0,
    credit_limit        DECIMAL(20, 4),
    last_sync_at        TIMESTAMPTZ,
    is_primary          BOOLEAN NOT NULL DEFAULT FALSE,
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bank_accounts_company ON bank_accounts(company_id);
CREATE INDEX idx_bank_accounts_entity ON bank_accounts(entity_id);

CREATE TABLE bank_transactions (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id          UUID NOT NULL REFERENCES companies(id),
    bank_account_id     UUID NOT NULL REFERENCES bank_accounts(id),
    transaction_date    DATE NOT NULL,
    value_date          DATE,
    description         TEXT NOT NULL,
    reference           VARCHAR(200),
    amount              DECIMAL(20, 4) NOT NULL,
    currency            CHAR(3) NOT NULL REFERENCES currencies(code),
    transaction_type    VARCHAR(50) NOT NULL, -- CREDIT, DEBIT
    balance_after       DECIMAL(20, 4),
    counterparty        VARCHAR(200),
    bank_reference      VARCHAR(200),
    status              VARCHAR(30) NOT NULL DEFAULT 'UNRECONCILED', -- UNRECONCILED, RECONCILED, EXCLUDED
    source              VARCHAR(50) NOT NULL DEFAULT 'IMPORT', -- IMPORT, MANUAL, API
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bank_txn_account_date ON bank_transactions(bank_account_id, transaction_date DESC);
CREATE INDEX idx_bank_txn_status ON bank_transactions(status);
CREATE INDEX idx_bank_txn_reference ON bank_transactions(reference);

-- ============================================================
-- CUSTOMERS
-- ============================================================

CREATE TABLE customers (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    customer_code   VARCHAR(50) NOT NULL,
    name            VARCHAR(200) NOT NULL,
    country         VARCHAR(100) NOT NULL,
    currency        CHAR(3) NOT NULL REFERENCES currencies(code),
    email           VARCHAR(254),
    phone           VARCHAR(30),
    address         TEXT,
    payment_terms   INT NOT NULL DEFAULT 30, -- Days
    credit_limit    DECIMAL(20, 4) NOT NULL DEFAULT 0,
    risk_level      VARCHAR(20) NOT NULL DEFAULT 'MEDIUM', -- LOW, MEDIUM, HIGH, CRITICAL
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, customer_code)
);

CREATE INDEX idx_customers_company ON customers(company_id);
CREATE INDEX idx_customers_name ON customers USING gin(name gin_trgm_ops);

-- ============================================================
-- SUPPLIERS
-- ============================================================

CREATE TABLE suppliers (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    supplier_code   VARCHAR(50) NOT NULL,
    name            VARCHAR(200) NOT NULL,
    country         VARCHAR(100) NOT NULL,
    currency        CHAR(3) NOT NULL REFERENCES currencies(code),
    email           VARCHAR(254),
    phone           VARCHAR(30),
    address         TEXT,
    payment_terms   INT NOT NULL DEFAULT 30,
    risk_level      VARCHAR(20) NOT NULL DEFAULT 'MEDIUM',
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, supplier_code)
);

CREATE INDEX idx_suppliers_company ON suppliers(company_id);
CREATE INDEX idx_suppliers_name ON suppliers USING gin(name gin_trgm_ops);

-- ============================================================
-- INVOICES (AR)
-- ============================================================

CREATE TABLE invoices (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id          UUID NOT NULL REFERENCES companies(id),
    entity_id           UUID NOT NULL REFERENCES entities(id),
    customer_id         UUID NOT NULL REFERENCES customers(id),
    invoice_number      VARCHAR(100) NOT NULL,
    invoice_date        DATE NOT NULL,
    due_date            DATE NOT NULL,
    currency            CHAR(3) NOT NULL REFERENCES currencies(code),
    subtotal            DECIMAL(20, 4) NOT NULL DEFAULT 0,
    tax_amount          DECIMAL(20, 4) NOT NULL DEFAULT 0,
    total_amount        DECIMAL(20, 4) NOT NULL DEFAULT 0,
    paid_amount         DECIMAL(20, 4) NOT NULL DEFAULT 0,
    outstanding_amount  DECIMAL(20, 4) GENERATED ALWAYS AS (total_amount - paid_amount) STORED,
    base_currency       CHAR(3) NOT NULL,
    base_amount         DECIMAL(20, 4),
    exchange_rate       DECIMAL(20, 8),
    rate_date           DATE,
    status              VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    -- DRAFT, ISSUED, PARTIALLY_PAID, PAID, OVERDUE, DISPUTED, CANCELLED, WRITTEN_OFF
    reference           VARCHAR(200),
    purchase_order      VARCHAR(100),
    notes               TEXT,
    created_by          UUID REFERENCES users(id),
    approved_by         UUID REFERENCES users(id),
    approved_at         TIMESTAMPTZ,
    idempotency_key     VARCHAR(100) UNIQUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, invoice_number)
);

CREATE INDEX idx_invoices_company ON invoices(company_id);
CREATE INDEX idx_invoices_customer ON invoices(customer_id);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_invoices_due_date ON invoices(due_date);
CREATE INDEX idx_invoices_number ON invoices USING gin(invoice_number gin_trgm_ops);

CREATE TABLE invoice_line_items (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    invoice_id      UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    description     TEXT NOT NULL,
    quantity        DECIMAL(15, 4) NOT NULL DEFAULT 1,
    unit_price      DECIMAL(20, 4) NOT NULL,
    tax_rate        DECIMAL(5, 4) NOT NULL DEFAULT 0,
    tax_amount      DECIMAL(20, 4) NOT NULL DEFAULT 0,
    line_total      DECIMAL(20, 4) NOT NULL,
    sort_order      INT NOT NULL DEFAULT 0
);

CREATE TABLE invoice_payments (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    invoice_id      UUID NOT NULL REFERENCES invoices(id),
    payment_date    DATE NOT NULL,
    amount          DECIMAL(20, 4) NOT NULL CHECK (amount > 0),
    currency        CHAR(3) NOT NULL REFERENCES currencies(code),
    base_amount     DECIMAL(20, 4),
    exchange_rate   DECIMAL(20, 8),
    payment_method  VARCHAR(50) NOT NULL DEFAULT 'BANK_TRANSFER',
    reference       VARCHAR(200),
    notes           TEXT,
    bank_account_id UUID REFERENCES bank_accounts(id),
    created_by      UUID REFERENCES users(id),
    idempotency_key VARCHAR(100) UNIQUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_invoice_payments_invoice ON invoice_payments(invoice_id);

-- ============================================================
-- BILLS (AP)
-- ============================================================

CREATE TABLE bills (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id          UUID NOT NULL REFERENCES companies(id),
    entity_id           UUID NOT NULL REFERENCES entities(id),
    supplier_id         UUID NOT NULL REFERENCES suppliers(id),
    bill_number         VARCHAR(100) NOT NULL,
    supplier_invoice_ref VARCHAR(100),
    invoice_date        DATE NOT NULL,
    due_date            DATE NOT NULL,
    currency            CHAR(3) NOT NULL REFERENCES currencies(code),
    subtotal            DECIMAL(20, 4) NOT NULL DEFAULT 0,
    tax_amount          DECIMAL(20, 4) NOT NULL DEFAULT 0,
    total_amount        DECIMAL(20, 4) NOT NULL DEFAULT 0,
    paid_amount         DECIMAL(20, 4) NOT NULL DEFAULT 0,
    outstanding_amount  DECIMAL(20, 4) GENERATED ALWAYS AS (total_amount - paid_amount) STORED,
    base_currency       CHAR(3) NOT NULL,
    base_amount         DECIMAL(20, 4),
    exchange_rate       DECIMAL(20, 8),
    rate_date           DATE,
    status              VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    -- DRAFT, APPROVED, SCHEDULED, PARTIALLY_PAID, PAID, OVERDUE, DISPUTED, CANCELLED
    reference           VARCHAR(200),
    notes               TEXT,
    created_by          UUID REFERENCES users(id),
    approved_by         UUID REFERENCES users(id),
    approved_at         TIMESTAMPTZ,
    idempotency_key     VARCHAR(100) UNIQUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, bill_number)
);

CREATE INDEX idx_bills_company ON bills(company_id);
CREATE INDEX idx_bills_supplier ON bills(supplier_id);
CREATE INDEX idx_bills_status ON bills(status);
CREATE INDEX idx_bills_due_date ON bills(due_date);

CREATE TABLE bill_payments (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    bill_id         UUID NOT NULL REFERENCES bills(id),
    payment_date    DATE NOT NULL,
    amount          DECIMAL(20, 4) NOT NULL CHECK (amount > 0),
    currency        CHAR(3) NOT NULL REFERENCES currencies(code),
    base_amount     DECIMAL(20, 4),
    exchange_rate   DECIMAL(20, 8),
    payment_method  VARCHAR(50) NOT NULL DEFAULT 'BANK_TRANSFER',
    reference       VARCHAR(200),
    notes           TEXT,
    bank_account_id UUID REFERENCES bank_accounts(id),
    created_by      UUID REFERENCES users(id),
    idempotency_key VARCHAR(100) UNIQUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bill_payments_bill ON bill_payments(bill_id);

-- ============================================================
-- PAYMENTS (General — AR + AP)
-- ============================================================

CREATE TABLE payments (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id          UUID NOT NULL REFERENCES companies(id),
    entity_id           UUID NOT NULL REFERENCES entities(id),
    payment_type        VARCHAR(20) NOT NULL, -- OUTGOING, INCOMING
    category            VARCHAR(50) NOT NULL, -- SUPPLIER, CUSTOMER, PAYROLL, TAX, LOAN, OTHER
    reference_type      VARCHAR(50), -- BILL, INVOICE, MANUAL
    reference_id        UUID,
    counterparty_name   VARCHAR(200) NOT NULL,
    amount              DECIMAL(20, 4) NOT NULL CHECK (amount > 0),
    currency            CHAR(3) NOT NULL REFERENCES currencies(code),
    base_amount         DECIMAL(20, 4),
    exchange_rate       DECIMAL(20, 8),
    bank_account_id     UUID REFERENCES bank_accounts(id),
    payment_date        DATE,
    value_date          DATE,
    status              VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    -- DRAFT, PENDING_APPROVAL, APPROVED, SCHEDULED, PROCESSING, COMPLETED, FAILED, CANCELLED
    batch_id            UUID,
    reference           VARCHAR(200),
    notes               TEXT,
    created_by          UUID NOT NULL REFERENCES users(id),
    idempotency_key     VARCHAR(100) UNIQUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_company ON payments(company_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_created_by ON payments(created_by);

CREATE TABLE payment_approvals (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    payment_id  UUID NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    approver_id UUID NOT NULL REFERENCES users(id),
    level       SMALLINT NOT NULL DEFAULT 1, -- Approval level
    decision    VARCHAR(20) NOT NULL, -- APPROVED, REJECTED
    reason      TEXT,
    decided_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- CASH MANAGEMENT
-- ============================================================

CREATE TABLE cash_accounts (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    entity_id       UUID NOT NULL REFERENCES entities(id),
    bank_account_id UUID NOT NULL REFERENCES bank_accounts(id),
    name            VARCHAR(200) NOT NULL,
    currency        CHAR(3) NOT NULL REFERENCES currencies(code),
    opening_balance DECIMAL(20, 4) NOT NULL DEFAULT 0,
    current_balance DECIMAL(20, 4) NOT NULL DEFAULT 0,
    balance_date    DATE NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE cash_movements (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    cash_account_id UUID NOT NULL REFERENCES cash_accounts(id),
    movement_date   DATE NOT NULL,
    movement_type   VARCHAR(30) NOT NULL, -- INFLOW, OUTFLOW, TRANSFER_IN, TRANSFER_OUT
    category        VARCHAR(100) NOT NULL, -- COLLECTION, PAYMENT, LOAN, INVESTMENT, PAYROLL, etc.
    amount          DECIMAL(20, 4) NOT NULL CHECK (amount > 0),
    currency        CHAR(3) NOT NULL REFERENCES currencies(code),
    base_amount     DECIMAL(20, 4),
    exchange_rate   DECIMAL(20, 8),
    reference_type  VARCHAR(50),
    reference_id    UUID,
    description     TEXT NOT NULL,
    is_forecast     BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cash_movements_account_date ON cash_movements(cash_account_id, movement_date DESC);
CREATE INDEX idx_cash_movements_date ON cash_movements(movement_date);

-- ============================================================
-- RECONCILIATION
-- ============================================================

CREATE TABLE reconciliations (
    id                      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id              UUID NOT NULL REFERENCES companies(id),
    bank_account_id         UUID NOT NULL REFERENCES bank_accounts(id),
    bank_transaction_id     UUID REFERENCES bank_transactions(id),
    ledger_reference_type   VARCHAR(50), -- INVOICE, BILL, PAYMENT, CASH_MOVEMENT
    ledger_reference_id     UUID,
    match_type              VARCHAR(30) NOT NULL DEFAULT 'AUTO', -- AUTO, MANUAL, PARTIAL
    status                  VARCHAR(30) NOT NULL DEFAULT 'MATCHED', -- MATCHED, PARTIAL, UNMATCHED, EXCEPTION
    match_amount            DECIMAL(20, 4),
    variance_amount         DECIMAL(20, 4) NOT NULL DEFAULT 0,
    notes                   TEXT,
    reconciled_by           UUID REFERENCES users(id),
    reconciled_at           TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_recon_bank_account ON reconciliations(bank_account_id);
CREATE INDEX idx_recon_status ON reconciliations(status);

-- ============================================================
-- FX EXPOSURE & DEALS
-- ============================================================

CREATE TABLE fx_exposures (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id          UUID NOT NULL REFERENCES companies(id),
    entity_id           UUID NOT NULL REFERENCES entities(id),
    exposure_date       DATE NOT NULL,
    currency            CHAR(3) NOT NULL REFERENCES currencies(code),
    exposure_type       VARCHAR(30) NOT NULL, -- RECEIVABLE, PAYABLE, CASH, INVESTMENT, LOAN
    gross_exposure      DECIMAL(20, 4) NOT NULL DEFAULT 0,
    hedged_amount       DECIMAL(20, 4) NOT NULL DEFAULT 0,
    net_exposure        DECIMAL(20, 4) GENERATED ALWAYS AS (gross_exposure - hedged_amount) STORED,
    base_currency       CHAR(3) NOT NULL,
    base_value          DECIMAL(20, 4),
    exchange_rate       DECIMAL(20, 8),
    risk_level          VARCHAR(20) NOT NULL DEFAULT 'MEDIUM',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fx_exposure_company_date ON fx_exposures(company_id, exposure_date DESC);
CREATE INDEX idx_fx_exposure_currency ON fx_exposures(currency);

CREATE TABLE fx_deals (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    entity_id       UUID NOT NULL REFERENCES entities(id),
    deal_number     VARCHAR(100) NOT NULL,
    deal_type       VARCHAR(30) NOT NULL, -- SPOT, FORWARD, SWAP, OPTION
    direction       VARCHAR(10) NOT NULL, -- BUY, SELL
    buy_currency    CHAR(3) NOT NULL REFERENCES currencies(code),
    sell_currency   CHAR(3) NOT NULL REFERENCES currencies(code),
    notional        DECIMAL(20, 4) NOT NULL CHECK (notional > 0),
    rate            DECIMAL(20, 8) NOT NULL,
    trade_date      DATE NOT NULL,
    settlement_date DATE,
    maturity_date   DATE,
    counterparty    VARCHAR(200),
    status          VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    -- DRAFT, PENDING_APPROVAL, APPROVED, CONFIRMED, SETTLED, CANCELLED
    notes           TEXT,
    created_by      UUID REFERENCES users(id),
    approved_by     UUID REFERENCES users(id),
    approved_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, deal_number)
);

CREATE INDEX idx_fx_deals_company ON fx_deals(company_id);
CREATE INDEX idx_fx_deals_status ON fx_deals(status);

CREATE TABLE hedges (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id          UUID NOT NULL REFERENCES companies(id),
    fx_deal_id          UUID REFERENCES fx_deals(id),
    fx_exposure_id      UUID REFERENCES fx_exposures(id),
    instrument_type     VARCHAR(30) NOT NULL, -- FORWARD, OPTION, SWAP
    currency            CHAR(3) NOT NULL REFERENCES currencies(code),
    notional            DECIMAL(20, 4) NOT NULL,
    hedge_rate          DECIMAL(20, 8) NOT NULL,
    maturity_date       DATE NOT NULL,
    direction           VARCHAR(10) NOT NULL, -- BUY, SELL
    counterparty        VARCHAR(200),
    status              VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, EXPIRED, CANCELLED
    coverage_ratio      DECIMAL(5, 4),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- RISK MANAGEMENT
-- ============================================================

CREATE TABLE risk_limits (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    entity_id       UUID REFERENCES entities(id),
    risk_type       VARCHAR(50) NOT NULL, -- FX_EXPOSURE, COUNTERPARTY, CONCENTRATION, LIQUIDITY
    currency        CHAR(3) REFERENCES currencies(code),
    limit_amount    DECIMAL(20, 4) NOT NULL,
    warning_amount  DECIMAL(20, 4),
    limit_currency  CHAR(3) NOT NULL REFERENCES currencies(code),
    description     TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_by      UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE risk_events (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    risk_limit_id   UUID REFERENCES risk_limits(id),
    event_type      VARCHAR(50) NOT NULL,
    severity        VARCHAR(20) NOT NULL, -- LOW, MEDIUM, HIGH, CRITICAL
    risk_score      DECIMAL(5, 2) NOT NULL DEFAULT 0,
    title           TEXT NOT NULL,
    description     TEXT NOT NULL,
    affected_entity VARCHAR(200),
    affected_amount DECIMAL(20, 4),
    currency        CHAR(3) REFERENCES currencies(code),
    status          VARCHAR(30) NOT NULL DEFAULT 'OPEN',
    -- OPEN, ACKNOWLEDGED, INVESTIGATING, RESOLVED, DISMISSED
    owner_id        UUID REFERENCES users(id),
    resolved_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_risk_events_company ON risk_events(company_id);
CREATE INDEX idx_risk_events_status ON risk_events(status);
CREATE INDEX idx_risk_events_severity ON risk_events(severity);

-- ============================================================
-- CASH FLOW FORECASTING
-- ============================================================

CREATE TABLE cash_forecasts (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    entity_id       UUID REFERENCES entities(id),
    forecast_date   DATE NOT NULL,
    horizon_days    INT NOT NULL, -- 7, 30, 60, 90, 180
    currency        CHAR(3) NOT NULL REFERENCES currencies(code),
    model_type      VARCHAR(50) NOT NULL DEFAULT 'WEIGHTED_AVERAGE',
    status          VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, ARCHIVED
    generated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    mae             DECIMAL(10, 4), -- Mean Absolute Error
    rmse            DECIMAL(10, 4), -- Root Mean Square Error
    mape            DECIMAL(10, 4), -- Mean Absolute Percentage Error
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE forecast_items (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    forecast_id     UUID NOT NULL REFERENCES cash_forecasts(id) ON DELETE CASCADE,
    item_date       DATE NOT NULL,
    category        VARCHAR(100) NOT NULL, -- COLLECTIONS, SUPPLIER_PAYMENTS, PAYROLL, etc.
    projected_inflow  DECIMAL(20, 4) NOT NULL DEFAULT 0,
    projected_outflow DECIMAL(20, 4) NOT NULL DEFAULT 0,
    actual_inflow     DECIMAL(20, 4),
    actual_outflow    DECIMAL(20, 4),
    projected_balance DECIMAL(20, 4) NOT NULL DEFAULT 0,
    confidence      DECIMAL(5, 4), -- 0.0 to 1.0
    notes           TEXT
);

CREATE INDEX idx_forecast_items_forecast ON forecast_items(forecast_id, item_date);

-- ============================================================
-- SCENARIOS
-- ============================================================

CREATE TABLE scenarios (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    name            VARCHAR(200) NOT NULL,
    description     TEXT,
    scenario_type   VARCHAR(50) NOT NULL DEFAULT 'CUSTOM',
    -- NORMAL, DELAYED_COLLECTIONS, FX_APPRECIATION, FX_DEPRECIATION, STRESS, CUSTOM
    parameters      JSONB NOT NULL DEFAULT '{}',
    is_predefined   BOOLEAN NOT NULL DEFAULT FALSE,
    created_by      UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE scenario_results (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    scenario_id         UUID NOT NULL REFERENCES scenarios(id) ON DELETE CASCADE,
    result_date         DATE NOT NULL,
    cash_balance        DECIMAL(20, 4),
    working_capital     DECIMAL(20, 4),
    fx_exposure         DECIMAL(20, 4),
    net_cash_flow       DECIMAL(20, 4),
    projected_inflows   DECIMAL(20, 4),
    projected_outflows  DECIMAL(20, 4),
    ccc                 DECIMAL(10, 4),
    dso                 DECIMAL(10, 4),
    dpo                 DECIMAL(10, 4),
    liquidity_ratio     DECIMAL(10, 4),
    calculated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- ALERTS & NOTIFICATIONS
-- ============================================================

CREATE TABLE alerts (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    alert_type      VARCHAR(50) NOT NULL,
    -- CASH_BELOW_THRESHOLD, INVOICE_OVERDUE, FX_LIMIT, FORECAST_SHORTAGE, etc.
    severity        VARCHAR(20) NOT NULL, -- LOW, MEDIUM, HIGH, CRITICAL
    title           TEXT NOT NULL,
    description     TEXT NOT NULL,
    source          VARCHAR(100) NOT NULL,
    entity_id       UUID REFERENCES entities(id),
    affected_amount DECIMAL(20, 4),
    currency        CHAR(3) REFERENCES currencies(code),
    reference_type  VARCHAR(50),
    reference_id    UUID,
    recommended_action TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, ACKNOWLEDGED, RESOLVED, DISMISSED
    acknowledged_by UUID REFERENCES users(id),
    acknowledged_at TIMESTAMPTZ,
    resolved_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_alerts_company ON alerts(company_id);
CREATE INDEX idx_alerts_status ON alerts(status);
CREATE INDEX idx_alerts_severity ON alerts(severity);

CREATE TABLE notifications (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    user_id         UUID NOT NULL REFERENCES users(id),
    alert_id        UUID REFERENCES alerts(id),
    category        VARCHAR(50) NOT NULL, -- CRITICAL, RISK, PAYMENT, COLLECTION, FORECAST, SYSTEM
    title           TEXT NOT NULL,
    message         TEXT NOT NULL,
    priority        SMALLINT NOT NULL DEFAULT 5, -- 1-10
    link_type       VARCHAR(50),
    link_id         UUID,
    link_url        VARCHAR(500),
    is_read         BOOLEAN NOT NULL DEFAULT FALSE,
    read_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_user ON notifications(user_id, is_read, created_at DESC);

-- ============================================================
-- AUDIT LOGGING (Immutable)
-- ============================================================

CREATE TABLE audit_logs (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    user_id         UUID REFERENCES users(id),
    user_email      VARCHAR(254), -- Denormalized for permanent record
    action          VARCHAR(100) NOT NULL,
    resource        VARCHAR(100) NOT NULL,
    resource_id     VARCHAR(100),
    old_value       JSONB,
    new_value       JSONB,
    ip_address      INET,
    user_agent      TEXT,
    request_id      UUID,
    severity        VARCHAR(20) NOT NULL DEFAULT 'INFO', -- INFO, WARNING, CRITICAL
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
    -- NOTE: No updated_at — audit logs are immutable
);

CREATE INDEX idx_audit_logs_company_created ON audit_logs(company_id, created_at DESC);
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id, created_at DESC);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource, resource_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);

-- Prevent deletion and updates of audit logs (trigger)
CREATE OR REPLACE FUNCTION prevent_audit_modification()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'Audit logs are immutable and cannot be modified or deleted';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER audit_logs_immutable_delete
    BEFORE DELETE ON audit_logs
    FOR EACH ROW EXECUTE FUNCTION prevent_audit_modification();

CREATE TRIGGER audit_logs_immutable_update
    BEFORE UPDATE ON audit_logs
    FOR EACH ROW EXECUTE FUNCTION prevent_audit_modification();

-- ============================================================
-- SYSTEM EVENTS & BACKGROUND JOBS
-- ============================================================

CREATE TABLE system_events (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_type      VARCHAR(100) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'PENDING', -- PENDING, RUNNING, COMPLETED, FAILED
    payload         JSONB,
    result          JSONB,
    error_message   TEXT,
    retry_count     INT NOT NULL DEFAULT 0,
    max_retries     INT NOT NULL DEFAULT 3,
    next_run_at     TIMESTAMPTZ,
    started_at      TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ,
    execution_ms    INT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_system_events_status ON system_events(status, next_run_at);
CREATE INDEX idx_system_events_type ON system_events(event_type);

-- ============================================================
-- APPROVAL WORKFLOWS
-- ============================================================

CREATE TABLE approval_workflows (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    name            VARCHAR(200) NOT NULL,
    resource_type   VARCHAR(50) NOT NULL, -- PAYMENT, FX_DEAL, RISK_LIMIT, WRITE_OFF
    min_amount      DECIMAL(20, 4),
    max_amount      DECIMAL(20, 4),
    currency        CHAR(3) REFERENCES currencies(code),
    required_levels SMALLINT NOT NULL DEFAULT 1,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- REPORT GENERATION
-- ============================================================

CREATE TABLE reports (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID NOT NULL REFERENCES companies(id),
    report_type     VARCHAR(100) NOT NULL,
    name            VARCHAR(200) NOT NULL,
    parameters      JSONB NOT NULL DEFAULT '{}',
    status          VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    file_url        VARCHAR(500),
    file_size       BIGINT,
    error_message   TEXT,
    requested_by    UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at    TIMESTAMPTZ
);

-- ============================================================
-- CONFIGURATION
-- ============================================================

CREATE TABLE company_settings (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id      UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    setting_key     VARCHAR(100) NOT NULL,
    setting_value   TEXT NOT NULL,
    setting_type    VARCHAR(20) NOT NULL DEFAULT 'STRING', -- STRING, NUMBER, BOOLEAN, JSON
    description     TEXT,
    updated_by      UUID REFERENCES users(id),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, setting_key)
);

-- ============================================================
-- TIMESTAMPS TRIGGER
-- ============================================================

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply to all tables with updated_at
DO $$
DECLARE
    t TEXT;
BEGIN
    FOR t IN
        SELECT table_name FROM information_schema.columns
        WHERE column_name = 'updated_at'
          AND table_schema = 'public'
    LOOP
        EXECUTE format('DROP TRIGGER IF EXISTS set_updated_at ON %I', t);
        EXECUTE format('
            CREATE TRIGGER set_updated_at
            BEFORE UPDATE ON %I
            FOR EACH ROW EXECUTE FUNCTION update_updated_at()', t);
    END LOOP;
END;
$$;
