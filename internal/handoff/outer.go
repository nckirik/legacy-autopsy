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
)

// OuterManifestArtifactType is the generic artifact type of the outer manifest.
const OuterManifestArtifactType = "OUTER-BUNDLE-MANIFEST"

// OuterManifestDefaultPath is the default package-relative artifact path.
const OuterManifestDefaultPath = "exit-e/OUTER-BUNDLE-MANIFEST.md"

// OuterOptions are the runtime-owned step-6 directives.
type OuterOptions struct {
	Namespace               string
	ProtocolVersion         string
	EnvironmentsSnapshots   string
	StorageBinding          string
	ReportPath              string
	ReportArtifactPath      string
	ManifestPath            string
	ManifestArtifactPath    string
	ReadinessPath           string
	ReadinessArtifactPath   string
	CertificatePath         string
	CertificateArtifactPath string
	HandbookViewDir         string
	HandbookSectionDir      string
	EquivalenceSuitePath    string
	EquivalenceReviewPath   string
	StoreDump               []byte
	ArtifactPath            string
}

// OuterBundle is one emitted outer manifest artifact.
type OuterBundle struct {
	Payload         []byte
	Fingerprint     string
	InstanceBinding string
	Members         []artifact.ManifestEntry
	// Files are the exact finalized member bytes, keyed by package-relative
	// path, so the physical package can be materialized without re-deriving.
	Files []artifact.Member
}

// BuildOuterBundle assembles the final member set and emits the authoritative
// outer manifest with EXIT-E-STATUS: Passed and the recomputed pre-signature
// input-set fingerprint. Signature attachment and §15.5 verification follow.
func BuildOuterBundle(doc *cdl.EIRDoc, opts OuterOptions) (OuterBundle, error) {
	var bundle OuterBundle
	if opts.ArtifactPath == "" {
		opts.ArtifactPath = OuterManifestDefaultPath
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
	if opts.CertificateArtifactPath == "" {
		opts.CertificateArtifactPath = ScopeCertificateDefaultPath
	}
	snapshot := SystemSnapshotValue(opts.Namespace, opts.ProtocolVersion, opts.EnvironmentsSnapshots)
	if strings.TrimSpace(snapshot) == "" {
		return bundle, fmt.Errorf("handoff: outer manifest requires a pinned system/protocol/snapshot identity")
	}
	reportBytes, err := os.ReadFile(opts.ReportPath)
	if err != nil {
		return bundle, err
	}
	manifestBytes, err := os.ReadFile(opts.ManifestPath)
	if err != nil {
		return bundle, err
	}
	readinessBytes, err := os.ReadFile(opts.ReadinessPath)
	if err != nil {
		return bundle, err
	}
	certificateBytes, err := os.ReadFile(opts.CertificatePath)
	if err != nil {
		return bundle, err
	}
	if result := reportField(readinessBytes, "content-readiness-result"); result != "Ready-For-Certificate" {
		return bundle, fmt.Errorf("handoff: outer manifest requires a Ready-For-Certificate readiness report, got %q", result)
	}
	readinessFingerprint, err := artifact.PayloadFingerprint(ContentReadinessArtifactType, opts.ReadinessArtifactPath, readinessBytes)
	if err != nil {
		return bundle, err
	}
	certificateFingerprint, err := artifact.PayloadFingerprint(ScopeCertificateArtifactType, opts.CertificateArtifactPath, certificateBytes)
	if err != nil {
		return bundle, err
	}

	members := []artifact.Member{
		{Path: opts.ReportArtifactPath, Bytes: reportBytes},
		{Path: opts.ManifestArtifactPath, Bytes: manifestBytes},
		{Path: opts.ReadinessArtifactPath, Bytes: readinessBytes},
		{Path: opts.CertificateArtifactPath, Bytes: certificateBytes},
	}
	if len(opts.StoreDump) > 0 {
		members = append(members, artifact.Member{Path: StoreMemberPath, Bytes: opts.StoreDump})
	}
	if opts.HandbookViewDir != "" {
		paths, err := directoryMembers(opts.HandbookViewDir, "handbook")
		if err != nil {
			return bundle, err
		}
		members = append(members, paths...)
	}
	if opts.HandbookSectionDir != "" {
		paths, err := directoryMembers(opts.HandbookSectionDir, "handbook/sections")
		if err != nil {
			return bundle, err
		}
		members = append(members, paths...)
	}
	if opts.EquivalenceSuitePath != "" {
		content, err := os.ReadFile(opts.EquivalenceSuitePath)
		if err != nil {
			return bundle, err
		}
		members = append(members, artifact.Member{Path: "equivalence/acceptance-suite.md", Bytes: content})
	}
	if opts.EquivalenceReviewPath != "" {
		content, err := os.ReadFile(opts.EquivalenceReviewPath)
		if err != nil {
			return bundle, err
		}
		members = append(members, artifact.Member{Path: "equivalence/equivalence-review.md", Bytes: content})
	}
	entries, err := artifact.PackageMembers(members)
	if err != nil {
		return bundle, err
	}
	innerMembers := make([]artifact.ManifestEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.Path == opts.ArtifactPath {
			continue
		}
		innerMembers = append(innerMembers, entry)
	}
	bundle.Members = innerMembers
	bundle.Files = members

	finalRegistry := artifact.RegistryFingerprint("EXIT-E-FINAL-CHECKS-v1", artifact.RegistryCheckIDs(doc, "FINAL-CHECK"))
	memberTable := renderMemberTable(innerMembers)
	preSignature, err := preSignatureFromTable(doc, snapshot, memberTable, readinessFingerprint, certificateFingerprint)
	if err != nil {
		return bundle, err
	}

	var b strings.Builder
	b.WriteString("<!-- ARTIFACT-PAYLOAD: " + OuterManifestArtifactType + " -->\n\n")
	fmt.Fprintf(&b, "- **ARTIFACT-TYPE:** %s\n", OuterManifestArtifactType)
	fmt.Fprintf(&b, "- **ARTIFACT-PATH:** %s\n", opts.ArtifactPath)
	fmt.Fprintf(&b, "- **PROTOCOL / ARTIFACT-SCHEMA VERSION:** %s / %s\n", strings.TrimSpace(opts.ProtocolVersion), artifact.ArtifactSchemaVersion)
	fmt.Fprintf(&b, "- **System / Protocol / Snapshot:** %s\n", snapshot)
	fmt.Fprintf(&b, "- **Final Verification Check Registry / Fingerprint:** EXIT-E-FINAL-CHECKS-v1 %s\n", finalRegistry)
	fmt.Fprintf(&b, "- **Pre-Signature Validation Input-Set Fingerprint:** %s\n", preSignature)
	fmt.Fprintf(&b, "- **Exit E Content-Readiness Report Fingerprint:** %s\n", readinessFingerprint)
	fmt.Fprintf(&b, "- **Scope Certificate Fingerprint:** %s\n", certificateFingerprint)
	b.WriteString("- **Package Members:**\n")
	b.WriteString(memberTable)
	b.WriteString("- **EXIT-E-STATUS:** Passed\n")
	b.WriteString("- **ARTIFACT-PAYLOAD-FINGERPRINT:** PENDING\n\n")
	b.WriteString("<!-- END-ARTIFACT-PAYLOAD: " + OuterManifestArtifactType + " -->\n\n")
	b.WriteString("<!-- ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: " + OuterManifestArtifactType + " -->\n\n")
	b.WriteString("- **Required Human Signatories / Signature Timestamps:** None\n")
	b.WriteString("- **Envelope Audit:** pending human signature\n")
	b.WriteString("<!-- END-ARTIFACT-DIGEST-SIGNATURE-ENVELOPE: " + OuterManifestArtifactType + " -->\n")
	payload := []byte(b.String())
	if err := artifact.ValidatePayload(doc, OuterManifestArtifactType, payload); err != nil {
		return bundle, err
	}
	fingerprint, err := artifact.PayloadFingerprint(OuterManifestArtifactType, opts.ArtifactPath, payload)
	if err != nil {
		return bundle, err
	}
	final := []byte(strings.Replace(string(payload),
		"- **ARTIFACT-PAYLOAD-FINGERPRINT:** PENDING", "- **ARTIFACT-PAYLOAD-FINGERPRINT:** "+fingerprint, 1))
	recomputed, err := artifact.PayloadFingerprint(OuterManifestArtifactType, opts.ArtifactPath, final)
	if err != nil {
		return bundle, err
	}
	if recomputed != fingerprint {
		return bundle, fmt.Errorf("handoff: outer manifest fingerprint does not recompute")
	}
	instance, err := artifact.InstanceBinding(OuterManifestArtifactType, opts.ArtifactPath, fingerprint)
	if err != nil {
		return bundle, err
	}
	bundle.Payload = final
	bundle.Fingerprint = fingerprint
	bundle.InstanceBinding = instance
	return bundle, nil
}

// preSignatureFromTable computes the §15.4 pre-signature input-set fingerprint
// from the exact rendered member table.
func preSignatureFromTable(doc *cdl.EIRDoc, snapshot, memberTable, readinessFingerprint, certificateFingerprint string) (string, error) {
	var preimage strings.Builder
	preimage.WriteString("EXIT-E-PRE-SIGNATURE-INPUTS|")
	preimage.WriteString(snapshot)
	preimage.WriteString("\n")
	preimage.WriteString("EXIT-E-CHECKS-v1 " + artifact.CheckRegistryFingerprint(doc))
	preimage.WriteString("\n")
	preimage.WriteString("EXIT-E-FINAL-CHECKS-v1 " + artifact.RegistryFingerprint("EXIT-E-FINAL-CHECKS-v1", artifact.RegistryCheckIDs(doc, "FINAL-CHECK")))
	preimage.WriteString("\n")
	preimage.WriteString(readinessFingerprint)
	preimage.WriteString("\n")
	preimage.WriteString(certificateFingerprint)
	preimage.WriteString("\n")
	canonical, err := capabilities.CanonicalMarkdown([]byte(memberTable))
	if err != nil {
		return "", err
	}
	preimage.WriteString(strings.TrimSuffix(string(canonical), "\n"))
	return capabilities.Hash([]byte(preimage.String())), nil
}

// renderMemberTable serializes the sorted member table with its header.
func renderMemberTable(members []artifact.ManifestEntry) string {
	var b strings.Builder
	b.WriteString("  | Artifact Type | Normalized Path | Package-Member File Fingerprint |\n")
	b.WriteString("  | --- | --- | --- |\n")
	for _, member := range members {
		fmt.Fprintf(&b, "  | %s | %s | %s |\n", memberType(member.Path), member.Path, member.Fingerprint)
	}
	return b.String()
}

// memberType classifies one package member path for the manifest table.
func memberType(path string) string {
	switch {
	case strings.HasPrefix(path, "exit-e/"):
		base := strings.TrimSuffix(filepath.Base(path), ".md")
		return strings.ToUpper(base)
	case path == StoreMemberPath:
		return "CANONICAL-STORE-DUMP"
	case strings.HasPrefix(path, "handbook/sections/"):
		return "HANDBOOK-SECTION"
	case strings.HasPrefix(path, "handbook/readability-"):
		return "HANDBOOK-READABILITY-REVIEW"
	case strings.HasPrefix(path, "handbook/"):
		return "HANDBOOK-VIEW"
	case strings.HasPrefix(path, "equivalence/acceptance"):
		return "EQUIVALENCE-SUITE"
	case strings.HasPrefix(path, "equivalence/"):
		return "EQUIVALENCE-REVIEW"
	default:
		return "HUMAN-ARTIFACT"
	}
}

// directoryMembers reads every Markdown file of a directory as package members
// under the given prefix.
func directoryMembers(dir, prefix string) ([]artifact.Member, error) {
	entries, err := os.ReadDir(dir)
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
	members := make([]artifact.Member, 0, len(names))
	for _, name := range names {
		content, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		members = append(members, artifact.Member{Path: prefix + "/" + name, Bytes: content})
	}
	return members, nil
}
