package capabilities

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
)

// expectedTypedID recomputes the §4.1 formula independently of the implementation.
func expectedTypedID(idType, namespace, owner, discriminator string) string {
	key := idType + "|" + namespace + "|" + owner + "|" + discriminator
	sum := sha256.Sum256([]byte(key))
	return idType + "-" + strings.ToUpper(hex.EncodeToString(sum[:]))[:12]
}

func TestTypedIDMatchesFormula(t *testing.T) {
	cases := []struct {
		idType, namespace, owner, discriminator string
	}{
		{"CMP", "example", "src/example.go#Run", "atomic-component"},
		{"SRC-FILE", "example", "src/example.go", "kind-file"},
		{"PRF", "example", "Q38FN", "persona-profile"},
		{"HBK", "example", "START-HERE.md#overview", "handbook-section"},
	}
	for _, tc := range cases {
		got, err := TypedID(tc.idType, tc.namespace, tc.owner, tc.discriminator)
		if err != nil {
			t.Fatalf("TypedID(%v): %v", tc, err)
		}
		want := expectedTypedID(tc.idType, tc.namespace, tc.owner, tc.discriminator)
		if got != want {
			t.Fatalf("TypedID(%v) = %s, want %s", tc, got, want)
		}
		if len(got) != len(tc.idType)+13 {
			t.Fatalf("TypedID(%v) = %s has unexpected length", tc, got)
		}
	}
}

func TestTypedIDPreservesUnicodeAndCase(t *testing.T) {
	got, err := TypedID("SRC-FILE", "example", "süß/Modules/Árvíztűrő.go", "unit")
	if err != nil {
		t.Fatal(err)
	}
	want := expectedTypedID("SRC-FILE", "example", "süß/Modules/Árvíztűrő.go", "unit")
	if got != want {
		t.Fatalf("unicode owner changed the identity: got %s want %s", got, want)
	}
}

func TestTypedIDRejectsInvalidInputs(t *testing.T) {
	cases := [][4]string{
		{"", "example", "owner", "discriminator"},
		{"cmp", "example", "owner", "discriminator"},
		{"CMP", "", "owner", "discriminator"},
		{"CMP", "example", "", "discriminator"},
		{"CMP", "example", "owner", ""},
	}
	for _, tc := range cases {
		if _, err := TypedID(tc[0], tc[1], tc[2], tc[3]); err == nil {
			t.Fatalf("TypedID(%v) unexpectedly succeeded", tc)
		}
	}
}
