import { useState, useEffect, useCallback } from 'react';
import { apiClient } from '../api/client';
import type { DashboardMetrics, PullRequest } from '../types/api';

const POLL_INTERVAL_MS = 30000; // 30 seconds

interface AuditStreamState {
  metrics: DashboardMetrics | null;
  pullRequests: PullRequest[];
  loading: boolean;
  lastUpdated: Date | null;
  refresh: () => void;
}

export function useAuditStream(): AuditStreamState {
  const [metrics, setMetrics] = useState<DashboardMetrics | null>(null);
  const [pullRequests, setPullRequests] = useState<PullRequest[]>([]);
  const [loading, setLoading] = useState(true);
  const [lastUpdated, setLastUpdated] = useState<Date | null>(null);

  const fetchData = useCallback(async () => {
    try {
      const [metricsRes, prsRes] = await Promise.all([
        apiClient.get('/api/metrics'),
        apiClient.get('/api/pull-requests'),
      ]);
      setMetrics(metricsRes.data);
      setPullRequests(prsRes.data ?? []);
      setLastUpdated(new Date());
    } catch (err) {
      console.error('Failed to fetch dashboard data:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  // Initial fetch
  useEffect(() => {
    fetchData();
  }, [fetchData]);

  // Polling interval
  useEffect(() => {
    const interval = setInterval(fetchData, POLL_INTERVAL_MS);
    return () => clearInterval(interval);
  }, [fetchData]);

  return {
    metrics,
    pullRequests,
    loading,
    lastUpdated,
    refresh: fetchData,
  };
}