package assurance

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// SECFinding is one committed synthesis security finding in bounded form.
type SECFinding struct {
	ID      string
	Payload string
	Claims  []string
}

// PrivacyPlan builds one bounded privacy unit that authors exactly one
// PRIVACY-FINDING (0C) record per committed SECURITY-FINDING, each declaring the
// exact SEC mirror that reciprocates it.
func PrivacyPlan(db string, doc *cdl.EIRDoc, namespace, snapshot string) ([]work.Unit, error) {
	names, _, _, err := work.FieldContract(doc, "SECURITY-FINDING")
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, "SECURITY-FINDING")
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}
	sort.Strings(ids)
	claimsByFinding, err := findingClaims(db, doc)
	if err != nil {
		return nil, err
	}
	claimNames, _, _, err := work.FieldContract(doc, "CLAIM")
	if err != nil {
		return nil, err
	}
	claimText := map[string]string{}
	for _, id := range ids {
		for _, claimID := range claimsByFinding[id] {
			if _, ok := claimText[claimID]; ok {
				continue
			}
			payload, err := store.Get(db, doc, "CLAIM", claimID)
			if err != nil {
				return nil, err
			}
			values, err := work.PayloadText(claimNames, payload)
			if err != nil {
				return nil, err
			}
			claimText[claimID] = strings.Join(strings.Fields(values["claim"]), " ")
		}
	}
	var b strings.Builder
	b.WriteString("# Privacy finding context (runtime-generated, bounded)\n\n")
	b.WriteString("Author exactly one PRIVACY-FINDING record per security finding below, in the same order.\n")
	b.WriteString("Each record is the 0C security/privacy review counterpart of the synthesis SEC finding; the SEC finding is its mirror.\n\n")
	b.WriteString("## Runtime directives (must be used exactly)\n\n")
	fmt.Fprintf(&b, "- **System Namespace:** %s\n", namespace)
	fmt.Fprintf(&b, "- **Environment/Snapshot:** %s\n", snapshot)
	b.WriteString("- `synthesis-sec-mirror` is the exact SEC ID shown for the finding. Never invent an ID.\n")
	b.WriteString("- `data-threat-review-category` uses the SEC category and one of: auth, authorization, transport, at-rest, in-transit, PII, secrets, audit, threat, privacy.\n")
	b.WriteString("- `affected-src-er-if-persona-environment` lists the affected record IDs and environments, grounded in the finding and its claims.\n")
	b.WriteString("- `observed-handling-risk-unknowns` states what was observed, the risk, and any unknowns. Never declare legal compliance.\n")
	b.WriteString("- `human-review-required` must state that human assessment is required.\n")
	b.WriteString("- `material-claims-version` lists the claim IDs shown for the finding, comma-separated. No other tokens.\n\n")
	reads := make([]work.Ref, 0, len(ids))
	for _, id := range ids {
		payload, err := store.Get(db, doc, "SECURITY-FINDING", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		reads = append(reads, work.Ref{FieldID: "SECURITY-FINDING", RecordID: id})
		fmt.Fprintf(&b, "## Security finding %s\n\n", id)
		fmt.Fprintf(&b, "- **Category:** %s\n- **Risk Level:** %s\n- **Affected IDs / Environments:** %s\n- **Observed Safeguard / Gap:** %s\n- **Human Review / Target Mitigation:** %s\n",
			values["category"], values["risk-level"], values["affected-ids-environments"], values["observed-safeguard-gap"], values["human-review-target-mitigation"])
		claims := claimsByFinding[id]
		if len(claims) == 0 {
			b.WriteString("- **Claims:** None\n\n")
			continue
		}
		b.WriteString("- **Claims:**\n")
		for _, claimID := range claims {
			fmt.Fprintf(&b, "  - %s: %s\n", claimID, claimText[claimID])
		}
		b.WriteString("\n")
		for _, claimID := range claims {
			reads = append(reads, work.Ref{FieldID: "CLAIM", RecordID: claimID})
		}
	}
	return []work.Unit{{
		ID:                   "U-ASSURANCE-PRIVACY",
		Mode:                 "Discovery",
		Scope:                namespace,
		IDType:               "PRV",
		Namespace:            namespace,
		IDFromBindings:       "material-claims-version",
		IDDiscriminatorField: "synthesis-sec-mirror",
		ExpectedRecords:      len(ids),
		RequiresReview:       true,
		ResultType:           "PRIVACY-FINDING",
		Reads:                reads,
		Inputs:               []work.Input{{Name: "privacy-context.md", Bytes: []byte(b.String())}},
		Sections:             []string{"workspace-registries", "nfr-security-catalog"},
		ReasoningEffort:      "off",
	}}, nil
}

// PrivacyStats reports 0C/SEC reciprocal mirror coverage.
type PrivacyStats struct {
	Findings      int
	Mirrored      int
	SECFindings   int
	SECWithOrigin int
}

// PrivacyCoverage counts committed PRV findings, how many declare a resolving
// SEC mirror, and how many SEC findings have a 0C origin.
func PrivacyCoverage(db string, doc *cdl.EIRDoc) (PrivacyStats, error) {
	var stats PrivacyStats
	names, _, _, err := work.FieldContract(doc, "PRIVACY-FINDING")
	if err != nil {
		return stats, err
	}
	ids, err := store.List(db, doc, "PRIVACY-FINDING")
	if err != nil {
		return stats, err
	}
	stats.Findings = len(ids)
	mirrors := map[string]int{}
	for _, id := range ids {
		payload, err := store.Get(db, doc, "PRIVACY-FINDING", id)
		if err != nil {
			return stats, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return stats, err
		}
		mirror := firstTokenField(values["synthesis-sec-mirror"])
		if mirror == "" || strings.EqualFold(mirror, "none") {
			continue
		}
		exists, err := store.Exists(db, doc, "SECURITY-FINDING", mirror)
		if err != nil {
			return stats, err
		}
		if !exists {
			continue
		}
		stats.Mirrored++
		mirrors[mirror]++
	}
	secIDs, err := store.List(db, doc, "SECURITY-FINDING")
	if err != nil {
		return stats, err
	}
	stats.SECFindings = len(secIDs)
	for _, id := range secIDs {
		if mirrors[id] > 0 {
			stats.SECWithOrigin++
		}
	}
	return stats, nil
}

// findingClaims maps each SEC finding to its committed claim IDs from the
// derived-from edges, falling back to its declared source bindings.
func findingClaims(db string, doc *cdl.EIRDoc) (map[string][]string, error) {
	names, _, _, err := work.FieldContract(doc, "SECURITY-FINDING")
	if err != nil {
		return nil, err
	}
	edgeNames, _, _, err := work.FieldContract(doc, "DERIVED-FROM-EDGE")
	if err != nil {
		return nil, err
	}
	edgeIDs, err := store.List(db, doc, "DERIVED-FROM-EDGE")
	if err != nil {
		return nil, err
	}
	claims := map[string][]string{}
	for _, id := range edgeIDs {
		payload, err := store.Get(db, doc, "DERIVED-FROM-EDGE", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(edgeNames, payload)
		if err != nil {
			return nil, err
		}
		if !strings.EqualFold(strings.TrimSpace(values["binding-kind"]), "CLAIM") {
			continue
		}
		claim := firstTokenField(values["claim-id"])
		derived := firstTokenField(values["derived-record-id"])
		if claim == "" || strings.EqualFold(claim, "none") || derived == "" {
			continue
		}
		claims[derived] = append(claims[derived], claim)
	}
	secIDs, err := store.List(db, doc, "SECURITY-FINDING")
	if err != nil {
		return nil, err
	}
	for _, id := range secIDs {
		if len(claims[id]) > 0 {
			continue
		}
		payload, err := store.Get(db, doc, "SECURITY-FINDING", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		for _, token := range strings.Split(values["source-bindings"], ",") {
			token = firstTokenField(strings.TrimSpace(token))
			if strings.HasPrefix(token, "CLM-") {
				claims[id] = append(claims[id], token)
			}
		}
	}
	for _, ids := range claims {
		sort.Strings(ids)
	}
	return claims, nil
}
