// Package llm adapts an OpenAI-compatible chat-completions endpoint to the
// bounded work executor and reviewer interfaces. It never logs or returns the
// API key and never receives raw SQL or filesystem access.
package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Environment variable names for connection details. Secrets are never taken
// from command-line flags or committed files.
const (
	EnvURL           = "LEGACY_AUTOPSY_LLM_URL"
	EnvKey           = "LEGACY_AUTOPSY_LLM_KEY"
	EnvModel         = "LEGACY_AUTOPSY_LLM_MODEL"
	EnvMaxTokens     = "LEGACY_AUTOPSY_LLM_MAX_TOKENS"
	EnvStream        = "LEGACY_AUTOPSY_LLM_STREAM"
	EnvStreamTimeout = "LEGACY_AUTOPSY_LLM_STREAM_TIMEOUT"
	// EnvReasoning optionally sets the OpenAI-compatible reasoning_effort side
	// parameter (for example "low", "medium", "high", "xhigh").
	EnvReasoning = "LEGACY_AUTOPSY_LLM_REASONING"
)

// Delta is one streamed response fragment: user-visible content and/or hidden
// reasoning. A callback receives both so debug surfaces and a future TUI can
// watch progress without waiting for the full response.
type Delta struct {
	Content   string
	Reasoning string
}

// Config describes one OpenAI-compatible endpoint.
type Config struct {
	BaseURL string
	APIKey  string
	Model   string
	// Timeout bounds time-to-first-byte and non-streaming completions.
	Timeout   time.Duration
	MaxTokens int
	// StreamTimeout bounds a complete streamed response; long reasoning runs
	// legitimately exceed the first-byte timeout.
	StreamTimeout time.Duration
	Reasoning     string
	DisableStream bool
}

// ConfigFromEnv loads connection details from the environment and fills
// non-secret defaults. The key is required but never echoed.
func ConfigFromEnv() (Config, error) {
	cfg := Config{
		BaseURL:       strings.TrimSpace(os.Getenv(EnvURL)),
		APIKey:        strings.TrimSpace(os.Getenv(EnvKey)),
		Model:         strings.TrimSpace(os.Getenv(EnvModel)),
		Reasoning:     strings.TrimSpace(os.Getenv(EnvReasoning)),
		Timeout:       120 * time.Second,
		MaxTokens:     16384,
		StreamTimeout: 45 * time.Minute,
	}
	if raw := strings.TrimSpace(os.Getenv(EnvStreamTimeout)); raw != "" {
		parsed, err := time.ParseDuration(raw)
		if err != nil || parsed <= 0 {
			return Config{}, fmt.Errorf("llm: %s %q is not a positive duration (for example 45m)", EnvStreamTimeout, raw)
		}
		cfg.StreamTimeout = parsed
	}
	if raw := strings.TrimSpace(os.Getenv(EnvMaxTokens)); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed <= 0 {
			return Config{}, fmt.Errorf("llm: %s %q is not a positive integer", EnvMaxTokens, raw)
		}
		cfg.MaxTokens = parsed
	}
	switch strings.ToLower(strings.TrimSpace(os.Getenv(EnvStream))) {
	case "0", "false", "off", "no":
		cfg.DisableStream = true
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
	Debug  io.Writer
}

// NewClient builds a client from validated configuration.
func NewClient(config Config) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	timeout := config.Timeout
	if timeout <= 0 {
		timeout = 120 * time.Second
		config.Timeout = timeout
	}
	if config.StreamTimeout <= 0 {
		config.StreamTimeout = 45 * time.Minute
	}
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ResponseHeaderTimeout: timeout,
	}
	return &Client{config: config, http: &http.Client{Transport: transport}}, nil
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
	Stream          bool          `json:"stream,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message      chatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

type chatStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content          string `json:"content"`
			ReasoningContent string `json:"reasoning_content"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// Complete sends one system/user exchange with the configured reasoning effort.
func (c *Client) Complete(ctx context.Context, system, user string) (string, error) {
	return c.CompleteEffort(ctx, system, user, c.config.Reasoning)
}

// CompleteEffort sends one exchange with an explicit reasoning effort; an empty
// effort omits the side parameter entirely.
func (c *Client) CompleteEffort(ctx context.Context, system, user, effort string) (string, error) {
	return c.CompleteEffortStream(ctx, system, user, effort, nil)
}

// CompleteEffortStream sends one exchange and, when streamed, reports every
// delta through onDelta as it arrives. Empty-content responses are retried a
// bounded number of times before failing closed.
func (c *Client) CompleteEffortStream(ctx context.Context, system, user, effort string, onDelta func(Delta)) (string, error) {
	maxTokens := c.config.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 4096
	}
	request := chatRequest{
		Model:           c.config.Model,
		Messages:        []chatMessage{{Role: "system", Content: system}, {Role: "user", Content: user}},
		Temperature:     0,
		MaxTokens:       maxTokens,
		ReasoningEffort: effort,
		Stream:          !c.config.DisableStream,
	}
	body, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("llm: encode request: %w", err)
	}
	endpoint := strings.TrimRight(c.config.BaseURL, "/") + "/chat/completions"
	if c.Debug != nil {
		fmt.Fprintf(c.Debug, "\n--- request body ---\n%s\n--- end request ---\n", string(body))
	}
	const maxEmptyRetries = 2
	var lastFinish string
	for attempt := 0; attempt <= maxEmptyRetries; attempt++ {
		var content string
		if request.Stream {
			content, lastFinish, err = c.streamOnce(ctx, endpoint, body, onDelta)
		} else {
			content, lastFinish, err = c.completeOnce(ctx, endpoint, body)
		}
		if err != nil {
			return "", err
		}
		if strings.TrimSpace(content) != "" {
			return content, nil
		}
		if c.Debug != nil {
			fmt.Fprintf(c.Debug, "\n--- empty provider response (attempt %d, finish_reason=%q) ---\n", attempt+1, lastFinish)
		}
	}
	return "", fmt.Errorf("llm: empty response after retries (finish_reason=%q); raise %s or lower reasoning effort", lastFinish, EnvMaxTokens)
}

func (c *Client) completeOnce(ctx context.Context, endpoint string, body []byte) (string, string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()
	response, err := c.do(ctx, endpoint, body)
	if err != nil {
		return "", "", err
	}
	defer response.Body.Close()
	payload, err := io.ReadAll(io.LimitReader(response.Body, 8<<20))
	if err != nil {
		return "", "", fmt.Errorf("llm: read response: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("llm: endpoint returned HTTP %d: %s", response.StatusCode, providerMessage(payload))
	}
	var decoded chatResponse
	if err := json.Unmarshal(payload, &decoded); err != nil {
		return "", "", fmt.Errorf("llm: decode response: %w", err)
	}
	if decoded.Error != nil {
		return "", "", fmt.Errorf("llm: provider error: %s", decoded.Error.Message)
	}
	if len(decoded.Choices) == 0 {
		return "", "", fmt.Errorf("llm: response contained no choices")
	}
	return decoded.Choices[0].Message.Content, decoded.Choices[0].FinishReason, nil
}

func (c *Client) streamOnce(ctx context.Context, endpoint string, body []byte, onDelta func(Delta)) (string, string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.StreamTimeout)
	defer cancel()
	response, err := c.do(ctx, endpoint, body)
	if err != nil {
		return "", "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		payload, _ := io.ReadAll(io.LimitReader(response.Body, 1<<20))
		return "", "", fmt.Errorf("llm: endpoint returned HTTP %d: %s", response.StatusCode, providerMessage(payload))
	}
	scanner := bufio.NewScanner(response.Body)
	scanner.Buffer(make([]byte, 64*1024), 8<<20)
	var content, reasoning strings.Builder
	finish := ""
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, ":") || !strings.HasPrefix(line, "data:") {
			continue
		}
		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "[DONE]" {
			break
		}
		var chunk chatStreamChunk
		if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
			return "", "", fmt.Errorf("llm: decode stream chunk: %w", err)
		}
		if chunk.Error != nil {
			return "", "", fmt.Errorf("llm: provider error: %s", chunk.Error.Message)
		}
		for _, choice := range chunk.Choices {
			if choice.Delta.ReasoningContent != "" {
				reasoning.WriteString(choice.Delta.ReasoningContent)
				if onDelta != nil {
					onDelta(Delta{Reasoning: choice.Delta.ReasoningContent})
				}
			}
			if choice.Delta.Content != "" {
				content.WriteString(choice.Delta.Content)
				if onDelta != nil {
					onDelta(Delta{Content: choice.Delta.Content})
				}
			}
			if choice.FinishReason != "" {
				finish = choice.FinishReason
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", "", fmt.Errorf("llm: read stream: %w", err)
	}
	return content.String(), finish, nil
}

func (c *Client) do(ctx context.Context, endpoint string, body []byte) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("llm: build request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "text/event-stream")
	request.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	response, err := c.http.Do(request)
	if err != nil {
		return nil, fmt.Errorf("llm: request to %s failed: %w", endpoint, err)
	}
	return response, nil
}

func providerMessage(payload []byte) string {
	var decoded struct {
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(payload, &decoded); err == nil && decoded.Error != nil {
		return decoded.Error.Message
	}
	text := strings.TrimSpace(string(payload))
	if len(text) > 200 {
		text = text[:200] + "..."
	}
	return text
}
