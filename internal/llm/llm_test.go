package llm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nckirik/legacy-autopsy/internal/work"
)

type captured struct {
	path    string
	auth    string
	model   string
	temp    float64
	systems []string
	users   []string
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
		recorded.temp = request.Temperature
		for _, message := range request.Messages {
			if message.Role == "system" {
				recorded.systems = append(recorded.systems, message.Content)
			} else {
				recorded.users = append(recorded.users, message.Content)
			}
		}
		content := reply(last(recorded.systems), last(recorded.users))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{{"message": map[string]string{"role": "assistant", "content": content}}},
		})
	}))
	t.Cleanup(server.Close)
	return server, recorded
}

func last(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[len(values)-1]
}

func testClient(t *testing.T, baseURL string) *Client {
	t.Helper()
	client, err := NewClient(Config{BaseURL: baseURL + "/v1", APIKey: "test-key", Model: "test-model", Timeout: 5 * time.Second, MaxTokens: 128})
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
	if recorded.model != "test-model" || recorded.temp != 0 {
		t.Fatalf("model=%q temp=%v", recorded.model, recorded.temp)
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
	config, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if config.Model != "model-x" || config.MaxTokens != 4096 {
		t.Fatalf("config = %+v", config)
	}
	if err := (Config{BaseURL: "x", APIKey: "super-secret"}).Validate(); err == nil || strings.Contains(err.Error(), "super-secret") {
		t.Fatalf("validation error leaks or misses: %v", err)
	}
}

func TestParseProposal(t *testing.T) {
	text := `RESULT: PROPOSAL
RECORD-ID: INV-0001
VERSION: 2
EFFECTS: none
---PAYLOAD---
- **Constraint:** totals match
- **Verification Vectors:** v1
- **Discovered In:** SRC-1
---END---
`
	proposal, err := ParseProposal(text)
	if err != nil {
		t.Fatal(err)
	}
	if proposal.RecordID != "INV-0001" || proposal.Version != 2 || len(proposal.Effects) != 0 {
		t.Fatalf("proposal = %+v", proposal)
	}
	if !strings.Contains(string(proposal.Payload), "**Constraint:**") {
		t.Fatalf("payload = %q", proposal.Payload)
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
		"RESULT: PROPOSAL\nRECORD-ID: X\n---PAYLOAD---\n- **A:** b\n---END---\n",
		"RESULT: PROPOSAL\nRECORD-ID: X\nVERSION: zero\n---PAYLOAD---\n- **A:** b\n---END---\n",
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
		return "RESULT: PROPOSAL\nRECORD-ID: INV-LLM\nVERSION: 1\nEFFECTS: none\n---PAYLOAD---\n- **Constraint:** model writes\n- **Verification Vectors:** v1\n- **Discovered In:** SRC-1\n---END---\n"
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
	if proposal.FieldID != "INVARIANT" || proposal.RecordID != "INV-LLM" {
		t.Fatalf("proposal = %+v", proposal)
	}
	if !strings.Contains(last(recorded.systems), "bounded semantic executor") {
		t.Fatalf("executor system prompt missing: %q", last(recorded.systems))
	}
	result, err := reviewer.Review(work.ReviewRequest{
		Unit: work.Unit{ID: "U", ResultType: "INVARIANT"}, Target: work.Ref{FieldID: "INVARIANT", RecordID: "INV-LLM"},
		Payload: proposal.Payload, Context: context,
	})
	if err != nil || result.Outcome != work.ReviewApproved {
		t.Fatalf("review = %+v err = %v", result, err)
	}
	if !strings.Contains(last(recorded.systems), "independent reviewer") {
		t.Fatalf("reviewer system prompt missing: %q", last(recorded.systems))
	}
}
