package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Divyansh670/opsmind-ai/backend/internal/models"
)

const costPredictorSystemPrompt = `You are a FinOps & Cloud Infrastructure Architect analyzing git diffs for cost impact.
Examine the provided code diff for changes to cloud infrastructure such as:
- Terraform (.tf) resource definitions
- CloudFormation templates
- Serverless framework configs (serverless.yml)
- Kubernetes resource limits/requests
- Database instance sizing
- Auto-scaling group configurations
- Storage provisioning (S3, EBS, etc.)

Estimate the monthly dollar cost impact of any infrastructure changes you find.
Consider instance type changes, added/removed resources, scaling changes, and storage size changes.

You must respond exclusively with a valid JSON object matching this exact schema, with no additional text, no markdown formatting, and no code fences:
{
  "has_drift": boolean,
  "drift_usd": number,
  "drift_explanation": "string",
  "affected_services": ["string"]
}

If no infrastructure changes are found, return {"has_drift": false, "drift_usd": 0, "drift_explanation": "No infrastructure changes detected", "affected_services": []}
A positive drift_usd means cost increase, negative means cost decrease.`

// CostPredictorAgent wraps the Groq client for cost analysis
type CostPredictorAgent struct {
	Client *GroqClient
}

// NewCostPredictorAgent creates a new Cost Predictor agent
func NewCostPredictorAgent(client *GroqClient) *CostPredictorAgent {
	return &CostPredictorAgent{Client: client}
}

// Analyze sends the diff to the LLM and parses the structured cost response
func (a *CostPredictorAgent) Analyze(ctx context.Context, diff string) (*models.CostPredictorResponse, error) {
	if strings.TrimSpace(diff) == "" {
		return &models.CostPredictorResponse{HasDrift: false, DriftUSD: 0}, nil
	}

	userPrompt := fmt.Sprintf("Analyze this git diff for infrastructure cost impact:\n\n%s", diff)

	rawResponse, err := a.Client.Complete(ctx, costPredictorSystemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("cost predictor LLM call failed: %w", err)
	}

	cleaned := cleanJSONResponse(rawResponse)

	var result models.CostPredictorResponse
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		log.Printf("WARN: failed to parse cost predictor response, raw: %s", rawResponse)
		return nil, fmt.Errorf("failed to parse cost predictor JSON: %w", err)
	}

	return &result, nil
}
