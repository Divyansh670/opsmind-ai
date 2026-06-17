package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Divyansh670/opsmind-ai/backend/internal/agents"
	"github.com/Divyansh670/opsmind-ai/backend/internal/api"
	"github.com/Divyansh670/opsmind-ai/backend/internal/config"
	"github.com/Divyansh670/opsmind-ai/backend/internal/db"
	"github.com/Divyansh670/opsmind-ai/backend/internal/models"
	"github.com/Divyansh670/opsmind-ai/backend/internal/webhook"
	"github.com/joho/godotenv"
)

func main() {
	envPaths := []string{".env", "../.env", "../../.env", "backend/.env"}
	envLoaded := false
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			log.Printf("Loaded .env from: %s", path)
			envLoaded = true
			break
		}
	}
	if !envLoaded {
		log.Println("No .env file found, reading from system environment")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	dbConfig := db.DefaultConfig(cfg.DSN())
	database, err := db.New(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create job queue (buffer of 100 jobs)
	jobQueue := make(chan models.WebhookPayload, 100)

	// Initialize webhook handler
	webhookHandler := webhook.NewHandler(cfg.GitHubWebhookSecret, jobQueue)

	// Simple job consumer (placeholder until agents are built)
	// Start worker pool
	groqClient := agents.NewGroqClient(cfg.GroqAPIKey, cfg.GroqModelID)
	githubClient := agents.NewGitHubClient(cfg.GitHubToken)
	repo := db.NewRepository(database)
	pool := agents.NewWorkerPool(cfg.MaxWorkers, jobQueue, groqClient, repo, githubClient)
	pool.Start()
	defer pool.Stop()
	// Set up HTTP router
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		dbErr := database.HealthCheck(ctx)
		status := "ok"
		dbStatus := "ok"
		httpCode := http.StatusOK

		if dbErr != nil {
			status = "degraded"
			dbStatus = "unreachable"
			httpCode = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpCode)
		json.NewEncoder(w).Encode(map[string]string{
			"status":   status,
			"service":  "opsmind-backend",
			"database": dbStatus,
		})
	})

	// GitHub webhook endpoint
	mux.HandleFunc("/webhook/github", webhookHandler.HandleWebhook)
	// Dashboard API endpoints
	dashboardHandler := api.NewDashboardHandler(repo)
	mux.HandleFunc("/api/metrics", dashboardHandler.HandleMetrics)
	mux.HandleFunc("/api/pull-requests", dashboardHandler.HandlePullRequests)
	mux.HandleFunc("/api/pull-requests/", dashboardHandler.HandleFindingsForPR)
	testHandler := webhook.NewTestTriggerHandler(jobQueue)
	mux.HandleFunc("/test/trigger", testHandler.HandleTestTrigger)

	// Build server
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      api.CORSMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		fmt.Printf("OpsMind backend starting on :%s\n", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown: %v", err)
	}
	fmt.Println("Server stopped.")
}
