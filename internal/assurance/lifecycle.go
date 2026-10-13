package assurance

import (
	"fmt"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// LifecyclePlan builds one lifecycle-eligibility unit per source-unit batch.
// Each unit authors one LIFECYCLE-ELIGIBILITY record per source unit with the
// classification, reason, and the claim evidence that supports it.
func LifecyclePlan(db string, doc *cdl.EIRDoc, namespace, snapshot string, batchSize int) ([]work.Unit, error) {
	if batchSize <= 0 {
		batchSize = 10
	}
	inventory, err := LoadInventory(db, doc)
	if err != nil {
		return nil, err
	}
	claimsBySubject, err := claimsBySubject(db, doc)
	if err != nil {
		return nil, err
	}
	claimNames, _, _, err := work.FieldContract(doc, "CLAIM")
	if err != nil {
		return nil, err
	}
	claimText := map[string]string{}
	for _, ids := range claimsBySubject {
		for _, id := range ids {
			payload, err := store.Get(db, doc, "CLAIM", id)
			if err != nil {
				return nil, err
			}
			values, err := work.PayloadText(claimNames, payload)
			if err != nil {
				return nil, err
			}
			claimText[id] = strings.Join(strings.Fields(values["claim"]), " ")
		}
	}
	var units []work.Unit
	for start := 0; start < len(inventory); start += batchSize {
		end := start + batchSize
		if end > len(inventory) {
			end = len(inventory)
		}
		batch := inventory[start:end]
		var b strings.Builder
		b.WriteString("# Lifecycle eligibility batch (runtime-generated, bounded)\n\n")
		b.WriteString("Author exactly one LIFECYCLE-ELIGIBILITY record per source unit, in order.\n\n")
		b.WriteString("## Runtime directives (must be used exactly)\n\n")
		fmt.Fprintf(&b, "- **System Namespace:** %s\n", namespace)
		fmt.Fprintf(&b, "- **Environment/Snapshot:** %s\n", snapshot)
		b.WriteString("- `scope-id` is the source unit ID exactly as shown.\n")
		b.WriteString("- `classification` is Eligible, Ineligible, or Unknown.\n")
		b.WriteString("- `reason` is one short sentence grounded in the cited claims.\n")
		b.WriteString("- `clm-evidence-bindings` lists the claim IDs you used, comma-separated, or None.\n")
		b.WriteString("- `classified-by-invocation` is runtime-assigned: emit the field with placeholder value PENDING.\n")
		b.WriteString("- `snapshot` is the snapshot above.\n\n")
		for _, item := range batch {
			fmt.Fprintf(&b, "## Source unit %s\n\n- **Concrete Coordinate:** %s\n", item.ID, item.Coordinate)
			claims := claimsBySubject[item.ID]
			if len(claims) == 0 {
				b.WriteString("- **Claims:** None\n\n")
				continue
			}
			b.WriteString("- **Claims:**\n")
			for _, claim := range claims {
				fmt.Fprintf(&b, "  - %s: %s\n", claim, claimText[claim])
			}
			b.WriteString("\n")
		}
		reads := make([]work.Ref, 0, len(batch))
		for _, item := range batch {
			reads = append(reads, work.Ref{FieldID: "SOURCE-INVENTORY-RECORD", RecordID: item.ID})
		}
		units = append(units, work.Unit{
			ID:              fmt.Sprintf("U-ASSURANCE-LIFECYCLE-%03d", start/batchSize+1),
			Mode:            "Discovery",
			Scope:           namespace,
			IDFromFields:    []string{"scope-id"},
			RuntimeFields:   []string{"classified-by-invocation"},
			ExpectedRecords: len(batch),
			ResultType:      "LIFECYCLE-ELIGIBILITY",
			Reads:           reads,
			Inputs:          []work.Input{{Name: "lifecycle-context.md", Bytes: []byte(b.String())}},
			Sections:        []string{"dependency-provenance"},
			ReasoningEffort: "off",
		})
	}
	return units, nil
}

// LifecycleCoverage reports how many source units have an eligibility record and
// how many remain Unknown.
func LifecycleCoverage(db string, doc *cdl.EIRDoc) (total, covered, unknown int, err error) {
	inventory, err := LoadInventory(db, doc)
	if err != nil {
		return 0, 0, 0, err
	}
	names, _, _, err := work.FieldContract(doc, "LIFECYCLE-ELIGIBILITY")
	if err != nil {
		return 0, 0, 0, err
	}
	ids, err := store.List(db, doc, "LIFECYCLE-ELIGIBILITY")
	if err != nil {
		return 0, 0, 0, err
	}
	byScope := map[string]string{}
	for _, id := range ids {
		payload, err := store.Get(db, doc, "LIFECYCLE-ELIGIBILITY", id)
		if err != nil {
			return 0, 0, 0, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return 0, 0, 0, err
		}
		byScope[firstTokenField(values["scope-id"])] = values["classification"]
	}
	total = len(inventory)
	for _, item := range inventory {
		classification, ok := byScope[item.ID]
		if !ok {
			continue
		}
		covered++
		if classification == "Unknown" || classification == "" {
			unknown++
		}
	}
	return total, covered, unknown, nil
}
