package handbook

import (
	"fmt"
	"strings"

	md "github.com/nckirik/legacy-autopsy/internal/markdown"
)

// diagramClass is the validator class of a declared Mermaid diagram head.
type diagramClass string

const (
	diagramFlowchart diagramClass = "flowchart"
	diagramState     diagramClass = "state"
	diagramER        diagramClass = "er"
	diagramSequence  diagramClass = "sequence"
)

// supportedDiagramHeads maps accepted Mermaid heads to their validator class.
// Unknown heads fail closed: an incapable parser reports a finding instead of
// silently accepting an unchecked diagram.
var supportedDiagramHeads = map[string]diagramClass{
	"flowchart":       diagramFlowchart,
	"graph":           diagramFlowchart,
	"statediagram":    diagramState,
	"statediagram-v2": diagramState,
	"erdiagram":       diagramER,
	"sequencediagram": diagramSequence,
}

// ValidateDiagrams checks every fenced `mermaid` block in a chapter with the
// deterministic diagram parser: the head must be declared and supported, the
// body must be non-empty, the diagram's declared structure must be present,
// and delimiters must balance. Other fenced blocks are ignored.
func ValidateDiagrams(subject, markdown string) []Finding {
	doc, err := md.Parse(markdown)
	if err != nil {
		return []Finding{{Code: "HBK-DIAGRAM-PARSE", Subject: subject, Detail: err.Error()}}
	}
	var findings []Finding
	index := 0
	for _, node := range doc.Nodes {
		if node.Kind != md.CodeBlock {
			continue
		}
		if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(node.Info)), "mermaid") {
			continue
		}
		index++
		findings = append(findings, validateDiagram(subject, index, node.Raw)...)
	}
	return findings
}

func validateDiagram(subject string, index int, raw string) []Finding {
	label := fmt.Sprintf("diagram %d", index)
	head := ""
	var body []string
	for _, line := range strings.Split(raw, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "%%") {
			continue
		}
		if head == "" {
			head = strings.ToLower(strings.Fields(trimmed)[0])
			continue
		}
		body = append(body, line)
	}
	if head == "" {
		return []Finding{{Code: "HBK-DIAGRAM-EMPTY", Subject: subject, Detail: label + " has no diagram head"}}
	}
	class, ok := supportedDiagramHeads[head]
	if !ok {
		return []Finding{{Code: "HBK-DIAGRAM-UNSUPPORTED-HEAD", Subject: subject, Detail: fmt.Sprintf("%s declares unsupported diagram head %q", label, head)}}
	}
	if len(body) == 0 {
		return []Finding{{Code: "HBK-DIAGRAM-EMPTY", Subject: subject, Detail: label + " has no diagram body"}}
	}
	if code, detail := unbalancedDelimiters(label, body); code != "" {
		return []Finding{{Code: code, Subject: subject, Detail: detail}}
	}
	joined := strings.Join(body, "\n")
	switch class {
	case diagramFlowchart:
		if !containsAny(joined, "-->", "---", "-.->", "==>", "--x", "--o", "<-->", "o--o", "x--x") {
			return []Finding{{Code: "HBK-DIAGRAM-NO-STRUCTURE", Subject: subject, Detail: label + " flowchart declares no links between nodes"}}
		}
	case diagramState:
		if !containsAny(joined, "-->") {
			return []Finding{{Code: "HBK-DIAGRAM-NO-STRUCTURE", Subject: subject, Detail: label + " state diagram declares no transitions"}}
		}
	case diagramER:
		if !containsAny(joined, "--") {
			return []Finding{{Code: "HBK-DIAGRAM-NO-STRUCTURE", Subject: subject, Detail: label + " entity-relationship diagram declares no relationships"}}
		}
	case diagramSequence:
		if !containsAny(joined, "->>", "-->>", "->", "-->", "-x", "--x") {
			return []Finding{{Code: "HBK-DIAGRAM-NO-STRUCTURE", Subject: subject, Detail: label + " sequence diagram declares no messages"}}
		}
	}
	return nil
}

// lintRequiredBlocks enforces the chapter content requirements: required
// diagrams (and the entity-relationship diagram) for flow-heavy chapters, and
// pseudocode for procedure-bearing chapters, each with a reason-bearing
// Not-Applicable escape.
func lintRequiredBlocks(subject, markdown string) []Finding {
	var findings []Finding
	if requiredDiagramChapters[subject] {
		classes := diagramClasses(markdown)
		if len(classes) == 0 && !NotApplicableWithReason(markdown) {
			findings = append(findings, Finding{Code: "HBK-DIAGRAM-REQUIRED", Subject: subject, Detail: "flow-heavy chapter needs a valid Mermaid diagram or a reason-bearing Not-Applicable"})
		}
		if erDiagramChapters[subject] && !containsClass(classes, diagramER) && !NotApplicableWithReason(markdown) {
			findings = append(findings, Finding{Code: "HBK-ER-DIAGRAM", Subject: subject, Detail: "data chapter needs an entity-relationship diagram of the entity/table inventory or a reason-bearing Not-Applicable"})
		}
	}
	if requiredPseudocodeChapters[subject] && !hasPseudocodeBlock(markdown) && !NotApplicableWithReason(markdown) {
		findings = append(findings, Finding{Code: "HBK-PSEUDOCODE", Subject: subject, Detail: "procedure-bearing chapter needs a language-agnostic pseudocode block or a reason-bearing Not-Applicable"})
	}
	return findings
}

// diagramClasses returns the validator classes of every supported Mermaid block
// in a chapter.
func diagramClasses(markdown string) []diagramClass {
	doc, err := md.Parse(markdown)
	if err != nil {
		return nil
	}
	var classes []diagramClass
	for _, node := range doc.Nodes {
		if node.Kind != md.CodeBlock || !strings.HasPrefix(strings.ToLower(strings.TrimSpace(node.Info)), "mermaid") {
			continue
		}
		for _, line := range strings.Split(node.Raw, "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "%%") {
				continue
			}
			if class, ok := supportedDiagramHeads[strings.ToLower(strings.Fields(trimmed)[0])]; ok {
				classes = append(classes, class)
			}
			break
		}
	}
	return classes
}

// hasPseudocodeBlock reports whether a chapter carries a non-empty
// language-agnostic pseudocode block (fenced `text` or `pseudocode`).
func hasPseudocodeBlock(markdown string) bool {
	doc, err := md.Parse(markdown)
	if err != nil {
		return false
	}
	for _, node := range doc.Nodes {
		if node.Kind != md.CodeBlock || strings.TrimSpace(node.Raw) == "" {
			continue
		}
		info := strings.ToLower(strings.TrimSpace(node.Info))
		if strings.HasPrefix(info, "text") || strings.HasPrefix(info, "pseudocode") {
			return true
		}
	}
	return false
}

func containsClass(classes []diagramClass, want diagramClass) bool {
	for _, class := range classes {
		if class == want {
			return true
		}
	}
	return false
}

func unbalancedDelimiters(label string, lines []string) (string, string) {
	counts := map[string]int{}
	for _, line := range lines {
		// Entity-relationship cardinality markers use brace-like tokens
		// (`o{`, `|{`, `}o`, `}|`) that are not block delimiters.
		normalized := line
		for _, marker := range []string{"o{", "|{", "}o", "}|"} {
			normalized = strings.ReplaceAll(normalized, marker, "")
		}
		counts["["] += strings.Count(normalized, "[")
		counts["]"] += strings.Count(normalized, "]")
		counts["{"] += strings.Count(normalized, "{")
		counts["}"] += strings.Count(normalized, "}")
		counts["("] += strings.Count(normalized, "(")
		counts[")"] += strings.Count(normalized, ")")
		if strings.Count(line, "\"")%2 != 0 {
			return "HBK-DIAGRAM-UNBALANCED", fmt.Sprintf("%s has an unterminated quote on %q", label, strings.TrimSpace(line))
		}
	}
	for _, pair := range [][2]string{{"[", "]"}, {"{", "}"}, {"(", ")"}} {
		if counts[pair[0]] != counts[pair[1]] {
			return "HBK-DIAGRAM-UNBALANCED", fmt.Sprintf("%s has unbalanced %s%s delimiters", label, pair[0], pair[1])
		}
	}
	return "", ""
}

func containsAny(value string, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(value, needle) {
			return true
		}
	}
	return false
}
