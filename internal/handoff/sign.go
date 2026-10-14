package handoff

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/artifact"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/handbook"
)

const (
	outerEnvelopeStart = "<!-- ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: OUTER-BUNDLE-MANIFEST -->"
	outerEnvelopeEnd   = "<!-- END-ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: OUTER-BUNDLE-MANIFEST -->"
)

// SignOuter replaces only the outer manifest signature envelope with the
// operator's approval. The hashed payload and its fingerprint never change.
func SignOuter(outerPath, artifactPath, approver, reference, timestamp, audit string) (string, string, error) {
	content, err := os.ReadFile(outerPath)
	if err != nil {
		return "", "", err
	}
	text := string(content)
	start := strings.Index(text, outerEnvelopeStart)
	end := strings.Index(text, outerEnvelopeEnd)
	if start < 0 || end < start {
		return "", "", fmt.Errorf("handoff: outer manifest has no signature envelope")
	}
	if reportField(content, "exit-e-status") != "Passed" {
		return "", "", fmt.Errorf("handoff: outer manifest is not Passed")
	}
	declared := reportField(content, "artifact-payload-fingerprint")
	recomputed, err := artifact.PayloadFingerprint(OuterManifestArtifactType, artifactPath, content)
	if err != nil {
		return "", "", err
	}
	if declared != recomputed {
		return "", "", fmt.Errorf("handoff: outer payload fingerprint does not recompute")
	}
	envelope := outerEnvelopeStart + "\n\n" +
		fmt.Sprintf("- **Required Human Signatories / Signature Timestamps:** %s at %s\n", approver, timestamp) +
		fmt.Sprintf("- **Envelope Audit:** %s\n", audit) +
		outerEnvelopeEnd
	signed := text[:start] + envelope + text[end+len(outerEnvelopeEnd):]
	if err := os.WriteFile(outerPath, []byte(signed), 0o644); err != nil {
		return "", "", err
	}
	after, err := artifact.PayloadFingerprint(OuterManifestArtifactType, artifactPath, []byte(signed))
	if err != nil {
		return "", "", err
	}
	if after != recomputed {
		return "", "", fmt.Errorf("handoff: signing changed the outer payload fingerprint")
	}
	envelopeFingerprint := outerEnvelopeFingerprint(artifactPath, envelope)
	return recomputed, envelopeFingerprint, nil
}

// outerEnvelopeFingerprint is the pilot canonical fingerprint of the outer
// signature envelope region.
func outerEnvelopeFingerprint(artifactPath, envelope string) string {
	canonical, err := capabilities.CanonicalMarkdown([]byte(envelope))
	if err != nil {
		return ""
	}
	return capabilities.Hash([]byte("OUTER-ENVELOPE|" + artifactPath + "\n" + string(canonical)))
}

// ReceiptCheck is one §15.5 final check row.
type ReceiptCheck struct {
	ID        string
	Result    string
	Offenders []string
}

// FinalReceipt is the non-authoritative final verification receipt.
type FinalReceipt struct {
	Payload     []byte
	Fingerprint string
	Overall     string
	Checks      []ReceiptCheck
}

// SignAndVerifyOptions configures the step-6 signing and receipt emission.
type SignAndVerifyOptions struct {
	ProtocolVersion         string
	OuterPath               string
	OuterArtifactPath       string
	ReportPath              string
	ReportArtifactPath      string
	ManifestPath            string
	ManifestArtifactPath    string
	ReadinessPath           string
	ReadinessArtifactPath   string
	CertificatePath         string
	CertificateArtifactPath string
	PackageDir              string
	ValidatorVersion        string
}

// BuildFinalReceipt runs the deterministic §15.5 final checks over the signed
// bundle bytes and emits the non-authoritative verification receipt.
func BuildFinalReceipt(doc *cdl.EIRDoc, opts SignAndVerifyOptions, timestamp string) (FinalReceipt, error) {
	var receipt FinalReceipt
	// Bundle-bytes-only: when a materialized package is available, every
	// verification input is read from the package itself, never from external
	// paths.
	load := func(external, artifactPath string) ([]byte, error) {
		if opts.PackageDir != "" {
			return os.ReadFile(filepath.Join(opts.PackageDir, filepath.FromSlash(artifactPath)))
		}
		return os.ReadFile(external)
	}
	outer, err := load(opts.OuterPath, opts.OuterArtifactPath)
	if err != nil {
		return receipt, err
	}
	report, err := load(opts.ReportPath, opts.ReportArtifactPath)
	if err != nil {
		return receipt, err
	}
	manifest, err := load(opts.ManifestPath, opts.ManifestArtifactPath)
	if err != nil {
		return receipt, err
	}
	readiness, err := load(opts.ReadinessPath, opts.ReadinessArtifactPath)
	if err != nil {
		return receipt, err
	}
	certificate, err := load(opts.CertificatePath, opts.CertificateArtifactPath)
	if err != nil {
		return receipt, err
	}
	snapshot := reportField(outer, "system-protocol-snapshot")
	add := func(id string, fine string, offenders ...string) {
		result := "Pass"
		if len(offenders) > 0 {
			result = "Fail"
		}
		receipt.Checks = append(receipt.Checks, ReceiptCheck{ID: id, Result: result, Offenders: offenders})
		_ = fine
	}

	// FINAL-01: outer schema, payload fingerprint, Passed literal.
	var outerProblems []string
	if err := artifact.ValidatePayload(doc, OuterManifestArtifactType, outer); err != nil {
		outerProblems = append(outerProblems, err.Error())
	}
	if reportField(outer, "exit-e-status") != "Passed" {
		outerProblems = append(outerProblems, "EXIT-E-STATUS is not Passed")
	}
	declared := reportField(outer, "artifact-payload-fingerprint")
	recomputed, err := artifact.PayloadFingerprint(OuterManifestArtifactType, opts.OuterArtifactPath, outer)
	if err != nil || recomputed != declared {
		outerProblems = append(outerProblems, "payload fingerprint does not recompute")
	}
	add("FINAL-01-OUTER-PAYLOAD", "", outerProblems...)

	// FINAL-02: exact snapshot equality across snapshot-bearing artifacts.
	var snapshotProblems []string
	for name, payload := range map[string][]byte{
		"exit-e-candidate-report":         report,
		"exit-e-content-readiness-report": readiness,
		"outer-bundle-manifest":           outer,
	} {
		declared := reportField(payload, "system-protocol-snapshot")
		if declared != snapshot {
			snapshotProblems = append(snapshotProblems, name)
		}
	}
	add("FINAL-02-SNAPSHOT", "", snapshotProblems...)

	// FINAL-03/04: member set and fingerprints from the materialized package.
	members, err := packageMembers(opts.PackageDir)
	if err != nil {
		return receipt, err
	}
	var setProblems, fingerprintProblems []string
	declaredMembers, err := declaredMemberEntries(outer)
	if err != nil {
		setProblems = append(setProblems, err.Error())
	} else {
		declaredMap := map[string]string{}
		for _, entry := range declaredMembers {
			if entry.Path == opts.OuterArtifactPath {
				setProblems = append(setProblems, "outer manifest lists itself")
			}
			declaredMap[entry.Path] = entry.Fingerprint
		}
		actual, err := artifact.PackageMembers(members)
		if err != nil {
			return receipt, err
		}
		actualMap := map[string]string{}
		for _, entry := range actual {
			if entry.Path == opts.OuterArtifactPath {
				continue
			}
			actualMap[entry.Path] = entry.Fingerprint
		}
		for path, want := range declaredMap {
			got, ok := actualMap[path]
			if !ok {
				setProblems = append(setProblems, "missing "+path)
				continue
			}
			if got != want {
				fingerprintProblems = append(fingerprintProblems, path)
			}
		}
		for path := range actualMap {
			if _, ok := declaredMap[path]; !ok {
				setProblems = append(setProblems, "extra "+path)
			}
		}
	}
	add("FINAL-03-MEMBER-SET", "", setProblems...)
	add("FINAL-04-MEMBER-FINGERPRINTS", "", fingerprintProblems...)

	// FINAL-05: envelopes: HBK certification statuses, CNF references, outer signature.
	var envelopeProblems []string
	for _, member := range members {
		if !strings.HasPrefix(member.Path, "handbook/sections/") {
			continue
		}
		status, confirmation := handbook.EnvelopeStatus(member.Bytes)
		if status != "B-CONFIRMED" || confirmation == "" {
			envelopeProblems = append(envelopeProblems, member.Path)
		}
	}
	outerEnvelope := outerSignatureEnvelope(outer)
	if outerEnvelope == "" {
		envelopeProblems = append(envelopeProblems, "outer signature envelope")
	}
	add("FINAL-05-ENVELOPES", "", envelopeProblems...)

	// FINAL-06: chain references recompute.
	var chainProblems []string
	stages := []artifact.ChainStage{
		{Type: CandidateReportArtifactType, Path: opts.ReportArtifactPath, Payload: report},
		{Type: CandidateManifestArtifactType, Path: opts.ManifestArtifactPath, Payload: manifest},
		{Type: ContentReadinessArtifactType, Path: opts.ReadinessArtifactPath, Payload: readiness},
		{Type: ScopeCertificateArtifactType, Path: opts.CertificateArtifactPath, Payload: certificate},
		{Type: OuterManifestArtifactType, Path: opts.OuterArtifactPath, Payload: outer},
	}
	if err := artifact.VerifyChain(doc, stages); err != nil {
		chainProblems = append(chainProblems, err.Error())
	}
	add("FINAL-06-CHAIN", "", chainProblems...)

	// FINAL-07: final registry and pre-signature inputs recompute.
	var registryProblems []string
	finalRegistry := artifact.RegistryFingerprint("EXIT-E-FINAL-CHECKS-v1", artifact.RegistryCheckIDs(doc, "FINAL-CHECK"))
	if !strings.Contains(string(outer), finalRegistry) {
		registryProblems = append(registryProblems, "final registry fingerprint")
	}
	preSignature, err := artifact.PreSignatureInputs(doc, outer, stages[2], stages[3])
	if err != nil || !strings.Contains(string(outer), preSignature) {
		registryProblems = append(registryProblems, "pre-signature input-set fingerprint")
	}
	add("FINAL-07-REGISTRY-AND-INPUTS", "", registryProblems...)

	// FINAL-08: readiness conditions.
	var conditionProblems []string
	if result := reportField(readiness, "content-readiness-result"); result != "Ready-For-Certificate" {
		conditionProblems = append(conditionProblems, "readiness "+result)
	}
	for _, check := range receipt.Checks {
		if check.Result != "Pass" {
			conditionProblems = append(conditionProblems, check.ID)
		}
	}
	add("FINAL-08-CONDITIONS", "", conditionProblems...)

	// FINAL-09: Exit E remained Pending through step 5.
	var sequenceProblems []string
	if state := reportField(readiness, "exit-e-prepackage-state"); state != "Pending" {
		sequenceProblems = append(sequenceProblems, "readiness state "+state)
	}
	add("FINAL-09-SEQUENCE", "", sequenceProblems...)

	overall := "Pass"
	for _, check := range receipt.Checks {
		if check.Result != "Pass" {
			overall = "Fail"
		}
	}
	receipt.Overall = overall

	registryIdentity := "EXIT-E-CHECKS-v1 " + artifact.CheckRegistryFingerprint(doc) + "; EXIT-E-FINAL-CHECKS-v1 " + finalRegistry
	outerEnvelopeFp := outerEnvelopeFingerprint(opts.OuterArtifactPath, outerSignatureEnvelope(outer))
	var bindings []string
	for _, member := range members {
		bindings = append(bindings, member.Path)
	}
	sort.Strings(bindings)

	var b strings.Builder
	fmt.Fprintf(&b, "- **Protocol / Artifact / Check Registry Versions and Fingerprints:** %s; artifact-schema %s\n", registryIdentity, artifact.ArtifactSchemaVersion)
	fmt.Fprintf(&b, "- **Normalized Bundle Identity:** %s\n", opts.PackageDir)
	fmt.Fprintf(&b, "- **Outer Payload / Canonical Outer-Envelope Fingerprints:** %s %s\n", recomputed, outerEnvelopeFp)
	fmt.Fprintf(&b, "- **System / Protocol / Environment / Snapshot:** %s\n", snapshot)
	fmt.Fprintf(&b, "- **Package-Member / Input Bindings and Input-Set Fingerprint:** %s\n", strings.Join(bindings, ", "))
	b.WriteString("- **Check Results:**\n")
	b.WriteString("  | Check ID | Result | Offending Bindings |\n")
	b.WriteString("  | --- | --- | --- |\n")
	for _, check := range receipt.Checks {
		fmt.Fprintf(&b, "  | %s | %s | %s |\n", check.ID, check.Result, strings.Join(check.Offenders, "; "))
	}
	fmt.Fprintf(&b, "- **Overall Result:** %s\n", overall)
	fmt.Fprintf(&b, "- **Validator Implementation Version / Timestamp:** %s at %s\n", opts.ValidatorVersion, timestamp)
	b.WriteString("- **Deterministic Result Fingerprint:** PENDING\n")
	payload := []byte(b.String())
	fingerprint := capabilities.Hash([]byte("FINAL-VERIFY-RESULT|" + string(payload)))
	payload = []byte(strings.Replace(string(payload),
		"- **Deterministic Result Fingerprint:** PENDING", "- **Deterministic Result Fingerprint:** "+fingerprint, 1))
	receipt.Payload = payload
	receipt.Fingerprint = fingerprint
	return receipt, nil
}

// packageMembers reads the materialized package files.
func packageMembers(dir string) ([]artifact.Member, error) {
	if dir == "" {
		return nil, nil
	}
	var members []artifact.Member
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		members = append(members, artifact.Member{Path: filepath.ToSlash(relative), Bytes: content})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(members, func(i, j int) bool { return members[i].Path < members[j].Path })
	return members, nil
}

// declaredMemberEntries parses the outer Package Members table.
func declaredMemberEntries(outer []byte) ([]artifact.ManifestEntry, error) {
	entries := []artifact.ManifestEntry{}
	inTable := false
	header := true
	for _, line := range strings.Split(string(outer), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- **Package Members:**") {
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
		if len(cells) != 3 || strings.Trim(cells[0], "-: ") == "" {
			continue
		}
		entries = append(entries, artifact.ManifestEntry{Path: cells[1], Fingerprint: cells[2]})
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("handoff: outer manifest has no Package Members table")
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
	return entries, nil
}

// outerSignatureEnvelope returns the outer signature envelope region, or "" when
// it still says None or lacks an approver.
func outerSignatureEnvelope(outer []byte) string {
	text := string(outer)
	start := strings.Index(text, outerEnvelopeStart)
	end := strings.Index(text, outerEnvelopeEnd)
	if start < 0 || end < start {
		return ""
	}
	region := text[start : end+len(outerEnvelopeEnd)]
	if strings.Contains(region, "Signature Timestamps:** None") {
		return ""
	}
	return region
}
