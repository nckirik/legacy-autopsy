package capabilities

import "testing"

func TestDefaultRegistry(t *testing.T) {
	r := DefaultRegistry()
	want := []string{"hash.sha256", "table.serialize", "identity.typed-id", "path.normalize", "canonical.markdown"}
	if len(r.IDs()) != len(want) {
		t.Fatalf("ids = %v", r.IDs())
	}
	for _, id := range want {
		entry, ok := r.Lookup(id)
		if !ok || !entry.Available || entry.Kind != KindDeterministic || entry.Version != 1 {
			t.Fatalf("%s = %+v (ok=%v)", id, entry, ok)
		}
		if err := r.Require(id, 1); err != nil {
			t.Fatalf("%s: %v", id, err)
		}
	}
}

func TestRegistryFailsClosed(t *testing.T) {
	r := DefaultRegistry()
	if err := r.Require("missing", 1); err == nil {
		t.Fatal("unknown capability accepted")
	}
	if err := r.Require("hash.sha256", 2); err == nil {
		t.Fatal("version mismatch accepted")
	}
	if err := r.SetAvailable("hash.sha256", false); err != nil {
		t.Fatal(err)
	}
	if err := r.Require("hash.sha256", 1); err == nil {
		t.Fatal("unavailable capability accepted")
	}
	if err := r.SetAvailable("missing", true); err == nil {
		t.Fatal("unknown capability set")
	}
}

func TestRegistrySetFromIDs(t *testing.T) {
	r := DefaultRegistry()
	if err := r.SetFromIDs([]string{"hash.sha256", "missing"}); err == nil {
		t.Fatal("unknown id accepted")
	}
	if !r.Available("hash.sha256") {
		t.Fatal("failed set must not change availability")
	}
	if err := r.SetFromIDs([]string{"hash.sha256"}); err != nil {
		t.Fatal(err)
	}
	if !r.Available("hash.sha256") || r.Available("path.normalize") {
		t.Fatalf("unexpected availability: %v", r.AvailableIDs())
	}
	if len(Deterministic()) != 5 {
		t.Fatalf("Deterministic() = %v", Deterministic())
	}
}
