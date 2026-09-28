package cdl

import (
	"os"
	"strings"
	"testing"
)

func TestLoadEIRGolden(t *testing.T) {
	data, err := os.ReadFile("../protocol/golden/protocol.eir.json")
	if err != nil {
		t.Fatal(err)
	}
	doc, err := LoadEIR(data)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Envelope.EIRFormat != EIRFormat {
		t.Fatalf("format %d", doc.Envelope.EIRFormat)
	}
	if len(doc.Sections) == 0 {
		t.Fatal("no sections loaded")
	}
}

func TestLoadEIRRejectsTampering(t *testing.T) {
	data, err := os.ReadFile("../protocol/golden/protocol.eir.json")
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)

	if _, err := LoadEIR([]byte(strings.Replace(text, `"eir-format":3`, `"eir-format":9`, 1))); err == nil {
		t.Fatal("unknown format accepted")
	}
	if _, err := LoadEIR([]byte(strings.Replace(text, `"eir-format":3`, `"eir-format":3,"extra":1`, 1))); err == nil {
		t.Fatal("unknown field accepted")
	}
	tampered := strings.Replace(text, `"eir-sha256":"sha256:`, `"eir-sha256":"sha256:f`, 1)
	if _, err := LoadEIR([]byte(tampered)); err == nil {
		t.Fatal("hash mismatch accepted")
	}
}

func minimalEIR(sections []EIRSection) []byte {
	return minimalEIRWithParts(sections, nil)
}

func minimalEIRWithParts(sections []EIRSection, parts []EIRPart) []byte {
	doc := &EIRDoc{
		Envelope: EIREnvelope{
			EIRFormat:         EIRFormat,
			Language:          LanguageVersion,
			Stdlib:            StdlibVersion,
			Protocol:          ProtocolVersion,
			SourceFingerprint: "sha256:test",
			Generator:         GeneratorVersion,
		},
	}
	doc.Sections = sections
	doc.Declarations.Parts = parts
	input := *doc
	input.Envelope.EIRHash = ""
	canonical, err := CanonicalBytes(input)
	if err != nil {
		panic(err)
	}
	doc.Envelope.EIRHash = Fingerprint(canonical)
	out, err := CanonicalBytes(doc)
	if err != nil {
		panic(err)
	}
	return out
}

func TestLoadEIRRejectsInvalidStructure(t *testing.T) {
	cases := map[string][]EIRSection{
		"duplicate number": {
			{ID: "a", Number: "0.1", Title: "A", Goal: "g"},
			{ID: "b", Number: "0.1", Title: "B", Goal: "g"},
		},
		"numbering gap": {
			{ID: "a", Number: "0.1", Title: "A", Goal: "g"},
			{ID: "b", Number: "0.3", Title: "B", Goal: "g"},
		},
		"unknown part": {
			{ID: "a", Number: "2.1", Title: "A", Goal: "g"},
		},
		"missing prose": {
			{ID: "a", Number: "0.1", Title: "A"},
		},
		"unknown rule": {
			{
				ID: "a", Number: "0.1", Title: "A", Goal: "g",
				Steps: []EIRStep{{ID: "s", Owner: "MACHINE", UsesRules: []string{"missing"}}},
			},
		},
	}
	for name, sections := range cases {
		if _, err := LoadEIR(minimalEIR(sections)); err == nil {
			t.Fatalf("%s accepted", name)
		}
	}
}

func TestLoadEIRAcceptsMinimalDocument(t *testing.T) {
	sections := []EIRSection{
		{ID: "a", Number: "0.1", Title: "A", Goal: "prose"},
		{ID: "b", Number: "1.1", Title: "B", Goal: "prose"},
		{ID: "c", Number: "1.1.1", Title: "C", Goal: "prose"},
		{ID: "d", Number: "1.2", Title: "D", Goal: "prose"},
	}
	parts := []EIRPart{{Number: "1", Title: "Part"}}
	if _, err := LoadEIR(minimalEIRWithParts(sections, parts)); err != nil {
		t.Fatalf("valid document rejected: %v", err)
	}
}
