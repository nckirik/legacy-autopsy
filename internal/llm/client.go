// Package llm adapts an OpenAI-compatible chat-completions endpoint to the
// bounded work executor and reviewer interfaces. It never logs or returns the
// API key and never receives raw SQL or filesystem access.
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Environment variable names for connection details. Secrets are never taken
// from command-line flags or committed files.
const (
	EnvURL   = "LEGACY_AUTOPSY_LLM_URL"
	EnvKey   = "LEGACY_AUTOPSY_LLM_KEY"
	EnvModel = "LEGACY_AUTOPSY_LLM_MODEL"
	// EnvReasoning optionally sets the OpenAI-compatible reasoning_effort side
	// parameter (for example "low", "medium", "high", "xhigh").
	EnvReasoning = "LEGACY_AUTOPSY_LLM_REASONING"
)

// Config describes one OpenAI-compatible endpoint.
type Config struct {
	BaseURL   string
	APIKey    string
	Model     string
	Reasoning string
	Timeout   time.Duration
	MaxTokens int
}

// ConfigFromEnv loads connection details from the environment and fills
// non-secret defaults. The key is required but never echoed.
func ConfigFromEnv() (Config, error) {
	cfg := Config{
		BaseURL:   strings.TrimSpace(os.Getenv(EnvURL)),
		APIKey:    strings.TrimSpace(os.Getenv(EnvKey)),
		Model:     strings.TrimSpace(os.Getenv(EnvModel)),
		Reasoning: strings.TrimSpace(os.Getenv(EnvReasoning)),
		Timeout:   120 * time.Second,
		MaxTokens: 4096,
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Validate fails closed on missing connection details without printing the key.
func (c Config) Validate() error {
	var missing []string
	if c.BaseURL == "" {
		missing = append(missing, EnvURL)
	}
	if c.APIKey == "" {
		missing = append(missing, EnvKey)
	}
	if c.Model == "" {
		missing = append(missing, EnvModel)
	}
	if strings.ContainsAny(c.Reasoning, "\r\n") {
		return fmt.Errorf("llm: reasoning effort contains line breaks")
	}
	if len(missing) > 0 {
		return fmt.Errorf("llm: missing configuration: %s", strings.Join(missing, ", "))
	}
	return nil
}

// Client is a minimal chat-completions client.
type Client struct {
	config Config
	http   *http.Client
}

// NewClient builds a client from validated configuration.
func NewClient(config Config) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	timeout := config.Timeout
	if timeout <= 0 {
		timeout = 120 * time.Second
	}
	return &Client{config: config, http: &http.Client{Timeout: timeout}}, nil
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model           string        `json:"model"`
	Messages        []chatMessage `json:"messages"`
	Temperature     float64       `json:"temperature"`
	MaxTokens       int           `json:"max_tokens,omitempty"`
	ReasoningEffort string        `json:"reasoning_effort,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// Complete sends one system/user exchange and returns the assistant text. Error
// messages never include the API key or request body.
func (c *Client) Complete(ctx context.Context, system, user string) (string, error) {
	maxTokens := c.config.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 4096
	}
	body, err := json.Marshal(chatRequest{
		Model:           c.config.Model,
		Messages:        []chatMessage{{Role: "system", Content: system}, {Role: "user", Content: user}},
		Temperature:     0,
		MaxTokens:       maxTokens,
		ReasoningEffort: c.config.Reasoning,
	})
	if err != nil {
		return "", fmt.Errorf("llm: encode request: %w", err)
	}
	endpoint := strings.TrimRight(c.config.BaseURL, "/") + "/chat/completions"
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("llm: build request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	response, err := c.http.Do(request)
	if err != nil {
		return "", fmt.Errorf("llm: request to %s failed: %w", endpoint, err)
	}
	defer response.Body.Close()
	payload, err := io.ReadAll(io.LimitReader(response.Body, 8<<20))
	if err != nil {
		return "", fmt.Errorf("llm: read response: %w", err)
	}
	var decoded chatResponse
	if err := json.Unmarshal(payload, &decoded); err != nil {
		return "", fmt.Errorf("llm: decode response: %w", err)
	}
	if decoded.Error != nil {
		return "", fmt.Errorf("llm: provider error: %s", decoded.Error.Message)
	}
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("llm: endpoint returned HTTP %d", response.StatusCode)
	}
	if len(decoded.Choices) == 0 {
		return "", fmt.Errorf("llm: response contained no choices")
	}
	return decoded.Choices[0].Message.Content, nil
}
