package cli

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const specRepoRoot = "../.."

func runSpecCLI(t *testing.T, args ...string) (string, string, int) {
	t.Helper()
	var out, errOut bytes.Buffer
	code := Run(args, &out, &errOut)
	return out.String(), errOut.String(), code
}

func TestSpecCompileCommand(t *testing.T) {
	out, errOut, code := runSpecCLI(t, "spec", "compile",
		"--repo", specRepoRoot,
		"--golden", "protocol/golden/protocol.eir.json")
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	for _, want := range []string{
		"compiled protocol/assembly.json",
		"source:     sha256:",
		"eir:        sha256:",
		"identities:",
		"golden:     protocol/golden/protocol.eir.json ok",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in output:\n%s", want, out)
		}
	}
}

func TestSpecRenderCommand(t *testing.T) {
	out, errOut, code := runSpecCLI(t, "spec", "render",
		"--repo", specRepoRoot,
		"--generated-at", "2026-09-25T00:00:00Z")
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	golden, err := os.ReadFile(filepath.Join(specRepoRoot, "protocol/golden/protocol.pilot.prompt.md"))
	if err != nil {
		t.Fatal(err)
	}
	if out != string(golden) {
		t.Fatalf("render output does not match golden\ngot:\n%s", out)
	}
}

func TestSpecRunCommand(t *testing.T) {
	out, errOut, code := runSpecCLI(t, "spec", "run",
		"--repo", specRepoRoot,
		"--fixture", "protocol/fixtures/incomplete.json")
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	for _, want := range []string{"reconcile.guard", "blocked", "GUARD_BLOCKED", "trace: sha256:"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in output:\n%s", want, out)
		}
	}
}

func TestSpecRunCommandJSON(t *testing.T) {
	out, errOut, code := runSpecCLI(t, "spec", "run",
		"--repo", specRepoRoot,
		"--fixture", "protocol/fixtures/reconciled.json",
		"--json")
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	if !strings.HasPrefix(strings.TrimSpace(out), `{"steps":`) {
		t.Fatalf("expected canonical trace JSON, got:\n%s", out)
	}
}

func TestSpecRunRequiresFixture(t *testing.T) {
	_, errOut, code := runSpecCLI(t, "spec", "run", "--repo", specRepoRoot)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if !strings.Contains(errOut, "requires --fixture") {
		t.Fatalf("missing fixture error: %s", errOut)
	}
}

func TestSpecRecordCommands(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	dir := t.TempDir()
	db := filepath.Join(dir, "store.db")
	payloadPath := filepath.Join(dir, "invariant.md")
	payload := "- **Constraint:** c\n- **Verification Vectors:** v1\n- **Discovered In:** SRC-1\n"
	if err := os.WriteFile(payloadPath, []byte(payload), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, errOut, code := runSpecCLI(t, "spec", "schema", "--repo", specRepoRoot, "--apply", db); code != 0 {
		t.Fatalf("schema apply exit %d: %s", code, errOut)
	}
	out, errOut, code := runSpecCLI(t, "spec", "record", "put",
		"--repo", specRepoRoot, "--db", db, "--field", "INVARIANT", "--id", "INV-CLI", "--payload", payloadPath)
	if code != 0 {
		t.Fatalf("record put exit %d: %s", code, errOut)
	}
	if !strings.Contains(out, "put INVARIANT INV-CLI") {
		t.Fatalf("record put output: %s", out)
	}
	out, errOut, code = runSpecCLI(t, "spec", "record", "list",
		"--repo", specRepoRoot, "--db", db, "--field", "INVARIANT")
	if code != 0 || strings.TrimSpace(out) != "INV-CLI" {
		t.Fatalf("record list exit %d out %q err %s", code, out, errOut)
	}
	out, errOut, code = runSpecCLI(t, "spec", "record", "get",
		"--repo", specRepoRoot, "--db", db, "--field", "INVARIANT", "--id", "INV-CLI")
	if code != 0 {
		t.Fatalf("record get exit %d: %s", code, errOut)
	}
	if strings.TrimSpace(out) != strings.TrimSpace(payload) {
		t.Fatalf("record get = %q", out)
	}
	out, errOut, code = runSpecCLI(t, "spec", "record", "inspect",
		"--repo", specRepoRoot, "--db", db, "--field", "INVARIANT")
	if code != 0 || !strings.Contains(out, "DEBUG RENDER - NON-AUTHORITATIVE") {
		t.Fatalf("record inspect exit %d out %q err %s", code, out, errOut)
	}
}
