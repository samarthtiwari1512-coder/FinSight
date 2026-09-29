import axios, { AxiosError, AxiosInstance, InternalAxiosRequestConfig } from 'axios';

const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

// ─── Token Management ───────────────────────────────────────────────────────

let accessToken: string | null = null;
let refreshPromise: Promise<string | null> | null = null;

export const setAccessToken = (token: string | null) => {
  accessToken = token;
};

export const getAccessToken = () => accessToken;

// ─── Axios Instance ──────────────────────────────────────────────────────────

const api: AxiosInstance = axios.create({
  baseURL: `${BASE_URL}/api/v1`,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
});

// ─── Request Interceptor: Attach token ───────────────────────────────────────

api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    if (accessToken) {
      config.headers.Authorization = `Bearer ${accessToken}`;
    }

    // Idempotency key for mutations
    if (['post', 'put', 'patch'].includes(config.method || '')) {
      if (!config.headers['X-Idempotency-Key']) {
        config.headers['X-Idempotency-Key'] = generateIdempotencyKey();
      }
    }

    return config;
  },
  (error) => Promise.reject(error)
);

// ─── Response Interceptor: Token refresh ─────────────────────────────────────

api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        // Deduplicate refresh calls
        if (!refreshPromise) {
          refreshPromise = refreshAccessToken();
        }
        const newToken = await refreshPromise;
        refreshPromise = null;

        if (newToken) {
          originalRequest.headers.Authorization = `Bearer ${newToken}`;
          return api(originalRequest);
        }
      } catch {
        refreshPromise = null;
        // Token refresh failed — redirect to login
        window.location.href = '/login';
      }
    }

    return Promise.reject(error);
  }
);

async function refreshAccessToken(): Promise<string | null> {
  const refreshToken = localStorage.getItem('refresh_token');
  if (!refreshToken) {
    return null;
  }

  const response = await axios.post(`${BASE_URL}/api/v1/auth/refresh`, {
    refresh_token: refreshToken,
  });

  const { access_token, refresh_token: newRefreshToken } = response.data.data;
  setAccessToken(access_token);
  localStorage.setItem('refresh_token', newRefreshToken);
  return access_token;
}

function generateIdempotencyKey(): string {
  return `${Date.now()}-${Math.random().toString(36).slice(2)}`;
}

// ─── Typed API response handler ───────────────────────────────────────────────

export interface APIResponse<T> {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
    request_id: string;
  };
  meta?: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
  };
}

// ─── Auth API ────────────────────────────────────────────────────────────────

export const authApi = {
  login: (email: string, password: string) =>
    api.post<APIResponse<LoginResponse>>('/auth/login', { email, password }),

  refresh: (refreshToken: string) =>
    api.post<APIResponse<LoginResponse>>('/auth/refresh', { refresh_token: refreshToken }),

  logout: (refreshToken: string) =>
    api.post('/auth/logout', { refresh_token: refreshToken }),

  me: () =>
    api.get<APIResponse<User>>('/auth/me'),
};

// ─── Dashboard API ────────────────────────────────────────────────────────────

export const dashboardApi = {
  getSummary: () =>
    api.get<APIResponse<DashboardSummary>>('/dashboard/summary'),
};

// ─── Cash API ────────────────────────────────────────────────────────────────

export const cashApi = {
  getPosition: (params?: CashPositionParams) =>
    api.get<APIResponse<CashPosition>>('/cash/position', { params }),

  getPositionByBank: () =>
    api.get<APIResponse<BankBalance[]>>('/cash/position/by-bank'),

  getPositionByCurrency: () =>
    api.get<APIResponse<CurrencyBalance[]>>('/cash/position/by-currency'),

  getPositionByEntity: () =>
    api.get<APIResponse<EntityBalance[]>>('/cash/position/by-entity'),

  getMovements: (params?: PaginationParams & { from?: string; to?: string }) =>
    api.get<APIResponse<CashMovement[]>>('/cash/movements', { params }),
};

// ─── Customers API ────────────────────────────────────────────────────────────

export const customersApi = {
  list: (params?: ListParams) =>
    api.get<APIResponse<Customer[]>>('/customers', { params }),

  get: (id: string) =>
    api.get<APIResponse<Customer>>(`/customers/${id}`),

  create: (data: CreateCustomerRequest) =>
    api.post<APIResponse<Customer>>('/customers', data),

  update: (id: string, data: Partial<CreateCustomerRequest>) =>
    api.put<APIResponse<Customer>>(`/customers/${id}`, data),

  getInvoices: (id: string, params?: ListParams) =>
    api.get<APIResponse<Invoice[]>>(`/customers/${id}/invoices`, { params }),

  getConcentration: () =>
    api.get<APIResponse<ConcentrationData>>('/customers/concentration'),
};

// ─── Suppliers API ───────────────────────────────────────────────────────────

export const suppliersApi = {
  list: (params?: ListParams) =>
    api.get<APIResponse<Supplier[]>>('/suppliers', { params }),

  get: (id: string) =>
    api.get<APIResponse<Supplier>>(`/suppliers/${id}`),

  create: (data: CreateSupplierRequest) =>
    api.post<APIResponse<Supplier>>('/suppliers', data),

  update: (id: string, data: Partial<CreateSupplierRequest>) =>
    api.put<APIResponse<Supplier>>(`/suppliers/${id}`, data),

  getConcentration: () =>
    api.get<APIResponse<ConcentrationData>>('/suppliers/concentration'),
};

// ─── Invoices API ────────────────────────────────────────────────────────────

export const invoicesApi = {
  list: (params?: InvoiceListParams) =>
    api.get<APIResponse<Invoice[]>>('/invoices', { params }),

  get: (id: string) =>
    api.get<APIResponse<Invoice>>(`/invoices/${id}`),

  create: (data: CreateInvoiceRequest) =>
    api.post<APIResponse<Invoice>>('/invoices', data),

  update: (id: string, data: Partial<CreateInvoiceRequest>) =>
    api.put<APIResponse<Invoice>>(`/invoices/${id}`, data),

  approve: (id: string) =>
    api.post<APIResponse<Invoice>>(`/invoices/${id}/approve`),

  recordPayment: (id: string, data: RecordPaymentRequest) =>
    api.post<APIResponse<InvoicePayment>>(`/invoices/${id}/payment`, data),

  dispute: (id: string, reason: string) =>
    api.post<APIResponse<Invoice>>(`/invoices/${id}/dispute`, { reason }),

  getAging: (params?: { as_of?: string }) =>
    api.get<APIResponse<AgingReport>>('/invoices/aging', { params }),

  getOverdue: (params?: ListParams) =>
    api.get<APIResponse<Invoice[]>>('/invoices/overdue', { params }),

  export: (params?: InvoiceListParams) =>
    api.get('/invoices/export', { params, responseType: 'blob' }),
};

// ─── Bills API ───────────────────────────────────────────────────────────────

export const billsApi = {
  list: (params?: BillListParams) =>
    api.get<APIResponse<Bill[]>>('/bills', { params }),

  get: (id: string) =>
    api.get<APIResponse<Bill>>(`/bills/${id}`),

  create: (data: CreateBillRequest) =>
    api.post<APIResponse<Bill>>('/bills', data),

  approve: (id: string) =>
    api.post<APIResponse<Bill>>(`/bills/${id}/approve`),

  recordPayment: (id: string, data: RecordPaymentRequest) =>
    api.post<APIResponse<BillPayment>>(`/bills/${id}/payment`, data),

  getAging: () =>
    api.get<APIResponse<AgingReport>>('/bills/aging'),
};

// ─── Working Capital API ──────────────────────────────────────────────────────

export const workingCapitalApi = {
  getKPIs: (params?: { period_days?: number }) =>
    api.get<APIResponse<WorkingCapitalKPIs>>('/working-capital/kpis', { params }),

  getTrend: (params?: { months?: number }) =>
    api.get<APIResponse<CCCTrendPoint[]>>('/working-capital/trend', { params }),

  getOpportunities: () =>
    api.get<APIResponse<WCOpportunity[]>>('/working-capital/opportunities'),
};

// ─── Forecast API ─────────────────────────────────────────────────────────────

export const forecastApi = {
  get: (params?: { horizon?: number }) =>
    api.get<APIResponse<ForecastOutput>>('/forecast', { params }),

  get30d: () => api.get<APIResponse<ForecastOutput>>('/forecast/30d'),
  get90d: () => api.get<APIResponse<ForecastOutput>>('/forecast/90d'),

  getAccuracy: () =>
    api.get<APIResponse<ForecastAccuracy>>('/forecast/accuracy'),
};

// ─── FX API ──────────────────────────────────────────────────────────────────

export const fxApi = {
  getRates: () =>
    api.get<APIResponse<FXRate[]>>('/fx/rates'),

  getRate: (from: string, to: string) =>
    api.get<APIResponse<FXRate>>(`/fx/rates/${from}/${to}`),

  getRateHistory: (params: { from: string; to: string; days?: number }) =>
    api.get<APIResponse<FXRate[]>>('/fx/rates/history', { params }),

  getExposure: () =>
    api.get<APIResponse<FXExposure[]>>('/fx/exposure'),

  getExposureSummary: () =>
    api.get<APIResponse<FXSummary>>('/fx/exposure/summary'),

  runSensitivity: (data: SensitivityRequest) =>
    api.post<APIResponse<SensitivityResult>>('/fx/sensitivity', data),

  runScenario: (data: FXScenarioRequest) =>
    api.post<APIResponse<ScenarioResult>>('/fx/scenarios', data),

  listDeals: (params?: ListParams) =>
    api.get<APIResponse<FXDeal[]>>('/fx/deals', { params }),

  createDeal: (data: CreateDealRequest) =>
    api.post<APIResponse<FXDeal>>('/fx/deals', data),

  approveDeal: (id: string) =>
    api.post<APIResponse<FXDeal>>(`/fx/deals/${id}/approve`),

  getHedgeCoverage: () =>
    api.get<APIResponse<HedgeCoverage>>('/fx/hedges/coverage'),
};

// ─── Payments API ─────────────────────────────────────────────────────────────

export const paymentsApi = {
  list: (params?: PaymentListParams) =>
    api.get<APIResponse<Payment[]>>('/payments', { params }),

  get: (id: string) =>
    api.get<APIResponse<Payment>>(`/payments/${id}`),

  create: (data: CreatePaymentRequest) =>
    api.post<APIResponse<Payment>>('/payments', data),

  approve: (id: string, reason?: string) =>
    api.post<APIResponse<Payment>>(`/payments/${id}/approve`, { reason }),

  reject: (id: string, reason: string) =>
    api.post<APIResponse<Payment>>(`/payments/${id}/reject`, { reason }),

  cancel: (id: string) =>
    api.delete<APIResponse<void>>(`/payments/${id}`),
};

// ─── Banks API ───────────────────────────────────────────────────────────────

export const banksApi = {
  list: () =>
    api.get<APIResponse<BankAccount[]>>('/banks'),

  get: (id: string) =>
    api.get<APIResponse<BankAccount>>(`/banks/${id}`),

  getTransactions: (id: string, params?: ListParams) =>
    api.get<APIResponse<BankTransaction[]>>(`/banks/${id}/transactions`, { params }),

  importStatement: (id: string, file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    return api.post<APIResponse<ImportResult>>(`/banks/${id}/import`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
  },
};

// ─── Risk API ────────────────────────────────────────────────────────────────

export const riskApi = {
  listEvents: (params?: RiskEventListParams) =>
    api.get<APIResponse<RiskEvent[]>>('/risk/events', { params }),

  getSummary: () =>
    api.get<APIResponse<RiskSummary>>('/risk/summary'),

  listLimits: () =>
    api.get<APIResponse<RiskLimit[]>>('/risk/limits'),

  acknowledgeEvent: (id: string) =>
    api.post<APIResponse<RiskEvent>>(`/risk/events/${id}/acknowledge`),
};

// ─── Scenarios API ────────────────────────────────────────────────────────────

export const scenariosApi = {
  list: () =>
    api.get<APIResponse<Scenario[]>>('/scenarios'),

  get: (id: string) =>
    api.get<APIResponse<Scenario>>(`/scenarios/${id}`),

  create: (data: CreateScenarioRequest) =>
    api.post<APIResponse<Scenario>>('/scenarios', data),

  run: (id: string) =>
    api.post<APIResponse<ScenarioResults>>(`/scenarios/${id}/run`),

  getResults: (id: string) =>
    api.get<APIResponse<ScenarioResults>>(`/scenarios/${id}/results`),

  runStressTest: (data: StressTestRequest) =>
    api.post<APIResponse<StressTestResult>>('/scenarios/stress-test', data),
};

// ─── Alerts API ───────────────────────────────────────────────────────────────

export const alertsApi = {
  list: (params?: AlertListParams) =>
    api.get<APIResponse<Alert[]>>('/alerts', { params }),

  getCounts: () =>
    api.get<APIResponse<AlertCounts>>('/alerts/counts'),

  acknowledge: (id: string) =>
    api.post(`/alerts/${id}/acknowledge`),

  dismiss: (id: string) =>
    api.post(`/alerts/${id}/dismiss`),

  listNotifications: () =>
    api.get<APIResponse<Notification[]>>('/notifications'),

  markRead: (id: string) =>
    api.post(`/notifications/${id}/read`),

  markAllRead: () =>
    api.post('/notifications/read-all'),

  getUnreadCount: () =>
    api.get<APIResponse<{ count: number }>>('/notifications/unread-count'),
};

// ─── Reconciliation API ───────────────────────────────────────────────────────

export const reconciliationApi = {
  getDashboard: () =>
    api.get<APIResponse<ReconciliationDashboard>>('/reconciliation/dashboard'),

  listUnmatched: () =>
    api.get<APIResponse<UnmatchedTransaction[]>>('/reconciliation/unmatched'),

  manualMatch: (data: ManualMatchRequest) =>
    api.post<APIResponse<Reconciliation>>('/reconciliation/match', data),
};

// ─── Audit API ────────────────────────────────────────────────────────────────

export const auditApi = {
  list: (params?: AuditListParams) =>
    api.get<APIResponse<AuditLog[]>>('/audit-logs', { params }),

  getResourceHistory: (resource: string, id: string) =>
    api.get<APIResponse<AuditLog[]>>(`/audit-logs/${resource}/${id}`),
};

// ─── Reports API ──────────────────────────────────────────────────────────────

export const reportsApi = {
  generateARAgeing: () =>
    api.post<APIResponse<Report>>('/reports/ar-aging'),

  generateAPAgeing: () =>
    api.post<APIResponse<Report>>('/reports/ap-aging'),

  generateWorkingCapital: () =>
    api.post<APIResponse<Report>>('/reports/working-capital'),

  generateCashPosition: () =>
    api.post<APIResponse<Report>>('/reports/cash-position'),

  generateFXExposure: () =>
    api.post<APIResponse<Report>>('/reports/fx-exposure'),

  download: (id: string) =>
    api.get(`/reports/${id}/download`, { responseType: 'blob' }),
};

// ─── Type Definitions ─────────────────────────────────────────────────────────

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  user: User;
}

export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  company_id: string;
  roles: Role[];
  is_active: boolean;
}

export interface Role {
  id: string;
  name: string;
  code: string;
  permissions: Permission[];
}

export interface Permission {
  resource: string;
  action: string;
}

export interface DashboardSummary {
  cash_position: CashPosition;
  working_capital: WorkingCapitalKPIs;
  ar_summary: ARSummary;
  ap_summary: APSummary;
  fx_summary: FXSummary;
  active_alerts: Alert[];
  liquidity_risk: string;
  fx_risk: string;
  receivables_risk: string;
  generated_at: string;
}

export interface CashPosition {
  as_of_date: string;
  opening_balance: number;
  total_inflows: number;
  total_outflows: number;
  closing_balance: number;
  net_movement: number;
  currency: string;
  by_bank?: BankBalance[];
  by_currency?: CurrencyBalance[];
  by_entity?: EntityBalance[];
}

export interface BankBalance {
  bank_name: string;
  account_number_masked: string;
  currency: string;
  available_balance: number;
  ledger_balance: number;
  base_balance: number;
  exchange_rate: number;
}

export interface CurrencyBalance {
  currency: string;
  amount: number;
  base_amount: number;
  exchange_rate: number;
  account_count: number;
}

export interface EntityBalance {
  entity_id: string;
  entity_name: string;
  country: string;
  balance: number;
  currency: string;
  base_balance: number;
}

export interface WorkingCapitalKPIs {
  dso: number;
  dpo: number;
  dio: number;
  ccc: number;
  total_receivables: number;
  total_payables: number;
  current_assets: number;
  current_liabilities: number;
  net_working_capital: number;
  current_ratio: number;
  quick_ratio: number;
  total_revenue: number;
  total_cogs: number;
  period: string;
  period_days: number;
  as_of_date: string;
}

export interface ARSummary {
  total_outstanding: number;
  total_overdue: number;
  invoice_count: number;
  overdue_count: number;
  aging_buckets: AgingBucket[];
  largest_customer?: string;
  concentration_pct: number;
}

export interface APSummary {
  total_outstanding: number;
  total_overdue: number;
  bill_count: number;
  overdue_count: number;
  due_in_3_days: number;
  aging_buckets: AgingBucket[];
  largest_supplier?: string;
  concentration_pct: number;
}

export interface AgingBucket {
  label: string;
  min_days: number;
  max_days?: number;
  amount: number;
  count: number;
  percent: number;
}

export interface FXSummary {
  total_gross_exposure: number;
  total_net_exposure: number;
  total_hedged: number;
  coverage_ratio: number;
  base_currency: string;
  by_currency: ExposureSummary[];
  highest_risk_currency?: string;
}

export interface ExposureSummary {
  currency: string;
  net_exposure: number;
  base_value: number;
  exchange_rate: number;
  risk_level: string;
  direction: string;
}

export interface Customer {
  id: string;
  customer_code: string;
  name: string;
  country: string;
  currency: string;
  payment_terms: number;
  credit_limit: number;
  risk_level: string;
  is_active: boolean;
  outstanding_amount?: number;
  overdue_amount?: number;
  avg_payment_days?: number;
  invoice_count?: number;
}

export interface Supplier {
  id: string;
  supplier_code: string;
  name: string;
  country: string;
  currency: string;
  payment_terms: number;
  risk_level: string;
  is_active: boolean;
  outstanding_payable?: number;
  overdue_amount?: number;
  bill_count?: number;
}

export interface Invoice {
  id: string;
  invoice_number: string;
  customer_id: string;
  customer_name?: string;
  invoice_date: string;
  due_date: string;
  currency: string;
  total_amount: number;
  paid_amount: number;
  outstanding_amount: number;
  status: InvoiceStatus;
  days_overdue?: number;
  age_bucket?: string;
}

export type InvoiceStatus =
  | 'DRAFT'
  | 'ISSUED'
  | 'PARTIALLY_PAID'
  | 'PAID'
  | 'OVERDUE'
  | 'DISPUTED'
  | 'CANCELLED'
  | 'WRITTEN_OFF';

export interface Bill {
  id: string;
  bill_number: string;
  supplier_id: string;
  supplier_name?: string;
  invoice_date: string;
  due_date: string;
  currency: string;
  total_amount: number;
  paid_amount: number;
  outstanding_amount: number;
  status: BillStatus;
  days_overdue?: number;
}

export type BillStatus =
  | 'DRAFT'
  | 'APPROVED'
  | 'SCHEDULED'
  | 'PARTIALLY_PAID'
  | 'PAID'
  | 'OVERDUE'
  | 'DISPUTED'
  | 'CANCELLED';

export interface Payment {
  id: string;
  payment_type: string;
  category: string;
  counterparty_name: string;
  amount: number;
  currency: string;
  status: PaymentStatus;
  payment_date?: string;
  created_by: string;
  created_by_name?: string;
  approvals?: PaymentApproval[];
}

export type PaymentStatus =
  | 'DRAFT'
  | 'PENDING_APPROVAL'
  | 'APPROVED'
  | 'SCHEDULED'
  | 'PROCESSING'
  | 'COMPLETED'
  | 'FAILED'
  | 'CANCELLED';

export interface PaymentApproval {
  id: string;
  approver_id: string;
  level: number;
  decision: string;
  reason?: string;
  decided_at: string;
  approver_name?: string;
}

export interface BankAccount {
  id: string;
  bank_name: string;
  account_number_masked: string;
  account_type: string;
  currency: string;
  country: string;
  account_holder: string;
  status: string;
  available_balance: number;
  ledger_balance: number;
  restricted_amount: number;
  last_sync_at?: string;
  is_primary: boolean;
}

export interface BankTransaction {
  id: string;
  transaction_date: string;
  description: string;
  reference?: string;
  amount: number;
  currency: string;
  transaction_type: string;
  balance_after?: number;
  status: string;
}

export interface FXRate {
  from_currency: string;
  to_currency: string;
  rate: number;
  source: string;
  provider?: string;
  rate_date: string;
  is_stale: boolean;
}

export interface FXExposure {
  id: string;
  currency: string;
  exposure_type: string;
  gross_exposure: number;
  hedged_amount: number;
  net_exposure: number;
  base_value?: number;
  exchange_rate?: number;
  risk_level: string;
}

export interface FXDeal {
  id: string;
  deal_number: string;
  deal_type: string;
  direction: string;
  buy_currency: string;
  sell_currency: string;
  notional: number;
  rate: number;
  trade_date: string;
  settlement_date?: string;
  maturity_date?: string;
  counterparty?: string;
  status: string;
}

export interface Alert {
  id: string;
  alert_type: string;
  severity: 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL';
  title: string;
  description: string;
  source: string;
  affected_amount?: number;
  currency?: string;
  recommended_action?: string;
  status: string;
  created_at: string;
}

export interface AlertCounts {
  total: number;
  critical: number;
  high: number;
  medium: number;
  low: number;
  unread: number;
}

export interface Notification {
  id: string;
  category: string;
  title: string;
  message: string;
  priority: number;
  link_url?: string;
  is_read: boolean;
  created_at: string;
}

export interface AuditLog {
  id: string;
  user_email?: string;
  action: string;
  resource: string;
  resource_id?: string;
  old_value?: Record<string, unknown>;
  new_value?: Record<string, unknown>;
  ip_address?: string;
  severity: string;
  created_at: string;
}

export interface RiskEvent {
  id: string;
  event_type: string;
  severity: string;
  risk_score: number;
  title: string;
  description: string;
  affected_entity?: string;
  affected_amount?: number;
  currency?: string;
  status: string;
  created_at: string;
}

export interface RiskLimit {
  id: string;
  risk_type: string;
  currency?: string;
  limit_amount: number;
  warning_amount?: number;
  limit_currency: string;
  description?: string;
  is_active: boolean;
}

export interface RiskSummary {
  liquidity_risk: string;
  fx_risk: string;
  receivables_risk: string;
  payables_risk: string;
  overall_risk: string;
  open_events: number;
  critical_events: number;
}

export interface Scenario {
  id: string;
  name: string;
  description?: string;
  scenario_type: string;
  parameters: Record<string, unknown>;
  is_predefined: boolean;
  created_at: string;
}

export interface ScenarioResults {
  scenario_id: string;
  baseline: Record<string, number>;
  scenario: Record<string, number>;
  variance: Record<string, number>;
  calculated_at: string;
  disclaimer: string;
}

export interface ForecastOutput {
  forecast_id: string;
  start_date: string;
  end_date: string;
  horizon_days: number;
  currency: string;
  opening_balance: number;
  projected_inflows: number;
  projected_outflows: number;
  net_cash_flow: number;
  projected_balance: number;
  min_projected_balance: number;
  min_balance_date?: string;
  liquidity_warning: boolean;
  liquidity_threshold: number;
  confidence: number;
  daily_items: ForecastDayItem[];
  generated_at: string;
  disclaimer: string;
}

export interface ForecastDayItem {
  date: string;
  projected_inflow: number;
  projected_outflow: number;
  net_flow: number;
  closing_balance: number;
  confidence: number;
  is_scheduled: boolean;
}

export interface ForecastAccuracy {
  mae: number;
  rmse: number;
  mape_percent: number;
  period: string;
  sample_size: number;
  note: string;
  disclaimer: string;
}

export interface WCOpportunity {
  category: string;
  title: string;
  description: string;
  detail: string;
  estimated_impact: number;
  currency: string;
  priority: string;
  disclaimer: string;
}

export interface CCCTrendPoint {
  date: string;
  dso: number;
  dpo: number;
  dio: number;
  ccc: number;
}

export interface ReconciliationDashboard {
  total_transactions: number;
  matched: number;
  unmatched: number;
  exceptions: number;
  match_rate: number;
}

export interface UnmatchedTransaction {
  id: string;
  transaction_date: string;
  description: string;
  amount: number;
  currency: string;
  transaction_type: string;
  bank_name: string;
}

export interface Reconciliation {
  id: string;
  status: string;
  match_amount: number;
  variance_amount: number;
  notes?: string;
}

export interface SensitivityRequest {
  currency: string;
  changes: number[];
}

export interface SensitivityResult {
  currency: string;
  base_currency: string;
  current_rate: number;
  rate_date: string;
  is_rate_stale: boolean;
  scenarios: SensitivityScenario[];
  disclaimer: string;
}

export interface SensitivityScenario {
  change_percent: number;
  label: string;
  new_rate: number;
  current_value: number;
  new_value: number;
  impact_amount: number;
}

export interface HedgeCoverage {
  total_exposure: number;
  hedged_amount: number;
  unhedged_amount: number;
  coverage_ratio: number;
  currency: string;
  by_currency: Array<{
    currency: string;
    exposure: number;
    hedged: number;
    coverage_ratio: number;
  }>;
}

export interface Report {
  id: string;
  report_type: string;
  name: string;
  status: string;
  file_url?: string;
  created_at: string;
  completed_at?: string;
}

export interface AgingReport {
  total_amount: number;
  currency: string;
  buckets: AgingBucket[];
  as_of_date: string;
}

export interface ImportResult {
  imported: number;
  rejected: number;
  duplicates: number;
  errors: string[];
}

export interface ConcentrationData {
  top5_percent: number;
  top10_percent: number;
  largest_name: string;
  largest_percent: number;
  items: Array<{
    name: string;
    amount: number;
    percent: number;
    rank: number;
  }>;
}

// Request types
export interface CreateCustomerRequest {
  name: string;
  country: string;
  currency: string;
  email?: string;
  phone?: string;
  payment_terms: number;
  credit_limit: number;
  risk_level: string;
}

export interface CreateSupplierRequest {
  name: string;
  country: string;
  currency: string;
  email?: string;
  phone?: string;
  payment_terms: number;
  risk_level: string;
}

export interface CreateInvoiceRequest {
  customer_id: string;
  entity_id: string;
  invoice_date: string;
  due_date: string;
  currency: string;
  line_items: InvoiceLineItem[];
  reference?: string;
  notes?: string;
}

export interface InvoiceLineItem {
  description: string;
  quantity: number;
  unit_price: number;
  tax_rate: number;
}

export interface CreateBillRequest {
  supplier_id: string;
  entity_id: string;
  invoice_date: string;
  due_date: string;
  currency: string;
  total_amount: number;
  tax_amount?: number;
  supplier_invoice_ref?: string;
  notes?: string;
}

export interface RecordPaymentRequest {
  payment_date: string;
  amount: number;
  currency: string;
  payment_method?: string;
  reference?: string;
  notes?: string;
}

export interface CreatePaymentRequest {
  payment_type: string;
  category: string;
  counterparty_name: string;
  amount: number;
  currency: string;
  payment_date?: string;
  bank_account_id?: string;
  notes?: string;
  idempotency_key?: string;
}

export interface CreateDealRequest {
  deal_type: string;
  direction: string;
  buy_currency: string;
  sell_currency: string;
  notional: number;
  rate: number;
  trade_date: string;
  settlement_date?: string;
  maturity_date?: string;
  counterparty?: string;
  notes?: string;
}

export interface CreateScenarioRequest {
  name: string;
  description?: string;
  scenario_type: string;
  parameters: Record<string, unknown>;
}

export interface StressTestRequest {
  revenue_change_percent?: number;
  collection_delay_days?: number;
  inventory_change_percent?: number;
  fx_shocks?: Array<{ currency: string; change_percent: number }>;
}

export interface StressTestResult {
  baseline: Record<string, number>;
  stressed: Record<string, number>;
  changes: Record<string, number>;
  narrative: string;
  disclaimer: string;
}

export interface FXScenarioRequest {
  currency: string;
  change_percent: number;
}

export interface ScenarioResult {
  currency: string;
  current_rate: number;
  new_rate: number;
  exposure: number;
  current_value: number;
  new_value: number;
  estimated_gain_loss: number;
  disclaimer: string;
}

export interface ManualMatchRequest {
  bank_transaction_id: string;
  ledger_reference_type: string;
  ledger_reference_id: string;
  notes?: string;
}

// Param types
export interface PaginationParams {
  page?: number;
  page_size?: number;
  sort_by?: string;
  sort_dir?: 'asc' | 'desc';
}

export interface ListParams extends PaginationParams {
  search?: string;
  is_active?: boolean;
}

export interface InvoiceListParams extends PaginationParams {
  status?: string;
  customer_id?: string;
  from_date?: string;
  to_date?: string;
  currency?: string;
  overdue_only?: boolean;
}

export interface BillListParams extends PaginationParams {
  status?: string;
  supplier_id?: string;
  from_date?: string;
  to_date?: string;
}

export interface PaymentListParams extends PaginationParams {
  status?: string;
  payment_type?: string;
}

export interface AlertListParams extends PaginationParams {
  severity?: string;
  status?: string;
}

export interface RiskEventListParams extends PaginationParams {
  severity?: string;
  status?: string;
}

export interface AuditListParams extends PaginationParams {
  resource?: string;
  action?: string;
  user_id?: string;
  from_date?: string;
  to_date?: string;
}

export interface CashPositionParams {
  period?: '1d' | '7d' | '30d' | 'mtd' | 'custom';
  from_date?: string;
  to_date?: string;
  currency?: string;
}

export interface CashMovement {
  id: string;
  bank_account_id: string;
  amount: number;
  currency: string;
  type: string;
  date: string;
  description: string;
}

export interface InvoicePayment {
  id: string;
  invoice_id: string;
  amount: number;
  currency: string;
  date: string;
}

export interface BillPayment {
  id: string;
  bill_id: string;
  amount: number;
  currency: string;
  date: string;
}

export default api;
