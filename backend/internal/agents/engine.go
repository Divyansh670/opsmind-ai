package agents

import (
	"context"
	"log"
	"sync"

	"github.com/Divyansh670/opsmind-ai/backend/internal/models"
)

// Job represents a single PR analysis task
type Job struct {
	Payload models.WebhookPayload
}

// WorkerPool manages concurrent agent execution
type WorkerPool struct {
	jobChannel chan models.WebhookPayload
	maxWorkers int
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(maxWorkers int, jobChannel chan models.WebhookPayload) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		jobChannel: jobChannel,
		maxWorkers: maxWorkers,
		ctx:        ctx,
		cancel:     cancel,
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

// processJob runs all 3 agents concurrently for a single PR
func (wp *WorkerPool) processJob(workerID int, payload models.WebhookPayload) {
	log.Printf("INFO: worker %d processing PR #%d from %s",
		workerID,
		payload.Number,
		payload.Repository.FullName,
	)

	// Run all 3 agents concurrently using a WaitGroup
	var agentWg sync.WaitGroup
	results := make(chan AgentResult, 3)

	// Agent 1: Security Sentinel
	agentWg.Add(1)
	go func() {
		defer agentWg.Done()
		result := runSecuritySentinel(payload)
		results <- result
	}()

	// Agent 2: Cost Predictor
	agentWg.Add(1)
	go func() {
		defer agentWg.Done()
		result := runCostPredictor(payload)
		results <- result
	}()

	// Agent 3: Architecture Supervisor
	agentWg.Add(1)
	go func() {
		defer agentWg.Done()
		result := runArchitectureSupervisor(payload)
		results <- result
	}()

	// Close results channel when all agents finish
	go func() {
		agentWg.Wait()
		close(results)
	}()

	// Collect results
	for result := range results {
		if result.Error != nil {
			log.Printf("ERROR: agent %s failed for PR #%d: %v",
				result.AgentName,
				payload.Number,
				result.Error,
			)
			continue
		}
		log.Printf("INFO: agent %s completed for PR #%d — findings: %d",
			result.AgentName,
			payload.Number,
			result.FindingCount,
		)
	}

	log.Printf("INFO: worker %d finished PR #%d", workerID, payload.Number)
}

// AgentResult holds the outcome of a single agent run
type AgentResult struct {
	AgentName    string
	FindingCount int
	Error        error
}

// Placeholder agent functions — will be replaced with real LLM calls in next steps
func runSecuritySentinel(payload models.WebhookPayload) AgentResult {
	log.Printf("INFO: SecuritySentinel analyzing PR #%d", payload.Number)
	// Real implementation coming in Step 12
	return AgentResult{
		AgentName:    models.AgentSecuritySentinel,
		FindingCount: 0,
		Error:        nil,
	}
}

func runCostPredictor(payload models.WebhookPayload) AgentResult {
	log.Printf("INFO: CostPredictor analyzing PR #%d", payload.Number)
	// Real implementation coming in Step 13
	return AgentResult{
		AgentName:    models.AgentCostPredictor,
		FindingCount: 0,
		Error:        nil,
	}
}

func runArchitectureSupervisor(payload models.WebhookPayload) AgentResult {
	log.Printf("INFO: ArchitectureSupervisor analyzing PR #%d", payload.Number)
	// Real implementation coming in Step 14
	return AgentResult{
		AgentName:    models.AgentArchitectureSupervisor,
		FindingCount: 0,
		Error:        nil,
	}
}
