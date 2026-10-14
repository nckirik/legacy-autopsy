package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/equivalence"
	"github.com/nckirik/legacy-autopsy/internal/handoff"
	"github.com/nckirik/legacy-autopsy/internal/store"
)

// specHandoff dispatches the §15.1/§15.4 handoff operations.
func specHandoff(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: legacy-autopsy spec handoff <candidate|report|manifest|readiness|equivalence|review|certificate|outer> [flags]")
	}
	switch args[0] {
	case "candidate":
		return specHandoffCandidate(args[1:], stdout, stderr)
	case "report":
		return specHandoffReport(args[1:], stdout, stderr)
	case "manifest":
		return specHandoffManifest(args[1:], stdout, stderr)
	case "readiness":
		return specHandoffReadiness(args[1:], stdout, stderr)
	case "equivalence":
		return specHandoffEquivalence(args[1:], stdout, stderr)
	case "review":
		return specHandoffReview(args[1:], stdout, stderr)
	case "certificate":
		return specHandoffCertificate(args[1:], stdout, stderr)
	case "outer":
		return specHandoffOuter(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("spec handoff operation %q is unsupported; supported operations: candidate, report, manifest, readiness, equivalence, review, certificate, outer", args[0])
	}
}

// specHandoffCandidate prints the §15.1 candidate condition report.
func specHandoffCandidate(args []string, stdout, stderr io.Writer) error {
	set := flags("spec handoff candidate", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	handbookDir := set.String("handbook-dir", "", "rendered handbook directory for the quality check")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" {
		return fmt.Errorf("spec handoff candidate requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	checks, err := handoff.EvaluateCandidateChecks(*db, res.EIR, handoff.Options{HandbookDir: *handbookDir})
	if err != nil {
		return err
	}
	failures := 0
	unsupported := 0
	for _, check := range checks {
		fmt.Fprintf(stdout, "%-40s %-12s %s\n", check.ID, check.Result, check.Detail)
		switch check.Result {
		case "Fail":
			failures++
		case "Unsupported":
			unsupported++
		}
	}
	fmt.Fprintf(stdout, "candidate: %d pass, %d fail, %d unsupported\n", len(checks)-failures-unsupported, failures, unsupported)
	if failures > 0 || unsupported > 0 {
		return fmt.Errorf("candidate conditions are not ready")
	}
	return nil
}

// specHandoffReport emits the immutable §15.4 step-2 Exit E Candidate Report.
func specHandoffReport(args []string, stdout, stderr io.Writer) error {
	set := flags("spec handoff report", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	handbookDir := set.String("handbook-dir", "", "rendered handbook directory for the quality check")
	out := set.String("out", "preview/package/EXIT-E-CANDIDATE-REPORT.md", "output file for the artifact payload")
	artifactPath := set.String("artifact-path", handoff.CandidateReportDefaultPath, "package-relative artifact path")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" {
		return fmt.Errorf("spec handoff report requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	header, err := preflightHeaderValues(*db, res.EIR)
	if err != nil {
		return err
	}
	report, err := handoff.BuildCandidateReport(*db, res.EIR, handoff.ReportOptions{
		Namespace:             header["system-namespace"],
		ProtocolVersion:       header["protocol-version"],
		EnvironmentsSnapshots: header["included-environments-snapshots"],
		StorageBinding:        header["canonical-storage-binding"],
		HandbookDir:           *handbookDir,
		ArtifactPath:          *artifactPath,
	})
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(*out, report.Payload, 0o644); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "artifact: %s\n", *artifactPath)
	fmt.Fprintf(stdout, "fingerprint: %s\n", report.Fingerprint)
	fmt.Fprintf(stdout, "instance: %s\n", report.InstanceBinding)
	fmt.Fprintf(stdout, "candidate result: %s\n", report.Result)
	fmt.Fprintf(stdout, "written: %s\n", *out)
	if report.Result != "Ready-For-Human-Review" {
		return fmt.Errorf("candidate report is Blocked")
	}
	return nil
}

// specHandoffManifest emits the immutable §15.4 step-2 candidate payload manifest.
func specHandoffManifest(args []string, stdout, stderr io.Writer) error {
	set := flags("spec handoff manifest", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	report := set.String("report", "preview/package/EXIT-E-CANDIDATE-REPORT.md", "emitted candidate report payload file")
	reportPath := set.String("report-artifact-path", handoff.CandidateReportDefaultPath, "candidate report package-relative artifact path")
	handbookDir := set.String("handbook-dir", "preview/package/handbook/sections", "minted HBK section payload directory")
	handbookPrefix := set.String("handbook-prefix", "handbook/sections", "package-relative prefix of HBK section payloads")
	out := set.String("out", "preview/package/CANDIDATE-PAYLOAD-MANIFEST.md", "output file for the artifact payload")
	artifactPath := set.String("artifact-path", handoff.CandidateManifestDefaultPath, "package-relative artifact path")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" {
		return fmt.Errorf("spec handoff manifest requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	header, err := preflightHeaderValues(*db, res.EIR)
	if err != nil {
		return err
	}
	manifest, err := handoff.BuildCandidatePayloadManifest(*db, res.EIR, handoff.ManifestOptions{
		ProtocolVersion:    header["protocol-version"],
		ReportPath:         *report,
		ReportArtifactPath: *reportPath,
		HandbookDir:        *handbookDir,
		HandbookPrefix:     *handbookPrefix,
		ArtifactPath:       *artifactPath,
	})
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(*out, manifest.Payload, 0o644); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "artifact: %s\n", *artifactPath)
	fmt.Fprintf(stdout, "fingerprint: %s\n", manifest.Fingerprint)
	fmt.Fprintf(stdout, "instance: %s\n", manifest.InstanceBinding)
	fmt.Fprintf(stdout, "semantic bindings: %d, DEC bindings: %d\n", manifest.SemanticBindings, manifest.DECBindings)
	fmt.Fprintf(stdout, "written: %s\n", *out)
	return nil
}

// specHandoffReadiness emits the immutable §15.4 step-4 content-readiness report.
func specHandoffReadiness(args []string, stdout, stderr io.Writer) error {
	set := flags("spec handoff readiness", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	report := set.String("report", "preview/package/EXIT-E-CANDIDATE-REPORT.md", "emitted candidate report payload file")
	reportPath := set.String("report-artifact-path", handoff.CandidateReportDefaultPath, "candidate report package-relative artifact path")
	manifest := set.String("manifest", "preview/package/CANDIDATE-PAYLOAD-MANIFEST.md", "emitted candidate payload manifest file")
	manifestPath := set.String("manifest-artifact-path", handoff.CandidateManifestDefaultPath, "manifest package-relative artifact path")
	handbookDir := set.String("handbook-dir", "preview/package/handbook/sections", "minted HBK section payload directory")
	handbookPrefix := set.String("handbook-prefix", "handbook/sections", "package-relative prefix of HBK section payloads")
	handbookViewDir := set.String("handbook-view-dir", "preview", "clean human handbook view directory")
	equivalencePath := set.String("equivalence", "", "generated acceptance-suite file")
	out := set.String("out", "preview/package/EXIT-E-CONTENT-READINESS-REPORT.md", "output file for the artifact payload")
	artifactPath := set.String("artifact-path", handoff.ContentReadinessDefaultPath, "package-relative artifact path")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" {
		return fmt.Errorf("spec handoff readiness requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	header, err := preflightHeaderValues(*db, res.EIR)
	if err != nil {
		return err
	}
	dump, _, err := store.CanonicalDump(*db, res.EIR, store.Provenance{
		Generator:         "legacy-autopsy content readiness",
		SourceFingerprint: header["system-namespace"],
	})
	if err != nil {
		return err
	}
	readiness, err := handoff.BuildContentReadiness(*db, res.EIR, handoff.ReadinessOptions{
		Namespace:             header["system-namespace"],
		ProtocolVersion:       header["protocol-version"],
		EnvironmentsSnapshots: header["included-environments-snapshots"],
		StorageBinding:        header["canonical-storage-binding"],
		ReportPath:            *report,
		ReportArtifactPath:    *reportPath,
		ManifestPath:          *manifest,
		ManifestArtifactPath:  *manifestPath,
		HandbookDir:           *handbookDir,
		HandbookPrefix:        *handbookPrefix,
		HandbookViewDir:       *handbookViewDir,
		EquivalencePath:       *equivalencePath,
		StoreDump:             dump,
		ArtifactPath:          *artifactPath,
	})
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(*out, readiness.Payload, 0o644); err != nil {
		return err
	}
	failures := 0
	for _, check := range readiness.Checks {
		fmt.Fprintf(stdout, "%-42s %s\n", check.ID, check.Result)
		if check.Result != "Pass" {
			failures++
		}
	}
	fmt.Fprintf(stdout, "artifact: %s\n", *artifactPath)
	fmt.Fprintf(stdout, "fingerprint: %s\n", readiness.Fingerprint)
	fmt.Fprintf(stdout, "instance: %s\n", readiness.InstanceBinding)
	fmt.Fprintf(stdout, "content readiness: %s (%d failing checks)\n", readiness.Result, failures)
	fmt.Fprintf(stdout, "written: %s\n", *out)
	return nil
}

// specHandoffEquivalence generates the deterministic acceptance-test suite.
func specHandoffEquivalence(args []string, stdout, stderr io.Writer) error {
	set := flags("spec handoff equivalence", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	out := set.String("out", "preview/package/equivalence/acceptance-suite.md", "output file for the suite")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" {
		return fmt.Errorf("spec handoff equivalence requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	suite, err := equivalence.Build(*db, res.EIR)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(*out, suite.Markdown, 0o644); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "suite mappings: %d\n", len(suite.Mappings))
	fmt.Fprintf(stdout, "written: %s\n", *out)
	return nil
}

// specHandoffCertificate emits the immutable §15.4 step-5 scope certificate.
func specHandoffCertificate(args []string, stdout, stderr io.Writer) error {
	set := flags("spec handoff certificate", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	report := set.String("report", "preview/package/EXIT-E-CANDIDATE-REPORT.md", "emitted candidate report payload file")
	reportPath := set.String("report-artifact-path", handoff.CandidateReportDefaultPath, "candidate report package-relative artifact path")
	manifest := set.String("manifest", "preview/package/CANDIDATE-PAYLOAD-MANIFEST.md", "emitted candidate payload manifest file")
	manifestPath := set.String("manifest-artifact-path", handoff.CandidateManifestDefaultPath, "manifest package-relative artifact path")
	readiness := set.String("readiness", "preview/package/EXIT-E-CONTENT-READINESS-REPORT.md", "emitted content-readiness report file")
	readinessPath := set.String("readiness-artifact-path", handoff.ContentReadinessDefaultPath, "readiness package-relative artifact path")
	out := set.String("out", "preview/package/exit-e/SCOPE-CERTIFICATE.md", "output file for the artifact payload")
	artifactPath := set.String("artifact-path", handoff.ScopeCertificateDefaultPath, "package-relative artifact path")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" {
		return fmt.Errorf("spec handoff certificate requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	header, err := preflightHeaderValues(*db, res.EIR)
	if err != nil {
		return err
	}
	certificate, err := handoff.BuildScopeCertificate(*db, res.EIR, handoff.CertificateOptions{
		Namespace:             header["system-namespace"],
		ProtocolVersion:       header["protocol-version"],
		EnvironmentsSnapshots: header["included-environments-snapshots"],
		StorageBinding:        header["canonical-storage-binding"],
		ReportPath:            *report,
		ReportArtifactPath:    *reportPath,
		ManifestPath:          *manifest,
		ManifestArtifactPath:  *manifestPath,
		ReadinessPath:         *readiness,
		ReadinessArtifactPath: *readinessPath,
		ArtifactPath:          *artifactPath,
	})
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(*out, certificate.Payload, 0o644); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "artifact: %s\n", *artifactPath)
	fmt.Fprintf(stdout, "fingerprint: %s\n", certificate.Fingerprint)
	fmt.Fprintf(stdout, "instance: %s\n", certificate.InstanceBinding)
	fmt.Fprintf(stdout, "written: %s\n", *out)
	return nil
}

// specHandoffOuter emits the step-6 outer bundle manifest.
func specHandoffOuter(args []string, stdout, stderr io.Writer) error {
	set := flags("spec handoff outer", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	report := set.String("report", "preview/package/EXIT-E-CANDIDATE-REPORT.md", "candidate report payload file")
	manifest := set.String("manifest", "preview/package/CANDIDATE-PAYLOAD-MANIFEST.md", "candidate payload manifest file")
	readiness := set.String("readiness", "preview/package/EXIT-E-CONTENT-READINESS-REPORT.md", "content-readiness report file")
	certificate := set.String("certificate", "preview/package/exit-e/SCOPE-CERTIFICATE.md", "scope certificate file")
	handbookViewDir := set.String("handbook-view-dir", "preview", "clean human handbook view directory")
	handbookDir := set.String("handbook-dir", "preview/package/handbook/sections", "minted HBK section payload directory")
	suite := set.String("equivalence", "preview/package/equivalence/acceptance-suite.md", "acceptance-suite file")
	suiteReview := set.String("equivalence-review", "preview/package/equivalence/equivalence-review.md", "acceptance-suite review file")
	out := set.String("out", "preview/package/exit-e/OUTER-BUNDLE-MANIFEST.md", "output file for the artifact payload")
	dumpOut := set.String("dump-out", "", "optional path to write the canonical store dump package member")
	packageDir := set.String("package-dir", "", "optional directory to materialize the complete package member layout")
	artifactPath := set.String("artifact-path", handoff.OuterManifestDefaultPath, "package-relative artifact path")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" {
		return fmt.Errorf("spec handoff outer requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	header, err := preflightHeaderValues(*db, res.EIR)
	if err != nil {
		return err
	}
	dump, _, err := store.CanonicalDump(*db, res.EIR, store.Provenance{
		Generator:         "legacy-autopsy outer manifest",
		SourceFingerprint: header["system-namespace"],
	})
	if err != nil {
		return err
	}
	bundle, err := handoff.BuildOuterBundle(res.EIR, handoff.OuterOptions{
		Namespace:             header["system-namespace"],
		ProtocolVersion:       header["protocol-version"],
		EnvironmentsSnapshots: header["included-environments-snapshots"],
		StorageBinding:        header["canonical-storage-binding"],
		ReportPath:            *report,
		ManifestPath:          *manifest,
		ReadinessPath:         *readiness,
		CertificatePath:       *certificate,
		HandbookViewDir:       *handbookViewDir,
		HandbookSectionDir:    *handbookDir,
		EquivalenceSuitePath:  *suite,
		EquivalenceReviewPath: *suiteReview,
		StoreDump:             dump,
		ArtifactPath:          *artifactPath,
	})
	if err != nil {
		return err
	}
	if *dumpOut != "" {
		if err := os.MkdirAll(filepath.Dir(*dumpOut), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(*dumpOut, dump, 0o644); err != nil {
			return err
		}
	}
	if *packageDir != "" {
		for _, member := range bundle.Files {
			path := filepath.Join(*packageDir, filepath.FromSlash(member.Path))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return err
			}
			if err := os.WriteFile(path, member.Bytes, 0o644); err != nil {
				return err
			}
		}
		fmt.Fprintf(stdout, "package members written: %d\n", len(bundle.Files))
	}
	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(*out, bundle.Payload, 0o644); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "artifact: %s\n", *artifactPath)
	fmt.Fprintf(stdout, "fingerprint: %s\n", bundle.Fingerprint)
	fmt.Fprintf(stdout, "instance: %s\n", bundle.InstanceBinding)
	fmt.Fprintf(stdout, "members: %d\n", len(bundle.Members))
	fmt.Fprintf(stdout, "written: %s\n", *out)
	return nil
}
