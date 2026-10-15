package work

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/store"
)

func TestRuntimeAnnotatedFieldsRewriteFingerprints(t *testing.T) {
	db, doc := testDoc(t)
	seedInvariant(t, db, doc, "INV-AAAAAAAAAAAA", "seeded")
	seedInvariant(t, db, doc, "INV-DDDDDDDDDDDD", "second")
	fingerprint, err := store.FingerprintOf(db, doc, "INVARIANT", "INV-AAAAAAAAAAAA")
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.FingerprintOf(db, doc, "INVARIANT", "INV-DDDDDDDDDDDD")
	if err != nil {
		t.Fatal(err)
	}
	unit := Unit{
		ID:                     "U-ANNOTATE",
		Mode:                   "Discovery",
		RecordID:               "INV-BBBBBBBBBBBB",
		ResultType:             "INVARIANT",
		RuntimeAnnotatedFields: []string{"verification-vectors"},
	}
	payload := "- **Constraint:** c\n- **Verification Vectors:** INV-AAAAAAAAAAAA@sha256:0000 (prose suffix); INV-DDDDDDDDDDDD\n- **Discovered In:** SRC-1\n"
	if _, err := CommitWithIDs(NewCaps(db, doc), unit, Proposal{FieldID: "INVARIANT", Records: []Record{{Payload: []byte(payload)}}}, false); err != nil {
		t.Fatal(err)
	}
	stored, err := store.Get(db, doc, "INVARIANT", "INV-BBBBBBBBBBBB")
	if err != nil {
		t.Fatal(err)
	}
	want := "INV-AAAAAAAAAAAA@version 1 " + fingerprint + " (prose suffix), INV-DDDDDDDDDDDD@version 1 " + second
	if !strings.Contains(string(stored), want) {
		t.Fatalf("stored payload lacks exact annotations:\n%s\nwant: %s", stored, want)
	}
	if strings.Contains(string(stored), "@sha256:0000") {
		t.Fatalf("model-authored annotation survived:\n%s", stored)
	}
}

func TestRuntimeAnnotatedFieldsFailClosed(t *testing.T) {
	db, doc := testDoc(t)
	unit := Unit{
		ID:                     "U-ANNOTATE-FAIL",
		Mode:                   "Discovery",
		RecordID:               "INV-CCCCCCCCCCCC",
		ResultType:             "INVARIANT",
		RuntimeAnnotatedFields: []string{"verification-vectors"},
	}
	payload := "- **Constraint:** c\n- **Verification Vectors:** INV-MISSING000000\n- **Discovered In:** SRC-1\n"
	if _, err := CommitWithIDs(NewCaps(db, doc), unit, Proposal{FieldID: "INVARIANT", Records: []Record{{Payload: []byte(payload)}}}, false); err == nil || !strings.Contains(err.Error(), "INV-MISSING000000") {
		t.Fatalf("unresolved annotation accepted: %v", err)
	}
	missing := "- **Constraint:** c\n- **Discovered In:** SRC-1\n"
	if _, err := CommitWithIDs(NewCaps(db, doc), unit, Proposal{FieldID: "INVARIANT", Records: []Record{{Payload: []byte(missing)}}}, false); err == nil || !strings.Contains(err.Error(), "is missing") {
		t.Fatalf("missing annotated field accepted: %v", err)
	}
}

func TestRuntimeFingerprintCoversAnnotations(t *testing.T) {
	db, doc := testDoc(t)
	seedInvariant(t, db, doc, "INV-AAAAAAAAAAAA", "seeded")
	unit := Unit{
		ID:                       "U-ANNOTATE-FINGERPRINT",
		Mode:                     "Discovery",
		RecordID:                 "INV-EEEEEEEEEEEE",
		ResultType:               "INVARIANT",
		RuntimeAnnotatedFields:   []string{"constraint"},
		RuntimeFingerprintFields: []string{"verification-vectors"},
	}
	payload := "- **Constraint:** INV-AAAAAAAAAAAA\n- **Verification Vectors:** PENDING\n- **Discovered In:** SRC-1\n"
	if _, err := CommitWithIDs(NewCaps(db, doc), unit, Proposal{FieldID: "INVARIANT", Records: []Record{{Payload: []byte(payload)}}}, false); err != nil {
		t.Fatal(err)
	}
	stored, err := store.Get(db, doc, "INVARIANT", "INV-EEEEEEEEEEEE")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(stored), "INV-AAAAAAAAAAAA@version 1 sha256:") {
		t.Fatalf("annotation missing from stored payload:\n%s", stored)
	}
	values, err := PayloadText([]string{"verification-vectors"}, stored)
	if err != nil {
		t.Fatal(err)
	}
	recomputed, err := capabilities.SemanticContentFingerprint(stored, "verification-vectors")
	if err != nil {
		t.Fatal(err)
	}
	if values["verification-vectors"] != recomputed {
		t.Fatalf("stored fingerprint %q does not cover the annotated payload (recomputed %q)", values["verification-vectors"], recomputed)
	}
}
