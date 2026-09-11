package protocol

import (
	"path/filepath"
	"testing"
)

func TestLoadNormativeProtocol(t *testing.T) {
	model, err := Load(filepath.Join("..", "..", "protocol.md"))
	if err != nil {
		t.Fatal(err)
	}
	if model.Version != "4.1.1" {
		t.Fatalf("version %q", model.Version)
	}
	if len(model.Modes) != 13 {
		t.Fatalf("mode count %d", len(model.Modes))
	}
	if !model.HasSection("§8.1. Resume identity header") {
		t.Fatal("required section not found")
	}
}
