import { Link } from 'react-router-dom'

function LegalSection({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div style={{ marginBottom: 48 }}>
      <h2 style={{ fontSize: 18, fontWeight: 600, color: 'var(--fs-text)', marginBottom: 16, paddingBottom: 12, borderBottom: '1px solid var(--fs-border)' }}>{title}</h2>
      <div style={{ fontSize: 15, color: 'var(--fs-text-secondary)', lineHeight: 1.85 }}>{children}</div>
    </div>
  )
}

export function PrivacyPage() {
  return (
    <div>
      <section style={{ padding: '80px 0 60px', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 16 }}>Legal</p>
          <h1 className="font-display" style={{ fontSize: 'clamp(32px, 4vw, 52px)', color: 'var(--fs-text)', lineHeight: 1.1, marginBottom: 12 }}>Privacy Policy</h1>
          <p style={{ fontSize: 13, color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace' }}>Last updated: October 2025</p>
        </div>
      </section>
      <section style={{ padding: '60px 0 80px' }}>
        <div className="pub-container" style={{ maxWidth: 780 }}>
          <div style={{ padding: '16px 20px', background: 'var(--fs-surface)', border: '1px solid var(--fs-border)', marginBottom: 40, fontSize: 14, color: 'var(--fs-text-secondary)', lineHeight: 1.7 }}>
            FinSight is a demonstration and simulation platform. No real financial transactions are processed. Data entered is used solely within the simulation environment.
          </div>
          <LegalSection title="What data FinSight collects">
            <p style={{ marginBottom: 16 }}>FinSight collects the following categories of data:</p>
            <ul style={{ listStyle: 'none', padding: 0 }}>
              {['Account information: email address, hashed password, name, role assignment', 'Session data: JWT tokens stored in browser memory and localStorage for refresh purposes', 'Usage data: actions performed within the platform, recorded in audit logs', 'Form submissions: demo requests and contact form entries'].map(item => (
                <li key={item} style={{ display: 'flex', gap: 12, marginBottom: 12 }}>
                  <span style={{ color: 'var(--fs-accent)', fontWeight: 700, flexShrink: 0 }}>—</span>
                  <span>{item}</span>
                </li>
              ))}
            </ul>
          </LegalSection>
          <LegalSection title="How data is used">
            <p style={{ marginBottom: 16 }}>Data is used exclusively to operate the FinSight simulation platform. It is not sold to third parties, used for advertising, or shared with external services beyond infrastructure providers required to operate the platform.</p>
          </LegalSection>
          <LegalSection title="Data retention">
            <p>Account data is retained for as long as the account is active. Audit log data is retained to maintain the integrity of the audit trail. You may request deletion of your account data by contacting us.</p>
          </LegalSection>
          <LegalSection title="Security">
            <p>Passwords are stored as bcrypt hashes. Access tokens have short expiry windows. All connections use HTTPS. For the complete list of security controls, see our <Link to="/security" style={{ color: 'var(--fs-accent)' }}>Security page</Link>.</p>
          </LegalSection>
          <LegalSection title="Contact">
            <p>For privacy-related enquiries, use the <Link to="/contact" style={{ color: 'var(--fs-accent)' }}>Contact page</Link>.</p>
          </LegalSection>
        </div>
      </section>
    </div>
  )
}

export function TermsPage() {
  return (
    <div>
      <section style={{ padding: '80px 0 60px', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 16 }}>Legal</p>
          <h1 className="font-display" style={{ fontSize: 'clamp(32px, 4vw, 52px)', color: 'var(--fs-text)', lineHeight: 1.1, marginBottom: 12 }}>Terms of Use</h1>
          <p style={{ fontSize: 13, color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace' }}>Last updated: October 2025</p>
        </div>
      </section>
      <section style={{ padding: '60px 0 80px' }}>
        <div className="pub-container" style={{ maxWidth: 780 }}>
          <LegalSection title="Nature of the platform">
            <p>FinSight is a financial intelligence simulation platform. No real financial transactions are executed. All financial data within the platform is synthetic demonstration data. FinSight makes no representations about the accuracy of demonstration data for any real-world financial decision.</p>
          </LegalSection>
          <LegalSection title="Acceptable use">
            <p style={{ marginBottom: 16 }}>You may use FinSight to evaluate, demonstrate and explore the functionality of a corporate treasury intelligence platform. You may not:</p>
            <ul style={{ listStyle: 'none', padding: 0 }}>
              {[
                'Attempt to access other users\' accounts or data',
                'Use the platform to process real financial transactions',
                'Attempt to reverse-engineer or extract proprietary systems',
                'Use automated tools to scrape or abuse the platform',
              ].map(item => (
                <li key={item} style={{ display: 'flex', gap: 12, marginBottom: 12 }}>
                  <span style={{ color: 'var(--fs-negative)', fontWeight: 700, flexShrink: 0 }}>—</span>
                  <span>{item}</span>
                </li>
              ))}
            </ul>
          </LegalSection>
          <LegalSection title="Limitation of liability">
            <p>FinSight is provided as a demonstration platform. No warranty is made regarding its fitness for any particular purpose. FinSight and its creators accept no liability for decisions made based on data within the simulation environment.</p>
          </LegalSection>
          <LegalSection title="Changes to these terms">
            <p>These terms may be updated. Continued use of the platform after any update constitutes acceptance of the revised terms.</p>
          </LegalSection>
        </div>
      </section>
    </div>
  )
}

export function CookiesPage() {
  return (
    <div>
      <section style={{ padding: '80px 0 60px', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 16 }}>Legal</p>
          <h1 className="font-display" style={{ fontSize: 'clamp(32px, 4vw, 52px)', color: 'var(--fs-text)', lineHeight: 1.1, marginBottom: 12 }}>Cookie Policy</h1>
          <p style={{ fontSize: 13, color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace' }}>Last updated: October 2025</p>
        </div>
      </section>
      <section style={{ padding: '60px 0 80px' }}>
        <div className="pub-container" style={{ maxWidth: 780 }}>
          <LegalSection title="What FinSight uses">
            <p style={{ marginBottom: 16 }}>FinSight uses browser localStorage (not cookies in the traditional sense) for session management purposes only. Specifically:</p>
            <ul style={{ listStyle: 'none', padding: 0 }}>
              {[
                'refresh_token: A refresh token stored in localStorage to maintain your login session. This token is rotated on each use.',
                'Theme preference: Your light/dark theme preference, stored locally in the browser.',
              ].map(item => (
                <li key={item} style={{ display: 'flex', gap: 12, marginBottom: 12 }}>
                  <span style={{ color: 'var(--fs-accent)', fontWeight: 700, flexShrink: 0 }}>—</span>
                  <span>{item}</span>
                </li>
              ))}
            </ul>
          </LegalSection>
          <LegalSection title="What FinSight does not use">
            <p>FinSight does not use tracking cookies, advertising cookies, analytics third-party cookies, or any form of cross-site tracking. There are no advertising or marketing pixels on this platform.</p>
          </LegalSection>
          <LegalSection title="Clearing your data">
            <p>You can clear all locally stored session data by logging out or clearing your browser's localStorage for this domain.</p>
          </LegalSection>
        </div>
      </section>
    </div>
  )
}

export function AccessibilityPage() {
  return (
    <div>
      <section style={{ padding: '80px 0 60px', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 16 }}>Legal</p>
          <h1 className="font-display" style={{ fontSize: 'clamp(32px, 4vw, 52px)', color: 'var(--fs-text)', lineHeight: 1.1, marginBottom: 12 }}>Accessibility</h1>
          <p style={{ fontSize: 13, color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace' }}>Last updated: October 2025</p>
        </div>
      </section>
      <section style={{ padding: '60px 0 80px' }}>
        <div className="pub-container" style={{ maxWidth: 780 }}>
          <LegalSection title="Our commitment">
            <p>FinSight is designed with accessibility in mind. We aim to ensure that the platform is usable by as many people as possible, including those who use assistive technologies.</p>
          </LegalSection>
          <LegalSection title="Implementation">
            <p style={{ marginBottom: 16 }}>Current accessibility implementations include:</p>
            <ul style={{ listStyle: 'none', padding: 0 }}>
              {[
                'Semantic HTML5 elements used throughout the public website and application',
                'ARIA labels on interactive elements where HTML semantics are insufficient',
                'Keyboard navigation supported on all navigation and form elements',
                'Visible focus states on all focusable elements',
                'Sufficient color contrast ratios between text and background (target: WCAG AA)',
                'Form labels properly associated with their controls',
                'Error messages accessible to screen readers via role="alert"',
                'Reduced-motion support via prefers-reduced-motion media query',
                'Status is not communicated by color alone — text labels accompany color indicators',
              ].map(item => (
                <li key={item} style={{ display: 'flex', gap: 12, marginBottom: 12 }}>
                  <span style={{ color: 'var(--fs-positive)', fontWeight: 700, flexShrink: 0 }}>✓</span>
                  <span>{item}</span>
                </li>
              ))}
            </ul>
          </LegalSection>
          <LegalSection title="Known limitations">
            <p>The platform is under active development. Some complex data table interactions within the authenticated application may not yet be fully optimized for all screen reader workflows. We continue to improve this.</p>
          </LegalSection>
          <LegalSection title="Feedback">
            <p>If you encounter an accessibility barrier, please tell us through the <Link to="/contact" style={{ color: 'var(--fs-accent)' }}>Contact page</Link>. Accessibility feedback is taken seriously.</p>
          </LegalSection>
        </div>
      </section>
    </div>
  )
}

export function NotFoundPage() {
  return (
    <div style={{ padding: '120px 0', textAlign: 'center' }}>
      <div className="pub-container">
        <div style={{ fontFamily: 'IBM Plex Mono, monospace', fontSize: 11, color: 'var(--fs-text-muted)', letterSpacing: '0.08em', marginBottom: 24 }}>
          404 — PAGE NOT FOUND
        </div>
        <h1 className="font-display" style={{ fontSize: 'clamp(36px, 5vw, 64px)', color: 'var(--fs-text)', marginBottom: 24 }}>
          This page doesn't exist.
        </h1>
        <p style={{ fontSize: 17, color: 'var(--fs-text-secondary)', marginBottom: 48, lineHeight: 1.7 }}>
          The address you entered doesn't match any page on FinSight.
        </p>
        <div style={{ display: 'flex', gap: 12, justifyContent: 'center', flexWrap: 'wrap' }}>
          <Link to="/" className="btn btn-primary">Go to Homepage</Link>
          <Link to="/login" className="btn btn-secondary">Sign In</Link>
        </div>
      </div>
    </div>
  )
}
