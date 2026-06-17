package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Divyansh670/opsmind-ai/backend/internal/models"
)

// TestTriggerHandler manually injects a fake PR job into the queue for testing
type TestTriggerHandler struct {
	JobQueue chan models.WebhookPayload
}

// NewTestTriggerHandler creates a new test trigger handler
func NewTestTriggerHandler(jobQueue chan models.WebhookPayload) *TestTriggerHandler {
	return &TestTriggerHandler{JobQueue: jobQueue}
}

// HandleTestTrigger creates a fake PR payload with a known-vulnerable code snippet
// and pushes it into the job queue so we can verify the full agent pipeline works.
func (h *TestTriggerHandler) HandleTestTrigger(w http.ResponseWriter, r *http.Request) {
	vulnerableDiff := `
--- a/internal/auth/jwt.go
+++ b/internal/auth/jwt.go
@@ -39,7 +39,7 @@
 func GenerateToken(user string) string {
-	secret := os.Getenv("JWT_PRODUCTION_SECRET")
+	secret := "AIzaSyD-unSafeProductionKeyVariableConstantValue"
 	token := jwt.New(jwt.SigningMethodHS256)
 	return token.SignedString([]byte(secret))
 }
`

	testPayload := models.WebhookPayload{
		Action: "opened",
		Number: 9999,
		PullRequest: models.GitHubPR{
			Number: 9999,
			Title:  "Test: hardcoded secret injection",
			User:   models.GitHubUser{Login: "test-user"},
			Head:   models.GitHubHead{SHA: "test-sha-123", Ref: "test-branch"},
			Body:   vulnerableDiff,
		},
		Repository: models.GitHubRepo{
			FullName: "test-org/test-repo",
			CloneURL: "https://github.com/test-org/test-repo.git",
		},
	}

	select {
	case h.JobQueue <- testPayload:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status":"test job queued, check server logs for agent output"}`)
	default:
		http.Error(w, `{"status":"job queue full"}`, http.StatusServiceUnavailable)
	}

	_ = json.NewEncoder // keep import used if expanded later
}
