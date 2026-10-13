package handbook

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/spec"
	"github.com/nckirik/legacy-autopsy/internal/store"
)

func testDocDB(t *testing.T) (*cdl.EIRDoc, string) {
	t.Helper()
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	res, err := spec.Compile("../..")
	if err != nil {
		t.Fatal(err)
	}
	db := filepath.Join(t.TempDir(), "store.db")
	ddl, err := store.Schema(res.EIR, store.Provenance{Generator: "handbook/test", SourceFingerprint: "sha256:test"})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Apply(db, ddl); err != nil {
		t.Fatal(err)
	}
	return res.EIR, db
}

func TestDraftRendersChaptersAndLint(t *testing.T) {
	doc, db := testDocDB(t)
	for i := 1; i <= 8; i++ {
		id := "MOD-00000000000" + string(rune('0'+i))
		payload := "- **source-bindings:** CLM-00000000000" + string(rune('0'+i)) + "\n" +
			"- **block-id:** " + id + "\n" +
			"- **record-title:** Module " + string(rune('A'+i-1)) + "\n" +
			"- **domain-path:** app/m" + string(rune('0'+i)) + "\n" +
			"- **responsibility-domain-boundary:** The module owns one boundary.\n" +
			"- **public-interfaces:** None\n" +
			"- **encapsulated-state-entities:** None\n" +
			"- **dependencies-dependents:** None\n" +
			"- **data-and-events-passed:** None\n" +
			"- **failure-propagation:** None\n" +
			"- **invariants-rules-owned:** None\n" +
			"- **observed-legacy-placement:** app\n" +
			"- **target-decision:** Pending\n" +
			"- **personas-use-cases-capabilities:** None\n"
		if err := store.Ingest(db, doc, "ARCHITECTURE-MODULE", id, 1, []byte(payload)); err != nil {
			t.Fatal(err)
		}
	}
	out, findings, err := Draft(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	text := string(out)
	for _, marker := range []string{"# Legacy System Reconstruction Handbook", "## 01. System Overview", "## 16. Glossary and Reference", "### Audience Purpose", "### Observed Legacy Behavior", "### Navigation"} {
		if !strings.Contains(text, marker) {
			t.Fatalf("draft missing %q", marker)
		}
	}
	if len(findings) != 0 {
		t.Fatalf("draft lint findings: %v", findings)
	}
}
