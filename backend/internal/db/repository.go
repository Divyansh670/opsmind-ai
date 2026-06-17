package db

import (
	"context"
	"fmt"

	"github.com/Divyansh670/opsmind-ai/backend/internal/models"
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
