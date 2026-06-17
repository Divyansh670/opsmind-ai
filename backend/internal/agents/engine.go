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

// processJob runs the security agent (and placeholders for the other 2) for a single PR
func (wp *WorkerPool) processJob(workerID int, payload models.WebhookPayload) {
	log.Printf("INFO: worker %d processing PR #%d from %s",
		workerID,
		payload.Number,
		payload.Repository.FullName,
	)

	// NOTE: real diff fetching from GitHub comes in a later step.
	// For now we use the PR body as a placeholder diff source for testing.
	diff := payload.PullRequest.Body

	secResult, err := wp.securityAgent.Analyze(wp.ctx, diff)
	if err != nil {
		log.Printf("ERROR: SecuritySentinel failed for PR #%d: %v", payload.Number, err)
	} else {
		log.Printf("INFO: SecuritySentinel completed for PR #%d — has_vulnerability=%v, findings=%d",
			payload.Number,
			secResult.HasVulnerability,
			len(secResult.Vulnerabilities),
		)
		for _, v := range secResult.Vulnerabilities {
			log.Printf("  -> [%s] %s:%d — %s", v.Severity, v.FilePath, v.LineNumber, v.ExploitExplanation)
		}
	}

	// CostPredictor and ArchitectureSupervisor wired in upcoming steps
	log.Printf("INFO: worker %d finished PR #%d", workerID, payload.Number)
}
