package agents

import (
	"context"
	"log"
	"sync"

	"github.com/Divyansh670/opsmind-ai/backend/internal/models"
)

// WorkerPool manages concurrent agent execution
type WorkerPool struct {
	jobChannel chan models.WebhookPayload
	maxWorkers int
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc

	securityAgent *SecuritySentinelAgent
	costAgent     *CostPredictorAgent
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(maxWorkers int, jobChannel chan models.WebhookPayload, groqClient *GroqClient) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		jobChannel:    jobChannel,
		maxWorkers:    maxWorkers,
		ctx:           ctx,
		cancel:        cancel,
		securityAgent: NewSecuritySentinelAgent(groqClient),
		costAgent:     NewCostPredictorAgent(groqClient),
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

// worker is a long-lived goroutine that processes jobs
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

// processJob runs Security Sentinel and Cost Predictor concurrently for a single PR
func (wp *WorkerPool) processJob(workerID int, payload models.WebhookPayload) {
	log.Printf("INFO: worker %d processing PR #%d from %s",
		workerID,
		payload.Number,
		payload.Repository.FullName,
	)

	diff := payload.PullRequest.Body

	var agentWg sync.WaitGroup

	// Run Security Sentinel
	agentWg.Add(1)
	go func() {
		defer agentWg.Done()
		secResult, err := wp.securityAgent.Analyze(wp.ctx, diff)
		if err != nil {
			log.Printf("ERROR: SecuritySentinel failed for PR #%d: %v", payload.Number, err)
			return
		}
		log.Printf("INFO: SecuritySentinel completed for PR #%d — has_vulnerability=%v, findings=%d",
			payload.Number,
			secResult.HasVulnerability,
			len(secResult.Vulnerabilities),
		)
		for _, v := range secResult.Vulnerabilities {
			log.Printf("  -> [SECURITY][%s] %s:%d — %s", v.Severity, v.FilePath, v.LineNumber, v.ExploitExplanation)
		}
	}()

	// Run Cost Predictor
	agentWg.Add(1)
	go func() {
		defer agentWg.Done()
		costResult, err := wp.costAgent.Analyze(wp.ctx, diff)
		if err != nil {
			log.Printf("ERROR: CostPredictor failed for PR #%d: %v", payload.Number, err)
			return
		}
		log.Printf("INFO: CostPredictor completed for PR #%d — has_drift=%v, drift_usd=%.2f",
			payload.Number,
			costResult.HasDrift,
			costResult.DriftUSD,
		)
		if costResult.HasDrift {
			log.Printf("  -> [COST] $%.2f/mo — %s (services: %v)",
				costResult.DriftUSD,
				costResult.DriftExplanation,
				costResult.AffectedServices,
			)
		}
	}()

	agentWg.Wait()

	log.Printf("INFO: worker %d finished PR #%d", workerID, payload.Number)
}
