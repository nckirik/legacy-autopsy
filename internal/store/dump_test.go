package store

import (
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/spec"
)

func buildStore(t *testing.T, insertOrder []string) (string, string, string) {
	t.Helper()
	res, err := spec.Compile(repoRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	db := filepath.Join(t.TempDir(), "store.db")
	prov := Provenance{Generator: "store/test", SourceFingerprint: "sha256:test"}
	ddl, err := Schema(res.EIR, prov)
	if err != nil {
		t.Fatal(err)
	}
	if err := Apply(db, ddl); err != nil {
		t.Fatal(err)
	}
	for _, id := range insertOrder {
		sql := `insert into decision_content_count (la_record_id, la_canonical_payload, decision_type, count) values ('` + id + `','{}','scope',1);`
		if _, err := sqlite(t, db, sql); err != nil {
			t.Fatal(err)
		}
	}
	dump, fp, err := CanonicalDump(db, res.EIR, prov)
	if err != nil {
		t.Fatal(err)
	}
	return db, string(dump), fp
}

func dumpFixture(t *testing.T, insertOrder []string) (string, string) {
	t.Helper()
	_, dump, fp := buildStore(t, insertOrder)
	return dump, fp
}

func TestCanonicalDumpDeterministic(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	first, fpFirst := dumpFixture(t, []string{"rec-a", "rec-b"})
	second, fpSecond := dumpFixture(t, []string{"rec-b", "rec-a"})
	if first != second || fpFirst != fpSecond {
		t.Fatal("canonical dump depends on insertion order")
	}
	if strings.Index(first, "rec-a") > strings.Index(first, "rec-b") {
		t.Fatal("rows are not sorted by la_record_id")
	}
	if !strings.Contains(first, "BEGIN;") || !strings.Contains(first, "COMMIT;") {
		t.Fatal("dump is not a transaction")
	}
	if !strings.Contains(first, `CREATE TABLE "decision_content_count"`) {
		t.Fatal("dump does not carry the relational schema")
	}
	third, fpThird := dumpFixture(t, []string{"rec-a", "rec-b", "rec-c"})
	if third == first || fpThird == fpFirst {
		t.Fatal("new row did not change the dump and fingerprint")
	}
}

func TestCanonicalDumpRestorable(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	_, dump, _ := buildStore(t, []string{"rec-a"})
	restored := filepath.Join(t.TempDir(), "restored.db")
	cmd := exec.Command("sqlite3", "--bail", restored)
	cmd.Stdin = strings.NewReader(dump)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("dump is not restorable: %v: %s", err, out)
	}
	wantTables := strings.Count(dump, "CREATE TABLE ")
	out, err := sqlite(t, restored, "select count(*) from sqlite_master where type='table';")
	if err != nil {
		t.Fatal(err, out)
	}
	if out != strconv.Itoa(wantTables) {
		t.Fatalf("restored %s tables, dump declares %d", out, wantTables)
	}
	rows, err := sqlite(t, restored, "select count(*) from decision_content_count;")
	if err != nil {
		t.Fatal(err, rows)
	}
	if rows != "1" {
		t.Fatalf("restored row count = %s, want 1", rows)
	}
}

func TestCanonicalDumpRejectsSchemaDrift(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	db, _, _ := buildStore(t, []string{"rec-a"})
	if _, err := sqlite(t, db, `create table "zz_extra" ("x" text);`); err != nil {
		t.Fatal(err)
	}
	res, err := spec.Compile(repoRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = CanonicalDump(db, res.EIR, Provenance{Generator: "store/test", SourceFingerprint: "sha256:test"})
	if err == nil || !strings.Contains(err.Error(), "not present in the generated EIR schema") {
		t.Fatalf("extra schema object not rejected: %v", err)
	}
}

func TestCanonicalDumpRejectsIndexDrift(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	db, _, _ := buildStore(t, []string{"rec-a"})
	if _, err := sqlite(t, db, `create index "zz_idx" on "decision_content_count" ("count");`); err != nil {
		t.Fatal(err)
	}
	res, err := spec.Compile(repoRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = CanonicalDump(db, res.EIR, Provenance{Generator: "store/test", SourceFingerprint: "sha256:test"})
	if err == nil || !strings.Contains(err.Error(), "not present in the generated EIR schema") {
		t.Fatalf("index drift not rejected: %v", err)
	}
}

func TestCanonicalDumpRejectsMissingSchema(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	db, _, _ := buildStore(t, []string{"rec-a"})
	if _, err := sqlite(t, db, `drop table "invocation_log";`); err != nil {
		t.Fatal(err)
	}
	res, err := spec.Compile(repoRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = CanonicalDump(db, res.EIR, Provenance{Generator: "store/test", SourceFingerprint: "sha256:test"})
	if err == nil || !strings.Contains(err.Error(), "missing") {
		t.Fatalf("missing schema object not rejected: %v", err)
	}
}

func TestCanonicalDumpFingerprintExcludesProvenance(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
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
	provA := Provenance{Generator: "store/1", SourceFingerprint: "sha256:a"}
	provB := Provenance{Generator: "store/2", SourceFingerprint: "sha256:b"}
	dumpA, fpA, err := CanonicalDump(db, res.EIR, provA)
	if err != nil {
		t.Fatal(err)
	}
	dumpB, fpB, err := CanonicalDump(db, res.EIR, provB)
	if err != nil {
		t.Fatal(err)
	}
	if fpA != fpB {
		t.Fatal("content fingerprint changed with provenance")
	}
	if string(dumpA) == string(dumpB) {
		t.Fatal("provenance header missing from dump bytes")
	}
}
