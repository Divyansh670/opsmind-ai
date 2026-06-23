package agents

import (
	"strings"
	"testing"

	"github.com/Divyansh670/opsmind-ai/backend/internal/models"
)

func TestFormatFindingsComment_NoIssues(t *testing.T) {
	secResult := &models.SecuritySentinelResponse{HasVulnerability: false}
	costResult := &models.CostPredictorResponse{HasDrift: false}
	archResult := &models.ArchitectureSupervisorResponse{HasIssues: false}

	comment := FormatFindingsComment(secResult, costResult, archResult)

	if !strings.Contains(comment, "No issues found") {
		t.Errorf("expected clean PR comment to contain 'No issues found', got: %s", comment)
	}
	if !strings.Contains(comment, "OpsMind AI") {
		t.Errorf("expected comment to contain 'OpsMind AI' branding, got: %s", comment)
	}
}

func TestFormatFindingsComment_WithSecurityFinding(t *testing.T) {
	secResult := &models.SecuritySentinelResponse{
		HasVulnerability: true,
		Vulnerabilities: []models.SecurityVulnerability{
			{
				Severity:           models.SeverityCritical,
				CWEID:              "CWE-798",
				FilePath:           "internal/auth/jwt.go",
				LineNumber:         39,
				ExploitExplanation: "Hardcoded secret key detected",
				RemediationSnippet: "Use os.Getenv instead",
			},
		},
	}
	costResult := &models.CostPredictorResponse{HasDrift: false}
	archResult := &models.ArchitectureSupervisorResponse{HasIssues: false}

	comment := FormatFindingsComment(secResult, costResult, archResult)

	if !strings.Contains(comment, "Security") {
		t.Errorf("expected security section in comment, got: %s", comment)
	}
	if !strings.Contains(comment, "CRITICAL") {
		t.Errorf("expected CRITICAL severity in comment, got: %s", comment)
	}
	if !strings.Contains(comment, "CWE-798") {
		t.Errorf("expected CWE ID in comment, got: %s", comment)
	}
}

func TestFormatFindingsComment_WithCostDrift(t *testing.T) {
	secResult := &models.SecuritySentinelResponse{HasVulnerability: false}
	costResult := &models.CostPredictorResponse{
		HasDrift:         true,
		DriftUSD:         1200.00,
		DriftExplanation: "RDS instance upgraded to db.r5.4xlarge",
		AffectedServices: []string{"RDS"},
	}
	archResult := &models.ArchitectureSupervisorResponse{HasIssues: false}

	comment := FormatFindingsComment(secResult, costResult, archResult)

	if !strings.Contains(comment, "Cost") {
		t.Errorf("expected cost section in comment, got: %s", comment)
	}
	if !strings.Contains(comment, "1200.00") {
		t.Errorf("expected cost amount in comment, got: %s", comment)
	}
}

func TestFormatFindingsComment_NilResults(t *testing.T) {
	comment := FormatFindingsComment(nil, nil, nil)
	if !strings.Contains(comment, "No issues found") {
		t.Errorf("expected nil results to produce clean comment, got: %s", comment)
	}
}

func TestSeverityEmoji(t *testing.T) {
	cases := []struct {
		severity string
		expected string
	}{
		{models.SeverityCritical, "🔴"},
		{models.SeverityHigh, "🟠"},
		{models.SeverityMedium, "🟡"},
		{models.SeverityLow, "🟢"},
		{"UNKNOWN", "⚪"},
	}

	for _, c := range cases {
		got := severityEmoji(c.severity)
		if got != c.expected {
			t.Errorf("severityEmoji(%s) = %s, want %s", c.severity, got, c.expected)
		}
	}
}
