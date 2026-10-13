package handbook

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/store"
)

func TestEvidenceRegisterExplainsClaims(t *testing.T) {
	doc, db := testDocDB(t)
	payload := "- **subject-id-record-version:** SRC-FILE-AAAAAAAAAAAA\n- **field-step:** file\n- **claim:** app/routes.js registers the public and API routes\n- **materiality:** Material\n- **evidence-level:** [E-DIRECT]\n- **anchor:** app/routes.js:1\n- **source-fingerprint-snapshot:** sha256:aaa\n- **evidence-route-lineage-id:** static-filesystem\n- **projection-classification:** raw-visible\n- **independent-supporting-claims:** None\n- **independence-basis:** None\n- **contradiction-ids:** None\n- **created-last-revalidated-by:** None\n"
	if err := store.Ingest(db, doc, "CLAIM", "CLM-000000000001", 1, []byte(payload)); err != nil {
		t.Fatal(err)
	}
	out, err := EvidenceRegister(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	for _, marker := range []string{
		"# Evidence Register",
		`<a id="clm-000000000001"></a>`,
		"## `CLM-000000000001`",
		"- **Claim:** app/routes.js registers the public and API routes",
		"- **Anchor:** `app/routes.js:1`",
		"- **Evidence Level:** [E-DIRECT]",
	} {
		if !strings.Contains(string(out), marker) {
			t.Fatalf("register missing %q:\n%s", marker, out)
		}
	}
}
