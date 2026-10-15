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
		{ID: "THING-COLUMN", Storage: "child", Parent: "THING", Key: []string{"name"}, Fields: []cdl.EIRFieldSpec{{Name: "name", Type: "string", Required: true}, {Name: "thing-id", Type: "string", Required: true}}},
		{ID: "THING-REPORT", Storage: "artifact", Fields: []cdl.EIRFieldSpec{{Name: "body", Type: "string", Required: true}}},
		{ID: "THING-FRAGMENT", Storage: "value", Fields: []cdl.EIRFieldSpec{{Name: "note", Type: "string"}}},
	}
	doc.Declarations.Fields[1].Refs = []cdl.EIRFieldRef{{Column: "thing-id", Target: "THING.thing-id"}}
	ddl, err := Schema(doc, Provenance{Generator: "store/test"})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`CREATE TABLE "thing"`,
		`UNIQUE ("thing_id")`,
		`CREATE TABLE "thing_column"`,
		`"la_parent_id" TEXT NOT NULL`,
		`UNIQUE ("la_parent_id", "name")`,
		`FOREIGN KEY ("la_parent_id") REFERENCES "thing" ("la_record_id")`,
		`FOREIGN KEY ("thing_id") REFERENCES "thing" ("thing_id")`,
	} {
		if !strings.Contains(ddl, want) {
			t.Errorf("schema missing %q", want)
		}
	}
	for _, absent := range []string{`CREATE TABLE "thing_report"`, `CREATE TABLE "thing_fragment"`} {
		if strings.Contains(ddl, absent) {
			t.Errorf("non-table class generated table %q", absent)
		}
	}
	if strings.Count(ddl, "CREATE TABLE ") != 2 {
		t.Fatalf("expected exactly two tables, got %d", strings.Count(ddl, "CREATE TABLE "))
	}
}

func TestSchemaChildAndRefFailClosed(t *testing.T) {
	base := func() *cdl.EIRDoc {
		doc := &cdl.EIRDoc{}
		doc.Declarations.Fields = []cdl.EIRField{
			{ID: "THING", Storage: "record", Fields: []cdl.EIRFieldSpec{{Name: "thing-id", Type: "string", Required: true}}},
		}
		return doc
	}
	child := cdl.EIRField{ID: "THING-COLUMN", Storage: "child", Key: []string{"name"}, Fields: []cdl.EIRFieldSpec{{Name: "name", Type: "string", Required: true}}}
	withParent := child
	withParent.Parent = "THING"
	missingParent := child
	missingParent.Parent = "NOPE"
	noKey := withParent
	noKey.Key = nil
	badRef := withParent
	badRef.Refs = []cdl.EIRFieldRef{{Column: "name", Target: "THING.missing-column"}}
	valueRef := withParent
	valueRef.Refs = []cdl.EIRFieldRef{{Column: "name", Target: "THING-FRAGMENT.note"}}
	cases := []struct {
		name  string
		field cdl.EIRField
		want  string
	}{
		{"missing parent", missingParent, "PARENT"},
		{"no key", noKey, "KEY"},
		{"unknown ref column", badRef, "not a declared field"},
	}
	for _, tc := range cases {
		doc := base()
		doc.Declarations.Fields = append(doc.Declarations.Fields, tc.field)
		if _, err := Schema(doc, Provenance{Generator: "store/test"}); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("%s: err = %v, want %q", tc.name, err, tc.want)
		}
	}
	doc := base()
	doc.Declarations.Fields = append(doc.Declarations.Fields,
		cdl.EIRField{ID: "THING-FRAGMENT", Storage: "value", Fields: []cdl.EIRFieldSpec{{Name: "note", Type: "string"}}},
		valueRef,
	)
	if _, err := Schema(doc, Provenance{Generator: "store/test"}); err == nil || !strings.Contains(err.Error(), "not a referenceable class") {
		t.Fatalf("value-class ref target accepted: %v", err)
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
	res2, err := spec.Compile(repoRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	expected := 0
	for _, field := range res2.EIR.Declarations.Fields {
		switch field.Storage {
		case "record", "registry", "singleton", "log", "child":
			expected++
		}
	}
	tables := strings.Count(ddl, "CREATE TABLE ")
	if tables != expected {
		t.Fatalf("generated %d tables, expected %d table-class FIELD blocks", tables, expected)
	}
	if !strings.Contains(ddl, `CREATE TABLE "entity_column"`) || !strings.Contains(ddl, `FOREIGN KEY ("la_parent_id") REFERENCES "entity_record" ("la_record_id")`) {
		t.Error("child projection missing for ENTITY-COLUMN")
	}
	if !strings.Contains(ddl, `FOREIGN KEY ("cluster_id") REFERENCES "entry_cluster" ("cluster_id")`) {
		t.Error("declared REF foreign key missing")
	}
	for _, absent := range []string{`CREATE TABLE "decision_content_count"`, `CREATE TABLE "block_evidence_summary"`} {
		if strings.Contains(ddl, absent) {
			t.Errorf("non-table class generated table %q", absent)
		}
	}
	if !strings.Contains(ddl, "-- schema-version: 3") {
		t.Error("schema version provenance missing")
	}
}
