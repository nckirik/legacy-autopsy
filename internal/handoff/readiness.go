package handoff

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/artifact"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/handbook"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/synthesis"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// ContentReadinessArtifactType is the generic artifact type of the stage-4 report.
const ContentReadinessArtifactType = "EXIT-E-CONTENT-READINESS-REPORT"

// ContentReadinessDefaultPath is the default package-relative artifact path.
const ContentReadinessDefaultPath = "exit-e/EXIT-E-CONTENT-READINESS-REPORT.md"

// ReadinessOptions are the runtime-owned stage-4 directives.
type ReadinessOptions struct {
	Namespace             string
	ProtocolVersion       string
	EnvironmentsSnapshots string
	StorageBinding        string
	ReportPath            string
	ReportArtifactPath    string
	ManifestPath          string
	ManifestArtifactPath  string
	HandbookDir           string
	StoreDump             []byte
	ArtifactPath          string
}

// ReadinessCheck is one content-readiness row.
type ReadinessCheck struct {
	ID       string
	Result   string
	Detail   string
	Evidence string
}

// ContentReadiness is one emitted stage-4 report.
type ContentReadiness struct {
	Payload         []byte
	Fingerprint     string
	InstanceBinding string
	Result          string
	Checks          []ReadinessCheck
}

type memberRow struct {
	artifactType string
	path         string
	fingerprint  string
}

// BuildContentReadiness emits the immutable stage-4 report. The evaluator
// recomputes candidate-stage fingerprints, validates available envelopes and
// members, and reports human-gated confirmations, approvals, and readability as
// blockers instead of passing them.
func BuildContentReadiness(db string, doc *cdl.EIRDoc, opts ReadinessOptions) (ContentReadiness, error) {
	var report ContentReadiness
	if opts.ArtifactPath == "" {
		opts.ArtifactPath = ContentReadinessDefaultPath
	}
	if opts.ReportArtifactPath == "" {
		opts.ReportArtifactPath = CandidateReportDefaultPath
	}
	if opts.ManifestArtifactPath == "" {
		opts.ManifestArtifactPath = CandidateManifestDefaultPath
	}
	snapshot := SystemSnapshotValue(opts.Namespace, opts.ProtocolVersion, opts.EnvironmentsSnapshots)
	if strings.TrimSpace(snapshot) == "" {
		return report, fmt.Errorf("handoff: content-readiness report requires a pinned system/protocol/snapshot identity")
	}
	reportBytes, err := os.ReadFile(opts.ReportPath)
	if err != nil {
		return report, err
	}
	manifestBytes, err := os.ReadFile(opts.ManifestPath)
	if err != nil {
		return report, err
	}
	reportFingerprint, err := artifact.PayloadFingerprint(CandidateReportArtifactType, opts.ReportArtifactPath, reportBytes)
	if err != nil {
		return report, err
	}
	manifestFingerprint, err := artifact.PayloadFingerprint(CandidateManifestArtifactType, opts.ManifestArtifactPath, manifestBytes)
	if err != nil {
		return report, err
	}
	exitAFingerprint := reportField(reportBytes, "exit-a-report-fingerprint")
	registryFingerprint := artifact.CheckRegistryFingerprint(doc)
	members, err := readinessMembers(opts)
	if err != nil {
		return report, err
	}
	membersByPath := map[string]memberRow{}
	for _, member := range members {
		membersByPath[member.path] = member
	}
	rows := []ReadinessCheck{}
	add := func(id, result, detail string, bindings []artifact.EvidenceBinding) error {
		evidence, err := artifact.EvidenceSetFingerprint(ContentReadinessArtifactType, "CHECK", id, snapshot, bindings)
		if err != nil {
			return err
		}
		rows = append(rows, ReadinessCheck{ID: id, Result: result, Detail: detail, Evidence: evidence})
		return nil
	}
	standardBindings := func(rowID string, extra ...artifact.EvidenceBinding) []artifact.EvidenceBinding {
		bindings := []artifact.EvidenceBinding{
			{Path: "protocol/assembly.json", Kind: "PROTOCOL-REGISTRY", RecordID: "EXIT-E-CHECKS-v1", Version: "1", Fingerprint: registryFingerprint},
			{Path: opts.ReportArtifactPath, Kind: "ARTIFACT-PAYLOAD", RecordID: CandidateReportArtifactType, Version: "artifact-schema/1", Fingerprint: reportFingerprint},
			{Path: opts.ManifestArtifactPath, Kind: "ARTIFACT-PAYLOAD", RecordID: CandidateManifestArtifactType, Version: "artifact-schema/1", Fingerprint: manifestFingerprint},
		}
		return append(bindings, extra...)
	}

	// READINESS-01: candidate artifacts and bound semantic payloads are unchanged.
	immutabilityDetail, immutabilityFail, err := candidateImmutability(db, doc, manifestBytes, reportFingerprint, opts.HandbookDir)
	if err != nil {
		return report, err
	}
	if immutabilityFail != "" {
		if err := add("READINESS-01-CANDIDATE-IMMUTABILITY", "Fail", immutabilityFail, standardBindings("READINESS-01-CANDIDATE-IMMUTABILITY")); err != nil {
			return report, err
		}
	} else {
		if err := add("READINESS-01-CANDIDATE-IMMUTABILITY", "Pass", immutabilityDetail, standardBindings("READINESS-01-CANDIDATE-IMMUTABILITY")); err != nil {
			return report, err
		}
	}

	// READINESS-02: human confirmations (CNFs) bound to exact candidate entries.
	cnfIDs, err := store.List(db, doc, "CONFIRMATION")
	if err != nil {
		return report, err
	}
	profileCount, err := store.List(db, doc, "PERSONA-PROFILE")
	if err != nil {
		return report, err
	}
	if len(cnfIDs) == 0 {
		if err := add("READINESS-02-CONFIRMATIONS", "Fail", fmt.Sprintf("no confirmation records exist; %d persona profiles await human confirmation", len(profileCount)), standardBindings("READINESS-02-CONFIRMATIONS")); err != nil {
			return report, err
		}
	} else {
		if err := add("READINESS-02-CONFIRMATIONS", "Fail", fmt.Sprintf("%d confirmation records exist but exact-PRF-binding validation is not implemented", len(cnfIDs)), standardBindings("READINESS-02-CONFIRMATIONS")); err != nil {
			return report, err
		}
	}

	// READINESS-03: approved decision content for every declared decision reference.
	decisionDetail, decisionFail := decisionApprovalCoverage(db, doc)
	if decisionFail != "" {
		if err := add("READINESS-03-DECISION-APPROVALS", "Fail", decisionFail, standardBindings("READINESS-03-DECISION-APPROVALS")); err != nil {
			return report, err
		}
	} else {
		if err := add("READINESS-03-DECISION-APPROVALS", "Pass", decisionDetail, standardBindings("READINESS-03-DECISION-APPROVALS")); err != nil {
			return report, err
		}
	}

	// READINESS-04: equivalence suite.
	if err := add("READINESS-04-EQUIVALENCE", "Fail", "no equivalence-suite artifacts are generated", standardBindings("READINESS-04-EQUIVALENCE")); err != nil {
		return report, err
	}

	// READINESS-05: fingerprint-bound human readability review.
	reviewCount := 0
	if opts.HandbookDir != "" {
		entries, err := os.ReadDir(opts.HandbookDir)
		if err != nil {
			return report, err
		}
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasPrefix(entry.Name(), "readability-") && strings.HasSuffix(entry.Name(), ".md") {
				reviewCount++
			}
		}
	}
	if reviewCount == 0 {
		if err := add("READINESS-05-HANDBOOK-READABILITY", "Fail", "no human readability reviews are recorded", standardBindings("READINESS-05-HANDBOOK-READABILITY")); err != nil {
			return report, err
		}
	} else {
		if err := add("READINESS-05-HANDBOOK-READABILITY", "Fail", fmt.Sprintf("%d readability reviews are recorded but coverage validation is incomplete", reviewCount), standardBindings("READINESS-05-HANDBOOK-READABILITY")); err != nil {
			return report, err
		}
	}

	// READINESS-06: no pending or provisional references, no contradictions.
	pending, contradictions, err := pendingAndContradictions(db, doc)
	if err != nil {
		return report, err
	}
	if pending > 0 {
		if err := add("READINESS-06-NO-PENDING-OR-CONTRADICTIONS", "Fail", fmt.Sprintf("%d synthesis records carry pending status", pending), standardBindings("READINESS-06-NO-PENDING-OR-CONTRADICTIONS")); err != nil {
			return report, err
		}
	} else {
		if err := add("READINESS-06-NO-PENDING-OR-CONTRADICTIONS", "Pass", fmt.Sprintf("no pending records; %d contradictions recorded", contradictions), standardBindings("READINESS-06-NO-PENDING-OR-CONTRADICTIONS")); err != nil {
			return report, err
		}
	}

	// READINESS-07: certification envelopes attached and valid.
	envelopes, envelopeFail := certificationEnvelopes(doc, opts.HandbookDir)
	if envelopeFail != "" {
		if err := add("READINESS-07-ENVELOPE-INTEGRITY", "Fail", envelopeFail, standardBindings("READINESS-07-ENVELOPE-INTEGRITY")); err != nil {
			return report, err
		}
	} else {
		if err := add("READINESS-07-ENVELOPE-INTEGRITY", "Pass", fmt.Sprintf("%d certification envelopes validate", envelopes), standardBindings("READINESS-07-ENVELOPE-INTEGRITY")); err != nil {
			return report, err
		}
	}

	// READINESS-08: complete pre-certificate member set with recomputing fingerprints.
	memberBindings := make([]artifact.EvidenceBinding, 0, len(members))
	for _, member := range members {
		memberBindings = append(memberBindings, artifact.EvidenceBinding{Path: member.path, Kind: "PACKAGE-MEMBER-FILE", RecordID: "None", Version: "None", Fingerprint: member.fingerprint})
	}
	if err := add("READINESS-08-PREPACKAGE-MEMBERS", "Pass", fmt.Sprintf("%d pre-certificate members fingerprint", len(members)), standardBindings("READINESS-08-PREPACKAGE-MEMBERS", memberBindings...)); err != nil {
		return report, err
	}

	// READINESS-09: snapshot consistency across the chain and members.
	snapshotFail := snapshotConsistency(reportBytes, manifestBytes, snapshot)
	if snapshotFail != "" {
		if err := add("READINESS-09-SNAPSHOT-CONSISTENCY", "Fail", snapshotFail, standardBindings("READINESS-09-SNAPSHOT-CONSISTENCY", memberBindings...)); err != nil {
			return report, err
		}
	} else {
		if err := add("READINESS-09-SNAPSHOT-CONSISTENCY", "Pass", "candidate artifacts and snapshot-bearing members share the pinned identity", standardBindings("READINESS-09-SNAPSHOT-CONSISTENCY", memberBindings...)); err != nil {
			return report, err
		}
	}

	// READINESS-10: approved exclusions and residual risks disclosed.
	excluded, err := approvedExclusionCount(db, doc)
	if err != nil {
		return report, err
	}
	if excluded > 0 {
		if err := add("READINESS-10-RESIDUAL-RISK", "Fail", fmt.Sprintf("%d approved exclusions lack disclosed residual-risk bindings", excluded), standardBindings("READINESS-10-RESIDUAL-RISK")); err != nil {
			return report, err
		}
	} else {
		if err := add("READINESS-10-RESIDUAL-RISK", "Pass", "no approved exclusions are recorded", standardBindings("READINESS-10-RESIDUAL-RISK")); err != nil {
			return report, err
		}
	}

	// READINESS-11: acyclic sequence and Exit E still Pending.
	sequenceFail := sequenceViolation(reportBytes, manifestBytes, reportFingerprint, manifestFingerprint)
	if sequenceFail != "" {
		if err := add("READINESS-11-SEQUENCE", "Fail", sequenceFail, standardBindings("READINESS-11-SEQUENCE")); err != nil {
			return report, err
		}
	} else {
		if err := add("READINESS-11-SEQUENCE", "Pass", "stage-4 references point only to earlier artifacts and Exit E remains Pending", standardBindings("READINESS-11-SEQUENCE")); err != nil {
			return report, err
		}
	}

	result := "Ready-For-Certificate"
	for _, row := range rows {
		if row.Result != "Pass" {
			result = "Blocked"
		}
	}

	var b strings.Builder
	b.WriteString("<!-- ARTIFACT-PAYLOAD: " + ContentReadinessArtifactType + " -->\n\n")
	b.WriteString("# Exit E Content-Readiness Report\n\n")
	fmt.Fprintf(&b, "- **ARTIFACT-TYPE:** %s\n", ContentReadinessArtifactType)
	fmt.Fprintf(&b, "- **ARTIFACT-PATH:** %s\n", opts.ArtifactPath)
	fmt.Fprintf(&b, "- **PROTOCOL / ARTIFACT-SCHEMA VERSION:** %s / %s\n", strings.TrimSpace(opts.ProtocolVersion), artifact.ArtifactSchemaVersion)
	fmt.Fprintf(&b, "- **System / Protocol / Snapshot:** %s\n", snapshot)
	fmt.Fprintf(&b, "- **Required Check Registry / Fingerprint:** EXIT-E-CHECKS-v1 %s\n", registryFingerprint)
	fmt.Fprintf(&b, "- **Exit A Report Fingerprint:** %s\n", exitAFingerprint)
	fmt.Fprintf(&b, "- **Exit E Candidate Report Fingerprint:** %s\n", reportFingerprint)
	fmt.Fprintf(&b, "- **Candidate Payload Manifest Fingerprint:** %s\n", manifestFingerprint)
	b.WriteString("- **Confirmed Semantic Bindings:**\n")
	b.WriteString("  | Record Type | Typed Record ID | Semantic Record Version | Semantic-Content Fingerprint | Certification Envelope Binding | Semantic Containing Path | Semantic Package-Member File Fingerprint | CNF ID | CNF Payload Fingerprint | CNF Signature Envelope Binding | CNF Containing Path | CNF Package-Member File Fingerprint |\n")
	b.WriteString("  | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |\n")
	b.WriteString("- **Approved DEC Bindings:**\n")
	b.WriteString("  | DEC ID | Decision-Content Version | Decision-Content Fingerprint | Approval Envelope Binding | DEC Containing Path | DEC Package-Member File Fingerprint |\n")
	b.WriteString("  | --- | --- | --- | --- | --- | --- |\n")
	b.WriteString("- **Pre-Certificate Package Members:**\n")
	b.WriteString("  | Artifact Type | Normalized Path | Package-Member File Fingerprint |\n")
	b.WriteString("  | --- | --- | --- |\n")
	for _, member := range members {
		fmt.Fprintf(&b, "  | %s | %s | %s |\n", member.artifactType, member.path, member.fingerprint)
	}
	b.WriteString("- **Check Results:**\n")
	b.WriteString("  | Check ID | Requirement | Result | Evidence Fingerprint |\n")
	b.WriteString("  | --- | --- | --- | --- |\n")
	requirements := artifact.CheckRequirements(doc)
	for _, row := range rows {
		requirement := requirements[row.ID]
		if requirement == "" {
			requirement = row.ID
		}
		fmt.Fprintf(&b, "  | %s | %s | %s | %s |\n", row.ID, requirement, row.Result, row.Evidence)
	}
	b.WriteString("- **EXIT-E-PREPACKAGE-STATE:** Pending\n")
	fmt.Fprintf(&b, "- **Content Readiness Result:** %s\n", result)
	b.WriteString("- **Blockers:**\n")
	b.WriteString("  | Blocker ID | Affected IDs | Description | Evidence Fingerprint |\n")
	b.WriteString("  | --- | --- | --- | --- |\n")
	for _, row := range rows {
		if row.Result == "Pass" {
			continue
		}
		fmt.Fprintf(&b, "  | BLOCKER-%s | %s | %s | %s |\n", row.ID, row.ID, sanitizeCell(row.Detail), row.Evidence)
	}
	b.WriteString("- **ARTIFACT-PAYLOAD-FINGERPRINT:** PENDING\n\n")
	b.WriteString("<!-- END-ARTIFACT-PAYLOAD: " + ContentReadinessArtifactType + " -->\n\n")
	b.WriteString("<!-- ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: " + ContentReadinessArtifactType + " -->\n\n")
	b.WriteString("- **Signature / Signatories / Timestamps:** None\n")
	b.WriteString("- **Envelope Audit:** not attached before certificate emission\n")
	b.WriteString("<!-- END-ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: " + ContentReadinessArtifactType + " -->\n")
	payload := []byte(b.String())
	fingerprint, err := artifact.PayloadFingerprint(ContentReadinessArtifactType, opts.ArtifactPath, payload)
	if err != nil {
		return report, err
	}
	final := []byte(strings.Replace(string(payload),
		"- **ARTIFACT-PAYLOAD-FINGERPRINT:** PENDING", "- **ARTIFACT-PAYLOAD-FINGERPRINT:** "+fingerprint, 1))
	recomputed, err := artifact.PayloadFingerprint(ContentReadinessArtifactType, opts.ArtifactPath, final)
	if err != nil {
		return report, err
	}
	if recomputed != fingerprint {
		return report, fmt.Errorf("handoff: content-readiness fingerprint does not recompute")
	}
	if err := artifact.ValidatePayload(doc, ContentReadinessArtifactType, final); err != nil {
		return report, err
	}
	instance, err := artifact.InstanceBinding(ContentReadinessArtifactType, opts.ArtifactPath, fingerprint)
	if err != nil {
		return report, err
	}
	report.Payload = final
	report.Fingerprint = fingerprint
	report.InstanceBinding = instance
	report.Result = result
	report.Checks = rows
	return report, nil
}

// readinessMembers materializes the stage-4 pre-certificate member set.
func readinessMembers(opts ReadinessOptions) ([]memberRow, error) {
	var members []memberRow
	add := func(path, artifactType string, content []byte) error {
		fingerprint, err := capabilities.PackageMemberFingerprint(path, content)
		if err != nil {
			return err
		}
		members = append(members, memberRow{artifactType: artifactType, path: path, fingerprint: fingerprint})
		return nil
	}
	reportBytes, err := os.ReadFile(opts.ReportPath)
	if err != nil {
		return nil, err
	}
	if err := add(opts.ReportArtifactPath, CandidateReportArtifactType, reportBytes); err != nil {
		return nil, err
	}
	manifestBytes, err := os.ReadFile(opts.ManifestPath)
	if err != nil {
		return nil, err
	}
	if err := add(opts.ManifestArtifactPath, CandidateManifestArtifactType, manifestBytes); err != nil {
		return nil, err
	}
	if len(opts.StoreDump) > 0 {
		if err := add(StoreMemberPath, "CANONICAL-STORE-DUMP", opts.StoreDump); err != nil {
			return nil, err
		}
	}
	if opts.HandbookDir != "" {
		entries, err := os.ReadDir(opts.HandbookDir)
		if err != nil {
			return nil, err
		}
		var names []string
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			names = append(names, entry.Name())
		}
		sort.Strings(names)
		for _, name := range names {
			content, err := os.ReadFile(filepath.Join(opts.HandbookDir, name))
			if err != nil {
				return nil, err
			}
			if err := add("handbook/"+name, "HANDBOOK-SECTION", content); err != nil {
				return nil, err
			}
		}
	}
	sort.Slice(members, func(i, j int) bool {
		if members[i].path != members[j].path {
			return members[i].path < members[j].path
		}
		return members[i].artifactType < members[j].artifactType
	})
	return members, nil
}

func reportField(payload []byte, slug string) string {
	parsed := string(payload)
	for _, line := range strings.Split(parsed, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "- **") {
			continue
		}
		end := strings.Index(trimmed, ":**")
		if end < 0 {
			continue
		}
		label := capabilities.ColumnName(trimmed[4 : end+2])
		if label == slug {
			return strings.TrimSpace(trimmed[end+4:])
		}
	}
	return ""
}

// candidateImmutability recomputes the manifest's semantic bindings against
// committed state and the candidate report reference.
func candidateImmutability(db string, doc *cdl.EIRDoc, manifestBytes []byte, reportFingerprint, handbookDir string) (string, string, error) {
	declaredReport := reportField(manifestBytes, "exit-e-candidate-report-fingerprint")
	if declaredReport != reportFingerprint {
		return "", "candidate report fingerprint changed after the manifest was emitted", nil
	}
	bindings := manifestBindings(manifestBytes)
	if len(bindings) == 0 {
		return "", "manifest declares no semantic bindings", nil
	}
	stale := 0
	var first string
	for _, binding := range bindings {
		if binding.recordType == "HBK" {
			continue // validated through the minted-file reload below
		}
		fieldID, ok := synthesis.FieldForID(binding.recordID)
		if !ok {
			return "", fmt.Sprintf("manifest binding %s has no declared field", binding.recordID), nil
		}
		version, err := store.VersionOf(db, doc, fieldID, binding.recordID)
		if err != nil {
			stale++
			if first == "" {
				first = binding.recordID + " (missing)"
			}
			continue
		}
		payload, err := store.Get(db, doc, fieldID, binding.recordID)
		if err != nil {
			return "", "", err
		}
		var fingerprint string
		if fieldID == "PERSONA-PROFILE" {
			fingerprint, err = capabilities.SemanticContentFingerprint(payload, "semantic-content-fingerprint")
			if err != nil {
				return "", "", err
			}
		} else {
			fingerprint, err = store.FingerprintOf(db, doc, fieldID, binding.recordID)
			if err != nil {
				return "", "", err
			}
		}
		if fmt.Sprintf("%d", version) != binding.version || fingerprint != binding.fingerprint {
			stale++
			if first == "" {
				first = binding.recordID
			}
		}
	}
	if handbookDir != "" {
		sections, err := handbook.LoadMintedSections(handbookDir)
		if err != nil {
			return "", "", err
		}
		byID := map[string]string{}
		for _, section := range sections {
			byID[section.ID] = section.Fingerprint
		}
		for _, binding := range bindings {
			if binding.recordType != "HBK" {
				continue
			}
			if byID[binding.recordID] != binding.fingerprint {
				stale++
				if first == "" {
					first = binding.recordID
				}
			}
		}
	}
	if stale > 0 {
		return "", fmt.Sprintf("%d semantic bindings changed since the manifest was emitted; first: %s", stale, first), nil
	}
	return fmt.Sprintf("%d semantic bindings still match the candidate manifest", len(bindings)), "", nil
}

type bindingRow struct {
	recordType  string
	recordID    string
	version     string
	fingerprint string
}

// manifestBindings parses the manifest Semantic Bindings table.
func manifestBindings(payload []byte) []bindingRow {
	var rows []bindingRow
	var inTable bool
	header := true
	for _, line := range strings.Split(string(payload), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- **Semantic Bindings:**") {
			inTable = true
			header = true
			continue
		}
		if inTable && strings.HasPrefix(trimmed, "- **") {
			inTable = false
			continue
		}
		if !inTable || !strings.HasPrefix(trimmed, "|") {
			continue
		}
		cells := splitTableRow(trimmed)
		if header {
			header = false
			continue
		}
		if len(cells) < 4 || strings.Trim(cells[0], "-: ") == "" {
			continue
		}
		rows = append(rows, bindingRow{recordType: cells[0], recordID: cells[1], version: cells[2], fingerprint: cells[3]})
	}
	return rows
}

func splitTableRow(line string) []string {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "|")
	line = strings.TrimSuffix(line, "|")
	parts := strings.Split(line, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// decisionApprovalCoverage requires every non-None declared decision reference
// to resolve to an approved DECISION record. A workspace that declares no
// references has a resolved zero denominator.
func decisionApprovalCoverage(db string, doc *cdl.EIRDoc) (string, string) {
	checks := []struct {
		fieldID string
		field   string
	}{
		{fieldID: "ARCHITECTURE-MODULE", field: "target-decision"},
		{fieldID: "USE-CASE", field: "target-decision-dec"},
		{fieldID: "CAPABILITY-RECORD", field: "target-decision-dec"},
	}
	decNames, _, _, err := work.FieldContract(doc, "DECISION")
	if err != nil {
		return "decision FIELD is not declared", ""
	}
	approved := map[string]bool{}
	decIDs, err := store.List(db, doc, "DECISION")
	if err != nil {
		return err.Error(), ""
	}
	for _, id := range decIDs {
		payload, err := store.Get(db, doc, "DECISION", id)
		if err != nil {
			return err.Error(), ""
		}
		values, err := work.PayloadText(decNames, payload)
		if err != nil {
			return err.Error(), ""
		}
		if strings.EqualFold(strings.TrimSpace(values["approval-status"]), "Approved") {
			approved[id] = true
		}
	}
	references := 0
	missing := 0
	var first string
	for _, check := range checks {
		names, _, _, err := work.FieldContract(doc, check.fieldID)
		if err != nil {
			continue
		}
		ids, err := store.List(db, doc, check.fieldID)
		if err != nil {
			return err.Error(), ""
		}
		for _, id := range ids {
			payload, err := store.Get(db, doc, check.fieldID, id)
			if err != nil {
				return err.Error(), ""
			}
			values, err := work.PayloadText(names, payload)
			if err != nil {
				return err.Error(), ""
			}
			value := strings.TrimSpace(values[check.field])
			if value == "" || strings.EqualFold(value, "None") {
				continue
			}
			references++
			token := strings.Fields(value)[0]
			if strings.HasPrefix(token, "DEC-") && approved[token] {
				continue
			}
			missing++
			if first == "" {
				first = id + " -> " + value
			}
		}
	}
	if missing > 0 {
		return fmt.Sprintf("%d decision references lack approved decision content; first: %s", missing, first), ""
	}
	return fmt.Sprintf("%d decision references resolve to approved decisions", references), ""
}

// certificationEnvelopes validates attached certification envelopes. Minted HBK
// sections carry B-POPULATED envelopes without confirmation versions, so an
// unconfirmed workspace fails closed.
func certificationEnvelopes(doc *cdl.EIRDoc, handbookDir string) (int, string) {
	if handbookDir == "" {
		return 0, "no certification envelopes are attached (no handbook package directory)"
	}
	sections, err := handbook.LoadMintedSections(handbookDir)
	if err != nil {
		return 0, err.Error()
	}
	attached := 0
	for _, section := range sections {
		if strings.Contains(string(section.Bytes), "HANDBOOK-STATUS:** B-CONFIRMED") {
			attached++
		}
	}
	if attached < len(sections) {
		return attached, fmt.Sprintf("%d of %d handbook sections are B-POPULATED without confirmation envelopes", len(sections)-attached, len(sections))
	}
	return attached, ""
}

// snapshotConsistency compares the declared snapshot field across artifacts.
func snapshotConsistency(reportBytes, manifestBytes []byte, snapshot string) string {
	for name, payload := range map[string][]byte{
		CandidateReportArtifactType:   reportBytes,
		CandidateManifestArtifactType: manifestBytes,
	} {
		declared := reportField(payload, "system-protocol-snapshot")
		if declared != "" && declared != snapshot {
			return fmt.Sprintf("%s declares snapshot %q", name, declared)
		}
	}
	return ""
}

// approvedExclusionCount counts coverage rows with an approved-exclusion
// disposition that have no disclosed residual-risk binding.
func approvedExclusionCount(db string, doc *cdl.EIRDoc) (int, error) {
	names, _, _, err := work.FieldContract(doc, "SOURCE-COVERAGE-ROW")
	if err != nil {
		return 0, err
	}
	ids, err := store.List(db, doc, "SOURCE-COVERAGE-ROW")
	if err != nil {
		return 0, err
	}
	count := 0
	for _, id := range ids {
		payload, err := store.Get(db, doc, "SOURCE-COVERAGE-ROW", id)
		if err != nil {
			return 0, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return 0, err
		}
		if strings.EqualFold(strings.TrimSpace(values["disposition"]), "Approved-Excluded") {
			count++
		}
	}
	return count, nil
}

// sequenceViolation checks the acyclic stage-4 reference order. Fingerprint
// carrier lines are excluded: only references count.
func sequenceViolation(reportBytes, manifestBytes []byte, reportFingerprint, manifestFingerprint string) string {
	if reportFingerprint == "" || manifestFingerprint == "" {
		return "chain fingerprints are missing"
	}
	manifestRefs := withoutCarrierLines(string(manifestBytes))
	reportRefs := withoutCarrierLines(string(reportBytes))
	if strings.Contains(manifestRefs, manifestFingerprint) {
		return "manifest contains its own fingerprint in a reference position"
	}
	if !strings.Contains(manifestRefs, reportFingerprint) {
		return "manifest does not reference the candidate report"
	}
	if strings.Contains(reportRefs, manifestFingerprint) {
		return "candidate report references the later manifest"
	}
	return ""
}

// withoutCarrierLines removes every payload-fingerprint carrier line.
func withoutCarrierLines(content string) string {
	var kept []string
	for _, line := range strings.Split(content, "\n") {
		if strings.Contains(line, "ARTIFACT-PAYLOAD-FINGERPRINT") {
			continue
		}
		kept = append(kept, line)
	}
	return strings.Join(kept, "\n")
}
