import { useEffect, useState } from 'react';
import { Shield, DollarSign, Building2 } from 'lucide-react';
import { apiClient } from '../api/client';
import type { AgentFinding, PullRequest } from '../types/api';

interface FindingDetailsProps {
  selectedPR: PullRequest | null;
}

export default function FindingDetails({ selectedPR }: FindingDetailsProps) {
  const [findings, setFindings] = useState<AgentFinding[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (selectedPR) {
      loadFindings(selectedPR.id);
    } else {
      setFindings([]);
    }
  }, [selectedPR]);

  async function loadFindings(prId: number) {
    setLoading(true);
    try {
      const res = await apiClient.get(`/api/pull-requests/${prId}/findings`);
      setFindings(res.data ?? []);
    } catch (err) {
      console.error('Failed to load findings:', err);
      setFindings([]);
    } finally {
      setLoading(false);
    }
  }

  if (!selectedPR) {
    return (
      <div style={styles.container}>
        <p style={styles.emptyText}>Select a pull request above to view detailed findings.</p>
      </div>
    );
  }

  return (
    <div style={styles.container}>
      <h2 style={styles.title}>
        CONTEXTUAL FINDINGS PANEL: PR #{selectedPR.pr_number} ({selectedPR.repo_name})
      </h2>

      {loading && <p style={styles.emptyText}>Loading findings...</p>}

      {!loading && findings.length === 0 && (
        <p style={styles.emptyText}>✅ No findings for this PR — clean audit.</p>
      )}

      {!loading &&
        findings.map((finding) => (
          <FindingCard key={finding.id} finding={finding} />
        ))}
    </div>
  );
}

function FindingCard({ finding }: { finding: AgentFinding }) {
  const agentMeta = getAgentMeta(finding.agent_name);
  const severityColor = getSeverityColor(finding.severity);

  return (
    <div style={styles.findingCard}>
      <div style={styles.findingHeader}>
        <div style={styles.agentBadge}>
          {agentMeta.icon}
          <span style={{ color: agentMeta.color }}>{finding.agent_name}</span>
        </div>
        <span style={{ ...styles.severityBadge, borderColor: severityColor, color: severityColor }}>
          {finding.severity}
        </span>
        {finding.cwe_id && <span style={styles.cweTag}>{finding.cwe_id}</span>}
      </div>

      {finding.file_path && (
        <div style={styles.fileLine}>
          File: <code style={styles.code}>{finding.file_path}</code>
          {finding.line_number > 0 && <> | Line: <code style={styles.code}>{finding.line_number}</code></>}
        </div>
      )}

      <p style={styles.description}>{finding.description}</p>

      {finding.remediation && (
        <div style={styles.remediationBox}>
          <strong style={{ color: '#34d399' }}>Suggested fix:</strong>{' '}
          <span style={styles.remediationText}>{finding.remediation}</span>
        </div>
      )}
    </div>
  );
}

function getAgentMeta(agentName: string) {
  switch (agentName) {
    case 'SecuritySentinel':
      return { icon: <Shield size={14} color="#f87171" />, color: '#f87171' };
    case 'CostPredictor':
      return { icon: <DollarSign size={14} color="#fbbf24" />, color: '#fbbf24' };
    case 'ArchitectureSupervisor':
      return { icon: <Building2 size={14} color="#60a5fa" />, color: '#60a5fa' };
    default:
      return { icon: null, color: '#9ca3af' };
  }
}

function getSeverityColor(severity: string) {
  switch (severity) {
    case 'CRITICAL':
      return '#f87171';
    case 'HIGH':
      return '#fb923c';
    case 'MEDIUM':
      return '#fbbf24';
    case 'LOW':
      return '#34d399';
    default:
      return '#9ca3af';
  }
}

const styles: { [key: string]: React.CSSProperties } = {
  container: {
    backgroundColor: '#13151f',
    border: '1px solid #1f2330',
    borderRadius: 10,
    padding: 20,
  },
  title: {
    fontSize: 13,
    fontWeight: 700,
    color: '#9ca3af',
    letterSpacing: 0.5,
    marginBottom: 16,
  },
  emptyText: {
    color: '#6b7280',
    fontSize: 13,
  },
  findingCard: {
    backgroundColor: '#0f1117',
    border: '1px solid #1f2330',
    borderRadius: 8,
    padding: 14,
    marginBottom: 12,
  },
  findingHeader: {
    display: 'flex',
    alignItems: 'center',
    gap: 10,
    marginBottom: 10,
  },
  agentBadge: {
    display: 'flex',
    alignItems: 'center',
    gap: 6,
    fontSize: 13,
    fontWeight: 600,
  },
  severityBadge: {
    padding: '2px 8px',
    borderRadius: 6,
    border: '1px solid',
    fontSize: 11,
    fontWeight: 700,
  },
  cweTag: {
    fontSize: 11,
    color: '#6b7280',
    backgroundColor: '#1a1d29',
    padding: '2px 8px',
    borderRadius: 6,
  },
  fileLine: {
    fontSize: 12,
    color: '#9ca3af',
    marginBottom: 8,
  },
  code: {
    backgroundColor: '#1a1d29',
    padding: '1px 6px',
    borderRadius: 4,
    color: '#e5e7eb',
  },
  description: {
    fontSize: 13,
    color: '#d1d5db',
    lineHeight: 1.5,
    marginBottom: 8,
},
  remediationBox: {
    fontSize: 12,
    backgroundColor: '#0d1f17',
    border: '1px solid #1f3328',
    borderRadius: 6,
    padding: '8px 10px',
    marginTop: 8,
  },
  remediationText: {
    color: '#a7f3d0',
  },
};