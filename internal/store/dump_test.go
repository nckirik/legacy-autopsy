package store

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/spec"
)

func dumpFixture(t *testing.T, insertOrder []string) (string, string) {
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
	return string(dump), fp
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
	third, fpThird := dumpFixture(t, []string{"rec-a", "rec-b", "rec-c"})
	if third == first || fpThird == fpFirst {
		t.Fatal("new row did not change the dump and fingerprint")
	}
}
