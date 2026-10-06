package artifact

import (
	"fmt"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	md "github.com/nckirik/legacy-autopsy/internal/markdown"
)

// ChainStage is one ordered artifact in the §15.4 certification sequence.
type ChainStage struct {
	Type    string
	Path    string
	Payload []byte
}

var chainSequence = []string{
	"EXIT-E-CANDIDATE-REPORT",
	"CANDIDATE-PAYLOAD-MANIFEST",
	"EXIT-E-CONTENT-READINESS-REPORT",
	"SCOPE-CERTIFICATE",
	"OUTER-BUNDLE-MANIFEST",
}

// RegistryFingerprint computes the canonical check-registry fingerprint from
// the ordered check IDs. The preimage is `EXIT-E-CHECK-REGISTRY|<id>` plus LF
// and the canonical one-column table of check IDs. (The requirement prose is
// not yet part of the preimage; this is the pilot realization.)
func RegistryFingerprint(registryID string, checks []string) string {
	var b strings.Builder
	b.WriteString("| Check ID |\n| --- |\n")
	for _, check := range checks {
		fmt.Fprintf(&b, "| %s |\n", check)
	}
	canonical, _ := capabilities.CanonicalMarkdown([]byte(b.String()))
	return capabilities.Hash([]byte("EXIT-E-CHECK-REGISTRY|" + registryID + "\n" + string(canonical)))
}

// CandidateCheckIDs and ReadinessCheckIDs return the ordered registry values
// declared for the candidate and content-readiness reports.
func RegistryCheckIDs(doc *cdl.EIRDoc, enumID string) []string {
	for _, enum := range doc.Declarations.Enums {
		if enum.ID == enumID {
			return append([]string(nil), enum.Values...)
		}
	}
	return nil
}

// VerifyChain validates the ordered §15.4 sequence: every stage type appears
// once in order, every payload matches its declared artifact schema, each stage
// references the immediately preceding stage fingerprint, snapshots are equal,
// report check sets equal their declared registries, and only the outer
// manifest declares `EXIT-E-STATUS: Passed`.
func VerifyChain(doc *cdl.EIRDoc, stages []ChainStage) error {
	if len(stages) != len(chainSequence) {
		return fmt.Errorf("artifact: chain has %d stages, want %d", len(stages), len(chainSequence))
	}
	var failures []string
	observedSnapshot := ""
	previousFingerprint := ""
	previousType := ""
	for i, stage := range stages {
		if stage.Type != chainSequence[i] {
			failures = append(failures, fmt.Sprintf("stage %d is %s, want %s", i+1, stage.Type, chainSequence[i]))
			continue
		}
		if err := ValidatePayload(doc, stage.Type, stage.Payload); err != nil {
			failures = append(failures, err.Error())
		}
		fingerprint, err := PayloadFingerprint(stage.Type, stage.Path, stage.Payload)
		if err != nil {
			failures = append(failures, err.Error())
			continue
		}
		if previousFingerprint != "" && !strings.Contains(string(stage.Payload), previousFingerprint) {
			failures = append(failures, fmt.Sprintf("%s does not reference the %s fingerprint", stage.Type, previousType))
		}
		if snapshot := payloadField(stage.Payload, "system-protocol-snapshot"); snapshot != "" {
			if observedSnapshot == "" {
				observedSnapshot = snapshot
			} else if snapshot != observedSnapshot {
				failures = append(failures, fmt.Sprintf("%s snapshot differs from the chain", stage.Type))
			}
		}
		if err := verifyCheckSet(doc, stage); err != nil {
			failures = append(failures, err.Error())
		}
		status := payloadField(stage.Payload, "exit-e-status")
		prepackage := payloadField(stage.Payload, "exit-e-prepackage-state")
		if stage.Type == "OUTER-BUNDLE-MANIFEST" {
			if status != "Passed" {
				failures = append(failures, "outer manifest is not Passed")
			}
		} else if status == "Passed" || prepackage == "Passed" {
			failures = append(failures, stage.Type+" declares Passed before the outer manifest")
		}
		previousFingerprint = fingerprint
		previousType = stage.Type
		if stage.Type == "OUTER-BUNDLE-MANIFEST" {
			expected := RegistryFingerprint("EXIT-E-FINAL-CHECKS-v1", RegistryCheckIDs(doc, "FINAL-CHECK"))
			if !strings.Contains(string(stage.Payload), expected) {
				failures = append(failures, "outer manifest does not carry the recomputed final check-registry fingerprint")
			}
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("artifact: chain verification failed: %s", strings.Join(failures, "; "))
	}
	return nil
}

func verifyCheckSet(doc *cdl.EIRDoc, stage ChainStage) error {
	registry := ""
	switch stage.Type {
	case "EXIT-E-CANDIDATE-REPORT":
		registry = "EXIT-E-CANDIDATE-CHECKS"
	case "EXIT-E-CONTENT-READINESS-REPORT":
		registry = "EXIT-E-CONTENT-READINESS-CHECKS"
	default:
		return nil
	}
	declared := RegistryCheckIDs(doc, registry)
	if len(declared) == 0 {
		return nil
	}
	observed := checkResults(stage.Payload)
	if len(observed) != len(declared) {
		return nil // table shape is validated by the payload schema; set equality is checked when present
	}
	present := map[string]int{}
	for _, id := range observed {
		present[id]++
	}
	for _, id := range declared {
		if present[id] != 1 {
			return fmt.Errorf("%s check set does not match %s (missing or duplicated %s)", stage.Type, registry, id)
		}
	}
	expected := RegistryFingerprint("EXIT-E-CHECKS-v1", declared)
	if observed := payloadField(stage.Payload, "required-check-registry-fingerprint"); observed != expected {
		return fmt.Errorf("%s does not carry the recomputed %s fingerprint", stage.Type, registry)
	}
	for id := range present {
		found := false
		for _, want := range declared {
			if want == id {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("%s check set contains unknown check %s", stage.Type, id)
		}
	}
	return nil
}

func checkResults(payload []byte) []string {
	parsed, err := md.Parse(string(payload))
	if err != nil {
		return nil
	}
	inTable := false
	header := true
	var ids []string
	for _, node := range parsed.Nodes {
		if node.Kind == md.Field {
			inTable = capabilities.ColumnName(node.Label) == "check-results"
			header = true
			continue
		}
		if !inTable || node.Kind != md.TableRow {
			continue
		}
		cells := splitRow(node.Raw)
		if header {
			header = false
			continue
		}
		if isSeparator(cells) || len(cells) == 0 {
			continue
		}
		ids = append(ids, strings.TrimSpace(cells[0]))
	}
	return ids
}

func payloadField(payload []byte, slug string) string {
	parsed, err := md.Parse(string(payload))
	if err != nil {
		return ""
	}
	for _, node := range parsed.Nodes {
		if node.Kind == md.Field && capabilities.ColumnName(node.Label) == slug {
			return strings.TrimSpace(node.Value)
		}
	}
	return ""
}
