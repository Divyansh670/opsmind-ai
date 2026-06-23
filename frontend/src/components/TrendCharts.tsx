import { useEffect, useState } from 'react';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  BarChart,
  Bar,
  Cell,
} from 'recharts';
import { fetchTrend } from '../api/client';
import type { PRTrendPoint } from '../types/api';

export default function TrendCharts() {
  const [data, setData] = useState<PRTrendPoint[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchTrend()
      .then(setData)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  if (loading) return null;
  if (data.length === 0) return null;

  return (
    <div style={styles.grid}>
      {/* Security Score Trend */}
      <div style={styles.card}>
        <h3 style={styles.title}>SECURITY SCORE TREND</h3>
        <ResponsiveContainer width="100%" height={180}>
          <LineChart data={data} margin={{ top: 5, right: 10, left: -20, bottom: 5 }}>
            <CartesianGrid strokeDasharray="3 3" stroke="#1f2330" />
            <XAxis
              dataKey="pr_number"
              tickFormatter={(v) => `#${v}`}
              tick={{ fill: '#6b7280', fontSize: 11 }}
            />
            <YAxis domain={[0, 100]} tick={{ fill: '#6b7280', fontSize: 11 }} />
            <Tooltip
              contentStyle={{ backgroundColor: '#13151f', border: '1px solid #1f2330', borderRadius: 6 }}
              labelFormatter={(v) => `PR #${v}`}
              formatter={(value: unknown) => [value as number, 'Security Score']}
            />
            <Line
              type="monotone"
              dataKey="security_score"
              stroke="#f87171"
              strokeWidth={2}
              dot={{ fill: '#f87171', r: 3 }}
              activeDot={{ r: 5 }}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>

      {/* Cost Drift Trend */}
      <div style={styles.card}>
        <h3 style={styles.title}>COST DRIFT PER PR (USD/mo)</h3>
        <ResponsiveContainer width="100%" height={180}>
          <BarChart data={data} margin={{ top: 5, right: 10, left: -10, bottom: 5 }}>
            <CartesianGrid strokeDasharray="3 3" stroke="#1f2330" />
            <XAxis
              dataKey="pr_number"
              tickFormatter={(v) => `#${v}`}
              tick={{ fill: '#6b7280', fontSize: 11 }}
            />
            <YAxis tick={{ fill: '#6b7280', fontSize: 11 }} />
            <Tooltip
              contentStyle={{ backgroundColor: '#13151f', border: '1px solid #1f2330', borderRadius: 6 }}
              labelFormatter={(v) => `PR #${v}`}
              formatter={(value: unknown) => [`$${(value as number).toFixed(0)}`, 'Cost Drift/mo']}
            />
            <Bar dataKey="cost_drift_usd" radius={[4, 4, 0, 0]}>
              {data.map((entry, index) => (
                <Cell
                  key={`cell-${index}`}
                  fill={entry.cost_drift_usd > 500 ? '#f87171' : entry.cost_drift_usd > 0 ? '#fbbf24' : '#34d399'}
                />
              ))}
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}

const styles: { [key: string]: React.CSSProperties } = {
  grid: {
    display: 'grid',
    gridTemplateColumns: '1fr 1fr',
    gap: 16,
    marginBottom: 28,
  },
  card: {
    backgroundColor: '#13151f',
    border: '1px solid #1f2330',
    borderRadius: 10,
    padding: '16px 20px',
  },
  title: {
    fontSize: 12,
    fontWeight: 700,
    color: '#6b7280',
    letterSpacing: 0.5,
    marginBottom: 12,
  },
};