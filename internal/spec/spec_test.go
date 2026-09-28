package spec

import "testing"

func TestCompileAssembly(t *testing.T) {
	res, err := Compile("../..")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.EIR.Declarations.Modes) != 13 {
		t.Fatalf("%d modes, want 13", len(res.EIR.Declarations.Modes))
	}
	found := false
	for _, section := range res.EIR.Sections {
		if section.Number == "7.7" {
			found = true
		}
	}
	if !found {
		t.Fatal("section 7.7 missing from the assembly")
	}
	if res.EIR.Envelope.SourceFingerprint == "" || res.EIRHash == "" {
		t.Fatal("assembly envelope lacks fingerprints")
	}
}

func TestModes(t *testing.T) {
	modes, err := Modes("../..")
	if err != nil {
		t.Fatal(err)
	}
	strict := 0
	for _, mode := range modes {
		if mode.Strict {
			strict++
		}
		if len(mode.Reads) == 0 {
			t.Fatalf("mode %s has no materialized read set", mode.ID)
		}
	}
	if strict != 3 {
		t.Fatalf("%d strict modes, want 3", strict)
	}
}
