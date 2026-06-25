package agents

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const groqAPIURL = "https://api.groq.com/openai/v1/chat/completions"

// GroqClient wraps calls to the Groq API
type GroqClient struct {
	APIKey     string
	Model      string
	HTTPClient *http.Client
}

// NewGroqClient creates a new Groq API client
func NewGroqClient(apiKey, model string) *GroqClient {
	return &GroqClient{
		APIKey: apiKey,
		Model:  model,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// groqMessage represents a single chat message
type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// groqRequest is the request body sent to Groq
type groqRequest struct {
	Model       string        `json:"model"`
	Messages    []groqMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens"`
}

// groqResponse is the response body from Groq
type groqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Complete sends a system + user prompt to Groq and returns the raw text response
func (c *GroqClient) Complete(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	reqBody := groqRequest{
		Model: c.Model,
		Messages: []groqMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.1, // low temperature for consistent, deterministic JSON output
		MaxTokens:   2048,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, groqAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Groq API: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var groqResp groqResponse
	if err := json.Unmarshal(bodyBytes, &groqResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if groqResp.Error != nil {
		return "", fmt.Errorf("groq API error: %s", groqResp.Error.Message)
	}

	if len(groqResp.Choices) == 0 {
		return "", fmt.Errorf("groq API returned no choices")
	}

	return groqResp.Choices[0].Message.Content, nil
}

// groqStreamRequest is the request body for streaming calls
type groqStreamRequest struct {
	Model       string        `json:"model"`
	Messages    []groqMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens"`
	Stream      bool          `json:"stream"`
}

// groqStreamChunk is a single SSE data chunk from Groq
type groqStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

// CompleteStream sends a prompt to Groq and streams the response token by token.
// It calls onToken for each token received, and returns when the stream is done.
func (c *GroqClient) CompleteStream(ctx context.Context, systemPrompt, userPrompt string, onToken func(string)) error {
	reqBody := groqStreamRequest{
		Model: c.Model,
		Messages: []groqMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.3,
		MaxTokens:   1024,
		Stream:      true,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal stream request: %w", err)
	}

	// Use a client without timeout for streaming
	streamClient := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, groqAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create stream request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "text/event-stream")

	resp, err := streamClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to start stream: %w", err)
	}
	defer resp.Body.Close()

	buf := make([]byte, 4096)
	var leftover string

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		n, err := resp.Body.Read(buf)
		if n > 0 {
			chunk := leftover + string(buf[:n])
			lines := strings.Split(chunk, "\n")
			leftover = ""

			for i, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" || line == "data: [DONE]" {
					continue
				}
				if !strings.HasPrefix(line, "data: ") {
					if i == len(lines)-1 {
						leftover = line
					}
					continue
				}

				jsonData := strings.TrimPrefix(line, "data: ")
				var streamChunk groqStreamChunk
				if jsonErr := json.Unmarshal([]byte(jsonData), &streamChunk); jsonErr != nil {
					continue
				}

				if len(streamChunk.Choices) > 0 {
					token := streamChunk.Choices[0].Delta.Content
					if token != "" {
						onToken(token)
					}
					if streamChunk.Choices[0].FinishReason != nil {
						return nil
					}
				}
			}
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("stream read error: %w", err)
		}
	}
}
