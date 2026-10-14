package handoff

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/artifact"
)

func TestBuildContentReadinessEmptyState(t *testing.T) {
	doc, db := testDocDB(t)
	report, err := BuildCandidateReport(db, doc, ReportOptions{
		Namespace:             "demo",
		ProtocolVersion:       "4.6",
		EnvironmentsSnapshots: "dev",
		ArtifactPath:          "exit-e/EXIT-E-CANDIDATE-REPORT.md",
	})
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	reportPath := filepath.Join(dir, "report.md")
	if err := os.WriteFile(reportPath, report.Payload, 0o644); err != nil {
		t.Fatal(err)
	}
	manifest, err := BuildCandidatePayloadManifest(db, doc, ManifestOptions{
		ProtocolVersion: "4.6",
		ReportPath:      reportPath,
	})
	if err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(dir, "manifest.md")
	if err := os.WriteFile(manifestPath, manifest.Payload, 0o644); err != nil {
		t.Fatal(err)
	}
	readiness, err := BuildContentReadiness(db, doc, ReadinessOptions{
		Namespace:             "demo",
		ProtocolVersion:       "4.6",
		EnvironmentsSnapshots: "dev",
		ReportPath:            reportPath,
		ManifestPath:          manifestPath,
	})
	if err != nil {
		t.Fatal(err)
	}
	if readiness.Result != "Blocked" {
		t.Fatalf("result = %q, want Blocked on empty state", readiness.Result)
	}
	if len(readiness.Checks) != 11 {
		t.Fatalf("checks = %d, want 11", len(readiness.Checks))
	}
	recomputed, err := artifact.PayloadFingerprint(ContentReadinessArtifactType, ContentReadinessDefaultPath, readiness.Payload)
	if err != nil {
		t.Fatal(err)
	}
	if recomputed != readiness.Fingerprint {
		t.Fatalf("fingerprint %s does not recompute to %s", readiness.Fingerprint, recomputed)
	}
	if err := artifact.ValidatePayload(doc, ContentReadinessArtifactType, readiness.Payload); err != nil {
		t.Fatalf("payload schema: %v", err)
	}
}
