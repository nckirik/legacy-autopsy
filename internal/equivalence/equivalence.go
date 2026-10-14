// Package equivalence generates the deterministic acceptance-test skeleton
// suite: one reviewable skeleton per use case plus property and baseline stubs
// for business rules, state machines, and external interfaces, with an explicit
// mapping table.
package equivalence

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// SuitePath is the package-relative path of the generated suite.
const SuitePath = "equivalence/acceptance-suite.md"

// ReviewPath is the package-relative path of the human review record.
const ReviewPath = "equivalence/equivalence-review.md"

var typedTokenPattern = regexp.MustCompile(`\b(CLM|BR|UC|IF|SCHED|DEP|CMP|ER|MOD|SM|FLT|SEC|NFR|REL|CFG|PRF|PRV|CAP|DR)-[A-F0-9]{12,64}\b`)

// Mapping is one suite mapping row with its skeleton body.
type Mapping struct {
	RecordType string
	RecordID   string
	Skeleton   string
	Body       string
}

// Suite is the generated acceptance-test skeleton document.
type Suite struct {
	Markdown []byte
	Mappings []Mapping
}

// Build generates the acceptance-test skeleton suite from committed records.
func Build(db string, doc *cdl.EIRDoc) (Suite, error) {
	type family struct {
		fieldID  string
		idType   string
		skeleton func(id string, values map[string]string) string
	}
	families := []family{
		{fieldID: "USE-CASE", idType: "UC", skeleton: useCaseSkeleton},
		{fieldID: "BUSINESS-RULE-RECORD", idType: "BR", skeleton: businessRuleSkeleton},
		{fieldID: "STATE-MACHINE", idType: "SM", skeleton: stateMachineSkeleton},
		{fieldID: "INTERFACE-RECORD", idType: "IF", skeleton: interfaceSkeleton},
	}
	var suite Suite
	var body strings.Builder
	body.WriteString("# Acceptance and Equivalence Suite\n\n")
	body.WriteString("Illustrative skeletons derived from committed records. Human review replaces the placeholders with executable checks or records an approved exclusion. The mapping table is the coverage denominator: every use case, business rule, state machine, and external interface appears exactly once.\n\n")
	body.WriteString("## Coverage Mapping\n\n")
	body.WriteString("| Record Type | Typed Record ID | Skeleton |\n")
	body.WriteString("| --- | --- | --- |\n")
	var stubs []string
	for _, item := range families {
		names, _, _, err := work.FieldContract(doc, item.fieldID)
		if err != nil {
			return suite, err
		}
		ids, err := store.List(db, doc, item.fieldID)
		if err != nil {
			return suite, err
		}
		sort.Strings(ids)
		for _, id := range ids {
			payload, err := store.Get(db, doc, item.fieldID, id)
			if err != nil {
				return suite, err
			}
			values, err := work.PayloadText(names, payload)
			if err != nil {
				return suite, err
			}
			anchor := strings.ToLower(item.idType + "-" + strings.TrimPrefix(id, item.idType+"-"))
			suite.Mappings = append(suite.Mappings, Mapping{RecordType: item.idType, RecordID: id, Skeleton: anchor})
			stubs = append(stubs, item.skeleton(anchor, values))
		}
	}
	for _, mapping := range suite.Mappings {
		fmt.Fprintf(&body, "| %s | %s | #%s |\n", mapping.RecordType, mapping.RecordID, mapping.Skeleton)
	}
	body.WriteString("\n")
	for _, stub := range stubs {
		body.WriteString(stub)
	}
	suite.Markdown = []byte(body.String())
	return suite, nil
}

func useCaseSkeleton(anchor string, values map[string]string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "### %s\n\n", anchor)
	writeSubject(&b, "Use case", values)
	b.WriteString("- **Goal:** " + condense(values["goal"]) + "\n")
	b.WriteString("- **Preconditions / Postconditions:** " + condense(values["preconditions-postconditions"]) + "\n")
	b.WriteString("- **Observable Contract:** " + condense(values["behavioral-contract"]) + "\n")
	b.WriteString("- **Related IDs:** " + joinIDs(relatedIDs(values)) + "\n\n")
	b.WriteString("- [ ] Happy path matches the observable contract.\n")
	b.WriteString("- [ ] Each precondition is enforced; each postcondition is observable.\n")
	b.WriteString("- [ ] Failure path returns the documented outcome.\n\n")
	return b.String()
}

// writeSubject renders the record identity for the human reviewer: typed ID and
// record title, never a claim token.
func writeSubject(b *strings.Builder, label string, values map[string]string) {
	id := firstToken(values["block-id"])
	if id == "None" || id == "" {
		id = joinIDs(relatedIDs(values))
	}
	title := condense(values["record-title"])
	fmt.Fprintf(b, "- **%s:** %s - %s\n", label, id, title)
}

func businessRuleSkeleton(anchor string, values map[string]string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "### %s\n\n", anchor)
	writeSubject(&b, "Business rule", values)
	b.WriteString("- **Condition:** " + condense(values["condition"]) + "\n")
	b.WriteString("- **Related IDs:** " + joinIDs(relatedIDs(values)) + "\n\n")
	b.WriteString("- [ ] Property holds for every accepted input.\n")
	b.WriteString("- [ ] Rejected input produces the documented rule outcome.\n\n")
	return b.String()
}

func stateMachineSkeleton(anchor string, values map[string]string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "### %s\n\n", anchor)
	writeSubject(&b, "State machine", values)
	b.WriteString("- **States:** " + condense(values["states-meaning"]) + "\n")
	b.WriteString("- **Related IDs:** " + joinIDs(relatedIDs(values)) + "\n\n")
	b.WriteString("- [ ] Every eligible transition is exercised once.\n")
	b.WriteString("- [ ] Invalid transitions are rejected.\n\n")
	return b.String()
}

func interfaceSkeleton(anchor string, values map[string]string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "### %s\n\n", anchor)
	writeSubject(&b, "Interface", values)
	b.WriteString("- **Baseline:** " + condense(values["compatibility-baseline"]) + "\n")
	b.WriteString("- **Related IDs:** " + joinIDs(relatedIDs(values)) + "\n\n")
	b.WriteString("- [ ] Request/response shape matches the baseline.\n")
	b.WriteString("- [ ] Migration rule is exercised.\n\n")
	return b.String()
}

func relatedIDs(values map[string]string) []string {
	seen := map[string]bool{}
	var ids []string
	for _, value := range values {
		for _, token := range typedTokenPattern.FindAllString(value, -1) {
			if seen[token] {
				continue
			}
			seen[token] = true
			ids = append(ids, token)
		}
	}
	sort.Strings(ids)
	return ids
}

func firstToken(value string) string {
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return "None"
	}
	return fields[0]
}

func condense(value string) string {
	value = strings.Join(strings.Fields(value), " ")
	if value == "" {
		return "Unknown"
	}
	if len(value) > 300 {
		value = value[:300] + "…"
	}
	return value
}

func joinIDs(ids []string) string {
	if len(ids) == 0 {
		return "None"
	}
	return strings.Join(ids, ", ")
}

// ParseMappings reads the suite mapping table and skeleton bodies back for
// coverage evaluation and human review.
func ParseMappings(content []byte) []Mapping {
	var mappings []Mapping
	bodies := map[string][]string{}
	var current string
	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "### ") {
			current = strings.TrimSpace(strings.TrimPrefix(trimmed, "### "))
			continue
		}
		if current != "" {
			bodies[current] = append(bodies[current], line)
		}
	}
	inTable := false
	header := true
	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "| Record Type |") {
			inTable = true
			header = true
			continue
		}
		if inTable && !strings.HasPrefix(trimmed, "|") {
			inTable = false
			continue
		}
		if !inTable {
			continue
		}
		if header {
			header = false
			continue
		}
		cells := strings.Split(strings.Trim(strings.TrimSpace(trimmed), "|"), "|")
		if len(cells) < 3 {
			continue
		}
		anchor := strings.TrimPrefix(strings.TrimSpace(cells[2]), "#")
		mappings = append(mappings, Mapping{
			RecordType: strings.TrimSpace(cells[0]),
			RecordID:   strings.TrimSpace(cells[1]),
			Skeleton:   "#" + anchor,
			Body:       strings.TrimSpace(strings.Join(bodies[anchor], "\n")),
		})
	}
	return mappings
}
