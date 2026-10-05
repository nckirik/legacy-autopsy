package store

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
)

// payloadFor builds a minimal valid payload for a table-class FIELD by using
// each declared spec name as its own label (ColumnName is idempotent on the
// derived kebab names).
func payloadFor(t *testing.T, doc *cdl.EIRDoc, fieldID string) []byte {
	t.Helper()
	field, err := tableField(doc, fieldID)
	if err != nil {
		t.Fatal(err)
	}
	var b strings.Builder
	for _, spec := range field.Fields {
		switch spec.Type {
		case "int":
			b.WriteString("- **" + spec.Name + ":** 1\n")
		case "set<string>":
			b.WriteString("- **" + spec.Name + ":** entry\n")
		default:
			b.WriteString("- **" + spec.Name + ":** value\n")
		}
	}
	return []byte(b.String())
}

func TestVerifyMonotonicLogs(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	db, res := ingestStore(t)
	doc := res.EIR
	if err := Ingest(db, doc, "INVOCATION-LOG", "LOG-1", 1, payloadFor(t, doc, "INVOCATION-LOG")); err != nil {
		t.Fatal(err)
	}
	before, err := LogState(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	if len(before["INVOCATION-LOG"]) != 1 {
		t.Fatalf("state = %+v", before)
	}
	if err := Ingest(db, doc, "INVOCATION-LOG", "LOG-2", 1, payloadFor(t, doc, "INVOCATION-LOG")); err != nil {
		t.Fatal(err)
	}
	after, err := LogState(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	if err := VerifyMonotonicLogs(before, after); err != nil {
		t.Fatalf("append rejected: %v", err)
	}
	// The trigger blocks deletes; drop it to prove the runtime post-condition
	// catches a deletion even if the database-level guard is bypassed.
	if _, err := sqlite(t, db, `drop trigger "invocation_log_no_delete"; delete from "invocation_log" where "la_record_id" = 'LOG-1';`); err != nil {
		t.Fatal(err)
	}
	deleted, err := LogState(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	if err := VerifyMonotonicLogs(before, deleted); err == nil || !strings.Contains(err.Error(), "append-only") {
		t.Fatalf("deletion accepted: %v", err)
	}
}
