import { Link } from 'react-router-dom'

const ARTICLES = [
  {
    slug: 'understanding-cash-conversion-cycle',
    category: 'Working Capital',
    title: 'Understanding the Cash Conversion Cycle',
    summary: 'The CCC is one of the most revealing metrics in corporate finance. Here\'s what it measures, why it matters, and how to use it operationally.',
    date: 'October 2025',
    readTime: '6 min read',
    body: `The Cash Conversion Cycle (CCC) measures how many days it takes a company to convert its investments in inventory and other resources into cash flows from sales.

The formula is simple: CCC = DSO + DIO − DPO

Where:
- DSO (Days Sales Outstanding): How long it takes customers to pay you
- DIO (Days Inventory Outstanding): How long inventory sits before it's sold
- DPO (Days Payable Outstanding): How long you take to pay your suppliers

A lower CCC is generally better — it means you're turning investment into cash faster. A company with a CCC of 20 days is much more capital-efficient than one with 80 days.

**Why it matters operationally**

A rising CCC is often the earliest signal of business stress. If DSO is climbing, customers may be struggling to pay. If DIO is rising, inventory may be building up unsold. If DPO is falling, you may be paying suppliers too quickly.

**Benchmarks vary significantly by industry**

Retail often operates with near-zero or negative CCC (supermarkets collect from customers before paying suppliers). Manufacturing and services companies typically run 30–90 days. Professional services can have CCCs above 60 days if invoicing and collection processes are slow.

**How FinSight tracks CCC**

FinSight calculates CCC components from actual transaction data — not from manually entered figures. DSO is computed from invoice dates and payment receipt dates. DPO from bill dates and payment sent dates. DIO from inventory movement data where available.

This means the CCC in FinSight reflects reality, not plan.`,
  },
  {
    slug: 'fx-exposure-primer',
    category: 'FX',
    title: 'FX Exposure: A Practical Primer for Finance Teams',
    summary: 'Currency risk doesn\'t announce itself. This article explains how FX exposure builds up, how to measure it, and what your options are.',
    date: 'September 2025',
    readTime: '8 min read',
    body: `FX exposure arises when a company has assets, liabilities, income or costs denominated in a foreign currency. The exposure becomes a problem when exchange rates move in an unfavorable direction before those positions are settled.

**Three types of FX exposure**

Transaction exposure: The risk of exchange rate changes affecting the value of confirmed contracts and payments. A USD invoice from an Indian exporter is subject to transaction exposure until it's paid and converted to INR.

Translation exposure: The risk of exchange rate changes affecting the reported value of foreign subsidiary assets and liabilities when consolidated into the parent company's financial statements.

Economic exposure: The longer-term risk that sustained exchange rate changes will affect a company's competitiveness and future cash flows.

**Why many companies underestimate their exposure**

FX exposure often accumulates across multiple systems — receivables in one, payables in another, cash in a third. Without a consolidated view, the net position is unclear. A company might have $10M in USD receivables and $8M in USD payables, giving a net exposure of $2M — but if these are tracked in separate systems, the $2M net figure never gets computed.

**Hedging is a risk management decision, not a profit strategy**

Hedging FX exposure reduces volatility — it doesn't generate returns. A forward contract locks in an exchange rate, eliminating uncertainty about a future payment. Whether hedging makes sense depends on the size of the exposure, the cost of hedging instruments, and the company's risk tolerance.

**FinSight's approach**

FinSight aggregates FX exposure from receivables, payables and cash balances by currency, showing the net position in base currency. It tracks hedging coverage and flags currencies where exposure exceeds hedged amounts.`,
  },
  {
    slug: 'treasury-technology-landscape',
    category: 'Treasury Technology',
    title: 'The Treasury Technology Landscape in 2025',
    summary: 'From spreadsheets to TMS platforms to modern API-first tools — how corporate treasury technology has evolved and where it\'s headed.',
    date: 'August 2025',
    readTime: '7 min read',
    body: `Most corporate treasury functions operate on a spectrum from spreadsheet-heavy to ERP-integrated to dedicated Treasury Management Systems (TMS). Each has trade-offs.

**The spreadsheet era**

Many mid-market finance teams still manage cash positions, working capital and FX exposure in spreadsheets. This works — until it doesn't. Spreadsheets break under scale, introduce version control problems, require manual reconciliation, and provide no real-time visibility.

**ERP-integrated treasury**

Enterprise systems like SAP and Oracle include treasury modules. These are powerful but expensive, require significant implementation effort, and often lag behind in usability. They're typically adopted by large enterprises with dedicated treasury IT teams.

**Dedicated TMS platforms**

Treasury Management Systems like Kyriba, FIS Quantum and TreasuryXpress sit between ERP suites and simple tools. They specialize in cash positioning, bank connectivity, FX and payments. Implementation is complex and licensing is significant.

**Modern API-first approaches**

More recent tools — including FinSight — take a different approach: build on modern web APIs, connect to banking infrastructure programmatically, and present the treasury picture through a clean, purpose-built interface. The goal is to make treasury visibility accessible without the implementation overhead of traditional TMS.

**What hasn't changed**

The underlying financial requirements haven't changed. Treasury teams still need to know their cash position, their working capital cycle, their FX exposure, and their liquidity outlook. The question is just how accurately and how quickly they can access these numbers.`,
  },
  {
    slug: 'dso-reduction-strategies',
    category: 'Working Capital',
    title: 'Practical Strategies for Reducing DSO',
    summary: 'Days Sales Outstanding is one of the highest-leverage levers in working capital management. These are the approaches that actually work.',
    date: 'July 2025',
    readTime: '5 min read',
    body: `DSO (Days Sales Outstanding) measures how many days it takes to collect payment after an invoice is issued. A DSO of 45 days means, on average, customers take 45 days to pay.

Reducing DSO directly improves cash flow. If a company with ₹500M in annual revenue reduces DSO from 60 to 45 days, it unlocks approximately ₹20M in cash.

**1. Invoice promptly and accurately**

Late or incorrect invoices delay payment. Invoicing within 24 hours of service delivery, with all required reference numbers and billing details, eliminates one of the most common payment delays.

**2. Offer early payment incentives**

A 1% discount for payment within 10 days (1/10 net 30) can accelerate cash flow significantly. Whether this makes financial sense depends on the cost of capital versus the cost of the discount.

**3. Identify slow-paying customers early**

Aging analysis at the customer level reveals which customers consistently pay late. Early identification allows proactive follow-up before invoices become significantly overdue.

**4. Automate reminders**

Manual follow-up on overdue invoices is inconsistent and time-consuming. Automated reminders at 7, 14 and 30 days past due, with escalation rules, improve collection rates without increasing headcount.

**5. Review credit terms regularly**

Credit terms set at account opening often don't get reviewed as relationships evolve. A customer who now represents 15% of AR might warrant tighter terms or a credit limit review.

**How FinSight supports DSO reduction**

FinSight tracks aging at the invoice and customer level, shows DSO trend over time, and alerts when specific customers' payment patterns deteriorate. The goal is to make overdue situations visible before they become bad debt.`,
  },
]

export default function InsightsPage() {
  return (
    <div>
      {/* Hero */}
      <section style={{ padding: '100px 0 80px', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 16 }}>
            Insights
          </p>
          <h1 className="font-display" style={{ fontSize: 'clamp(36px, 5vw, 64px)', color: 'var(--fs-text)', lineHeight: 1.05, marginBottom: 24 }}>
            Writing on corporate finance and treasury.
          </h1>
          <p style={{ fontSize: 18, color: 'var(--fs-text-secondary)', maxWidth: '52ch', lineHeight: 1.7 }}>
            Practical thinking on working capital, cash management, FX, risk and treasury technology.
          </p>
        </div>
      </section>

      {/* Article list */}
      <section style={{ padding: '80px 0' }}>
        <div className="pub-container">
          <div style={{ display: 'flex', flexDirection: 'column', gap: 1, border: '1px solid var(--fs-border)', background: 'var(--fs-border)' }}>
            {ARTICLES.map(article => (
              <Link key={article.slug} to={`/insights/${article.slug}`} style={{ textDecoration: 'none' }}>
                <article style={{
                  padding: '40px 48px',
                  background: 'var(--fs-surface)',
                  transition: 'background 0.15s',
                }}
                  onMouseEnter={e => (e.currentTarget as HTMLElement).style.background = 'var(--fs-bg)'}
                  onMouseLeave={e => (e.currentTarget as HTMLElement).style.background = 'var(--fs-surface)'}
                >
                  <div style={{ display: 'flex', gap: 16, alignItems: 'center', marginBottom: 16 }}>
                    <span style={{
                      fontSize: 11, fontWeight: 700,
                      color: 'var(--fs-accent)',
                      letterSpacing: '0.05em',
                      textTransform: 'uppercase',
                      fontFamily: 'IBM Plex Mono, monospace',
                    }}>
                      {article.category}
                    </span>
                    <span style={{ color: 'var(--fs-border-strong)', userSelect: 'none' }}>·</span>
                    <span style={{ fontSize: 12, color: 'var(--fs-text-muted)' }}>{article.date}</span>
                    <span style={{ color: 'var(--fs-border-strong)', userSelect: 'none' }}>·</span>
                    <span style={{ fontSize: 12, color: 'var(--fs-text-muted)' }}>{article.readTime}</span>
                  </div>
                  <h2 style={{ fontSize: 22, fontWeight: 600, color: 'var(--fs-text)', marginBottom: 12, lineHeight: 1.3 }}>
                    {article.title}
                  </h2>
                  <p style={{ fontSize: 15, color: 'var(--fs-text-secondary)', lineHeight: 1.7 }}>
                    {article.summary}
                  </p>
                  <div style={{ marginTop: 20, fontSize: 13, color: 'var(--fs-accent)', fontWeight: 500 }}>
                    Read →
                  </div>
                </article>
              </Link>
            ))}
          </div>
        </div>
      </section>
    </div>
  )
}

// Individual article page
export function InsightArticlePage({ slug }: { slug: string }) {
  const article = ARTICLES.find(a => a.slug === slug)
  if (!article) {
    return (
      <div className="pub-container" style={{ padding: '100px 0', textAlign: 'center' }}>
        <h1 style={{ color: 'var(--fs-text)', marginBottom: 24 }}>Article not found</h1>
        <Link to="/insights" className="btn btn-secondary">Back to Insights</Link>
      </div>
    )
  }

  return (
    <div>
      <section style={{ padding: '80px 0', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container" style={{ maxWidth: 780 }}>
          <div style={{ display: 'flex', gap: 16, alignItems: 'center', marginBottom: 32 }}>
            <Link to="/insights" style={{ fontSize: 13, color: 'var(--fs-text-muted)', textDecoration: 'none' }}>Insights</Link>
            <span style={{ color: 'var(--fs-text-muted)' }}>/</span>
            <span style={{
              fontSize: 11, fontWeight: 700,
              color: 'var(--fs-accent)',
              letterSpacing: '0.05em',
              textTransform: 'uppercase',
            }}>
              {article.category}
            </span>
          </div>

          <h1 className="font-display" style={{ fontSize: 'clamp(28px, 4vw, 48px)', color: 'var(--fs-text)', lineHeight: 1.15, marginBottom: 20 }}>
            {article.title}
          </h1>

          <p style={{ fontSize: 18, color: 'var(--fs-text-secondary)', lineHeight: 1.7, marginBottom: 32, borderLeft: '2px solid var(--fs-accent)', paddingLeft: 20 }}>
            {article.summary}
          </p>

          <div style={{ display: 'flex', gap: 24, fontSize: 12, color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 60, paddingBottom: 40, borderBottom: '1px solid var(--fs-border)' }}>
            <span>{article.date}</span>
            <span>{article.readTime}</span>
          </div>

          <div style={{ fontSize: 16, color: 'var(--fs-text-secondary)', lineHeight: 1.95 }}>
            {article.body.split('\n\n').map((para, i) => {
              if (para.startsWith('**') && para.endsWith('**')) {
                return (
                  <h3 key={i} style={{ fontSize: 18, fontWeight: 600, color: 'var(--fs-text)', margin: '36px 0 16px' }}>
                    {para.slice(2, -2)}
                  </h3>
                )
              }
              if (para.startsWith('- ')) {
                return (
                  <ul key={i} style={{ listStyle: 'none', padding: 0, marginBottom: 20 }}>
                    {para.split('\n').map((line, j) => (
                      <li key={j} style={{ display: 'flex', gap: 12, marginBottom: 8 }}>
                        <span style={{ color: 'var(--fs-accent)', fontWeight: 700, flexShrink: 0 }}>—</span>
                        <span>{line.slice(2)}</span>
                      </li>
                    ))}
                  </ul>
                )
              }
              // Handle **bold** inside paragraph
              const parts = para.split(/(\*\*[^*]+\*\*)/g)
              return (
                <p key={i} style={{ marginBottom: 20 }}>
                  {parts.map((part, j) =>
                    part.startsWith('**') ? <strong key={j} style={{ color: 'var(--fs-text)', fontWeight: 600 }}>{part.slice(2, -2)}</strong> : part
                  )}
                </p>
              )
            })}
          </div>
        </div>
      </section>

      <section style={{ padding: '60px 0' }}>
        <div className="pub-container" style={{ maxWidth: 780 }}>
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: 16, alignItems: 'center', justifyContent: 'space-between' }}>
            <Link to="/insights" className="pub-link">← All articles</Link>
            <div style={{ display: 'flex', gap: 12 }}>
              <Link to="/demo" className="btn btn-primary">Request a Demo</Link>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}
