package parity

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

// TestSynthesisCatalogsDrift checks every §11.1–§11.13 record schema, table shape,
// and closed enum against protocol.md.
func TestSynthesisCatalogsDrift(t *testing.T) {
	res := compileWithRegistry(t, "examples/spec/synthesis-catalogs.cdl")
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

	text111 := load("11.1. Partial versus Final Synthesis")
	compareSets(t, "SYNTHESIS-GAP", fields["SYNTHESIS-GAP"], sectionLabels(text111, false))
	compareSets(t, "GAP-MATERIALITY", enums["GAP-MATERIALITY"], enumValuesFromFieldLine(text111, "materiality"))
	compareSets(t, "GAP-DISPOSITION", enums["GAP-DISPOSITION"], enumValuesFromFieldLine(text111, "disposition"))

	text112 := load("11.2. Common synthesis block header")
	payload, envelope := splitEnvelope(text112)
	compareSets(t, "SYNTHESIS-BLOCK-HEADER", fields["SYNTHESIS-BLOCK-HEADER"], payload)
	compareSets(t, "CERTIFICATION-ENVELOPE", fields["CERTIFICATION-ENVELOPE"], envelope)

	text113 := load("11.3. Architecture blueprint — `90-ARCH-BLUEPRINT.md`")
	compareSets(t, "ARCHITECTURE-MODULE", fields["ARCHITECTURE-MODULE"], labelsUnder(text113, "MOD-ID"))

	text114 := load("11.4. Entity — `91-DATA-MODEL.md`")
	compareSets(t, "ENTITY-RECORD", fields["ENTITY-RECORD"], labelsUnder(text114, "ER-ID"))
	compareSets(t, "ENTITY-COLUMN", fields["ENTITY-COLUMN"], labeledTable(text114, "columns-fields"))
	compareSets(t, "STORAGE-KIND", enums["STORAGE-KIND"], enumValuesFromFieldLine(text114, "storage"))

	text115 := load("11.5. Relationship — `91-DATA-MODEL.md`")
	compareSets(t, "RELATIONSHIP-RECORD", fields["RELATIONSHIP-RECORD"], labelsUnder(text115, "REL-ID"))
	compareSets(t, "RELATIONSHIP-CARDINALITY", enums["RELATIONSHIP-CARDINALITY"], enumValuesFromFieldLine(text115, "cardinality"))
	compareSets(t, "RELATIONSHIP-ENFORCEMENT", enums["RELATIONSHIP-ENFORCEMENT"], enumValuesFromFieldLine(text115, "enforcement"))
	compareSets(t, "DELETE-UPDATE-SEMANTICS", enums["DELETE-UPDATE-SEMANTICS"], enumValuesFromFieldLine(text115, "delete-update-semantics"))

	text116 := load("11.6. State machine — `91-DATA-MODEL.md`")
	compareSets(t, "STATE-MACHINE", fields["STATE-MACHINE"], labelsUnder(text116, "SM-ID"))
	compareSets(t, "STATE-TRANSITION", fields["STATE-TRANSITION"], labeledTable(text116, "transitions"))

	text117 := load("11.7. Database routine — `91-DATA-MODEL.md`")
	compareSets(t, "DATABASE-ROUTINE", fields["DATABASE-ROUTINE"], labelsUnder(text117, "DR-ID"))
	compareSets(t, "DB-ROUTINE-KIND", enums["DB-ROUTINE-KIND"], enumValuesFromFieldLine(text117, "kind"))
	compareSets(t, "DB-TARGET-DECISION", enums["DB-TARGET-DECISION"], enumValuesFromFieldLine(text117, "target-decision"))

	text118 := load("11.8. Business rule — `92-BUSINESS-RULES.md`")
	compareSets(t, "BUSINESS-RULE-RECORD", fields["BUSINESS-RULE-RECORD"], labelsUnder(text118, "BR-ID"))
	compareSets(t, "VIOLATION-POLICY", enums["VIOLATION-POLICY"], enumValuesFromFieldLine(text118, "violation-policy"))
	compareSets(t, "BUSINESS-RULE-CAPABILITY-STATE", enums["BUSINESS-RULE-CAPABILITY-STATE"], enumValuesFromFieldLine(text118, "capability-state"))

	text119 := load("11.9. Use case — `93-USE-CASES.md`")
	compareSets(t, "USE-CASE", fields["USE-CASE"], labelsUnder(text119, "UC-ID"))
	compareSets(t, "USE-CASE-STEP", fields["USE-CASE-STEP"], tableAfterHeading(text119, "#### Main Flow"))
	compareSets(t, "USE-CASE-CAPABILITY-CLASSIFICATION", enums["USE-CASE-CAPABILITY-CLASSIFICATION"], enumValuesFromFieldLine(text119, "capability-classification"))
	compareSets(t, "USE-CASE-TARGET-DECISION", enums["USE-CASE-TARGET-DECISION"], enumValuesFromFieldLine(text119, "target-decision"))

	text1110 := load("11.10. Interface — `94-INTERFACES.md`")
	compareSets(t, "INTERFACE-RECORD", fields["INTERFACE-RECORD"], labelsUnder(text1110, "IF-ID"))
	compareSets(t, "INTERFACE-REQUEST-FIELD", fields["INTERFACE-REQUEST-FIELD"], labeledTable(text1110, "request-schema"))
	compareSets(t, "INTERFACE-RESPONSE-FIELD", fields["INTERFACE-RESPONSE-FIELD"], labeledTable(text1110, "response-schema"))

	text1111 := load("11.11. Deployment, configuration, and scheduling — `95-DEPLOYMENT.md`")
	compareSets(t, "DEPLOYMENT-ELEMENT", fields["DEPLOYMENT-ELEMENT"], labelsUnder(text1111, "DEP-ID"))
	compareSets(t, "CONFIGURATION-RECORD", fields["CONFIGURATION-RECORD"], labelsUnder(text1111, "CFG-ID"))
	compareSets(t, "SCHEDULE-RECORD", fields["SCHEDULE-RECORD"], labelsUnder(text1111, "SCHED-ID"))
	compareSets(t, "CONFIG-MEDIUM", enums["CONFIG-MEDIUM"], enumValuesFromFieldLine(text1111, "medium"))

	text1112 := load("11.12. NFR and security — `96-NON-FUNCTIONAL-SECURITY.md`")
	compareSets(t, "NFR-RECORD", fields["NFR-RECORD"], labelsUnder(text1112, "NFR-ID"))
	compareSets(t, "SECURITY-FINDING", fields["SECURITY-FINDING"], labelsUnder(text1112, "SEC-ID"))
	compareSets(t, "FAULT-RECORD", fields["FAULT-RECORD"], labelsUnder(text1112, "FLT-ID"))
	compareSets(t, "NFR-ATTRIBUTE", enums["NFR-ATTRIBUTE"], enumValuesFromFieldLine(text1112, "attribute"))
	compareSets(t, "SECURITY-CATEGORY", enums["SECURITY-CATEGORY"], enumValuesFromFieldLine(text1112, "category"))
	compareSets(t, "RISK-LEVEL", enums["RISK-LEVEL"], enumValuesFromFieldLine(text1112, "risk-level"))

	text1113 := load("11.13. Persona profile")
	profile, profileEnvelope := splitEnvelope(text1113)
	compareSets(t, "PERSONA-PROFILE", fields["PERSONA-PROFILE"], profile)
	compareSets(t, "CERTIFICATION-ENVELOPE", fields["CERTIFICATION-ENVELOPE"], profileEnvelope)
	compareSets(t, "AUTHORIZATION-MATRIX-ROW", fields["AUTHORIZATION-MATRIX-ROW"], labeledTable(text1113, "authorization-matrix"))
}

// labelsUnder collects labels inside a "### [TOKEN]" block until the next heading.
func labelsUnder(text, token string) []string {
	prefix := "### [" + token + "]"
	var labels []string
	active := false
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "### [") {
			active = strings.HasPrefix(trimmed, prefix)
			continue
		}
		if !active {
			continue
		}
		if match := logFieldLine.FindStringSubmatch(trimmed); match != nil {
			label := strings.ReplaceAll(match[1], "(s)", "s")
			labels = append(labels, normalizeLabel(label))
		}
	}
	if len(labels) == 0 {
		panic("no labels under " + token)
	}
	return labels
}

// labeledTable returns the table header columns following a normalized field label.
func labeledTable(text, label string) []string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		match := logFieldLine.FindStringSubmatch(trimmed)
		if match == nil || normalizeLabel(match[1]) != label {
			continue
		}
		for _, next := range lines[i+1:] {
			next = strings.TrimSpace(next)
			if strings.HasPrefix(next, "|") {
				return tableColumns(next)
			}
			if strings.HasPrefix(next, "- **") || strings.HasPrefix(next, "###") {
				break
			}
		}
	}
	panic("no table under label " + label)
}

// tableAfterHeading returns the first table header following an exact heading line.
func tableAfterHeading(text, heading string) []string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) != heading {
			continue
		}
		for _, next := range lines[i+1:] {
			next = strings.TrimSpace(next)
			if strings.HasPrefix(next, "|") {
				return tableColumns(next)
			}
		}
	}
	panic("no table after heading " + heading)
}

// splitEnvelope splits field labels into payload and certification-envelope parts.
func splitEnvelope(text string) ([]string, []string) {
	var payload, envelope []string
	inEnvelope := false
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "CERTIFICATION-ENVELOPE") && strings.HasPrefix(trimmed, "<!--") && !strings.Contains(trimmed, "END-CERTIFICATION-ENVELOPE") {
			inEnvelope = true
			continue
		}
		if match := logFieldLine.FindStringSubmatch(trimmed); match != nil {
			if inEnvelope {
				envelope = append(envelope, normalizeLabel(match[1]))
			} else {
				payload = append(payload, normalizeLabel(match[1]))
			}
		}
	}
	if len(payload) == 0 || len(envelope) == 0 {
		panic("envelope split failed")
	}
	return payload, envelope
}
