package agents

import (
	"context"
	"log"
	"sync"

	"github.com/Divyansh670/opsmind-ai/backend/internal/db"
	"github.com/Divyansh670/opsmind-ai/backend/internal/models"
)

// WorkerPool manages concurrent agent execution
type WorkerPool struct {
	jobChannel chan models.WebhookPayload
	maxWorkers int
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc

	securityAgent     *SecuritySentinelAgent
	costAgent         *CostPredictorAgent
	architectureAgent *ArchitectureSupervisorAgent
	repo              *db.Repository
	githubClient      *GitHubClient
	geminiClient      *GeminiClient
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(maxWorkers int, jobChannel chan models.WebhookPayload, groqClient *GroqClient, geminiClient *GeminiClient, repo *db.Repository, githubClient *GitHubClient) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		jobChannel:        jobChannel,
		maxWorkers:        maxWorkers,
		ctx:               ctx,
		cancel:            cancel,
		securityAgent:     NewSecuritySentinelAgent(groqClient),
		costAgent:         NewCostPredictorAgent(groqClient),
		architectureAgent: NewArchitectureSupervisorAgent(groqClient, geminiClient, repo),
		repo:              repo,
		githubClient:      githubClient,
		geminiClient:      geminiClient,
	}
}

// Start launches all worker goroutines
func (wp *WorkerPool) Start() {
	log.Printf("INFO: starting worker pool with %d workers", wp.maxWorkers)
	for i := 0; i < wp.maxWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// Stop gracefully shuts down all workers
func (wp *WorkerPool) Stop() {
	log.Println("INFO: stopping worker pool...")
	wp.cancel()
	wp.wg.Wait()
	log.Println("INFO: worker pool stopped")
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	log.Printf("INFO: worker %d started", id)

	for {
		select {
		case payload, ok := <-wp.jobChannel:
			if !ok {
				log.Printf("INFO: worker %d shutting down — job channel closed", id)
				return
			}
			wp.processJob(id, payload)
		case <-wp.ctx.Done():
			log.Printf("INFO: worker %d shutting down — context cancelled", id)
			return
		}
	}
}

// processJob runs all 3 agents concurrently, persists findings, and posts a GitHub comment
func (wp *WorkerPool) processJob(workerID int, payload models.WebhookPayload) {
	log.Printf("INFO: worker %d processing PR #%d from %s",
		workerID, payload.Number, payload.Repository.FullName)

	ctx := wp.ctx
	diff := payload.PullRequest.Body

	realDiff, err := wp.githubClient.FetchPRDiff(ctx, payload.Repository.FullName, payload.Number)
	if err != nil {
		log.Printf("WARN: failed to fetch real diff for PR #%d, falling back to PR body: %v", payload.Number, err)
	} else if realDiff != "" {
		diff = realDiff
		log.Printf("INFO: fetched real diff for PR #%d (%d bytes)", payload.Number, len(diff))
	}

	repoID, err := wp.repo.UpsertRepository(ctx, payload.Repository.FullName, payload.Repository.CloneURL)
	if err != nil {
		log.Printf("ERROR: failed to upsert repository for PR #%d: %v", payload.Number, err)
		return
	}

	prID, err := wp.repo.UpsertPullRequest(ctx, repoID, payload.Number, payload.PullRequest.Head.SHA, payload.PullRequest.User.Login)
	if err != nil {
		log.Printf("ERROR: failed to upsert PR #%d: %v", payload.Number, err)
		return
	}

	var agentWg sync.WaitGroup
	var mu sync.Mutex
	highestSeverityScore := 0
	totalCostDrift := 0.0
	hasCritical := false

	severityToScore := map[string]int{
		models.SeverityCritical: 90,
		models.SeverityHigh:     70,
		models.SeverityMedium:   45,
		models.SeverityLow:      20,
	}

	var secResult *models.SecuritySentinelResponse
	var costResult *models.CostPredictorResponse
	var archResult *models.ArchitectureSupervisorResponse

	// Security Sentinel
	agentWg.Add(1)
	go func() {
		defer agentWg.Done()
		result, err := wp.securityAgent.Analyze(ctx, diff)
		if err != nil {
			log.Printf("ERROR: SecuritySentinel failed for PR #%d: %v", payload.Number, err)
			return
		}
		log.Printf("INFO: SecuritySentinel completed for PR #%d — findings=%d", payload.Number, len(result.Vulnerabilities))

		mu.Lock()
		secResult = result
		mu.Unlock()

		for _, v := range result.Vulnerabilities {
			finding := models.AgentFinding{
				PRID: prID, AgentName: models.AgentSecuritySentinel, Severity: v.Severity,
				CWEID: v.CWEID, FilePath: v.FilePath, LineNumber: v.LineNumber,
				Description: v.ExploitExplanation, Remediation: v.RemediationSnippet,
			}
			if err := wp.repo.InsertFinding(ctx, finding); err != nil {
				log.Printf("ERROR: failed to save security finding: %v", err)
			}
			mu.Lock()
			if score := severityToScore[v.Severity]; score > highestSeverityScore {
				highestSeverityScore = score
			}
			if v.Severity == models.SeverityCritical {
				hasCritical = true
			}
			mu.Unlock()
		}
	}()

	// Cost Predictor
	agentWg.Add(1)
	go func() {
		defer agentWg.Done()
		result, err := wp.costAgent.Analyze(ctx, diff)
		if err != nil {
			log.Printf("ERROR: CostPredictor failed for PR #%d: %v", payload.Number, err)
			return
		}
		log.Printf("INFO: CostPredictor completed for PR #%d — drift_usd=%.2f", payload.Number, result.DriftUSD)

		mu.Lock()
		costResult = result
		mu.Unlock()

		if result.HasDrift {
			finding := models.AgentFinding{
				PRID: prID, AgentName: models.AgentCostPredictor, Severity: models.SeverityMedium,
				Description: result.DriftExplanation,
			}
			if err := wp.repo.InsertFinding(ctx, finding); err != nil {
				log.Printf("ERROR: failed to save cost finding: %v", err)
			}
			mu.Lock()
			totalCostDrift += result.DriftUSD
			mu.Unlock()
		}
	}()

	// Architecture Supervisor
	agentWg.Add(1)
	go func() {
		defer agentWg.Done()
		result, err := wp.architectureAgent.Analyze(ctx, diff)
		if err != nil {
			log.Printf("ERROR: ArchitectureSupervisor failed for PR #%d: %v", payload.Number, err)
			return
		}
		log.Printf("INFO: ArchitectureSupervisor completed for PR #%d — issues=%d", payload.Number, len(result.Issues))

		mu.Lock()
		archResult = result
		mu.Unlock()

		for _, issue := range result.Issues {
			finding := models.AgentFinding{
				PRID: prID, AgentName: models.AgentArchitectureSupervisor, Severity: models.SeverityLow,
				FilePath: issue.FilePath, LineNumber: issue.LineNumber,
				Description: issue.Description, Remediation: issue.Suggestion,
			}
			if err := wp.repo.InsertFinding(ctx, finding); err != nil {
				log.Printf("ERROR: failed to save architecture finding: %v", err)
			}
		}
	}()

	agentWg.Wait()

	finalStatus := models.PRStatusApproved
	if hasCritical {
		finalStatus = models.PRStatusFlagged
	} else if highestSeverityScore > 0 {
		finalStatus = models.PRStatusPending
	}

	if err := wp.repo.UpdatePRStatus(ctx, prID, finalStatus, highestSeverityScore, totalCostDrift); err != nil {
		log.Printf("ERROR: failed to update PR status for PR #%d: %v", payload.Number, err)
	}

	commentBody := FormatFindingsComment(secResult, costResult, archResult)
	if err := wp.githubClient.PostPRComment(ctx, payload.Repository.FullName, payload.Number, commentBody); err != nil {
		log.Printf("WARN: failed to post GitHub comment for PR #%d: %v", payload.Number, err)
	} else {
		log.Printf("INFO: posted summary comment to GitHub for PR #%d", payload.Number)
	}

	log.Printf("INFO: worker %d finished PR #%d — status=%s, score=%d, cost_drift=$%.2f",
		workerID, payload.Number, finalStatus, highestSeverityScore, totalCostDrift)
}
