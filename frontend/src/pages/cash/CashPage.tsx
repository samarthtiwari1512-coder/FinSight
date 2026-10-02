import { useQuery } from '@tanstack/react-query'
import { cashApi } from '../../api'
import { Landmark, ArrowUpRight, ArrowDownRight, Wallet, RefreshCw } from 'lucide-react'

export default function CashPage() {
  const { data: positionData, isLoading: isLoadingPos } = useQuery({
    queryKey: ['cash', 'position'],
    queryFn: () => cashApi.getPosition().then(res => res.data.data)
  })

  const { data: bankData, isLoading: isLoadingBanks } = useQuery({
    queryKey: ['cash', 'by-bank'],
    queryFn: () => cashApi.getPositionByBank().then(res => res.data.data)
  })

  return (
    <div className="page fade-in">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <div>
          <h1 style={{ fontSize: 24, fontWeight: 700, color: 'var(--fs-text)', marginBottom: 4 }}>
            Cash Management
          </h1>
          <p style={{ color: 'var(--fs-text-muted)', fontSize: 14 }}>
            Monitor global liquidity, bank accounts, and recent cash movements.
          </p>
        </div>
        <button className="btn-secondary" style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          <RefreshCw size={14} />
          <span>Refresh</span>
        </button>
      </div>

      {/* KPI Cards */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))', gap: 16, marginBottom: 24 }}>
        <div className="glass-panel" style={{ padding: 20 }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 16 }}>
            <div style={{ width: 40, height: 40, borderRadius: 8, backgroundColor: 'rgba(59, 130, 246, 0.1)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#3b82f6' }}>
              <Wallet size={20} />
            </div>
            <h3 style={{ fontSize: 14, fontWeight: 600, color: 'var(--fs-text-secondary)' }}>Total Cash (INR)</h3>
          </div>
          <div style={{ fontSize: 28, fontWeight: 700, color: 'var(--fs-text)', fontFamily: 'JetBrains Mono, monospace' }}>
            {isLoadingPos ? '...' : `₹${positionData?.closing_balance?.toLocaleString() || '0'}`}
          </div>
        </div>

        <div className="glass-panel" style={{ padding: 20 }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 16 }}>
            <div style={{ width: 40, height: 40, borderRadius: 8, backgroundColor: 'rgba(16, 185, 129, 0.1)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#10b981' }}>
              <ArrowUpRight size={20} />
            </div>
            <h3 style={{ fontSize: 14, fontWeight: 600, color: 'var(--fs-text-secondary)' }}>Total Inflows</h3>
          </div>
          <div style={{ fontSize: 28, fontWeight: 700, color: 'var(--fs-text)', fontFamily: 'JetBrains Mono, monospace' }}>
            {isLoadingPos ? '...' : `₹${positionData?.total_inflows?.toLocaleString() || '0'}`}
          </div>
        </div>

        <div className="glass-panel" style={{ padding: 20 }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 16 }}>
            <div style={{ width: 40, height: 40, borderRadius: 8, backgroundColor: 'rgba(239, 68, 68, 0.1)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#ef4444' }}>
              <ArrowDownRight size={20} />
            </div>
            <h3 style={{ fontSize: 14, fontWeight: 600, color: 'var(--fs-text-secondary)' }}>Total Outflows</h3>
          </div>
          <div style={{ fontSize: 28, fontWeight: 700, color: 'var(--fs-text)', fontFamily: 'JetBrains Mono, monospace' }}>
            {isLoadingPos ? '...' : `₹${positionData?.total_outflows?.toLocaleString() || '0'}`}
          </div>
        </div>
      </div>

      {/* Bank Balances Table */}
      <div className="glass-panel" style={{ padding: 20 }}>
        <h2 style={{ fontSize: 16, fontWeight: 600, color: 'var(--fs-text)', marginBottom: 16, display: 'flex', alignItems: 'center', gap: 8 }}>
          <Landmark size={18} />
          Balances by Bank
        </h2>
        
        {isLoadingBanks ? (
          <div style={{ padding: 40, textAlign: 'center', color: 'var(--fs-text-muted)' }}>Loading bank balances...</div>
        ) : (
          <div style={{ overflowX: 'auto' }}>
            <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: 14 }}>
              <thead>
                <tr style={{ borderBottom: '1px solid var(--fs-border)', color: 'var(--fs-text-secondary)', textAlign: 'left' }}>
                  <th style={{ padding: '12px 16px', fontWeight: 600 }}>Bank Name</th>
                  <th style={{ padding: '12px 16px', fontWeight: 600 }}>Account No</th>
                  <th style={{ padding: '12px 16px', fontWeight: 600 }}>Currency</th>
                  <th style={{ padding: '12px 16px', fontWeight: 600, textAlign: 'right' }}>Balance</th>
                  <th style={{ padding: '12px 16px', fontWeight: 600, textAlign: 'right' }}>Equivalent (INR)</th>
                </tr>
              </thead>
              <tbody>
                {bankData && bankData.length > 0 ? (
                  bankData.map(bank => (
                    <tr key={bank.bank_name + bank.currency} style={{ borderBottom: '1px solid var(--fs-border)' }}>
                      <td style={{ padding: '12px 16px', color: 'var(--fs-text)', fontWeight: 500 }}>{bank.bank_name}</td>
                      <td style={{ padding: '12px 16px', color: 'var(--fs-text-secondary)', fontFamily: 'JetBrains Mono, monospace' }}>{bank.account_number_masked}</td>
                      <td style={{ padding: '12px 16px', color: 'var(--fs-text-secondary)' }}>
                        <span style={{ padding: '2px 8px', borderRadius: 12, backgroundColor: 'var(--fs-surface-hover)', fontSize: 12, fontWeight: 600 }}>
                          {bank.currency}
                        </span>
                      </td>
                      <td style={{ padding: '12px 16px', color: 'var(--fs-text)', textAlign: 'right', fontFamily: 'JetBrains Mono, monospace' }}>
                        {bank.available_balance.toLocaleString()}
                      </td>
                      <td style={{ padding: '12px 16px', color: 'var(--fs-text)', textAlign: 'right', fontFamily: 'JetBrains Mono, monospace', fontWeight: 600 }}>
                        ₹{bank.base_balance.toLocaleString()}
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan={5} style={{ padding: '24px 16px', textAlign: 'center', color: 'var(--fs-text-muted)' }}>
                      No bank balances available.
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
