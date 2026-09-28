package capabilities

import (
	"strings"
	"testing"
)

func TestExtendCollisionPrefixes(t *testing.T) {
	base := "ABCDEF123456"
	a := base + strings.Repeat("0", 52)
	b := base + strings.Repeat("1", 52)
	gotA, gotB, err := ExtendCollisionPrefixes(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if len(gotA) != 16 || len(gotB) != 16 || gotA == gotB {
		t.Fatalf("extended to %q/%q", gotA, gotB)
	}
	if !strings.HasPrefix(gotA, base) || !strings.HasPrefix(gotB, base) {
		t.Fatalf("extension changed the colliding prefix: %q/%q", gotA, gotB)
	}

	// Differ later: extension continues to 20.
	c := base + "0000" + strings.Repeat("0", 48)
	d := base + "0000" + strings.Repeat("2", 48)
	gotC, gotD, err := ExtendCollisionPrefixes(c, d)
	if err != nil {
		t.Fatal(err)
	}
	if len(gotC) != 20 || len(gotD) != 20 {
		t.Fatalf("extended to %d/%d, want 20", len(gotC), len(gotD))
	}

	if _, _, err := ExtendCollisionPrefixes(base, base); err == nil {
		t.Fatal("identical hashes accepted")
	}
	if _, _, err := ExtendCollisionPrefixes("ABCDEF123456", "ABCDEF123457"+strings.Repeat("0", 52)); err == nil {
		t.Fatal("non-colliding hashes accepted")
	}
	if _, _, err := ExtendCollisionPrefixes("lower", "lower"); err == nil {
		t.Fatal("invalid hex accepted")
	}
}

func TestExtendCollisionKeysDiffer(t *testing.T) {
	if _, _, err := ExtendCollision("same", "same"); err == nil {
		t.Fatal("identical keys accepted")
	}
	// Distinct keys are astronomically unlikely to collide; the capability must
	// fail closed rather than fabricate an extension.
	if _, _, err := ExtendCollision("key-a", "key-b"); err == nil {
		t.Fatal("non-colliding keys produced an extension")
	}
}
