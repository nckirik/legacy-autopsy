package llm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nckirik/legacy-autopsy/internal/work"
)

type captured struct {
	path      string
	auth      string
	model     string
	reasoning string
	temp      float64
	stream    bool
	systems   []string
	users     []string
}

func fakeServer(t *testing.T, reply func(system, user string) string) (*httptest.Server, *captured) {
	t.Helper()
	recorded := &captured{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorded.path = r.URL.Path
		recorded.auth = r.Header.Get("Authorization")
		var request chatRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("decode request: %v", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		recorded.model = request.Model
		recorded.reasoning = request.ReasoningEffort
		recorded.temp = request.Temperature
		recorded.stream = request.Stream
		for _, message := range request.Messages {
			if message.Role == "system" {
				recorded.systems = append(recorded.systems, message.Content)
			} else {
				recorded.users = append(recorded.users, message.Content)
			}
		}
		content := reply(last(recorded.systems), last(recorded.users))
		if request.Stream {
			w.Header().Set("Content-Type", "text/event-stream")
			writeSSE(w, map[string]any{"choices": []map[string]any{{"delta": map[string]string{"reasoning_content": "think "}}}})
			half := len(content) / 2
			writeSSE(w, map[string]any{"choices": []map[string]any{{"delta": map[string]string{"content": content[:half]}}}})
			writeSSE(w, map[string]any{"choices": []map[string]any{{"delta": map[string]string{"content": content[half:]}, "finish_reason": "stop"}}})
			fmt.Fprint(w, "data: [DONE]\n\n")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{{"message": map[string]string{"role": "assistant", "content": content}}},
		})
	}))
	t.Cleanup(server.Close)
	return server, recorded
}

func writeSSE(w http.ResponseWriter, chunk map[string]any) {
	payload, _ := json.Marshal(chunk)
	fmt.Fprintf(w, "data: %s\n\n", payload)
}

func last(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[len(values)-1]
}

func testClient(t *testing.T, baseURL string) *Client {
	t.Helper()
	client, err := NewClient(Config{BaseURL: baseURL + "/v1", APIKey: "test-key", Model: "test-model", Reasoning: "high", Timeout: 5 * time.Second, MaxTokens: 128})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestClientComplete(t *testing.T) {
	server, recorded := fakeServer(t, func(system, user string) string { return "hello " + user })
	client := testClient(t, server.URL)
	text, err := client.Complete(t.Context(), "system prompt", "user prompt")
	if err != nil {
		t.Fatal(err)
	}
	if text != "hello user prompt" {
		t.Fatalf("text = %q", text)
	}
	if recorded.path != "/v1/chat/completions" {
		t.Fatalf("path = %q", recorded.path)
	}
	if recorded.auth != "Bearer test-key" {
		t.Fatalf("auth = %q", recorded.auth)
	}
	if recorded.model != "test-model" || recorded.temp != 0 || recorded.reasoning != "high" {
		t.Fatalf("model=%q temp=%v reasoning=%q", recorded.model, recorded.temp, recorded.reasoning)
	}
}

func TestClientErrorsDoNotLeakKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusInternalServerError)
	}))
	defer server.Close()
	client := testClient(t, server.URL)
	_, err := client.Complete(t.Context(), "s", "u")
	if err == nil {
		t.Fatal("HTTP 500 accepted")
	}
	if strings.Contains(err.Error(), "test-key") {
		t.Fatalf("error leaks the API key: %v", err)
	}
}

func TestConfigFromEnv(t *testing.T) {
	t.Setenv(EnvURL, "")
	t.Setenv(EnvKey, "")
	t.Setenv(EnvModel, "")
	if _, err := ConfigFromEnv(); err == nil {
		t.Fatal("missing configuration accepted")
	}
	t.Setenv(EnvURL, "https://example.test/v1")
	t.Setenv(EnvKey, "secret-value")
	t.Setenv(EnvModel, "model-x")
	t.Setenv(EnvReasoning, "xhigh")
	config, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if config.Model != "model-x" || config.Reasoning != "xhigh" || config.MaxTokens != 16384 {
		t.Fatalf("config = %+v", config)
	}
	if config.ReviewReasoning != "medium" {
		t.Fatalf("reviewer default = %q, want medium", config.ReviewReasoning)
	}
	t.Setenv(EnvReviewReasoning, "low")
	if overridden, err := ConfigFromEnv(); err != nil || overridden.ReviewReasoning != "low" {
		t.Fatalf("reviewer override = %+v err = %v", overridden, err)
	}
	t.Setenv(EnvReviewReasoning, "")
	t.Setenv(EnvReasoning, "bad\nvalue")
	if _, err := ConfigFromEnv(); err == nil {
		t.Fatal("line break in reasoning accepted")
	}
	t.Setenv(EnvReasoning, "")
	if err := (Config{BaseURL: "x", APIKey: "super-secret"}).Validate(); err == nil || strings.Contains(err.Error(), "super-secret") {
		t.Fatalf("validation error leaks or misses: %v", err)
	}
}

func TestParseProposal(t *testing.T) {
	text := `RESULT: PROPOSAL
EFFECTS: none
---RECORD---
- **Constraint:** totals match
- **Verification Vectors:** v1
- **Discovered In:** SRC-1
---END---
`
	proposal, err := ParseProposal(text)
	if err != nil {
		t.Fatal(err)
	}
	if len(proposal.Effects) != 0 || proposal.FieldID != "" || len(proposal.Records) != 1 {
		t.Fatalf("proposal = %+v", proposal)
	}
	if !strings.Contains(string(proposal.Records[0].Payload), "**Constraint:**") {
		t.Fatalf("payload = %q", proposal.Records[0].Payload)
	}
	zero, err := ParseProposal("RESULT: ZERO-DOMAIN\nREASON: no schedules found in the window\nEVIDENCE: SRC-1, SRC-2\n")
	if err != nil || zero.ZeroDomain == nil || len(zero.ZeroDomain.Evidence) != 2 {
		t.Fatalf("zero-domain parse = %+v err = %v", zero, err)
	}
	blocker, err := ParseProposal("RESULT: BLOCKER\nKIND: uncertainty\nMESSAGE: not enough evidence\n")
	if err == nil {
		t.Fatal("blocker parsed as a proposal")
	}
	var blockerErr *work.BlockerError
	if !asBlocker(err, &blockerErr) || blockerErr.Blocker.Kind != "uncertainty" {
		t.Fatalf("blocker error = %v", err)
	}
	_ = blocker
	for _, bad := range []string{
		"no markers at all",
		"RESULT: PROPOSAL\nno record block\n",
		"RESULT: SOMETHING\n---RECORD---\n- **A:** b\n---END---\n",
		"RESULT: PROPOSAL\n---RECORD---\n- **A:** b\n",
	} {
		if _, err := ParseProposal(bad); err == nil {
			t.Fatalf("invalid response accepted: %q", bad)
		}
	}
}

func asBlocker(err error, target **work.BlockerError) bool {
	blocker, ok := err.(*work.BlockerError)
	if ok {
		*target = blocker
	}
	return ok
}

func TestParseReview(t *testing.T) {
	result, err := ParseReview("OUTCOME: rejected\nFINDINGS:\n- missing evidence\n- wrong scope\n")
	if err != nil {
		t.Fatal(err)
	}
	if result.Outcome != work.ReviewRejected || len(result.Findings) != 2 {
		t.Fatalf("result = %+v", result)
	}
	result, err = ParseReview("OUTCOME: approved\nFINDINGS:\n- none\n")
	if err != nil || result.Outcome != work.ReviewApproved || len(result.Findings) != 0 {
		t.Fatalf("result = %+v err = %v", result, err)
	}
	if _, err := ParseReview("looks fine to me"); err == nil {
		t.Fatal("unstructured review accepted")
	}
}

func TestExecutorAndReviewerThroughFakeModel(t *testing.T) {
	server, recorded := fakeServer(t, func(system, user string) string {
		if strings.Contains(system, "independent reviewer") {
			return "OUTCOME: approved\nFINDINGS:\n- none\n"
		}
		return "RESULT: PROPOSAL\nEFFECTS: none\n---RECORD---\n- **Constraint:** model writes\n- **Verification Vectors:** v1\n- **Discovered In:** SRC-1\n---END---\n"
	})
	client := testClient(t, server.URL)
	executor := NewExecutor(client)
	reviewer := NewReviewer(client)

	context := work.Context{UnitID: "U", Bytes: []byte("# Work unit U\n")}
	proposal, err := executor.Execute(work.Invocation{
		Unit:    work.Unit{ID: "U", Mode: "Discovery", ResultType: "INVARIANT"},
		Context: context,
	})
	if err != nil {
		t.Fatal(err)
	}
	if proposal.FieldID != "INVARIANT" {
		t.Fatalf("proposal = %+v", proposal)
	}
	if !strings.Contains(last(recorded.systems), "bounded semantic executor") {
		t.Fatalf("executor system prompt missing: %q", last(recorded.systems))
	}
	result, err := reviewer.Review(work.ReviewRequest{
		Unit:     work.Unit{ID: "U", ResultType: "INVARIANT"},
		Targets:  []work.Ref{{FieldID: "INVARIANT", RecordID: "INV-LLM"}},
		Payloads: [][]byte{proposal.Records[0].Payload}, Context: context,
	})
	if err != nil || result.Outcome != work.ReviewApproved {
		t.Fatalf("review = %+v err = %v", result, err)
	}
	if !strings.Contains(last(recorded.systems), "independent reviewer") {
		t.Fatalf("reviewer system prompt missing: %q", last(recorded.systems))
	}
}

func TestClientRetriesEmptyContent(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.Header().Set("Content-Type", "text/event-stream")
		if attempts >= 3 {
			writeSSE(w, map[string]any{"choices": []map[string]any{{"delta": map[string]string{"content": "OK"}, "finish_reason": "stop"}}})
		}
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer server.Close()
	client := testClient(t, server.URL)
	text, err := client.Complete(t.Context(), "s", "u")
	if err != nil {
		t.Fatal(err)
	}
	if text != "OK" || attempts != 3 {
		t.Fatalf("text=%q attempts=%d", text, attempts)
	}
}

func TestClientStreamsDeltas(t *testing.T) {
	server, recorded := fakeServer(t, func(system, user string) string { return "streamed reply" })
	client := testClient(t, server.URL)
	var content, reasoning strings.Builder
	text, err := client.CompleteEffortStream(t.Context(), "s", "u", "xhigh", func(delta Delta) {
		content.WriteString(delta.Content)
		reasoning.WriteString(delta.Reasoning)
	})
	if err != nil {
		t.Fatal(err)
	}
	if text != "streamed reply" || content.String() != "streamed reply" {
		t.Fatalf("text=%q deltas=%q", text, content.String())
	}
	if !strings.Contains(reasoning.String(), "think") {
		t.Fatalf("reasoning deltas = %q", reasoning.String())
	}
	if !recorded.stream {
		t.Fatal("client did not request streaming")
	}
}

func TestEffortResolution(t *testing.T) {
	client := &Client{config: Config{Reasoning: "xhigh"}}
	executor := &Executor{client: client}
	if got := executor.effortFor(work.Unit{}); got != "xhigh" {
		t.Fatalf("default = %q", got)
	}
	if got := executor.effortFor(work.Unit{ReasoningEffort: "off"}); got != "" {
		t.Fatalf("unit off = %q", got)
	}
	if got := executor.effortFor(work.Unit{ReasoningEffort: "medium"}); got != "medium" {
		t.Fatalf("unit medium = %q", got)
	}
	force := "high"
	executor.ForceReasoning = &force
	if got := executor.effortFor(work.Unit{ReasoningEffort: "off"}); got != "high" {
		t.Fatalf("force = %q", got)
	}
	off := "none"
	executor.ForceReasoning = &off
	if got := executor.effortFor(work.Unit{ReasoningEffort: "high"}); got != "" {
		t.Fatalf("force off = %q", got)
	}
}
