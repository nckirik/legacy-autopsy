package analysis

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/markdown"
)

var update = flag.Bool("update", false, "rewrite analysis/metrics.golden.json")

const (
	protocolPath       = "../../protocol.md"
	classificationPath = "../../analysis/protocol-classification.json"
	metricsPath        = "../../analysis/metrics.golden.json"
)

func loadProtocol(t *testing.T) (*markdown.Document, string) {
	t.Helper()
	b, err := os.ReadFile(protocolPath)
	if err != nil {
		t.Fatal(err)
	}
	doc, err := markdown.Parse(string(b))
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(b)
	return doc, "sha256:" + hex.EncodeToString(sum[:])
}

func TestClassificationCoverage(t *testing.T) {
	doc, fingerprint := loadProtocol(t)
	c, err := Load(classificationPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := Validate(doc, fingerprint, c); err != nil {
		t.Fatal(err)
	}
}

func TestMetricsGolden(t *testing.T) {
	doc, _ := loadProtocol(t)
	c, err := Load(classificationPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := Validate(doc, c.SourceFingerprint, c); err != nil {
		t.Fatal(err)
	}
	metrics, err := Compute(doc, c)
	if err != nil {
		t.Fatal(err)
	}
	got, err := metrics.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if *update {
		if err := os.MkdirAll(filepath.Dir(metricsPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(metricsPath, got, 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(metricsPath)
	if err != nil {
		t.Fatalf("read metrics golden (run: go test ./internal/analysis -update): %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("metrics golden mismatch\ngot:\n%s", got)
	}
}
