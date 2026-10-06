package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/artifact"
)

func TestSpecArtifactValidateAndVerify(t *testing.T) {
	dir := t.TempDir()
	payloadPath := filepath.Join(dir, "envelope.md")
	if err := os.WriteFile(payloadPath, []byte("- **Required Human Signatories / Signature Timestamps:** Alice\n- **Envelope Audit:** reviewed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, errOut, code := runSpecCLI(t, "spec", "artifact", "validate",
		"--repo", specRepoRoot, "--field", "OUTER-SIGNATURE-ENVELOPE", "--payload", payloadPath,
		"--type", "OUTER-BUNDLE-MANIFEST", "--path", "handoff/27.md")
	if code != 0 {
		t.Fatalf("validate exit %d: %s", code, errOut)
	}
	if !strings.Contains(out, "valid OUTER-SIGNATURE-ENVELOPE payload") || !strings.Contains(out, "payload fingerprint:") {
		t.Fatalf("validate output: %s", out)
	}

	bundle := filepath.Join(dir, "bundle")
	if err := os.MkdirAll(filepath.Join(bundle, "handbook"), 0o755); err != nil {
		t.Fatal(err)
	}
	snapshot := "- **System / Protocol / Snapshot:** demo | 4.2.1 | dev | snap-1\n"
	if err := os.WriteFile(filepath.Join(bundle, "handbook", "01.md"), []byte(snapshot), 0o644); err != nil {
		t.Fatal(err)
	}
	others := []artifact.Member{{Path: "handbook/01.md", Bytes: []byte(snapshot)}}
	manifest, err := artifact.PackageMembers(others)
	if err != nil {
		t.Fatal(err)
	}
	var b strings.Builder
	b.WriteString("- **System / Protocol / Snapshot:** demo | 4.2.1 | dev | snap-1\n")
	b.WriteString("- **Package Members:**\n")
	b.WriteString("  | Artifact Type | Normalized Path | Package-Member File Fingerprint |\n")
	b.WriteString("  | :------------ | :-------------- | :------------------------------ |\n")
	for _, entry := range manifest {
		fmt.Fprintf(&b, "  | OUTER | %s | %s |\n", entry.Path, entry.Fingerprint)
	}
	b.WriteString("- **EXIT-E-STATUS:** Passed\n")
	outerPath := filepath.Join(bundle, "27.md")
	if err := os.WriteFile(outerPath, []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	out, errOut, code = runSpecCLI(t, "spec", "artifact", "verify",
		"--repo", specRepoRoot, "--type", "OUTER-BUNDLE-MANIFEST", "--path", "27.md", "--dir", bundle)
	if code != 0 {
		t.Fatalf("verify exit %d: %s", code, errOut)
	}
	if !strings.Contains(out, "bundle verified: 2 member(s)") {
		t.Fatalf("verify output: %s", out)
	}
	if err := os.WriteFile(filepath.Join(bundle, "handbook", "01.md"), []byte(snapshot+"x"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, _, code = runSpecCLI(t, "spec", "artifact", "verify",
		"--repo", specRepoRoot, "--type", "OUTER-BUNDLE-MANIFEST", "--path", "27.md", "--dir", bundle)
	if code == 0 {
		t.Fatal("tampered bundle verified")
	}
}

func TestSpecPreflightScan(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module demo\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, errOut, code := runSpecCLI(t, "spec", "preflight", "scan", "--repo", specRepoRoot, "--target", dir)
	if code != 0 {
		t.Fatalf("scan exit %d: %s", code, errOut)
	}
	for _, marker := range []string{"# Repository scan", "- **Languages:** Go", "- **Entry Points:**", "fingerprint: sha256:"} {
		if !strings.Contains(out, marker) {
			t.Fatalf("scan output missing %q:\n%s", marker, out)
		}
	}
}

func TestSpecDiscoveryRequiresPreflight(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	dir := t.TempDir()
	db := filepath.Join(dir, "store.db")
	if _, errOut, code := runSpecCLI(t, "spec", "schema", "--repo", specRepoRoot, "--apply", db); code != 0 {
		t.Fatalf("schema apply exit %d: %s", code, errOut)
	}
	_, errOut, code := runSpecCLI(t, "spec", "discovery", "run",
		"--repo", specRepoRoot, "--target", dir, "--db", db)
	if code == 0 || !strings.Contains(errOut, "requires a completed preflight") {
		t.Fatalf("exit %d err %s", code, errOut)
	}
}

func TestSpecAssuranceRequiresInventory(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	dir := t.TempDir()
	db := filepath.Join(dir, "store.db")
	if _, errOut, code := runSpecCLI(t, "spec", "schema", "--repo", specRepoRoot, "--apply", db); code != 0 {
		t.Fatalf("schema apply exit %d: %s", code, errOut)
	}
	_, errOut, code := runSpecCLI(t, "spec", "assurance", "claims",
		"--repo", specRepoRoot, "--target", dir, "--db", db)
	if code == 0 || !strings.Contains(errOut, "requires a completed preflight") {
		t.Fatalf("exit %d err %s", code, errOut)
	}
}
