import { useQuery } from '@tanstack/react-query'
import { suppliersApi } from '../../api'
import { Truck, Plus, AlertCircle, Filter } from 'lucide-react'

export default function SuppliersPage() {
  const { data: suppliersRes, isLoading, error } = useQuery({
    queryKey: ['suppliers', 'list'],
    queryFn: () => suppliersApi.list().then(res => res.data.data)
  })

  return (
    <div className="page fade-in">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <div>
          <h1 style={{ fontSize: 24, fontWeight: 700, color: 'var(--fs-text)', marginBottom: 4 }}>
            Suppliers (Accounts Payable)
          </h1>
          <p style={{ color: 'var(--fs-text-muted)', fontSize: 14 }}>
            Manage vendors, payment terms, and outstanding payables.
          </p>
        </div>
        <div style={{ display: 'flex', gap: 12 }}>
          <button className="btn-secondary" style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <Filter size={14} />
            <span>Filter</span>
          </button>
          <button className="btn-primary" style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <Plus size={14} />
            <span>Add Supplier</span>
          </button>
        </div>
      </div>

      <div className="glass-panel" style={{ padding: 20 }}>
        {isLoading ? (
          <div style={{ padding: 40, textAlign: 'center', color: 'var(--fs-text-muted)' }}>Loading suppliers...</div>
        ) : error ? (
          <div style={{ padding: 40, textAlign: 'center', color: 'var(--danger)' }}>
            <AlertCircle size={24} style={{ margin: '0 auto 12px' }} />
            Failed to load suppliers.
          </div>
        ) : (
          <div style={{ overflowX: 'auto' }}>
            <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: 14 }}>
              <thead>
                <tr style={{ borderBottom: '1px solid var(--fs-border)', color: 'var(--fs-text-secondary)', textAlign: 'left' }}>
                  <th style={{ padding: '12px 16px', fontWeight: 600 }}>Supplier Name</th>
                  <th style={{ padding: '12px 16px', fontWeight: 600 }}>Code</th>
                  <th style={{ padding: '12px 16px', fontWeight: 600 }}>Country</th>
                  <th style={{ padding: '12px 16px', fontWeight: 600, textAlign: 'center' }}>Terms (Days)</th>
                  <th style={{ padding: '12px 16px', fontWeight: 600, textAlign: 'right' }}>Outstanding</th>
                  <th style={{ padding: '12px 16px', fontWeight: 600 }}>Risk Level</th>
                </tr>
              </thead>
              <tbody>
                {suppliersRes && suppliersRes.length > 0 ? (
                  suppliersRes.map(supplier => (
                    <tr key={supplier.id} style={{ borderBottom: '1px solid var(--fs-border)' }}>
                      <td style={{ padding: '12px 16px', color: 'var(--fs-text)', fontWeight: 500, display: 'flex', alignItems: 'center', gap: 8 }}>
                        <div style={{ width: 32, height: 32, borderRadius: '50%', backgroundColor: 'var(--fs-surface-hover)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--fs-text-secondary)' }}>
                          <Truck size={14} />
                        </div>
                        {supplier.name}
                      </td>
                      <td style={{ padding: '12px 16px', color: 'var(--fs-text-secondary)', fontFamily: 'JetBrains Mono, monospace' }}>
                        {supplier.supplier_code}
                      </td>
                      <td style={{ padding: '12px 16px', color: 'var(--fs-text-secondary)' }}>
                        {supplier.country}
                      </td>
                      <td style={{ padding: '12px 16px', color: 'var(--fs-text-secondary)', textAlign: 'center' }}>
                        Net {supplier.payment_terms}
                      </td>
                      <td style={{ padding: '12px 16px', color: 'var(--fs-text)', textAlign: 'right', fontFamily: 'JetBrains Mono, monospace', fontWeight: 600 }}>
                        {supplier.currency} {supplier.outstanding_payable?.toLocaleString() || '0'}
                      </td>
                      <td style={{ padding: '12px 16px' }}>
                        <span style={{ 
                          padding: '4px 8px', 
                          borderRadius: 4, 
                          fontSize: 12, 
                          fontWeight: 500,
                          backgroundColor: supplier.risk_level === 'HIGH' ? 'rgba(239, 68, 68, 0.1)' : supplier.risk_level === 'MEDIUM' ? 'rgba(245, 158, 11, 0.1)' : 'rgba(16, 185, 129, 0.1)',
                          color: supplier.risk_level === 'HIGH' ? '#ef4444' : supplier.risk_level === 'MEDIUM' ? '#f59e0b' : '#10b981'
                        }}>
                          {supplier.risk_level}
                        </span>
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan={6} style={{ padding: '24px 16px', textAlign: 'center', color: 'var(--fs-text-muted)' }}>
                      No suppliers found.
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
