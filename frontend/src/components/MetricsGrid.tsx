import { AlertTriangle, DollarSign, CheckCircle2 } from 'lucide-react';
import type { DashboardMetrics } from '../types/api';

interface MetricsGridProps {
  metrics: DashboardMetrics | null;
  loading: boolean;
}

export default function MetricsGrid({ metrics, loading }: MetricsGridProps) {
  return (
    <div style={styles.grid}>
      <MetricCard
        icon={<AlertTriangle size={18} color="#f87171" />}
        label="CRITICAL OPEN FLAWS"
        value={loading ? '—' : `${metrics?.critical_open_flaws ?? 0} Active`}
        accentColor="#f87171"
      />
      <MetricCard
        icon={<DollarSign size={18} color="#fbbf24" />}
        label="MONTHLY COST DRIFT"
        value={
          loading
            ? '—'
            : `${metrics && metrics.monthly_cost_drift >= 0 ? '+' : ''}$${(
                metrics?.monthly_cost_drift ?? 0
              ).toFixed(2)} / mo`
        }
        accentColor="#fbbf24"
      />
      <MetricCard
        icon={<CheckCircle2 size={18} color="#34d399" />}
        label="PIPELINE PASS RATE"
        value={loading ? '—' : `${(metrics?.pipeline_pass_rate ?? 0).toFixed(1)}%`}
        accentColor="#34d399"
      />
    </div>
  );
}

interface MetricCardProps {
  icon: React.ReactNode;
  label: string;
  value: string;
  accentColor: string;
}

function MetricCard({ icon, label, value, accentColor }: MetricCardProps) {
  return (
    <div style={styles.card}>
      <div style={styles.cardHeader}>
        {icon}
        <span style={styles.cardLabel}>{label}</span>
      </div>
      <div style={{ ...styles.cardValue, color: accentColor }}>{value}</div>
    </div>
  );
}

const styles: { [key: string]: React.CSSProperties } = {
  grid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(3, 1fr)',
    gap: 16,
    marginBottom: 28,
  },
  card: {
    backgroundColor: '#13151f',
    border: '1px solid #1f2330',
    borderRadius: 10,
    padding: '18px 20px',
  },
  cardHeader: {
    display: 'flex',
    alignItems: 'center',
    gap: 8,
    marginBottom: 12,
  },
  cardLabel: {
    fontSize: 12,
    fontWeight: 600,
    color: '#9ca3af',
    letterSpacing: 0.5,
  },
  cardValue: {
    fontSize: 26,
    fontWeight: 700,
  },
};