# FinSight
**Corporate Treasury, Working Capital & FX Risk Intelligence Platform**

> A production-style enterprise treasury management system built with Go (Gin) + React (TypeScript)

---

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.22+
- Node.js 20+
- PostgreSQL 16 (or use Docker)
- Redis 7 (or use Docker)

### Run with Docker Compose

```bash
# Clone and start all services
git clone <repo>
cd FinSight
docker compose up -d

# The app will be available at:
# Frontend:   http://localhost
# Backend:    http://localhost:8080
# Prometheus: http://localhost:9090 (with --profile monitoring)
# Grafana:    http://localhost:3001 (with --profile monitoring)
```

### Run Locally (Development)

**Backend:**
```bash
cd backend
cp .env.example .env         # Edit with your config
go mod download
go run ./cmd/server/main.go
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev                  # Starts on http://localhost:5173
```

---

## 🔐 Demo Credentials

All users belong to **Acme Global Manufacturing Ltd** (INR base currency):

| Role               | Email                          | Password   | Access Level         |
|--------------------|--------------------------------|------------|----------------------|
| CFO                | amelia.chen@acmeglobal.com     | Demo@2024  | Full read + approve  |
| Treasury Manager   | treasury@acmeglobal.com        | Demo@2024  | Treasury + FX + Cash |
| Finance Manager    | finance@acmeglobal.com         | Demo@2024  | AR + AP + Payments   |
| AR User            | ar@acmeglobal.com              | Demo@2024  | Customers + Invoices |
| AP User            | ap@acmeglobal.com              | Demo@2024  | Suppliers + Bills    |
| Risk Analyst       | risk@acmeglobal.com            | Demo@2024  | FX + Risk + Scenarios|
| Auditor            | auditor@acmeglobal.com         | Demo@2024  | Read-only everywhere |
| Admin              | admin@acmeglobal.com           | Demo@2024  | Full system access   |

---

## 🏗️ Architecture

```
FinSight/
├── backend/                    # Go API server
│   ├── cmd/server/main.go      # Entry point, DI, routing
│   ├── internal/
│   │   ├── auth/               # JWT + session management
│   │   ├── config/             # Configuration management
│   │   ├── handlers/           # HTTP request handlers
│   │   ├── middleware/         # Auth, RBAC, logging, security headers
│   │   ├── models/             # Domain models (decimal arithmetic)
│   │   ├── repositories/       # Database access layer (pgx)
│   │   ├── services/           # Business logic
│   │   └── workers/            # Background jobs (FX sync, alerts)
│   ├── migrations/
│   │   ├── 001_initial_schema.sql  # 31-table normalized schema
│   │   └── 002_seed_data.sql       # Demo company + realistic data
│   └── pkg/response/           # Standardized API responses
├── frontend/                   # React 18 + TypeScript + Vite
│   ├── src/
│   │   ├── api/                # Axios client + all typed endpoints
│   │   ├── components/layout/  # AppLayout, sidebar, header
│   │   ├── pages/              # Module pages (dashboard, cash, FX, etc.)
│   │   └── stores/             # Zustand state management
│   └── ...
├── infra/
│   ├── prometheus/             # Metrics scraping config
│   └── grafana/                # Dashboard provisioning
└── docker-compose.yml
```

### Tech Stack

| Layer          | Technology                           |
|----------------|--------------------------------------|
| Backend        | Go 1.22, Gin, pgx/v5                 |
| Database       | PostgreSQL 16 (DECIMAL for money)    |
| Cache          | Redis 7                              |
| Auth           | JWT (RS256), bcrypt, refresh tokens  |
| Frontend       | React 18, TypeScript, Vite           |
| State          | Zustand + TanStack Query             |
| Charts         | Recharts                             |
| FX Provider    | ExchangeRate-API (configurable)      |
| Monitoring     | Prometheus + Grafana                 |
| Containerization| Docker, Docker Compose              |

---

## 📊 Features

### Treasury
- **Cash Position**: Real-time view by bank, currency, and entity
- **Bank Accounts**: 24 accounts across 6 entities in 4 countries
- **Cash Flow Forecast**: 7/30/60/90-day forecast with confidence intervals

### Accounts Receivable
- Invoice lifecycle (Draft → Issued → Approved → Paid)
- AR aging analysis (Current / 30 / 60 / 90+ days)
- Customer credit limit monitoring
- Collection probability modeling

### Accounts Payable
- Bill management with approval workflow
- AP aging and DPO analysis
- Payment scheduling with bank connectivity design
- Supplier concentration risk

### FX Management
- Real-time rates with staleness detection
- Net exposure by currency with risk levels
- FX deal management (spot, forward, swap)
- Sensitivity analysis and scenario modeling
- Hedge coverage tracking

### Working Capital
- DSO, DPO, DIO, CCC with trend analysis
- Net working capital and liquidity ratios
- Optimization opportunity generation

### Risk & Compliance
- Automated alerts for threshold breaches
- Risk event tracking and escalation
- Configurable risk limits by category
- Stress testing with scenario analysis

### Enterprise Controls
- **RBAC**: 8 roles with granular permissions per resource/action
- **Audit Logs**: Immutable audit trail for all mutations
- **Idempotency**: All mutations require unique idempotency keys
- **Reconciliation**: Bank statement matching with exception management

---

## 🔒 Security Design

1. **JWT Authentication**: Short-lived access tokens (8h) + rotating refresh tokens (30d)
2. **RBAC**: Permission checked per-endpoint at middleware layer
3. **Password Security**: bcrypt with configurable cost factor (default: 12)
4. **Audit Trail**: Immutable trigger-protected audit_logs table
5. **No Floating-Point Money**: All amounts stored and computed as DECIMAL/NUMERIC
6. **Idempotency**: Mutations protected with idempotency keys
7. **Rate Limiting**: Redis-backed rate limiting on auth endpoints
8. **Secure Headers**: CSP, X-Frame-Options, X-Content-Type-Options on all responses

---

## ⚠️ Disclaimers

This is a **simulation/analytics platform**:
- No real financial transactions are executed
- No real payment accounts are connected by default  
- Cash flow forecasts are statistical estimates, not guarantees
- FX scenario results are labeled as estimates, not financial advice
- Working capital opportunity descriptions include appropriate disclaimers

---

## 🛠️ Development

### Backend (Go)

```bash
# Run tests
cd backend && go test ./...

# Check with vet
go vet ./...

# Build binary
go build -o finsight ./cmd/server/main.go
```

### Frontend

```bash
cd frontend
npm run dev      # Development
npm run build    # Production build
npm run lint     # ESLint
npm test         # Vitest
```

### Database Migrations

Migrations run automatically from `backend/migrations/` on container start.
Run manually:
```bash
psql -U finsight -d finsight -f backend/migrations/001_initial_schema.sql
psql -U finsight -d finsight -f backend/migrations/002_seed_data.sql
```

---

## 📈 API Documentation

Base URL: `http://localhost:8080/api/v1`

| Category         | Endpoints |
|------------------|-----------|
| Auth             | POST /auth/login, /auth/refresh, /auth/logout, GET /auth/me |
| Dashboard        | GET /dashboard, /dashboard/summary |
| Cash             | GET /cash/position, /cash/movements, /cash/accounts |
| Banks            | GET/POST /banks, GET /banks/:id/transactions |
| Customers        | GET/POST/PUT /customers, GET /customers/concentration |
| Suppliers        | GET/POST/PUT /suppliers |
| Invoices (AR)    | GET/POST/PUT /invoices, POST /:id/approve, /:id/payment |
| Bills (AP)       | GET/POST/PUT /bills, POST /:id/approve |
| Working Capital  | GET /working-capital/kpis, /trend, /opportunities |
| Cash Forecast    | GET /forecast, /forecast/30d, /forecast/90d |
| FX               | GET /fx/rates, /fx/exposure, POST /fx/sensitivity |
| FX Deals         | GET/POST /fx/deals, POST /:id/approve |
| Payments         | GET/POST /payments, POST /:id/approve, /:id/reject |
| Risk             | GET /risk/events, /risk/limits |
| Scenarios        | GET/POST /scenarios, POST /:id/run |
| Reconciliation   | GET /reconciliation/dashboard, POST /reconciliation/match |
| Alerts           | GET /alerts, GET /alerts/counts |
| Reports          | POST /reports/*, GET /reports/:id/download |
| Audit            | GET /audit-logs |
| Admin            | GET/POST /admin/users, /admin/roles |
| Health           | GET /health, /ready |

All responses follow the standard envelope:
```json
{
  "success": true,
  "data": { ... },
  "error": null,
  "meta": { "page": 1, "page_size": 25, "total": 150, "total_pages": 6 }
}
```
