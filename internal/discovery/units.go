package discovery

import (
	"fmt"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/work"
)

// Cluster is the bounded cluster context handed to discovery units.
type Cluster struct {
	ID         string
	Title      string
	Owner      string
	ScopePaths []string
}

// Directives are runtime-owned discovery values.
type Directives struct {
	Namespace      string
	Snapshot       string
	StorageBinding string
	BatchSize      int
}

// Plan builds one bounded inventory unit per candidate batch. Each unit must
// author exactly one SOURCE-INVENTORY-RECORD per candidate, in order; the
// runtime computes every SRC ID from the payload kind and coordinate.
func Plan(candidates []Candidate, clusters []Cluster, directives Directives) []work.Unit {
	batchSize := directives.BatchSize
	if batchSize <= 0 {
		batchSize = 25
	}
	var units []work.Unit
	for start := 0; start < len(candidates); start += batchSize {
		end := start + batchSize
		if end > len(candidates) {
			end = len(candidates)
		}
		batch := candidates[start:end]
		units = append(units, work.Unit{
			ID:               fmt.Sprintf("U-DISCOVERY-INVENTORY-%03d", start/batchSize+1),
			Mode:             "Discovery",
			Scope:            directives.Namespace,
			IDFromCoordinate: "concrete-coordinate",
			IDType:           "SRC",
			Namespace:        directives.Namespace,
			ExpectedRecords:  len(batch),
			ResultType:       "SOURCE-INVENTORY-RECORD",
			Reads:            []work.Ref{{FieldID: "PREFLIGHT-HEADER", RecordID: "0A"}},
			Inputs:           []work.Input{{Name: "discovery-context.md", Bytes: contextMarkdown(batch, clusters, directives)}},
			Sections:         []string{"atomic-unit-rule", "source-inventory"},
			ReasoningEffort:  "off",
		})
	}
	return units
}

func contextMarkdown(batch []Candidate, clusters []Cluster, directives Directives) []byte {
	var b strings.Builder
	b.WriteString("# Discovery inventory batch (runtime-generated, bounded)\n\n")
	b.WriteString("## Candidates - author exactly one SOURCE-INVENTORY-RECORD per row, in this order\n\n")
	b.WriteString("| Concrete Coordinate | Size (bytes) | File Fingerprint |\n| :-- | --: | :-- |\n")
	for _, candidate := range batch {
		fmt.Fprintf(&b, "| %s | %d | %s |\n", candidate.Path, candidate.Size, candidate.Fingerprint)
	}
	b.WriteString("\n## Entry clusters (use the matching cluster's owner persona and scope)\n\n")
	if len(clusters) == 0 {
		b.WriteString("- None\n")
	} else {
		for _, cluster := range clusters {
			fmt.Fprintf(&b, "- %s (%s): owner %s, scope %s\n", cluster.ID, cluster.Title, cluster.Owner, strings.Join(cluster.ScopePaths, ", "))
		}
	}
	b.WriteString("\n## Runtime directives (must be used exactly)\n\n")
	fmt.Fprintf(&b, "- **System Namespace:** %s\n", directives.Namespace)
	fmt.Fprintf(&b, "- **Environment/Snapshot:** %s\n", directives.Snapshot)
	fmt.Fprintf(&b, "- **Canonical Storage Binding:** %s\n", directives.StorageBinding)
	b.WriteString("- `concrete-coordinate` must be copied exactly from the candidate row.\n")
	b.WriteString("- `kind` is FILE for ordinary files; use the most specific SRC kind only with evidence.\n")
	b.WriteString("- `discovery-method` is `filesystem` for every row.\n")
	b.WriteString("- `source-fingerprint-snapshot-src-version` is the row fingerprint plus ` version 1`.\n")
	b.WriteString("- `parent-container-src`, `mappings`, `frontier-references`, `coverage-rows`, and `claims` are `None` unless evidence-backed.\n")
	b.WriteString("- `canonical-owner` is the owner persona prefix of the cluster whose scope contains the coordinate; otherwise `None`.\n")
	b.WriteString("- `disposition-summary` is `Unknown` (fail-inclusive) until swept.\n")
	b.WriteString("- `coverage-summary` is the closed enum value `C-GAP` for newly discovered unswept units.\n")
	b.WriteString("- `traversal-tracks` is `Normal` unless the scan evidence shows otherwise.\n")
	return []byte(b.String())
}
