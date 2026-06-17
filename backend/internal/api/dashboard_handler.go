package api

import (
	"encoding/json"
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
