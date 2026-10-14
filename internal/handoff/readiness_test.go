package handoff

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/artifact"
	"github.com/nckirik/legacy-autopsy/internal/store"
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

func TestDecisionApprovalCoverageRequiresCurrentBinding(t *testing.T) {
	decisionPayload := func(envelopeVersion string, binding string) []byte {
		return []byte("- **decision-content-version:** 2\n- **decision-content-fingerprint:** sha256:new\n- **type:** target\n- **observed-inputs:** ARCHITECTURE-MODULE/MOD-AAAAAAAAAAAA version 1 sha256:aaa\n- **options-decision-rationale-risks:** keep\n- **exact-scope-and-denominator-effect:** none\n- **affected-records-rollout-rollback:** none\n- **content-supersession:** None\n\n<!-- DECISION-APPROVAL-ENVELOPE: excluded from decision-content version/fingerprint -->\n\n- **approval-status:** Approved\n- **approval-envelope-version:** " + envelopeVersion + "\n- **owner-approver-authority-basis-approval-timestamp:** operator 2026-10-14\n- **approved-content-binding:** " + binding + "\n- **approval-signature-audit:** decision review\n<!-- END-DECISION-APPROVAL-ENVELOPE -->\n")
	}
	seedModule := func(t *testing.T, doc interface{}, db string) {
		t.Helper()
		module := []byte("- **source-bindings:** CLM-AAAAAAAAAAAA\n- **block-id:** MOD-AAAAAAAAAAAA\n- **record-title:** Request validation\n- **domain-path:** demo/core\n- **responsibility-domain-boundary:** Validates requests.\n- **public-interfaces:** None\n- **encapsulated-state-entities:** None\n- **dependencies-dependents:** None\n- **data-and-events-passed:** None\n- **invariants-rules-owned:** Unknown\n- **personas-use-cases-capabilities:** None\n- **failure-propagation:** Unknown\n- **observed-legacy-placement:** inline\n- **target-decision:** DEC-AAAAAAAAAAAA\n")
		if err := store.Ingest(db, doc.(*cdl.EIRDoc), "ARCHITECTURE-MODULE", "MOD-AAAAAAAAAAAA", 1, module); err != nil {
			t.Fatal(err)
		}
	}

	doc, db := testDocDB(t)
	seedModule(t, doc, db)
	if err := store.Ingest(db, doc, "DECISION", "DEC-AAAAAAAAAAAA", 1, decisionPayload("1", "DEC-AAAAAAAAAAAA version 1 sha256:old")); err != nil {
		t.Fatal(err)
	}
	if _, fail := decisionApprovalCoverage(db, doc); !strings.Contains(fail, "stale decision content") {
		t.Fatalf("stale approval binding was not rejected: %q", fail)
	}

	doc, db = testDocDB(t)
	seedModule(t, doc, db)
	if err := store.Ingest(db, doc, "DECISION", "DEC-AAAAAAAAAAAA", 1, decisionPayload("1", "DEC-AAAAAAAAAAAA version 2 sha256:new")); err != nil {
		t.Fatal(err)
	}
	detail, fail := decisionApprovalCoverage(db, doc)
	if fail != "" {
		t.Fatalf("current approval binding rejected: %s", fail)
	}
	if !strings.Contains(detail, "1 decision references") {
		t.Fatalf("detail = %q", detail)
	}

	doc, db = testDocDB(t)
	seedModule(t, doc, db)
	if err := store.Ingest(db, doc, "DECISION", "DEC-AAAAAAAAAAAA", 1, decisionPayload("0", "DEC-AAAAAAAAAAAA version 2 sha256:new")); err != nil {
		t.Fatal(err)
	}
	if _, fail := decisionApprovalCoverage(db, doc); !strings.Contains(fail, "approval envelope missing") {
		t.Fatalf("missing envelope was not rejected: %q", fail)
	}
}
