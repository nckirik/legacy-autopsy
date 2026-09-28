package capabilities

import "testing"

func TestNormalizeRelativePath(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want string
	}{
		{"./src//Café/./File.go", "src/Café/File.go"},
		{`src\payments\service.go`, "src/payments/service.go"},
	} {
		got, err := NormalizeRelativePath(tc.in)
		if err != nil {
			t.Fatalf("%q: %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("%q: got %q want %q", tc.in, got, tc.want)
		}
	}
}

func TestNormalizeRelativePathRejects(t *testing.T) {
	for _, in := range []string{"", "/etc/passwd", "C:/secret", "src/../secret", "src/./a/../b", "./..", "."} {
		if _, err := NormalizeRelativePath(in); err == nil {
			t.Fatalf("%q accepted", in)
		}
	}
}

func TestTypedIDVectors(t *testing.T) {
	id, err := TypedID("CMP", "billing-v1", "src/payments/service.go#Charge", "atomic-component")
	if err != nil {
		t.Fatal(err)
	}
	if id != "CMP-08ED474400F4" {
		t.Fatalf("got %s", id)
	}
	id, err = TypedID("SRC-FILE", "billing-v1", "src/payments/service.go", "source-unit")
	if err != nil {
		t.Fatal(err)
	}
	if id != "SRC-FILE-7A5180F89B03" {
		t.Fatalf("got %s", id)
	}
}
