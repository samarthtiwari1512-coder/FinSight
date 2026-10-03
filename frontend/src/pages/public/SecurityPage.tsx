import { Link } from 'react-router-dom'

const CONTROLS = [
  {
    category: 'Authentication',
    items: [
      { label: 'JWT-based authentication', detail: 'Access tokens signed with HS256. Short expiry windows (15 minutes) limit exposure if a token is intercepted.' },
      { label: 'Refresh token rotation', detail: 'Refresh tokens are rotated on every use. A reused refresh token invalidates the entire session family.' },
      { label: 'bcrypt password hashing', detail: 'All passwords are hashed using bcrypt with a configurable cost factor. Plain-text passwords are never stored.' },
      { label: 'Secure session management', detail: 'Sessions are invalidated on logout. All tokens are cryptographically random and not predictable.' },
    ],
  },
  {
    category: 'Authorization',
    items: [
      { label: 'Role-based access control', detail: 'Four primary roles: CFO, Treasury, Finance, Auditor. Each role carries specific permissions enforced at the API layer.' },
      { label: 'Tenant isolation', detail: 'All data queries are scoped to the authenticated user\'s company tenant. Cross-tenant data access is architecturally prevented.' },
      { label: 'API-level enforcement', detail: 'Authorization checks occur at every API endpoint, not just at the UI level. Route guards in the frontend are a secondary layer.' },
    ],
  },
  {
    category: 'Data Protection',
    items: [
      { label: 'Input validation', detail: 'All incoming data is validated at the API layer for type, length, format and range. Malformed input is rejected before processing.' },
      { label: 'Parameterized queries', detail: 'All database queries use parameterized statements to prevent SQL injection. Dynamic query construction is not used.' },
      { label: 'Transport security', detail: 'All connections use HTTPS/TLS. Unencrypted connections are not accepted in production.' },
    ],
  },
  {
    category: 'Monitoring & Audit',
    items: [
      { label: 'Structured audit logging', detail: 'Every sensitive action — logins, data changes, permission updates — is recorded with user identity, timestamp and context.' },
      { label: 'Immutable audit trail', detail: 'Audit log records are append-only. Existing records cannot be modified or deleted by application users.' },
      { label: 'Rate limiting', detail: 'API endpoints implement rate limiting to prevent brute-force attacks and abuse. Authentication endpoints have stricter limits.' },
    ],
  },
  {
    category: 'Infrastructure',
    items: [
      { label: 'Environment-level secrets', detail: 'Database credentials, JWT secrets and API keys are provided via environment variables. They are never hardcoded in source code.' },
      { label: 'Health monitoring', detail: 'A dedicated health endpoint monitors database connectivity, cache availability and service status.' },
      { label: 'Error handling', detail: 'Errors are caught and logged server-side. API responses return structured error codes without leaking implementation details.' },
    ],
  },
]

export default function SecurityPage() {
  return (
    <div>
      {/* Hero */}
      <section style={{ padding: '100px 0 80px', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 16 }}>
            Security
          </p>
          <h1 className="font-display" style={{ fontSize: 'clamp(36px, 5vw, 64px)', color: 'var(--fs-text)', lineHeight: 1.05, marginBottom: 24 }}>
            Financial data deserves careful protection.
          </h1>
          <p style={{ fontSize: 18, color: 'var(--fs-text-secondary)', maxWidth: '56ch', lineHeight: 1.7 }}>
            This page documents the security controls implemented in FinSight. It is written to be accurate, not aspirational. We describe what is built, not what is planned.
          </p>

          <div style={{
            marginTop: 40, padding: '16px 20px',
            border: '1px solid var(--fs-border-strong)',
            background: 'var(--fs-surface)',
            fontSize: 13, color: 'var(--fs-text-secondary)',
            lineHeight: 1.7,
            maxWidth: 600,
          }}>
            <strong style={{ color: 'var(--fs-text)' }}>Note on certifications:</strong> FinSight does not currently hold SOC 2, ISO 27001, PCI DSS or GDPR certifications. This is a demonstration platform. The controls described below are those implemented in the current codebase.
          </div>
        </div>
      </section>

      {/* Controls */}
      <section style={{ padding: '80px 0' }}>
        <div className="pub-container">
          {CONTROLS.map(({ category, items }) => (
            <div key={category} style={{ marginBottom: 64 }}>
              <h2 style={{
                fontSize: 13, fontWeight: 700,
                letterSpacing: '0.06em', textTransform: 'uppercase',
                color: 'var(--fs-text-muted)',
                fontFamily: 'IBM Plex Mono, monospace',
                marginBottom: 24,
                paddingBottom: 12,
                borderBottom: '1px solid var(--fs-border)',
              }}>
                {category}
              </h2>
              <div style={{ display: 'flex', flexDirection: 'column', gap: 1, border: '1px solid var(--fs-border)', background: 'var(--fs-border)' }}>
                {items.map(({ label, detail }) => (
                  <div key={label} style={{
                    padding: '24px 28px',
                    background: 'var(--fs-surface)',
                    display: 'grid',
                    gridTemplateColumns: '240px 1fr',
                    gap: 40,
                    alignItems: 'start',
                  }} className="pub-two-col">
                    <div style={{ display: 'flex', alignItems: 'flex-start', gap: 12 }}>
                      <span style={{ width: 8, height: 8, background: 'var(--fs-positive)', display: 'inline-block', flexShrink: 0, marginTop: 5 }} />
                      <span style={{ fontSize: 14, fontWeight: 600, color: 'var(--fs-text)' }}>{label}</span>
                    </div>
                    <p style={{ fontSize: 14, color: 'var(--fs-text-secondary)', lineHeight: 1.75 }}>{detail}</p>
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Responsible disclosure */}
      <section style={{ padding: '60px 0', borderTop: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <h2 style={{ fontSize: 20, fontWeight: 600, color: 'var(--fs-text)', marginBottom: 16 }}>Responsible Disclosure</h2>
          <p style={{ fontSize: 15, color: 'var(--fs-text-secondary)', lineHeight: 1.8, maxWidth: '56ch', marginBottom: 24 }}>
            If you discover a security issue in FinSight, please report it responsibly. Do not publish details publicly before we have had the opportunity to investigate and respond. Contact us through the form on our Contact page.
          </p>
          <Link to="/contact" className="btn btn-secondary">Contact us</Link>
        </div>
      </section>
    </div>
  )
}
