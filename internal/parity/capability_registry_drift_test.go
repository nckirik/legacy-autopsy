package parity

import (
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/spec"
)

// TestCapabilityRegistryMatchesDeclarations binds the deterministic capability
// registry to the capabilities declared in CDL. A registry entry without a
// declaration, a declared capability without an implementation, a kind
// mismatch, or an insufficient version is a defect.
func TestCapabilityRegistryMatchesDeclarations(t *testing.T) {
	res, err := spec.Compile(root(t))
	if err != nil {
		t.Fatal(err)
	}
	registry := capabilities.DefaultRegistry()
	declared := map[string]bool{}
	for _, capability := range res.EIR.Declarations.Capabilities {
		if declared[capability.ID] {
			t.Fatalf("capability %s declared more than once", capability.ID)
		}
		declared[capability.ID] = true
		entry, ok := registry.Lookup(capability.ID)
		if !ok {
			t.Fatalf("declared capability %s has no registry implementation", capability.ID)
		}
		if entry.Kind != capability.Kind {
			t.Fatalf("capability %s kind: registry=%s declared=%s", capability.ID, entry.Kind, capability.Kind)
		}
		if err := registry.Require(capability.ID, capability.Version); err != nil {
			t.Fatalf("capability %s: %v", capability.ID, err)
		}
	}
	for _, id := range registry.IDs() {
		if !declared[id] {
			t.Fatalf("registered capability %s is not declared in CDL", id)
		}
	}
}
