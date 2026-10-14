package equivalence

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/spec"
	"github.com/nckirik/legacy-autopsy/internal/store"
)

func TestBuildEmptyStateAndParseMappings(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	res, err := spec.Compile("../..")
	if err != nil {
		t.Fatal(err)
	}
	db := t.TempDir() + "/store.db"
	ddl, err := store.Schema(res.EIR, store.Provenance{Generator: "equivalence/test", SourceFingerprint: "sha256:test"})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Apply(db, ddl); err != nil {
		t.Fatal(err)
	}
	suite, err := Build(db, res.EIR)
	if err != nil {
		t.Fatal(err)
	}
	if len(suite.Mappings) != 0 {
		t.Fatalf("mappings = %d, want 0 on empty state", len(suite.Mappings))
	}
	synthetic := []byte("| Record Type | Typed Record ID | Skeleton |\n| --- | --- | --- |\n| UC | UC-ABCDEF012345 | #uc-abcdef012345 |\n\n### uc-abcdef012345\n\n- **Use case:** UC-ABCDEF012345 - Deploy trigger\n")
	mappings := ParseMappings(synthetic)
	if len(mappings) != 1 || mappings[0].RecordID != "UC-ABCDEF012345" || mappings[0].RecordType != "UC" {
		t.Fatalf("parsed mappings = %+v", mappings)
	}
	if !strings.Contains(mappings[0].Body, "Deploy trigger") {
		t.Fatalf("mapping body not captured: %+v", mappings[0])
	}
}
