// Package runtime is the native MACHINE VM and step executor for the frozen
// pilot subset. It consumes EIR and never reinterprets CDL source.
package runtime

import (
	"encoding/json"
	"os"

	"github.com/nckirik/legacy-autopsy/cdl"
)

// Fixture is deterministic execution input: capability availability, registry
// data, and AGENT stub results.
type Fixture struct {
	Capabilities []string                             `json:"capabilities"`
	Registries   map[string][]map[string]any          `json:"registries"`
	Agent        map[string]map[string]map[string]any `json:"agent"`
}

// LoadFixture reads a synthetic execution fixture.
func LoadFixture(path string) (Fixture, error) {
	var fx Fixture
	b, err := os.ReadFile(path)
	if err != nil {
		return fx, err
	}
	if err := json.Unmarshal(b, &fx); err != nil {
		return fx, err
	}
	return fx, nil
}

func (f Fixture) capabilitySet() map[string]bool {
	out := map[string]bool{}
	for _, c := range f.Capabilities {
		out[c] = true
	}
	return out
}

// Effect is one committed or drafted state change.
type Effect struct {
	Target    string `json:"target"`
	Operation string `json:"operation"`
	Value     string `json:"value,omitempty"`
}

// Diagnostic is one stable execution diagnostic.
type Diagnostic struct {
	Code    string `json:"code"`
	Subject string `json:"subject"`
}

// StepRecord is one step's canonical execution record.
type StepRecord struct {
	Step        string       `json:"step"`
	Owner       string       `json:"owner"`
	Outcome     string       `json:"outcome"`
	Authority   string       `json:"authority"`
	Effects     []Effect     `json:"effects"`
	Diagnostics []Diagnostic `json:"diagnostics"`
	EIRHash     string       `json:"eir-hash"`
	Channel     string       `json:"channel"`
}

// Trace is a canonical execution trace.
type Trace struct {
	Steps []StepRecord `json:"steps"`
}

// Hash returns the canonical trace hash used by backend-conformance test A.
func (t Trace) Hash() (string, error) {
	b, err := cdl.CanonicalBytes(t)
	if err != nil {
		return "", err
	}
	return cdl.Fingerprint(b), nil
}
