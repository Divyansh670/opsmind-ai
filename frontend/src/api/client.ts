import axios from 'axios';
import type { PullRequest, AgentFinding } from '../types/api';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080';

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 60000,
});

// Fetches health status of backend + database
export async function fetchHealth() {
  const res = await apiClient.get('/health');
  return res.data;
}

// These endpoints don't exist on the backend yet — we'll build them in Step 23
export async function fetchPullRequests(): Promise<PullRequest[]> {
  const res = await apiClient.get('/api/pull-requests');
  return res.data;
}

export async function fetchFindingsForPR(prId: number): Promise<AgentFinding[]> {
  const res = await apiClient.get(`/api/pull-requests/${prId}/findings`);
  return res.data;
}
export async function dismissFinding(
  findingId: number,
  action: 'DISMISSED' | 'APPROVED_EXCEPTION',
  reason: string
): Promise<void> {
  await apiClient.post(`/api/findings/${findingId}/dismiss`, {
    action,
    reason,
  });
}
import type { PRTrendPoint } from '../types/api';

export async function fetchTrend(): Promise<PRTrendPoint[]> {
  const res = await apiClient.get('/api/trend');
  return res.data ?? [];
}