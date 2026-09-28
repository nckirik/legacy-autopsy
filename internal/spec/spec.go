// Package spec compiles the checked-in CDL document assembly for in-process
// runtime consumption. Legacy Autopsy depends on cdl; cdl never depends on this
// package or any other runtime package.
package spec

import (
	"path/filepath"

	"github.com/nckirik/legacy-autopsy/cdl"
)

const (
	assemblyPath = "examples/spec/assembly.json"
	ledgerPath   = "examples/spec/identity-ledger.json"
)

// Compile compiles the repository's ordered CDL assembly with its committed
// identity ledger.
func Compile(repoRoot string) (*cdl.Result, error) {
	sources, err := cdl.LoadAssembly(repoRoot, assemblyPath)
	if err != nil {
		return nil, err
	}
	ledger, err := cdl.LoadLedger(filepath.Join(repoRoot, filepath.FromSlash(ledgerPath)))
	if err != nil {
		return nil, err
	}
	return cdl.Compile(cdl.CompileInput{Sources: sources, Ledger: ledger})
}

// Modes returns the compiled invocation modes with materialized read sets.
func Modes(repoRoot string) ([]cdl.EIRMode, error) {
	res, err := Compile(repoRoot)
	if err != nil {
		return nil, err
	}
	return res.EIR.Declarations.Modes, nil
}
