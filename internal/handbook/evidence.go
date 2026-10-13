package handbook

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// EvidenceRegister renders the handbook evidence appendix: every committed claim
// with its text, subject source unit, anchor, materiality, and evidence level,
// so chapter footnotes resolve to a readable register.
func EvidenceRegister(db string, doc *cdl.EIRDoc) ([]byte, error) {
	names, _, _, err := work.FieldContract(doc, "CLAIM")
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, "CLAIM")
	if err != nil {
		return nil, err
	}
	var b strings.Builder
	b.WriteString("# Evidence Register\n\n")
	b.WriteString("This appendix explains the claim IDs that chapter footnotes cite.\n")
	b.WriteString("Each `CLM-*` record is one evidence-backed fact. The source fingerprint and\n")
	b.WriteString("anchor make the fact reproducible from the original evidence.\n\n")
	if len(ids) == 0 {
		b.WriteString("None\n")
		return []byte(b.String()), nil
	}
	sort.Strings(ids)
	for _, id := range ids {
		payload, err := store.Get(db, doc, "CLAIM", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(&b, "## `%s`\n\n", id)
		fmt.Fprintf(&b, "- **Claim:** %s\n", strings.Join(strings.Fields(values["claim"]), " "))
		fmt.Fprintf(&b, "- **Subject:** %s\n", strings.Fields(values["subject-id-record-version"])[0])
		fmt.Fprintf(&b, "- **Anchor:** `%s`\n", values["anchor"])
		fmt.Fprintf(&b, "- **Materiality:** %s\n", values["materiality"])
		fmt.Fprintf(&b, "- **Evidence Level:** %s\n", values["evidence-level"])
		fmt.Fprintf(&b, "- **Source Fingerprint:** `%s`\n", values["source-fingerprint-snapshot"])
		b.WriteString("\n")
	}
	return []byte(b.String()), nil
}
