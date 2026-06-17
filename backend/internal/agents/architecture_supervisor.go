package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Divyansh670/opsmind-ai/backend/internal/models"
)

const architectureSupervisorSystemPrompt = `You are an Enterprise System Architect reviewing git diffs for architectural pattern violations.
Examine the provided code diff for structural anomalies such as:
- Direct database access from controller/handler layers (should use repository pattern)
- Business logic leaking into presentation/API layers
- Missing separation of concerns (e.g., mixing HTTP handling with data access)
- Tight coupling between unrelated modules
- Violation of dependency injection patterns
- Circular dependencies between packages
- God objects/functions doing too many unrelated things

You must respond exclusively with a valid JSON object matching this exact schema, with no additional text, no markdown formatting, and no code fences:
{
  "has_issues": boolean,
  "issues": [
    {
      "file_path": "string",
      "line_number": integer,
      "description": "string",
      "suggestion": "string"
    }
  ]
}

If no architectural issues are found, return {"has_issues": false, "issues": []}`

// ArchitectureSupervisorAgent wraps the Groq client for architecture analysis
type ArchitectureSupervisorAgent struct {
	Client *GroqClient
}

// NewArchitectureSupervisorAgent creates a new Architecture Supervisor agent
func NewArchitectureSupervisorAgent(client *GroqClient) *ArchitectureSupervisorAgent {
	return &ArchitectureSupervisorAgent{Client: client}
}

// Analyze sends the diff to the LLM and parses the structured architecture response
func (a *ArchitectureSupervisorAgent) Analyze(ctx context.Context, diff string) (*models.ArchitectureSupervisorResponse, error) {
	if strings.TrimSpace(diff) == "" {
		return &models.ArchitectureSupervisorResponse{HasIssues: false, Issues: nil}, nil
	}

	userPrompt := fmt.Sprintf("Analyze this git diff for architectural pattern violations:\n\n%s", diff)

	rawResponse, err := a.Client.Complete(ctx, architectureSupervisorSystemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("architecture supervisor LLM call failed: %w", err)
	}

	cleaned := cleanJSONResponse(rawResponse)

	var result models.ArchitectureSupervisorResponse
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		log.Printf("WARN: failed to parse architecture supervisor response, raw: %s", rawResponse)
		return nil, fmt.Errorf("failed to parse architecture supervisor JSON: %w", err)
	}

	return &result, nil
}
