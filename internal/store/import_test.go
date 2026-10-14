package store

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestImportRoundTrip(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	db, res := ingestStore(t)
	payload := "- **Constraint:** a\n- **Verification Vectors:** v1\n- **Discovered In:** SRC-1\n"
	if err := Ingest(db, res.EIR, "INVARIANT", "INV-A", 1, []byte(payload)); err != nil {
		t.Fatal(err)
	}
	prov := Provenance{Generator: "store/test", SourceFingerprint: "sha256:test"}
	dump, fingerprint, err := CanonicalDump(db, res.EIR, prov)
	if err != nil {
		t.Fatal(err)
	}
	fresh := filepath.Join(t.TempDir(), "fresh.db")
	got, err := Import(fresh, res.EIR, dump, fingerprint)
	if err != nil {
		t.Fatal(err)
	}
	if got != fingerprint {
		t.Fatalf("import fingerprint = %s, want %s", got, fingerprint)
	}
	again, againFingerprint, err := CanonicalDump(fresh, res.EIR, prov)
	if err != nil {
		t.Fatal(err)
	}
	if againFingerprint != fingerprint {
		t.Fatalf("round-trip fingerprint = %s, want %s", againFingerprint, fingerprint)
	}
	if string(again) != string(dump) {
		t.Fatal("round-trip dump is not byte-identical")
	}
}

func TestImportRejects(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	db, res := ingestStore(t)
	payload := "- **Constraint:** a\n- **Verification Vectors:** v1\n- **Discovered In:** SRC-1\n"
	if err := Ingest(db, res.EIR, "INVARIANT", "INV-A", 1, []byte(payload)); err != nil {
		t.Fatal(err)
	}
	prov := Provenance{Generator: "store/test", SourceFingerprint: "sha256:test"}
	dump, fingerprint, err := CanonicalDump(db, res.EIR, prov)
	if err != nil {
		t.Fatal(err)
	}

	// A wrong expected fingerprint fails before anything is applied.
	fresh := filepath.Join(t.TempDir(), "fresh.db")
	if _, err := Import(fresh, res.EIR, dump, "sha256:0000"); err == nil || !strings.Contains(err.Error(), "does not match the expected") {
		t.Fatalf("wrong fingerprint accepted: %v", err)
	}
	if empty, err := isEmptyDatabase(fresh); err != nil || !empty {
		t.Fatalf("target was modified: empty=%v err=%v", empty, err)
	}

	// A tampered content block fails against the original fingerprint.
	tampered := strings.Replace(string(dump), "INV-A", "INV-Z", 1)
	if _, err := Import(filepath.Join(t.TempDir(), "tampered.db"), res.EIR, []byte(tampered), fingerprint); err == nil || !strings.Contains(err.Error(), "does not match the expected") {
		t.Fatalf("tampered dump accepted: %v", err)
	}

	// A non-empty target database is refused.
	if _, err := Import(db, res.EIR, dump, fingerprint); err == nil || !strings.Contains(err.Error(), "not empty") {
		t.Fatalf("non-empty target accepted: %v", err)
	}

	// Content after the canonical COMMIT block is refused.
	if _, err := Import(filepath.Join(t.TempDir(), "trailing.db"), res.EIR, append(dump, []byte("INSERT INTO \"invariant\" VALUES ('x');\n")...), fingerprint); err == nil || !strings.Contains(err.Error(), "after the canonical COMMIT") {
		t.Fatalf("trailing content accepted: %v", err)
	}
}

func TestImportRejectsBadHeader(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	_, res := ingestStore(t)
	bad := []byte("-- GENERATED canonical store dump by test - do not edit\n-- schema-version: 99\n-- dump-format: 2\nBEGIN;\nCOMMIT;\n")
	if _, err := Import(filepath.Join(t.TempDir(), "x.db"), res.EIR, bad, ""); err == nil || !strings.Contains(err.Error(), "schema version") {
		t.Fatalf("bad schema version accepted: %v", err)
	}
}
