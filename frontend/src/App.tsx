import { Routes, Route, Navigate, useParams } from 'react-router-dom'
import { useEffect } from 'react'
import { useAuthStore } from './stores/authStore'
import AppLayout from './components/layout/AppLayout'
import PublicLayout from './components/layout/PublicLayout'

// App pages (existing — preserved)
import LoginPage from './pages/LoginPage'
import DashboardPage from './pages/dashboard/DashboardPage'
import CashPage from './pages/cash/CashPage'
import BanksPage from './pages/banks/BanksPage'
import CustomersPage from './pages/customers/CustomersPage'
import SuppliersPage from './pages/suppliers/SuppliersPage'
import InvoicesPage from './pages/invoices/InvoicesPage'
import ReceivablesPage from './pages/receivables/ReceivablesPage'
import PayablesPage from './pages/payables/PayablesPage'
import BillsPage from './pages/bills/BillsPage'
import PaymentsPage from './pages/payments/PaymentsPage'
import WorkingCapitalPage from './pages/working-capital/WorkingCapitalPage'
import ForecastPage from './pages/forecast/ForecastPage'
import FXPage from './pages/fx/FXPage'
import FXExposurePage from './pages/fx/FXExposurePage'
import FXDealsPage from './pages/fx/FXDealsPage'
import RiskPage from './pages/risk/RiskPage'
import ScenariosPage from './pages/scenarios/ScenariosPage'
import ReconciliationPage from './pages/reconciliation/ReconciliationPage'
import AlertsPage from './pages/alerts/AlertsPage'
import ReportsPage from './pages/reports/ReportsPage'
import AuditPage from './pages/audit/AuditPage'
import AdminPage from './pages/admin/AdminPage'
import SystemPage from './pages/system/SystemPage'

// Public website pages (new)
import HomePage from './pages/public/HomePage'
import SolutionsIndexPage, { SolutionDetailPage } from './pages/public/SolutionsPage'
import PlatformPage from './pages/public/PlatformPage'
import SecurityPage from './pages/public/SecurityPage'
import AboutPage from './pages/public/AboutPage'
import InsightsPage, { InsightArticlePage } from './pages/public/InsightsPage'
import ContactPage from './pages/public/ContactPage'
import DemoPage from './pages/public/DemoPage'
import {
  PrivacyPage,
  TermsPage,
  CookiesPage,
  AccessibilityPage,
  NotFoundPage,
} from './pages/public/LegalPages'

// Guard for protected routes
function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuthStore()
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }
  return <>{children}</>
}

export default function App() {
  const { isAuthenticated, refreshSession } = useAuthStore()

  // On app start, restore session from stored refresh token
  useEffect(() => {
    const refreshToken = localStorage.getItem('refresh_token')
    if (refreshToken && !isAuthenticated) {
      refreshSession()
    }
  }, [])

  return (
    <Routes>
      {/* ── Login ── */}
      <Route path="/login" element={
        isAuthenticated ? <Navigate to="/dashboard" replace /> : <LoginPage />
      } />

      {/* ── Protected application routes ── (no layout conflict with public routes) */}
      <Route element={
        <ProtectedRoute>
          <AppLayout />
        </ProtectedRoute>
      }>
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/cash" element={<CashPage />} />
        <Route path="/banks" element={<BanksPage />} />
        <Route path="/customers" element={<CustomersPage />} />
        <Route path="/suppliers" element={<SuppliersPage />} />
        <Route path="/invoices" element={<InvoicesPage />} />
        <Route path="/receivables" element={<ReceivablesPage />} />
        <Route path="/payables" element={<PayablesPage />} />
        <Route path="/bills" element={<BillsPage />} />
        <Route path="/payments" element={<PaymentsPage />} />
        <Route path="/working-capital" element={<WorkingCapitalPage />} />
        <Route path="/forecast" element={<ForecastPage />} />
        <Route path="/fx" element={<FXPage />} />
        <Route path="/fx/exposure" element={<FXExposurePage />} />
        <Route path="/fx/deals" element={<FXDealsPage />} />
        <Route path="/risk" element={<RiskPage />} />
        <Route path="/scenarios" element={<ScenariosPage />} />
        <Route path="/reconciliation" element={<ReconciliationPage />} />
        <Route path="/alerts" element={<AlertsPage />} />
        <Route path="/reports" element={<ReportsPage />} />
        <Route path="/audit" element={<AuditPage />} />
        <Route path="/admin" element={<AdminPage />} />
        <Route path="/system-health" element={<SystemPage />} />
      </Route>

      {/* ── Public website routes ── wrapped in PublicLayout (nav + footer) ── */}
      <Route element={<PublicLayout />}>
        <Route path="/" element={<HomePage />} />
        <Route path="/solutions" element={<SolutionsIndexPage />} />
        <Route path="/solutions/working-capital" element={<SolutionDetailPage slug="working-capital" />} />
        <Route path="/solutions/cash-flow" element={<SolutionDetailPage slug="cash-flow" />} />
        <Route path="/solutions/fx" element={<SolutionDetailPage slug="fx" />} />
        <Route path="/solutions/risk" element={<SolutionDetailPage slug="risk" />} />
        <Route path="/platform" element={<PlatformPage />} />
        <Route path="/insights" element={<InsightsPage />} />
        <Route path="/insights/:slug" element={<InsightArticleRoute />} />
        <Route path="/security" element={<SecurityPage />} />
        <Route path="/about" element={<AboutPage />} />
        <Route path="/contact" element={<ContactPage />} />
        <Route path="/demo" element={<DemoPage />} />
        <Route path="/privacy" element={<PrivacyPage />} />
        <Route path="/terms" element={<TermsPage />} />
        <Route path="/cookies" element={<CookiesPage />} />
        <Route path="/accessibility" element={<AccessibilityPage />} />
        {/* 404 catch-all within public layout */}
        <Route path="*" element={<NotFoundPage />} />
      </Route>
    </Routes>
  )
}

// Wrapper to pull the slug from URL params
function InsightArticleRoute() {
  const { slug } = useParams<{ slug: string }>()
  return <InsightArticlePage slug={slug || ''} />
}
