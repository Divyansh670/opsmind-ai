import { useEffect, useState } from 'react';
import { GitBranch, AlertTriangle, CheckCircle2, DollarSign } from 'lucide-react';
import { fetchRepos } from '../api/client';
import type { RepoStats } from '../types/api';

export default function RepositoriesPage() {
  const [repos, setRepos] = useState<RepoStats[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchRepos()
      .then(setRepos)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  return (
    <div>
      <div style={styles.pageHeader}>
        <h1 style={styles.pageTitle}>REPOSITORIES</h1>
      </div>

      {loading && (
        <div style={styles.empty}>Loading repositories...</div>
      )}

      {!loading && repos.length === 0 && (
        <div style={styles.emptyCard}>
          <GitBranch size={32} color="#374151" />
          <p style={styles.emptyText}>No repositories tracked yet.</p>
          <p style={styles.emptySubtext}>
            Repositories appear here automatically when their PRs are analyzed via webhook.
          </p>
        </div>
      )}

      <div style={styles.grid}>
        {repos.map(repo => (
          <RepoCard key={repo.id} repo={repo} />
        ))}
      </div>
    </div>
  );
}

function RepoCard({ repo }: { repo: RepoStats }) {
  const passRate = repo.total_prs > 0
    ? Math.round((repo.approved_prs / repo.total_prs) * 100)
    : 0;

  const riskColor = repo.flagged_prs > 0 ? '#f87171' : '#34d399';

  return (
    <div style={styles.card}>
      <div style={styles.cardHeader}>
        <div style={styles.repoName}>
          <GitBranch size={16} color="#60a5fa" />
          <span style={styles.repoNameText}>{repo.repo_name}</span>
        </div>
        <span style={{
          ...styles.riskBadge,
          borderColor: riskColor,
          color: riskColor,
        }}>
          {repo.flagged_prs > 0 ? 'AT RISK' : 'HEALTHY'}
        </span>
      </div>

      <div style={styles.statsGrid}>
        <div style={styles.stat}>
          <span style={styles.statValue}>{repo.total_prs}</span>
          <span style={styles.statLabel}>Total PRs</span>
        </div>
        <div style={styles.stat}>
          <span style={{ ...styles.statValue, color: '#f87171' }}>{repo.flagged_prs}</span>
          <span style={styles.statLabel}>Flagged</span>
        </div>
        <div style={styles.stat}>
          <span style={{ ...styles.statValue, color: '#34d399' }}>{repo.approved_prs}</span>
          <span style={styles.statLabel}>Approved</span>
        </div>
        <div style={styles.stat}>
          <span style={{ ...styles.statValue, color: '#fbbf24' }}>{passRate}%</span>
          <span style={styles.statLabel}>Pass Rate</span>
        </div>
      </div>

      <div style={styles.divider} />

      <div style={styles.bottomRow}>
        <div style={styles.bottomStat}>
          <AlertTriangle size={13} color="#f87171" />
          <span style={styles.bottomStatText}>
            Avg Score: <strong>{Math.round(repo.avg_security_score)}</strong>
          </span>
        </div>
        <div style={styles.bottomStat}>
          <DollarSign size={13} color="#fbbf24" />
          <span style={styles.bottomStatText}>
            Drift: <strong>${repo.total_cost_drift_usd.toFixed(0)}/mo</strong>
          </span>
        </div>
        <div style={styles.bottomStat}>
          <CheckCircle2 size={13} color="#6b7280" />
          <span style={styles.bottomStatText}>
            Updated: <strong>{repo.last_updated ?? 'Never'}</strong>
          </span>
        </div>
      </div>
    </div>
  );
}

const styles: { [key: string]: React.CSSProperties } = {
  pageHeader: {
    marginBottom: 20,
  },
  pageTitle: {
    fontSize: 16,
    fontWeight: 700,
    color: '#9ca3af',
    letterSpacing: 0.5,
  },
  empty: {
    color: '#6b7280',
    fontSize: 13,
  },
  emptyCard: {
    backgroundColor: '#13151f',
    border: '1px solid #1f2330',
    borderRadius: 10,
    padding: 40,
    textAlign: 'center' as const,
    display: 'flex',
    flexDirection: 'column' as const,
    alignItems: 'center',
    gap: 12,
  },
  emptyText: {
    color: '#6b7280',
    fontSize: 14,
    fontWeight: 600,
  },
  emptySubtext: {
    color: '#4b5563',
    fontSize: 13,
  },
  grid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fill, minmax(340px, 1fr))',
    gap: 16,
  },
  card: {
    backgroundColor: '#13151f',
    border: '1px solid #1f2330',
    borderRadius: 10,
    padding: 20,
  },
  cardHeader: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    marginBottom: 16,
  },
  repoName: {
    display: 'flex',
    alignItems: 'center',
    gap: 8,
  },
  repoNameText: {
    fontSize: 14,
    fontWeight: 600,
    color: '#e5e7eb',
  },
  riskBadge: {
    fontSize: 10,
    fontWeight: 700,
    padding: '3px 8px',
    borderRadius: 6,
    border: '1px solid',
    letterSpacing: 0.5,
  },
  statsGrid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(4, 1fr)',
    gap: 8,
    marginBottom: 16,
  },
  stat: {
    display: 'flex',
    flexDirection: 'column' as const,
    alignItems: 'center',
    backgroundColor: '#0f1117',
    borderRadius: 6,
    padding: '8px 4px',
  },
  statValue: {
    fontSize: 18,
    fontWeight: 700,
    color: '#e5e7eb',
  },
  statLabel: {
    fontSize: 10,
    color: '#6b7280',
    marginTop: 2,
  },
  divider: {
    borderTop: '1px solid #1f2330',
    marginBottom: 12,
  },
  bottomRow: {
    display: 'flex',
    justifyContent: 'space-between',
  },
  bottomStat: {
    display: 'flex',
    alignItems: 'center',
    gap: 4,
  },
  bottomStatText: {
    fontSize: 11,
    color: '#6b7280',
  },
};