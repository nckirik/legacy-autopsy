package assurance

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
)

func seedCoverageRow(t *testing.T, db string, doc *cdl.EIRDoc, id, src, status string) {
	t.Helper()
	payload := "- **cov-id:** " + id + "\n- **src-version:** " + src + " version 1\n- **persona-owner:** WEB\n- **track:** Normal\n- **environment:** dev\n- **snapshot:** dev\n- **scope:** In-Scope\n- **disposition:** Mapped\n- **c-status:** " + status + "\n- **exclusion-dec:** None\n- **mappings:** None\n- **frontier-ticket-claim-evidence:** None\n"
	if err := store.Ingest(db, doc, "SOURCE-COVERAGE-ROW", id, 1, []byte(payload)); err != nil {
		t.Fatal(err)
	}
}

func seedEvidenceClaim(t *testing.T, db string, doc *cdl.EIRDoc, id, subject, materiality, level string) {
	t.Helper()
	payload := "- **subject-id-record-version:** " + subject + " version 1\n- **field-step:** lock evidence\n- **claim:** The component behavior is observed.\n- **materiality:** " + materiality + "\n- **evidence-level:** " + level + "\n- **anchor:** app/a.py:1\n- **source-fingerprint-snapshot:** sha256:aaa\n- **evidence-route-lineage-id:** static-filesystem\n- **projection-classification:** raw-visible\n- **independence-basis:** single static route\n- **created-last-revalidated-by:** ALFA-NONE-NONE-001\n"
	if err := store.Ingest(db, doc, "CLAIM", id, 1, []byte(payload)); err != nil {
		t.Fatal(err)
	}
}

func promotionComponent(personas, staging, source, claims, depends string) []byte {
	return []byte("- **record-kind:** ATOMIC\n" +
		"- **source-ids:** " + source + "\n" +
		"- **file-path-line-range:** app/a.py:1\n" +
		"- **virtual-coordinate:** app/a.py :: job :: scan\n" +
		"- **pov-owner:** POV-3\n" +
		"- **general-intent:** Background scan.\n" +
		"- **inputs-outputs:** None\n" +
		"- **auth-binding:** None\n" +
		"- **staging-status:** " + staging + "\n" +
		"- **coverage-status:** [C-COVERED]\n" +
		"- **personas:** " + personas + "\n" +
		"- **traversal-tracks:** Normal\n" +
		"- **capability-state:** Active\n" +
		"- **runtime-state-matrix:** dev\n" +
		"- **depends-on-depended-on-by:** " + depends + "\n" +
		"- **synthesis-references:** None\n" +
		"- **side-effects:** None\n" +
		"- **exceptional-behaviors:** Unknown\n" +
		"- **edge-and-odd-cases:** Unknown\n" +
		"- **behavioral-notes-for-reconstruction:** Unknown\n" +
		"- **material-claims:** " + claims + "\n" +
		"- **block-confidence-rank:** High\n" +
		"- **block-evidence-profile:** [E-DIRECT] app/a.py:1\n" +
		"- **evidence-anchors:** sha256:aaa\n" +
		"- **record-version-supersession:** version 1\n")
}

func TestPromotionLocks(t *testing.T) {
	doc, db := testDocDB(t)
	seedSource(t, db, doc, "SRC-FILE-AAAAAAAAAAAA")
	seedCoverageRow(t, db, doc, "COV-1", "SRC-FILE-AAAAAAAAAAAA", "[C-COVERED]")
	seedEvidenceClaim(t, db, doc, "CLM-AAAAAAAAAAAA", "SRC-FILE-AAAAAAAAAAAA", "Material", "[E-DIRECT]")
	if err := store.Ingest(db, doc, "ATOMIC-COMPONENT", "CMP-AAAAAAAAAAAA", 1, promotionComponent("WEB", "[S-STAGING-VERIFIED]", "SRC-FILE-AAAAAAAAAAAA", "CLM-AAAAAAAAAAAA", "None")); err != nil {
		t.Fatal(err)
	}
	if err := ValidatePromotionLocks(db, doc, "CMP-AAAAAAAAAAAA"); err != nil {
		t.Fatalf("valid locks rejected: %v", err)
	}

	// Usage uniqueness: two personas fail.
	if err := store.ReplaceRecord(db, doc, "ATOMIC-COMPONENT", "CMP-AAAAAAAAAAAA", 2, promotionComponent("WEB, AD", "[S-STAGING-VERIFIED]", "SRC-FILE-AAAAAAAAAAAA", "CLM-AAAAAAAAAAAA", "None")); err != nil {
		t.Fatal(err)
	}
	if err := ValidatePromotionLocks(db, doc, "CMP-AAAAAAAAAAAA"); err == nil || !strings.Contains(err.Error(), "usage uniqueness") {
		t.Fatalf("two personas accepted: %v", err)
	}

	// Evidence: a material [E-UNKNOWN] claim fails.
	seedEvidenceClaim(t, db, doc, "CLM-BBBBBBBBBBBB", "SRC-FILE-AAAAAAAAAAAA", "Material", "[E-UNKNOWN]")
	if err := store.ReplaceRecord(db, doc, "ATOMIC-COMPONENT", "CMP-AAAAAAAAAAAA", 3, promotionComponent("WEB", "[S-STAGING-VERIFIED]", "SRC-FILE-AAAAAAAAAAAA", "CLM-BBBBBBBBBBBB", "None")); err != nil {
		t.Fatal(err)
	}
	if err := ValidatePromotionLocks(db, doc, "CMP-AAAAAAAAAAAA"); err == nil || !strings.Contains(err.Error(), "evidence") {
		t.Fatalf("material unknown claim accepted: %v", err)
	}

	// Coverage: a C-GAP row fails.
	if err := store.ReplaceRecord(db, doc, "SOURCE-COVERAGE-ROW", "COV-1", 2, []byte("- **cov-id:** COV-1\n- **src-version:** SRC-FILE-AAAAAAAAAAAA version 1\n- **persona-owner:** WEB\n- **track:** Normal\n- **environment:** dev\n- **snapshot:** dev\n- **scope:** In-Scope\n- **disposition:** Unmapped\n- **c-status:** [C-GAP]\n- **exclusion-dec:** None\n- **mappings:** None\n- **frontier-ticket-claim-evidence:** None\n")); err != nil {
		t.Fatal(err)
	}
	if err := store.ReplaceRecord(db, doc, "ATOMIC-COMPONENT", "CMP-AAAAAAAAAAAA", 4, promotionComponent("WEB", "[S-STAGING-VERIFIED]", "SRC-FILE-AAAAAAAAAAAA", "CLM-AAAAAAAAAAAA", "None")); err != nil {
		t.Fatal(err)
	}
	if err := ValidatePromotionLocks(db, doc, "CMP-AAAAAAAAAAAA"); err == nil || !strings.Contains(err.Error(), "coverage") {
		t.Fatalf("uncovered source accepted: %v", err)
	}

	// Dependency: another component depending on it fails.
	if err := store.Ingest(db, doc, "ATOMIC-COMPONENT", "CMP-BBBBBBBBBBBB", 1, promotionComponent("AD", "[S-DISCOVERED]", "SRC-FILE-AAAAAAAAAAAA", "None", "CMP-AAAAAAAAAAAA")); err != nil {
		t.Fatal(err)
	}
	if err := store.ReplaceRecord(db, doc, "SOURCE-COVERAGE-ROW", "COV-1", 3, []byte("- **cov-id:** COV-1\n- **src-version:** SRC-FILE-AAAAAAAAAAAA version 1\n- **persona-owner:** WEB\n- **track:** Normal\n- **environment:** dev\n- **snapshot:** dev\n- **scope:** In-Scope\n- **disposition:** Mapped\n- **c-status:** [C-COVERED]\n- **exclusion-dec:** None\n- **mappings:** None\n- **frontier-ticket-claim-evidence:** None\n")); err != nil {
		t.Fatal(err)
	}
	if err := ValidatePromotionLocks(db, doc, "CMP-AAAAAAAAAAAA"); err == nil || !strings.Contains(err.Error(), "dependency") {
		t.Fatalf("dependent component accepted: %v", err)
	}

	// Ambiguity: an active ticket fails.
	ticket := ticketTestPayload("CMP-AAAAAAAAAAAA", "CLM-AAAAAAAAAAAA")
	if err := store.Ingest(db, doc, "TICKET", "TKT-AD-POV-3-001", 1, ticket); err != nil {
		t.Fatal(err)
	}
	if err := ValidatePromotionLocks(db, doc, "CMP-AAAAAAAAAAAA"); err == nil || !strings.Contains(err.Error(), "ambiguity") {
		t.Fatalf("active ticket accepted: %v", err)
	}

	// The diagnostics surface staged components with broken locks.
	diagnostics, err := PromotionLockDiagnostics(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	if len(diagnostics) != 1 || !strings.Contains(diagnostics[0], "CMP-AAAAAAAAAAAA") {
		t.Fatalf("diagnostics = %v", diagnostics)
	}
}
