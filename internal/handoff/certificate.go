package handoff

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/artifact"
	"github.com/nckirik/legacy-autopsy/internal/assurance"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// ScopeCertificateArtifactType is the generic artifact type of the certificate.
const ScopeCertificateArtifactType = "SCOPE-CERTIFICATE"

// ScopeCertificateDefaultPath is the default package-relative artifact path.
const ScopeCertificateDefaultPath = "exit-e/SCOPE-CERTIFICATE.md"

// CertificateOptions are the runtime-owned certificate directives.
type CertificateOptions struct {
	Namespace             string
	ProtocolVersion       string
	EnvironmentsSnapshots string
	StorageBinding        string
	ReportPath            string
	ReportArtifactPath    string
	ManifestPath          string
	ManifestArtifactPath  string
	ReadinessPath         string
	ReadinessArtifactPath string
	ArtifactPath          string
}

// ScopeCertificate is one emitted certificate artifact.
type ScopeCertificate struct {
	Payload         []byte
	Fingerprint     string
	InstanceBinding string
}

// BuildScopeCertificate emits the stage-5 certificate after the content-readiness
// report is Ready-For-Certificate. A blocked readiness report refuses emission.
func BuildScopeCertificate(db string, doc *cdl.EIRDoc, opts CertificateOptions) (ScopeCertificate, error) {
	var certificate ScopeCertificate
	if opts.ArtifactPath == "" {
		opts.ArtifactPath = ScopeCertificateDefaultPath
	}
	if opts.ReportArtifactPath == "" {
		opts.ReportArtifactPath = CandidateReportDefaultPath
	}
	if opts.ManifestArtifactPath == "" {
		opts.ManifestArtifactPath = CandidateManifestDefaultPath
	}
	if opts.ReadinessArtifactPath == "" {
		opts.ReadinessArtifactPath = ContentReadinessDefaultPath
	}
	snapshot := SystemSnapshotValue(opts.Namespace, opts.ProtocolVersion, opts.EnvironmentsSnapshots)
	if strings.TrimSpace(snapshot) == "" {
		return certificate, fmt.Errorf("handoff: scope certificate requires a pinned system/protocol/snapshot identity")
	}
	reportBytes, err := os.ReadFile(opts.ReportPath)
	if err != nil {
		return certificate, err
	}
	manifestBytes, err := os.ReadFile(opts.ManifestPath)
	if err != nil {
		return certificate, err
	}
	readinessBytes, err := os.ReadFile(opts.ReadinessPath)
	if err != nil {
		return certificate, err
	}
	if result := reportField(readinessBytes, "content-readiness-result"); result != "Ready-For-Certificate" {
		return certificate, fmt.Errorf("handoff: scope certificate requires a Ready-For-Certificate readiness report, got %q", result)
	}
	reportFingerprint, err := artifact.PayloadFingerprint(CandidateReportArtifactType, opts.ReportArtifactPath, reportBytes)
	if err != nil {
		return certificate, err
	}
	manifestFingerprint, err := artifact.PayloadFingerprint(CandidateManifestArtifactType, opts.ManifestArtifactPath, manifestBytes)
	if err != nil {
		return certificate, err
	}
	readinessFingerprint, err := artifact.PayloadFingerprint(ContentReadinessArtifactType, opts.ReadinessArtifactPath, readinessBytes)
	if err != nil {
		return certificate, err
	}
	if declared := reportField(readinessBytes, "exit-e-candidate-report-fingerprint"); declared != reportFingerprint {
		return certificate, fmt.Errorf("handoff: readiness report references a different candidate report")
	}
	if declared := reportField(readinessBytes, "candidate-payload-manifest-fingerprint"); declared != manifestFingerprint {
		return certificate, fmt.Errorf("handoff: readiness report references a different candidate manifest")
	}
	if declared := reportField(readinessBytes, "exit-e-prepackage-state"); declared != "Pending" {
		return certificate, fmt.Errorf("handoff: readiness report does not leave Exit E Pending")
	}
	registryFingerprint := artifact.CheckRegistryFingerprint(doc)
	exitAFingerprint := reportField(reportBytes, "exit-a-report-fingerprint")

	cnfIDs, err := store.List(db, doc, "CONFIRMATION")
	if err != nil {
		return certificate, err
	}
	type confirmedBinding struct {
		recordType  string
		recordID    string
		version     string
		fingerprint string
		cnfID       string
		cnfPayload  string
	}
	bindings := ParseManifestBindings(manifestBytes)
	byToken := map[string]ManifestBinding{}
	for _, binding := range bindings {
		byToken[binding.BindingToken()] = binding
	}
	var confirmed []confirmedBinding
	seen := map[string]string{}
	for _, id := range cnfIDs {
		payload, err := store.Get(db, doc, "CONFIRMATION", id)
		if err != nil {
			return certificate, err
		}
		cnfPayload, err := store.FingerprintOf(db, doc, "CONFIRMATION", id)
		if err != nil {
			return certificate, err
		}
		for _, token := range ConfirmationBindings(payload) {
			binding, ok := byToken[token]
			if !ok {
				continue
			}
			if existing, ok := seen[binding.RecordID]; ok {
				return certificate, fmt.Errorf("handoff: binding %s is confirmed by both %s and %s", binding.RecordID, existing, id)
			}
			seen[binding.RecordID] = id
			confirmed = append(confirmed, confirmedBinding{
				recordType:  binding.RecordType,
				recordID:    binding.RecordID,
				version:     binding.Version,
				fingerprint: binding.Fingerprint,
				cnfID:       id,
				cnfPayload:  cnfPayload,
			})
		}
	}
	sort.Slice(confirmed, func(i, j int) bool {
		if confirmed[i].recordType != confirmed[j].recordType {
			return confirmed[i].recordType < confirmed[j].recordType
		}
		return confirmed[i].recordID < confirmed[j].recordID
	})

	counts, err := sourceCounts(db, doc)
	if err != nil {
		return certificate, err
	}
	summary, err := assurance.CoverageSummary(db, doc)
	if err != nil {
		return certificate, err
	}
	mapped := summary.CoverageByDisposition["Mapped"]
	excluded := summary.CoverageByDisposition["Approved-Excluded"]
	total := 0
	for _, count := range summary.CoverageByDisposition {
		total += count
	}
	reconciliationFingerprint := artifact.ValidationSummaryFingerprint("export-reconciliation", snapshot, [][2]string{
		{"mapped", fmt.Sprintf("%d", mapped)},
		{"approved-excluded", fmt.Sprintf("%d", excluded)},
		{"total", fmt.Sprintf("%d", total)},
	})

	var b strings.Builder
	b.WriteString("<!-- ARTIFACT-PAYLOAD: " + ScopeCertificateArtifactType + " -->\n\n")
	b.WriteString("# Scope Certificate\n\n")
	fmt.Fprintf(&b, "- **ARTIFACT-TYPE:** %s\n", ScopeCertificateArtifactType)
	fmt.Fprintf(&b, "- **ARTIFACT-PATH:** %s\n", opts.ArtifactPath)
	fmt.Fprintf(&b, "- **PROTOCOL / ARTIFACT-SCHEMA VERSION:** %s / %s\n", strings.TrimSpace(opts.ProtocolVersion), artifact.ArtifactSchemaVersion)
	fmt.Fprintf(&b, "- **System / Protocol / Snapshot:** %s\n", snapshot)
	b.WriteString("- **Included Scope Dimensions:**\n")
	b.WriteString("  | Dimension | Value |\n  | --- | --- |\n")
	fmt.Fprintf(&b, "  | system-namespace | %s |\n", opts.Namespace)
	fmt.Fprintf(&b, "  | protocol-version | %s |\n", strings.TrimSpace(opts.ProtocolVersion))
	fmt.Fprintf(&b, "  | environments-snapshots | %s |\n", opts.EnvironmentsSnapshots)
	fmt.Fprintf(&b, "  | canonical-storage-binding | %s |\n", opts.StorageBinding)
	b.WriteString("- **Source Counts by Kind:**\n")
	b.WriteString("  | SRC Kind | Gross Discovered | Approved Excluded | Effective Denominator | Covered |\n")
	b.WriteString("  | --- | --- | --- | --- | --- |\n")
	for _, row := range counts {
		fmt.Fprintf(&b, "  | %s | %d | %d | %d | %d |\n", row[0], row[1], row[2], row[3], row[4])
	}
	b.WriteString("- **Approved Exclusions and Residual Risks:**\n")
	b.WriteString("  | DEC ID | Excluded Unit ID | Decision-Content Version | Decision-Content Fingerprint | Approval Envelope Binding | Risk ID | Exact Residual Risk |\n")
	b.WriteString("  | --- | --- | --- | --- | --- | --- | --- |\n")
	b.WriteString("- **Export Reconciliation:** ")
	fmt.Fprintf(&b, "%d mapped, %d approved-excluded, %d total normalized units; %s\n", mapped, excluded, total, reconciliationFingerprint)
	fmt.Fprintf(&b, "- **Exit A Report Fingerprint:** %s\n", exitAFingerprint)
	fmt.Fprintf(&b, "- **Exit E Candidate Report Fingerprint:** %s\n", reportFingerprint)
	fmt.Fprintf(&b, "- **Candidate Payload Manifest Fingerprint:** %s\n", manifestFingerprint)
	fmt.Fprintf(&b, "- **Exit E Content-Readiness Report Fingerprint:** %s\n", readinessFingerprint)
	b.WriteString("- **Confirmed Semantic Payload Bindings:**\n")
	b.WriteString("  | Record Type | Typed Record ID | Semantic Record Version | Semantic-Content Fingerprint | CNF ID | CNF Payload Fingerprint |\n")
	b.WriteString("  | --- | --- | --- | --- | --- | --- |\n")
	for _, binding := range confirmed {
		fmt.Fprintf(&b, "  | %s | %s | %s | %s | %s | %s |\n",
			binding.recordType, binding.recordID, binding.version, binding.fingerprint, binding.cnfID, binding.cnfPayload)
	}
	b.WriteString("- **Known Residual Risks:**\n")
	b.WriteString("  | Risk ID | Affected IDs | Source DEC ID | Exact Residual Risk | Disclosure Path |\n")
	b.WriteString("  | --- | --- | --- | --- | --- |\n")
	b.WriteString("- **ARTIFACT-PAYLOAD-FINGERPRINT:** PENDING\n\n")
	b.WriteString("<!-- END-ARTIFACT-PAYLOAD: " + ScopeCertificateArtifactType + " -->\n\n")
	b.WriteString("<!-- ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: " + ScopeCertificateArtifactType + " -->\n\n")
	b.WriteString("- **Signature / Signatories / Timestamps:** None\n")
	b.WriteString("- **Envelope Audit:** not attached before certificate emission\n")
	b.WriteString("<!-- END-ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: " + ScopeCertificateArtifactType + " -->\n")
	_ = registryFingerprint
	payload := []byte(b.String())
	fingerprint, err := artifact.PayloadFingerprint(ScopeCertificateArtifactType, opts.ArtifactPath, payload)
	if err != nil {
		return certificate, err
	}
	final := []byte(strings.Replace(string(payload),
		"- **ARTIFACT-PAYLOAD-FINGERPRINT:** PENDING", "- **ARTIFACT-PAYLOAD-FINGERPRINT:** "+fingerprint, 1))
	recomputed, err := artifact.PayloadFingerprint(ScopeCertificateArtifactType, opts.ArtifactPath, final)
	if err != nil {
		return certificate, err
	}
	if recomputed != fingerprint {
		return certificate, fmt.Errorf("handoff: certificate fingerprint does not recompute")
	}
	if err := artifact.ValidatePayload(doc, ScopeCertificateArtifactType, final); err != nil {
		return certificate, err
	}
	instance, err := artifact.InstanceBinding(ScopeCertificateArtifactType, opts.ArtifactPath, fingerprint)
	if err != nil {
		return certificate, err
	}
	certificate.Payload = final
	certificate.Fingerprint = fingerprint
	certificate.InstanceBinding = instance
	return certificate, nil
}

// sourceCounts aggregates inventory and coverage rows by SRC kind.
func sourceCounts(db string, doc *cdl.EIRDoc) ([][5]any, error) {
	inventoryNames, _, _, err := work.FieldContract(doc, "SOURCE-INVENTORY-RECORD")
	if err != nil {
		return nil, err
	}
	inventoryIDs, err := store.List(db, doc, "SOURCE-INVENTORY-RECORD")
	if err != nil {
		return nil, err
	}
	gross := map[string]int{}
	for _, id := range inventoryIDs {
		payload, err := store.Get(db, doc, "SOURCE-INVENTORY-RECORD", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(inventoryNames, payload)
		if err != nil {
			return nil, err
		}
		kind := strings.ToUpper(strings.TrimSpace(values["kind"]))
		if kind == "" {
			kind = "OTHER"
		}
		gross[kind]++
	}
	coverageNames, _, _, err := work.FieldContract(doc, "SOURCE-COVERAGE-ROW")
	if err != nil {
		return nil, err
	}
	coverageIDs, err := store.List(db, doc, "SOURCE-COVERAGE-ROW")
	if err != nil {
		return nil, err
	}
	covered := map[string]int{}
	for _, id := range coverageIDs {
		payload, err := store.Get(db, doc, "SOURCE-COVERAGE-ROW", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(coverageNames, payload)
		if err != nil {
			return nil, err
		}
		src := strings.Fields(values["src-version"])
		if len(src) == 0 {
			continue
		}
		kind := srcKind(src[0])
		if strings.EqualFold(strings.TrimSpace(values["disposition"]), "Mapped") {
			covered[kind]++
		}
	}
	kinds := make([]string, 0, len(gross))
	for kind := range gross {
		kinds = append(kinds, kind)
	}
	sort.Strings(kinds)
	rows := make([][5]any, 0, len(kinds))
	for _, kind := range kinds {
		rows = append(rows, [5]any{kind, gross[kind], 0, gross[kind], covered[kind]})
	}
	return rows, nil
}

func srcKind(id string) string {
	parts := strings.SplitN(id, "-", 3)
	if len(parts) < 3 {
		return "OTHER"
	}
	return strings.ToUpper(parts[1])
}
