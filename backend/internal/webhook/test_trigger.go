package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Divyansh670/opsmind-ai/backend/internal/models"
)

// TestTriggerHandler manually injects a fake PR job into the queue for testing
type TestTriggerHandler struct {
	JobQueue chan models.WebhookPayload
	mu       sync.Mutex
	lastHit  time.Time
}

// NewTestTriggerHandler creates a new test trigger handler
func NewTestTriggerHandler(jobQueue chan models.WebhookPayload) *TestTriggerHandler {
	return &TestTriggerHandler{JobQueue: jobQueue}
}

const testTriggerCooldown = 30 * time.Second

// HandleTestTrigger creates a fake PR payload with a known-vulnerable code snippet
// and pushes it into the job queue so we can verify the full agent pipeline works.
func (h *TestTriggerHandler) HandleTestTrigger(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	if time.Since(h.lastHit) < testTriggerCooldown {
		remaining := testTriggerCooldown - time.Since(h.lastHit)
		h.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintf(w, `{"error":"rate limited, try again in %.0f seconds"}`, remaining.Seconds())
		return
	}
	h.lastHit = time.Now()
	h.mu.Unlock()

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

--- a/infra/main.tf
+++ b/infra/main.tf
@@ -10,7 +10,7 @@
 resource "aws_db_instance" "primary" {
-  instance_class = "db.t3.medium"
+  instance_class = "db.r5.4xlarge"
   allocated_storage = 100
   engine = "postgres"
 }

--- a/internal/handlers/user_handler.go
+++ b/internal/handlers/user_handler.go
@@ -15,6 +15,10 @@
 func GetUserHandler(w http.ResponseWriter, r *http.Request) {
+	db, _ := sql.Open("postgres", connStr)
+	rows, _ := db.Query("SELECT * FROM users WHERE id = $1", r.URL.Query().Get("id"))
+	defer rows.Close()
+	// directly querying DB from handler instead of using repository layer
 	userID := r.URL.Query().Get("id")
 	json.NewEncoder(w).Encode(userID)
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
