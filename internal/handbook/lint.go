// Package handbook implements the deterministic handbook quality gates: the
// PI-11 substance lint (required sections plus claim-backed assertions) and the
// PI-12 readability review schema (no relay-only attestation).
package handbook

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/ste"
)

// RequiredSections are the chapter sections every handbook chapter must carry.
var RequiredSections = []string{"summary", "observed legacy behavior", "target decisions", "navigation"}

// MinimumAssertions is the PI-11 claim-backed assertion floor per chapter.
const MinimumAssertions = 8

// Finding is one deterministic lint finding.
type Finding struct {
	Code    string
	Subject string
	Detail  string
}

func (f Finding) String() string { return fmt.Sprintf("%s %s: %s", f.Code, f.Subject, f.Detail) }

// AssertionCount counts claim-backed assertion lines: statement lines that
// carry an inline claim citation or a numbered footnote reference. Footnote
// definitions, fenced diagram/code blocks, and the collapsed record appendix
// are not assertions, so the count is identical before and after citation
// numbering.
func AssertionCount(markdown string) int {
	if idx := strings.Index(markdown, "\n### Record Details"); idx >= 0 {
		markdown = markdown[:idx]
	}
	count := 0
	inFence := false
	for _, line := range strings.Split(markdown, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
			continue
		}
		if inFence || trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "- [") || strings.HasPrefix(trimmed, "[^") {
			continue
		}
		if claimReference(trimmed) || footnoteReference(trimmed) {
			count++
		}
	}
	return count
}

var footnoteReferencePattern = regexp.MustCompile(`\[\^[0-9]+\]`)

// footnoteReference reports whether a statement line carries a numbered
// citation reference.
func footnoteReference(line string) bool {
	return footnoteReferencePattern.MatchString(line)
}

// claimReference reports whether a line cites at least one CLM-* claim ID.
func claimReference(line string) bool {
	fields := strings.FieldsFunc(line, func(r rune) bool {
		return !(r >= 'A' && r <= 'Z') && !(r >= '0' && r <= '9') && r != '-'
	})
	for _, field := range fields {
		if strings.HasPrefix(field, "CLM-") && len(field) > len("CLM-") {
			return true
		}
	}
	return false
}

// NotApplicableWithReason reports a reason-bearing Not-Applicable exemption.
func NotApplicableWithReason(markdown string) bool {
	for _, line := range strings.Split(markdown, "\n") {
		_, reason, ok := strings.Cut(line, "Not-Applicable")
		if !ok {
			_, reason, ok = strings.Cut(line, "N/A")
			if !ok {
				continue
			}
		}
		reason = strings.TrimSpace(strings.TrimLeft(reason, ":]|-"))
		if len(strings.Fields(reason)) >= 3 {
			return true
		}
	}
	return false
}

// requiredDiagramChapters are the designated flow-heavy chapters that must
// carry at least one valid Mermaid diagram or a reason-bearing Not-Applicable.
var requiredDiagramChapters = map[string]bool{"02": true, "03": true, "05": true, "08": true}

// erDiagramChapters must additionally carry an entity-relationship diagram.
var erDiagramChapters = map[string]bool{"03": true}

// requiredPseudocodeChapters carry procedures by default (database routines,
// faults, or schedules in their catalogs) and must include a language-agnostic
// pseudocode block or a reason-bearing Not-Applicable.
var requiredPseudocodeChapters = map[string]bool{"03": true, "08": true}

// LintChapter applies the deterministic PI-11 substance gate to one chapter.
func LintChapter(subject, markdown string) []Finding {
	var findings []Finding
	lower := strings.ToLower(markdown)
	for _, section := range RequiredSections {
		if !strings.Contains(lower, section) {
			findings = append(findings, Finding{Code: "HBK-MISSING-SECTION", Subject: subject, Detail: "missing required section " + section})
		}
	}
	for _, finding := range ste.Lint(subject, proseView(markdown)) {
		findings = append(findings, Finding{Code: finding.Code, Subject: finding.Subject, Detail: finding.Detail})
	}
	for _, finding := range ValidateDiagrams(subject, markdown) {
		findings = append(findings, finding)
	}
	findings = append(findings, lintRequiredBlocks(subject, markdown)...)
	if count := AssertionCount(markdown); count < MinimumAssertions {
		if !NotApplicableWithReason(markdown) {
			findings = append(findings, Finding{
				Code: "HBK-SUBSTANCE", Subject: subject,
				Detail: fmt.Sprintf("%d claim-backed assertions, want at least %d or a reason-bearing Not-Applicable", count, MinimumAssertions),
			})
		}
	}
	return findings
}

// proseView removes mechanical reference lines (identity bindings) and fenced
// diagram/code blocks so the STE lint sees human-facing prose only.
func proseView(markdown string) string {
	var lines []string
	inFence := false
	for _, line := range strings.Split(markdown, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		if strings.HasPrefix(trimmed, "- **source-bindings:**") {
			continue
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// ReadabilityReview is the deterministic schema of a PI-12 readability CNF.
type ReadabilityReview struct {
	Value       string // Pass | Fail | Not-Applicable
	Scope       string
	Limitations string
	Timestamp   string
	HBKBindings []string
}

// ValidateReadabilityReview checks the PI-12 CNF schema. Pass and Fail require a
// real recorded read: non-empty scope, limitations, timestamp, and at least one
// exact reviewed HBK binding. Not-Applicable is the only case allowed to omit
// bindings, and it still needs scope.
func ValidateReadabilityReview(subject string, review ReadabilityReview) []Finding {
	var findings []Finding
	switch review.Value {
	case "Pass", "Fail":
		if strings.TrimSpace(review.Scope) == "" {
			findings = append(findings, Finding{Code: "CNF-SCOPE", Subject: subject, Detail: "readability review scope is empty"})
		}
		if strings.TrimSpace(review.Limitations) == "" {
			findings = append(findings, Finding{Code: "CNF-LIMITATIONS", Subject: subject, Detail: "readability review limitations are empty"})
		}
		if strings.TrimSpace(review.Timestamp) == "" {
			findings = append(findings, Finding{Code: "CNF-TIMESTAMP", Subject: subject, Detail: "readability review timestamp is empty"})
		}
		if len(review.HBKBindings) == 0 {
			findings = append(findings, Finding{Code: "CNF-BINDING", Subject: subject, Detail: "readability review binds no reviewed HBK payload"})
		}
	case "Not-Applicable":
		if strings.TrimSpace(review.Scope) == "" {
			findings = append(findings, Finding{Code: "CNF-SCOPE", Subject: subject, Detail: "Not-Applicable review still requires scope"})
		}
	default:
		findings = append(findings, Finding{Code: "CNF-VALUE", Subject: subject, Detail: fmt.Sprintf("invalid Handbook Readability Review value %q", review.Value)})
	}
	for i, binding := range review.HBKBindings {
		if strings.TrimSpace(binding) == "" {
			findings = append(findings, Finding{Code: "CNF-BINDING", Subject: subject, Detail: fmt.Sprintf("empty HBK binding at position %d", i+1)})
		}
	}
	return findings
}
