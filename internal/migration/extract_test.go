package migration

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"os"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/markdown"
)

var update = flag.Bool("update", false, "rewrite generated migration sources")

const (
	protocolPath = "../../protocol/legacy/protocol-4.1.2.md"
	outputPath   = "../../protocol/status-taxonomy.cdl"
)

func generate(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile(protocolPath)
	if err != nil {
		t.Fatal(err)
	}
	doc, err := markdown.Parse(string(b))
	if err != nil {
		t.Fatal(err)
	}
	taxonomy, err := ExtractStatusTaxonomy(doc)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(b)
	return []byte(RenderStatusTaxonomy(taxonomy, "sha256:"+hex.EncodeToString(sum[:])))
}

// TestStatusTaxonomyGeneration requires the checked-in generated source to match a
// fresh deterministic extraction byte-for-byte.
func TestStatusTaxonomyGeneration(t *testing.T) {
	got := generate(t)
	if *update {
		if err := os.WriteFile(outputPath, got, 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read generated source (run: go test ./internal/migration -update): %v", err)
	}
	if string(got) != string(want) {
		t.Fatal("generated status taxonomy drifted from protocol.md; run: go test ./internal/migration -update")
	}
}
