package handbook

import (
	"fmt"
	"sort"
	"strings"
)

// mechanicalFields are never shown in the human view.
var mechanicalFields = map[string]bool{
	"source-bindings": true,
	"block-id":        true,
}

// summaryFields maps a catalog to the ordered prose fields shown as the record
// summary. Other fields go to a collapsed audit block.
var summaryFields = map[string][]string{
	"ARCHITECTURE-MODULE":  {"responsibility-domain-boundary", "public-interfaces", "dependencies-dependents", "observed-legacy-placement"},
	"ENTITY-RECORD":        {"domain-context-purpose", "storage", "relationships", "invariants"},
	"RELATIONSHIP-RECORD":  {"cardinality", "enforcement", "business-constraint", "reciprocal-entity-references"},
	"STATE-MACHINE":        {"scope-owner-entity-or-aggregate", "states-meaning", "initial-terminal-states", "invalid-refused-transitions"},
	"DATABASE-ROUTINE":     {"atomic-body-intent", "inputs-outputs", "entities-read-written", "error-handling-semantics"},
	"BUSINESS-RULE-RECORD": {"applies-when", "condition", "action-effect", "precedence-conflicts"},
	"USE-CASE":             {"goal", "trigger-entry", "observed-legacy-behavior", "behavioral-contract"},
	"INTERFACE-RECORD":     {"channel-transport-direction", "protocol-content-type", "request-schema", "response-schema", "error-contract"},
	"DEPLOYMENT-ELEMENT":   {"component-placement-replication", "connections", "single-point-of-failure", "environment-snapshot-matrix"},
	"CONFIGURATION-RECORD": {"exact-key", "type-default-allowed-values", "behavioral-effect", "enable-disable-semantics"},
	"SCHEDULE-RECORD":      {"type-trigger-frequency-window", "owner-persona", "actions", "failure-retry-dead-letter-idempotency"},
	"NFR-RECORD":           {"attribute", "observed-inferred-value", "scope-environment", "driver-constraint"},
	"SECURITY-FINDING":     {"category", "finding-data-classification", "risk-level", "affected-ids-environments", "observed-safeguard-gap"},
	"FAULT-RECORD":         {"failure-class-trigger-detection-point", "propagation-path", "rollback-compensating-behavior", "recovery-procedure"},
}

// recordEvidence extracts the claim IDs a record was derived from.
func recordEvidence(values map[string]string) []string {
	var claims []string
	for _, token := range strings.Split(values["source-bindings"], ",") {
		token = strings.TrimSpace(token)
		if strings.HasPrefix(token, "CLM-") {
			claims = append(claims, token)
		}
	}
	sort.Strings(claims)
	return claims
}

// renderRecord renders one committed record as a readable handbook block: a
// title, summary prose, an evidence line, and a collapsed audit table with the
// remaining fields. Mechanical identity fields never appear in the human view.
func renderRecord(fieldID, id string, values map[string]string) string {
	title := strings.TrimSpace(values["record-title"])
	if title == "" {
		title = id
	}
	var b strings.Builder
	if title == id {
		fmt.Fprintf(&b, "#### `%s`\n\n", id)
	} else {
		fmt.Fprintf(&b, "#### %s (`%s`)\n\n", title, id)
	}
	summary := summaryFields[fieldID]
	if len(summary) == 0 {
		for _, name := range orderedFieldNames(values) {
			if mechanicalFields[name] || name == "record-title" || name == "domain-path" {
				continue
			}
			summary = append(summary, name)
			if len(summary) == 3 {
				break
			}
		}
	}
	shown := map[string]bool{"record-title": true}
	for _, name := range summary {
		value := strings.TrimSpace(values[name])
		if value == "" || strings.EqualFold(value, "none") {
			continue
		}
		shown[name] = true
		fmt.Fprintf(&b, "- **%s:** %s\n", humanLabel(name), value)
	}
	b.WriteString("\n")
	if claims := recordEvidence(values); len(claims) > 0 {
		fmt.Fprintf(&b, "- **Evidence claims:** %s\n", strings.Join(claims, ", "))
	}
	if path := strings.TrimSpace(values["domain-path"]); path != "" {
		fmt.Fprintf(&b, "- **Location:** `%s`\n", path)
	}
	fmt.Fprintf(&b, "\n<details><summary>All record fields</summary>\n\n")
	for _, name := range orderedFieldNames(values) {
		if mechanicalFields[name] || shown[name] {
			continue
		}
		value := strings.TrimSpace(values[name])
		if value == "" {
			value = "None"
		}
		fmt.Fprintf(&b, "- **%s:** %s\n", humanLabel(name), value)
	}
	fmt.Fprintf(&b, "\n</details>\n")
	return b.String()
}

func orderedFieldNames(values map[string]string) []string {
	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// humanLabel converts a kebab-case field name into a readable label.
func humanLabel(name string) string {
	words := strings.Split(name, "-")
	for i, word := range words {
		if word == "" {
			continue
		}
		words[i] = strings.ToUpper(word[:1]) + word[1:]
	}
	return strings.Join(words, " ")
}
