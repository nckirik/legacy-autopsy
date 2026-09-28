package parityreport

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "regenerate the crown parity report")

const (
	jsonPath = "analysis/parity-report.json"
	mdPath   = "analysis/parity-report.md"
)

func TestParityReport(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	report, err := Build(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, section := range report.Sections {
		if section.Status == "missing" {
			t.Fatalf("protocol section %s %q has no CDL source", section.Number, section.Heading)
		}
		if section.Status == "implemented" && section.CDLSource == "" {
			t.Fatalf("protocol section %s %q is implemented but has no CDL source", section.Number, section.Heading)
		}
	}

	jsonBytes, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	jsonBytes = append(jsonBytes, '\n')
	markdownBytes := []byte(RenderMarkdown(report))

	if *update {
		if err := os.WriteFile(filepath.Join(root, jsonPath), jsonBytes, 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, mdPath), markdownBytes, 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}

	for path, want := range map[string][]byte{jsonPath: jsonBytes, mdPath: markdownBytes} {
		got, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			t.Fatalf("%s: %v (run: go test ./internal/parityreport -update)", path, err)
		}
		if string(got) != string(want) {
			t.Fatalf("%s is stale (run: go test ./internal/parityreport -update)", path)
		}
	}
}
