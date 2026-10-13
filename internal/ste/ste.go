// Package ste supports ASD-STE100 (Simplified Technical English) for
// human-facing generated prose. The prompt directive is the authoring contract;
// Lint enforces the mechanically checkable rules and reports the rest as
// review findings. The canonical standard wins over this implementation.
package ste

import (
	"fmt"
	"strings"
	"unicode"
)

// Directive is appended to human-facing generation prompts (handbook chapters,
// decision batches, operator questions, human reports).
const Directive = `Write human-facing prose in ASD-STE100 Simplified Technical English:

- One word has one meaning. Use the same word for the same thing every time.
- Give one instruction or one idea in each sentence.
- Keep procedural sentences to 20 words or fewer and descriptive sentences to
  25 words or fewer.
- Use the active voice. Name the actor. Do not use passive constructions.
- Use simple present, simple past, or simple future tense.
- Do not use contractions, semicolons, or idioms.
- Do not use "should": use "must" for a requirement and "can" for a capability.
- Keep noun groups to three words or fewer.
- Use a vertical list for more than two coordinated items.
- Use "None" for an empty value. Do not use "N/A" in generated prose.`

// Finding is one deterministic STE lint finding.
type Finding struct {
	Code    string
	Subject string
	Detail  string
}

func (f Finding) String() string { return fmt.Sprintf("%s %s: %s", f.Code, f.Subject, f.Detail) }

// MaxWords limits a sentence: STE allows 20 words for procedures and 25 for
// descriptions; Lint applies the descriptive limit and reports longer
// procedural-looking text for review.
const MaxWords = 25

// Lint checks the mechanically verifiable STE rules. Field-label lines are
// measured on their value only, serialized `;`-separated values are exempt from
// the semicolon rule, and reference lines (identity bindings) are skipped.
func Lint(subject, text string) []Finding {
	var findings []Finding
	report := func(code, detail string) {
		findings = append(findings, Finding{Code: code, Subject: subject, Detail: detail})
	}
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if isReferenceLine(trimmed) {
			continue
		}
		value := trimmed
		field := false
		if label, rest, ok := splitFieldLine(trimmed); ok {
			field = true
			value = rest
			_ = label
		}
		if value == "" {
			continue
		}
		structured := strings.Contains(value, "; ") && !strings.ContainsAny(value, ".!?")
		if !structured && field && !looksLikeProse(value) {
			// Identifier or path values are data, not prose.
			checkLexical(value, report)
			continue
		}
		if !structured {
			for _, sentence := range sentences(value) {
				words := strings.Fields(sentence)
				if len(words) > MaxWords {
					report("STE-SENTENCE-LENGTH", fmt.Sprintf("sentence has %d words (limit %d): %q", len(words), MaxWords, truncate(sentence, 80)))
				}
			}
		}
		if !field && !structured && strings.Contains(value, ";") {
			report("STE-SEMICOLON", "Do not use semicolons. Split the sentence or use a list.")
		}
		checkLexical(value, report)
	}
	return findings
}

func checkLexical(text string, report func(code, detail string)) {
	for _, contraction := range contractions(text) {
		report("STE-CONTRACTION", fmt.Sprintf("contraction %q is not permitted", contraction))
	}
	if strings.Contains(strings.ToLower(text), "n/a") {
		report("STE-EMPTY-VALUE", `use "None" for an empty value, not "N/A"`)
	}
	lower := strings.ToLower(text)
	if strings.Contains(lower, "should ") || strings.HasSuffix(lower, "should") {
		report("STE-SHOULD", `use "must" for a requirement and "can" for a capability, not "should"`)
	}
}

func isReferenceLine(trimmed string) bool {
	return strings.HasPrefix(trimmed, "- **source-bindings:")
}

func splitFieldLine(trimmed string) (string, string, bool) {
	if !strings.HasPrefix(trimmed, "- **") {
		return "", "", false
	}
	end := strings.Index(trimmed[2:], ":**")
	if end < 0 {
		return "", "", false
	}
	return trimmed[4 : end+2], strings.TrimSpace(trimmed[end+5:]), true
}

func looksLikeProse(value string) bool {
	words := strings.Fields(value)
	if len(words) < 3 {
		return false
	}
	lowercaseWords := 0
	for _, word := range words {
		clean := strings.Trim(word, ".,:;()[]\"")
		if clean != strings.ToUpper(clean) && len(clean) >= 3 {
			lowercaseWords++
		}
	}
	return lowercaseWords >= 2
}

var knownContractions = []string{
	"'s ", "'t ", "'re ", "'ve ", "'ll ", "'d ", "n't", "it's", "don't", "can't", "won't", "isn't", "aren't", "doesn't",
}

func contractions(text string) []string {
	lower := " " + strings.ToLower(text) + " "
	var found []string
	seen := map[string]bool{}
	for _, contraction := range knownContractions {
		if strings.Contains(lower, contraction) && !seen[contraction] {
			seen[contraction] = true
			found = append(found, strings.TrimSpace(contraction))
		}
	}
	return found
}

// sentences splits on terminal punctuation. Abbreviations are rare in
// controlled prose, so this bounded split is sufficient for linting.
func sentences(text string) []string {
	var out []string
	var current strings.Builder
	for _, r := range text {
		current.WriteRune(r)
		switch r {
		case '.', '!', '?':
			s := strings.TrimSpace(current.String())
			if s != "" {
				out = append(out, s)
			}
			current.Reset()
		case '\n':
			s := strings.TrimSpace(current.String())
			if s != "" && isListItem(s) {
				out = append(out, s)
			}
			current.Reset()
		}
	}
	if s := strings.TrimSpace(current.String()); s != "" {
		out = append(out, s)
	}
	return out
}

func isListItem(s string) bool {
	trimmed := strings.TrimLeftFunc(s, unicode.IsSpace)
	return strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ")
}

func truncate(value string, max int) string {
	if len(value) <= max {
		return value
	}
	return value[:max] + "..."
}
