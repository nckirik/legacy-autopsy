// Package assurance runs the Assurance pipeline: evidence-bound claims over
// committed source-inventory records with runtime-computed CLM identities.
package assurance

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// Inventory is one committed source unit in bounded form.
type Inventory struct {
	ID          string
	Coordinate  string
	Fingerprint string
}

// LoadInventory reads every committed SOURCE-INVENTORY-RECORD into its bounded
// coordinate/fingerprint form.
func LoadInventory(db string, doc *cdl.EIRDoc) ([]Inventory, error) {
	names, _, _, err := work.FieldContract(doc, "SOURCE-INVENTORY-RECORD")
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, "SOURCE-INVENTORY-RECORD")
	if err != nil {
		return nil, err
	}
	inventory := make([]Inventory, 0, len(ids))
	for _, id := range ids {
		payload, err := store.Get(db, doc, "SOURCE-INVENTORY-RECORD", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		inventory = append(inventory, Inventory{
			ID:          id,
			Coordinate:  values["concrete-coordinate"],
			Fingerprint: firstTokenField(values["source-fingerprint-snapshot-src-version"]),
		})
	}
	return inventory, nil
}

// Directives are runtime-owned assurance values.
type Directives struct {
	Namespace       string
	Snapshot        string
	StorageBinding  string
	BatchSize       int
	MaxExcerptBytes int
}

// Plan builds one bounded claim unit per inventory batch. Each unit must author
// exactly one CLAIM per row, in order, anchored at the row coordinate.
func Plan(root string, inventory []Inventory, directives Directives) ([]work.Unit, error) {
	batchSize := directives.BatchSize
	if batchSize <= 0 {
		batchSize = 10
	}
	maxExcerpt := directives.MaxExcerptBytes
	if maxExcerpt <= 0 {
		maxExcerpt = 1200
	}
	var units []work.Unit
	for start := 0; start < len(inventory); start += batchSize {
		end := start + batchSize
		if end > len(inventory) {
			end = len(inventory)
		}
		batch := inventory[start:end]
		context, err := claimContext(root, batch, directives, maxExcerpt)
		if err != nil {
			return nil, err
		}
		reads := make([]work.Ref, 0, len(batch))
		for _, item := range batch {
			reads = append(reads, work.Ref{FieldID: "SOURCE-INVENTORY-RECORD", RecordID: item.ID})
		}
		units = append(units, work.Unit{
			ID:           fmt.Sprintf("U-ASSURANCE-CLAIMS-%03d", start/batchSize+1),
			Mode:         "Discovery",
			Scope:        directives.Namespace,
			IDType:       "CLM",
			Namespace:    directives.Namespace,
			IDOwnerField: "subject-id-record-version",
			// Identity anchors in evidence: the claim's source coordinate (its
			// anchor), not its prose. Rewording a claim preserves its ID; a
			// different claim at the same anchor conflicts explicitly.
			IDDiscriminatorField: "anchor",
			ExpectedRecords:      len(batch),
			ResultType:           "CLAIM",
			Reads:                reads,
			Inputs:               []work.Input{{Name: "claim-evidence.md", Bytes: context}},
			Sections:             []string{"claim-evidence"},
			ReasoningEffort:      "off",
		})
	}
	return units, nil
}

func claimContext(root string, batch []Inventory, directives Directives, maxExcerpt int) ([]byte, error) {
	var b strings.Builder
	b.WriteString("# Claim evidence batch (runtime-generated, bounded)\n\n")
	b.WriteString("Author exactly one CLAIM per source row, in order. Each claim states one falsifiable fact supported by the excerpt.\n\n")
	for _, item := range batch {
		fmt.Fprintf(&b, "## Source unit %s\n\n", item.ID)
		fmt.Fprintf(&b, "- **Concrete Coordinate:** %s\n", item.Coordinate)
		fmt.Fprintf(&b, "- **Fingerprint:** %s\n", item.Fingerprint)
		excerpt := boundedExcerpt(root, item.Coordinate, maxExcerpt)
		fmt.Fprintf(&b, "\n```text\n%s\n```\n\n", excerpt)
	}
	b.WriteString("## Runtime directives (must be used exactly)\n\n")
	fmt.Fprintf(&b, "- **System Namespace:** %s\n", directives.Namespace)
	fmt.Fprintf(&b, "- **Environment/Snapshot:** %s\n", directives.Snapshot)
	b.WriteString("- `subject-id-record-version` is the source unit ID exactly as shown.\n")
	b.WriteString("- `anchor` is `coordinate:line` where the line is the excerpt line (1-based) that supports the claim.\n")
	b.WriteString("- `source-fingerprint-snapshot` is the row fingerprint exactly as shown.\n")
	b.WriteString("- `materiality` is Material or Supporting; `evidence-level` is one of [E-DIRECT], [E-RUNTIME-OBSERVED], [E-CROSS-VALIDATED], [E-INFERRED], [E-UNKNOWN] with brackets.\n")
	b.WriteString("- `evidence-route-lineage-id` is `static-filesystem`.\n")
	b.WriteString("- `projection-classification` is `raw-visible`.\n")
	b.WriteString("- `independent-supporting-claims`, `contradiction-ids` are `None`; `independence-basis` is `None`.\n")
	b.WriteString("- `created-last-revalidated-by` is `None`.\n")
	return []byte(b.String()), nil
}

func boundedExcerpt(root, coordinate string, maxBytes int) string {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(coordinate)))
	if err != nil {
		return "(unreadable: " + err.Error() + ")"
	}
	text := string(data)
	if len(text) > maxBytes {
		text = text[:maxBytes] + "\n...[truncated]"
	}
	return text
}

func firstTokenField(value string) string {
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}
