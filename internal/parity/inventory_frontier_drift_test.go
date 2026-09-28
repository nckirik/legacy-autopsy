package parity

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

// TestInventoryAndFrontierDrift checks the §4.3/§4.4/§4.5/§4.6 record schemas and
// closed enums against protocol.md.
func TestInventoryAndFrontierDrift(t *testing.T) {
	res := compileWithRegistry(t, "protocol/inventory-and-frontier.cdl")
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

	text43 := load("4.3. Source inventory denominator — `10-SOURCE-INVENTORY.md`")
	compareSets(t, "SOURCE-INVENTORY-RECORD", fields["SOURCE-INVENTORY-RECORD"], sectionLabels(text43, true))
	compareSets(t, "DISCOVERY-METHOD", enums["DISCOVERY-METHOD"], enumValuesFromFieldLine(text43, "discovery-method"))
	compareSets(t, "COVERAGE-SUMMARY", enums["COVERAGE-SUMMARY"], enumValuesFromFieldLine(text43, "coverage-summary"))
	compareSets(t, "INVENTORY-KIND", enums["INVENTORY-KIND"], inventoryKindsFromSentence(t, text43))

	text44 := load("4.4. Traversal frontier — `11-TRAVERSAL-FRONTIER.md`")
	compareSets(t, "FRONTIER-RECORD", fields["FRONTIER-RECORD"], sectionLabels(text44, false))
	compareSets(t, "FRONTIER-DISCOVERY-METHOD", enums["FRONTIER-DISCOVERY-METHOD"], enumValuesFromFieldLine(text44, "discovery-method"))
	compareSets(t, "FRONTIER-STATE", enums["FRONTIER-STATE"], enumValuesFromFieldLine(text44, "state"))

	text45 := load("4.5. Scope-specific source coverage — `16-SOURCE-COVERAGE.md`")
	var row []string
	for _, line := range strings.Split(text45, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "|") {
			cols := tableColumns(trimmed)
			if len(cols) > 0 && cols[0] == "cov-id" {
				row = cols
			}
		}
	}
	if len(row) == 0 {
		t.Fatal("§4.5 coverage table not found")
	}
	compareSets(t, "SOURCE-COVERAGE-ROW", fields["SOURCE-COVERAGE-ROW"], row)
	compareSets(t, "COVERAGE-SCOPE", enums["COVERAGE-SCOPE"], backtickedAfter(text45, "`Scope` is"))
	compareSets(t, "COVERAGE-DISPOSITION", enums["COVERAGE-DISPOSITION"], backtickedAfter(text45, "disposition is"))

	text46 := load("4.6. Stable supersession and tombstones")
	compareSets(t, "SUPERSESSION-TOMBSTONE", fields["SUPERSESSION-TOMBSTONE"], sectionLabels(text46, false))
}

// sectionLabels collects normalized field labels in a section, optionally mapping
// the "Track(s)" plural to a stable field name.
func sectionLabels(text string, inventory bool) []string {
	var labels []string
	for _, line := range strings.Split(text, "\n") {
		if match := logFieldLine.FindStringSubmatch(strings.TrimSpace(line)); match != nil {
			label := strings.ReplaceAll(match[1], "(s)", "s")
			labels = append(labels, normalizeLabel(label))
		}
	}
	_ = inventory
	return labels
}

// enumValuesFromFieldLine returns the pipe-separated values after a field label.
func enumValuesFromFieldLine(text, label string) []string {
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		match := logFieldLine.FindStringSubmatch(trimmed)
		if match == nil || normalizeLabel(match[1]) != label {
			continue
		}
		idx := strings.Index(trimmed, ":**")
		if idx < 0 {
			continue
		}
		tail := strings.NewReplacer("[", "", "]", "").Replace(trimmed[idx+3:])
		var values []string
		for _, value := range strings.Split(tail, "|") {
			values = append(values, strings.TrimSpace(value))
		}
		return values
	}
	return nil
}

func inventoryKindsFromSentence(t *testing.T, text string) []string {
	t.Helper()
	const prefix = "Inventory MUST cover, where applicable: "
	idx := strings.Index(text, prefix)
	if idx < 0 {
		t.Fatal("inventory kinds sentence not found")
	}
	tail := text[idx+len(prefix):]
	if cut := strings.Index(tail, "."); cut >= 0 {
		tail = tail[:cut]
	}
	values := strings.Split(tail, "; ")
	if n := len(values); n > 0 {
		values[n-1] = strings.TrimPrefix(values[n-1], "and ")
	}
	return values
}

func backtickedAfter(text, marker string) []string {
	idx := strings.Index(text, marker)
	if idx < 0 {
		return nil
	}
	tail := text[idx+len(marker):]
	if cut := strings.IndexAny(tail, ".;"); cut >= 0 {
		tail = tail[:cut]
	}
	var values []string
	for _, match := range backtickedValue.FindAllStringSubmatch(tail, -1) {
		values = append(values, match[1])
	}
	return values
}
