export interface Repository {
  id: number;
  repo_name: string;
  git_url: string;
  created_at: string;
}

export interface PullRequest {
  id: number;
  repo_id: number;
  pr_number: number;
  head_commit: string;
  author: string;
  status: 'PENDING' | 'APPROVED' | 'FLAGGED';
  security_score: number;
  cost_drift_usd: number;
  created_at: string;
  updated_at: string;
  repo_name?: string;
}

export interface AgentFinding {
  id: number;
  pr_id: number;
  agent_name: 'SecuritySentinel' | 'CostPredictor' | 'ArchitectureSupervisor';
  severity: 'CRITICAL' | 'HIGH' | 'MEDIUM' | 'LOW';
  cwe_id: string;
  file_path: string;
  line_number: number;
  description: string;
  remediation: string;
  dismissed: boolean;
  dismiss_reason: string;
  created_at: string;
}

export interface DashboardMetrics {
  critical_open_flaws: number;
  monthly_cost_drift: number;
  pipeline_pass_rate: number;
}