package handoff

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/artifact"
)

func TestBuildCandidatePayloadManifestEmptyState(t *testing.T) {
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
	reportPath := filepath.Join(t.TempDir(), "report.md")
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
	if manifest.SemanticBindings != 0 {
		t.Fatalf("semantic bindings = %d, want 0 on empty state", manifest.SemanticBindings)
	}
	if !strings.Contains(string(manifest.Payload), report.Fingerprint) {
		t.Fatalf("manifest does not reference the candidate report fingerprint")
	}
	recomputed, err := artifact.PayloadFingerprint(CandidateManifestArtifactType, CandidateManifestDefaultPath, manifest.Payload)
	if err != nil {
		t.Fatal(err)
	}
	if recomputed != manifest.Fingerprint {
		t.Fatalf("fingerprint %s does not recompute to %s", manifest.Fingerprint, recomputed)
	}
	if err := artifact.ValidatePayload(doc, CandidateManifestArtifactType, manifest.Payload); err != nil {
		t.Fatalf("payload schema: %v", err)
	}
}
