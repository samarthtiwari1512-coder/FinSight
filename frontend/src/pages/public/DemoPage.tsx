import { useState } from 'react'
import { Loader2, CheckCircle2 } from 'lucide-react'

type FormStatus = 'idle' | 'loading' | 'success' | 'error'

const COMPANY_SIZES = ['1–10', '11–50', '51–200', '201–500', '501–1000', '1000+']
const INTERESTS = ['Working Capital', 'Cash Flow Management', 'FX Exposure', 'Risk Management', 'Full Platform Overview']
const ROLES = ['CFO / Finance Director', 'Treasury Manager', 'Finance Controller', 'Financial Analyst', 'Finance Operations', 'Other']

function validate(data: Record<string, string>) {
  const errors: Record<string, string> = {}
  ;['name', 'company', 'email', 'role'].forEach(f => {
    if (!data[f] || !data[f].trim()) errors[f] = 'Required'
  })
  if (data.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(data.email)) {
    errors.email = 'Enter a valid work email address'
  }
  return errors
}

export default function DemoPage() {
  const [form, setForm] = useState({ name: '', company: '', email: '', role: '', companySize: '', interest: '', message: '' })
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [status, setStatus] = useState<FormStatus>('idle')

  const set = (field: string) => (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    setForm(f => ({ ...f, [field]: e.target.value }))
    setErrors(err => { const next = { ...err }; delete next[field]; return next })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const errs = validate(form)
    if (Object.keys(errs).length > 0) { setErrors(errs); return }
    setStatus('loading')
    await new Promise(r => setTimeout(r, 1400))
    setStatus('success')
  }

  return (
    <div>
      <section style={{ padding: '100px 0 60px', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 16 }}>
            Request a Demo
          </p>
          <h1 className="font-display" style={{ fontSize: 'clamp(36px, 5vw, 64px)', color: 'var(--fs-text)', lineHeight: 1.05, marginBottom: 20 }}>
            See FinSight in practice.
          </h1>
          <p style={{ fontSize: 17, color: 'var(--fs-text-secondary)', maxWidth: '52ch', lineHeight: 1.7 }}>
            Walk through the platform with realistic synthetic data. No commitment required. Typically 30–45 minutes.
          </p>
        </div>
      </section>

      <section style={{ padding: '80px 0' }}>
        <div className="pub-container">
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 2fr', gap: 80, alignItems: 'start' }} className="pub-two-col">
            {/* What to expect */}
            <div>
              <h2 style={{ fontSize: 16, fontWeight: 600, color: 'var(--fs-text)', marginBottom: 24 }}>What to expect</h2>
              <div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
                {[
                  { step: '01', label: 'Live platform walkthrough', detail: 'We\'ll walk through the modules most relevant to your situation — cash, working capital, FX or risk.' },
                  { step: '02', label: 'Realistic demo data', detail: 'All data shown is pre-seeded with realistic synthetic financial records — not generic placeholder data.' },
                  { step: '03', label: 'Your questions answered', detail: 'Bring specific scenarios or use cases. This is a working session, not a slide deck presentation.' },
                ].map(({ step, label, detail }) => (
                  <div key={step} style={{ display: 'flex', gap: 16 }}>
                    <span style={{ fontFamily: 'IBM Plex Mono, monospace', fontSize: 11, color: 'var(--fs-text-muted)', fontWeight: 700, flexShrink: 0, marginTop: 3 }}>{step}</span>
                    <div>
                      <div style={{ fontSize: 14, fontWeight: 600, color: 'var(--fs-text)', marginBottom: 6 }}>{label}</div>
                      <div style={{ fontSize: 13, color: 'var(--fs-text-secondary)', lineHeight: 1.7 }}>{detail}</div>
                    </div>
                  </div>
                ))}
              </div>

              <div style={{ marginTop: 48, padding: '20px', border: '1px solid var(--fs-border)', background: 'var(--fs-surface)' }}>
                <div style={{ fontSize: 12, fontWeight: 700, color: 'var(--fs-text-muted)', letterSpacing: '0.05em', textTransform: 'uppercase', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 12 }}>
                  Or sign in directly
                </div>
                <p style={{ fontSize: 13, color: 'var(--fs-text-secondary)', lineHeight: 1.7, marginBottom: 16 }}>
                  The demo environment is available immediately. Use one of the pre-set roles on the login page.
                </p>
                <a href="/login" className="btn btn-secondary" style={{ fontSize: 13 }}>Go to Login →</a>
              </div>
            </div>

            {/* Form */}
            {status === 'success' ? (
              <div style={{
                padding: '60px',
                border: '1px solid var(--fs-border)',
                background: 'var(--fs-surface)',
                textAlign: 'center',
              }}>
                <CheckCircle2 size={40} style={{ color: 'var(--fs-positive)', marginBottom: 16 }} />
                <h2 style={{ fontSize: 22, fontWeight: 600, color: 'var(--fs-text)', marginBottom: 12 }}>Request received</h2>
                <p style={{ fontSize: 15, color: 'var(--fs-text-secondary)', lineHeight: 1.7, marginBottom: 24 }}>
                  Thank you. We'll be in touch shortly to confirm a time that works.
                </p>
                <p style={{ fontSize: 14, color: 'var(--fs-text-muted)' }}>
                  In the meantime, you can explore the demo environment directly from the login page.
                </p>
              </div>
            ) : (
              <form onSubmit={handleSubmit} noValidate>
                <div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 16 }} className="pub-two-col">
                    <DemoField label="Full Name" id="demo-name" error={errors.name} required>
                      <input id="demo-name" className="input" value={form.name} onChange={set('name')} placeholder="Your name" style={{ height: 44 }} />
                    </DemoField>
                    <DemoField label="Company" id="demo-company" error={errors.company} required>
                      <input id="demo-company" className="input" value={form.company} onChange={set('company')} placeholder="Company name" style={{ height: 44 }} />
                    </DemoField>
                  </div>
                  <DemoField label="Work Email" id="demo-email" error={errors.email} required>
                    <input id="demo-email" type="email" className="input" value={form.email} onChange={set('email')} placeholder="you@company.com" style={{ height: 44 }} />
                  </DemoField>
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 16 }} className="pub-two-col">
                    <DemoField label="Role" id="demo-role" error={errors.role} required>
                      <select id="demo-role" className="input" value={form.role} onChange={set('role')} style={{ height: 44, cursor: 'pointer' }}>
                        <option value="">Select role</option>
                        {ROLES.map(r => <option key={r} value={r}>{r}</option>)}
                      </select>
                    </DemoField>
                    <DemoField label="Company Size" id="demo-size" error={errors.companySize}>
                      <select id="demo-size" className="input" value={form.companySize} onChange={set('companySize')} style={{ height: 44, cursor: 'pointer' }}>
                        <option value="">Select size</option>
                        {COMPANY_SIZES.map(s => <option key={s} value={s}>{s} employees</option>)}
                      </select>
                    </DemoField>
                  </div>
                  <DemoField label="Primary Interest" id="demo-interest" error={errors.interest}>
                    <select id="demo-interest" className="input" value={form.interest} onChange={set('interest')} style={{ height: 44, cursor: 'pointer' }}>
                      <option value="">What would you like to focus on?</option>
                      {INTERESTS.map(i => <option key={i} value={i}>{i}</option>)}
                    </select>
                  </DemoField>
                  <DemoField label="Anything else?" id="demo-message" error={errors.message}>
                    <textarea id="demo-message" className="input" value={form.message} onChange={set('message')} placeholder="Tell us about your current setup or what you're trying to solve (optional)" rows={4} style={{ resize: 'vertical' }} />
                  </DemoField>

                  {status === 'error' && (
                    <div style={{ padding: '12px 16px', background: 'var(--fs-negative-bg)', border: '1px solid var(--fs-negative)', color: 'var(--fs-negative)', fontSize: 14 }}>
                      Something went wrong. Please try again.
                    </div>
                  )}

                  <div>
                    <button type="submit" className="btn btn-primary btn-lg" disabled={status === 'loading'} style={{ minWidth: 160 }}>
                      {status === 'loading'
                        ? <><Loader2 size={15} style={{ animation: 'spin 0.6s linear infinite' }} /> Submitting...</>
                        : 'Request a Demo'
                      }
                    </button>
                  </div>
                </div>
              </form>
            )}
          </div>
        </div>
      </section>
    </div>
  )
}

function DemoField({ label, id, error, required, children }: {
  label: string; id: string; error?: string; required?: boolean; children: React.ReactNode
}) {
  return (
    <div>
      <label htmlFor={id} style={{
        display: 'block', fontSize: 12, fontWeight: 600,
        color: 'var(--fs-text-secondary)',
        letterSpacing: '0.04em', textTransform: 'uppercase',
        marginBottom: 8,
      }}>
        {label}{required && <span style={{ color: 'var(--fs-negative)', marginLeft: 4 }}>*</span>}
      </label>
      {children}
      {error && <p role="alert" style={{ fontSize: 12, color: 'var(--fs-negative)', marginTop: 5 }}>{error}</p>}
    </div>
  )
}
