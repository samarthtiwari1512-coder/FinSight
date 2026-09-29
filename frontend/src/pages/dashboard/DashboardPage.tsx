import { useQuery } from '@tanstack/react-query'
import {
  AreaChart, Area, BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer,
  PieChart, Pie, Cell
} from 'recharts'
import {
  Wallet, TrendingUp, TrendingDown, AlertTriangle, Globe,
  ArrowUpRight, ArrowDownRight, RefreshCw, Clock, Shield,
  CircleDollarSign, Building2, Zap, Activity
} from 'lucide-react'
import { dashboardApi, forecastApi } from '../../api'
import { format } from 'date-fns'
import type { DashboardSummary, ForecastOutput } from '../../api'

import type { LucideIcon } from 'lucide-react'

// ─── Utility ────────────────────────────────────────────────────────────────

const fmtCr = (n: number) => {
  if (Math.abs(n) >= 10000000) return `₹${(n / 10000000).toFixed(2)} Cr`
  if (Math.abs(n) >= 100000) return `₹${(n / 100000).toFixed(1)} L`
  return `₹${n.toLocaleString('en-IN')}`
}

const fmtPct = (n: number) => `${n >= 0 ? '+' : ''}${n.toFixed(1)}%`

// ─── KPI Card ────────────────────────────────────────────────────────────────

interface KPIProps {
  label: string
  value: string
  change?: number
  changeLabel?: string
  icon: LucideIcon
  color?: string
  risk?: string
  sublabel?: string
}

function KPICard({ label, value, change, changeLabel, icon: Icon, color = '#3b82f6', risk, sublabel }: KPIProps) {
  const isUp = change !== undefined && change >= 0
  const riskColors: Record<string, string> = {
    LOW: '#10b981', MEDIUM: '#f59e0b', HIGH: '#f97316', CRITICAL: '#ef4444',
  }

  return (
    <div className="kpi-card fade-in">
      <div style={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', marginBottom: 14 }}>
        <div style={{
          width: 38,
          height: 38,
          background: `${color}18`,
          border: `1px solid ${color}30`,
          borderRadius: 10,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}>
          <Icon size={17} color={color} />
        </div>
        {risk && (
          <span style={{
            fontSize: 10,
            fontWeight: 700,
            padding: '2px 8px',
            borderRadius: 100,
            background: `${riskColors[risk] || '#64748b'}18`,
            color: riskColors[risk] || '#64748b',
            border: `1px solid ${riskColors[risk] || '#64748b'}30`,
          }}>
            {risk}
          </span>
        )}
      </div>

      <div style={{ fontSize: 11, color: 'var(--text-muted)', fontWeight: 600, letterSpacing: '0.05em', textTransform: 'uppercase', marginBottom: 4 }}>
        {label}
      </div>

      <div style={{ fontSize: 22, fontWeight: 800, color: 'var(--text-primary)', letterSpacing: '-0.01em', marginBottom: 6 }}>
        {value}
      </div>

      {sublabel && (
        <div style={{ fontSize: 11, color: 'var(--text-muted)', marginBottom: 4 }}>{sublabel}</div>
      )}

      {change !== undefined && (
        <div style={{ display: 'flex', alignItems: 'center', gap: 4 }}>
          {isUp ? <ArrowUpRight size={13} style={{ color: '#10b981' }} /> : <ArrowDownRight size={13} style={{ color: '#ef4444' }} />}
          <span style={{ fontSize: 12, fontWeight: 600, color: isUp ? '#10b981' : '#ef4444' }}>
            {fmtPct(Math.abs(change))}
          </span>
          {changeLabel && (
            <span style={{ fontSize: 11, color: 'var(--text-muted)' }}>{changeLabel}</span>
          )}
        </div>
      )}
    </div>
  )
}

// ─── Aging Bar ───────────────────────────────────────────────────────────────

function AgingBar({ buckets }: { buckets: Array<{ label: string; amount: number; percent: number }> }) {
  const COLORS = ['#10b981', '#60a5fa', '#f59e0b', '#f97316', '#ef4444']
  return (
    <div>
      <div style={{ display: 'flex', height: 8, borderRadius: 4, overflow: 'hidden', gap: 2, marginBottom: 12 }}>
        {buckets.map((b, i) => (
          <div key={b.label} style={{
            flex: b.percent,
            background: COLORS[i % COLORS.length],
            opacity: b.percent > 0 ? 1 : 0,
          }} />
        ))}
      </div>
      <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap' }}>
        {buckets.map((b, i) => (
          <div key={b.label} style={{ display: 'flex', alignItems: 'center', gap: 5 }}>
            <div style={{ width: 8, height: 8, borderRadius: 2, background: COLORS[i % COLORS.length], flexShrink: 0 }} />
            <span style={{ fontSize: 11, color: 'var(--text-secondary)' }}>{b.label}</span>
            <span style={{ fontSize: 11, color: 'var(--text-muted)' }}>({b.percent.toFixed(0)}%)</span>
          </div>
        ))}
      </div>
    </div>
  )
}

// ─── Alert Row ───────────────────────────────────────────────────────────────

function AlertRow({ alert }: { alert: any }) {
  const colors: Record<string, string> = {
    CRITICAL: '#ef4444', HIGH: '#f97316', MEDIUM: '#f59e0b', LOW: '#10b981',
  }
  const color = colors[alert.severity] || '#64748b'

  return (
    <div style={{
      display: 'flex',
      alignItems: 'flex-start',
      gap: 10,
      padding: '10px 0',
      borderBottom: '1px solid var(--border-subtle)',
    }}>
      <div style={{ width: 6, height: 6, borderRadius: '50%', background: color, marginTop: 5, flexShrink: 0 }} />
      <div style={{ flex: 1, minWidth: 0 }}>
        <div style={{ fontSize: 12, fontWeight: 600, color: 'var(--text-primary)', marginBottom: 2 }}>
          {alert.title}
        </div>
        <div style={{ fontSize: 11, color: 'var(--text-muted)', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
          {alert.description}
        </div>
      </div>
      <span style={{
        fontSize: 10, fontWeight: 700, padding: '2px 6px', borderRadius: 100,
        background: `${color}18`, color, border: `1px solid ${color}25`, whiteSpace: 'nowrap',
      }}>
        {alert.severity}
      </span>
    </div>
  )
}

// ─── FX Rate Row ─────────────────────────────────────────────────────────────

function FXRow({ curr, rate, stale }: { curr: string; rate: number; stale: boolean }) {
  return (
    <div style={{
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'space-between',
      padding: '8px 0',
      borderBottom: '1px solid var(--border-subtle)',
    }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
        <div style={{
          width: 28, height: 20, background: 'rgba(37,99,235,0.1)',
          borderRadius: 4, display: 'flex', alignItems: 'center',
          justifyContent: 'center', fontSize: 9, fontWeight: 700, color: '#60a5fa',
        }}>{curr}</div>
        <span style={{ fontSize: 12, color: 'var(--text-secondary)' }}>{curr}/INR</span>
      </div>
      <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
        <span style={{ fontFamily: 'JetBrains Mono, monospace', fontSize: 12, fontWeight: 600, color: 'var(--text-primary)' }}>
          ₹{rate.toFixed(4)}
        </span>
        {stale && <Clock size={11} style={{ color: '#f59e0b' }} />}
      </div>
    </div>
  )
}

// ─── Dashboard Page ───────────────────────────────────────────────────────────

export default function DashboardPage() {
  const { data: summaryResp, isLoading: loadingSummary, dataUpdatedAt } = useQuery({
    queryKey: ['dashboard-summary'],
    queryFn: () => dashboardApi.getSummary().then(r => r.data),
    refetchInterval: 5 * 60 * 1000, // Refresh every 5 min
  })

  const { data: forecast30Resp } = useQuery({
    queryKey: ['forecast-30d'],
    queryFn: () => forecastApi.get30d().then(r => r.data),
  })

  const summary: DashboardSummary | undefined = summaryResp?.data
  const forecast: ForecastOutput | undefined = forecast30Resp?.data

  if (loadingSummary) {
    return (
      <div className="page">
        <div style={{ display: 'flex', gap: 16, marginBottom: 24 }}>
          {[...Array(5)].map((_, i) => (
            <div key={i} className="skeleton" style={{ height: 110, flex: 1 }} />
          ))}
        </div>
        <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: 20 }}>
          <div className="skeleton" style={{ height: 280 }} />
          <div className="skeleton" style={{ height: 280 }} />
        </div>
      </div>
    )
  }

  const cash = summary?.cash_position
  const wc = summary?.working_capital
  const ar = summary?.ar_summary
  const ap = summary?.ap_summary
  const fx = summary?.fx_summary

  // Build forecast chart data
  const forecastChartData = forecast?.daily_items?.slice(0, 30).map(d => ({
    date: format(new Date(d.date), 'MMM d'),
    inflow: Math.round(d.projected_inflow / 100000),
    outflow: Math.round(d.projected_outflow / 100000),
    balance: Math.round(d.closing_balance / 10000000),
  })) || []

  // Aging donut data
  const arAgingData = (ar?.aging_buckets || []).map((b, i) => ({
    name: b.label,
    value: b.amount,
  }))

  const AGING_COLORS = ['#10b981', '#60a5fa', '#f59e0b', '#f97316', '#ef4444']

  return (
    <div className="page fade-in">
      {/* Page Header */}
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 24 }}>
        <div>
          <h1 style={{ fontSize: 20, fontWeight: 700, color: 'var(--text-primary)', marginBottom: 2 }}>
            Executive Dashboard
          </h1>
          <div style={{ fontSize: 12, color: 'var(--text-muted)', display: 'flex', alignItems: 'center', gap: 8 }}>
            <Activity size={12} />
            <span>
              As of {dataUpdatedAt ? format(new Date(dataUpdatedAt), 'dd MMM yyyy, HH:mm') : 'Loading...'}
            </span>
            <span style={{ padding: '1px 8px', borderRadius: 100, background: 'rgba(16,185,129,0.1)', color: '#10b981', fontSize: 10, fontWeight: 600 }}>
              LIVE
            </span>
          </div>
        </div>
        <div style={{ display: 'flex', gap: 8 }}>
          <select className="select" style={{ width: 140 }}>
            <option>Last 30 days</option>
            <option>Last 90 days</option>
            <option>FY 2024-25</option>
          </select>
          <button className="btn btn-secondary btn-sm">
            <RefreshCw size={13} /> Refresh
          </button>
        </div>
      </div>

      {/* Risk Banner */}
      {summary?.liquidity_risk === 'HIGH' || summary?.fx_risk === 'HIGH' ? (
        <div style={{
          background: 'rgba(239,68,68,0.05)',
          border: '1px solid rgba(239,68,68,0.2)',
          borderRadius: 10,
          padding: '10px 16px',
          marginBottom: 20,
          display: 'flex',
          alignItems: 'center',
          gap: 10,
        }}>
          <AlertTriangle size={15} style={{ color: '#ef4444', flexShrink: 0 }} />
          <span style={{ fontSize: 13, color: '#ef4444', fontWeight: 500 }}>
            Elevated risk detected: {summary?.liquidity_risk === 'HIGH' ? 'Liquidity risk is HIGH. ' : ''}
            {summary?.fx_risk === 'HIGH' ? 'FX exposure risk is HIGH.' : ''}
            Review alerts for recommended actions.
          </span>
        </div>
      ) : null}

      {/* KPI Row */}
      <div className="kpi-grid" style={{ marginBottom: 20 }}>
        <KPICard
          label="Total Cash Position"
          value={fmtCr(cash?.closing_balance || 0)}
          change={2.3}
          changeLabel="vs last week"
          icon={Wallet}
          color="#3b82f6"
          sublabel={`${cash?.by_bank?.length || 0} bank accounts`}
        />
        <KPICard
          label="Net Receivables (AR)"
          value={fmtCr(ar?.total_outstanding || 0)}
          change={-3.1}
          changeLabel="vs last month"
          icon={TrendingUp}
          color="#10b981"
          sublabel={`${ar?.overdue_count || 0} invoices overdue`}
          risk={ar && ar.overdue_count > 50 ? 'HIGH' : 'MEDIUM'}
        />
        <KPICard
          label="Total Payables (AP)"
          value={fmtCr(ap?.total_outstanding || 0)}
          change={1.8}
          changeLabel="vs last month"
          icon={TrendingDown}
          color="#f59e0b"
          sublabel={`₹${((ap?.due_in_3_days || 0) / 10000000).toFixed(1)} Cr due in 3 days`}
          risk="MEDIUM"
        />
        <KPICard
          label="FX Net Exposure"
          value={`$${((fx?.total_net_exposure || 0) / 1000000).toFixed(1)}M`}
          change={-1.5}
          changeLabel="vs yesterday"
          icon={Globe}
          color="#a78bfa"
          sublabel={`${((fx?.coverage_ratio || 0) * 100).toFixed(0)}% hedged`}
          risk={summary?.fx_risk as string}
        />
        <KPICard
          label="Cash Conversion Cycle"
          value={`${(wc?.ccc || 0).toFixed(0)} days`}
          change={wc ? wc.ccc - 48 : 0}
          changeLabel="vs 48d target"
          icon={RefreshCw}
          color="#f97316"
          sublabel={`DSO ${(wc?.dso || 0).toFixed(0)}d · DPO ${(wc?.dpo || 0).toFixed(0)}d`}
        />
      </div>

      {/* Charts Row 1 */}
      <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: 20, marginBottom: 20 }}>
        {/* Cash Flow Forecast */}
        <div className="card">
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 16 }}>
            <div>
              <div style={{ fontSize: 14, fontWeight: 700, color: 'var(--text-primary)', marginBottom: 2 }}>
                30-Day Cash Flow Forecast
              </div>
              <div style={{ fontSize: 11, color: '#f59e0b', display: 'flex', alignItems: 'center', gap: 4 }}>
                <Zap size={11} />
                Statistical estimate — confidence: {((forecast?.confidence || 0) * 100).toFixed(0)}%
              </div>
            </div>
            {forecast?.liquidity_warning && (
              <div style={{
                display: 'flex',
                alignItems: 'center',
                gap: 6,
                padding: '4px 10px',
                borderRadius: 8,
                background: 'rgba(239,68,68,0.1)',
                border: '1px solid rgba(239,68,68,0.2)',
              }}>
                <AlertTriangle size={12} style={{ color: '#ef4444' }} />
                <span style={{ fontSize: 11, fontWeight: 600, color: '#ef4444' }}>Liquidity Warning</span>
              </div>
            )}
          </div>
          <ResponsiveContainer width="100%" height={200}>
            <AreaChart data={forecastChartData}>
              <defs>
                <linearGradient id="inflowGrad" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#10b981" stopOpacity={0.25} />
                  <stop offset="95%" stopColor="#10b981" stopOpacity={0} />
                </linearGradient>
                <linearGradient id="outflowGrad" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#ef4444" stopOpacity={0.2} />
                  <stop offset="95%" stopColor="#ef4444" stopOpacity={0} />
                </linearGradient>
              </defs>
              <XAxis
                dataKey="date"
                tick={{ fontSize: 10, fill: 'var(--text-muted)' }}
                tickLine={false}
                axisLine={false}
                interval={4}
              />
              <YAxis
                tick={{ fontSize: 10, fill: 'var(--text-muted)' }}
                tickLine={false}
                axisLine={false}
                tickFormatter={(v) => `₹${v}L`}
                width={50}
              />
              <Tooltip
                contentStyle={{
                  background: 'var(--bg-elevated)',
                  border: '1px solid var(--border-color)',
                  borderRadius: 8,
                  fontSize: 12,
                }}
                formatter={(value: number, name: string) => [
                  `₹${value}L`, name === 'inflow' ? 'Inflow' : 'Outflow'
                ]}
              />
              <Area type="monotone" dataKey="inflow" stroke="#10b981" strokeWidth={1.5} fill="url(#inflowGrad)" />
              <Area type="monotone" dataKey="outflow" stroke="#ef4444" strokeWidth={1.5} fill="url(#outflowGrad)" />
            </AreaChart>
          </ResponsiveContainer>
          {forecast?.disclaimer && (
            <div style={{ fontSize: 10, color: 'var(--text-muted)', marginTop: 10, borderTop: '1px solid var(--border-subtle)', paddingTop: 10 }}>
              ⚠ {forecast.disclaimer}
            </div>
          )}
        </div>

        {/* AR Aging Donut */}
        <div className="card">
          <div style={{ fontSize: 14, fontWeight: 700, color: 'var(--text-primary)', marginBottom: 4 }}>AR Aging Breakdown</div>
          <div style={{ fontSize: 11, color: 'var(--text-muted)', marginBottom: 16 }}>
            Total: {fmtCr(ar?.total_outstanding || 0)}
          </div>
          {arAgingData.length > 0 ? (
            <>
              <ResponsiveContainer width="100%" height={140}>
                <PieChart>
                  <Pie
                    data={arAgingData}
                    cx="50%"
                    cy="50%"
                    innerRadius={40}
                    outerRadius={65}
                    paddingAngle={2}
                    dataKey="value"
                  >
                    {arAgingData.map((_, i) => (
                      <Cell key={i} fill={AGING_COLORS[i % AGING_COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip
                    contentStyle={{
                      background: 'var(--bg-elevated)',
                      border: '1px solid var(--border-color)',
                      borderRadius: 8,
                      fontSize: 11,
                    }}
                    formatter={(v: number) => fmtCr(v)}
                  />
                </PieChart>
              </ResponsiveContainer>
              <AgingBar buckets={ar?.aging_buckets || []} />
            </>
          ) : (
            <div className="empty-state" style={{ padding: 32 }}>No AR data</div>
          )}
        </div>
      </div>

      {/* Charts Row 2 */}
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: 20 }}>
        {/* Working Capital KPIs */}
        <div className="card">
          <div style={{ fontSize: 14, fontWeight: 700, color: 'var(--text-primary)', marginBottom: 16 }}>
            Working Capital Metrics
          </div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
            {[
              { label: 'Days Sales Outstanding (DSO)', value: `${(wc?.dso || 0).toFixed(1)}d`, target: 30, max: 90, color: '#10b981' },
              { label: 'Days Payable Outstanding (DPO)', value: `${(wc?.dpo || 0).toFixed(1)}d`, target: 45, max: 90, color: '#3b82f6' },
              { label: 'Cash Conversion Cycle', value: `${(wc?.ccc || 0).toFixed(1)}d`, target: 30, max: 120, color: '#f59e0b' },
              { label: 'Current Ratio', value: `${(wc?.current_ratio || 0).toFixed(2)}x`, target: null, max: null, color: '#a78bfa' },
            ].map((m) => (
              <div key={m.label}>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 4 }}>
                  <span style={{ fontSize: 11, color: 'var(--text-muted)' }}>{m.label}</span>
                  <span style={{ fontSize: 12, fontWeight: 700, color: 'var(--text-primary)', fontFamily: 'JetBrains Mono, monospace' }}>
                    {m.value}
                  </span>
                </div>
                {m.max && (
                  <div className="progress-bar">
                    <div
                      className="progress-fill"
                      style={{
                        width: `${Math.min(100, ((parseFloat(m.value) || 0) / m.max) * 100)}%`,
                        background: m.color,
                      }}
                    />
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>

        {/* FX Exposure by Currency */}
        <div className="card">
          <div style={{ fontSize: 14, fontWeight: 700, color: 'var(--text-primary)', marginBottom: 4 }}>FX Exposure</div>
          <div style={{ fontSize: 11, color: 'var(--text-muted)', marginBottom: 16 }}>
            Net exposure by currency
          </div>
          {(fx?.by_currency || []).map((exp) => {
            const riskColors: Record<string, string> = { LOW: '#10b981', MEDIUM: '#f59e0b', HIGH: '#f97316', CRITICAL: '#ef4444' }
            const color = riskColors[exp.risk_level] || '#64748b'
            return (
              <div key={exp.currency} style={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                padding: '8px 0',
                borderBottom: '1px solid var(--border-subtle)',
              }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                  <div style={{
                    width: 28, height: 20,
                    background: `${color}18`,
                    border: `1px solid ${color}30`,
                    borderRadius: 4,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    fontSize: 9, fontWeight: 700, color,
                  }}>{exp.currency}</div>
                  <div>
                    <div style={{ fontSize: 12, fontWeight: 500, color: 'var(--text-primary)' }}>{exp.currency}/INR</div>
                    <div style={{ fontSize: 10, color: 'var(--text-muted)' }}>{exp.direction || 'NET'}</div>
                  </div>
                </div>
                <div style={{ textAlign: 'right' }}>
                  <div style={{ fontSize: 12, fontWeight: 700, color: 'var(--text-primary)', fontFamily: 'JetBrains Mono, monospace' }}>
                    {fmtCr(exp.base_value || 0)}
                  </div>
                  <span style={{ fontSize: 9, fontWeight: 700, padding: '1px 5px', borderRadius: 100, background: `${color}18`, color, border: `1px solid ${color}25` }}>
                    {exp.risk_level}
                  </span>
                </div>
              </div>
            )
          })}
          {fx?.by_currency?.length === 0 && (
            <div className="empty-state" style={{ padding: 24 }}>No FX exposure data</div>
          )}
        </div>

        {/* Active Alerts */}
        <div className="card">
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 12 }}>
            <div style={{ fontSize: 14, fontWeight: 700, color: 'var(--text-primary)' }}>Active Alerts</div>
            <span style={{
              fontSize: 11, color: '#ef4444', fontWeight: 600,
              padding: '2px 8px', borderRadius: 100,
              background: 'rgba(239,68,68,0.1)', border: '1px solid rgba(239,68,68,0.2)',
            }}>
              {summary?.active_alerts?.length || 0} active
            </span>
          </div>
          <div style={{ overflowY: 'auto', maxHeight: 260 }}>
            {(summary?.active_alerts || []).length > 0 ? (
              summary!.active_alerts.map((alert) => (
                <AlertRow key={alert.id} alert={alert} />
              ))
            ) : (
              <div className="empty-state">
                <Shield size={32} />
                <div style={{ marginTop: 8 }}>No active alerts</div>
              </div>
            )}
          </div>
          <button
            className="btn btn-secondary btn-sm"
            style={{ width: '100%', marginTop: 12, justifyContent: 'center' }}
            onClick={() => window.location.href = '/alerts'}
          >
            View All Alerts
          </button>
        </div>
      </div>
    </div>
  )
}
