import type { ReactNode } from 'react';
import { LayoutDashboard, GitBranch, Settings, Bot } from 'lucide-react';

interface LayoutProps {
  children: ReactNode;
}

export default function Layout({ children }: LayoutProps) {
  return (
    <div style={styles.page}>
      <header style={styles.header}>
        <div style={styles.brand}>
          <Bot size={22} color="#60a5fa" />
          <span style={styles.brandText}>OpsMind AI</span>
        </div>

        <nav style={styles.nav}>
          <a href="#" style={styles.navItemActive}>
            <LayoutDashboard size={16} />
            Dashboard
          </a>
          <a href="#" style={styles.navItem}>
            <GitBranch size={16} />
            Repositories
          </a>
          <a href="#" style={styles.navItem}>
            <Settings size={16} />
            Settings
          </a>
        </nav>

        <div style={styles.userBadge}>User: Admin</div>
      </header>

      <main style={styles.main}>{children}</main>
    </div>
  );
}

const styles: { [key: string]: React.CSSProperties } = {
  page: {
    minHeight: '100vh',
    backgroundColor: '#0f1117',
    color: '#e5e7eb',
    fontFamily: "'Inter', system-ui, sans-serif",
  },
  header: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    padding: '14px 24px',
    borderBottom: '1px solid #1f2330',
    backgroundColor: '#13151f',
  },
  brand: {
    display: 'flex',
    alignItems: 'center',
    gap: 8,
  },
  brandText: {
    fontWeight: 700,
    fontSize: 16,
    letterSpacing: 0.3,
  },
  nav: {
    display: 'flex',
    gap: 24,
  },
  navItem: {
    display: 'flex',
    alignItems: 'center',
    gap: 6,
    color: '#9ca3af',
    textDecoration: 'none',
    fontSize: 14,
  },
  navItemActive: {
    display: 'flex',
    alignItems: 'center',
    gap: 6,
    color: '#60a5fa',
    textDecoration: 'none',
    fontSize: 14,
    fontWeight: 600,
  },
  userBadge: {
    fontSize: 13,
    color: '#9ca3af',
  },
  main: {
    padding: '24px',
    maxWidth: 1400,
    margin: '0 auto',
  },
};