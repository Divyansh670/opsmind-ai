package models

import "time"

// ==========================================
// Repository — maps to repositories table
// ==========================================
type Repository struct {
	ID        int       `json:"id"`
	RepoName  string    `json:"repo_name"`
	GitURL    string    `json:"git_url"`
	CreatedAt time.Time `json:"created_at"`
}

// ==========================================
// PullRequest — maps to pull_requests table
// ==========================================
type PullRequest struct {
	ID            int       `json:"id"`
	RepoID        int       `json:"repo_id"`
	PRNumber      int       `json:"pr_number"`
	HeadCommit    string    `json:"head_commit"`
	Author        string    `json:"author"`
	Status        string    `json:"status"`
	SecurityScore int       `json:"security_score"`
	CostDriftUSD  float64   `json:"cost_drift_usd"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PRStatus constants
const (
	PRStatusPending  = "PENDING"
	PRStatusApproved = "APPROVED"
	PRStatusFlagged  = "FLAGGED"
)

// ==========================================
// AgentFinding — maps to agent_findings table
// ==========================================
type AgentFinding struct {
	ID            int       `json:"id"`
	PRID          int       `json:"pr_id"`
	AgentName     string    `json:"agent_name"`
	Severity      string    `json:"severity"`
	CWEID         string    `json:"cwe_id"`
	FilePath      string    `json:"file_path"`
	LineNumber    int       `json:"line_number"`
	Description   string    `json:"description"`
	Remediation   string    `json:"remediation"`
	Dismissed     bool      `json:"dismissed"`
	DismissReason string    `json:"dismiss_reason"`
	CreatedAt     time.Time `json:"created_at"`
}

// Severity constants
const (
	SeverityCritical = "CRITICAL"
	SeverityHigh     = "HIGH"
	SeverityMedium   = "MEDIUM"
	SeverityLow      = "LOW"
)

// AgentName constants
const (
	AgentSecuritySentinel       = "SecuritySentinel"
	AgentCostPredictor          = "CostPredictor"
	AgentArchitectureSupervisor = "ArchitectureSupervisor"
)

// ==========================================
// FeedbackLog — maps to feedback_logs table
// ==========================================
type FeedbackLog struct {
	ID               int       `json:"id"`
	FindingID        int       `json:"finding_id"`
	DiffSnippet      string    `json:"diff_snippet"`
	AgentExplanation string    `json:"agent_explanation"`
	EngineerReason   string    `json:"engineer_reason"`
	Action           string    `json:"action"`
	CreatedAt        time.Time `json:"created_at"`
}

// FeedbackAction constants
const (
	FeedbackActionDismissed         = "DISMISSED"
	FeedbackActionApprovedException = "APPROVED_EXCEPTION"
)

// ==========================================
// ArchitectureRule — maps to architecture_rules table
// ==========================================
type ArchitectureRule struct {
	ID        int       `json:"id"`
	RuleText  string    `json:"rule_text"`
	CreatedAt time.Time `json:"created_at"`
}

// ==========================================
// Agent response structures (for LLM output)
// ==========================================

// SecurityVulnerability is what the Security Sentinel agent returns
type SecurityVulnerability struct {
	FilePath           string `json:"file_path"`
	LineNumber         int    `json:"line_number"`
	Severity           string `json:"severity"`
	CWEID              string `json:"cwe_id"`
	ExploitExplanation string `json:"exploit_explanation"`
	RemediationSnippet string `json:"remediation_snippet"`
}

// SecuritySentinelResponse is the full JSON response from Security Sentinel
type SecuritySentinelResponse struct {
	HasVulnerability bool                    `json:"has_vulnerability"`
	Vulnerabilities  []SecurityVulnerability `json:"vulnerabilities"`
}

// CostPredictorResponse is the full JSON response from Cost Predictor
type CostPredictorResponse struct {
	HasDrift         bool     `json:"has_drift"`
	DriftUSD         float64  `json:"drift_usd"`
	DriftExplanation string   `json:"drift_explanation"`
	AffectedServices []string `json:"affected_services"`
}

// ArchitectureIssue is a single issue found by Architecture Supervisor
type ArchitectureIssue struct {
	FilePath    string `json:"file_path"`
	LineNumber  int    `json:"line_number"`
	Description string `json:"description"`
	Suggestion  string `json:"suggestion"`
}

// ArchitectureSupervisorResponse is the full JSON response from Architecture Supervisor
type ArchitectureSupervisorResponse struct {
	HasIssues bool                `json:"has_issues"`
	Issues    []ArchitectureIssue `json:"issues"`
}

// ==========================================
// WebhookPayload — GitHub webhook structure
// ==========================================
type WebhookPayload struct {
	Action      string     `json:"action"`
	Number      int        `json:"number"`
	PullRequest GitHubPR   `json:"pull_request"`
	Repository  GitHubRepo `json:"repository"`
}

type GitHubPR struct {
	Number int        `json:"number"`
	Title  string     `json:"title"`
	User   GitHubUser `json:"user"`
	Head   GitHubHead `json:"head"`
	Body   string     `json:"body"`
}

type GitHubUser struct {
	Login string `json:"login"`
}

type GitHubHead struct {
	SHA string `json:"sha"`
	Ref string `json:"ref"`
}

type GitHubRepo struct {
	FullName string `json:"full_name"`
	CloneURL string `json:"clone_url"`
}
