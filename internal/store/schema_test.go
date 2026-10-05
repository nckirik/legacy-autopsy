package store

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/spec"
)

func schema(t *testing.T) string {
	t.Helper()
	res, err := spec.Compile(repoRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	ddl, err := Schema(res.EIR, Provenance{Generator: "store/test", SourceFingerprint: "sha256:test"})
	if err != nil {
		t.Fatal(err)
	}
	return ddl
}

func repoRoot(t *testing.T) string {
	t.Helper()
	return "../.."
}

func TestSchemaDeterministic(t *testing.T) {
	first := schema(t)
	second := schema(t)
	if first != second {
		t.Fatal("schema generation is not deterministic")
	}
}

func TestSchemaRejectsUnusableIdentifiers(t *testing.T) {
	doc := &cdl.EIRDoc{}
	doc.Declarations.Fields = []cdl.EIRField{{
		ID:     "!!!",
		Fields: []cdl.EIRFieldSpec{{Name: "ok", Type: "string"}},
	}}
	if _, err := Schema(doc, Provenance{Generator: "store/test"}); err == nil {
		t.Fatal("field without a usable table name accepted")
	}
	doc.Declarations.Fields = []cdl.EIRField{{
		ID:     "VALID",
		Fields: []cdl.EIRFieldSpec{{Name: "...", Type: "string"}},
	}}
	if _, err := Schema(doc, Provenance{Generator: "store/test"}); err == nil {
		t.Fatal("column without a usable name accepted")
	}
}

func TestSchemaTablesAndChecks(t *testing.T) {
	ddl := schema(t)
	for _, want := range []string{
		`CREATE TABLE "source_inventory_record"`,
		`"la_record_id" TEXT PRIMARY KEY`,
		`"concrete_coordinate" TEXT NOT NULL`,
		`CREATE TABLE "invocation_log"`,
		`BEFORE UPDATE ON "invocation_log"`,
		`BEFORE DELETE ON "invocation_log"`,
		`RAISE(ABORT, '0G is append-only')`,
		`"count" INTEGER`,
	} {
		if !strings.Contains(ddl, want) {
			t.Errorf("schema missing %q", want)
		}
	}
	if !strings.Contains(ddl, `CHECK ("discovery_method" IN ('filesystem'`) {
		t.Error("enum CHECK does not list enum values")
	}
	tables := strings.Count(ddl, "CREATE TABLE ")
	if tables < 80 {
		t.Fatalf("generated %d tables, expected one per FIELD block", tables)
	}
	if !strings.Contains(ddl, "-- schema-version: 1") {
		t.Error("schema version provenance missing")
	}
}
