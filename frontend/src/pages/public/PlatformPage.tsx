import { Link } from 'react-router-dom'

const MODULES = [
  {
    id: 'cash',
    name: 'Cash Management',
    summary: 'Real-time consolidated cash positions across all bank accounts and entities.',
    detail: 'Aggregate balances from multiple banking relationships, view inflows and outflows by day, and track cash by entity, currency or legal structure. All positions update from transaction data — not from manual entry.',
  },
  {
    id: 'banks',
    name: 'Bank Accounts',
    summary: 'Manage all bank accounts, account types and currencies in one registry.',
    detail: 'Maintain a complete registry of bank accounts with IFSC codes, account types, entity assignments and currency designations. Each account links to transactions, cash positions and reconciliation.',
  },
  {
    id: 'receivables',
    name: 'Accounts Receivable',
    summary: 'Invoice-level tracking of all outstanding customer receivables.',
    detail: 'Track every customer invoice from creation to settlement. Aging analysis shows current, 30d, 60d, 90d+ buckets. Overdue alerts trigger automatically when payment is past due date.',
  },
  {
    id: 'payables',
    name: 'Accounts Payable',
    summary: 'Structured visibility into supplier payables and upcoming payment obligations.',
    detail: 'View all supplier bills, their due dates and outstanding amounts. Identify cash outflows due in 3, 7 and 30 days to plan liquidity. DPO analytics track your payable strategy effectiveness.',
  },
  {
    id: 'working-capital',
    name: 'Working Capital Analytics',
    summary: 'DSO, DPO, DIO and CCC — calculated from actual transaction data.',
    detail: 'Key working capital metrics are computed continuously from live data. Trend charts show how the metrics move over time. Benchmarks highlight where performance deviates from targets.',
  },
  {
    id: 'forecast',
    name: 'Cash Flow Forecasting',
    summary: '30-day forward projection of inflows, outflows and closing balance.',
    detail: 'Statistical forecasting based on historical transaction patterns and confirmed receivables/payables. Confidence intervals and liquidity warnings are included. Disclaimer clearly identifies the forecast as a model estimate.',
  },
  {
    id: 'fx',
    name: 'FX Exposure',
    summary: 'Consolidated foreign currency exposure with hedging coverage tracking.',
    detail: 'Aggregate FX exposure by currency from receivables, payables and cash positions. Risk classification uses exposure size and hedging coverage. Net exposure versus hedged position shown side-by-side.',
  },
  {
    id: 'deals',
    name: 'FX Deal Management',
    summary: 'Record, track and manage FX forward deals and their settlement dates.',
    detail: 'Maintain a register of active FX deals including forward rates, settlement dates, notional amounts and counterparty banks. FX deals reduce your net unhedged exposure automatically.',
  },
  {
    id: 'risk',
    name: 'Risk Monitoring',
    summary: 'Threshold-based financial risk alerts with severity classification.',
    detail: 'Define thresholds for liquidity, FX exposure, receivables aging and other risk factors. The risk engine evaluates conditions and generates structured alerts classified as LOW, MEDIUM, HIGH or CRITICAL.',
  },
  {
    id: 'scenarios',
    name: 'Scenario Analysis',
    summary: 'Model the impact of hypothetical events on your financial position.',
    detail: 'Simulate scenarios such as a major customer defaulting, a currency devaluation or a large unexpected payable. Results show the impact on working capital, cash position and FX exposure.',
  },
  {
    id: 'payments',
    name: 'Payments',
    summary: 'Transaction-level payment management across all bank accounts.',
    detail: 'Record and categorize payments by type, entity, bank account and counterparty. Payments link automatically to receivables, payables and cash positions for accurate real-time balances.',
  },
  {
    id: 'reconciliation',
    name: 'Reconciliation',
    summary: 'Match bank transactions to system records and identify exceptions.',
    detail: 'Bank statement reconciliation against system-recorded transactions. Unmatched items are flagged for review. Reconciliation status provides assurance that your cash position is accurate.',
  },
  {
    id: 'reporting',
    name: 'Reports',
    summary: 'Structured financial reports exportable for management and compliance use.',
    detail: 'Generate treasury reports covering cash position, working capital, FX exposure and risk summary. Export to structured formats for distribution to management or external parties.',
  },
  {
    id: 'audit',
    name: 'Audit Logs',
    summary: 'Complete structured audit trail of all user actions and data changes.',
    detail: 'Every sensitive action in FinSight — data changes, user logins, permission changes — is recorded with timestamp, user identity and context. Audit logs are immutable and available for compliance review.',
  },
]

export default function PlatformPage() {
  return (
    <div>
      {/* Hero */}
      <section style={{ padding: '100px 0 80px', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 16 }}>
            Platform
          </p>
          <h1 className="font-display" style={{ fontSize: 'clamp(36px, 5vw, 64px)', color: 'var(--fs-text)', lineHeight: 1.05, marginBottom: 24 }}>
            One platform. The complete financial picture.
          </h1>
          <p style={{ fontSize: 18, color: 'var(--fs-text-secondary)', maxWidth: '56ch', lineHeight: 1.7 }}>
            FinSight's modules are built on a single shared data model. A payment recorded in Cash Management instantly updates working capital analytics, FX exposure and forecasts. Nothing falls through the cracks between disconnected tools.
          </p>
        </div>
      </section>

      {/* How modules connect */}
      <section style={{ padding: '80px 0', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 32 }}>
            How the modules connect
          </p>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))', gap: 1, background: 'var(--fs-border)', border: '1px solid var(--fs-border)' }}>
            {[
              { from: 'Bank Accounts', to: 'Cash Management', note: 'Balance feeds into position' },
              { from: 'Invoices', to: 'Accounts Receivable', note: 'Invoice → aging → DSO' },
              { from: 'Bills', to: 'Accounts Payable', note: 'Bill → payment timing → DPO' },
              { from: 'AR + AP + Cash', to: 'Working Capital', note: 'DSO, DPO, DIO, CCC computed' },
              { from: 'Working Capital', to: 'Cash Forecast', note: 'Transaction patterns drive projection' },
              { from: 'AR + AP + Cash', to: 'FX Exposure', note: 'Currency positions aggregated' },
              { from: 'FX Exposure', to: 'Risk Monitoring', note: 'Thresholds trigger risk events' },
              { from: 'Risk Events', to: 'Alerts', note: 'Prioritized with severity' },
            ].map(({ from, to, note }) => (
              <div key={from + to} style={{ padding: '20px 24px', background: 'var(--fs-surface)' }}>
                <div style={{ fontFamily: 'IBM Plex Mono, monospace', fontSize: 12, color: 'var(--fs-text)', marginBottom: 4 }}>
                  {from} → {to}
                </div>
                <div style={{ fontSize: 12, color: 'var(--fs-text-muted)' }}>{note}</div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Module index */}
      <section id="modules" style={{ padding: '80px 0' }}>
        <div className="pub-container">
          <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 32 }}>
            All modules
          </p>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 1, border: '1px solid var(--fs-border)', background: 'var(--fs-border)' }}>
            {MODULES.map(({ id, name, summary, detail }) => (
              <div key={id} id={id} style={{ padding: '32px 36px', background: 'var(--fs-surface)' }}>
                <div style={{ display: 'grid', gridTemplateColumns: '2fr 3fr', gap: 48 }} className="pub-two-col">
                  <div>
                    <h2 style={{ fontSize: 18, fontWeight: 600, color: 'var(--fs-text)', marginBottom: 10 }}>{name}</h2>
                    <p style={{ fontSize: 14, color: 'var(--fs-text-secondary)', lineHeight: 1.7, fontStyle: 'italic' }}>{summary}</p>
                  </div>
                  <p style={{ fontSize: 14, color: 'var(--fs-text-secondary)', lineHeight: 1.8 }}>{detail}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA */}
      <section style={{ padding: '80px 0', borderTop: '1px solid var(--fs-border)' }}>
        <div className="pub-container" style={{ display: 'flex', flexWrap: 'wrap', gap: 24, alignItems: 'center', justifyContent: 'space-between' }}>
          <div>
            <h2 className="font-display" style={{ fontSize: 32, color: 'var(--fs-text)', marginBottom: 8 }}>See the full platform in action.</h2>
            <p style={{ fontSize: 16, color: 'var(--fs-text-secondary)' }}>Request a demo walkthrough with realistic synthetic data.</p>
          </div>
          <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap' }}>
            <Link to="/demo" className="btn btn-primary">Request a Demo</Link>
            <Link to="/login" className="btn btn-secondary">Sign In</Link>
          </div>
        </div>
      </section>
    </div>
  )
}
