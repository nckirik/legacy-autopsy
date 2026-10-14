package handoff

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/artifact"
	"github.com/nckirik/legacy-autopsy/internal/store"
)

func TestProfileDependenciesParseFingerprintAnnotations(t *testing.T) {
	doc, db := testDocDB(t)
	registry := []byte("- **persona:** Admin\n- **prefix:** AD\n- **persona-slug:** admin\n- **persona-directory:** AD-admin\n- **actor-class:** human\n- **canonical-entry-ids:** None\n- **scope-paths:** app/\n- **registry-state:** Active\n")
	if err := store.Ingest(db, doc, "PERSONA-REGISTRY", "AD", 1, registry); err != nil {
		t.Fatal(err)
	}
	fingerprint, err := store.FingerprintOf(db, doc, "PERSONA-REGISTRY", "AD")
	if err != nil {
		t.Fatal(err)
	}
	values := map[string]string{"derived-from": "PERSONA-REGISTRY/AD@" + fingerprint, "material-claims": "None"}
	tuples, err := profileDependencies(db, doc, values)
	if err != nil {
		t.Fatal(err)
	}
	if len(tuples) != 1 || tuples[0][0] != "AD" || !strings.Contains(tuples[0][1], fingerprint) || tuples[0][2] != "REFERENCE" {
		t.Fatalf("tuples = %+v", tuples)
	}
	values["derived-from"] = "PERSONA-REGISTRY/AD@sha256:0000000000000000000000000000000000000000000000000000000000000000"
	if _, err := profileDependencies(db, doc, values); err == nil || !strings.Contains(err.Error(), "declares fingerprint") {
		t.Fatalf("mismatched declared fingerprint accepted: %v", err)
	}
	values["derived-from"] = "PERSONA-REGISTRY/AD@version 1"
	if _, err := profileDependencies(db, doc, values); err == nil || !strings.Contains(err.Error(), "declares no fingerprint") {
		t.Fatalf("annotation without a fingerprint accepted: %v", err)
	}
}

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
