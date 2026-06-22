import type { ReactNode } from 'react';
import { LayoutDashboard, GitBranch, Settings, Bot } from 'lucide-react';

type Page = 'dashboard' | 'settings';

interface LayoutProps {
  children: ReactNode;
  currentPage: Page;
  onNavigate: (page: Page) => void;
}

export default function Layout({ children, currentPage, onNavigate }: LayoutProps) {
  return (
    <div style={styles.page}>
      <header style={styles.header}>
        <div style={styles.brand}>
          <Bot size={22} color="#60a5fa" />
          <span style={styles.brandText}>OpsMind AI</span>
        </div>

        <nav style={styles.nav}>
          <button
            style={currentPage === 'dashboard' ? styles.navItemActive : styles.navItem}
            onClick={() => onNavigate('dashboard')}
          >
            <LayoutDashboard size={16} />
            Dashboard
          </button>
          <button
            style={styles.navItem}
            onClick={() => {}}
          >
            <GitBranch size={16} />
            Repositories
          </button>
          <button
            style={currentPage === 'settings' ? styles.navItemActive : styles.navItem}
            onClick={() => onNavigate('settings')}
          >
            <Settings size={16} />
            Settings
          </button>
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
    background: 'none',
    border: 'none',
    fontSize: 14,
    cursor: 'pointer',
    padding: 0,
  },
  navItemActive: {
    display: 'flex',
    alignItems: 'center',
    gap: 6,
    color: '#60a5fa',
    background: 'none',
    border: 'none',
    fontSize: 14,
    fontWeight: 600,
    cursor: 'pointer',
    padding: 0,
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