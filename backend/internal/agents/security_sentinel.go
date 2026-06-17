package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Divyansh670/opsmind-ai/backend/internal/models"
)

const securitySentinelSystemPrompt = `You are a strict automated application security engine auditing git diff inputs.
Analyze the provided code diff for critical code security vulnerabilities such as:
- Hardcoded secrets, API keys, or credentials
- SQL injection vectors
- Cross-site scripting (XSS) vulnerabilities
- Insecure deserialization
- Unsafe or outdated dependency usage
- Broken authentication or authorization logic

You must respond exclusively with a valid JSON object matching this exact schema, with no additional text, no markdown formatting, and no code fences:
{
  "has_vulnerability": boolean,
  "vulnerabilities": [
    {
      "file_path": "string",
      "line_number": integer,
      "severity": "CRITICAL" | "HIGH" | "MEDIUM" | "LOW",
      "cwe_id": "string",
      "exploit_explanation": "string",
      "remediation_snippet": "string"
    }
  ]
}

If no vulnerabilities are found, return {"has_vulnerability": false, "vulnerabilities": []}`

// SecuritySentinelAgent wraps the Groq client for security analysis
type SecuritySentinelAgent struct {
	Client *GroqClient
}

// NewSecuritySentinelAgent creates a new Security Sentinel agent
func NewSecuritySentinelAgent(client *GroqClient) *SecuritySentinelAgent {
	return &SecuritySentinelAgent{Client: client}
}

// Analyze sends the diff to the LLM and parses the structured response
func (a *SecuritySentinelAgent) Analyze(ctx context.Context, diff string) (*models.SecuritySentinelResponse, error) {
	if strings.TrimSpace(diff) == "" {
		return &models.SecuritySentinelResponse{HasVulnerability: false, Vulnerabilities: nil}, nil
	}

	userPrompt := fmt.Sprintf("Analyze this git diff:\n\n%s", diff)

	rawResponse, err := a.Client.Complete(ctx, securitySentinelSystemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("security sentinel LLM call failed: %w", err)
	}

	// Clean up potential markdown code fences the LLM might add despite instructions
	cleaned := cleanJSONResponse(rawResponse)

	var result models.SecuritySentinelResponse
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		log.Printf("WARN: failed to parse security sentinel response, raw: %s", rawResponse)
		return nil, fmt.Errorf("failed to parse security sentinel JSON: %w", err)
	}

	return &result, nil
}

// cleanJSONResponse strips markdown code fences (```json ... ```) if the LLM adds them
func cleanJSONResponse(raw string) string {
	cleaned := strings.TrimSpace(raw)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	return strings.TrimSpace(cleaned)
}
