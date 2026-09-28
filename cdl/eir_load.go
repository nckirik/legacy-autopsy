package cdl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// LoadEIR parses and verifies a canonical EIR document. It fails closed on an
// unknown format or protocol, unknown JSON fields, trailing content, a missing
// or mismatched eir-sha256, duplicate identities, and inconsistent derived
// numbering. Verification recomputes values; it never trusts stored ones.
func LoadEIR(data []byte) (*EIRDoc, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var doc EIRDoc
	if err := dec.Decode(&doc); err != nil {
		return nil, fmt.Errorf("eir: parse: %w", err)
	}
	if _, err := dec.Token(); err != io.EOF {
		return nil, fmt.Errorf("eir: trailing content after the EIR document")
	}
	if doc.Envelope.EIRFormat != EIRFormat {
		return nil, fmt.Errorf("eir: unsupported eir-format %d, want %d", doc.Envelope.EIRFormat, EIRFormat)
	}
	if doc.Envelope.Protocol != ProtocolVersion {
		return nil, fmt.Errorf("eir: protocol mismatch: got %q, want %q", doc.Envelope.Protocol, ProtocolVersion)
	}
	if doc.Envelope.Language == "" || doc.Envelope.Stdlib == "" || doc.Envelope.Generator == "" || doc.Envelope.SourceFingerprint == "" {
		return nil, fmt.Errorf("eir: envelope is missing language, stdlib, generator, or source fingerprint")
	}
	if doc.Envelope.EIRHash == "" {
		return nil, fmt.Errorf("eir: missing eir-sha256")
	}
	hashInput := doc
	hashInput.Envelope.EIRHash = ""
	canonical, err := CanonicalBytes(hashInput)
	if err != nil {
		return nil, fmt.Errorf("eir: canonicalize: %w", err)
	}
	if got := Fingerprint(canonical); got != doc.Envelope.EIRHash {
		return nil, fmt.Errorf("eir: hash mismatch: got %s, want %s", got, doc.Envelope.EIRHash)
	}
	if err := validateEIR(&doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func validateEIR(doc *EIRDoc) error {
	parts := map[string]bool{}
	for i, part := range doc.Declarations.Parts {
		want := strconv.Itoa(i + 1)
		if part.Number != want {
			return fmt.Errorf("eir: part %d has number %q, want %q", i, part.Number, want)
		}
		if parts[part.Number] {
			return fmt.Errorf("eir: duplicate part number %q", part.Number)
		}
		parts[part.Number] = true
	}

	ids := map[string]string{}
	addID := func(id, kind string) error {
		if id == "" {
			return fmt.Errorf("eir: empty %s id", kind)
		}
		if prior, ok := ids[id]; ok {
			return fmt.Errorf("eir: id %q is used by both %s and %s", id, prior, kind)
		}
		ids[id] = kind
		return nil
	}
	for _, m := range doc.Declarations.Modes {
		if err := addID(m.ID, "mode"); err != nil {
			return err
		}
	}
	for _, c := range doc.Declarations.Capabilities {
		if err := addID(c.ID, "capability"); err != nil {
			return err
		}
	}
	rules := map[string]bool{}
	for _, r := range doc.Declarations.Rules {
		if err := addID(r.ID, "rule"); err != nil {
			return err
		}
		rules[r.ID] = true
	}
	for _, e := range doc.Declarations.Enums {
		if err := addID(e.ID, "enum"); err != nil {
			return err
		}
	}
	for _, f := range doc.Declarations.Fields {
		if err := addID(f.ID, "field"); err != nil {
			return err
		}
	}
	for _, s := range doc.Declarations.States {
		if err := addID(s.ID, "state"); err != nil {
			return err
		}
	}
	for _, v := range doc.Declarations.Values {
		if err := addID(v.ID, "value"); err != nil {
			return err
		}
	}
	for _, ty := range doc.Declarations.Types {
		if err := addID(ty.ID, "type"); err != nil {
			return err
		}
	}

	numbers := map[string]bool{}
	preamble := 0
	perPart := map[string]int{}
	childCount := map[string]int{}
	completed := map[string]bool{}
	lastPart := ""
	for _, section := range doc.Sections {
		if err := addID(section.ID, "section"); err != nil {
			return err
		}
		if section.Goal == "" {
			return fmt.Errorf("eir: section %s has no normative prose", section.ID)
		}
		if section.Number == "" || numbers[section.Number] {
			return fmt.Errorf("eir: section %s has empty or duplicate number %q", section.ID, section.Number)
		}
		numbers[section.Number] = true
		segs := strings.Split(section.Number, ".")
		if segs[0] == "0" {
			if len(segs) != 2 {
				return fmt.Errorf("eir: preamble section %s has invalid number %q", section.ID, section.Number)
			}
			preamble++
			if section.Number != fmt.Sprintf("0.%d", preamble) {
				return fmt.Errorf("eir: preamble numbering gap at section %s (%q)", section.ID, section.Number)
			}
		} else {
			if !parts[segs[0]] {
				return fmt.Errorf("eir: section %s references unknown part %q", section.ID, segs[0])
			}
			if lastPart != segs[0] {
				if completed[segs[0]] {
					return fmt.Errorf("eir: sections of part %q are not contiguous", segs[0])
				}
				if lastPart != "" {
					completed[lastPart] = true
				}
				lastPart = segs[0]
			}
			switch len(segs) {
			case 2:
				perPart[segs[0]]++
				if segs[1] != strconv.Itoa(perPart[segs[0]]) {
					return fmt.Errorf("eir: part %s numbering gap at section %s (%q)", segs[0], section.ID, section.Number)
				}
			case 3:
				parent := segs[0] + "." + segs[1]
				if !numbers[parent] {
					return fmt.Errorf("eir: section %s has missing or later parent %q", section.ID, parent)
				}
				childCount[parent]++
				if segs[2] != strconv.Itoa(childCount[parent]) {
					return fmt.Errorf("eir: subsection numbering gap at section %s (%q)", section.ID, section.Number)
				}
			default:
				return fmt.Errorf("eir: section %s has unsupported number depth %q", section.ID, section.Number)
			}
		}
		for _, step := range section.Steps {
			if err := addID(step.ID, "step"); err != nil {
				return err
			}
			for _, use := range step.UsesRules {
				if !rules[use] {
					return fmt.Errorf("eir: step %s references unknown rule %q", step.ID, use)
				}
			}
		}
		for _, use := range section.Uses {
			if use.Kind == "RULE" && !rules[use.ID] {
				return fmt.Errorf("eir: section %s references unknown rule %q", section.ID, use.ID)
			}
		}
	}
	return nil
}
