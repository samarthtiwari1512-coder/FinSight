import { Link } from 'react-router-dom'
import { ArrowRight } from 'lucide-react'

const SOLUTIONS = [
  {
    slug: 'working-capital',
    label: 'Working Capital',
    tagline: 'Know exactly how fast your money moves.',
    description: 'Working capital is the difference between current assets and current liabilities. It determines whether you can pay your suppliers without borrowing, whether you can grow without straining cash, and whether your business is healthy day-to-day.',
    why: 'Most companies track working capital quarterly at best — which means problems are discovered after they become expensive. The components that drive working capital (receivables aging, payables timing, inventory levels) change daily.',
    approach: 'FinSight computes DSO, DPO, DIO, and CCC from actual transaction data, updated in real time. You see where the cycle is lengthening, which customers are paying slowly, and whether your payables strategy is optimal.',
    metrics: [
      { label: 'Days Sales Outstanding', value: '42d', vs: '+3d vs target', isNeg: true },
      { label: 'Days Payable Outstanding', value: '38d', vs: 'On target', isNeg: false },
      { label: 'Days Inventory Outstanding', value: '22d', vs: '−1d this month', isNeg: false },
      { label: 'Cash Conversion Cycle', value: '26d', vs: 'Improving', isNeg: false },
      { label: 'Current Ratio', value: '1.42×', vs: 'Healthy', isNeg: false },
    ],
    deeper: [
      { heading: 'Receivables aging', body: 'Track which invoices are current, 30-60 days, 60-90 days, or over 90 days. Identify customers with deteriorating payment patterns early.' },
      { heading: 'Payables optimization', body: 'Balance payment timing to maximize DPO without straining supplier relationships or missing early-payment discounts.' },
      { heading: 'CCC trend analysis', body: 'Watch the Cash Conversion Cycle over time. A rising CCC is often the earliest signal of business stress.' },
    ],
  },
  {
    slug: 'cash-flow',
    label: 'Cash Flow',
    tagline: 'Know when cash gets tight before it happens.',
    description: 'Cash flow is the movement of money into and out of your business. It is distinct from profit — a profitable company can run out of cash. Understanding your cash position today, and projecting it forward, is one of the most critical financial management tasks.',
    why: 'Most treasury teams spend too much time manually assembling cash position reports from multiple bank feeds, spreadsheets and ERP exports. By the time the report is ready, the numbers have changed.',
    approach: 'FinSight consolidates all bank account positions, maps incoming receivables and outgoing payables, and generates a 30-day forward projection using statistical modeling. Liquidity warnings trigger automatically when the projected position falls below threshold.',
    metrics: [
      { label: 'Total Cash (INR)', value: '₹82.4M', vs: '+4.8% vs last week', isNeg: false },
      { label: '7-Day Projected Inflow', value: '₹12.1M', vs: 'On forecast', isNeg: false },
      { label: '7-Day Projected Outflow', value: '₹9.3M', vs: '−₹0.7M vs forecast', isNeg: false },
      { label: 'Bank Accounts', value: '24', vs: 'Across 6 entities', isNeg: false },
      { label: 'Liquidity Status', value: 'Adequate', vs: 'All thresholds met', isNeg: false },
    ],
    deeper: [
      { heading: '30-day forecast', body: 'Statistical projection of daily inflows and outflows based on historical patterns and confirmed transactions.' },
      { heading: 'Liquidity alerts', body: 'Automated alerts when projected balance falls below minimum threshold. Configurable per entity and currency.' },
      { heading: 'Multi-entity consolidation', body: 'Aggregate cash positions across all legal entities, subsidiaries and bank accounts into a single consolidated view.' },
    ],
  },
  {
    slug: 'fx',
    label: 'FX Exposure',
    tagline: 'Understand your currency risk before it understands you.',
    description: 'If your company transacts in multiple currencies, every change in exchange rates affects your effective cash position, receivables value and payables cost. FX exposure management is about knowing your net position and deciding how much of it to hedge.',
    why: 'Companies often discover their FX exposure only when they reconcile quarterly financials — by which point currency movements have already had their impact. Real-time exposure tracking allows proactive hedging decisions.',
    approach: 'FinSight aggregates FX exposure by currency from all receivables, payables and cash balances. It shows net exposure, current hedging coverage, sensitivity to rate movements, and tracks active FX forward deals.',
    metrics: [
      { label: 'Total Net FX Exposure', value: '$45.2M', vs: 'Across 8 currencies', isNeg: false },
      { label: 'USD Exposure', value: '$28.4M', vs: 'Partially hedged', isNeg: false },
      { label: 'EUR Exposure', value: '€9.1M', vs: 'Fully hedged', isNeg: false },
      { label: 'GBP Exposure', value: '£4.8M', vs: 'Unhedged — HIGH', isNeg: true },
      { label: 'Hedging Coverage', value: '64%', vs: 'Target: 80%', isNeg: true },
    ],
    deeper: [
      { heading: 'Currency-level breakdown', body: 'See exposure broken down by currency, with risk ratings based on volatility and your hedging coverage.' },
      { heading: 'Sensitivity analysis', body: 'Model the impact of rate moves on your effective position. A 5% depreciation in GBP affects ₹X crore.' },
      { heading: 'Forward deal tracking', body: 'Record and track FX forward contracts alongside your live exposure to see true net hedged positions.' },
    ],
  },
  {
    slug: 'risk',
    label: 'Risk Management',
    tagline: 'Financial risks identified before they become problems.',
    description: 'Financial risk in treasury spans liquidity risk, counterparty risk, FX risk, interest rate risk and operational risk. Managing these is not about predicting the future — it is about knowing your current exposures and having alerts in place when thresholds are breached.',
    why: 'Risk tends to accumulate gradually and announce itself suddenly. A customer approaching their credit limit, a currency position growing beyond tolerance, a payable cluster creating a liquidity crunch — all of these can be detected early with the right monitoring.',
    approach: 'FinSight generates structured risk events from financial data, evaluates them against configurable thresholds, and delivers prioritized alerts. Every alert includes severity, affected entity, recommended action and audit trail.',
    metrics: [
      { label: 'Active Alerts', value: '3', vs: '1 CRITICAL, 2 MEDIUM', isNeg: true },
      { label: 'Liquidity Risk', value: 'MEDIUM', vs: 'Monitor closely', isNeg: true },
      { label: 'FX Risk', value: 'HIGH', vs: 'GBP exposure unhedged', isNeg: true },
      { label: 'Counterparty Risk', value: 'LOW', vs: 'All within limits', isNeg: false },
      { label: 'Risk Events (30d)', value: '12', vs: '2 escalated', isNeg: false },
    ],
    deeper: [
      { heading: 'Threshold-based alerting', body: 'Define minimum cash thresholds, maximum FX exposure limits, and overdue receivables triggers. Alerts fire automatically.' },
      { heading: 'Severity classification', body: 'Risk events are classified as LOW, MEDIUM, HIGH or CRITICAL. Each level triggers different notification and escalation behaviors.' },
      { heading: 'Scenario analysis', body: 'Model what happens to your liquidity and working capital under different scenarios — a major customer defaulting, a currency devaluation, a supplier demanding early payment.' },
    ],
  },
]

export default function SolutionsIndexPage() {
  return (
    <div>
      <section style={{ padding: '100px 0 80px', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 16 }}>
            Solutions
          </p>
          <h1 className="font-display" style={{ fontSize: 'clamp(36px, 5vw, 64px)', color: 'var(--fs-text)', lineHeight: 1.05, marginBottom: 24 }}>
            Financial clarity for every team.
          </h1>
          <p style={{ fontSize: 18, color: 'var(--fs-text-secondary)', maxWidth: '52ch', lineHeight: 1.7 }}>
            FinSight addresses the core financial management challenges facing corporate treasury and finance teams — not as separate tools, but as an integrated system.
          </p>
        </div>
      </section>

      <section style={{ padding: '80px 0' }}>
        <div className="pub-container">
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))', gap: 1, border: '1px solid var(--fs-border)', background: 'var(--fs-border)' }}>
            {SOLUTIONS.map(s => (
              <Link key={s.slug} to={`/solutions/${s.slug}`} style={{ textDecoration: 'none' }}>
                <div style={{
                  padding: '40px 36px',
                  background: 'var(--fs-surface)',
                  height: '100%',
                  transition: 'background 0.15s',
                }}
                  onMouseEnter={e => (e.currentTarget as HTMLElement).style.background = 'var(--fs-bg)'}
                  onMouseLeave={e => (e.currentTarget as HTMLElement).style.background = 'var(--fs-surface)'}
                >
                  <h2 style={{ fontSize: 22, fontWeight: 600, color: 'var(--fs-text)', marginBottom: 12 }}>{s.label}</h2>
                  <p style={{ fontSize: 15, color: 'var(--fs-text-secondary)', lineHeight: 1.7, marginBottom: 24 }}>{s.tagline}</p>
                  <span style={{ fontSize: 13, color: 'var(--fs-accent)', fontWeight: 500 }}>Explore →</span>
                </div>
              </Link>
            ))}
          </div>
        </div>
      </section>
    </div>
  )
}

// Individual Solution Page Component
export function SolutionDetailPage({ slug }: { slug: string }) {
  const solution = SOLUTIONS.find(s => s.slug === slug)
  if (!solution) return <NotFoundInline />

  return (
    <div>
      {/* Hero */}
      <section style={{ padding: '100px 0 80px', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 16 }}>
            <Link to="/solutions" style={{ fontSize: 13, color: 'var(--fs-text-muted)', textDecoration: 'none' }}>Solutions</Link>
            <span style={{ color: 'var(--fs-text-muted)' }}>/</span>
            <span style={{ fontSize: 13, color: 'var(--fs-text-secondary)' }}>{solution.label}</span>
          </div>
          <h1 className="font-display" style={{ fontSize: 'clamp(36px, 5vw, 64px)', color: 'var(--fs-text)', lineHeight: 1.05, marginBottom: 24 }}>
            {solution.label}
          </h1>
          <p style={{ fontSize: 20, color: 'var(--fs-accent)', fontWeight: 500, marginBottom: 24 }}>
            {solution.tagline}
          </p>
          <p style={{ fontSize: 17, color: 'var(--fs-text-secondary)', maxWidth: '58ch', lineHeight: 1.8 }}>
            {solution.description}
          </p>
        </div>
      </section>

      {/* Why it matters */}
      <section style={{ padding: '80px 0', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 80 }} className="pub-two-col">
            <div>
              <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 16 }}>
                Why it matters
              </p>
              <p style={{ fontSize: 16, color: 'var(--fs-text-secondary)', lineHeight: 1.85 }}>
                {solution.why}
              </p>
            </div>
            <div>
              <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 16 }}>
                How FinSight approaches it
              </p>
              <p style={{ fontSize: 16, color: 'var(--fs-text-secondary)', lineHeight: 1.85 }}>
                {solution.approach}
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Example metrics */}
      <section style={{ padding: '80px 0', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 24 }}>
            Example metrics · Synthetic demonstration data
          </p>
          <div style={{ border: '1px solid var(--fs-border)', maxWidth: 600 }}>
            {solution.metrics.map((row, i) => (
              <div key={row.label} style={{
                display: 'flex', justifyContent: 'space-between', alignItems: 'center',
                padding: '16px 20px',
                background: i % 2 === 0 ? 'var(--fs-surface)' : 'transparent',
                borderBottom: i < solution.metrics.length - 1 ? '1px solid var(--fs-border)' : 'none',
              }}>
                <span style={{ fontSize: 14, color: 'var(--fs-text-secondary)' }}>{row.label}</span>
                <div style={{ textAlign: 'right' }}>
                  <span style={{ fontFamily: 'IBM Plex Mono, monospace', fontSize: 16, fontWeight: 700, color: 'var(--fs-text)' }}>{row.value}</span>
                  <div style={{ fontSize: 11, color: row.isNeg ? 'var(--fs-negative)' : 'var(--fs-positive)', marginTop: 3 }}>{row.vs}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Deeper technical */}
      <section style={{ padding: '80px 0', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 32 }}>
            Deeper detail
          </p>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(260px, 1fr))', gap: 32 }}>
            {solution.deeper.map(({ heading, body }) => (
              <div key={heading} style={{ padding: '28px', border: '1px solid var(--fs-border)' }}>
                <h3 style={{ fontSize: 16, fontWeight: 600, color: 'var(--fs-text)', marginBottom: 12 }}>{heading}</h3>
                <p style={{ fontSize: 14, color: 'var(--fs-text-secondary)', lineHeight: 1.75 }}>{body}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA */}
      <section style={{ padding: '80px 0' }}>
        <div className="pub-container">
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: 16, alignItems: 'center', justifyContent: 'space-between', padding: '40px', border: '1px solid var(--fs-border)', background: 'var(--fs-surface)' }}>
            <div>
              <h2 style={{ fontSize: 22, fontWeight: 600, color: 'var(--fs-text)', marginBottom: 8 }}>Ready to see it in practice?</h2>
              <p style={{ fontSize: 15, color: 'var(--fs-text-secondary)' }}>Walk through {solution.label} with realistic demo data.</p>
            </div>
            <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap' }}>
              <Link to="/demo" className="btn btn-primary">Request a Demo</Link>
              <Link to="/login" className="btn btn-secondary">Sign In</Link>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}

function NotFoundInline() {
  return (
    <div className="pub-container" style={{ padding: '100px 0', textAlign: 'center' }}>
      <p style={{ fontSize: 48, marginBottom: 16 }}>—</p>
      <h1 style={{ fontSize: 24, color: 'var(--fs-text)', marginBottom: 12 }}>Solution not found</h1>
      <Link to="/solutions" className="btn btn-secondary">Back to Solutions</Link>
    </div>
  )
}
