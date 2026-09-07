package identity

import "testing"

func TestGenerateDeterministicID(t *testing.T) {
	result, err := Generate("CMP", "", "billing-v1", "src/payments/service.go#Charge", "atomic-component")
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "CMP-08ED474400F4" {
		t.Fatalf("unexpected ID %s", result.ID)
	}
}

func TestGenerateSourceID(t *testing.T) {
	result, err := Generate("SRC", "FILE", "billing-v1", "src/payments/service.go", "source-unit")
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "SRC-FILE-7A5180F89B03" {
		t.Fatalf("unexpected ID %s", result.ID)
	}
}

func TestGenerateSourceFileNormalizesOwnerPath(t *testing.T) {
	left, err := Generate("SRC", "FILE", "billing-v1", "./src//payments/service.go", "source-unit")
	if err != nil {
		t.Fatal(err)
	}
	right, err := Generate("SRC", "FILE", "billing-v1", "src/payments/service.go", "source-unit")
	if err != nil {
		t.Fatal(err)
	}
	if left != right {
		t.Fatalf("normalized identities differ: %#v %#v", left, right)
	}
}

func TestGenerateRejectsDeferredSpecializedIdentity(t *testing.T) {
	if _, err := Generate("PRF", "", "billing-v1", "ABC", "persona-profile"); err == nil {
		t.Fatal("expected specialized identity rejection")
	}
}

func TestNormalizeRelativePath(t *testing.T) {
	got, err := NormalizeRelativePath("./src//Café/./File.go")
	if err != nil {
		t.Fatal(err)
	}
	if got != "src/Café/File.go" {
		t.Fatalf("got %q", got)
	}
	if _, err := NormalizeRelativePath("src/../secret"); err == nil {
		t.Fatal("expected traversal rejection")
	}
}
