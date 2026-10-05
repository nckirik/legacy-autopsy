package store

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func sqlite(t *testing.T, db string, sql string) (string, error) {
	t.Helper()
	cmd := exec.Command("sqlite3", db, sql)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func TestApplySchema(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	ddl := schema(t)
	db := filepath.Join(t.TempDir(), "store.db")
	if err := Apply(db, ddl); err != nil {
		t.Fatal(err)
	}
	out, err := sqlite(t, db, "select count(*) from sqlite_master where type='table';")
	if err != nil {
		t.Fatal(err, out)
	}
	if strings.Count(ddl, "CREATE TABLE ") == 0 || out == "" {
		t.Fatal("no tables created")
	}
	out, err = sqlite(t, db, "select count(*) from sqlite_master where type='trigger' and name like 'invocation_log%';")
	if err != nil {
		t.Fatal(err, out)
	}
	if out != "2" {
		t.Fatalf("expected two append-only triggers, got %q", out)
	}
}

func TestApplyEnforcesEnums(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	db := filepath.Join(t.TempDir(), "store.db")
	if err := Apply(db, schema(t)); err != nil {
		t.Fatal(err)
	}
	valid := `insert into capability_state_block (la_record_id, la_canonical_payload, capability_state, runtime_state_matrix, historical_intended_personas, dependencies, last_known_state) values ('r1','{}','Active','m','h','d','l');`
	if _, err := sqlite(t, db, valid); err != nil {
		t.Fatalf("valid row rejected: %v", err)
	}
	invalid := `insert into capability_state_block (la_record_id, la_canonical_payload, capability_state, runtime_state_matrix, historical_intended_personas, dependencies, last_known_state) values ('r2','{}','Nonsense','m','h','d','l');`
	if _, err := sqlite(t, db, invalid); err == nil {
		t.Fatal("invalid enum accepted")
	}
}
