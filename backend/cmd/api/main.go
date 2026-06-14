package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Divyansh670/opsmind-ai/backend/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (only in development — ignored if vars already set)
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, reading from system environment")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status":"ok","service":"opsmind-backend"}`)
	})

	addr := ":" + cfg.ServerPort
	fmt.Printf("OpsMind backend starting on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
