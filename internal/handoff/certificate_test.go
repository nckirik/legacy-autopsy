package handoff

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/artifact"
)

func TestBuildScopeCertificateRefusesBlockedReadiness(t *testing.T) {
	doc, db := testDocDB(t)
	dir := t.TempDir()
	reportPath := filepath.Join(dir, "report.md")
	manifestPath := filepath.Join(dir, "manifest.md")
	readinessPath := filepath.Join(dir, "readiness.md")
	if err := os.WriteFile(reportPath, []byte("- **Exit A Report Fingerprint:** sha256:exit-a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifestPath, []byte("- **Exit E Candidate Report Fingerprint:** sha256:report\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	blocked := "- **Exit E Candidate Report Fingerprint:** sha256:report\n- **Candidate Payload Manifest Fingerprint:** sha256:manifest\n- **EXIT-E-PREPACKAGE-STATE:** Pending\n- **Content Readiness Result:** Blocked\n"
	if err := os.WriteFile(readinessPath, []byte(blocked), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := BuildScopeCertificate(db, doc, CertificateOptions{
		Namespace:       "demo",
		ProtocolVersion: "4.6",
		ReadinessPath:   readinessPath,
		ReportPath:      reportPath,
		ManifestPath:    manifestPath,
	}); err == nil {
		t.Fatalf("certificate accepted a blocked readiness report")
	}
}

func TestBuildScopeCertificateReady(t *testing.T) {
	doc, db := testDocDB(t)
	dir := t.TempDir()
	reportPath := filepath.Join(dir, "report.md")
	manifestPath := filepath.Join(dir, "manifest.md")
	readinessPath := filepath.Join(dir, "readiness.md")
	reportBytes := []byte("- **Exit A Report Fingerprint:** sha256:exit-a\n")
	manifestBytes := []byte("- **Exit E Candidate Report Fingerprint:** sha256:report\n")
	if err := os.WriteFile(reportPath, reportBytes, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifestPath, manifestBytes, 0o644); err != nil {
		t.Fatal(err)
	}
	reportFingerprint, err := artifact.PayloadFingerprint(CandidateReportArtifactType, CandidateReportDefaultPath, reportBytes)
	if err != nil {
		t.Fatal(err)
	}
	manifestFingerprint, err := artifact.PayloadFingerprint(CandidateManifestArtifactType, CandidateManifestDefaultPath, manifestBytes)
	if err != nil {
		t.Fatal(err)
	}
	ready := "- **Exit E Candidate Report Fingerprint:** " + reportFingerprint + "\n" +
		"- **Candidate Payload Manifest Fingerprint:** " + manifestFingerprint + "\n" +
		"- **EXIT-E-PREPACKAGE-STATE:** Pending\n" +
		"- **Content Readiness Result:** Ready-For-Certificate\n"
	if err := os.WriteFile(readinessPath, []byte(ready), 0o644); err != nil {
		t.Fatal(err)
	}
	certificate, err := BuildScopeCertificate(db, doc, CertificateOptions{
		Namespace:             "demo",
		ProtocolVersion:       "4.6",
		EnvironmentsSnapshots: "dev",
		StorageBinding:        "sqlite/1",
		ReportPath:            reportPath,
		ManifestPath:          manifestPath,
		ReadinessPath:         readinessPath,
	})
	if err != nil {
		t.Fatal(err)
	}
	recomputed, err := artifact.PayloadFingerprint(ScopeCertificateArtifactType, ScopeCertificateDefaultPath, certificate.Payload)
	if err != nil {
		t.Fatal(err)
	}
	if recomputed != certificate.Fingerprint {
		t.Fatalf("fingerprint %s does not recompute to %s", certificate.Fingerprint, recomputed)
	}
	if err := artifact.ValidatePayload(doc, ScopeCertificateArtifactType, certificate.Payload); err != nil {
		t.Fatalf("payload schema: %v", err)
	}
}
