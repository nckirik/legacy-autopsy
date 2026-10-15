package store

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/spec"
)

func ingestStore(t *testing.T) (string, *cdl.Result) {
	t.Helper()
	res, err := spec.Compile(repoRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	db := filepath.Join(t.TempDir(), "store.db")
	ddl, err := Schema(res.EIR, Provenance{Generator: "store/test", SourceFingerprint: "sha256:test"})
	if err != nil {
		t.Fatal(err)
	}
	if err := Apply(db, ddl); err != nil {
		t.Fatal(err)
	}
	return db, res
}

func TestIngestRecord(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	db, res := ingestStore(t)
	payload := "- **Constraint:** all totals match\n- **Verification Vectors:** zeta, alpha\n- **Discovered In:** SRC-0001\n"
	if err := Ingest(db, res.EIR, "INVARIANT", "INV-0001", 1, []byte(payload)); err != nil {
		t.Fatal(err)
	}
	set, err := sqlite(t, db, `select verification_vectors from invariant where la_record_id = 'INV-0001';`)
	if err != nil {
		t.Fatal(err, set)
	}
	if set != "alpha, zeta" {
		t.Fatalf("set column = %q, want canonical %q", set, "alpha, zeta")
	}
	fingerprint, err := sqlite(t, db, `select la_semantic_fingerprint from invariant where la_record_id = 'INV-0001';`)
	if err != nil {
		t.Fatal(err, fingerprint)
	}
	canonical, err := capabilities.CanonicalMarkdown([]byte(payload))
	if err != nil {
		t.Fatal(err)
	}
	if fingerprint != capabilities.Hash(canonical) {
		t.Fatal("stored fingerprint does not match the canonical payload")
	}
	stored, err := sqlite(t, db, `select la_canonical_payload from invariant where la_record_id = 'INV-0001';`)
	if err != nil {
		t.Fatal(err, stored)
	}
	if stored != strings.TrimSuffix(string(canonical), "\n") {
		t.Fatalf("payload not stored canonically: %q", stored)
	}
}

func TestIngestCanonicalizesTransportBytes(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	db, res := ingestStore(t)
	messy := "- **Constraint:** match  \r\n- **Verification Vectors:** v1\r\n- **Discovered In:** SRC-0002\r\n"
	if err := Ingest(db, res.EIR, "INVARIANT", "INV-0002", 1, []byte(messy)); err != nil {
		t.Fatal(err)
	}
	first, err := sqlite(t, db, `select la_semantic_fingerprint from invariant where la_record_id = 'INV-0002';`)
	if err != nil {
		t.Fatal(err, first)
	}
	clean := "- **Constraint:** match\n- **Verification Vectors:** v1\n- **Discovered In:** SRC-0002\n"
	if err := Ingest(db, res.EIR, "INVARIANT", "INV-0003", 1, []byte(clean)); err != nil {
		t.Fatal(err)
	}
	second, err := sqlite(t, db, `select la_semantic_fingerprint from invariant where la_record_id = 'INV-0003';`)
	if err != nil {
		t.Fatal(err, second)
	}
	if first != second {
		t.Fatal("transport variance changed the payload fingerprint")
	}
}

func TestIngestFailsClosed(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	db, res := ingestStore(t)
	valid := "- **Constraint:** c\n- **Verification Vectors:** v1\n- **Discovered In:** SRC-0001\n"
	cases := []struct {
		name    string
		field   string
		record  string
		payload string
		want    string
	}{
		{"missing required", "INVARIANT", "INV-1", "- **Verification Vectors:** v1\n- **Discovered In:** SRC-1\n", "missing required"},
		{"unknown field", "INVARIANT", "INV-2", valid + "- **Extra:** x\n", "unknown field"},
		{"duplicate set", "INVARIANT", "INV-3", "- **Constraint:** c\n- **Verification Vectors:** a, a\n- **Discovered In:** SRC-1\n", "duplicate"},
		{"non-table class", "BLOCK-EVIDENCE-SUMMARY", "BES-1", "- **block-confidence-rank:** Low\n", "not a table class"},
		{"child without parent", "ENTITY-COLUMN", "EC-1", "- **name:** n\n", "parent record ID is required"},
		{"unknown field block", "NOPE", "X-1", valid, "unknown FIELD"},
	}
	for _, tc := range cases {
		err := Ingest(db, res.EIR, tc.field, tc.record, 1, []byte(tc.payload))
		if err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("%s: err = %v, want %q", tc.name, err, tc.want)
		}
	}
	if err := Ingest(db, res.EIR, "INVARIANT", "INV-5", 1, []byte(valid)); err != nil {
		t.Fatal(err)
	}
	if err := Ingest(db, res.EIR, "INVARIANT", "INV-5", 2, []byte(valid)); err == nil {
		t.Fatal("duplicate record ID accepted")
	}
}

func TestIngestRejectsUnknownEnum(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	db, res := ingestStore(t)
	payload := "- **Capability State:** Nonsense\n- **Runtime State Matrix:** m\n- **Historical Intended Personas:** h\n- **Dependencies:** d\n- **Last Known State:** l\n"
	err := Ingest(db, res.EIR, "CAPABILITY-STATE-BLOCK", "CAP-1", 1, []byte(payload))
	if err == nil || !strings.Contains(err.Error(), "not in enum") {
		t.Fatalf("invalid registry value accepted: %v", err)
	}
}
