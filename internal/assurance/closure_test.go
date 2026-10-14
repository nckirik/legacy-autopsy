package assurance

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
)

func seedSource(t *testing.T, db string, doc *cdl.EIRDoc, id string) {
	t.Helper()
	payload := "- **kind:** FILE\n- **concrete-coordinate:** app/" + id + ".py\n- **parent-container-src:** None\n- **source-fingerprint-snapshot-src-version:** sha256:aaa version 1\n- **discovery-method:** filesystem\n- **scope-summary:** demo\n- **environment-runtime-state:** dev\n- **traversal-tracks:** Normal\n- **canonical-owner:** None\n- **disposition-summary:** Unknown\n- **mappings:** None\n- **frontier-references:** None\n- **coverage-rows:** None\n- **coverage-summary:** C-GAP\n- **claims:** None\n"
	if err := store.Ingest(db, doc, "SOURCE-INVENTORY-RECORD", id, 1, []byte(payload)); err != nil {
		t.Fatal(err)
	}
}

func componentPayload(id, status, dependencies string) []byte {
	return []byte("- **record-kind:** ATOMIC\n" +
		"- **source-ids:** SRC-FILE-AAAAAAAAAAAA\n" +
		"- **file-path-line-range:** app/a.py:1\n" +
		"- **virtual-coordinate:** app/a.py :: job :: scan\n" +
		"- **pov-owner:** POV-3\n" +
		"- **general-intent:** Background scan.\n" +
		"- **inputs-outputs:** None\n" +
		"- **auth-binding:** None\n" +
		"- **staging-status:** " + status + "\n" +
		"- **coverage-status:** [C-GAP]\n" +
		"- **personas:** WEB\n" +
		"- **traversal-tracks:** Normal, Shadow/Conditional, Disabled\n" +
		"- **capability-state:** Unknown\n" +
		"- **runtime-state-matrix:** dev/production\n" +
		"- **depends-on-depended-on-by:** " + dependencies + "\n" +
		"- **synthesis-references:** None\n" +
		"- **side-effects:** None\n" +
		"- **exceptional-behaviors:** Unknown\n" +
		"- **edge-and-odd-cases:** Unknown\n" +
		"- **behavioral-notes-for-reconstruction:** Unknown\n" +
		"- **material-claims:** None\n" +
		"- **block-confidence-rank:** Low\n" +
		"- **block-evidence-profile:** [E-DIRECT] app/a.py:1\n" +
		"- **evidence-anchors:** sha256:aaa\n" +
		"- **record-version-supersession:** version 1\n")
}

func TestValidateDeadCodePositiveAndNegative(t *testing.T) {
	doc, db := testDocDB(t)
	seedSource(t, db, doc, "SRC-FILE-AAAAAAAAAAAA")
	if err := store.Ingest(db, doc, "ATOMIC-COMPONENT", "CMP-AAAAAAAAAAAA", 1, componentPayload("CMP-AAAAAAAAAAAA", "[S-DEAD-CODE]", "None")); err != nil {
		t.Fatal(err)
	}
	if err := ValidateDeadCode(db, doc, "CMP-AAAAAAAAAAAA"); err != nil {
		t.Fatalf("valid dead-code proof rejected: %v", err)
	}

	// A committed caller in another component disproves the proof.
	if err := store.Ingest(db, doc, "ATOMIC-COMPONENT", "CMP-BBBBBBBBBBBB", 1, componentPayload("CMP-BBBBBBBBBBBB", "[S-DISCOVERED]", "CMP-AAAAAAAAAAAA")); err != nil {
		t.Fatal(err)
	}
	if err := ValidateDeadCode(db, doc, "CMP-AAAAAAAAAAAA"); err == nil || !strings.Contains(err.Error(), "discovered caller") {
		t.Fatalf("caller component accepted: %v", err)
	}
}

func TestValidateDeadCodeRejectsUnresolvedFrontier(t *testing.T) {
	doc, db := testDocDB(t)
	seedSource(t, db, doc, "SRC-FILE-AAAAAAAAAAAA")
	if err := store.Ingest(db, doc, "ATOMIC-COMPONENT", "CMP-AAAAAAAAAAAA", 1, componentPayload("CMP-AAAAAAAAAAAA", "[S-DEAD-CODE]", "None")); err != nil {
		t.Fatal(err)
	}
	frontier := []byte("- **origin-invocation-cluster-pov-track:** ALFA-WEB-POV-3-001 | WEB | POV-3 | Normal\n- **from-src-cmp:** CMP-AAAAAAAAAAAA\n- **target-concrete-coordinate-or-resolution-rule:** app/a.py\n- **target-src:** SRC-FILE-AAAAAAAAAAAA\n- **discovery-method:** other\n- **depth:** 1\n- **snapshot-fingerprint:** sha256:aaa\n- **state:** Open\n- **disposition-evidence:** pending\n- **child-frontier-ids:** None\n- **closing-invocation:** None\n")
	if err := store.Ingest(db, doc, "FRONTIER-RECORD", "FRT-AAAAAAAAAAAA", 1, frontier); err != nil {
		t.Fatal(err)
	}
	if err := ValidateDeadCode(db, doc, "CMP-AAAAAAAAAAAA"); err == nil || !strings.Contains(err.Error(), "unresolved frontier") {
		t.Fatalf("open frontier accepted: %v", err)
	}
	terminal := []byte("- **origin-invocation-cluster-pov-track:** ALFA-WEB-POV-3-001 | WEB | POV-3 | Normal\n- **from-src-cmp:** CMP-AAAAAAAAAAAA\n- **target-concrete-coordinate-or-resolution-rule:** app/a.py\n- **target-src:** SRC-FILE-AAAAAAAAAAAA\n- **discovery-method:** other\n- **depth:** 1\n- **snapshot-fingerprint:** sha256:aaa\n- **state:** Terminal-Mapped\n- **disposition-evidence:** mapped\n- **child-frontier-ids:** None\n- **closing-invocation:** ALFA-WEB-POV-3-001\n")
	if err := store.ReplaceRecord(db, doc, "FRONTIER-RECORD", "FRT-AAAAAAAAAAAA", 2, terminal); err != nil {
		t.Fatal(err)
	}
	if err := ValidateDeadCode(db, doc, "CMP-AAAAAAAAAAAA"); err != nil {
		t.Fatalf("terminal frontier rejected: %v", err)
	}
}

func TestValidateDeadCodeRequiresDenominatorAndInventory(t *testing.T) {
	doc, db := testDocDB(t)
	seedSource(t, db, doc, "SRC-FILE-AAAAAAAAAAAA")
	if err := store.Ingest(db, doc, "ATOMIC-COMPONENT", "CMP-AAAAAAAAAAAA", 1, componentPayload("CMP-AAAAAAAAAAAA", "[S-DEAD-CODE]", "None")); err != nil {
		t.Fatal(err)
	}
	if err := store.Ingest(db, doc, "ATOMIC-COMPONENT", "CMP-CCCCCCCCCCCC", 1, []byte(
		"- **record-kind:** ATOMIC\n- **source-ids:** SRC-FILE-MISSING00000\n- **file-path-line-range:** app/a.py:1\n- **virtual-coordinate:** app/a.py :: job :: scan\n- **pov-owner:** POV-3\n- **general-intent:** Background scan.\n- **inputs-outputs:** None\n- **auth-binding:** None\n- **staging-status:** [S-DEAD-CODE]\n- **coverage-status:** [C-GAP]\n- **personas:** WEB\n- **traversal-tracks:** None\n- **capability-state:** Unknown\n- **runtime-state-matrix:** None\n- **depends-on-depended-on-by:** None\n- **synthesis-references:** None\n- **side-effects:** None\n- **exceptional-behaviors:** Unknown\n- **edge-and-odd-cases:** Unknown\n- **behavioral-notes-for-reconstruction:** Unknown\n- **material-claims:** None\n- **block-confidence-rank:** Low\n- **block-evidence-profile:** [E-DIRECT] app/a.py:1\n- **evidence-anchors:** sha256:aaa\n- **record-version-supersession:** version 1\n")); err != nil {
		t.Fatal(err)
	}
	if err := ValidateDeadCode(db, doc, "CMP-CCCCCCCCCCCC"); err == nil || !strings.Contains(err.Error(), "denominator dimension") {
		t.Fatalf("missing denominator accepted: %v", err)
	}
}

func TestValidateLifecycleEligibility(t *testing.T) {
	doc, db := testDocDB(t)
	seedSource(t, db, doc, "SRC-FILE-AAAAAAAAAAAA")
	row := func(id, scope, classification, bindings, snapshot string) []byte {
		return []byte("- **scope-id:** " + scope + "\n- **classification:** " + classification + "\n- **reason:** bounded static evidence\n- **clm-evidence-bindings:** " + bindings + "\n- **classified-by-invocation:** ALFA-NONE-NONE-001\n- **snapshot:** " + snapshot + "\n")
	}
	if err := store.Ingest(db, doc, "LIFECYCLE-ELIGIBILITY", "LCE-1", 1, row("LCE-1", "SRC-FILE-AAAAAAAAAAAA", "Unknown", "CLM-AAAAAAAAAAAA", "dev")); err != nil {
		t.Fatal(err)
	}
	diagnostics, err := ValidateLifecycleEligibility(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	if len(diagnostics) != 1 || !strings.Contains(diagnostics[0], "without a ticket") {
		t.Fatalf("diagnostics = %v", diagnostics)
	}
	ticket := []byte("- **source-pov-file-invocation:** SRC-FILE-AAAAAAAAAAAA | 03-BACKGROUND-POV.md | ALFA-WEB-POV-3-001\n- **target-pov-ledger:** POV-4 | personas/_shared/04-DATA-POV.md\n- **status:** T-OPEN\n- **trace-history:** origin:ALFA -> T-OPEN\n- **affected-ids:** SRC-FILE-AAAAAAAAAAAA\n- **memory-anchor:** bounded static read\n- **description:** runtime behavior unknown\n- **static-investigation:** bounded static paths attempted\n- **escalation-reason:** RUNTIME-REQUIRED\n- **probe-feasibility:** safe isolated probe available\n- **coverage-promotion-effect:** blocks the affected row\n")
	if err := store.Ingest(db, doc, "TICKET", "TKT-WEB-POV-3-001", 1, ticket); err != nil {
		t.Fatal(err)
	}
	diagnostics, err = ValidateLifecycleEligibility(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %v", diagnostics)
	}
	// A duplicate scope row and a missing-provenance row are rejected.
	if err := store.Ingest(db, doc, "LIFECYCLE-ELIGIBILITY", "LCE-2", 1, row("LCE-2", "SRC-FILE-AAAAAAAAAAAA", "Eligible", "None", "production")); err != nil {
		t.Fatal(err)
	}
	diagnostics, err = ValidateLifecycleEligibility(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(diagnostics, "\n")
	if !strings.Contains(joined, "duplicates the scope row") || !strings.Contains(joined, "no claim/evidence provenance") || !strings.Contains(joined, "inconsistent with the header snapshot") {
		t.Fatalf("diagnostics = %v", diagnostics)
	}
}
