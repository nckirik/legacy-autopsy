// Package parity holds the two pilot parity tests: backend conformance (test A)
// and protocol parity against the frozen fixture oracle (test B).
package parity

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/reference"
	"github.com/nckirik/legacy-autopsy/internal/runtime"
)

var fixtureNames = []string{"reconciled.json", "incomplete.json", "container-only.json", "hash-unavailable.json"}

func root(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func compileExample(t *testing.T) *cdl.Result {
	t.Helper()
	return compileSection(t, "protocol/export-reconciliation.cdl")
}

func protocolPathIn(root string) string { return protocol.OraclePath(root) }

// compileSection compiles the shared globals plus the given CDL sources against
// the committed ledger.
func compileSection(t *testing.T, rels ...string) *cdl.Result {
	t.Helper()
	root := root(t)
	globals, err := os.ReadFile(filepath.Join(root, "protocol/globals.cdl"))
	if err != nil {
		t.Fatal(err)
	}
	ledger, err := cdl.LoadLedger(filepath.Join(root, "protocol/identity-ledger.json"))
	if err != nil {
		t.Fatal(err)
	}
	sources := []cdl.Source{{Path: "protocol/globals.cdl", Bytes: globals}}
	for _, rel := range rels {
		src, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatal(err)
		}
		sources = append(sources, cdl.Source{Path: rel, Bytes: src})
	}
	res, err := cdl.Compile(cdl.CompileInput{
		Sources: sources,
		Ledger:  ledger,
	})
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func fixturePath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join(root(t), "protocol/fixtures", name)
}

func capabilityMap(ids []string) map[string]bool {
	out := map[string]bool{}
	for _, id := range ids {
		out[id] = true
	}
	return out
}

// TestBackendConformance is test A: the native runtime and the reference VM must
// produce the same canonical execution trace for the same EIR and inputs.
func TestBackendConformance(t *testing.T) {
	res := compileExample(t)
	for _, name := range fixtureNames {
		t.Run(name, func(t *testing.T) {
			path := fixturePath(t, name)
			nativeFx, err := runtime.LoadFixture(path)
			if err != nil {
				t.Fatal(err)
			}
			referenceFx, err := reference.LoadFixture(path)
			if err != nil {
				t.Fatal(err)
			}
			caps := capabilityMap(nativeFx.Capabilities)
			nativeTrace, err := runtime.Execute(res.EIR, nativeFx, caps)
			if err != nil {
				t.Fatalf("native execute: %v", err)
			}
			referenceTrace, err := reference.Execute(res.EIR, referenceFx, caps)
			if err != nil {
				t.Fatalf("reference execute: %v", err)
			}
			nativeBytes, err := cdl.CanonicalBytes(nativeTrace)
			if err != nil {
				t.Fatal(err)
			}
			referenceBytes, err := cdl.CanonicalBytes(referenceTrace)
			if err != nil {
				t.Fatal(err)
			}
			nativeHash, err := nativeTrace.Hash()
			if err != nil {
				t.Fatal(err)
			}
			referenceHash, err := referenceTrace.Hash()
			if err != nil {
				t.Fatal(err)
			}
			if string(nativeBytes) != string(referenceBytes) {
				t.Fatalf("trace mismatch\nnative:    %s\nreference: %s", nativeBytes, referenceBytes)
			}
			if nativeHash != referenceHash {
				t.Fatalf("trace hash mismatch: %s vs %s", nativeHash, referenceHash)
			}
		})
	}
}
