package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Divyansh670/opsmind-ai/backend/internal/db"
)

// DashboardHandler serves dashboard-related read endpoints
type DashboardHandler struct {
	repo *db.Repository
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(repo *db.Repository) *DashboardHandler {
	return &DashboardHandler{repo: repo}
}

// HandleMetrics returns the 3 top-level dashboard metrics
func (h *DashboardHandler) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	criticalFlaws, costDrift, passRate, err := h.repo.GetDashboardMetrics(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"critical_open_flaws": criticalFlaws,
		"monthly_cost_drift":  costDrift,
		"pipeline_pass_rate":  passRate,
	})
}

// HandlePullRequests returns all PRs for the table view
func (h *DashboardHandler) HandlePullRequests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	prs, err := h.repo.GetAllPullRequests(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prs)
}

// HandleFindingsForPR returns findings for a specific PR ID, e.g. /api/pull-requests/5/findings
func (h *DashboardHandler) HandleFindingsForPR(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/pull-requests/")
	idStr := strings.TrimSuffix(path, "/findings")

	prID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid PR id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	findings, err := h.repo.GetFindingsForPR(ctx, prID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(findings)
}

// HandleDismissFinding handles POST /api/findings/{id}/dismiss
func (h *DashboardHandler) HandleDismissFinding(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/findings/")
	idStr := strings.TrimSuffix(path, "/dismiss")

	findingID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid finding id", http.StatusBadRequest)
		return
	}

	var body struct {
		Action string `json:"action"`
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if body.Action == "" {
		body.Action = "DISMISSED"
	}

	ctx := r.Context()
	if err := h.repo.DismissFinding(ctx, findingID, body.Action, body.Reason); err != nil {
		log.Printf("ERROR: failed to dismiss finding %d: %v", findingID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: finding %d dismissed with action=%s", findingID, body.Action)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "dismissed"})
}
