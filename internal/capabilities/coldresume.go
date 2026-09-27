package capabilities

import (
	"fmt"
	"strings"
)

// ColdResumeField is one cold-resume assessment or validation field.
type ColdResumeField struct {
	Label string
	Value string
}

// ColdResumeCarrier is the field excluded from its own preimage.
const ColdResumeCarrier = "Cold-Resume Check Fingerprint"

var coldResumeRequiredFields = []string{
	"Iteration / Invocation / Mode",
	"Persona / Cluster / Track / POV",
	"Relevant Active Tickets and Frontiers",
	"Last Mutation Affecting Targets",
	"Semantic Blockers and Contradictions",
	"Proposed Allowed Next Actions",
	"Explicit Context Exclusions and Reasons",
	"Runtime-Validated Loaded Files / Fingerprints",
	"Runtime Scope / Mode / Write-Right Check",
	"Runtime Stale Check",
	"Cold-Resume Result",
}

// ColdResumeFingerprint computes the §8.5 cold-resume check fingerprint. Preimage:
// "COLD-RESUME|" + namespace + "|" + invocation ID + LF + the canonical
// serialization of every field except the carrier.
func ColdResumeFingerprint(namespace, invocationID string, fields []ColdResumeField) (string, error) {
	if namespace == "" || invocationID == "" {
		return "", fmt.Errorf("cold-resume: namespace and invocation ID are required")
	}
	known := map[string]bool{ColdResumeCarrier: true}
	for _, label := range coldResumeRequiredFields {
		known[label] = true
	}
	seen := map[string]bool{}
	for _, field := range fields {
		if !known[field.Label] {
			return "", fmt.Errorf("cold-resume: unknown field %q", field.Label)
		}
		if seen[field.Label] {
			return "", fmt.Errorf("cold-resume: duplicate field %q", field.Label)
		}
		seen[field.Label] = true
	}
	for _, label := range coldResumeRequiredFields {
		if !seen[label] {
			return "", fmt.Errorf("cold-resume: missing field %q", label)
		}
	}
	var b strings.Builder
	for _, field := range fields {
		if field.Label == ColdResumeCarrier {
			continue
		}
		b.WriteString("- **" + field.Label + ":** " + field.Value + "\n")
	}
	canonical, err := CanonicalMarkdown([]byte(b.String()))
	if err != nil {
		return "", err
	}
	preimage := "COLD-RESUME|" + namespace + "|" + invocationID + "\n" + string(canonical)
	return Hash([]byte(preimage)), nil
}
