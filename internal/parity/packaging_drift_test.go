package parity

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

type artifactBlock struct {
	payload  []string
	envelope []string
	tables   map[string][]string
}

// TestPackagingDrift checks every §15 packaging schema, check registry, and table
// shape against protocol.md.
func TestPackagingDrift(t *testing.T) {
	res := compileWithRegistry(t, "protocol/packaging.cdl")
	fields := map[string][]string{}
	for _, f := range res.EIR.Declarations.Fields {
		names := make([]string, 0, len(f.Fields))
		for _, spec := range f.Fields {
			names = append(names, spec.Name)
		}
		fields[f.ID] = names
	}
	enums := map[string][]string{}
	for _, e := range res.EIR.Declarations.Enums {
		enums[e.ID] = e.Values
	}

	model, err := protocol.Load(protocolPathIn(root(t)))
	if err != nil {
		t.Fatal(err)
	}
	load := func(heading string) string {
		text, err := model.SectionText(heading)
		if err != nil {
			t.Fatal(err)
		}
		return text
	}

	artifacts := map[string]artifactBlock{
		"EXIT-E-CANDIDATE-REPORT":         parseArtifactBlock(t, load("15.1.2. Exit E Candidate Report ordered payload schema"), "EXIT-E-CANDIDATE-REPORT"),
		"CANDIDATE-PAYLOAD-MANIFEST":      parseArtifactBlock(t, load("15.1.3. Candidate Payload Manifest ordered payload schema"), "CANDIDATE-PAYLOAD-MANIFEST"),
		"EXIT-E-CONTENT-READINESS-REPORT": parseArtifactBlock(t, load("15.1.4. Exit E Content-Readiness Report ordered payload schema"), "EXIT-E-CONTENT-READINESS-REPORT"),
		"SCOPE-CERTIFICATE":               parseArtifactBlock(t, load("15.2. Scope certificate"), "SCOPE-CERTIFICATE"),
		"OUTER-BUNDLE-MANIFEST":           parseArtifactBlock(t, load("15.4. Acyclic certification sequence and hash domains"), "OUTER-BUNDLE-MANIFEST"),
	}
	for id, block := range artifacts {
		compareSets(t, id, fields[id], block.payload)
		if id == "OUTER-BUNDLE-MANIFEST" {
			compareSets(t, "OUTER-SIGNATURE-ENVELOPE", fields["OUTER-SIGNATURE-ENVELOPE"], block.envelope)
		} else {
			compareSets(t, "ARTIFACT-SIGNATURE-ENVELOPE", fields["ARTIFACT-SIGNATURE-ENVELOPE"], block.envelope)
		}
	}

	tableBlocks := map[string]string{
		"semantic-record-counts":                 "SEMANTIC-RECORD-COUNT",
		"decision-content-counts":                "DECISION-CONTENT-COUNT",
		"check-results":                          "CHECK-RESULT",
		"blockers":                               "BLOCKER",
		"semantic-bindings":                      "SEMANTIC-BINDING",
		"dec-content-bindings":                   "DEC-CONTENT-BINDING",
		"confirmed-semantic-bindings":            "CONFIRMED-SEMANTIC-BINDING",
		"approved-dec-bindings":                  "APPROVED-DEC-BINDING",
		"pre-certificate-package-members":        "PACKAGE-MEMBER",
		"included-scope-dimensions":              "INCLUDED-SCOPE-DIMENSION",
		"source-counts-by-kind":                  "SOURCE-COUNT",
		"approved-exclusions-and-residual-risks": "APPROVED-EXCLUSION-RESIDUAL-RISK",
		"confirmed-semantic-payload-bindings":    "CONFIRMED-SEMANTIC-PAYLOAD-BINDING",
		"known-residual-risks":                   "KNOWN-RESIDUAL-RISK",
		"package-members":                        "PACKAGE-MEMBER",
	}
	seenTables := 0
	for _, block := range artifacts {
		for label, columns := range block.tables {
			fieldID, ok := tableBlocks[label]
			if !ok {
				t.Fatalf("unmapped table label %q", label)
			}
			compareSets(t, fieldID, fields[fieldID], columns)
			seenTables++
		}
	}
	if seenTables == 0 {
		t.Fatal("no table shapes parsed")
	}

	text151 := load("15.1.1. Universal deterministic packaging rules")
	var candidateChecks, readinessChecks []string
	for _, row := range rawTableRows(text151) {
		if len(row) < 2 {
			continue
		}
		switch row[0] {
		case "CANDIDATE":
			candidateChecks = append(candidateChecks, row[1])
		case "CONTENT-READINESS":
			readinessChecks = append(readinessChecks, row[1])
		}
	}
	if len(candidateChecks) == 0 || len(readinessChecks) == 0 {
		t.Fatal("check registry rows not found in §15.1.1")
	}
	compareSets(t, "EXIT-E-CANDIDATE-CHECKS", enums["EXIT-E-CANDIDATE-CHECKS"], candidateChecks)
	compareSets(t, "EXIT-E-CONTENT-READINESS-CHECKS", enums["EXIT-E-CONTENT-READINESS-CHECKS"], readinessChecks)

	text155 := load("15.5. Reproducible final-bundle verification")
	var finalChecks []string
	for _, row := range rawTableRows(text155) {
		if len(row) >= 1 && strings.HasPrefix(row[0], "FINAL-") {
			finalChecks = append(finalChecks, row[0])
		}
	}
	if len(finalChecks) == 0 {
		t.Fatal("final check registry rows not found in §15.5")
	}
	compareSets(t, "FINAL-CHECK", enums["FINAL-CHECK"], finalChecks)
	compareSets(t, "FINAL-CHECK-RESULT", fields["FINAL-CHECK-RESULT"], tableColumns("| Check ID | Result | Offending Bindings |"))

	var receipt []string
	inFence := false
	for _, line := range strings.Split(text155, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inFence = !inFence
			continue
		}
		if !inFence {
			continue
		}
		if match := logFieldLine.FindStringSubmatch(trimmed); match != nil {
			receipt = append(receipt, normalizeLabel(match[1]))
		}
	}
	if len(receipt) == 0 {
		t.Fatal("final verification receipt fields not found")
	}
	compareSets(t, "FINAL-VERIFICATION-RECEIPT", fields["FINAL-VERIFICATION-RECEIPT"], receipt)
}

// parseArtifactBlock extracts the ordered payload fields, envelope fields, and
// labeled table shapes for one artifact type from a section's text.
func parseArtifactBlock(t *testing.T, text, artifactType string) artifactBlock {
	t.Helper()
	block := artifactBlock{tables: map[string][]string{}}
	mode := ""
	lastLabel := ""
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.Contains(trimmed, "ARTIFACT-PAYLOAD: "+artifactType):
			mode = "payload"
			continue
		case strings.HasPrefix(trimmed, "<!-- END-ARTIFACT-PAYLOAD"):
			mode = ""
			continue
		case strings.Contains(trimmed, "ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: "+artifactType):
			mode = "envelope"
			continue
		case strings.HasPrefix(trimmed, "<!-- END-ARTIFACT-DIGEST-SIGNATURE-ENVELOPE"):
			mode = ""
			continue
		}
		if mode == "" {
			continue
		}
		if match := logFieldLine.FindStringSubmatch(trimmed); match != nil {
			label := normalizeLabel(match[1])
			if mode == "payload" {
				block.payload = append(block.payload, label)
			} else {
				block.envelope = append(block.envelope, label)
			}
			lastLabel = label
			continue
		}
		if mode == "payload" && strings.HasPrefix(trimmed, "|") && lastLabel != "" {
			cols := tableColumns(trimmed)
			if len(cols) > 0 && cols[0] != "" && !strings.HasPrefix(cols[0], ":") {
				block.tables[lastLabel] = cols
				lastLabel = ""
			}
		}
	}
	if len(block.payload) == 0 || len(block.envelope) == 0 {
		t.Fatalf("artifact %s: parsed %d payload and %d envelope fields", artifactType, len(block.payload), len(block.envelope))
	}
	return block
}

func rawTableRows(text string) [][]string {
	var rows [][]string
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "|") || strings.Contains(trimmed, "---") {
			continue
		}
		var cells []string
		for _, cell := range strings.Split(strings.Trim(trimmed, "|"), "|") {
			cells = append(cells, strings.TrimSpace(cell))
		}
		rows = append(rows, cells)
	}
	return rows
}
