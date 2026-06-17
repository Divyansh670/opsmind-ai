import type { PullRequest } from '../types/api';

interface PullRequestTableProps {
  pullRequests: PullRequest[];
  loading: boolean;
  selectedPRId: number | null;
  onSelectPR: (pr: PullRequest) => void;
}

export default function PullRequestTable({
  pullRequests,
  loading,
  selectedPRId,
  onSelectPR,
}: PullRequestTableProps) {
  return (
    <div style={styles.container}>
      <h2 style={styles.title}>RECENT ACTIVE PULL REQUEST AUDITS</h2>

      <table style={styles.table}>
        <thead>
          <tr>
            <th style={styles.th}>PR #</th>
            <th style={styles.th}>REPOSITORY</th>
            <th style={styles.th}>AUTHOR</th>
            <th style={styles.th}>SECURITY SCORE</th>
            <th style={styles.th}>DRIFT STATUS</th>
          </tr>
        </thead>
        <tbody>
          {loading && (
            <tr>
              <td style={styles.td} colSpan={5}>
                Loading...
              </td>
            </tr>
          )}
          {!loading && pullRequests.length === 0 && (
            <tr>
              <td style={styles.td} colSpan={5}>
                No pull requests analyzed yet.
              </td>
            </tr>
          )}
          {pullRequests.map((pr) => (
            <tr
              key={pr.id}
              onClick={() => onSelectPR(pr)}
              style={{
                ...styles.row,
                backgroundColor: selectedPRId === pr.id ? '#1a1d29' : 'transparent',
              }}
            >
              <td style={styles.td}>#{pr.pr_number}</td>
              <td style={styles.td}>{pr.repo_name}</td>
              <td style={styles.td}>{pr.author}</td>
              <td style={styles.td}>
                <ScoreBadge score={pr.security_score} />
              </td>
              <td style={styles.td}>
                <StatusBadge status={pr.status} />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function ScoreBadge({ score }: { score: number }) {
  let label = 'LOW';
  let color = '#34d399';
  if (score >= 80) {
    label = 'HIGH';
    color = '#f87171';
  } else if (score >= 40) {
    label = 'MED';
    color = '#fbbf24';
  }

  return (
    <span style={{ ...styles.badge, borderColor: color, color }}>
      {score} - {label}
    </span>
  );
}

function StatusBadge({ status }: { status: string }) {
  const colorMap: { [key: string]: string } = {
    APPROVED: '#34d399',
    FLAGGED: '#f87171',
    PENDING: '#fbbf24',
  };
  const color = colorMap[status] ?? '#9ca3af';

  return (
    <span style={{ ...styles.badge, borderColor: color, color }}>
      {status}
    </span>
  );
}

const styles: { [key: string]: React.CSSProperties } = {
  container: {
    backgroundColor: '#13151f',
    border: '1px solid #1f2330',
    borderRadius: 10,
    padding: 20,
    marginBottom: 28,
  },
  title: {
    fontSize: 13,
    fontWeight: 700,
    color: '#9ca3af',
    letterSpacing: 0.5,
    marginBottom: 16,
  },
  table: {
    width: '100%',
    borderCollapse: 'collapse',
  },
  th: {
    textAlign: 'left',
    fontSize: 11,
    color: '#6b7280',
    fontWeight: 600,
    padding: '8px 12px',
    borderBottom: '1px solid #1f2330',
  },
  td: {
    fontSize: 13,
    padding: '12px 12px',
    borderBottom: '1px solid #1a1d29',
  },
  row: {
    cursor: 'pointer',
  },
  badge: {
    display: 'inline-block',
    padding: '3px 10px',
    borderRadius: 6,
    border: '1px solid',
    fontSize: 11,
    fontWeight: 600,
  },
};