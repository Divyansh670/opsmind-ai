package agents

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GitHubClient wraps calls to the GitHub REST API
type GitHubClient struct {
	Token      string
	HTTPClient *http.Client
}

// NewGitHubClient creates a new GitHub API client
func NewGitHubClient(token string) *GitHubClient {
	return &GitHubClient{
		Token: token,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// FetchPRDiff fetches the unified diff for a pull request.
// repoFullName is in the form "owner/repo", e.g. "Divyansh670/opsmind-ai"
func (c *GitHubClient) FetchPRDiff(ctx context.Context, repoFullName string, prNumber int) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%d", repoFullName, prNumber)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Asking GitHub for the diff format instead of JSON
	req.Header.Set("Accept", "application/vnd.github.v3.diff")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch PR diff: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("github API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	diffBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read diff body: %w", err)
	}

	return string(diffBytes), nil
}
