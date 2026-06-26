import { useState } from 'react';
import Layout from './components/Layout';
import MetricsGrid from './components/MetricsGrid';
import PullRequestTable from './components/PullRequestTable';
import FindingDetails from './components/FindingDetails';
import RulesManager from './components/RulesManager';
import { useAuditStream } from './hooks/useAuditStream';
import TrendCharts from './components/TrendCharts';
import RepositoriesPage from './components/RepositoriesPage';
import ChatPanel from './components/ChatPanel';
import { MessageSquare } from 'lucide-react';
import type { PullRequest } from './types/api';

type Page = 'dashboard' | 'repositories' | 'settings';

function App() {
  const { metrics, pullRequests, loading, lastUpdated, refresh } = useAuditStream();
  const [selectedPR, setSelectedPR] = useState<PullRequest | null>(null);
  const [page, setPage] = useState<Page>('dashboard');
  const [chatOpen, setChatOpen] = useState(false);

  function handleFindingDismissed() {
    refresh();
  }

  return (
    <Layout currentPage={page} onNavigate={(p) => { setSelectedPR(null); setPage(p); }}>
      <div
        key={page}
        style={{ animation: 'fadeIn 200ms ease' }}
      >
        {page === 'dashboard' && (
          <>
            <div style={styles.pageHeader}>
              <h1 style={styles.pageTitle}>SYSTEM OVERVIEW (GLOBAL POSTURE)</h1>
              <div style={styles.lastUpdated}>
                {lastUpdated && (
                  <span>Last updated: {lastUpdated.toLocaleTimeString()}</span>
                )}
                <button style={styles.refreshBtn} onClick={refresh}>
                  ↻ Refresh
                </button>
              </div>
            </div>

            <MetricsGrid metrics={metrics} loading={loading} />
            <TrendCharts />
            <PullRequestTable
              pullRequests={pullRequests}
              loading={loading}
              selectedPRId={selectedPR?.id ?? null}
              onSelectPR={setSelectedPR}
            />
            <FindingDetails
              selectedPR={selectedPR}
              onFindingDismissed={handleFindingDismissed}
            />
          </>
        )}

        {page === 'repositories' && <RepositoriesPage />}

        {page === 'settings' && (
          <>
            <div style={styles.pageHeader}>
              <h1 style={styles.pageTitle}>SETTINGS</h1>
            </div>
            <RulesManager />
          </>
        )}
      </div>
      {/* Floating chat button */}
      <button
        style={{
          position: 'fixed',
          bottom: 24,
          right: 24,
          width: 52,
          height: 52,
          borderRadius: '50%',
          backgroundColor: '#1d4ed8',
          border: 'none',
          color: '#fff',
          cursor: 'pointer',
          display: chatOpen ? 'none' : 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          boxShadow: '0 4px 24px rgba(29,78,216,0.4)',
          zIndex: 999,
        }}
        onClick={() => setChatOpen(true)}
        title="Ask OpsMind AI"
      >
        <MessageSquare size={22} />
      </button>

      <ChatPanel isOpen={chatOpen} onClose={() => setChatOpen(false)} />
    </Layout>
  );
}

const styles: { [key: string]: React.CSSProperties } = {
  pageHeader: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    marginBottom: 20,
  },
  pageTitle: {
    fontSize: 16,
    fontWeight: 700,
    color: '#9ca3af',
    letterSpacing: 0.5,
  },
  lastUpdated: {
    display: 'flex',
    alignItems: 'center',
    gap: 12,
    fontSize: 12,
    color: '#6b7280',
  },
  refreshBtn: {
    backgroundColor: 'transparent',
    border: '1px solid #374151',
    color: '#9ca3af',
    borderRadius: 6,
    padding: '4px 12px',
    fontSize: 12,
    cursor: 'pointer',
  },
};

export default App;