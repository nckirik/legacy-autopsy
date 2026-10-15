package store

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func entityPayload() []byte {
	return []byte("- **source-bindings:** CLM-AAAAAAAAAAAA\n" +
		"- **block-id:** ER-AAAAAAAAAAAA\n" +
		"- **record-title:** Malzemeler\n" +
		"- **domain-path:** demo/inventory/malzemeler\n" +
		"- **domain-context-purpose:** Stores materials.\n" +
		"- **storage:** table\n" +
		"- **physical-src-ids:** None\n" +
		"- **primary-key-secondary-keys-indexes:** id\n" +
		"- **columns-fields:** id; ad; kategori_id\n" +
		"- **relationships:** None\n" +
		"- **enums-lookup-domains:** None\n" +
		"- **lifecycle-state-fields:** None\n" +
		"- **derived-computed-fields:** None\n" +
		"- **invariants:** None\n" +
		"- **audit-soft-delete:** None\n" +
		"- **retention-lifecycle:** None\n" +
		"- **permissions:** None\n" +
		"- **capability-state-impact:** Unknown\n" +
		"- **schema-drift-matrix:** None\n")
}

func entityColumnPayload(name string) []byte {
	return []byte("- **name:** " + name + "\n" +
		"- **type:** INT\n" +
		"- **null:** no\n" +
		"- **default:** None\n" +
		"- **constraints:** primary key\n" +
		"- **unique:** yes\n" +
		"- **claims:** CLM-AAAAAAAAAAAA\n")
}

func TestIngestChildAndForeignKeys(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	db, res := ingestStore(t)
	if err := Ingest(db, res.EIR, "ENTITY-RECORD", "ER-AAAAAAAAAAAA", 1, entityPayload()); err != nil {
		t.Fatal(err)
	}
	childID, err := IngestChild(db, res.EIR, "ENTITY-COLUMN", "ER-AAAAAAAAAAAA", 1, entityColumnPayload("id"))
	if err != nil {
		t.Fatal(err)
	}
	if childID != "ER-AAAAAAAAAAAA|id" {
		t.Fatalf("child record ID = %q, want ER-AAAAAAAAAAAA|id", childID)
	}
	rows, err := queryRows(db, `select la_parent_id from "entity_column" where la_record_id = 'ER-AAAAAAAAAAAA|id';`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0]["la_parent_id"] != "ER-AAAAAAAAAAAA" {
		t.Fatalf("parent column = %+v", rows)
	}

	// The parent must exist: the declared foreign key is enforced.
	if _, err := IngestChild(db, res.EIR, "ENTITY-COLUMN", "ER-MISSING00000", 1, entityColumnPayload("id")); err == nil || !strings.Contains(err.Error(), "FOREIGN KEY") {
		t.Fatalf("unknown parent accepted: %v", err)
	}
	// The natural key must be unique per parent.
	if _, err := IngestChild(db, res.EIR, "ENTITY-COLUMN", "ER-AAAAAAAAAAAA", 1, entityColumnPayload("id")); err == nil || !strings.Contains(strings.ToLower(err.Error()), "unique") {
		t.Fatalf("duplicate child key accepted: %v", err)
	}
	// A key field must be present in the payload.
	if _, err := IngestChild(db, res.EIR, "ENTITY-COLUMN", "ER-AAAAAAAAAAAA", 1, []byte("- **type:** INT\n- **null:** no\n- **default:** None\n- **constraints:** none\n- **unique:** no\n- **claims:** None\n")); err == nil || !strings.Contains(err.Error(), "key field") {
		t.Fatalf("missing child key accepted: %v", err)
	}
	// ReplaceRecord preserves the parent link on version updates.
	if err := ReplaceRecord(db, res.EIR, "ENTITY-COLUMN", "ER-AAAAAAAAAAAA|id", 2, entityColumnPayload("id")); err != nil {
		t.Fatal(err)
	}
	rows, err = queryRows(db, `select la_parent_id, la_record_version from "entity_column" where la_record_id = 'ER-AAAAAAAAAAAA|id';`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0]["la_parent_id"] != "ER-AAAAAAAAAAAA" {
		t.Fatalf("parent column after replace = %+v", rows)
	}
}

func TestChildDumpRoundTrip(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	db, res := ingestStore(t)
	if err := Ingest(db, res.EIR, "ENTITY-RECORD", "ER-AAAAAAAAAAAA", 1, entityPayload()); err != nil {
		t.Fatal(err)
	}
	if _, err := IngestChild(db, res.EIR, "ENTITY-COLUMN", "ER-AAAAAAAAAAAA", 1, entityColumnPayload("id")); err != nil {
		t.Fatal(err)
	}
	if _, err := IngestChild(db, res.EIR, "ENTITY-COLUMN", "ER-AAAAAAAAAAAA", 1, entityColumnPayload("ad")); err != nil {
		t.Fatal(err)
	}
	prov := Provenance{Generator: "store/test", SourceFingerprint: "sha256:test"}
	dump, fingerprint, err := CanonicalDump(db, res.EIR, prov)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(dump), `"la_parent_id"`) {
		t.Fatalf("dump does not carry the child parent column:\n%s", dump)
	}
	fresh := filepath.Join(t.TempDir(), "fresh.db")
	if _, err := Import(fresh, res.EIR, dump, fingerprint); err != nil {
		t.Fatal(err)
	}
	again, againFingerprint, err := CanonicalDump(fresh, res.EIR, prov)
	if err != nil {
		t.Fatal(err)
	}
	if againFingerprint != fingerprint || string(again) != string(dump) {
		t.Fatal("child round-trip dump is not byte-identical")
	}
}

func TestRefForeignKeysAreEnforced(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	db, res := ingestStore(t)
	sweep := "- **cluster-id:** CL-1\n- **persona-owner:** WEB\n- **pov:** POV-1\n- **track:** Normal\n- **environment:** dev\n- **snapshot:** snap-1\n- **snapshot-environment:** snap-1/dev\n- **persona-cluster-pov-track:** WEB|CL-1|POV-1|Normal\n- **expanded-concrete-entry-coordinates:** api.js\n- **source-denominator:** 1\n- **expected-mapped-excluded-counts-by-kind:** 1|0\n- **open-terminal-frontier-counts:** 0|0\n- **active-terminal-ticket-counts:** 0|0\n- **export-reconciliation:** None\n- **coverage-result:** [R-SWEPT]\n- **closing-invocation:** ALFA-NONE-NONE-001\n- **evidence-claims:** CLM-AAAAAAAAAAAA\n"
	if err := Ingest(db, res.EIR, "SWEEP-RECORD", "CL-1|WEB|POV-1|Normal|dev|snap-1", 1, []byte(sweep)); err == nil || !strings.Contains(err.Error(), "FOREIGN KEY") {
		t.Fatalf("sweep without its referenced cluster accepted: %v", err)
	}
	cluster := "- **cluster-id:** CL-1\n- **cluster-title:** Web entry\n- **concrete-entry-coordinates:** api.js\n- **canonical-owner-persona:** WEB\n- **primary-pov:** POV-1\n- **traversal-track:** Normal\n- **capability-state:** Active\n- **environment-snapshot:** dev\n- **source-ids:** None\n- **max-depth:** 3\n- **traversal-status:** [R-UNSWEPT]\n- **required-traversal-matrix:** None\n"
	if err := Ingest(db, res.EIR, "ENTRY-CLUSTER", "CL-1", 1, []byte(cluster)); err != nil {
		t.Fatal(err)
	}
	if err := Ingest(db, res.EIR, "SWEEP-RECORD", "CL-1|WEB|POV-1|Normal|dev|snap-1", 1, []byte(sweep)); err != nil {
		t.Fatalf("sweep with its referenced cluster rejected: %v", err)
	}
}
