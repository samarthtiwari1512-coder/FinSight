/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: ['selector', '[data-theme="dark"]'],
  theme: {
    extend: {
      fontFamily: {
        sans: ['"IBM Plex Sans"', 'sans-serif'],
        mono: ['"IBM Plex Mono"', 'monospace'],
        display: ['"DM Serif Display"', 'serif'],
      },
      colors: {
        fs: {
          bg: 'var(--fs-bg)',
          sidebar: 'var(--fs-sidebar)',
          surface: 'var(--fs-surface)',
          'surface-elevated': 'var(--fs-surface-elevated)',
          'surface-secondary': 'var(--fs-surface-secondary)',
          text: 'var(--fs-text)',
          'text-secondary': 'var(--fs-text-secondary)',
          'text-muted': 'var(--fs-text-muted)',
          border: 'var(--fs-border)',
          'border-strong': 'var(--fs-border-strong)',
          accent: 'var(--fs-accent)',
          positive: 'var(--fs-positive)',
          'positive-bg': 'var(--fs-positive-bg)',
          warning: 'var(--fs-warning)',
          'warning-bg': 'var(--fs-warning-bg)',
          negative: 'var(--fs-negative)',
          'negative-bg': 'var(--fs-negative-bg)',
          neutral: 'var(--fs-neutral)',
        }
      },
      boxShadow: {
        'fs-subtle': '0 1px 2px rgba(25, 26, 23, 0.05)',
      },
      transitionDuration: {
        'fs': '200ms',
      },
      borderRadius: {
        'fs': '8px',
      }
    },
  },
  plugins: [],
}
