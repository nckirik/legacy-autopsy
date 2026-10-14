// Package handoff emits the immutable §15.4 step-2 candidate artifacts.
package handoff

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/artifact"
	"github.com/nckirik/legacy-autopsy/internal/assurance"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/synthesis"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// CandidateReportArtifactType is the generic artifact type of the report.
const CandidateReportArtifactType = "EXIT-E-CANDIDATE-REPORT"

// CandidateReportDefaultPath is the default package-relative artifact path.
const CandidateReportDefaultPath = "exit-e/EXIT-E-CANDIDATE-REPORT.md"

// ReportOptions are the runtime-owned candidate report directives.
type ReportOptions struct {
	Namespace             string
	ProtocolVersion       string
	EnvironmentsSnapshots string
	StorageBinding        string
	HandbookDir           string
	ArtifactPath          string
}

// CandidateReport is one emitted candidate report artifact.
type CandidateReport struct {
	Payload         []byte
	Fingerprint     string
	InstanceBinding string
	Result          string
}

// SystemSnapshotValue is the canonical System / Protocol / Snapshot identity
// carried by every snapshot-bearing packaging artifact.
func SystemSnapshotValue(namespace, protocolVersion, environments string) string {
	return strings.Join([]string{namespace, protocolVersion, environments}, " / ")
}

// BuildCandidateReport validates the committed state, computes every
// deterministic summary and evidence-set fingerprint, and serializes the
// immutable Exit E Candidate Report payload.
func BuildCandidateReport(db string, doc *cdl.EIRDoc, opts ReportOptions) (CandidateReport, error) {
	var report CandidateReport
	if opts.ArtifactPath == "" {
		opts.ArtifactPath = CandidateReportDefaultPath
	}
	snapshot := SystemSnapshotValue(opts.Namespace, opts.ProtocolVersion, opts.EnvironmentsSnapshots)
	if strings.TrimSpace(snapshot) == "" {
		return report, fmt.Errorf("handoff: candidate report requires a pinned system/protocol/snapshot identity")
	}
	checks, err := EvaluateCandidateChecks(db, doc, Options{HandbookDir: opts.HandbookDir})
	if err != nil {
		return report, err
	}
	gate, err := assurance.EvaluateExitA(db, doc)
	if err != nil {
		return report, err
	}
	_, storeFingerprint, err := store.CanonicalDump(db, doc, store.Provenance{
		Generator:         "legacy-autopsy candidate report",
		SourceFingerprint: capabilities.Hash([]byte(snapshot)),
	})
	if err != nil {
		return report, err
	}
	handbookFiles, handbookFingerprint, err := handbookFileFingerprints(opts.HandbookDir)
	if err != nil {
		return report, err
	}

	registryFingerprint := artifact.CheckRegistryFingerprint(doc)
	exitAFingerprint := exitAReportFingerprint(snapshot, gate)
	inputFingerprint, err := candidateInputFingerprint(snapshot, storeFingerprint, handbookFiles, registryFingerprint)
	if err != nil {
		return report, err
	}
	readiness := artifact.ArtifactSchemaVersion
	version := strings.TrimSpace(opts.ProtocolVersion) + " / " + readiness

	coverageText, _ := coverageSummary(db, doc, snapshot, gate)
	traceabilityText, _, err := traceabilitySummary(db, doc, snapshot, storeFingerprint)
	if err != nil {
		return report, err
	}
	handbookText, _ := handbookSummary(opts.HandbookDir, checks)

	blocked := false
	for _, check := range checks {
		if check.Result != "Pass" {
			blocked = true
		}
	}
	result := "Ready-For-Human-Review"
	if blocked {
		result = "Blocked"
	}

	var b strings.Builder
	b.WriteString("<!-- ARTIFACT-PAYLOAD: " + CandidateReportArtifactType + " -->\n\n")
	b.WriteString("# Exit E Candidate Report\n\n")
	fmt.Fprintf(&b, "- **ARTIFACT-TYPE:** %s\n", CandidateReportArtifactType)
	fmt.Fprintf(&b, "- **ARTIFACT-PATH:** %s\n", opts.ArtifactPath)
	fmt.Fprintf(&b, "- **PROTOCOL / ARTIFACT-SCHEMA VERSION:** %s\n", version)
	fmt.Fprintf(&b, "- **System / Protocol / Snapshot:** %s\n", snapshot)
	fmt.Fprintf(&b, "- **Required Check Registry / Fingerprint:** EXIT-E-CHECKS-v1 %s\n", registryFingerprint)
	fmt.Fprintf(&b, "- **Exit A Report Fingerprint:** %s\n", exitAFingerprint)
	fmt.Fprintf(&b, "- **Deterministic Candidate Input-Set Fingerprint:** %s\n", inputFingerprint)
	fmt.Fprintf(&b, "- **Coverage Validation Summary:** %s\n", coverageText)
	fmt.Fprintf(&b, "- **Traceability Validation Summary:** %s\n", traceabilityText)
	fmt.Fprintf(&b, "- **Handbook Validation Summary:** %s\n", handbookText)
	b.WriteString("- **Semantic Record Counts:**\n")
	b.WriteString("  | Record Type | Count |\n  | --- | --- |\n")
	for _, row := range semanticRecordCounts(db, doc) {
		fmt.Fprintf(&b, "  | %s | %s |\n", row[0], row[1])
	}
	b.WriteString("- **Decision Content Counts:**\n")
	b.WriteString("  | Decision Type | Count |\n  | --- | --- |\n")
	for _, row := range decisionContentCounts(db, doc) {
		fmt.Fprintf(&b, "  | %s | %s |\n", row[0], row[1])
	}
	requirements := artifact.CheckRequirements(doc)
	b.WriteString("- **Check Results:**\n")
	b.WriteString("  | Check ID | Requirement | Result | Evidence Fingerprint |\n  | --- | --- | --- | --- |\n")
	type rowResult struct {
		check       CandidateCheck
		result      string
		evidence    string
		requirement string
	}
	rows := make([]rowResult, 0, len(checks))
	for _, check := range checks {
		result := check.Result
		if result != "Pass" {
			result = "Fail"
		}
		requirement := requirements[check.ID]
		if requirement == "" {
			requirement = check.ID
		}
		evidence, err := checkEvidenceFingerprint(doc, snapshot, check.ID, storeFingerprint, handbookFingerprint, inputFingerprint, registryFingerprint)
		if err != nil {
			return report, err
		}
		rows = append(rows, rowResult{check: check, result: result, evidence: evidence, requirement: requirement})
	}
	for _, row := range rows {
		fmt.Fprintf(&b, "  | %s | %s | %s | %s |\n", row.check.ID, row.requirement, row.result, row.evidence)
	}
	b.WriteString("- **Blockers:**\n")
	b.WriteString("  | Blocker ID | Affected IDs | Description | Evidence Fingerprint |\n  | --- | --- | --- | --- |\n")
	for _, row := range rows {
		if row.result == "Pass" {
			continue
		}
		fmt.Fprintf(&b, "  | BLOCKER-%s | %s | %s | %s |\n", row.check.ID, row.check.ID, sanitizeCell(row.check.Detail), row.evidence)
	}
	fmt.Fprintf(&b, "- **Candidate Result:** %s\n", result)
	b.WriteString("- **ARTIFACT-PAYLOAD-FINGERPRINT:** PENDING\n\n")
	b.WriteString("<!-- END-ARTIFACT-PAYLOAD: " + CandidateReportArtifactType + " -->\n\n")
	b.WriteString("<!-- ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: " + CandidateReportArtifactType + " -->\n\n")
	b.WriteString("- **Signature / Signatories / Timestamps:** None\n")
	b.WriteString("- **Envelope Audit:** not attached at the candidate stage\n")
	b.WriteString("<!-- END-ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: " + CandidateReportArtifactType + " -->\n")

	payload := []byte(b.String())
	fingerprint, err := artifact.PayloadFingerprint(CandidateReportArtifactType, opts.ArtifactPath, payload)
	if err != nil {
		return report, err
	}
	withCarrier := strings.Replace(string(payload),
		"- **ARTIFACT-PAYLOAD-FINGERPRINT:** PENDING", "- **ARTIFACT-PAYLOAD-FINGERPRINT:** "+fingerprint, 1)
	final := []byte(withCarrier)
	recomputed, err := artifact.PayloadFingerprint(CandidateReportArtifactType, opts.ArtifactPath, final)
	if err != nil {
		return report, err
	}
	if recomputed != fingerprint {
		return report, fmt.Errorf("handoff: candidate report fingerprint does not recompute")
	}
	if err := artifact.ValidatePayload(doc, CandidateReportArtifactType, final); err != nil {
		return report, err
	}
	instance, err := artifact.InstanceBinding(CandidateReportArtifactType, opts.ArtifactPath, fingerprint)
	if err != nil {
		return report, err
	}
	report.Payload = final
	report.Fingerprint = fingerprint
	report.InstanceBinding = instance
	report.Result = result
	return report, nil
}

func exitAReportFingerprint(snapshot string, gate assurance.GateResult) string {
	rows := make([]string, 0, len(gate.Checks))
	for _, check := range gate.Checks {
		rows = append(rows, check.ID+"="+check.Result+": "+check.Detail)
	}
	return artifact.CanonicalSetFingerprint("EXIT-A-REPORT", snapshot, rows)
}

func candidateInputFingerprint(snapshot, storeFingerprint string, handbook []artifact.EvidenceBinding, registryFingerprint string) (string, error) {
	bindings := []artifact.EvidenceBinding{
		{Path: "protocol/assembly.json", Kind: "PROTOCOL-REGISTRY", RecordID: "EXIT-E-CHECKS-v1", Version: "1", Fingerprint: registryFingerprint},
		{Path: "store/legacy-autopsy.sql", Kind: "DENOMINATOR-QUERY", RecordID: "candidate-store-dump", Version: "dump/2", Fingerprint: storeFingerprint},
	}
	bindings = append(bindings, handbook...)
	return artifact.BindingTableFingerprint("CANDIDATE-INPUTS", snapshot, bindings)
}

func coverageSummary(db string, doc *cdl.EIRDoc, snapshot string, gate assurance.GateResult) (string, string) {
	passed, failed := 0, 0
	for _, check := range gate.Checks {
		if check.Result == "Pass" {
			passed++
		} else {
			failed++
		}
	}
	summary, err := assurance.CoverageSummary(db, doc)
	if err != nil {
		return "unavailable", ""
	}
	mapped := summary.CoverageByDisposition["Mapped"] + summary.CoverageByDisposition["Approved-Excluded"]
	total := 0
	for _, count := range summary.CoverageByDisposition {
		total += count
	}
	fingerprint := artifact.ValidationSummaryFingerprint("coverage", snapshot, [][2]string{
		{"gate-checks-failed", fmt.Sprintf("%d", failed)},
		{"gate-checks-passed", fmt.Sprintf("%d", passed)},
		{"mapped-or-excluded", fmt.Sprintf("%d", mapped)},
		{"total", fmt.Sprintf("%d", total)},
	})
	return fmt.Sprintf("%d passed, %d failed; %d/%d mapped or excluded; %s", passed, failed, mapped, total, fingerprint), fingerprint
}

func traceabilitySummary(db string, doc *cdl.EIRDoc, snapshot, storeFingerprint string) (string, string, error) {
	result := "Pass"
	if err := referentialIntegrity(db, doc); err != nil {
		result = "Fail"
	}
	fingerprint := artifact.ValidationSummaryFingerprint("traceability", snapshot, [][2]string{
		{"result", result},
		{"store-content", storeFingerprint},
	})
	return fmt.Sprintf("%s; %s", result, fingerprint), fingerprint, nil
}

func handbookSummary(dir string, checks []CandidateCheck) (string, string) {
	chapters := 0
	findings := 0
	if dir != "" {
		entries, err := os.ReadDir(dir)
		if err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasPrefix(entry.Name(), "chapter-") && strings.HasSuffix(entry.Name(), ".md") {
					chapters++
				}
			}
		}
	}
	for _, check := range checks {
		if check.ID == "CANDIDATE-04-HANDBOOK-QUALITY" && check.Result != "Pass" {
			findings = 1
		}
	}
	fingerprint := artifact.ValidationSummaryFingerprint("handbook", dir, [][2]string{
		{"chapters", fmt.Sprintf("%d", chapters)},
		{"failed", fmt.Sprintf("%d", findings)},
	})
	return fmt.Sprintf("%d chapters, %d failed checks; %s", chapters, findings, fingerprint), fingerprint
}

func checkEvidenceFingerprint(doc *cdl.EIRDoc, snapshot, checkID, storeFingerprint, handbookFingerprint, inputFingerprint, registryFingerprint string) (string, error) {
	denominator := capabilities.Hash([]byte("DENOMINATOR-QUERY|" + checkID + "|" + snapshot + "|" + storeFingerprint + "|" + handbookFingerprint + "|" + inputFingerprint))
	return artifact.EvidenceSetFingerprint(CandidateReportArtifactType, "CHECK", checkID, snapshot, []artifact.EvidenceBinding{
		{Path: "protocol/assembly.json", Kind: "PROTOCOL-REGISTRY", RecordID: "EXIT-E-CHECKS-v1", Version: "1", Fingerprint: registryFingerprint},
		{Path: "store/legacy-autopsy.sql", Kind: "DENOMINATOR-QUERY", RecordID: checkID, Version: "candidate/1", Fingerprint: denominator},
	})
}

// handbookFileFingerprints reads every handbook file in sorted order and
// returns its FILE-TRANSPORT bindings plus one aggregate fingerprint.
func handbookFileFingerprints(dir string) ([]artifact.EvidenceBinding, string, error) {
	if dir == "" {
		return nil, capabilities.Hash([]byte("FILE-TRANSPORT|none")), nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, "", err
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	bindings := make([]artifact.EvidenceBinding, 0, len(names))
	var digests []string
	for _, name := range names {
		content, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, "", err
		}
		path := "handbook/" + name
		fingerprint, err := artifact.FileTransportFingerprint(path, content)
		if err != nil {
			return nil, "", err
		}
		bindings = append(bindings, artifact.EvidenceBinding{Path: path, Kind: "FILE-TRANSPORT", RecordID: "None", Version: "None", Fingerprint: fingerprint})
		digests = append(digests, path+"="+fingerprint)
	}
	return bindings, artifact.CanonicalSetFingerprint("HANDBOOK-FILES", "handbook-dir", digests), nil
}

func semanticRecordCounts(db string, doc *cdl.EIRDoc) [][2]string {
	type entry struct {
		fieldID string
		label   string
	}
	var entries []entry
	for _, catalog := range synthesis.Catalogs {
		entries = append(entries, entry{fieldID: catalog.FieldID, label: catalog.IDType})
	}
	entries = append(entries, entry{fieldID: "PERSONA-PROFILE", label: "PRF"})
	sort.Slice(entries, func(i, j int) bool { return entries[i].fieldID < entries[j].fieldID })
	rows := make([][2]string, 0, len(entries))
	for _, item := range entries {
		ids, err := store.List(db, doc, item.fieldID)
		if err != nil {
			continue
		}
		rows = append(rows, [2]string{item.fieldID, fmt.Sprintf("%d", len(ids))})
	}
	return rows
}

func decisionContentCounts(db string, doc *cdl.EIRDoc) [][2]string {
	names, _, _, err := work.FieldContract(doc, "DECISION")
	if err != nil {
		return nil
	}
	ids, err := store.List(db, doc, "DECISION")
	if err != nil {
		return nil
	}
	counts := map[string]int{}
	for _, id := range ids {
		payload, err := store.Get(db, doc, "DECISION", id)
		if err != nil {
			continue
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			continue
		}
		counts[values["type"]]++
	}
	var rows [][2]string
	for decisionType, count := range counts {
		rows = append(rows, [2]string{decisionType, fmt.Sprintf("%d", count)})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i][0] < rows[j][0] })
	return rows
}

func sanitizeCell(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
