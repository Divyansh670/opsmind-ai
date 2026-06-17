import { useEffect, useState } from 'react';
import Layout from './components/Layout';
import MetricsGrid from './components/MetricsGrid';
import PullRequestTable from './components/PullRequestTable';
import FindingDetails from './components/FindingDetails';
import { apiClient } from './api/client';
import type { DashboardMetrics, PullRequest } from './types/api';

function App() {
  const [metrics, setMetrics] = useState<DashboardMetrics | null>(null);
  const [pullRequests, setPullRequests] = useState<PullRequest[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedPR, setSelectedPR] = useState<PullRequest | null>(null);

  useEffect(() => {
    loadDashboard();
  }, []);

  async function loadDashboard() {
    setLoading(true);
    try {
      const [metricsRes, prsRes] = await Promise.all([
        apiClient.get('/api/metrics'),
        apiClient.get('/api/pull-requests'),
      ]);
      setMetrics(metricsRes.data);
      setPullRequests(prsRes.data ?? []);
    } catch (err) {
      console.error('Failed to load dashboard data:', err);
    } finally {
      setLoading(false);
    }
  }

  return (
    <Layout>
      <h1 style={{ fontSize: 16, fontWeight: 700, color: '#9ca3af', letterSpacing: 0.5, marginBottom: 20 }}>
        SYSTEM OVERVIEW (GLOBAL POSTURE)
      </h1>

      <MetricsGrid metrics={metrics} loading={loading} />

      <PullRequestTable
        pullRequests={pullRequests}
        loading={loading}
        selectedPRId={selectedPR?.id ?? null}
        onSelectPR={setSelectedPR}
      />

     <FindingDetails selectedPR={selectedPR} />
    </Layout>
  );
}

export default App;