import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'
import { Eye, EyeOff, Loader2 } from 'lucide-react'
import toast from 'react-hot-toast'

const DEMO_CREDENTIALS = [
  { role: 'CFO', email: 'amelia.chen@acmeglobal.com', password: 'Demo@2024' },
  { role: 'Treasury', email: 'treasury@acmeglobal.com', password: 'Demo@2024' },
  { role: 'Finance', email: 'finance@acmeglobal.com', password: 'Demo@2024' },
  { role: 'Auditor', email: 'auditor@acmeglobal.com', password: 'Demo@2024' },
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
      background: 'var(--fs-bg)',
    }}>
      {/* Left Panel — Branding (Editorial / Financial Ledger Style) */}
      <div style={{
        flex: 1,
        display: 'flex',
        flexDirection: 'column',
        padding: '60px',
        justifyContent: 'space-between',
        borderRight: '1px solid var(--fs-border)',
        background: 'var(--fs-sidebar)',
      }}>
        {/* Logo */}
        <div style={{ display: 'flex', alignItems: 'center', gap: 14 }}>
          <div style={{
            width: 44,
            height: 44,
            background: 'var(--fs-text)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontSize: 18,
            fontWeight: 800,
            color: 'var(--fs-bg)',
          }}>
            FS
          </div>
          <div>
            <div style={{ fontSize: 22, fontWeight: 800, color: 'var(--fs-text)', letterSpacing: '-0.01em' }}>
              FinSight
            </div>
            <div style={{ fontSize: 11, color: 'var(--fs-text-muted)', letterSpacing: '0.08em' }}>
              TREASURY INTELLIGENCE
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
            border: '1px solid var(--fs-border-strong)',
            marginBottom: 24,
            fontSize: 11, 
            fontWeight: 600, 
            color: 'var(--fs-text-secondary)', 
            letterSpacing: '0.05em',
            textTransform: 'uppercase'
          }}>
            ACME GLOBAL • DEMO ENVIRONMENT
          </div>

          <h1 className="font-display" style={{
            fontSize: 48,
            color: 'var(--fs-text)',
            lineHeight: 1.1,
            marginBottom: 20,
          }}>
            Corporate Treasury. <br />
            <span style={{ color: 'var(--fs-text-secondary)' }}>Quantified.</span>
          </h1>

          <p style={{ fontSize: 15, color: 'var(--fs-text-secondary)', lineHeight: 1.7, marginBottom: 40, borderLeft: '2px solid var(--fs-accent)', paddingLeft: 16 }}>
            Monitor cash positions, FX exposure, working capital, and financial risk
            across all entities in real time — with the precision of an enterprise treasury workstation.
          </p>
        </div>

        <div style={{ fontSize: 11, color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace' }}>
          VER 2.1.0-SIMULATION
        </div>
      </div>

      {/* Right Panel — Login Form */}
      <div style={{
        width: 520,
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        padding: '60px',
        background: 'var(--fs-bg)',
      }}>
        <div className="fade-in">
          <div style={{ marginBottom: 36, borderBottom: '1px solid var(--fs-border)', paddingBottom: 24 }}>
            <h2 className="font-display" style={{ fontSize: 28, color: 'var(--fs-text)', marginBottom: 8 }}>
              Sign In
            </h2>
            <p style={{ fontSize: 14, color: 'var(--fs-text-secondary)' }}>
              Enter your credentials to access the treasury terminal
            </p>
          </div>

          <form onSubmit={handleLogin} style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
            <div>
              <label style={{ display: 'block', fontSize: 12, fontWeight: 600, color: 'var(--fs-text-secondary)', marginBottom: 6, textTransform: 'uppercase', letterSpacing: '0.05em' }}>
                Email Address
              </label>
              <input
                id="email"
                type="email"
                className="input"
                placeholder="you@company.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                autoComplete="username"
                autoFocus
                style={{ height: 44 }}
              />
            </div>

            <div>
              <label style={{ display: 'block', fontSize: 12, fontWeight: 600, color: 'var(--fs-text-secondary)', marginBottom: 6, textTransform: 'uppercase', letterSpacing: '0.05em' }}>
                Password
              </label>
              <div style={{ position: 'relative' }}>
                <input
                  id="password"
                  type={showPassword ? 'text' : 'password'}
                  className="input"
                  placeholder="••••••••"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  autoComplete="current-password"
                  style={{ paddingRight: 40, height: 44 }}
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  style={{
                    position: 'absolute',
                    right: 12,
                    top: '50%',
                    transform: 'translateY(-50%)',
                    background: 'none',
                    border: 'none',
                    cursor: 'pointer',
                    color: 'var(--fs-text-muted)',
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
              className="btn btn-primary"
              disabled={isLoading}
              style={{ marginTop: 8, height: 44, width: '100%' }}
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
          <div style={{ marginTop: 48, paddingTop: 24, borderTop: '1px solid var(--fs-border)' }}>
            <div style={{ fontSize: 11, color: 'var(--fs-text-muted)', fontWeight: 600, letterSpacing: '0.05em', marginBottom: 16 }}>
              QUICK DEMO ACCESS
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12 }}>
              {DEMO_CREDENTIALS.map((cred) => (
                <button
                  key={cred.role}
                  onClick={() => fillDemo(cred.email, cred.password)}
                  className="btn btn-secondary"
                  style={{ justifyContent: 'center' }}
                >
                  {cred.role}
                </button>
              ))}
            </div>

            <p style={{
              fontSize: 12,
              color: 'var(--fs-text-muted)',
              marginTop: 16,
              lineHeight: 1.5,
            }}>
              Select a role to pre-fill credentials. All accounts are provisioned with simulated financial records for demonstration purposes.
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
