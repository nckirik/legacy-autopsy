package cdl

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFormatGoldenAllForms(t *testing.T) {
	src, err := os.ReadFile("testdata/all-forms.cdl")
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("testdata/all-forms.golden.cdl")
	if err != nil {
		t.Fatal(err)
	}
	got, diags := FormatSource("testdata/all-forms.cdl", src)
	if len(diags) > 0 {
		t.Fatalf("%v", diags)
	}
	if string(got) != string(want) {
		t.Fatalf("canonical form mismatch; regenerate with: go run ./cmd/cdl fmt cdl/testdata/all-forms.cdl > cdl/testdata/all-forms.golden.cdl\n--- got ---\n%s", got)
	}
}

func TestProtocolSourcesCanonical(t *testing.T) {
	files, err := filepath.Glob("../protocol/*.cdl")
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range files {
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		formatted, diags := FormatSource(path, src)
		if len(diags) > 0 {
			t.Fatalf("%s: %v", path, diags)
		}
		if string(formatted) != string(src) {
			t.Fatalf("%s is not canonically formatted; run: go run ./cmd/cdl fmt -w %s", path, path)
		}
	}
}
