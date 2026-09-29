import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { setAccessToken, authApi, User } from '../api'

interface AuthState {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean

  login: (email: string, password: string) => Promise<void>
  logout: () => Promise<void>
  refreshSession: () => Promise<void>
  hasPermission: (resource: string, action: string) => boolean
  hasRole: (role: string) => boolean
}

// Permission mapping matching backend rolePermissions
const ROLE_PERMISSIONS: Record<string, string[]> = {
  CFO: [
    'dashboard:read', 'cash:read', 'cash:export', 'bank_accounts:read',
    'customers:read', 'suppliers:read', 'invoices:read', 'invoices:export',
    'invoices:approve', 'bills:read', 'bills:approve', 'payments:read',
    'payments:approve', 'fx:read', 'fx:approve', 'forecast:read',
    'risk:read', 'risk:manage', 'scenarios:read', 'scenarios:create',
    'reconciliation:read', 'alerts:read', 'alerts:manage',
    'reports:read', 'reports:create', 'audit_logs:read', 'working_capital:read',
  ],
  TREASURY_MANAGER: [
    'dashboard:read', 'cash:read', 'bank_accounts:read', 'bank_accounts:create',
    'fx:read', 'fx:create', 'fx:approve', 'forecast:read', 'forecast:create',
    'scenarios:read', 'scenarios:create', 'alerts:read', 'reports:read',
    'reconciliation:read', 'reconciliation:create', 'working_capital:read',
  ],
  FINANCE_MANAGER: [
    'dashboard:read', 'customers:read', 'customers:create', 'customers:update',
    'suppliers:read', 'suppliers:create', 'suppliers:update',
    'invoices:read', 'invoices:create', 'invoices:update', 'invoices:approve',
    'bills:read', 'bills:create', 'bills:update', 'bills:approve',
    'payments:read', 'payments:create', 'payments:approve',
    'reconciliation:read', 'working_capital:read', 'reports:read',
  ],
  AR_USER: [
    'customers:read', 'customers:create', 'invoices:read', 'invoices:create',
    'invoices:update', 'alerts:read', 'working_capital:read',
  ],
  AP_USER: [
    'suppliers:read', 'suppliers:create', 'bills:read', 'bills:create',
    'payments:read', 'payments:create', 'alerts:read',
  ],
  RISK_ANALYST: [
    'dashboard:read', 'fx:read', 'risk:read', 'risk:manage',
    'scenarios:read', 'scenarios:create', 'alerts:read', 'working_capital:read',
  ],
  AUDITOR: [
    'dashboard:read', 'cash:read', 'bank_accounts:read', 'customers:read',
    'suppliers:read', 'invoices:read', 'bills:read', 'payments:read',
    'fx:read', 'reconciliation:read', 'audit_logs:read', 'alerts:read',
    'reports:read', 'working_capital:read',
  ],
  ADMIN: ['*:*'],
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      isAuthenticated: false,
      isLoading: false,

      login: async (email, password) => {
        set({ isLoading: true })
        try {
          const response = await authApi.login(email, password)
          const { access_token, refresh_token, user } = response.data.data!

          setAccessToken(access_token)
          localStorage.setItem('refresh_token', refresh_token)

          set({ user, isAuthenticated: true, isLoading: false })
        } catch (error) {
          set({ isLoading: false })
          throw error
        }
      },

      logout: async () => {
        const refreshToken = localStorage.getItem('refresh_token')
        try {
          if (refreshToken) {
            await authApi.logout(refreshToken)
          }
        } catch {
          // Silent fail on logout
        }
        setAccessToken(null)
        localStorage.removeItem('refresh_token')
        set({ user: null, isAuthenticated: false })
      },

      refreshSession: async () => {
        const refreshToken = localStorage.getItem('refresh_token')
        if (!refreshToken) {
          set({ user: null, isAuthenticated: false })
          return
        }
        try {
          const response = await authApi.refresh(refreshToken)
          const { access_token, refresh_token: newRefresh, user } = response.data.data!
          setAccessToken(access_token)
          localStorage.setItem('refresh_token', newRefresh)
          set({ user, isAuthenticated: true })
        } catch {
          setAccessToken(null)
          localStorage.removeItem('refresh_token')
          set({ user: null, isAuthenticated: false })
        }
      },

      hasPermission: (resource, action) => {
        const { user } = get()
        if (!user) return false

        const key = `${resource}:${action}`

        for (const role of user.roles) {
          const perms = ROLE_PERMISSIONS[role.code] || []
          if (perms.includes('*:*') || perms.includes(key) || perms.includes(`${resource}:*`)) {
            return true
          }
        }
        return false
      },

      hasRole: (roleCode) => {
        const { user } = get()
        if (!user) return false
        return user.roles.some((r) => r.code === roleCode || r.code === 'ADMIN')
      },
    }),
    {
      name: 'finsight-auth',
      partialize: (state) => ({
        user: state.user,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)
