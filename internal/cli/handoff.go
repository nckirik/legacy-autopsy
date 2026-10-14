package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/handoff"
)

// specHandoff dispatches the §15.1/§15.4 handoff operations.
func specHandoff(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: legacy-autopsy spec handoff <candidate|report|manifest> [flags]")
	}
	switch args[0] {
	case "candidate":
		return specHandoffCandidate(args[1:], stdout, stderr)
	case "report":
		return specHandoffReport(args[1:], stdout, stderr)
	case "manifest":
		return specHandoffManifest(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("spec handoff operation %q is unsupported; supported operations: candidate, report, manifest", args[0])
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
	handbookDir := set.String("handbook-dir", "", "minted handbook directory to bind HBK sections")
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
