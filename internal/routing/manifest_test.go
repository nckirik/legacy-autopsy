package routing

import (
	"path/filepath"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
)

func TestManifestMatchesProtocol(t *testing.T) {
	root := filepath.Join("..", "..")
	model, err := protocol.Load(filepath.Join(root, "protocol.md"))
	if err != nil {
		t.Fatal(err)
	}
	manifest, err := Load(filepath.Join(root, "skill", "modes.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := Validate(root, manifest, model); err != nil {
		t.Fatal(err)
	}
}
