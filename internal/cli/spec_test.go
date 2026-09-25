package cli

import (
	"bytes"
	"os"
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
		"--golden", "examples/spec/golden/export-reconciliation.eir.json")
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	for _, want := range []string{
		"compiled examples/spec/export-reconciliation.cdl",
		"source:     sha256:",
		"eir:        sha256:",
		"identities:",
		"golden:     examples/spec/golden/export-reconciliation.eir.json ok",
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
	golden, err := os.ReadFile(filepath.Join(specRepoRoot, "examples/spec/golden/export-reconciliation.light.prompt.md"))
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
		"--fixture", "examples/spec/fixtures/incomplete.json")
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
		"--fixture", "examples/spec/fixtures/reconciled.json",
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
