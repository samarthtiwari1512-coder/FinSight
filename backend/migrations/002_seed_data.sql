-- FinSight Demo Seed Data
-- Company: Acme Global Manufacturing Ltd
-- Realistic enterprise data for demo/development

-- ============================================================
-- CURRENCIES
-- ============================================================

INSERT INTO currencies (code, name, symbol, decimal_places) VALUES
('INR', 'Indian Rupee',       '₹',  2),
('USD', 'US Dollar',          '$',  2),
('EUR', 'Euro',               '€',  2),
('GBP', 'British Pound',      '£',  2),
('JPY', 'Japanese Yen',       '¥',  0),
('CHF', 'Swiss Franc',        'Fr', 2),
('AUD', 'Australian Dollar',  'A$', 2),
('CAD', 'Canadian Dollar',    'C$', 2),
('SGD', 'Singapore Dollar',   'S$', 2)
ON CONFLICT DO NOTHING;

-- ============================================================
-- COMPANY & ENTITIES
-- ============================================================

INSERT INTO companies (id, name, code, base_currency, country, fiscal_year_start, timezone, address)
VALUES (
    'a0000000-0000-0000-0000-000000000001',
    'Acme Global Manufacturing Ltd',
    'ACME',
    'INR',
    'India',
    4,
    'Asia/Kolkata',
    '14th Floor, One BKC, Bandra Kurla Complex, Mumbai 400051, India'
);

INSERT INTO entities (id, company_id, name, code, country, currency, entity_type, tax_id) VALUES
('b0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001', 'Acme Manufacturing India Pvt Ltd', 'ACME-IN', 'India',     'INR', 'HQ',         'AABCA1234C'),
('b0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001', 'Acme Manufacturing USA Inc',       'ACME-US', 'USA',       'USD', 'SUBSIDIARY', 'EIN-12-3456789'),
('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', 'Acme Manufacturing GmbH',          'ACME-DE', 'Germany',   'EUR', 'SUBSIDIARY', 'DE123456789'),
('b0000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000001', 'Acme Manufacturing Pte Ltd',       'ACME-SG', 'Singapore', 'SGD', 'SUBSIDIARY', 'SG-123456789A');

-- ============================================================
-- ROLES & PERMISSIONS
-- ============================================================

INSERT INTO roles (id, company_id, name, code, description, is_system) VALUES
('c0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001', 'CFO',                    'CFO',               'Chief Financial Officer - full read access + approvals', TRUE),
('c0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001', 'Treasury Manager',       'TREASURY_MANAGER',  'Manage cash, FX, bank accounts, and forecasting', TRUE),
('c0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', 'Finance Manager',        'FINANCE_MANAGER',   'AR, AP, invoices, reconciliation, payments', TRUE),
('c0000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000001', 'AR User',                'AR_USER',           'Accounts receivable and collections', TRUE),
('c0000000-0000-0000-0000-000000000005', 'a0000000-0000-0000-0000-000000000001', 'AP User',                'AP_USER',           'Accounts payable and supplier payments', TRUE),
('c0000000-0000-0000-0000-000000000006', 'a0000000-0000-0000-0000-000000000001', 'Risk Analyst',           'RISK_ANALYST',      'FX exposure, scenarios, risk management', TRUE),
('c0000000-0000-0000-0000-000000000007', 'a0000000-0000-0000-0000-000000000001', 'Auditor',                'AUDITOR',           'Read-only access including audit logs', TRUE),
('c0000000-0000-0000-0000-000000000008', 'a0000000-0000-0000-0000-000000000001', 'Administrator',          'ADMIN',             'Full system administration', TRUE);

INSERT INTO permissions (resource, action, description) VALUES
('dashboard',    'read',    'View executive dashboard'),
('cash',         'read',    'View cash positions'),
('cash',         'export',  'Export cash reports'),
('bank_accounts','read',    'View bank accounts'),
('bank_accounts','create',  'Create bank accounts'),
('bank_accounts','update',  'Update bank accounts'),
('customers',    'read',    'View customers'),
('customers',    'create',  'Create customers'),
('customers',    'update',  'Update customers'),
('suppliers',    'read',    'View suppliers'),
('suppliers',    'create',  'Create suppliers'),
('suppliers',    'update',  'Update suppliers'),
('invoices',     'read',    'View invoices'),
('invoices',     'create',  'Create invoices'),
('invoices',     'update',  'Update invoices'),
('invoices',     'approve', 'Approve invoices'),
('invoices',     'export',  'Export invoice data'),
('bills',        'read',    'View bills'),
('bills',        'create',  'Create bills'),
('bills',        'update',  'Update bills'),
('bills',        'approve', 'Approve bills'),
('payments',     'read',    'View payments'),
('payments',     'create',  'Create payments'),
('payments',     'approve', 'Approve payments'),
('payments',     'export',  'Export payment data'),
('fx',           'read',    'View FX rates and exposure'),
('fx',           'create',  'Create FX deals'),
('fx',           'approve', 'Approve FX deals'),
('forecast',     'read',    'View cash forecasts'),
('forecast',     'create',  'Generate forecasts'),
('risk',         'read',    'View risk events'),
('risk',         'manage',  'Manage risk limits'),
('scenarios',    'read',    'View scenarios'),
('scenarios',    'create',  'Create scenarios'),
('reconciliation','read',   'View reconciliations'),
('reconciliation','create', 'Perform reconciliations'),
('alerts',       'read',    'View alerts'),
('alerts',       'manage',  'Manage alert rules'),
('reports',      'read',    'View reports'),
('reports',      'create',  'Generate reports'),
('audit_logs',   'read',    'View audit logs'),
('users',        'read',    'View users'),
('users',        'create',  'Create users'),
('users',        'update',  'Update users'),
('settings',     'read',    'View system settings'),
('settings',     'update',  'Update system settings'),
('working_capital','read',  'View working capital KPIs');

-- ============================================================
-- DEMO USERS
-- All passwords are: Acme@12345
-- Hash: bcrypt of "Acme@12345"
-- ============================================================

INSERT INTO users (id, company_id, entity_id, email, password_hash, first_name, last_name, is_active, is_verified) VALUES
('d0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001', 'cfo@acmeglobal.com',       '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LeuvRbBJTRJK9bxnS', 'Priya',   'Sharma',    TRUE, TRUE),
('d0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001', 'treasury@acmeglobal.com',  '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LeuvRbBJTRJK9bxnS', 'Rahul',   'Mehta',     TRUE, TRUE),
('d0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001', 'finance@acmeglobal.com',   '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LeuvRbBJTRJK9bxnS', 'Anita',   'Desai',     TRUE, TRUE),
('d0000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001', 'ar@acmeglobal.com',        '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LeuvRbBJTRJK9bxnS', 'Suresh',  'Kumar',     TRUE, TRUE),
('d0000000-0000-0000-0000-000000000005', 'a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001', 'ap@acmeglobal.com',        '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LeuvRbBJTRJK9bxnS', 'Meena',   'Pillai',    TRUE, TRUE),
('d0000000-0000-0000-0000-000000000006', 'a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001', 'risk@acmeglobal.com',      '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LeuvRbBJTRJK9bxnS', 'Arun',    'Nair',      TRUE, TRUE),
('d0000000-0000-0000-0000-000000000007', 'a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001', 'auditor@acmeglobal.com',   '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LeuvRbBJTRJK9bxnS', 'Kavita',  'Joshi',     TRUE, TRUE),
('d0000000-0000-0000-0000-000000000008', 'a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001', 'admin@acmeglobal.com',     '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LeuvRbBJTRJK9bxnS', 'Dev',     'Admin',     TRUE, TRUE);

-- Role assignments
INSERT INTO user_roles (user_id, role_id) VALUES
('d0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000001'), -- CFO
('d0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000002'), -- Treasury Manager
('d0000000-0000-0000-0000-000000000003', 'c0000000-0000-0000-0000-000000000003'), -- Finance Manager
('d0000000-0000-0000-0000-000000000004', 'c0000000-0000-0000-0000-000000000004'), -- AR User
('d0000000-0000-0000-0000-000000000005', 'c0000000-0000-0000-0000-000000000005'), -- AP User
('d0000000-0000-0000-0000-000000000006', 'c0000000-0000-0000-0000-000000000006'), -- Risk Analyst
('d0000000-0000-0000-0000-000000000007', 'c0000000-0000-0000-0000-000000000007'), -- Auditor
('d0000000-0000-0000-0000-000000000008', 'c0000000-0000-0000-0000-000000000008'); -- Admin

-- ============================================================
-- FX RATES (Seed historical + current rates)
-- Base: INR
-- ============================================================

INSERT INTO fx_rates (from_currency, to_currency, rate, source, provider, rate_date, is_latest) VALUES
('USD', 'INR', 83.42, 'MANUAL', 'SEED', CURRENT_DATE, TRUE),
('EUR', 'INR', 90.15, 'MANUAL', 'SEED', CURRENT_DATE, TRUE),
('GBP', 'INR', 105.32, 'MANUAL', 'SEED', CURRENT_DATE, TRUE),
('JPY', 'INR', 0.5523, 'MANUAL', 'SEED', CURRENT_DATE, TRUE),
('CHF', 'INR', 93.80, 'MANUAL', 'SEED', CURRENT_DATE, TRUE),
('AUD', 'INR', 54.21, 'MANUAL', 'SEED', CURRENT_DATE, TRUE),
('SGD', 'INR', 61.85, 'MANUAL', 'SEED', CURRENT_DATE, TRUE),
('CAD', 'INR', 60.94, 'MANUAL', 'SEED', CURRENT_DATE, TRUE),
('INR', 'USD', 0.011987, 'MANUAL', 'SEED', CURRENT_DATE, TRUE),
('INR', 'EUR', 0.011093, 'MANUAL', 'SEED', CURRENT_DATE, TRUE),
('INR', 'GBP', 0.009495, 'MANUAL', 'SEED', CURRENT_DATE, TRUE);

-- Historical rates (last 30 days simulation)
INSERT INTO fx_rates (from_currency, to_currency, rate, source, provider, rate_date, is_latest)
SELECT
    'USD', 'INR',
    83.42 + (RANDOM() * 2 - 1),
    'MANUAL', 'SEED',
    CURRENT_DATE - (n || ' days')::INTERVAL,
    FALSE
FROM generate_series(1, 30) AS n;

INSERT INTO fx_rates (from_currency, to_currency, rate, source, provider, rate_date, is_latest)
SELECT
    'EUR', 'INR',
    90.15 + (RANDOM() * 2 - 1),
    'MANUAL', 'SEED',
    CURRENT_DATE - (n || ' days')::INTERVAL,
    FALSE
FROM generate_series(1, 30) AS n;

INSERT INTO fx_rates (from_currency, to_currency, rate, source, provider, rate_date, is_latest)
SELECT
    'GBP', 'INR',
    105.32 + (RANDOM() * 2 - 1),
    'MANUAL', 'SEED',
    CURRENT_DATE - (n || ' days')::INTERVAL,
    FALSE
FROM generate_series(1, 30) AS n;

-- ============================================================
-- BANK ACCOUNTS
-- ============================================================

INSERT INTO bank_accounts (id, company_id, entity_id, bank_name, account_number_masked, account_number_hash,
    account_type, currency, country, bank_code, account_holder, available_balance, ledger_balance, is_primary) VALUES
-- India accounts
('e0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001',
    'HDFC Bank', 'XXXX-XXXX-4521', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8',
    'CURRENT', 'INR', 'India', 'HDFCINBB', 'Acme Manufacturing India Pvt Ltd', 182500000.00, 182500000.00, TRUE),

('e0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001',
    'ICICI Bank', 'XXXX-XXXX-7832', '6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b',
    'CURRENT', 'INR', 'India', 'ICICINBB', 'Acme Manufacturing India Pvt Ltd', 95300000.00, 95300000.00, FALSE),

('e0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001',
    'State Bank of India', 'XXXX-XXXX-1109', 'd4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35',
    'OVERDRAFT', 'INR', 'India', 'SBININBB', 'Acme Manufacturing India Pvt Ltd', 50000000.00, 50000000.00, FALSE),

-- USD account (India entity)
('e0000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001',
    'HDFC Bank', 'XXXX-XXXX-5678', '4e07408562bedb8b60ce05c1decfe3ad16b1a0e6b12cd81b24a7f4a1d7d64c0',
    'NOSTRO', 'USD', 'India', 'HDFCINBB', 'Acme Manufacturing India Pvt Ltd', 3500000.00, 3500000.00, FALSE),

-- USA accounts
('e0000000-0000-0000-0000-000000000005', 'a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000002',
    'JPMorgan Chase', 'XXXX-XXXX-2345', 'ef2d127de37b942baad06145e54b0c619a1f22327b2ebbcfbec78f5564afe39d',
    'CURRENT', 'USD', 'USA', 'CHASUS33', 'Acme Manufacturing USA Inc', 4200000.00, 4200000.00, TRUE),

('e0000000-0000-0000-0000-000000000006', 'a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000002',
    'Bank of America', 'XXXX-XXXX-8891', 'e7f6c011776e8db7cd330b54174fd76f7d0216b612387a5ffcfb81e6f0919683',
    'SAVINGS', 'USD', 'USA', 'BOFAUS3N', 'Acme Manufacturing USA Inc', 1800000.00, 1800000.00, FALSE),

-- Germany accounts
('e0000000-0000-0000-0000-000000000007', 'a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000003',
    'Deutsche Bank', 'XXXX-XXXX-6634', '19581e27de7ced00ff1ce50b2047e7a567c76b1cbaebabe5ef03f7c3017bb5b7',
    'CURRENT', 'EUR', 'Germany', 'DEUTDEDB', 'Acme Manufacturing GmbH', 2100000.00, 2100000.00, TRUE),

-- Singapore accounts
('e0000000-0000-0000-0000-000000000008', 'a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000004',
    'DBS Bank', 'XXXX-XXXX-3310', '4a44dc15364204a80fe80e9039455cc1608281820fe2b24f1e5233ade6af1dd5',
    'CURRENT', 'SGD', 'Singapore', 'DBSSSGSG', 'Acme Manufacturing Pte Ltd', 1200000.00, 1200000.00, TRUE);

-- ============================================================
-- CUSTOMERS (60 customers)
-- ============================================================

INSERT INTO customers (company_id, customer_code, name, country, currency, payment_terms, credit_limit, risk_level) VALUES
('a0000000-0000-0000-0000-000000000001', 'CUST-001', 'Global Tech Solutions Inc', 'USA', 'USD', 45, 50000000.00, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'CUST-002', 'Precision Parts Europe GmbH', 'Germany', 'EUR', 30, 35000000.00, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'CUST-003', 'AutoMotive UK Ltd', 'UK', 'GBP', 60, 40000000.00, 'MEDIUM'),
('a0000000-0000-0000-0000-000000000001', 'CUST-004', 'Pacific Machinery Co', 'Japan', 'JPY', 30, 20000000.00, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'CUST-005', 'TechMart Asia Pte Ltd', 'Singapore', 'SGD', 30, 15000000.00, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'CUST-006', 'National Auto Components', 'India', 'INR', 45, 25000000.00, 'MEDIUM'),
('a0000000-0000-0000-0000-000000000001', 'CUST-007', 'Bharat Industries Ltd', 'India', 'INR', 30, 30000000.00, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'CUST-008', 'Sunrise Export Co', 'India', 'INR', 60, 10000000.00, 'HIGH'),
('a0000000-0000-0000-0000-000000000001', 'CUST-009', 'Steel Fabricators USA LLC', 'USA', 'USD', 30, 20000000.00, 'MEDIUM'),
('a0000000-0000-0000-0000-000000000001', 'CUST-010', 'European Industrial AG', 'Switzerland', 'CHF', 30, 25000000.00, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'CUST-011', 'AustroParts GmbH', 'Germany', 'EUR', 45, 18000000.00, 'MEDIUM'),
('a0000000-0000-0000-0000-000000000001', 'CUST-012', 'Southeast Motors Pte', 'Singapore', 'SGD', 30, 12000000.00, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'CUST-013', 'North America Parts Corp', 'USA', 'USD', 45, 32000000.00, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'CUST-014', 'Kansei Electronics Japan', 'Japan', 'JPY', 60, 15000000.00, 'MEDIUM'),
('a0000000-0000-0000-0000-000000000001', 'CUST-015', 'Global Components France', 'France', 'EUR', 30, 20000000.00, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'CUST-016', 'Maharashtra Auto Parts', 'India', 'INR', 30, 8000000.00, 'MEDIUM'),
('a0000000-0000-0000-0000-000000000001', 'CUST-017', 'Gujarat Metal Works', 'India', 'INR', 45, 12000000.00, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'CUST-018', 'Delhi Electronics Pvt Ltd', 'India', 'INR', 30, 9000000.00, 'HIGH'),
('a0000000-0000-0000-0000-000000000001', 'CUST-019', 'Midlands Manufacturing UK', 'UK', 'GBP', 45, 22000000.00, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'CUST-020', 'Nordic Components AB', 'Sweden', 'EUR', 30, 16000000.00, 'LOW');

-- ============================================================
-- SUPPLIERS (30 suppliers)
-- ============================================================

INSERT INTO suppliers (company_id, supplier_code, name, country, currency, payment_terms, risk_level) VALUES
('a0000000-0000-0000-0000-000000000001', 'SUPP-001', 'RawMat Industries USA', 'USA', 'USD', 30, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'SUPP-002', 'European Steel AG', 'Germany', 'EUR', 45, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'SUPP-003', 'Nippon Components Co', 'Japan', 'JPY', 60, 'MEDIUM'),
('a0000000-0000-0000-0000-000000000001', 'SUPP-004', 'Chemical Suppliers Ltd', 'India', 'INR', 30, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'SUPP-005', 'Metal Works Singapore', 'Singapore', 'SGD', 30, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'SUPP-006', 'British Industrial Parts', 'UK', 'GBP', 45, 'MEDIUM'),
('a0000000-0000-0000-0000-000000000001', 'SUPP-007', 'Tata Steel Ltd', 'India', 'INR', 30, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'SUPP-008', 'JSPL Raw Materials', 'India', 'INR', 45, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'SUPP-009', 'Asian Plastic Co Ltd', 'China', 'USD', 30, 'HIGH'),
('a0000000-0000-0000-0000-000000000001', 'SUPP-010', 'Swiss Precision Parts AG', 'Switzerland', 'CHF', 30, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'SUPP-011', 'Energy Solutions Corp', 'USA', 'USD', 15, 'MEDIUM'),
('a0000000-0000-0000-0000-000000000001', 'SUPP-012', 'Logistics Partners Pvt', 'India', 'INR', 15, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'SUPP-013', 'IT Services GmbH', 'Germany', 'EUR', 30, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'SUPP-014', 'Maintenance Corp USA', 'USA', 'USD', 30, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'SUPP-015', 'Packaging Solutions Ltd', 'India', 'INR', 30, 'LOW');

-- ============================================================
-- COMPANY SETTINGS
-- ============================================================

INSERT INTO company_settings (company_id, setting_key, setting_value, setting_type, description) VALUES
('a0000000-0000-0000-0000-000000000001', 'min_cash_threshold', '50000000', 'NUMBER', 'Minimum cash threshold in INR before liquidity warning'),
('a0000000-0000-0000-0000-000000000001', 'target_cash_buffer', '100000000', 'NUMBER', 'Target cash buffer in INR'),
('a0000000-0000-0000-0000-000000000001', 'emergency_cash_threshold', '25000000', 'NUMBER', 'Emergency cash threshold in INR'),
('a0000000-0000-0000-0000-000000000001', 'payment_approval_l1_limit', '1000000', 'NUMBER', 'Payment amount requiring L1 approval (INR)'),
('a0000000-0000-0000-0000-000000000001', 'payment_approval_l2_limit', '10000000', 'NUMBER', 'Payment amount requiring L2 approval (INR)'),
('a0000000-0000-0000-0000-000000000001', 'payment_approval_senior_limit', '100000000', 'NUMBER', 'Payment amount requiring senior approval (INR)'),
('a0000000-0000-0000-0000-000000000001', 'fx_rate_staleness_minutes', '60', 'NUMBER', 'FX rate considered stale after N minutes'),
('a0000000-0000-0000-0000-000000000001', 'forecast_horizon_days', '90', 'NUMBER', 'Default forecast horizon in days'),
('a0000000-0000-0000-0000-000000000001', 'invoice_overdue_alert_days', '7', 'NUMBER', 'Alert when invoice overdue more than N days'),
('a0000000-0000-0000-0000-000000000001', 'ar_concentration_threshold', '25', 'NUMBER', 'Alert when single customer > N% of AR'),
('a0000000-0000-0000-0000-000000000001', 'ap_concentration_threshold', '30', 'NUMBER', 'Alert when single supplier > N% of AP');

-- ============================================================
-- RISK LIMITS
-- ============================================================

INSERT INTO risk_limits (company_id, risk_type, currency, limit_amount, warning_amount, limit_currency, description) VALUES
('a0000000-0000-0000-0000-000000000001', 'FX_EXPOSURE', 'USD', 500000000.00, 400000000.00, 'INR', 'USD FX exposure limit in INR'),
('a0000000-0000-0000-0000-000000000001', 'FX_EXPOSURE', 'EUR', 300000000.00, 250000000.00, 'INR', 'EUR FX exposure limit in INR'),
('a0000000-0000-0000-0000-000000000001', 'FX_EXPOSURE', 'GBP', 200000000.00, 160000000.00, 'INR', 'GBP FX exposure limit in INR'),
('a0000000-0000-0000-0000-000000000001', 'LIQUIDITY', NULL, 50000000.00, 75000000.00, 'INR', 'Minimum liquidity threshold'),
('a0000000-0000-0000-0000-000000000001', 'CONCENTRATION', NULL, 25.00, 20.00, 'INR', 'Max single customer AR concentration (%)');

-- ============================================================
-- PREDEFINED SCENARIOS
-- ============================================================

INSERT INTO scenarios (id, company_id, name, description, scenario_type, parameters, is_predefined) VALUES
('f0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001',
    'Baseline', 'Normal business conditions based on current trends',
    'NORMAL', '{"description": "No adjustments applied"}', TRUE),

('f0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001',
    'Delayed Customer Collections', 'Customer payments delayed by 15 days on average',
    'DELAYED_COLLECTIONS', '{"delay_days": 15, "applies_to": "all_customers"}', TRUE),

('f0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001',
    'USD Appreciation (+5%)', 'USD appreciates 5% against INR',
    'FX_APPRECIATION', '{"currency": "USD", "change_percent": 5}', TRUE),

('f0000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000001',
    'USD Depreciation (-5%)', 'USD depreciates 5% against INR',
    'FX_DEPRECIATION', '{"currency": "USD", "change_percent": -5}', TRUE),

('f0000000-0000-0000-0000-000000000005', 'a0000000-0000-0000-0000-000000000001',
    'Inventory Build-Up', 'Inventory increases by 10%, requiring additional working capital',
    'CUSTOM', '{"inventory_change_percent": 10, "working_capital_impact": "NEGATIVE"}', TRUE),

('f0000000-0000-0000-0000-000000000006', 'a0000000-0000-0000-0000-000000000001',
    'Supplier Payment Acceleration', 'Supplier payment terms reduced by 15 days',
    'CUSTOM', '{"dpo_change_days": -15, "applies_to": "all_suppliers"}', TRUE),

('f0000000-0000-0000-0000-000000000007', 'a0000000-0000-0000-0000-000000000001',
    'Combined Stress Scenario', 'USD +5% + Collections delayed 15 days + Revenue -10%',
    'STRESS', '{"fx_shock": {"currency": "USD", "change_percent": 5}, "collection_delay_days": 15, "revenue_change_percent": -10}', TRUE);

-- ============================================================
-- APPROVAL WORKFLOWS
-- ============================================================

INSERT INTO approval_workflows (company_id, name, resource_type, min_amount, max_amount, currency, required_levels) VALUES
('a0000000-0000-0000-0000-000000000001', 'Small Payment Approval',    'PAYMENT', 0, 1000000.00, 'INR', 1),
('a0000000-0000-0000-0000-000000000001', 'Medium Payment Approval',   'PAYMENT', 1000001.00, 10000000.00, 'INR', 2),
('a0000000-0000-0000-0000-000000000001', 'Large Payment Approval',    'PAYMENT', 10000001.00, NULL, 'INR', 3),
('a0000000-0000-0000-0000-000000000001', 'FX Deal Approval',          'FX_DEAL', 0, NULL, 'INR', 2),
('a0000000-0000-0000-0000-000000000001', 'Risk Limit Change',         'RISK_LIMIT', 0, NULL, 'INR', 2),
('a0000000-0000-0000-0000-000000000001', 'Invoice Write-Off',         'WRITE_OFF', 0, NULL, 'INR', 2);

-- ============================================================
-- GENERATE INVOICES (Using PL/pgSQL for bulk realistic data)
-- ============================================================

DO $$
DECLARE
    customer_ids UUID[];
    entity_ids UUID[];
    cust_id UUID;
    ent_id UUID;
    inv_date DATE;
    due_date DATE;
    amount DECIMAL(20,4);
    currencies CHAR(3)[];
    curr CHAR(3);
    status VARCHAR(30);
    inv_num INT;
    paid DECIMAL(20,4);
BEGIN
    -- Get customer IDs
    SELECT ARRAY(SELECT id FROM customers WHERE company_id = 'a0000000-0000-0000-0000-000000000001' LIMIT 20)
    INTO customer_ids;

    -- Get entity IDs
    entity_ids := ARRAY[
        'b0000000-0000-0000-0000-000000000001'::UUID,
        'b0000000-0000-0000-0000-000000000002'::UUID,
        'b0000000-0000-0000-0000-000000000003'::UUID
    ];

    currencies := ARRAY['INR', 'USD', 'EUR', 'GBP', 'JPY'];

    FOR inv_num IN 1..200 LOOP
        cust_id := customer_ids[1 + (floor(random() * array_length(customer_ids, 1)))::INT];
        ent_id := entity_ids[1 + (floor(random() * 3))::INT];

        -- Random invoice date in last 6 months
        inv_date := CURRENT_DATE - (floor(random() * 180) || ' days')::INTERVAL;
        due_date := inv_date + (floor(30 + random() * 30) || ' days')::INTERVAL;

        -- Random amount between 1L and 50L INR equivalent
        amount := floor(100000 + random() * 4900000) / 100.0 * 100;

        -- Currency from customer's currency
        SELECT c.currency INTO curr FROM customers c WHERE c.id = cust_id;

        -- Determine status based on dates
        IF due_date > CURRENT_DATE THEN
            IF random() > 0.7 THEN
                status := 'PAID';
                paid := amount;
            ELSIF random() > 0.5 THEN
                status := 'PARTIALLY_PAID';
                paid := amount * (0.3 + random() * 0.5);
            ELSE
                status := 'ISSUED';
                paid := 0;
            END IF;
        ELSE
            IF random() > 0.6 THEN
                status := 'PAID';
                paid := amount;
            ELSIF random() > 0.3 THEN
                status := 'OVERDUE';
                paid := 0;
            ELSIF random() > 0.5 THEN
                status := 'PARTIALLY_PAID';
                paid := amount * 0.5;
            ELSE
                status := 'DISPUTED';
                paid := 0;
            END IF;
        END IF;

        INSERT INTO invoices (
            company_id, entity_id, customer_id, invoice_number, invoice_date, due_date,
            currency, subtotal, tax_amount, total_amount, paid_amount,
            base_currency, base_amount, exchange_rate, rate_date,
            status, created_by
        ) VALUES (
            'a0000000-0000-0000-0000-000000000001',
            ent_id,
            cust_id,
            'INV-2024-' || LPAD(inv_num::TEXT, 5, '0'),
            inv_date,
            due_date,
            curr,
            amount,
            amount * 0.18,
            amount * 1.18,
            paid * 1.18,
            'INR',
            amount * 1.18 * CASE curr WHEN 'USD' THEN 83.42 WHEN 'EUR' THEN 90.15 WHEN 'GBP' THEN 105.32 WHEN 'JPY' THEN 0.5523 ELSE 1.0 END,
            CASE curr WHEN 'USD' THEN 83.42 WHEN 'EUR' THEN 90.15 WHEN 'GBP' THEN 105.32 WHEN 'JPY' THEN 0.5523 ELSE 1.0 END,
            inv_date,
            status,
            'd0000000-0000-0000-0000-000000000004'
        );
    END LOOP;
END;
$$;

-- ============================================================
-- GENERATE BILLS (200 bills)
-- ============================================================

DO $$
DECLARE
    supplier_ids UUID[];
    supp_id UUID;
    ent_id UUID;
    inv_date DATE;
    due_date DATE;
    amount DECIMAL(20,4);
    curr CHAR(3);
    status VARCHAR(30);
    bill_num INT;
    paid DECIMAL(20,4);
BEGIN
    SELECT ARRAY(SELECT id FROM suppliers WHERE company_id = 'a0000000-0000-0000-0000-000000000001' LIMIT 15)
    INTO supplier_ids;

    FOR bill_num IN 1..200 LOOP
        supp_id := supplier_ids[1 + (floor(random() * array_length(supplier_ids, 1)))::INT];
        ent_id := 'b0000000-0000-0000-0000-000000000001';

        inv_date := CURRENT_DATE - (floor(random() * 150) || ' days')::INTERVAL;
        due_date := inv_date + (floor(15 + random() * 45) || ' days')::INTERVAL;
        amount := floor(50000 + random() * 2000000) / 100.0 * 100;

        SELECT s.currency INTO curr FROM suppliers s WHERE s.id = supp_id;

        IF due_date > CURRENT_DATE THEN
            IF random() > 0.6 THEN
                status := 'PAID'; paid := amount;
            ELSIF random() > 0.4 THEN
                status := 'APPROVED'; paid := 0;
            ELSE
                status := 'SCHEDULED'; paid := 0;
            END IF;
        ELSE
            IF random() > 0.5 THEN
                status := 'PAID'; paid := amount;
            ELSIF random() > 0.3 THEN
                status := 'OVERDUE'; paid := 0;
            ELSE
                status := 'PARTIALLY_PAID'; paid := amount * 0.5;
            END IF;
        END IF;

        INSERT INTO bills (
            company_id, entity_id, supplier_id, bill_number, invoice_date, due_date,
            currency, subtotal, tax_amount, total_amount, paid_amount,
            base_currency, base_amount, exchange_rate, rate_date,
            status, created_by
        ) VALUES (
            'a0000000-0000-0000-0000-000000000001',
            ent_id,
            supp_id,
            'BILL-2024-' || LPAD(bill_num::TEXT, 5, '0'),
            inv_date,
            due_date,
            curr,
            amount,
            amount * 0.18,
            amount * 1.18,
            paid * 1.18,
            'INR',
            amount * 1.18 * CASE curr WHEN 'USD' THEN 83.42 WHEN 'EUR' THEN 90.15 WHEN 'GBP' THEN 105.32 WHEN 'JPY' THEN 0.5523 ELSE 1.0 END,
            CASE curr WHEN 'USD' THEN 83.42 WHEN 'EUR' THEN 90.15 WHEN 'GBP' THEN 105.32 WHEN 'JPY' THEN 0.5523 ELSE 1.0 END,
            inv_date,
            status,
            'd0000000-0000-0000-0000-000000000005'
        );
    END LOOP;
END;
$$;

-- ============================================================
-- INITIAL ALERTS (Realistic sample)
-- ============================================================

INSERT INTO alerts (company_id, alert_type, severity, title, description, source, affected_amount, currency, recommended_action) VALUES
('a0000000-0000-0000-0000-000000000001', 'INVOICE_OVERDUE', 'HIGH',
    'Large Invoice Overdue — Global Tech Solutions',
    'Invoice INV-2024-00043 for $287,500 USD from Global Tech Solutions is 47 days overdue.',
    'AR_ENGINE', 287500.00, 'USD',
    'Escalate to senior relationship manager. Review credit limit utilization.'),

('a0000000-0000-0000-0000-000000000001', 'FX_LIMIT', 'HIGH',
    'USD Exposure Approaching Risk Limit',
    'USD net exposure is ₹4.28 Cr, 85.6% of the configured risk limit of ₹5 Cr.',
    'FX_ENGINE', 42800000.00, 'INR',
    'Review open USD receivables and consider hedging strategy.'),

('a0000000-0000-0000-0000-000000000001', 'FORECAST_SHORTAGE', 'MEDIUM',
    'Projected Cash Below Threshold in 18 Days',
    'Cash flow forecast shows projected balance of ₹4.2 Cr in 18 days, below minimum threshold of ₹5 Cr.',
    'FORECAST_ENGINE', 42000000.00, 'INR',
    'Accelerate collections or arrange short-term credit facility.'),

('a0000000-0000-0000-0000-000000000001', 'PAYMENT_DUE', 'MEDIUM',
    'Large Supplier Payment Due in 3 Days',
    'Payment to European Steel AG of €185,000 EUR due in 3 days.',
    'AP_ENGINE', 185000.00, 'EUR',
    'Ensure sufficient EUR balance or prepare payment instruction.'),

('a0000000-0000-0000-0000-000000000001', 'RECONCILIATION_EXCEPTION', 'LOW',
    'Bank Statement Reconciliation Exceptions',
    '7 HDFC Bank transactions unmatched for 48+ hours.',
    'RECONCILIATION_ENGINE', NULL, NULL,
    'Review unmatched transactions and match or create manual reconciliation.');

-- ============================================================
-- FX EXPOSURES (Current)
-- ============================================================

INSERT INTO fx_exposures (company_id, entity_id, exposure_date, currency, exposure_type,
    gross_exposure, hedged_amount, base_currency, base_value, exchange_rate, risk_level) VALUES
('a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001',
    CURRENT_DATE, 'USD', 'RECEIVABLE', 6200000, 2000000, 'INR', 517204000, 83.42, 'HIGH'),
('a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001',
    CURRENT_DATE, 'USD', 'PAYABLE', 4800000, 1500000, 'INR', 400416000, 83.42, 'HIGH'),
('a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001',
    CURRENT_DATE, 'EUR', 'RECEIVABLE', 3100000, 1000000, 'INR', 279465000, 90.15, 'MEDIUM'),
('a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001',
    CURRENT_DATE, 'EUR', 'PAYABLE', 2200000, 800000, 'INR', 198330000, 90.15, 'MEDIUM'),
('a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001',
    CURRENT_DATE, 'GBP', 'RECEIVABLE', 1500000, 0, 'INR', 157980000, 105.32, 'MEDIUM'),
('a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001',
    CURRENT_DATE, 'GBP', 'PAYABLE', 800000, 0, 'INR', 84256000, 105.32, 'LOW'),
('a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001',
    CURRENT_DATE, 'JPY', 'RECEIVABLE', 850000000, 0, 'INR', 469455000, 0.5523, 'MEDIUM');
