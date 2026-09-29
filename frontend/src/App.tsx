import { Routes, Route, Navigate } from 'react-router-dom'
import { useEffect } from 'react'
import { useAuthStore } from './stores/authStore'
import { setAccessToken } from './api'
import AppLayout from './components/layout/AppLayout'
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
      <Route path="/login" element={
        isAuthenticated ? <Navigate to="/dashboard" replace /> : <LoginPage />
      } />

      <Route path="/" element={
        <ProtectedRoute>
          <AppLayout />
        </ProtectedRoute>
      }>
        <Route index element={<Navigate to="/dashboard" replace />} />
        <Route path="dashboard" element={<DashboardPage />} />
        <Route path="cash" element={<CashPage />} />
        <Route path="banks" element={<BanksPage />} />
        <Route path="customers" element={<CustomersPage />} />
        <Route path="suppliers" element={<SuppliersPage />} />
        <Route path="invoices" element={<InvoicesPage />} />
        <Route path="receivables" element={<ReceivablesPage />} />
        <Route path="payables" element={<PayablesPage />} />
        <Route path="bills" element={<BillsPage />} />
        <Route path="payments" element={<PaymentsPage />} />
        <Route path="working-capital" element={<WorkingCapitalPage />} />
        <Route path="forecast" element={<ForecastPage />} />
        <Route path="fx" element={<FXPage />} />
        <Route path="fx/exposure" element={<FXExposurePage />} />
        <Route path="fx/deals" element={<FXDealsPage />} />
        <Route path="risk" element={<RiskPage />} />
        <Route path="scenarios" element={<ScenariosPage />} />
        <Route path="reconciliation" element={<ReconciliationPage />} />
        <Route path="alerts" element={<AlertsPage />} />
        <Route path="reports" element={<ReportsPage />} />
        <Route path="audit" element={<AuditPage />} />
        <Route path="admin" element={<AdminPage />} />
        <Route path="system-health" element={<SystemPage />} />
      </Route>

      <Route path="*" element={<Navigate to="/dashboard" replace />} />
    </Routes>
  )
}
