package parity

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

func compileWithRegistry(t *testing.T, rel string) *cdl.Result {
	t.Helper()
	root := root(t)
	read := func(rel string) []byte {
		b, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	ledger, err := cdl.LoadLedger(filepath.Join(root, "protocol/identity-ledger.json"))
	if err != nil {
		t.Fatal(err)
	}
	res, err := cdl.Compile(cdl.CompileInput{
		Sources: []cdl.Source{
			{Path: "protocol/globals.cdl", Bytes: read("protocol/globals.cdl")},
			{Path: "protocol/status-taxonomy.cdl", Bytes: read("protocol/status-taxonomy.cdl")},
			{Path: rel, Bytes: read(rel)},
		},
		Ledger: ledger,
	})
	if err != nil {
		t.Fatal(err)
	}
	return res
}

// TestEvidenceAndDecisionsDrift checks the §5.2, §5.4, §5.5, and §5.7 record field
// lists against protocol.md and requires §5.3/§5.6 normative token coverage.
func TestEvidenceAndDecisionsDrift(t *testing.T) {
	res := compileWithRegistry(t, "protocol/evidence-and-decisions.cdl")
	fields := map[string][]string{}
	for _, f := range res.EIR.Declarations.Fields {
		names := make([]string, 0, len(f.Fields))
		for _, spec := range f.Fields {
			names = append(names, spec.Name)
		}
		fields[f.ID] = names
	}
	cdlText := ""
	for _, sec := range res.EIR.Sections {
		cdlText += sec.Title + "\n" + sec.Goal + "\n"
	}
	for _, r := range res.EIR.Declarations.Rules {
		cdlText += r.Goal + "\n"
	}

	model, err := protocol.Load(protocolPathIn(root(t)))
	if err != nil {
		t.Fatal(err)
	}

	// §5.2 splits into the claim record and the block evidence summary.
	text52, err := model.SectionText("5.2. Claim-level evidence - `12-CLAIM-EVIDENCE.md`")
	if err != nil {
		t.Fatal(err)
	}
	var claimLabels, blockLabels []string
	inBlockSummary := false
	for _, line := range strings.Split(text52, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "### Block evidence summary") {
			inBlockSummary = true
			continue
		}
		match := logFieldLine.FindStringSubmatch(trimmed)
		if match == nil {
			continue
		}
		if inBlockSummary {
			blockLabels = append(blockLabels, normalizeLabel(match[1]))
		} else {
			claimLabels = append(claimLabels, normalizeLabel(match[1]))
		}
	}
	compareSets(t, "CLAIM", fields["CLAIM"], claimLabels)
	compareSets(t, "BLOCK-EVIDENCE-SUMMARY", fields["BLOCK-EVIDENCE-SUMMARY"], blockLabels)

	for ref, fieldID := range map[string]string{
		"5.4. Contradictions - `14-CONTRADICTIONS.md`":     "CONTRADICTION",
		"5.5. Authoritative decisions - `15-DECISIONS.md`": "DECISION",
		"5.7. Human confirmations - `18-CONFIRMATIONS.md`": "CONFIRMATION",
	} {
		text, err := model.SectionText(ref)
		if err != nil {
			t.Fatal(err)
		}
		var labels []string
		for _, line := range strings.Split(text, "\n") {
			if match := logFieldLine.FindStringSubmatch(strings.TrimSpace(line)); match != nil {
				labels = append(labels, normalizeLabel(match[1]))
			}
		}
		compareSets(t, fieldID, fields[fieldID], labels)
	}

	for _, ref := range []string{
		"5.3. Sanitized projections and helper limitations",
		"5.6. Acquisition-candidate reconciliation - `17-ACQUISITION-CANDIDATES.md`",
	} {
		text, err := model.SectionText(ref)
		if err != nil {
			t.Fatal(err)
		}
		for _, token := range hyphenatedToken.FindAllString(text, -1) {
			if !strings.Contains(cdlText, token) {
				t.Fatalf("%s token %s is missing from the CDL section text", ref, token)
			}
		}
	}
}
