import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'
import { Eye, EyeOff, Shield, TrendingUp, Globe, BarChart3, Loader2 } from 'lucide-react'
import toast from 'react-hot-toast'

const DEMO_CREDENTIALS = [
  { role: 'CFO', email: 'amelia.chen@acmeglobal.com', password: 'Demo@2024' },
  { role: 'Treasury', email: 'treasury@acmeglobal.com', password: 'Demo@2024' },
  { role: 'Finance', email: 'finance@acmeglobal.com', password: 'Demo@2024' },
  { role: 'Auditor', email: 'auditor@acmeglobal.com', password: 'Demo@2024' },
]

const FEATURES = [
  { icon: TrendingUp, label: 'Real-time Cash Position', description: '24 bank accounts across 6 entities' },
  { icon: Globe, label: 'FX Exposure & Hedging', description: '$45M+ exposure across 8 currencies' },
  { icon: BarChart3, label: 'Working Capital KPIs', description: 'DSO, DPO, CCC with trend analysis' },
  { icon: Shield, label: 'Risk & Compliance', description: 'Automated alerts with audit trails' },
]

export default function LoginPage() {
  const navigate = useNavigate()
  const { login } = useAuthStore()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!email || !password) {
      toast.error('Please enter your credentials')
      return
    }

    setIsLoading(true)
    try {
      await login(email, password)
      navigate('/dashboard')
      toast.success('Welcome back!')
    } catch (err: any) {
      const msg = err?.response?.data?.error?.message || 'Invalid email or password'
      toast.error(msg)
    } finally {
      setIsLoading(false)
    }
  }

  const fillDemo = (email: string, password: string) => {
    setEmail(email)
    setPassword(password)
  }

  return (
    <div style={{
      minHeight: '100vh',
      display: 'flex',
      background: 'var(--bg-primary)',
      position: 'relative',
      overflow: 'hidden',
    }}>
      {/* Background glow effects */}
      <div style={{
        position: 'absolute',
        top: '-20%',
        left: '-10%',
        width: 600,
        height: 600,
        background: 'radial-gradient(circle, rgba(37,99,235,0.08) 0%, transparent 70%)',
        pointerEvents: 'none',
      }} />
      <div style={{
        position: 'absolute',
        bottom: '-20%',
        right: '-10%',
        width: 500,
        height: 500,
        background: 'radial-gradient(circle, rgba(124,58,237,0.06) 0%, transparent 70%)',
        pointerEvents: 'none',
      }} />

      {/* Left Panel — Branding */}
      <div style={{
        flex: 1,
        display: 'flex',
        flexDirection: 'column',
        padding: '60px',
        justifyContent: 'space-between',
        borderRight: '1px solid var(--border-color)',
        background: 'linear-gradient(135deg, rgba(15,22,41,0.8) 0%, rgba(10,15,30,0.9) 100%)',
        backdropFilter: 'blur(10px)',
        position: 'relative',
        zIndex: 1,
      }}>
        {/* Logo */}
        <div style={{ display: 'flex', alignItems: 'center', gap: 14 }}>
          <div style={{
            width: 44,
            height: 44,
            background: 'linear-gradient(135deg, #1d4ed8, #60a5fa)',
            borderRadius: 12,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontSize: 18,
            fontWeight: 800,
            color: 'white',
            boxShadow: '0 4px 20px rgba(37, 99, 235, 0.3)',
          }}>FS</div>
          <div>
            <div style={{ fontSize: 22, fontWeight: 800, color: 'white', letterSpacing: '-0.01em' }}>
              FinSight
            </div>
            <div style={{ fontSize: 11, color: 'rgba(148,163,184,0.8)', letterSpacing: '0.08em' }}>
              TREASURY INTELLIGENCE PLATFORM
            </div>
          </div>
        </div>

        {/* Hero Content */}
        <div style={{ maxWidth: 480 }}>
          <div style={{
            display: 'inline-flex',
            alignItems: 'center',
            gap: 8,
            padding: '6px 14px',
            borderRadius: 100,
            background: 'rgba(37, 99, 235, 0.1)',
            border: '1px solid rgba(37, 99, 235, 0.2)',
            marginBottom: 24,
          }}>
            <div style={{ width: 6, height: 6, borderRadius: '50%', background: '#10b981' }} />
            <span style={{ fontSize: 11, fontWeight: 600, color: '#60a5fa', letterSpacing: '0.05em' }}>
              ACME GLOBAL MANUFACTURING LTD • DEMO
            </span>
          </div>

          <h1 style={{
            fontSize: 38,
            fontWeight: 800,
            color: 'white',
            lineHeight: 1.15,
            letterSpacing: '-0.02em',
            marginBottom: 20,
          }}>
            Corporate Treasury<br />
            <span style={{
              background: 'linear-gradient(90deg, #60a5fa, #a78bfa)',
              WebkitBackgroundClip: 'text',
              WebkitTextFillColor: 'transparent',
            }}>at Your Fingertips</span>
          </h1>

          <p style={{ fontSize: 15, color: 'rgba(148,163,184,0.9)', lineHeight: 1.7, marginBottom: 40 }}>
            Monitor cash positions, FX exposure, working capital, and financial risk
            across all entities in real time — with the precision of an enterprise treasury system.
          </p>

          {/* Feature highlights */}
          <div style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
            {FEATURES.map((f) => (
              <div key={f.label} style={{ display: 'flex', alignItems: 'center', gap: 14 }}>
                <div style={{
                  width: 36,
                  height: 36,
                  background: 'rgba(37, 99, 235, 0.1)',
                  border: '1px solid rgba(37, 99, 235, 0.15)',
                  borderRadius: 8,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  flexShrink: 0,
                }}>
                  <f.icon size={16} style={{ color: '#60a5fa' }} />
                </div>
                <div>
                  <div style={{ fontSize: 13, fontWeight: 600, color: 'rgba(241,245,249,0.9)' }}>{f.label}</div>
                  <div style={{ fontSize: 11, color: 'rgba(148,163,184,0.7)' }}>{f.description}</div>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div style={{ fontSize: 11, color: 'rgba(100,116,139,0.6)' }}>
          FinSight is a simulation platform. No real financial transactions are executed.
        </div>
      </div>

      {/* Right Panel — Login Form */}
      <div style={{
        width: 460,
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        padding: '60px 48px',
        position: 'relative',
        zIndex: 1,
      }}>
        <div className="fade-in">
          <div style={{ marginBottom: 36 }}>
            <h2 style={{ fontSize: 26, fontWeight: 700, color: 'var(--text-primary)', marginBottom: 8 }}>
              Sign in to FinSight
            </h2>
            <p style={{ fontSize: 14, color: 'var(--text-secondary)' }}>
              Use your credentials or select a demo role below
            </p>
          </div>

          <form onSubmit={handleLogin} style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
            <div className="form-group">
              <label className="form-label">Email Address</label>
              <input
                id="email"
                type="email"
                className="input"
                placeholder="you@company.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                autoComplete="username"
                autoFocus
              />
            </div>

            <div className="form-group">
              <label className="form-label">Password</label>
              <div style={{ position: 'relative' }}>
                <input
                  id="password"
                  type={showPassword ? 'text' : 'password'}
                  className="input"
                  placeholder="••••••••"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  autoComplete="current-password"
                  style={{ paddingRight: 40 }}
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  style={{
                    position: 'absolute',
                    right: 10,
                    top: '50%',
                    transform: 'translateY(-50%)',
                    background: 'none',
                    border: 'none',
                    cursor: 'pointer',
                    color: 'var(--text-muted)',
                    display: 'flex',
                  }}
                >
                  {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                </button>
              </div>
            </div>

            <button
              id="login-btn"
              type="submit"
              className="btn btn-primary btn-lg"
              disabled={isLoading}
              style={{ marginTop: 4 }}
            >
              {isLoading ? (
                <>
                  <Loader2 size={16} style={{ animation: 'spin 0.6s linear infinite' }} />
                  Signing in...
                </>
              ) : (
                'Sign in'
              )}
            </button>
          </form>

          {/* Demo Credentials */}
          <div style={{ marginTop: 36 }}>
            <div style={{
              display: 'flex',
              alignItems: 'center',
              gap: 12,
              marginBottom: 16,
            }}>
              <div style={{ flex: 1, height: 1, background: 'var(--border-color)' }} />
              <span style={{ fontSize: 11, color: 'var(--text-muted)', fontWeight: 600, whiteSpace: 'nowrap' }}>
                QUICK DEMO ACCESS
              </span>
              <div style={{ flex: 1, height: 1, background: 'var(--border-color)' }} />
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8 }}>
              {DEMO_CREDENTIALS.map((cred) => (
                <button
                  key={cred.role}
                  onClick={() => fillDemo(cred.email, cred.password)}
                  className="btn btn-secondary btn-sm"
                  style={{ justifyContent: 'center' }}
                >
                  {cred.role}
                </button>
              ))}
            </div>

            <p style={{
              fontSize: 11,
              color: 'var(--text-muted)',
              marginTop: 12,
              textAlign: 'center',
              lineHeight: 1.5,
            }}>
              Click a role to fill credentials, then click Sign In.
              <br />All demo data is pre-seeded with realistic financial records.
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
