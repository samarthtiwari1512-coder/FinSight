import { Link } from 'react-router-dom'

export default function AboutPage() {
  return (
    <div>
      {/* Hero */}
      <section style={{ padding: '100px 0 80px', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 16 }}>
            About
          </p>
          <h1 className="font-display" style={{ fontSize: 'clamp(36px, 5vw, 64px)', color: 'var(--fs-text)', lineHeight: 1.05, marginBottom: 24 }}>
            FinSight is a treasury intelligence platform built to bring clarity to corporate finance.
          </h1>
        </div>
      </section>

      {/* What it is */}
      <section style={{ padding: '80px 0', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 80 }} className="pub-two-col">
            <div>
              <h2 style={{ fontSize: 22, fontWeight: 600, color: 'var(--fs-text)', marginBottom: 24 }}>What FinSight is</h2>
              <p style={{ fontSize: 16, color: 'var(--fs-text-secondary)', lineHeight: 1.85, marginBottom: 20 }}>
                FinSight is a full-stack financial management platform covering corporate treasury, working capital, accounts receivable, accounts payable, FX exposure and risk monitoring.
              </p>
              <p style={{ fontSize: 16, color: 'var(--fs-text-secondary)', lineHeight: 1.85 }}>
                It is designed for finance teams — CFOs, treasury managers, finance controllers and auditors — who need operational clarity on where the company's money is, where it's going, and what risks it faces.
              </p>
            </div>
            <div>
              <h2 style={{ fontSize: 22, fontWeight: 600, color: 'var(--fs-text)', marginBottom: 24 }}>Why it exists</h2>
              <p style={{ fontSize: 16, color: 'var(--fs-text-secondary)', lineHeight: 1.85, marginBottom: 20 }}>
                The financial picture inside most companies is distributed across spreadsheets, banking portals, ERP exports and email threads. Nothing gives a coherent view across all of these at once.
              </p>
              <p style={{ fontSize: 16, color: 'var(--fs-text-secondary)', lineHeight: 1.85 }}>
                FinSight exists to be that coherent view. One system where cash positions, receivables, payables, FX exposure, working capital metrics and forecasts all live together and update in real time.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Problems it addresses */}
      <section style={{ padding: '80px 0', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <h2 style={{ fontSize: 22, fontWeight: 600, color: 'var(--fs-text)', marginBottom: 40 }}>
            The financial problems it addresses
          </h2>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 1, border: '1px solid var(--fs-border)', background: 'var(--fs-border)' }}>
            {[
              {
                problem: 'Delayed cash visibility',
                detail: 'Treasury teams often know their cash position hours or days after the fact. FinSight consolidates bank account data into a real-time view across all entities.',
              },
              {
                problem: 'Working capital blind spots',
                detail: 'DSO, DPO and CCC are often computed manually at month-end. FinSight calculates them continuously from transaction data so trends are visible as they develop.',
              },
              {
                problem: 'Undiscovered FX exposure',
                detail: 'Companies with international receivables and payables accumulate FX exposure without a clear view of the net position. FinSight aggregates this across all currencies.',
              },
              {
                problem: 'Reactive risk management',
                detail: 'Financial risks tend to be identified only after they crystallize. FinSight monitors configurable thresholds and generates alerts when conditions warrant attention.',
              },
              {
                problem: 'Disconnected financial tools',
                detail: 'Cash data lives in one system, invoices in another, payments in a third. FinSight provides a shared data model across all treasury functions.',
              },
            ].map(({ problem, detail }) => (
              <div key={problem} style={{
                padding: '28px 36px',
                background: 'var(--fs-surface)',
                display: 'grid',
                gridTemplateColumns: '260px 1fr',
                gap: 40,
                alignItems: 'start',
              }} className="pub-two-col">
                <h3 style={{ fontSize: 15, fontWeight: 600, color: 'var(--fs-text)' }}>{problem}</h3>
                <p style={{ fontSize: 14, color: 'var(--fs-text-secondary)', lineHeight: 1.8 }}>{detail}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Philosophy */}
      <section style={{ padding: '80px 0', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <div style={{ maxWidth: '64ch', margin: '0 auto', textAlign: 'center' }}>
            <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 24 }}>
              Product philosophy
            </p>
            <blockquote className="font-display" style={{ fontSize: 'clamp(24px, 3.5vw, 36px)', color: 'var(--fs-text)', lineHeight: 1.3, marginBottom: 32 }}>
              "Financial software should tell you what is actually happening — not what you want to hear."
            </blockquote>
            <p style={{ fontSize: 16, color: 'var(--fs-text-secondary)', lineHeight: 1.85 }}>
              FinSight is designed around accuracy over optimism. Forecasts come with confidence intervals. Risk alerts include severity classification. Metrics are computed from transaction data, not from manually entered targets.
            </p>
          </div>
        </div>
      </section>

      {/* Transparency */}
      <section style={{ padding: '80px 0', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <h2 style={{ fontSize: 22, fontWeight: 600, color: 'var(--fs-text)', marginBottom: 24 }}>Approach to transparency</h2>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 40 }} className="pub-two-col">
            {[
              {
                heading: 'Honest data labelling',
                body: 'Demo data is clearly labelled as synthetic. Forecast data includes disclaimers. Model estimates are distinguished from confirmed figures.',
              },
              {
                heading: 'Accurate security claims',
                body: 'The Security page describes implemented controls only — not aspirational ones. We do not claim certifications we do not hold.',
              },
              {
                heading: 'No invented social proof',
                body: 'FinSight does not invent customer counts, revenue figures or years of experience. What is described is what exists.',
              },
              {
                heading: 'Simulation disclosure',
                body: 'FinSight is clearly disclosed as a simulation platform. No real financial transactions are executed through the system.',
              },
            ].map(({ heading, body }) => (
              <div key={heading} style={{ padding: '28px', border: '1px solid var(--fs-border)' }}>
                <h3 style={{ fontSize: 15, fontWeight: 600, color: 'var(--fs-text)', marginBottom: 12 }}>{heading}</h3>
                <p style={{ fontSize: 14, color: 'var(--fs-text-secondary)', lineHeight: 1.75 }}>{body}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA */}
      <section style={{ padding: '80px 0' }}>
        <div className="pub-container">
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: 16, alignItems: 'center' }}>
            <Link to="/demo" className="btn btn-primary btn-lg">Request a Demo</Link>
            <Link to="/contact" className="btn btn-secondary btn-lg">Get in Touch</Link>
          </div>
        </div>
      </section>
    </div>
  )
}
