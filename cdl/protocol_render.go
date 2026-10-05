package cdl

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// RenderProtocol renders the canonical Markdown edition from a compiled EIR
// document. The render is deterministic and covers every section, declaration,
// rule, field, enum, state, and mode in the document.
func RenderProtocol(doc *EIRDoc) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "# Protocol: Legacy System Deconstruction, Assurance, and Reconstruction\n\n")
	fmt.Fprintf(&b, "**Version:** 4.2 (Canonical Reconstruction-Ready Edition)\n")
	fmt.Fprintf(&b, "**Status:** Normative\n")
	fmt.Fprintf(&b, "**Purpose:** Produce an evidence-grounded, complete, framework-agnostic description of a legacy system and a separately reviewed reconstruction package without requiring downstream readers to reopen the legacy source.\n")
	fmt.Fprintf(&b, "**Generated-From:** %s (%s)\n", doc.Envelope.Generator, doc.Envelope.SourceFingerprint)
	fmt.Fprintf(&b, "**Language:** %s (eir-format %d)\n", doc.Envelope.Language, doc.Envelope.EIRFormat)
	fmt.Fprintf(&b, "**Authority:** generated render of the canonical CDL sources; do not edit.\n\n")
	fmt.Fprintf(&b, "---\n\n")

	renderGlobals(&b, doc)

	partTitles := map[string]string{}
	for _, part := range doc.Declarations.Parts {
		partTitles[part.Number] = part.Title
	}
	sections := append([]EIRSection(nil), doc.Sections...)
	sort.SliceStable(sections, func(i, j int) bool {
		return compareNumbers(sections[i].Number, sections[j].Number) < 0
	})
	currentPart := ""
	for _, section := range sections {
		top := strings.SplitN(section.Number, ".", 2)[0]
		if title, ok := partTitles[top]; ok && top != currentPart {
			fmt.Fprintf(&b, "# Part %s. %s\n\n", top, title)
			currentPart = top
		}
		renderSection(&b, doc, section)
	}

	fmt.Fprintf(&b, "---\n\n")
	fmt.Fprintf(&b, "**End of Canonical Deconstruction Protocol v4.2**\n")
	return []byte(b.String())
}

func renderGlobals(b *strings.Builder, doc *EIRDoc) {
	fmt.Fprintf(b, "## Document Declarations\n\n")
	if len(doc.Declarations.Capabilities) > 0 {
		fmt.Fprintf(b, "### Capabilities\n\n| Capability | Kind | Version |\n| :-- | :-- | --: |\n")
		for _, capability := range doc.Declarations.Capabilities {
			fmt.Fprintf(b, "| %s | %s | %d |\n", capability.ID, capability.Kind, capability.Version)
		}
		fmt.Fprintf(b, "\n")
	}
	renderStringList(b, "Registries", doc.Declarations.Registries)
	renderStringList(b, "Artifacts", doc.Declarations.Artifacts)
	renderStringList(b, "Gates", doc.Declarations.Gates)
	renderStringList(b, "Workflow targets", doc.Declarations.WorkflowTargets)
	renderStringList(b, "Base reads", doc.Declarations.BaseReads)
	if len(doc.Declarations.Modes) > 0 {
		fmt.Fprintf(b, "### Modes\n\n| Mode | Strict | Reads |\n| :-- | :-- | --: |\n")
		for _, mode := range doc.Declarations.Modes {
			fmt.Fprintf(b, "| %s | %t | %d |\n", mode.ID, mode.Strict, len(mode.Reads))
		}
		fmt.Fprintf(b, "\n")
	}
}

func renderStringList(b *strings.Builder, title string, values []string) {
	if len(values) == 0 {
		return
	}
	fmt.Fprintf(b, "### %s\n\n", title)
	for _, value := range values {
		fmt.Fprintf(b, "- %s\n", value)
	}
	fmt.Fprintf(b, "\n")
}

func renderSection(b *strings.Builder, doc *EIRDoc, section EIRSection) {
	fmt.Fprintf(b, "## %s. %s\n\n", section.Number, section.Title)
	if section.Artifact != "" {
		fmt.Fprintf(b, "**Artifact:** %s\n\n", section.Artifact)
	}
	if goal := strings.TrimSpace(section.Goal); goal != "" {
		fmt.Fprintf(b, "%s\n\n", goal)
		for _, rule := range doc.Declarations.Rules {
			if rule.Section != section.ID || rule.Predicate == "" {
				continue
			}
			fmt.Fprintf(b, "### Rule: %s\n\nPredicate: %s\n\n", rule.ID, rule.Predicate)
		}
		return
	}

	for _, rule := range doc.Declarations.Rules {
		if rule.Section != section.ID {
			continue
		}
		fmt.Fprintf(b, "### Rule: %s\n\n", rule.ID)
		if goal := strings.TrimSpace(rule.Goal); goal != "" {
			fmt.Fprintf(b, "%s\n\n", goal)
		} else if rule.Predicate != "" {
			fmt.Fprintf(b, "Predicate: %s\n\n", rule.Predicate)
		}
	}
	for _, enum := range doc.Declarations.Enums {
		if enum.Section != section.ID {
			continue
		}
		fmt.Fprintf(b, "### Enum: %s\n\n", enum.ID)
		for _, value := range enum.Values {
			fmt.Fprintf(b, "- %s\n", value)
		}
		fmt.Fprintf(b, "\n")
	}
	for _, state := range doc.Declarations.States {
		if state.Section != section.ID {
			continue
		}
		fmt.Fprintf(b, "### State: %s\n\n", state.ID)
		fmt.Fprintf(b, "- Values: %s\n- Initial: %s\n", strings.Join(state.Values, ", "), state.Initial)
		for _, allows := range state.Allows {
			fmt.Fprintf(b, "- %s -> %s\n", allows[0], allows[1])
		}
		fmt.Fprintf(b, "\n")
	}
	for _, field := range doc.Declarations.Fields {
		if field.Section != section.ID {
			continue
		}
		fmt.Fprintf(b, "### Field: %s\n\n| Field | Type | Required |\n| :-- | :-- | :-- |\n", field.ID)
		for _, spec := range field.Fields {
			fmt.Fprintf(b, "| %s | %s | %t |\n", spec.Name, spec.Type, spec.Required)
		}
		fmt.Fprintf(b, "\n")
	}
	for _, table := range doc.Declarations.Tables {
		if table.Section != section.ID {
			continue
		}
		fmt.Fprintf(b, "### Table: %s\n\nRows: %s, key: %s.\n\n", table.ID, table.RowType, table.Key)
	}
}

func compareNumbers(a, b string) int {
	ap := strings.Split(a, ".")
	bp := strings.Split(b, ".")
	for i := 0; i < len(ap) && i < len(bp); i++ {
		ai, aerr := strconv.Atoi(ap[i])
		bi, berr := strconv.Atoi(bp[i])
		if aerr != nil || berr != nil {
			return strings.Compare(ap[i], bp[i])
		}
		if ai != bi {
			return ai - bi
		}
	}
	return len(ap) - len(bp)
}
