import { useQuery } from '@tanstack/react-query'
import { banksApi } from '../../api'
import { Landmark, Plus, RefreshCw, AlertCircle } from 'lucide-react'

export default function BanksPage() {
  const { data: banksRes, isLoading, error } = useQuery({
    queryKey: ['banks', 'list'],
    queryFn: () => banksApi.list().then(res => res.data.data)
  })

  return (
    <div className="page fade-in">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <div>
          <h1 style={{ fontSize: 24, fontWeight: 700, color: 'var(--text-primary)', marginBottom: 4 }}>
            Bank Accounts
          </h1>
          <p style={{ color: 'var(--text-muted)', fontSize: 14 }}>
            Manage corporate bank accounts, connections, and statement syncing.
          </p>
        </div>
        <div style={{ display: 'flex', gap: 12 }}>
          <button className="btn-secondary" style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <RefreshCw size={14} />
            <span>Sync All</span>
          </button>
          <button className="btn-primary" style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <Plus size={14} />
            <span>Add Account</span>
          </button>
        </div>
      </div>

      <div className="glass-panel" style={{ padding: 20 }}>
        {isLoading ? (
          <div style={{ padding: 40, textAlign: 'center', color: 'var(--text-muted)' }}>Loading bank accounts...</div>
        ) : error ? (
          <div style={{ padding: 40, textAlign: 'center', color: 'var(--danger)' }}>
            <AlertCircle size={24} style={{ margin: '0 auto 12px' }} />
            Failed to load bank accounts.
          </div>
        ) : (
          <div style={{ overflowX: 'auto' }}>
            <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: 14 }}>
              <thead>
                <tr style={{ borderBottom: '1px solid var(--border-light)', color: 'var(--text-secondary)', textAlign: 'left' }}>
                  <th style={{ padding: '12px 16px', fontWeight: 600 }}>Bank</th>
                  <th style={{ padding: '12px 16px', fontWeight: 600 }}>Account No.</th>
                  <th style={{ padding: '12px 16px', fontWeight: 600 }}>Type</th>
                  <th style={{ padding: '12px 16px', fontWeight: 600 }}>Currency</th>
                  <th style={{ padding: '12px 16px', fontWeight: 600, textAlign: 'right' }}>Available Balance</th>
                  <th style={{ padding: '12px 16px', fontWeight: 600 }}>Status</th>
                </tr>
              </thead>
              <tbody>
                {banksRes && banksRes.length > 0 ? (
                  banksRes.map(bank => (
                    <tr key={bank.id} style={{ borderBottom: '1px solid var(--border-light)' }}>
                      <td style={{ padding: '12px 16px', color: 'var(--text-primary)', fontWeight: 500, display: 'flex', alignItems: 'center', gap: 8 }}>
                        <div style={{ width: 32, height: 32, borderRadius: '50%', backgroundColor: 'var(--bg-tertiary)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--text-secondary)' }}>
                          <Landmark size={14} />
                        </div>
                        {bank.bank_name}
                      </td>
                      <td style={{ padding: '12px 16px', color: 'var(--text-secondary)', fontFamily: 'JetBrains Mono, monospace' }}>
                        {bank.account_number_masked}
                      </td>
                      <td style={{ padding: '12px 16px', color: 'var(--text-secondary)' }}>
                        {bank.account_type}
                      </td>
                      <td style={{ padding: '12px 16px', color: 'var(--text-secondary)' }}>
                        <span style={{ padding: '2px 8px', borderRadius: 12, backgroundColor: 'var(--bg-tertiary)', fontSize: 12, fontWeight: 600 }}>
                          {bank.currency}
                        </span>
                      </td>
                      <td style={{ padding: '12px 16px', color: 'var(--text-primary)', textAlign: 'right', fontFamily: 'JetBrains Mono, monospace', fontWeight: 600 }}>
                        {bank.available_balance.toLocaleString()}
                      </td>
                      <td style={{ padding: '12px 16px' }}>
                        <span style={{ 
                          padding: '4px 8px', 
                          borderRadius: 4, 
                          fontSize: 12, 
                          fontWeight: 500,
                          backgroundColor: bank.status === 'ACTIVE' ? 'rgba(16, 185, 129, 0.1)' : 'rgba(239, 68, 68, 0.1)',
                          color: bank.status === 'ACTIVE' ? '#10b981' : '#ef4444'
                        }}>
                          {bank.status}
                        </span>
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan={6} style={{ padding: '24px 16px', textAlign: 'center', color: 'var(--text-muted)' }}>
                      No bank accounts found.
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
