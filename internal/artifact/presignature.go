package artifact

import (
	"fmt"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	md "github.com/nckirik/legacy-autopsy/internal/markdown"
)

// PreSignatureInputs recomputes the §15.4 pre-signature authoritative input-set
// fingerprint from the outer manifest payload and the two prior payloads. The
// preimage is `EXIT-E-PRE-SIGNATURE-INPUTS|`, the canonical outer snapshot, the
// two registry identities and fingerprints, the readiness and certificate
// artifact-payload fingerprints, and the canonical Package Members table.
func PreSignatureInputs(doc *cdl.EIRDoc, outerPayload []byte, readiness, certificate ChainStage) (string, error) {
	snapshot := payloadField(outerPayload, "system-protocol-snapshot")
	if snapshot == "" {
		return "", fmt.Errorf("artifact: outer manifest has no system-protocol-snapshot")
	}
	readinessFingerprint, err := PayloadFingerprint(readiness.Type, readiness.Path, readiness.Payload)
	if err != nil {
		return "", err
	}
	certificateFingerprint, err := PayloadFingerprint(certificate.Type, certificate.Path, certificate.Payload)
	if err != nil {
		return "", err
	}
	table := canonicalMemberTable(outerPayload)
	if table == "" {
		return "", fmt.Errorf("artifact: outer manifest has no Package Members table")
	}
	var preimage strings.Builder
	preimage.WriteString("EXIT-E-PRE-SIGNATURE-INPUTS|")
	preimage.WriteString(snapshot)
	preimage.WriteString("\n")
	preimage.WriteString("EXIT-E-CHECKS-v1 " + CheckRegistryFingerprint(doc))
	preimage.WriteString("\n")
	preimage.WriteString("EXIT-E-FINAL-CHECKS-v1 " + RegistryFingerprint("EXIT-E-FINAL-CHECKS-v1", RegistryCheckIDs(doc, "FINAL-CHECK")))
	preimage.WriteString("\n")
	preimage.WriteString(readinessFingerprint)
	preimage.WriteString("\n")
	preimage.WriteString(certificateFingerprint)
	preimage.WriteString("\n")
	preimage.WriteString(table)
	return capabilities.Hash([]byte(preimage.String())), nil
}

func verifyPreSignature(doc *cdl.EIRDoc, outer ChainStage, stages []ChainStage) error {
	declared := payloadField(outer.Payload, "pre-signature-validation-input-set-fingerprint")
	if declared == "" {
		return fmt.Errorf("artifact: outer manifest has no pre-signature input-set fingerprint")
	}
	readiness := stages[2]
	certificate := stages[3]
	expected, err := PreSignatureInputs(doc, outer.Payload, readiness, certificate)
	if err != nil {
		return err
	}
	if declared != expected {
		return fmt.Errorf("artifact: pre-signature input-set fingerprint does not recompute")
	}
	return nil
}

// canonicalMemberTable extracts the outer Package Members table and returns its
// canonical §4.1.2 serialization with the declared header and separator.
func canonicalMemberTable(payload []byte) string {
	parsed, err := md.Parse(string(payload))
	if err != nil {
		return ""
	}
	inTable := false
	var rows []string
	for _, node := range parsed.Nodes {
		if node.Kind == md.Field {
			inTable = capabilities.ColumnName(node.Label) == "package-members"
			continue
		}
		if !inTable || node.Kind != md.TableRow {
			continue
		}
		rows = append(rows, strings.TrimSpace(node.Raw))
	}
	if len(rows) < 2 {
		return ""
	}
	canonical, err := capabilities.CanonicalMarkdown([]byte(strings.Join(rows, "\n") + "\n"))
	if err != nil {
		return ""
	}
	return strings.TrimSuffix(string(canonical), "\n")
}
