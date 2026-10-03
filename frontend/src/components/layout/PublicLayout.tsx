import { useState, useEffect } from 'react'
import { Link, useLocation, Outlet } from 'react-router-dom'
import { Menu, X, ChevronDown } from 'lucide-react'

const NAV_LINKS = [
  { label: 'Home', href: '/' },
  {
    label: 'Solutions', href: '/solutions',
    children: [
      { label: 'Working Capital', href: '/solutions/working-capital' },
      { label: 'Cash Flow', href: '/solutions/cash-flow' },
      { label: 'FX Exposure', href: '/solutions/fx' },
      { label: 'Risk Management', href: '/solutions/risk' },
    ]
  },
  { label: 'Platform', href: '/platform' },
  { label: 'Insights', href: '/insights' },
  { label: 'Security', href: '/security' },
  { label: 'About', href: '/about' },
  { label: 'Contact', href: '/contact' },
]

export default function PublicLayout() {
  const location = useLocation()
  const [mobileOpen, setMobileOpen] = useState(false)
  const [openDropdown, setOpenDropdown] = useState<string | null>(null)
  const [scrolled, setScrolled] = useState(false)

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 20)
    window.addEventListener('scroll', onScroll, { passive: true })
    return () => window.removeEventListener('scroll', onScroll)
  }, [])

  useEffect(() => {
    setMobileOpen(false)
    setOpenDropdown(null)
  }, [location.pathname])

  return (
    <div style={{ minHeight: '100vh', background: 'var(--fs-bg)' }}>
      {/* ── Public Nav ── */}
      <header
        style={{
          position: 'sticky',
          top: 0,
          zIndex: 100,
          background: scrolled ? 'var(--fs-bg)' : 'transparent',
          borderBottom: scrolled ? '1px solid var(--fs-border)' : '1px solid transparent',
          transition: 'background 0.2s ease, border-color 0.2s ease',
        }}
        role="banner"
      >
        <div className="pub-container" style={{ display: 'flex', alignItems: 'center', gap: 0, height: 64 }}>
          {/* Logo */}
          <Link to="/" className="pub-logo" aria-label="FinSight home" style={{ display: 'flex', alignItems: 'center', gap: 10, textDecoration: 'none', marginRight: 40 }}>
            <span style={{
              width: 32, height: 32,
              background: 'var(--fs-text)',
              color: 'var(--fs-bg)',
              display: 'flex', alignItems: 'center', justifyContent: 'center',
              fontSize: 14, fontWeight: 800,
              flexShrink: 0,
            }}>FS</span>
            <span style={{ fontSize: 17, fontWeight: 800, color: 'var(--fs-text)', letterSpacing: '-0.01em' }}>FinSight</span>
          </Link>

          {/* Desktop Nav */}
          <nav aria-label="Primary navigation" style={{ display: 'flex', alignItems: 'center', gap: 2, flex: 1 }} className="pub-nav-desktop">
            {NAV_LINKS.map((link) =>
              link.children ? (
                <div key={link.label} style={{ position: 'relative' }}
                  onMouseEnter={() => setOpenDropdown(link.label)}
                  onMouseLeave={() => setOpenDropdown(null)}
                >
                  <button
                    aria-haspopup="true"
                    aria-expanded={openDropdown === link.label}
                    style={{
                      display: 'flex', alignItems: 'center', gap: 4,
                      padding: '6px 12px',
                      background: 'none', border: 'none', cursor: 'pointer',
                      fontSize: 14, fontWeight: 500,
                      color: location.pathname.startsWith(link.href) ? 'var(--fs-text)' : 'var(--fs-text-secondary)',
                      transition: 'color 0.15s',
                    }}
                    onFocus={() => setOpenDropdown(link.label)}
                    onBlur={() => setTimeout(() => setOpenDropdown(null), 150)}
                  >
                    {link.label}
                    <ChevronDown size={14} style={{ transition: 'transform 0.15s', transform: openDropdown === link.label ? 'rotate(180deg)' : 'none' }} />
                  </button>
                  {openDropdown === link.label && (
                    <div role="menu" style={{
                      position: 'absolute', top: '100%', left: 0,
                      background: 'var(--fs-surface-elevated)',
                      border: '1px solid var(--fs-border)',
                      padding: '8px 0',
                      minWidth: 200,
                      boxShadow: '0 8px 24px rgba(0,0,0,0.06)',
                    }}>
                      {link.children.map((child) => (
                        <Link key={child.href} to={child.href} role="menuitem"
                          style={{
                            display: 'block', padding: '8px 16px',
                            fontSize: 14, color: 'var(--fs-text-secondary)',
                            textDecoration: 'none',
                            transition: 'color 0.1s, background 0.1s',
                          }}
                          onMouseEnter={e => { (e.target as HTMLElement).style.color = 'var(--fs-text)'; (e.target as HTMLElement).style.background = 'var(--fs-surface)' }}
                          onMouseLeave={e => { (e.target as HTMLElement).style.color = 'var(--fs-text-secondary)'; (e.target as HTMLElement).style.background = 'transparent' }}
                        >
                          {child.label}
                        </Link>
                      ))}
                    </div>
                  )}
                </div>
              ) : (
                <Link key={link.href} to={link.href}
                  style={{
                    padding: '6px 12px',
                    fontSize: 14, fontWeight: 500,
                    color: location.pathname === link.href ? 'var(--fs-text)' : 'var(--fs-text-secondary)',
                    textDecoration: 'none',
                    transition: 'color 0.15s',
                    borderBottom: location.pathname === link.href ? '1px solid var(--fs-text)' : '1px solid transparent',
                  }}
                >
                  {link.label}
                </Link>
              )
            )}
          </nav>

          {/* Right side CTAs */}
          <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginLeft: 'auto' }} className="pub-nav-ctas">
            <Link to="/login"
              style={{
                padding: '6px 14px',
                fontSize: 14, fontWeight: 500,
                color: 'var(--fs-text-secondary)',
                textDecoration: 'none',
              }}
            >
              Sign In
            </Link>
            <Link to="/demo"
              className="btn btn-primary"
              style={{ fontSize: 13 }}
            >
              Request a Demo
            </Link>
          </div>

          {/* Mobile hamburger */}
          <button
            aria-label={mobileOpen ? 'Close menu' : 'Open menu'}
            aria-expanded={mobileOpen}
            onClick={() => setMobileOpen(!mobileOpen)}
            style={{
              display: 'none', background: 'none', border: 'none',
              cursor: 'pointer', color: 'var(--fs-text)', padding: 8,
              marginLeft: 8,
            }}
            className="pub-hamburger"
          >
            {mobileOpen ? <X size={20} /> : <Menu size={20} />}
          </button>
        </div>

        {/* Mobile menu */}
        {mobileOpen && (
          <div className="pub-mobile-menu" role="navigation" aria-label="Mobile navigation">
            {NAV_LINKS.map((link) => (
              <div key={link.label}>
                <Link to={link.href}
                  style={{
                    display: 'block', padding: '14px 24px',
                    fontSize: 16, fontWeight: 500,
                    color: 'var(--fs-text)',
                    textDecoration: 'none',
                    borderBottom: '1px solid var(--fs-border)',
                  }}
                >
                  {link.label}
                </Link>
                {link.children?.map((child) => (
                  <Link key={child.href} to={child.href}
                    style={{
                      display: 'block', padding: '10px 40px',
                      fontSize: 14,
                      color: 'var(--fs-text-secondary)',
                      textDecoration: 'none',
                      borderBottom: '1px solid var(--fs-border)',
                    }}
                  >
                    {child.label}
                  </Link>
                ))}
              </div>
            ))}
            <div style={{ padding: '16px 24px', display: 'flex', flexDirection: 'column', gap: 10 }}>
              <Link to="/login" className="btn btn-secondary" style={{ justifyContent: 'center' }}>Sign In</Link>
              <Link to="/demo" className="btn btn-primary" style={{ justifyContent: 'center' }}>Request a Demo</Link>
            </div>
          </div>
        )}
      </header>

      {/* Page content */}
      <main>
        <Outlet />
      </main>

      {/* ── Footer ── */}
      <footer style={{
        borderTop: '1px solid var(--fs-border)',
        background: 'var(--fs-sidebar)',
        marginTop: 80,
      }} role="contentinfo">
        <div className="pub-container" style={{ padding: '60px 0 40px' }}>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(160px, 1fr))', gap: 40, marginBottom: 60 }}>
            {/* Brand */}
            <div>
              <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 12 }}>
                <span style={{ width: 24, height: 24, background: 'var(--fs-text)', color: 'var(--fs-bg)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 10, fontWeight: 800 }}>FS</span>
                <span style={{ fontSize: 15, fontWeight: 800, color: 'var(--fs-text)' }}>FinSight</span>
              </div>
              <p style={{ fontSize: 13, color: 'var(--fs-text-muted)', lineHeight: 1.7, maxWidth: 200 }}>
                Corporate treasury intelligence. Built for finance teams that need clarity.
              </p>
            </div>

            {/* Solutions */}
            <div>
              <div style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', marginBottom: 16 }}>Solutions</div>
              {[
                { label: 'Working Capital', href: '/solutions/working-capital' },
                { label: 'Cash Flow', href: '/solutions/cash-flow' },
                { label: 'FX Exposure', href: '/solutions/fx' },
                { label: 'Risk Management', href: '/solutions/risk' },
              ].map(l => <FooterLink key={l.href} {...l} />)}
            </div>

            {/* Platform */}
            <div>
              <div style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', marginBottom: 16 }}>Platform</div>
              {[
                { label: 'Overview', href: '/platform' },
                { label: 'Cash Management', href: '/platform#cash' },
                { label: 'Receivables', href: '/platform#receivables' },
                { label: 'FX & Risk', href: '/platform#fx' },
                { label: 'Reporting', href: '/platform#reporting' },
              ].map(l => <FooterLink key={l.href} {...l} />)}
            </div>

            {/* Resources */}
            <div>
              <div style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', marginBottom: 16 }}>Resources</div>
              {[
                { label: 'Insights', href: '/insights' },
                { label: 'Security', href: '/security' },
                { label: 'Request a Demo', href: '/demo' },
              ].map(l => <FooterLink key={l.href} {...l} />)}
            </div>

            {/* Company */}
            <div>
              <div style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', marginBottom: 16 }}>Company</div>
              {[
                { label: 'About', href: '/about' },
                { label: 'Contact', href: '/contact' },
              ].map(l => <FooterLink key={l.href} {...l} />)}
            </div>

            {/* Legal */}
            <div>
              <div style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', marginBottom: 16 }}>Legal</div>
              {[
                { label: 'Privacy', href: '/privacy' },
                { label: 'Terms', href: '/terms' },
                { label: 'Cookies', href: '/cookies' },
                { label: 'Accessibility', href: '/accessibility' },
              ].map(l => <FooterLink key={l.href} {...l} />)}
            </div>
          </div>

          <div style={{ borderTop: '1px solid var(--fs-border)', paddingTop: 24, display: 'flex', flexWrap: 'wrap', alignItems: 'center', justifyContent: 'space-between', gap: 16 }}>
            <p style={{ fontSize: 12, color: 'var(--fs-text-muted)' }}>
              © {new Date().getFullYear()} FinSight. All rights reserved. This is a simulation platform. No real financial transactions are executed.
            </p>
            <p style={{ fontSize: 12, color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace' }}>
              v2.1.0
            </p>
          </div>
        </div>
      </footer>
    </div>
  )
}

function FooterLink({ label, href }: { label: string; href: string }) {
  return (
    <Link to={href} style={{
      display: 'block',
      fontSize: 13,
      color: 'var(--fs-text-secondary)',
      textDecoration: 'none',
      marginBottom: 10,
      transition: 'color 0.15s',
    }}
      onMouseEnter={e => (e.target as HTMLElement).style.color = 'var(--fs-text)'}
      onMouseLeave={e => (e.target as HTMLElement).style.color = 'var(--fs-text-secondary)'}
    >
      {label}
    </Link>
  )
}
