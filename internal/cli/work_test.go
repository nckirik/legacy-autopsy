package cli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/llm"
)

func fakeModelServer(t *testing.T) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		system := ""
		for _, message := range request.Messages {
			if message.Role == "system" {
				system = message.Content
			}
		}
		content := "OK"
		if !strings.Contains(system, "connectivity probe") {
			content = "RESULT: PROPOSAL\nEFFECTS: none\n---RECORD---\n- **Constraint:** model writes\n- **Verification Vectors:** v1\n- **Discovered In:** SRC-1\n---END---\n"
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{{"message": map[string]string{"role": "assistant", "content": content}}},
		})
	}))
	t.Cleanup(server.Close)
	return server
}

func setLLMEnv(t *testing.T, server *httptest.Server) {
	t.Helper()
	t.Setenv(llm.EnvURL, server.URL+"/v1")
	t.Setenv(llm.EnvKey, "test-key")
	t.Setenv(llm.EnvModel, "test-model")
}

func TestSpecWorkPingWithFakeModel(t *testing.T) {
	server := fakeModelServer(t)
	setLLMEnv(t, server)
	out, errOut, code := runSpecCLI(t, "spec", "work", "ping", "--repo", specRepoRoot)
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	for _, marker := range []string{"model:    test-model", "reply:    OK", "endpoint:"} {
		if !strings.Contains(out, marker) {
			t.Fatalf("ping output missing %q:\n%s", marker, out)
		}
	}
	if strings.Contains(out, "test-key") {
		t.Fatal("ping output leaks the API key")
	}
}

func TestSpecWorkRunWithFakeModel(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	server := fakeModelServer(t)
	setLLMEnv(t, server)
	dir := t.TempDir()
	db := filepath.Join(dir, "store.db")
	unitsPath := filepath.Join(dir, "units.json")
	units := `[{"id":"U-1","mode":"Discovery","scope":"demo","record_id":"INV-CLI-LLM","result_type":"INVARIANT"}]`
	if err := os.WriteFile(unitsPath, []byte(units), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, errOut, code := runSpecCLI(t, "spec", "schema", "--repo", specRepoRoot, "--apply", db); code != 0 {
		t.Fatalf("schema apply exit %d: %s", code, errOut)
	}
	out, errOut, code := runSpecCLI(t, "spec", "work", "run",
		"--repo", specRepoRoot, "--db", db, "--units", unitsPath, "--review-mode", "critical")
	if code != 0 {
		t.Fatalf("work run exit %d: %s", code, errOut)
	}
	if !strings.Contains(out, "committed 1 unit") || !strings.Contains(out, "U-1") {
		t.Fatalf("work run output:\n%s", out)
	}
	getOut, errOut, code := runSpecCLI(t, "spec", "record", "get",
		"--repo", specRepoRoot, "--db", db, "--field", "INVARIANT", "--id", "INV-CLI-LLM")
	if code != 0 {
		t.Fatalf("record get exit %d: %s", code, errOut)
	}
	if !strings.Contains(getOut, "model writes") {
		t.Fatalf("committed record payload:\n%s", getOut)
	}
}
