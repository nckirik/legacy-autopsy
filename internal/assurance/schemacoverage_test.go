package assurance

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/store"
)

func TestSchemaCoverageDiagnostics(t *testing.T) {
	doc, db := testDocDB(t)
	table := seedInventory(t, db, doc, "TABLE", "schema.sql :: table :: malzemeler", "sha256:aaa")
	view := seedInventory(t, db, doc, "VIEW", "schema.sql :: view :: vw_ozet", "sha256:bbb")
	dr := seedInventory(t, db, doc, "DR", "schema.sql :: dr :: trg_touch", "sha256:ccc")
	action := seedInventory(t, db, doc, "ACTION", "pages/units.py :: action :: download", "sha256:ddd")
	_ = view
	_ = dr
	_ = action

	diagnostics, err := SchemaCoverageDiagnostics(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(diagnostics, "\n")
	for _, want := range []string{
		"declared table schema.sql :: table :: malzemeler has no committed claim",
		"declared view schema.sql :: view :: vw_ozet has no committed claim",
		"declared database routine schema.sql :: dr :: trg_touch has no committed claim",
		"discovered output action pages/units.py :: action :: download has no committed claim",
	} {
		if !strings.Contains(joined, want) {
			t.Fatalf("diagnostics missing %q:\n%s", want, joined)
		}
	}

	// Claims exist but no entity/routine records bind them yet.
	seedClaim(t, db, doc, "CLM-TABLE0000001", table)
	seedClaim(t, db, doc, "CLM-VIEW00000001", view)
	seedClaim(t, db, doc, "CLM-DR0000000001", dr)
	seedClaim(t, db, doc, "CLM-ACTION000001", action)
	diagnostics, err = SchemaCoverageDiagnostics(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	joined = strings.Join(diagnostics, "\n")
	if !strings.Contains(joined, "is not represented by an entity record") || !strings.Contains(joined, "is not represented by a database-routine record") {
		t.Fatalf("representation diagnostics missing:\n%s", joined)
	}
	if strings.Contains(joined, "action") {
		t.Fatalf("claimed action still diagnosed:\n%s", joined)
	}

	// Entity and routine records that bind the claims clear the diagnostics.
	entity := []byte("- **source-bindings:** CLM-TABLE0000001, CLM-VIEW00000001\n- **block-id:** ER-AAAAAAAAAAAA\n- **record-title:** Malzemeler\n- **domain-path:** demo/inventory/malzemeler\n- **domain-context-purpose:** Stores materials.\n- **storage:** table\n- **physical-src-ids:** None\n- **primary-key-secondary-keys-indexes:** id\n- **columns-fields:** id; ad; kategori_id\n- **relationships:** None\n- **enums-lookup-domains:** None\n- **lifecycle-state-fields:** None\n- **derived-computed-fields:** None\n- **invariants:** None\n- **audit-soft-delete:** None\n- **retention-lifecycle:** None\n- **permissions:** None\n- **capability-state-impact:** Unknown\n- **schema-drift-matrix:** None\n")
	if err := store.Ingest(db, doc, "ENTITY-RECORD", "ER-AAAAAAAAAAAA", 1, entity); err != nil {
		t.Fatal(err)
	}
	routine := []byte("- **source-bindings:** CLM-DR0000000001\n- **block-id:** DR-AAAAAAAAAAAA\n- **record-title:** Touch trigger\n- **domain-path:** demo/inventory/trigger\n- **kind:** trigger\n- **owner-table-entity:** malzemeler\n- **firing-event-timing-granularity:** before insert\n- **atomic-body-intent:** Sets timestamps.\n- **inputs-outputs:** None\n- **entities-read-written:** malzemeler\n- **invariants-rules-enforced:** None\n- **callers-firers:** None\n- **failure-transaction-semantics:** rollback\n- **canonical-db-source-drift:** schema.sql\n- **target-decision:** Undecided\n- **target-decision-dec:** None\n")
	if err := store.Ingest(db, doc, "DATABASE-ROUTINE", "DR-AAAAAAAAAAAA", 1, routine); err != nil {
		t.Fatal(err)
	}
	diagnostics, err = SchemaCoverageDiagnostics(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics remain:\n%s", strings.Join(diagnostics, "\n"))
	}
}
