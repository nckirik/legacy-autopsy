package cdl

import (
	"crypto/sha256"
	"encoding/hex"
)

// Source is one ordered source input.
type Source struct {
	Path  string
	Bytes []byte
}

// CompileInput is one compilation request. Ledger must be supplied (loaded from
// the committed identity ledger); nil fails the ledger invariant.
type CompileInput struct {
	Sources  []Source
	Ledger   map[string]IDRecord
	Versions Versions
}

// Result is a successful compilation.
type Result struct {
	Program           *Program
	EIR               *EIRDoc
	EIRBytes          []byte
	EIRHash           string
	SourceFingerprint string
	IDs               []IDRecord
}

// Compile runs parse, resolution, checks, ledger verification, and EIR emission.
func Compile(in CompileInput) (*Result, error) {
	if in.Versions.Language == "" {
		in.Versions = Versions{Language: LanguageVersion, Stdlib: StdlibVersion, Generator: GeneratorVersion}
	}
	var prog Program
	var diags Diagnostics
	for _, src := range in.Sources {
		parsed, d := Parse(src.Path, src.Bytes)
		diags = append(diags, d...)
		if parsed == nil {
			continue
		}
		prog.Globals.Capabilities = append(prog.Globals.Capabilities, parsed.Globals.Capabilities...)
		prog.Globals.Registries = append(prog.Globals.Registries, parsed.Globals.Registries...)
		prog.Globals.Artifacts = append(prog.Globals.Artifacts, parsed.Globals.Artifacts...)
		prog.Globals.Rules = append(prog.Globals.Rules, parsed.Globals.Rules...)
		prog.Globals.Gates = append(prog.Globals.Gates, parsed.Globals.Gates...)
		prog.Globals.WorkflowTargets = append(prog.Globals.WorkflowTargets, parsed.Globals.WorkflowTargets...)
		prog.Sections = append(prog.Sections, parsed.Sections...)
		prog.Projections = append(prog.Projections, parsed.Projections...)
	}
	if len(diags) > 0 {
		return nil, diags
	}
	semDiags, _ := resolve(&prog)
	diags = append(diags, semDiags...)
	ids := collectIDs(&prog)
	diags = append(diags, checkLedger(ids, in.Ledger)...)
	if len(diags) > 0 {
		return nil, diags
	}
	fingerprint := sourceFingerprint(in.Sources)
	doc, eirBytes, eirHash, err := BuildEIR(&prog, in.Versions, fingerprint)
	if err != nil {
		return nil, err
	}
	return &Result{
		Program:           &prog,
		EIR:               doc,
		EIRBytes:          eirBytes,
		EIRHash:           eirHash,
		SourceFingerprint: fingerprint,
		IDs:               ids,
	}, nil
}

func sourceFingerprint(sources []Source) string {
	h := sha256.New()
	for _, s := range sources {
		h.Write([]byte(s.Path))
		h.Write([]byte{0})
		h.Write(s.Bytes)
		h.Write([]byte{0})
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil))
}
