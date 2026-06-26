package db

import (
	"context"
	"fmt"
	"time"

	"github.com/Divyansh670/opsmind-ai/backend/internal/models"
	"github.com/pgvector/pgvector-go"
)

// Repository handles all database read/write operations
type Repository struct {
	db *DB
}

// NewRepository creates a new repository wrapping the DB connection pool
func NewRepository(db *DB) *Repository {
	return &Repository{db: db}
}

// UpsertRepository ensures a repository row exists, returns its ID
func (r *Repository) UpsertRepository(ctx context.Context, repoName, gitURL string) (int, error) {
	var id int
	query := `
		INSERT INTO repositories (repo_name, git_url)
		VALUES ($1, $2)
		ON CONFLICT (repo_name) DO UPDATE SET repo_name = EXCLUDED.repo_name
		RETURNING id
	`
	err := r.db.Pool.QueryRow(ctx, query, repoName, gitURL).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert repository: %w", err)
	}
	return id, nil
}

// UpsertPullRequest ensures a PR row exists, returns its ID
func (r *Repository) UpsertPullRequest(ctx context.Context, repoID, prNumber int, headCommit, author string) (int, error) {
	var id int
	query := `
		INSERT INTO pull_requests (repo_id, pr_number, head_commit, author, status)
		VALUES ($1, $2, $3, $4, 'PENDING')
		ON CONFLICT (repo_id, pr_number) DO UPDATE SET
			head_commit = EXCLUDED.head_commit,
			updated_at = NOW()
		RETURNING id
	`
	err := r.db.Pool.QueryRow(ctx, query, repoID, prNumber, headCommit, author).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert pull request: %w", err)
	}
	return id, nil
}

// UpdatePRStatus updates the status, security score, and cost drift on a PR
func (r *Repository) UpdatePRStatus(ctx context.Context, prID int, status string, securityScore int, costDriftUSD float64) error {
	query := `
		UPDATE pull_requests
		SET status = $1, security_score = $2, cost_drift_usd = $3, updated_at = NOW()
		WHERE id = $4
	`
	_, err := r.db.Pool.Exec(ctx, query, status, securityScore, costDriftUSD, prID)
	if err != nil {
		return fmt.Errorf("failed to update PR status: %w", err)
	}
	return nil
}

// InsertFinding saves a single agent finding to the database
func (r *Repository) InsertFinding(ctx context.Context, finding models.AgentFinding) error {
	query := `
		INSERT INTO agent_findings
			(pr_id, agent_name, severity, cwe_id, file_path, line_number, description, remediation)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Pool.Exec(ctx, query,
		finding.PRID,
		finding.AgentName,
		finding.Severity,
		finding.CWEID,
		finding.FilePath,
		finding.LineNumber,
		finding.Description,
		finding.Remediation,
	)
	if err != nil {
		return fmt.Errorf("failed to insert finding: %w", err)
	}
	return nil
}

// PullRequestWithRepo is a PR joined with its repository name, used for dashboard display
type PullRequestWithRepo struct {
	models.PullRequest
	RepoName string `json:"repo_name"`
}

// GetAllPullRequests fetches all PRs joined with their repo name, newest first
func (r *Repository) GetAllPullRequests(ctx context.Context) ([]PullRequestWithRepo, error) {
	query := `
		SELECT pr.id, pr.repo_id, pr.pr_number, pr.head_commit, pr.author, pr.status,
		       pr.security_score, pr.cost_drift_usd, pr.created_at, pr.updated_at,
		       repo.repo_name
		FROM pull_requests pr
		JOIN repositories repo ON repo.id = pr.repo_id
		ORDER BY pr.updated_at DESC
		LIMIT 50
	`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pull requests: %w", err)
	}
	defer rows.Close()

	var prs []PullRequestWithRepo
	for rows.Next() {
		var pr PullRequestWithRepo
		err := rows.Scan(
			&pr.ID, &pr.RepoID, &pr.PRNumber, &pr.HeadCommit, &pr.Author, &pr.Status,
			&pr.SecurityScore, &pr.CostDriftUSD, &pr.CreatedAt, &pr.UpdatedAt,
			&pr.RepoName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pull request row: %w", err)
		}
		prs = append(prs, pr)
	}
	return prs, nil
}

// GetFindingsForPR fetches all findings for a given PR ID
func (r *Repository) GetFindingsForPR(ctx context.Context, prID int) ([]models.AgentFinding, error) {
	query := `
		SELECT id, pr_id, agent_name, severity, cwe_id, file_path, line_number,
		       description, remediation, dismissed, dismiss_reason, created_at
		FROM agent_findings
		WHERE pr_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, prID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch findings: %w", err)
	}
	defer rows.Close()

	var findings []models.AgentFinding
	for rows.Next() {
		var f models.AgentFinding
		var cweID, filePath, dismissReason *string
		err := rows.Scan(
			&f.ID, &f.PRID, &f.AgentName, &f.Severity, &cweID, &filePath, &f.LineNumber,
			&f.Description, &f.Remediation, &f.Dismissed, &dismissReason, &f.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan finding row: %w", err)
		}
		if cweID != nil {
			f.CWEID = *cweID
		}
		if filePath != nil {
			f.FilePath = *filePath
		}
		if dismissReason != nil {
			f.DismissReason = *dismissReason
		}
		findings = append(findings, f)
	}
	return findings, nil
}

// GetDashboardMetrics computes the 3 top-level metrics for the dashboard
func (r *Repository) GetDashboardMetrics(ctx context.Context) (criticalFlaws int, costDrift float64, passRate float64, err error) {
	err = r.db.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM agent_findings WHERE severity = 'CRITICAL' AND dismissed = FALSE
	`).Scan(&criticalFlaws)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to count critical flaws: %w", err)
	}

	err = r.db.Pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(cost_drift_usd), 0) FROM pull_requests
		WHERE created_at > NOW() - INTERVAL '30 days'
	`).Scan(&costDrift)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to sum cost drift: %w", err)
	}

	var total, approved int
	err = r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM pull_requests`).Scan(&total)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to count total PRs: %w", err)
	}
	err = r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM pull_requests WHERE status = 'APPROVED'`).Scan(&approved)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to count approved PRs: %w", err)
	}

	if total > 0 {
		passRate = (float64(approved) / float64(total)) * 100
	}

	return criticalFlaws, costDrift, passRate, nil
}

// DismissFinding marks a finding as dismissed and logs the feedback
func (r *Repository) DismissFinding(ctx context.Context, findingID int, action, reason string) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Mark finding as dismissed
	_, err = tx.Exec(ctx, `
		UPDATE agent_findings
		SET dismissed = TRUE, dismiss_reason = $1
		WHERE id = $2
	`, reason, findingID)
	if err != nil {
		return fmt.Errorf("failed to dismiss finding: %w", err)
	}

	// Log to feedback_logs for MLOps retraining later
	_, err = tx.Exec(ctx, `
		INSERT INTO feedback_logs (finding_id, engineer_reason, action)
		VALUES ($1, $2, $3)
	`, findingID, reason, action)
	if err != nil {
		return fmt.Errorf("failed to insert feedback log: %w", err)
	}

	return tx.Commit(ctx)
}

// ArchitectureRuleInput is what's needed to create a new rule
type ArchitectureRuleInput struct {
	RuleText  string
	Embedding []float32
}

// ArchitectureRule represents a stored rule with its embedding
type ArchitectureRule struct {
	ID        int       `json:"id"`
	RuleText  string    `json:"rule_text"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateArchitectureRule inserts a new rule with its embedding
func (r *Repository) CreateArchitectureRule(ctx context.Context, input ArchitectureRuleInput) (int, error) {
	var id int
	vec := pgvector.NewVector(input.Embedding)
	err := r.db.Pool.QueryRow(ctx, `
		INSERT INTO architecture_rules (rule_text, embedding)
		VALUES ($1, $2)
		RETURNING id
	`, input.RuleText, vec).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create architecture rule: %w", err)
	}
	return id, nil
}

// GetAllArchitectureRules returns all rules (without embeddings, for display)
func (r *Repository) GetAllArchitectureRules(ctx context.Context) ([]ArchitectureRule, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, rule_text, created_at FROM architecture_rules ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch architecture rules: %w", err)
	}
	defer rows.Close()

	var rules []ArchitectureRule
	for rows.Next() {
		var rule ArchitectureRule
		if err := rows.Scan(&rule.ID, &rule.RuleText, &rule.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan architecture rule: %w", err)
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

// DeleteArchitectureRule removes a rule
func (r *Repository) DeleteArchitectureRule(ctx context.Context, id int) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM architecture_rules WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete architecture rule: %w", err)
	}
	return nil
}

// FindRelevantRules uses pgvector cosine similarity to find rules most relevant to a diff
func (r *Repository) FindRelevantRules(ctx context.Context, diffEmbedding []float32, limit int) ([]string, error) {
	vec := pgvector.NewVector(diffEmbedding)
	rows, err := r.db.Pool.Query(ctx, `
		SELECT rule_text FROM architecture_rules
		ORDER BY embedding <=> $1
		LIMIT $2
	`, vec, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find relevant rules: %w", err)
	}
	defer rows.Close()

	var rules []string
	for rows.Next() {
		var rule string
		if err := rows.Scan(&rule); err != nil {
			return nil, fmt.Errorf("failed to scan relevant rule: %w", err)
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

// PRTrendPoint represents a single data point for trend charts
type PRTrendPoint struct {
	Date          string  `json:"date"`
	SecurityScore int     `json:"security_score"`
	CostDrift     float64 `json:"cost_drift_usd"`
	Status        string  `json:"status"`
	PRNumber      int     `json:"pr_number"`
}

// GetPRTrend returns the last 20 PRs ordered by date for trend charting
func (r *Repository) GetPRTrend(ctx context.Context) ([]PRTrendPoint, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT 
			TO_CHAR(pr.created_at, 'Mon DD') as date,
			pr.security_score,
			pr.cost_drift_usd,
			pr.status,
			pr.pr_number
		FROM pull_requests pr
		ORDER BY pr.created_at ASC
		LIMIT 20
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PR trend: %w", err)
	}
	defer rows.Close()

	var points []PRTrendPoint
	for rows.Next() {
		var p PRTrendPoint
		if err := rows.Scan(&p.Date, &p.SecurityScore, &p.CostDrift, &p.Status, &p.PRNumber); err != nil {
			return nil, fmt.Errorf("failed to scan trend point: %w", err)
		}
		points = append(points, p)
	}
	return points, nil
}

// RepoStats represents a repository with aggregated PR statistics
type RepoStats struct {
	ID          int     `json:"id"`
	RepoName    string  `json:"repo_name"`
	TotalPRs    int     `json:"total_prs"`
	FlaggedPRs  int     `json:"flagged_prs"`
	ApprovedPRs int     `json:"approved_prs"`
	AvgScore    float64 `json:"avg_security_score"`
	TotalDrift  float64 `json:"total_cost_drift_usd"`
	LastUpdated string  `json:"last_updated"`
}

// GetRepoStats returns all repositories with aggregated statistics
func (r *Repository) GetRepoStats(ctx context.Context) ([]RepoStats, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT
			repo.id,
			repo.repo_name,
			COUNT(pr.id) as total_prs,
			COUNT(CASE WHEN pr.status = 'FLAGGED' THEN 1 END) as flagged_prs,
			COUNT(CASE WHEN pr.status = 'APPROVED' THEN 1 END) as approved_prs,
			COALESCE(AVG(pr.security_score), 0) as avg_security_score,
			COALESCE(SUM(pr.cost_drift_usd), 0) as total_cost_drift_usd,
			TO_CHAR(MAX(pr.updated_at), 'Mon DD, YYYY') as last_updated
		FROM repositories repo
		LEFT JOIN pull_requests pr ON pr.repo_id = repo.id
		GROUP BY repo.id, repo.repo_name
		ORDER BY total_prs DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repo stats: %w", err)
	}
	defer rows.Close()

	var stats []RepoStats
	for rows.Next() {
		var s RepoStats
		if err := rows.Scan(
			&s.ID, &s.RepoName, &s.TotalPRs, &s.FlaggedPRs, &s.ApprovedPRs,
			&s.AvgScore, &s.TotalDrift, &s.LastUpdated,
		); err != nil {
			return nil, fmt.Errorf("failed to scan repo stats: %w", err)
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// RAGContext represents a retrieved piece of context for the chatbot
type RAGContext struct {
	SourceType string  `json:"source_type"` // "finding", "rule", "pr"
	Content    string  `json:"content"`
	RepoName   string  `json:"repo_name"`
	PRNumber   int     `json:"pr_number"`
	Severity   string  `json:"severity"`
	FilePath   string  `json:"file_path"`
	Score      float64 `json:"score"`
}

// SearchRAGContext finds relevant findings and rules using hybrid search
// (pgvector similarity + PostgreSQL full-text search, merged and deduplicated)
func (r *Repository) SearchRAGContext(ctx context.Context, questionEmbedding []float32, question string, limit int) ([]RAGContext, error) {
	vec := pgvector.NewVector(questionEmbedding)
	seen := make(map[string]bool)
	var results []RAGContext

	// 1. Vector similarity search on architecture rules
	ruleRows, err := r.db.Pool.Query(ctx, `
		SELECT 
			'rule' as source_type,
			rule_text as content,
			'' as repo_name,
			0 as pr_number,
			'' as severity,
			'' as file_path,
			(1 - (embedding <=> $1)) as score
		FROM architecture_rules
		ORDER BY embedding <=> $1
		LIMIT $2
	`, vec, limit/3)
	if err != nil {
		return nil, fmt.Errorf("failed to search rules: %w", err)
	}
	defer ruleRows.Close()

	for ruleRows.Next() {
		var c RAGContext
		if err := ruleRows.Scan(&c.SourceType, &c.Content, &c.RepoName, &c.PRNumber, &c.Severity, &c.FilePath, &c.Score); err != nil {
			return nil, err
		}
		key := c.SourceType + ":" + c.Content[:min(len(c.Content), 50)]
		if !seen[key] {
			seen[key] = true
			results = append(results, c)
		}
	}

	// 2. Full-text search on findings (exact keyword matching)
	ftsRows, err := r.db.Pool.Query(ctx, `
		SELECT
			'finding' as source_type,
			af.description || ' Remediation: ' || af.remediation as content,
			r.repo_name,
			pr.pr_number,
			af.severity,
			COALESCE(af.file_path, '') as file_path,
			ts_rank(
				to_tsvector('english', af.description || ' ' || COALESCE(af.file_path, '') || ' ' || COALESCE(af.cwe_id, '')),
				plainto_tsquery('english', $1)
			) as score
		FROM agent_findings af
		JOIN pull_requests pr ON pr.id = af.pr_id
		JOIN repositories r ON r.id = pr.repo_id
		WHERE 
			af.dismissed = FALSE
			AND to_tsvector('english', af.description || ' ' || COALESCE(af.file_path, '') || ' ' || COALESCE(af.cwe_id, ''))
				@@ plainto_tsquery('english', $1)
		ORDER BY score DESC
		LIMIT $2
	`, question, limit/3)
	if err != nil {
		return nil, fmt.Errorf("failed to do full-text search: %w", err)
	}
	defer ftsRows.Close()

	for ftsRows.Next() {
		var c RAGContext
		if err := ftsRows.Scan(&c.SourceType, &c.Content, &c.RepoName, &c.PRNumber, &c.Severity, &c.FilePath, &c.Score); err != nil {
			return nil, err
		}
		key := fmt.Sprintf("finding:%s:pr%d", c.FilePath, c.PRNumber)
		if !seen[key] {
			seen[key] = true
			results = append(results, c)
		}
	}

	// 3. Recent findings fallback (ensures we always have context even with no keyword match)
	recentRows, err := r.db.Pool.Query(ctx, `
		SELECT
			'finding' as source_type,
			af.description || ' Remediation: ' || af.remediation as content,
			r.repo_name,
			pr.pr_number,
			af.severity,
			COALESCE(af.file_path, '') as file_path,
			0.3 as score
		FROM agent_findings af
		JOIN pull_requests pr ON pr.id = af.pr_id
		JOIN repositories r ON r.id = pr.repo_id
		WHERE af.dismissed = FALSE
		ORDER BY af.created_at DESC
		LIMIT $1
	`, limit/3)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch recent findings: %w", err)
	}
	defer recentRows.Close()

	for recentRows.Next() {
		var c RAGContext
		if err := recentRows.Scan(&c.SourceType, &c.Content, &c.RepoName, &c.PRNumber, &c.Severity, &c.FilePath, &c.Score); err != nil {
			return nil, err
		}
		key := fmt.Sprintf("finding:%s:pr%d", c.FilePath, c.PRNumber)
		if !seen[key] {
			seen[key] = true
			results = append(results, c)
		}
	}

	return results, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
