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
