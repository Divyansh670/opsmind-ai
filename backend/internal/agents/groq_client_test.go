package agents

import (
	"strings"
	"testing"
)

func TestCleanJSONResponse_AlreadyClean(t *testing.T) {
	input := `{"has_vulnerability": false, "vulnerabilities": []}`
	result := cleanJSONResponse(input)
	if result != input {
		t.Errorf("expected clean JSON to be unchanged, got: %s", result)
	}
}

func TestCleanJSONResponse_WithMarkdownFences(t *testing.T) {
	input := "```json\n{\"has_vulnerability\": false}\n```"
	result := cleanJSONResponse(input)
	if strings.Contains(result, "```") {
		t.Errorf("expected markdown fences to be removed, got: %s", result)
	}
	if !strings.Contains(result, "{") {
		t.Errorf("expected JSON content to remain, got: %s", result)
	}
}

func TestCleanJSONResponse_WithLeadingText(t *testing.T) {
	// cleanJSONResponse strips markdown fences but does not strip leading prose —
	// that's intentional since the LLM is prompted to return only JSON
	input := "```json\n{\"has_vulnerability\": false}\n```"
	result := cleanJSONResponse(input)
	if strings.Contains(result, "```") {
		t.Errorf("expected markdown fences to be removed, got: %s", result)
	}
	if !strings.Contains(result, "has_vulnerability") {
		t.Errorf("expected JSON content to remain after fence removal, got: %s", result)
	}
}
