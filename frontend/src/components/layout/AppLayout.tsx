import { Outlet, NavLink, useNavigate } from 'react-router-dom'
import { useState } from 'react'
import type { LucideIcon } from 'lucide-react'
import {
  LayoutDashboard, Landmark, Building2, Users, Truck, FileText,
  Receipt, ArrowLeftRight, TrendingUp, BarChart3, Globe, Shield,
  GitBranch, RefreshCcw, Bell, FileBarChart, ScrollText, Settings,
  UserCog, Activity, LogOut, ChevronDown, Search, Menu, X,
  Wallet, CreditCard, AlertTriangle, Zap
} from 'lucide-react'
import { useAuthStore } from '../../stores/authStore'
import { alertsApi } from '../../api'
import { useQuery } from '@tanstack/react-query'
import toast from 'react-hot-toast'

interface NavSection {
  label: string
  items: NavItem[]
}

interface NavItem {
  path: string
  label: string
  icon: LucideIcon
  badge?: string
}

const navSections: NavSection[] = [
  {
    label: 'Overview',
    items: [
      { path: '/dashboard', label: 'Executive Dashboard', icon: LayoutDashboard },
    ],
  },
  {
    label: 'Treasury',
    items: [
      { path: '/cash', label: 'Cash Management', icon: Wallet },
      { path: '/banks', label: 'Bank Accounts', icon: Landmark },
      { path: '/forecast', label: 'Cash Forecast', icon: TrendingUp },
    ],
  },
  {
    label: 'Receivables',
    items: [
      { path: '/receivables', label: 'AR Summary', icon: CreditCard },
      { path: '/customers', label: 'Customers', icon: Users },
      { path: '/invoices', label: 'Invoices', icon: FileText },
    ],
  },
  {
    label: 'Payables',
    items: [
      { path: '/payables', label: 'AP Summary', icon: Receipt },
      { path: '/suppliers', label: 'Suppliers', icon: Truck },
      { path: '/bills', label: 'Bills', icon: FileText },
    ],
  },
  {
    label: 'Finance',
    items: [
      { path: '/payments', label: 'Payments', icon: ArrowLeftRight },
      { path: '/working-capital', label: 'Working Capital', icon: BarChart3 },
      { path: '/reconciliation', label: 'Reconciliation', icon: RefreshCcw },
    ],
  },
  {
    label: 'FX & Risk',
    items: [
      { path: '/fx', label: 'FX Rates', icon: Globe },
      { path: '/fx/exposure', label: 'FX Exposure', icon: Shield },
      { path: '/fx/deals', label: 'FX Deals', icon: Zap },
      { path: '/risk', label: 'Risk Events', icon: AlertTriangle },
      { path: '/scenarios', label: 'Scenarios', icon: GitBranch },
    ],
  },
  {
    label: 'Intelligence',
    items: [
      { path: '/alerts', label: 'Alerts', icon: Bell },
      { path: '/reports', label: 'Reports', icon: FileBarChart },
    ],
  },
  {
    label: 'Compliance',
    items: [
      { path: '/audit', label: 'Audit Logs', icon: ScrollText },
    ],
  },
  {
    label: 'System',
    items: [
      { path: '/admin', label: 'Administration', icon: UserCog },
      { path: '/system-health', label: 'System Health', icon: Activity },
    ],
  },
]

export default function AppLayout() {
  const { user, logout } = useAuthStore()
  const navigate = useNavigate()
  const [sidebarOpen, setSidebarOpen] = useState(true)
  const [searchOpen, setSearchOpen] = useState(false)

  const { data: alertCounts } = useQuery({
    queryKey: ['alert-counts'],
    queryFn: () => alertsApi.getCounts().then(r => r.data.data),
    refetchInterval: 60000,
  })

  const handleLogout = async () => {
    await logout()
    navigate('/login')
    toast.success('Signed out successfully')
  }

  const criticalAlerts = alertCounts?.critical || 0

  return (
    <div style={{ display: 'flex', minHeight: '100vh' }}>
      {/* Sidebar */}
      {sidebarOpen && (
        <aside className="sidebar fade-in">
          {/* Logo */}
          <div style={{
            padding: '20px 20px 16px',
            borderBottom: '1px solid var(--fs-border)',
            display: 'flex',
            alignItems: 'center',
            gap: '10px',
          }}>
            <div style={{
              width: 32,
              height: 32,
              background: 'var(--fs-accent)',
              borderRadius: 8,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              fontSize: 14,
              fontWeight: 800,
              color: 'white',
              flexShrink: 0,
            }}>FS</div>
            <div>
              <div style={{ fontSize: 15, fontWeight: 700, color: 'var(--fs-text)', lineHeight: 1.2 }}>
                FinSight
              </div>
              <div style={{ fontSize: 10, color: 'var(--fs-text-muted)', letterSpacing: '0.05em' }}>
                TREASURY PLATFORM
              </div>
            </div>
          </div>

          {/* Navigation */}
          <nav style={{ flex: 1, padding: '8px 0', overflowY: 'auto' }}>
            {navSections.map((section) => (
              <div key={section.label}>
                <div className="nav-section-label">{section.label}</div>
                {section.items.map((item) => (
                  <NavLink
                    key={item.path}
                    to={item.path}
                    className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}
                  >
                    <item.icon size={15} />
                    <span style={{ flex: 1 }}>{item.label}</span>
                    {item.path === '/alerts' && criticalAlerts > 0 && (
                      <span style={{
                        background: 'var(--fs-negative)',
                        color: 'white',
                        fontSize: 10,
                        fontWeight: 700,
                        padding: '1px 6px',
                        borderRadius: 100,
                        minWidth: 18,
                        textAlign: 'center',
                      }}>
                        {criticalAlerts}
                      </span>
                    )}
                  </NavLink>
                ))}
              </div>
            ))}
          </nav>

          {/* User Footer */}
          <div style={{
            borderTop: '1px solid var(--fs-border)',
            padding: '12px',
          }}>
            <div style={{
              display: 'flex',
              alignItems: 'center',
              gap: 10,
              padding: '8px',
              borderRadius: 8,
              cursor: 'pointer',
            }}>
              <div style={{
                width: 30,
                height: 30,
                borderRadius: '50%',
                background: 'var(--fs-text)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                color: 'white',
                fontSize: 12,
                fontWeight: 700,
                flexShrink: 0,
              }}>
                {user?.first_name?.[0]}{user?.last_name?.[0]}
              </div>
              <div style={{ flex: 1, minWidth: 0 }}>
                <div style={{ fontSize: 12, fontWeight: 600, color: 'var(--fs-text)', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                  {user?.first_name} {user?.last_name}
                </div>
                <div style={{ fontSize: 10, color: 'var(--fs-text-muted)', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                  {user?.roles?.[0]?.name || 'User'}
                </div>
              </div>
              <button
                onClick={handleLogout}
                style={{
                  background: 'none',
                  border: 'none',
                  cursor: 'pointer',
                  color: 'var(--fs-text-muted)',
                  padding: 4,
                  borderRadius: 6,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                }}
                title="Sign out"
              >
                <LogOut size={14} />
              </button>
            </div>
          </div>
        </aside>
      )}

      {/* Main Area */}
      <div className="main-content" style={{ marginLeft: sidebarOpen ? 'var(--sidebar-width)' : 0, flex: 1 }}>
        {/* Header */}
        <header className="header">
          <button
            onClick={() => setSidebarOpen(!sidebarOpen)}
            style={{
              background: 'none',
              border: '1px solid var(--fs-border)',
              borderRadius: 8,
              padding: 8,
              cursor: 'pointer',
              color: 'var(--fs-text-secondary)',
              display: 'flex',
              alignItems: 'center',
            }}
          >
            {sidebarOpen ? <X size={16} /> : <Menu size={16} />}
          </button>

          {/* Search Bar */}
          <div style={{
            flex: 1,
            maxWidth: 400,
            position: 'relative',
          }}>
            <Search
              size={14}
              style={{
                position: 'absolute',
                left: 10,
                top: '50%',
                transform: 'translateY(-50%)',
                color: 'var(--fs-text-muted)',
              }}
            />
            <input
              placeholder="Search invoices, customers, payments..."
              className="input"
              style={{ paddingLeft: 32, fontSize: 13 }}
              onFocus={() => setSearchOpen(true)}
            />
          </div>

          <div style={{ marginLeft: 'auto', display: 'flex', alignItems: 'center', gap: 12 }}>
            {/* Company badge */}
            <div style={{
              background: 'transparent',
              border: '1px solid var(--fs-border)',
              borderRadius: 8,
              padding: '4px 10px',
              fontSize: 11,
              fontWeight: 600,
              color: 'var(--fs-accent)',
            }}>
              ACME GLOBAL
            </div>

            {/* FX Rate indicator */}
            <div style={{
              fontSize: 11,
              color: 'var(--fs-text-muted)',
              display: 'flex',
              alignItems: 'center',
              gap: 4,
            }}>
              <Globe size={12} />
              <span>₹83.42 / USD</span>
            </div>

            {/* Alert bell */}
            <NavLink to="/alerts" style={{ position: 'relative', color: 'var(--fs-text-muted)', display: 'flex' }}>
              <Bell size={18} />
              {criticalAlerts > 0 && (
                <span className="pulse-red" style={{
                  position: 'absolute',
                  top: -2,
                  right: -2,
                  width: 8,
                  height: 8,
                  borderRadius: '50%',
                  background: 'var(--fs-negative)',
                }} />
              )}
            </NavLink>

            {/* Settings */}
            <NavLink to="/admin" style={{ color: 'var(--fs-text-muted)', display: 'flex' }}>
              <Settings size={18} />
            </NavLink>
          </div>
        </header>

        {/* Page Content */}
        <main style={{ minHeight: 'calc(100vh - var(--header-height))' }}>
          <Outlet />
        </main>
      </div>
    </div>
  )
}
