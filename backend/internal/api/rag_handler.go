package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Divyansh670/opsmind-ai/backend/internal/agents"
	"github.com/Divyansh670/opsmind-ai/backend/internal/db"
)

// RAGHandler handles the chatbot question-answer endpoint
type RAGHandler struct {
	repo         *db.Repository
	geminiClient *agents.GeminiClient
	groqClient   *agents.GroqClient
}

// NewRAGHandler creates a new RAG handler
func NewRAGHandler(repo *db.Repository, gemini *agents.GeminiClient, groq *agents.GroqClient) *RAGHandler {
	return &RAGHandler{repo: repo, geminiClient: gemini, groqClient: groq}
}

type RAGRequest struct {
	Question string `json:"question"`
}

type RAGSource struct {
	Type     string `json:"type"`
	RepoName string `json:"repo_name"`
	PRNumber int    `json:"pr_number"`
	FilePath string `json:"file_path"`
	Severity string `json:"severity"`
	Snippet  string `json:"snippet"`
}

type RAGResponse struct {
	Answer  string      `json:"answer"`
	Sources []RAGSource `json:"sources"`
}

const ragSystemPrompt = `You are OpsMind AI Assistant, an expert DevSecOps engineer embedded in an engineering dashboard.

You help engineers understand security vulnerabilities, cloud cost impacts, and architectural issues found in their pull requests.

You are given CONTEXT from the system's database — real findings from real PRs, plus custom architecture rules defined by the team.

Rules:
- Answer ONLY based on the provided context. Do not make up findings.
- Be specific: reference file paths, line numbers, severity levels, and PR numbers when available.
- Be concise but actionable — give engineers something they can act on immediately.
- If the context doesn't contain enough information to answer, say so clearly.
- Format code suggestions in markdown code blocks.`

// HandleRAGQuery handles POST /api/chat
func (h *RAGHandler) HandleRAGQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RAGRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Question) == "" {
		http.Error(w, "question is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Step 1: Embed the question
	embedding, err := h.geminiClient.Embed(ctx, req.Question)
	if err != nil {
		log.Printf("ERROR: failed to embed question: %v", err)
		http.Error(w, "failed to process question", http.StatusInternalServerError)
		return
	}

	// Step 2: Retrieve relevant context from pgvector + findings
	contexts, err := h.repo.SearchRAGContext(ctx, embedding, 8)
	if err != nil {
		log.Printf("ERROR: failed to search RAG context: %v", err)
		http.Error(w, "failed to retrieve context", http.StatusInternalServerError)
		return
	}

	// Step 3: Build context string and sources list
	var contextBuilder strings.Builder
	var sources []RAGSource

	contextBuilder.WriteString("CONTEXT FROM OPSMIND DATABASE:\n\n")

	for i, c := range contexts {
		contextBuilder.WriteString(fmt.Sprintf("[Source %d] Type: %s", i+1, c.SourceType))
		if c.RepoName != "" {
			contextBuilder.WriteString(fmt.Sprintf(" | Repo: %s", c.RepoName))
		}
		if c.PRNumber > 0 {
			contextBuilder.WriteString(fmt.Sprintf(" | PR #%d", c.PRNumber))
		}
		if c.Severity != "" {
			contextBuilder.WriteString(fmt.Sprintf(" | Severity: %s", c.Severity))
		}
		if c.FilePath != "" {
			contextBuilder.WriteString(fmt.Sprintf(" | File: %s", c.FilePath))
		}
		contextBuilder.WriteString(fmt.Sprintf("\nContent: %s\n\n", c.Content))

		if c.Content != "" {
			snippet := c.Content
			if len(snippet) > 120 {
				snippet = snippet[:120] + "..."
			}
			sources = append(sources, RAGSource{
				Type:     c.SourceType,
				RepoName: c.RepoName,
				PRNumber: c.PRNumber,
				FilePath: c.FilePath,
				Severity: c.Severity,
				Snippet:  snippet,
			})
		}
	}

	// Step 4: Call Groq with context + question
	userPrompt := fmt.Sprintf("%s\n\nENGINEER QUESTION: %s", contextBuilder.String(), req.Question)

	answer, err := h.groqClient.Complete(context.Background(), ragSystemPrompt, userPrompt)
	if err != nil {
		log.Printf("ERROR: RAG Groq call failed: %v", err)
		http.Error(w, "failed to generate answer", http.StatusInternalServerError)
		return
	}

	// Step 5: Return answer + sources
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RAGResponse{
		Answer:  answer,
		Sources: sources,
	})
}
