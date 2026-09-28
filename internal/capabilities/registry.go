package capabilities

import (
	"fmt"
	"sort"
)

// Capability kinds, matching the declared capability taxonomy.
const (
	KindDeterministic = "deterministic"
	KindProposing     = "proposing"
)

// Entry is one registered capability implementation.
type Entry struct {
	ID        string
	Kind      string
	Version   int
	Available bool
}

// Registry is an ordered, versioned capability registry with explicit
// availability. It fails closed: unknown capabilities and version mismatches
// are errors, never silent degradation.
type Registry struct {
	order []string
	byID  map[string]Entry
}

// DefaultRegistry returns the registered deterministic capabilities. The set is
// drift-checked against the declared CDL capabilities by internal/parity.
func DefaultRegistry() *Registry {
	r := &Registry{byID: map[string]Entry{}}
	for _, entry := range []Entry{
		{ID: "hash.sha256", Kind: KindDeterministic, Version: 1, Available: true},
		{ID: "table.serialize", Kind: KindDeterministic, Version: 1, Available: true},
		{ID: "identity.typed-id", Kind: KindDeterministic, Version: 1, Available: true},
		{ID: "path.normalize", Kind: KindDeterministic, Version: 1, Available: true},
		{ID: "canonical.markdown", Kind: KindDeterministic, Version: 1, Available: true},
	} {
		r.order = append(r.order, entry.ID)
		r.byID[entry.ID] = entry
	}
	return r
}

// Lookup returns the registration for a capability.
func (r *Registry) Lookup(id string) (Entry, bool) {
	entry, ok := r.byID[id]
	return entry, ok
}

// IDs returns all registered capability IDs in registration order.
func (r *Registry) IDs() []string {
	out := make([]string, len(r.order))
	copy(out, r.order)
	return out
}

// AvailableIDs returns the available capability IDs, sorted for determinism.
func (r *Registry) AvailableIDs() []string {
	var out []string
	for _, id := range r.order {
		if r.byID[id].Available {
			out = append(out, id)
		}
	}
	sort.Strings(out)
	return out
}

// Available reports capability availability; unknown capabilities are not
// available.
func (r *Registry) Available(id string) bool {
	entry, ok := r.byID[id]
	return ok && entry.Available
}

// Version returns the registered version of a capability.
func (r *Registry) Version(id string) (int, bool) {
	entry, ok := r.byID[id]
	if !ok {
		return 0, false
	}
	return entry.Version, true
}

// SetAvailable toggles availability for a registered capability.
func (r *Registry) SetAvailable(id string, available bool) error {
	entry, ok := r.byID[id]
	if !ok {
		return fmt.Errorf("capabilities: unknown capability %q", id)
	}
	entry.Available = available
	r.byID[id] = entry
	return nil
}

// Require negotiates a capability at a minimum version and fails closed when it
// is unknown, unavailable, or too old.
func (r *Registry) Require(id string, minVersion int) error {
	entry, ok := r.byID[id]
	switch {
	case !ok:
		return fmt.Errorf("capabilities: unknown capability %q", id)
	case !entry.Available:
		return fmt.Errorf("capabilities: capability %q is unavailable", id)
	case entry.Version < minVersion:
		return fmt.Errorf("capabilities: capability %q version %d is below required %d", id, entry.Version, minVersion)
	}
	return nil
}

// SetFromIDs applies an availability set: every registered capability starts
// unavailable and each listed ID must be registered. All IDs are validated
// before any availability changes, so an unknown ID changes nothing.
func (r *Registry) SetFromIDs(ids []string) error {
	seen := map[string]bool{}
	for _, id := range ids {
		if _, ok := r.byID[id]; !ok {
			return fmt.Errorf("capabilities: unknown capability %q", id)
		}
		seen[id] = true
	}
	for id, entry := range r.byID {
		entry.Available = seen[id]
		r.byID[id] = entry
	}
	return nil
}
