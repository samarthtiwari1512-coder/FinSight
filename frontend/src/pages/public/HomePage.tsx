import { Link } from 'react-router-dom'
import { ArrowRight, TrendingUp, DollarSign, RefreshCw, Globe, Shield } from 'lucide-react'

// Financial background watermark terms
const BG_TERMS = [
  '₹', '$', '€', '£', '¥',
  'USD', 'EUR', 'GBP', 'INR', 'JPY',
  'CASH FLOW', 'LIQUIDITY', 'RECEIVABLES',
  'PAYABLES', 'FX', 'RISK', 'WORKING CAPITAL',
  'FORECAST', 'DSO', 'DPO', 'CCC',
  '₹82.4M', '$1.2M', '+4.8%', '90 DAYS',
  'HEDGE', 'TREASURY', 'NET POSITION',
]

const MODULES = [
  {
    icon: RefreshCw,
    label: 'Working Capital',
    desc: 'DSO, DPO, DIO, and Cash Conversion Cycle tracked in a single view so you know exactly how fast your money moves.',
    href: '/solutions/working-capital',
    stat: '48d CCC',
  },
  {
    icon: TrendingUp,
    label: 'Cash Flow',
    desc: 'Current and projected liquidity across all bank accounts and entities. Know when cash gets tight before it happens.',
    href: '/solutions/cash-flow',
    stat: '₹82.4M position',
  },
  {
    icon: Globe,
    label: 'FX Exposure',
    desc: 'Consolidated foreign currency exposure, hedging coverage and sensitivity analysis in one place.',
    href: '/solutions/fx',
    stat: '$45M+ exposure',
  },
  {
    icon: Shield,
    label: 'Risk & Alerts',
    desc: 'Automated financial risk alerts, threshold monitoring and structured audit trails for compliance.',
    href: '/solutions/risk',
    stat: 'Real-time alerts',
  },
]

const STORY_SECTIONS = [
  {
    heading: 'Customers pay at different times.',
    body: `Invoices go out. Some come back in 15 days, some in 90. Meanwhile, suppliers are waiting. The gap between what you're owed and what you owe is the heartbeat of your business.`,
  },
  {
    heading: 'Currencies move while you sleep.',
    body: `A contract denominated in USD looks different when you report in INR. A 3% shift in exchange rates can add or remove millions from your effective position.`,
  },
  {
    heading: 'Forecasts need to be honest.',
    body: `Most treasury tools show you what happened. FinSight helps you project forward — with actual transaction data, aging analysis and confidence intervals.`,
  },
]

export default function HomePage() {
  return (
    <div>
      {/* ── Hero ── */}
      <section
        aria-label="Hero"
        style={{
          position: 'relative',
          minHeight: '92vh',
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
          borderBottom: '1px solid var(--fs-border)',
          overflow: 'hidden',
          padding: '80px 0',
        }}
      >
        {/* Financial background watermark */}
        <div aria-hidden="true" style={{
          position: 'absolute', inset: 0,
          pointerEvents: 'none',
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fill, minmax(120px, 1fr))',
          gap: '8px',
          padding: '40px',
          overflow: 'hidden',
          zIndex: 0,
        }}>
          {BG_TERMS.concat(BG_TERMS).map((term, i) => (
            <span key={i} style={{
              fontFamily: 'IBM Plex Mono, monospace',
              fontSize: i % 5 === 0 ? 11 : 9,
              color: 'var(--fs-text)',
              opacity: 0.03 + (i % 3) * 0.015,
              whiteSpace: 'nowrap',
              userSelect: 'none',
              letterSpacing: '0.05em',
              fontWeight: i % 4 === 0 ? 700 : 400,
            }}>{term}</span>
          ))}
        </div>

        <div className="pub-container" style={{ position: 'relative', zIndex: 1 }}>
          {/* Pre-headline tag */}
          <div style={{
            display: 'inline-flex', alignItems: 'center', gap: 8,
            marginBottom: 32,
            padding: '5px 12px',
            border: '1px solid var(--fs-border-strong)',
            fontSize: 11, fontWeight: 600,
            color: 'var(--fs-text-secondary)',
            letterSpacing: '0.06em',
            textTransform: 'uppercase',
            fontFamily: 'IBM Plex Mono, monospace',
          }}>
            <span style={{ width: 6, height: 6, background: 'var(--fs-positive)', display: 'inline-block' }} />
            Corporate Treasury Intelligence
          </div>

          {/* Main headline */}
          <h1 className="font-display" style={{
            fontSize: 'clamp(48px, 7vw, 96px)',
            lineHeight: 1.0,
            letterSpacing: '-0.02em',
            color: 'var(--fs-text)',
            marginBottom: 28,
            maxWidth: '14ch',
          }}>
            Welcome to FinSight.
          </h1>

          {/* Sub-headline */}
          <p style={{
            fontSize: 'clamp(17px, 2.2vw, 22px)',
            color: 'var(--fs-text-secondary)',
            lineHeight: 1.6,
            maxWidth: '52ch',
            marginBottom: 48,
          }}>
            A clearer way to understand where your company's money is going.
          </p>

          {/* CTAs */}
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: 12, alignItems: 'center' }}>
            <Link to="/login" className="btn btn-primary btn-lg" style={{ fontSize: 15, padding: '12px 28px' }}>
              Explore FinSight
              <ArrowRight size={16} />
            </Link>
            <Link to="/demo" className="btn btn-secondary btn-lg" style={{ fontSize: 15, padding: '12px 28px' }}>
              Request a Demo
            </Link>
          </div>

          {/* Social proof strip */}
          <div style={{
            marginTop: 64,
            paddingTop: 32,
            borderTop: '1px solid var(--fs-border)',
            display: 'flex',
            flexWrap: 'wrap',
            gap: 40,
          }}>
            {[
              { stat: '24', label: 'Bank accounts tracked' },
              { stat: '$45M+', label: 'FX exposure monitored' },
              { stat: '8', label: 'Currencies tracked' },
              { stat: 'Real-time', label: 'Liquidity alerts' },
            ].map(item => (
              <div key={item.label}>
                <div style={{ fontFamily: 'IBM Plex Mono, monospace', fontSize: 22, fontWeight: 700, color: 'var(--fs-accent)', letterSpacing: '-0.01em' }}>
                  {item.stat}
                </div>
                <div style={{ fontSize: 12, color: 'var(--fs-text-muted)', marginTop: 4 }}>
                  {item.label}
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Story: financial picture ── */}
      <section style={{ padding: '100px 0', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <div style={{ maxWidth: '64ch', marginBottom: 64 }}>
            <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 16 }}>
              The problem
            </p>
            <h2 className="font-display" style={{ fontSize: 'clamp(32px, 4vw, 52px)', color: 'var(--fs-text)', lineHeight: 1.1, marginBottom: 24 }}>
              The financial picture is always moving.
            </h2>
            <p style={{ fontSize: 17, color: 'var(--fs-text-secondary)', lineHeight: 1.8 }}>
              No single spreadsheet captures the whole thing — and by the time one does, something has changed.
            </p>
          </div>

          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))', gap: 1, border: '1px solid var(--fs-border)', background: 'var(--fs-border)' }}>
            {STORY_SECTIONS.map((s, i) => (
              <div key={i} style={{
                padding: '40px 36px',
                background: 'var(--fs-surface)',
              }}>
                <div style={{
                  fontFamily: 'IBM Plex Mono, monospace',
                  fontSize: 10, fontWeight: 700,
                  color: 'var(--fs-text-muted)',
                  letterSpacing: '0.08em',
                  marginBottom: 16,
                }}>
                  {String(i + 1).padStart(2, '0')}
                </div>
                <h3 style={{ fontSize: 20, fontWeight: 600, color: 'var(--fs-text)', lineHeight: 1.3, marginBottom: 14 }}>
                  {s.heading}
                </h3>
                <p style={{ fontSize: 14, color: 'var(--fs-text-secondary)', lineHeight: 1.75 }}>
                  {s.body}
                </p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── See the whole picture ── */}
      <section style={{ padding: '100px 0', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 80, alignItems: 'center' }} className="pub-two-col">
            <div>
              <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 16 }}>
                The solution
              </p>
              <h2 className="font-display" style={{ fontSize: 'clamp(32px, 4vw, 48px)', color: 'var(--fs-text)', lineHeight: 1.1, marginBottom: 24 }}>
                See the whole picture.
              </h2>
              <p style={{ fontSize: 16, color: 'var(--fs-text-secondary)', lineHeight: 1.8, marginBottom: 32 }}>
                FinSight brings together your cash positions, receivables, payables, FX exposure and risk into a single operational view — updated continuously, not at month-end.
              </p>
              <p style={{ fontSize: 16, color: 'var(--fs-text-secondary)', lineHeight: 1.8, marginBottom: 40 }}>
                It is designed for finance teams that need to make decisions on real numbers, not approximations.
              </p>
              <Link to="/platform" className="pub-link">
                See how the platform works <ArrowRight size={14} style={{ display: 'inline', verticalAlign: 'middle', marginLeft: 4 }} />
              </Link>
            </div>

            {/* Demo data panel */}
            <div style={{
              background: 'var(--fs-surface)',
              border: '1px solid var(--fs-border)',
              padding: 32,
            }}>
              <div style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 20, display: 'flex', justifyContent: 'space-between' }}>
                <span>Working Capital · ACME GLOBAL</span>
                <span style={{ color: 'var(--fs-accent)' }}>DEMO</span>
              </div>
              {[
                { label: 'Days Sales Outstanding (DSO)', value: '42d', vs: '+3d vs target', isNeg: true },
                { label: 'Days Payable Outstanding (DPO)', value: '38d', vs: 'On target', isNeg: false },
                { label: 'Cash Conversion Cycle (CCC)', value: '48d', vs: '−2d this month', isNeg: false },
                { label: 'Current Ratio', value: '1.42×', vs: 'Healthy', isNeg: false },
                { label: 'Total Cash (INR)', value: '₹82.4M', vs: '+4.8% vs last week', isNeg: false },
              ].map(row => (
                <div key={row.label} style={{
                  display: 'flex', justifyContent: 'space-between', alignItems: 'center',
                  padding: '12px 0',
                  borderBottom: '1px solid var(--fs-border)',
                }}>
                  <span style={{ fontSize: 13, color: 'var(--fs-text-secondary)' }}>{row.label}</span>
                  <div style={{ textAlign: 'right' }}>
                    <span style={{ fontFamily: 'IBM Plex Mono, monospace', fontSize: 14, fontWeight: 600, color: 'var(--fs-text)' }}>{row.value}</span>
                    <div style={{ fontSize: 11, color: row.isNeg ? 'var(--fs-negative)' : 'var(--fs-positive)', marginTop: 2 }}>{row.vs}</div>
                  </div>
                </div>
              ))}
              <div style={{ fontSize: 11, color: 'var(--fs-text-muted)', marginTop: 16, textAlign: 'center' }}>
                Synthetic demonstration data
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* ── Module grid ── */}
      <section style={{ padding: '100px 0', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <div style={{ maxWidth: 560, marginBottom: 64 }}>
            <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 16 }}>
              What FinSight covers
            </p>
            <h2 className="font-display" style={{ fontSize: 'clamp(28px, 3.5vw, 44px)', color: 'var(--fs-text)', lineHeight: 1.1 }}>
              Integrated modules, one operational view.
            </h2>
          </div>

          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(260px, 1fr))', gap: 1, border: '1px solid var(--fs-border)', background: 'var(--fs-border)' }}>
            {MODULES.map(({ icon: Icon, label, desc, href, stat }) => (
              <Link key={href} to={href} style={{ textDecoration: 'none' }}>
                <div style={{
                  padding: '36px 32px',
                  background: 'var(--fs-surface)',
                  height: '100%',
                  transition: 'background 0.15s',
                  cursor: 'pointer',
                }}
                  onMouseEnter={e => (e.currentTarget as HTMLElement).style.background = 'var(--fs-bg)'}
                  onMouseLeave={e => (e.currentTarget as HTMLElement).style.background = 'var(--fs-surface)'}
                >
                  <div style={{ marginBottom: 20 }}>
                    <Icon size={20} style={{ color: 'var(--fs-accent)' }} />
                  </div>
                  <div style={{ fontFamily: 'IBM Plex Mono, monospace', fontSize: 12, fontWeight: 700, color: 'var(--fs-accent)', marginBottom: 8 }}>
                    {stat}
                  </div>
                  <h3 style={{ fontSize: 18, fontWeight: 600, color: 'var(--fs-text)', marginBottom: 12 }}>
                    {label}
                  </h3>
                  <p style={{ fontSize: 13, color: 'var(--fs-text-secondary)', lineHeight: 1.7 }}>
                    {desc}
                  </p>
                  <div style={{ marginTop: 20, fontSize: 12, color: 'var(--fs-accent)', fontWeight: 500 }}>
                    Learn more →
                  </div>
                </div>
              </Link>
            ))}
          </div>
        </div>
      </section>

      {/* ── Platform ── */}
      <section style={{ padding: '100px 0', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 80, alignItems: 'start' }} className="pub-two-col">
            <div>
              <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 16 }}>
                Platform
              </p>
              <h2 className="font-display" style={{ fontSize: 'clamp(28px, 3.5vw, 44px)', color: 'var(--fs-text)', lineHeight: 1.1, marginBottom: 24 }}>
                Modules that work together.
              </h2>
              <p style={{ fontSize: 16, color: 'var(--fs-text-secondary)', lineHeight: 1.8, marginBottom: 40 }}>
                FinSight is not a collection of separate dashboards. Every module shares the same data model — so a payment recorded in Cash Management is automatically reflected in your Working Capital KPIs, FX exposure and forecasts.
              </p>
              <Link to="/platform" className="btn btn-secondary">
                View Platform Overview
              </Link>
            </div>

            {/* Module connections */}
            <div style={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
              {[
                'Cash Management & Bank Accounts',
                'Accounts Receivable & Invoicing',
                'Accounts Payable & Bills',
                'Working Capital Analytics',
                'Cash Flow Forecasting',
                'FX Exposure & Deal Management',
                'Risk Monitoring & Scenario Analysis',
                'Payments & Reconciliation',
                'Reporting & Audit',
              ].map((mod, i) => (
                <div key={mod} style={{
                  display: 'flex', alignItems: 'center', gap: 12,
                  padding: '14px 16px',
                  background: i % 2 === 0 ? 'var(--fs-surface)' : 'transparent',
                  border: '1px solid var(--fs-border)',
                }}>
                  <span style={{
                    fontFamily: 'IBM Plex Mono, monospace',
                    fontSize: 10, color: 'var(--fs-text-muted)', fontWeight: 700,
                    minWidth: 20,
                  }}>{String(i + 1).padStart(2, '0')}</span>
                  <span style={{ fontSize: 13, color: 'var(--fs-text-secondary)' }}>{mod}</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* ── Security ── */}
      <section style={{ padding: '100px 0', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <div style={{ maxWidth: 640, marginBottom: 60 }}>
            <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 16 }}>
              Security
            </p>
            <h2 className="font-display" style={{ fontSize: 'clamp(28px, 3.5vw, 44px)', color: 'var(--fs-text)', lineHeight: 1.1, marginBottom: 24 }}>
              Financial data deserves financial-grade controls.
            </h2>
            <p style={{ fontSize: 16, color: 'var(--fs-text-secondary)', lineHeight: 1.8 }}>
              FinSight implements role-based access control, JWT-based authentication, bcrypt password hashing, secure session management, input validation, and structured audit logging.
            </p>
          </div>

          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: 24 }}>
            {[
              { label: 'Role-based access', desc: 'CFO, Treasury, Finance, Auditor roles with granular permissions' },
              { label: 'Secure sessions', desc: 'JWT access tokens with refresh token rotation' },
              { label: 'Password security', desc: 'bcrypt hashing with configurable cost factors' },
              { label: 'Audit logging', desc: 'Structured audit trail for every sensitive action' },
              { label: 'Input validation', desc: 'Server-side validation on all API endpoints' },
              { label: 'Tenant isolation', desc: 'Data segregated per company entity by design' },
            ].map(item => (
              <div key={item.label} style={{ padding: '24px', border: '1px solid var(--fs-border)' }}>
                <div style={{
                  width: 8, height: 8, background: 'var(--fs-positive)',
                  marginBottom: 16,
                }} />
                <div style={{ fontSize: 14, fontWeight: 600, color: 'var(--fs-text)', marginBottom: 8 }}>
                  {item.label}
                </div>
                <div style={{ fontSize: 13, color: 'var(--fs-text-secondary)', lineHeight: 1.6 }}>
                  {item.desc}
                </div>
              </div>
            ))}
          </div>

          <div style={{ marginTop: 40 }}>
            <Link to="/security" className="pub-link">
              Read our full security overview <ArrowRight size={14} style={{ display: 'inline', verticalAlign: 'middle', marginLeft: 4 }} />
            </Link>
          </div>
        </div>
      </section>

      {/* ── CTA ── */}
      <section style={{ padding: '100px 0' }}>
        <div className="pub-container" style={{ textAlign: 'center' }}>
          <h2 className="font-display" style={{ fontSize: 'clamp(32px, 5vw, 64px)', color: 'var(--fs-text)', lineHeight: 1.05, marginBottom: 24, letterSpacing: '-0.02em' }}>
            See what your treasury actually looks like.
          </h2>
          <p style={{ fontSize: 18, color: 'var(--fs-text-secondary)', marginBottom: 48, maxWidth: '40ch', margin: '0 auto 48px' }}>
            Request a demo and walk through the platform with real synthetic data — no commitment required.
          </p>
          <div style={{ display: 'flex', gap: 12, justifyContent: 'center', flexWrap: 'wrap' }}>
            <Link to="/demo" className="btn btn-primary btn-lg" style={{ fontSize: 16, padding: '14px 32px' }}>
              Request a Demo
            </Link>
            <Link to="/login" className="btn btn-secondary btn-lg" style={{ fontSize: 16, padding: '14px 32px' }}>
              Sign In
            </Link>
          </div>
        </div>
      </section>
    </div>
  )
}
