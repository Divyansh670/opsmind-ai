package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Divyansh670/opsmind-ai/backend/internal/models"
)

// Handler handles incoming GitHub webhook requests
type Handler struct {
	WebhookSecret string
	JobQueue      chan models.WebhookPayload
}

// NewHandler creates a new webhook handler
func NewHandler(secret string, jobQueue chan models.WebhookPayload) *Handler {
	return &Handler{
		WebhookSecret: secret,
		JobQueue:      jobQueue,
	}
}

// HandleWebhook is the HTTP handler for GitHub webhook POST requests
func (h *Handler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR: failed to read webhook body: %v", err)
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate HMAC-SHA256 signature if secret is configured
	if h.WebhookSecret != "" {
		signature := r.Header.Get("X-Hub-Signature-256")
		if !h.validateSignature(body, signature) {
			log.Printf("ERROR: invalid webhook signature")
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}
	}

	// Only process pull_request events
	eventType := r.Header.Get("X-GitHub-Event")
	if eventType != "pull_request" {
		log.Printf("INFO: ignoring event type: %s", eventType)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status":"ignored"}`)
		return
	}

	// Parse the webhook payload
	var payload models.WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("ERROR: failed to parse webhook payload: %v", err)
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// Only process opened or synchronize actions
	if payload.Action != "opened" && payload.Action != "synchronize" {
		log.Printf("INFO: ignoring PR action: %s", payload.Action)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status":"ignored"}`)
		return
	}

	log.Printf("INFO: received PR #%d from %s", payload.Number, payload.Repository.FullName)

	// Queue the job for async processing — return 200 immediately to GitHub
	select {
	case h.JobQueue <- payload:
		log.Printf("INFO: queued PR #%d for processing", payload.Number)
	default:
		log.Printf("WARN: job queue full, dropping PR #%d", payload.Number)
	}

	// Always return 200 OK immediately to GitHub
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status":"queued"}`)
}

// validateSignature verifies the HMAC-SHA256 signature from GitHub
func (h *Handler) validateSignature(body []byte, signature string) bool {
	if signature == "" {
		return false
	}

	// GitHub sends "sha256=<hash>"
	if len(signature) < 7 || signature[:7] != "sha256=" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(h.WebhookSecret))
	mac.Write(body)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	expectedSignature := "sha256=" + expectedMAC

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
