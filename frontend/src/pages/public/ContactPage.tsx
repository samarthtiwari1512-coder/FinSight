import { useState } from 'react'
import { Loader2, CheckCircle2 } from 'lucide-react'

type FormStatus = 'idle' | 'loading' | 'success' | 'error'

function validate(data: Record<string, string>, fields: string[]) {
  const errors: Record<string, string> = {}
  fields.forEach(f => {
    if (!data[f] || !data[f].trim()) errors[f] = 'This field is required'
  })
  if (data.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(data.email)) {
    errors.email = 'Enter a valid work email address'
  }
  return errors
}

export default function ContactPage() {
  const [form, setForm] = useState({ name: '', company: '', email: '', subject: '', message: '' })
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [status, setStatus] = useState<FormStatus>('idle')

  const set = (field: string) => (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    setForm(f => ({ ...f, [field]: e.target.value }))
    setErrors(err => { const next = { ...err }; delete next[field]; return next })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const errs = validate(form, ['name', 'email', 'subject', 'message'])
    if (Object.keys(errs).length > 0) { setErrors(errs); return }
    setStatus('loading')
    // Simulate submit (real backend integration would go here)
    await new Promise(r => setTimeout(r, 1200))
    setStatus('success')
  }

  return (
    <div>
      <section style={{ padding: '100px 0 60px', borderBottom: '1px solid var(--fs-border)' }}>
        <div className="pub-container">
          <p style={{ fontSize: 11, fontWeight: 700, letterSpacing: '0.06em', textTransform: 'uppercase', color: 'var(--fs-text-muted)', fontFamily: 'IBM Plex Mono, monospace', marginBottom: 16 }}>
            Contact
          </p>
          <h1 className="font-display" style={{ fontSize: 'clamp(36px, 5vw, 64px)', color: 'var(--fs-text)', lineHeight: 1.05, marginBottom: 16 }}>
            Get in touch.
          </h1>
          <p style={{ fontSize: 17, color: 'var(--fs-text-secondary)', maxWidth: '48ch', lineHeight: 1.7 }}>
            Questions about FinSight, the platform, security, or requesting a demo — use the form below.
          </p>
        </div>
      </section>

      <section style={{ padding: '80px 0' }}>
        <div className="pub-container">
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 2fr', gap: 80, alignItems: 'start' }} className="pub-two-col">
            {/* Contact details */}
            <div>
              <h2 style={{ fontSize: 16, fontWeight: 600, color: 'var(--fs-text)', marginBottom: 24 }}>Before you write</h2>
              <div style={{ display: 'flex', flexDirection: 'column', gap: 24 }}>
                {[
                  { label: 'General enquiries', detail: 'Questions about FinSight, how it works, or whether it\'s right for your team.' },
                  { label: 'Security disclosure', detail: 'Responsible disclosure of security issues. Please do not post details publicly.' },
                  { label: 'Request a demo', detail: 'For a walkthrough of the platform, use the Request a Demo form for a more structured response.' },
                ].map(({ label, detail }) => (
                  <div key={label}>
                    <div style={{ fontSize: 13, fontWeight: 600, color: 'var(--fs-text)', marginBottom: 6 }}>{label}</div>
                    <div style={{ fontSize: 13, color: 'var(--fs-text-secondary)', lineHeight: 1.7 }}>{detail}</div>
                  </div>
                ))}
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
                <h2 style={{ fontSize: 22, fontWeight: 600, color: 'var(--fs-text)', marginBottom: 12 }}>Message sent</h2>
                <p style={{ fontSize: 15, color: 'var(--fs-text-secondary)', lineHeight: 1.7 }}>
                  We've received your message. We'll respond as soon as possible.
                </p>
              </div>
            ) : (
              <form onSubmit={handleSubmit} noValidate>
                <div style={{ display: 'flex', flexDirection: 'column', gap: 24 }}>
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 20 }} className="pub-two-col">
                    <FormField label="Name" id="contact-name" error={errors.name}>
                      <input id="contact-name" className="input" value={form.name} onChange={set('name')} placeholder="Your name" style={{ height: 44 }} />
                    </FormField>
                    <FormField label="Company" id="contact-company" error={errors.company}>
                      <input id="contact-company" className="input" value={form.company} onChange={set('company')} placeholder="Company name" style={{ height: 44 }} />
                    </FormField>
                  </div>
                  <FormField label="Work Email" id="contact-email" error={errors.email} required>
                    <input id="contact-email" type="email" className="input" value={form.email} onChange={set('email')} placeholder="you@company.com" style={{ height: 44 }} />
                  </FormField>
                  <FormField label="Subject" id="contact-subject" error={errors.subject} required>
                    <input id="contact-subject" className="input" value={form.subject} onChange={set('subject')} placeholder="What is your message about?" style={{ height: 44 }} />
                  </FormField>
                  <FormField label="Message" id="contact-message" error={errors.message} required>
                    <textarea id="contact-message" className="input" value={form.message} onChange={set('message')} placeholder="Your message..." rows={6} style={{ resize: 'vertical' }} />
                  </FormField>

                  {status === 'error' && (
                    <div style={{ padding: '12px 16px', background: 'var(--fs-negative-bg)', border: '1px solid var(--fs-negative)', color: 'var(--fs-negative)', fontSize: 14 }}>
                      Something went wrong. Please try again.
                    </div>
                  )}

                  <div>
                    <button type="submit" className="btn btn-primary" disabled={status === 'loading'} style={{ height: 44, minWidth: 140 }}>
                      {status === 'loading' ? <><Loader2 size={15} style={{ animation: 'spin 0.6s linear infinite' }} /> Sending...</> : 'Send Message'}
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

function FormField({ label, id, error, required, children }: {
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
      {error && <p role="alert" style={{ fontSize: 12, color: 'var(--fs-negative)', marginTop: 6 }}>{error}</p>}
    </div>
  )
}
