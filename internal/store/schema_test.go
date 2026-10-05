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
		ID:      "!!!",
		Storage: "record",
		Fields:  []cdl.EIRFieldSpec{{Name: "ok", Type: "string"}},
	}}
	if _, err := Schema(doc, Provenance{Generator: "store/test"}); err == nil {
		t.Fatal("field without a usable table name accepted")
	}
	doc.Declarations.Fields = []cdl.EIRField{{
		ID:      "VALID",
		Storage: "record",
		Fields:  []cdl.EIRFieldSpec{{Name: "...", Type: "string"}},
	}}
	if _, err := Schema(doc, Provenance{Generator: "store/test"}); err == nil {
		t.Fatal("column without a usable name accepted")
	}
	doc.Declarations.Fields = []cdl.EIRField{{
		ID:     "UNCLASSIFIED",
		Fields: []cdl.EIRFieldSpec{{Name: "ok", Type: "string"}},
	}}
	if _, err := Schema(doc, Provenance{Generator: "store/test"}); err == nil {
		t.Fatal("field without a STORAGE class accepted")
	}
}

func TestSchemaKeysAndSkippedClasses(t *testing.T) {
	doc := &cdl.EIRDoc{}
	doc.Declarations.Fields = []cdl.EIRField{
		{ID: "THING", Storage: "record", Key: []string{"thing-id"}, Fields: []cdl.EIRFieldSpec{{Name: "thing-id", Type: "string", Required: true}}},
		{ID: "THING-COLUMN", Storage: "child", Parent: "THING", Key: []string{"name"}, Fields: []cdl.EIRFieldSpec{{Name: "name", Type: "string", Required: true}}},
		{ID: "THING-REPORT", Storage: "artifact", Fields: []cdl.EIRFieldSpec{{Name: "body", Type: "string", Required: true}}},
		{ID: "THING-FRAGMENT", Storage: "value", Fields: []cdl.EIRFieldSpec{{Name: "note", Type: "string"}}},
	}
	ddl, err := Schema(doc, Provenance{Generator: "store/test"})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`CREATE TABLE "thing"`, `UNIQUE ("thing_id")`} {
		if !strings.Contains(ddl, want) {
			t.Errorf("schema missing %q", want)
		}
	}
	for _, absent := range []string{`CREATE TABLE "thing_column"`, `CREATE TABLE "thing_report"`, `CREATE TABLE "thing_fragment"`} {
		if strings.Contains(ddl, absent) {
			t.Errorf("non-table class generated table %q", absent)
		}
	}
	if strings.Count(ddl, "CREATE TABLE ") != 1 {
		t.Fatalf("expected exactly one table, got %d", strings.Count(ddl, "CREATE TABLE "))
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
		`"constraint" TEXT NOT NULL`,
	} {
		if !strings.Contains(ddl, want) {
			t.Errorf("schema missing %q", want)
		}
	}
	if !strings.Contains(ddl, `CHECK ("discovery_method" IN ('filesystem'`) {
		t.Error("enum CHECK does not list enum values")
	}
	tables := strings.Count(ddl, "CREATE TABLE ")
	if tables != 48 {
		t.Fatalf("generated %d tables, expected 48 table-class FIELD blocks", tables)
	}
	for _, absent := range []string{`CREATE TABLE "decision_content_count"`, `CREATE TABLE "entity_column"`} {
		if strings.Contains(ddl, absent) {
			t.Errorf("non-table class generated table %q", absent)
		}
	}
	if !strings.Contains(ddl, "-- schema-version: 2") {
		t.Error("schema version provenance missing")
	}
}
