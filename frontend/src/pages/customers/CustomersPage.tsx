import { useQuery } from '@tanstack/react-query'
import { customersApi } from '../../api'
import { Users, Plus, AlertCircle, TrendingUp, Filter } from 'lucide-react'

export default function CustomersPage() {
  const { data: customersRes, isLoading, error } = useQuery({
    queryKey: ['customers', 'list'],
    queryFn: () => customersApi.list().then(res => res.data.data)
  })

  return (
    <div className="page fade-in">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <div>
          <h1 style={{ fontSize: 24, fontWeight: 700, color: 'var(--fs-text)', marginBottom: 4 }}>
            Customers (Accounts Receivable)
          </h1>
          <p style={{ color: 'var(--fs-text-muted)', fontSize: 14 }}>
            Manage client profiles, credit limits, and outstanding balances.
          </p>
        </div>
        <div style={{ display: 'flex', gap: 12 }}>
          <button className="btn-secondary" style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <Filter size={14} />
            <span>Filter</span>
          </button>
          <button className="btn-primary" style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <Plus size={14} />
            <span>Add Customer</span>
          </button>
        </div>
      </div>

      <div className="glass-panel" style={{ padding: 20 }}>
        {isLoading ? (
          <div style={{ padding: 40, textAlign: 'center', color: 'var(--fs-text-muted)' }}>Loading customers...</div>
        ) : error ? (
          <div style={{ padding: 40, textAlign: 'center', color: 'var(--danger)' }}>
            <AlertCircle size={24} style={{ margin: '0 auto 12px' }} />
            Failed to load customers.
          </div>
        ) : (
          <div style={{ overflowX: 'auto' }}>
            <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: 14 }}>
              <thead>
                <tr style={{ borderBottom: '1px solid var(--fs-border)', color: 'var(--fs-text-secondary)', textAlign: 'left' }}>
                  <th style={{ padding: '12px 16px', fontWeight: 600 }}>Customer Name</th>
                  <th style={{ padding: '12px 16px', fontWeight: 600 }}>Code</th>
                  <th style={{ padding: '12px 16px', fontWeight: 600 }}>Country</th>
                  <th style={{ padding: '12px 16px', fontWeight: 600, textAlign: 'right' }}>Credit Limit</th>
                  <th style={{ padding: '12px 16px', fontWeight: 600, textAlign: 'right' }}>Outstanding</th>
                  <th style={{ padding: '12px 16px', fontWeight: 600 }}>Risk Level</th>
                </tr>
              </thead>
              <tbody>
                {customersRes && customersRes.length > 0 ? (
                  customersRes.map(customer => (
                    <tr key={customer.id} style={{ borderBottom: '1px solid var(--fs-border)' }}>
                      <td style={{ padding: '12px 16px', color: 'var(--fs-text)', fontWeight: 500, display: 'flex', alignItems: 'center', gap: 8 }}>
                        <div style={{ width: 32, height: 32, borderRadius: '50%', backgroundColor: 'var(--fs-surface-hover)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--fs-text-secondary)' }}>
                          <Users size={14} />
                        </div>
                        {customer.name}
                      </td>
                      <td style={{ padding: '12px 16px', color: 'var(--fs-text-secondary)', fontFamily: 'JetBrains Mono, monospace' }}>
                        {customer.customer_code}
                      </td>
                      <td style={{ padding: '12px 16px', color: 'var(--fs-text-secondary)' }}>
                        {customer.country}
                      </td>
                      <td style={{ padding: '12px 16px', color: 'var(--fs-text)', textAlign: 'right', fontFamily: 'JetBrains Mono, monospace' }}>
                        {customer.currency} {customer.credit_limit.toLocaleString()}
                      </td>
                      <td style={{ padding: '12px 16px', color: 'var(--fs-text)', textAlign: 'right', fontFamily: 'JetBrains Mono, monospace', fontWeight: 600 }}>
                        {customer.currency} {customer.outstanding_amount?.toLocaleString() || '0'}
                      </td>
                      <td style={{ padding: '12px 16px' }}>
                        <span style={{ 
                          padding: '4px 8px', 
                          borderRadius: 4, 
                          fontSize: 12, 
                          fontWeight: 500,
                          backgroundColor: customer.risk_level === 'HIGH' ? 'rgba(239, 68, 68, 0.1)' : customer.risk_level === 'MEDIUM' ? 'rgba(245, 158, 11, 0.1)' : 'rgba(16, 185, 129, 0.1)',
                          color: customer.risk_level === 'HIGH' ? '#ef4444' : customer.risk_level === 'MEDIUM' ? '#f59e0b' : '#10b981'
                        }}>
                          {customer.risk_level}
                        </span>
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan={6} style={{ padding: '24px 16px', textAlign: 'center', color: 'var(--fs-text-muted)' }}>
                      No customers found.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}
