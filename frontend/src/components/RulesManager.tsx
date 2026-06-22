import { useEffect, useState } from 'react';
import { Trash2, Plus, Building2 } from 'lucide-react';
import { apiClient } from '../api/client';

interface ArchitectureRule {
  id: number;
  rule_text: string;
  created_at: string;
}

export default function RulesManager() {
  const [rules, setRules] = useState<ArchitectureRule[]>([]);
  const [loading, setLoading] = useState(true);
  const [newRule, setNewRule] = useState('');
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    loadRules();
  }, []);

  async function loadRules() {
    setLoading(true);
    try {
      const res = await apiClient.get('/api/rules');
      setRules(res.data ?? []);
    } catch (err) {
      console.error('Failed to load rules:', err);
    } finally {
      setLoading(false);
    }
  }

  async function handleAddRule() {
    if (!newRule.trim()) return;
    setSaving(true);
    setError('');
    try {
      await apiClient.post('/api/rules', { rule_text: newRule.trim() });
      setNewRule('');
      await loadRules();
    } catch (err) {
      setError('Failed to save rule. Please try again.');
      console.error('Failed to add rule:', err);
    } finally {
      setSaving(false);
    }
  }

  async function handleDeleteRule(id: number) {
    try {
      await apiClient.delete(`/api/rules/${id}`);
      setRules(prev => prev.filter(r => r.id !== id));
    } catch (err) {
      console.error('Failed to delete rule:', err);
    }
  }

  return (
    <div style={styles.container}>
      <div style={styles.header}>
        <Building2 size={18} color="#60a5fa" />
        <h2 style={styles.title}>CUSTOM ARCHITECTURE RULES</h2>
      </div>

      <p style={styles.description}>
        Rules defined here are embedded using AI and automatically injected into
        the Architecture Supervisor agent when analyzing relevant pull requests
        via semantic (pgvector) similarity search.
      </p>

      <div style={styles.addSection}>
        <textarea
          style={styles.textarea}
          placeholder="e.g. Database access must always go through the repository layer, never directly from HTTP handlers."
          value={newRule}
          onChange={e => setNewRule(e.target.value)}
          rows={3}
        />
        {error && <p style={styles.error}>{error}</p>}
        <button
          style={{
            ...styles.addBtn,
            opacity: saving || !newRule.trim() ? 0.5 : 1,
            cursor: saving || !newRule.trim() ? 'not-allowed' : 'pointer',
          }}
          onClick={handleAddRule}
          disabled={saving || !newRule.trim()}
        >
          <Plus size={14} />
          {saving ? 'Generating embedding...' : 'Add Rule'}
        </button>
      </div>

      <div style={styles.rulesList}>
        {loading && <p style={styles.empty}>Loading rules...</p>}
        {!loading && rules.length === 0 && (
          <p style={styles.empty}>
            No custom rules yet. Add your first rule above to start enforcing
            company-specific architecture standards.
          </p>
        )}
        {rules.map(rule => (
          <div key={rule.id} style={styles.ruleCard}>
            <p style={styles.ruleText}>{rule.rule_text}</p>
            <div style={styles.ruleFooter}>
              <span style={styles.ruleDate}>
                Added {new Date(rule.created_at).toLocaleDateString()}
              </span>
              <button
                style={styles.deleteBtn}
                onClick={() => handleDeleteRule(rule.id)}
                title="Delete rule"
              >
                <Trash2 size={13} />
                Delete
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

const styles: { [key: string]: React.CSSProperties } = {
  container: {
    backgroundColor: '#13151f',
    border: '1px solid #1f2330',
    borderRadius: 10,
    padding: 24,
  },
  header: {
    display: 'flex',
    alignItems: 'center',
    gap: 8,
    marginBottom: 10,
  },
  title: {
    fontSize: 13,
    fontWeight: 700,
    color: '#9ca3af',
    letterSpacing: 0.5,
  },
  description: {
    fontSize: 13,
    color: '#6b7280',
    lineHeight: 1.6,
    marginBottom: 20,
  },
  addSection: {
    marginBottom: 24,
  },
  textarea: {
    width: '100%',
    backgroundColor: '#0f1117',
    border: '1px solid #1f2330',
    borderRadius: 8,
    padding: '10px 12px',
    color: '#e5e7eb',
    fontSize: 13,
    lineHeight: 1.5,
    resize: 'vertical' as const,
    outline: 'none',
    marginBottom: 10,
    boxSizing: 'border-box' as const,
    fontFamily: 'inherit',
  },
  addBtn: {
    display: 'flex',
    alignItems: 'center',
    gap: 6,
    backgroundColor: '#1d4ed8',
    border: 'none',
    borderRadius: 6,
    color: '#fff',
    fontSize: 13,
    fontWeight: 600,
    padding: '8px 16px',
    cursor: 'pointer',
  },
  error: {
    color: '#f87171',
    fontSize: 12,
    marginBottom: 8,
  },
  rulesList: {
    display: 'flex',
    flexDirection: 'column' as const,
    gap: 10,
  },
  ruleCard: {
    backgroundColor: '#0f1117',
    border: '1px solid #1f2330',
    borderRadius: 8,
    padding: '12px 14px',
  },
  ruleText: {
    fontSize: 13,
    color: '#d1d5db',
    lineHeight: 1.5,
    marginBottom: 8,
  },
  ruleFooter: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  ruleDate: {
    fontSize: 11,
    color: '#4b5563',
  },
  deleteBtn: {
    display: 'flex',
    alignItems: 'center',
    gap: 4,
    backgroundColor: 'transparent',
    border: '1px solid #374151',
    borderRadius: 6,
    color: '#6b7280',
    fontSize: 11,
    padding: '3px 8px',
    cursor: 'pointer',
  },
  empty: {
    color: '#6b7280',
    fontSize: 13,
  },
};