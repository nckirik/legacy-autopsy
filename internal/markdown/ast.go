// Package markdown provides the small structural Markdown model used by the
// bootstrap validators. It deliberately exposes structure instead of treating
// protocol documents as unanchored text.
package markdown

import (
	"fmt"
	"strings"
)

type Kind string

const (
	Heading   Kind = "heading"
	Field     Kind = "field"
	CodeBlock Kind = "code-block"
	Comment   Kind = "comment"
	TableRow  Kind = "table-row"
	Paragraph Kind = "paragraph"
	Blank     Kind = "blank"
)

type Span struct {
	StartLine int
	EndLine   int
}

type Node struct {
	Kind  Kind
	Level int
	Text  string
	Label string
	Value string
	Info  string
	Raw   string
	Span  Span
}

type Document struct {
	Source string
	Nodes  []Node
}

// Parse builds a bounded structural representation of the Markdown constructs
// currently needed by the protocol model. It does not use regular expressions.
func Parse(source string) (*Document, error) {
	source = strings.ReplaceAll(strings.ReplaceAll(source, "\r\n", "\n"), "\r", "\n")
	lines := strings.Split(source, "\n")
	doc := &Document{Source: source}

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		lineNo := i + 1

		if marker, info, ok := fenceStart(trimmed); ok {
			start := i
			i++
			for ; i < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[i]), marker); i++ {
			}
			if i >= len(lines) {
				return nil, fmt.Errorf("line %d: unclosed fenced code block", lineNo)
			}
			raw := strings.Join(lines[start+1:i], "\n")
			doc.Nodes = append(doc.Nodes, Node{Kind: CodeBlock, Info: info, Raw: raw, Span: Span{StartLine: lineNo, EndLine: i + 1}})
			continue
		}

		if level, text, ok := heading(line); ok {
			doc.Nodes = append(doc.Nodes, Node{Kind: Heading, Level: level, Text: text, Span: Span{StartLine: lineNo, EndLine: lineNo}})
			continue
		}
		if label, value, ok := field(trimmed); ok {
			doc.Nodes = append(doc.Nodes, Node{Kind: Field, Label: label, Value: value, Raw: line, Span: Span{StartLine: lineNo, EndLine: lineNo}})
			continue
		}
		if strings.HasPrefix(trimmed, "<!--") {
			doc.Nodes = append(doc.Nodes, Node{Kind: Comment, Text: trimmed, Raw: line, Span: Span{StartLine: lineNo, EndLine: lineNo}})
			continue
		}
		if strings.HasPrefix(trimmed, "|") && strings.HasSuffix(trimmed, "|") {
			doc.Nodes = append(doc.Nodes, Node{Kind: TableRow, Raw: line, Span: Span{StartLine: lineNo, EndLine: lineNo}})
			continue
		}
		if trimmed == "" {
			doc.Nodes = append(doc.Nodes, Node{Kind: Blank, Raw: line, Span: Span{StartLine: lineNo, EndLine: lineNo}})
			continue
		}
		doc.Nodes = append(doc.Nodes, Node{Kind: Paragraph, Text: trimmed, Raw: line, Span: Span{StartLine: lineNo, EndLine: lineNo}})
	}
	return doc, nil
}

func fenceStart(line string) (marker, info string, ok bool) {
	for _, candidate := range []string{"```", "~~~"} {
		if strings.HasPrefix(line, candidate) {
			return candidate, strings.TrimSpace(strings.TrimPrefix(line, candidate)), true
		}
	}
	return "", "", false
}

func heading(line string) (int, string, bool) {
	if len(line) < 3 || line[0] != '#' {
		return 0, "", false
	}
	level := 0
	for level < len(line) && line[level] == '#' {
		level++
	}
	if level > 6 || level >= len(line) || line[level] != ' ' {
		return 0, "", false
	}
	text := strings.TrimSpace(line[level+1:])
	return level, text, text != ""
}

func field(line string) (string, string, bool) {
	line = strings.TrimSpace(strings.TrimPrefix(line, "- "))
	if !strings.HasPrefix(line, "**") {
		return "", "", false
	}
	end := strings.Index(line[2:], ":**")
	if end < 0 {
		return "", "", false
	}
	end += 2
	label := strings.TrimSpace(line[2:end])
	value := strings.TrimSpace(line[end+3:])
	return label, value, label != ""
}

func (d *Document) HeadingCount(title string) int {
	count := 0
	for _, node := range d.Nodes {
		if node.Kind == Heading && node.Text == title {
			count++
		}
	}
	return count
}

func (d *Document) Heading(title string) (Node, bool) {
	for _, node := range d.Nodes {
		if node.Kind == Heading && node.Text == title {
			return node, true
		}
	}
	return Node{}, false
}

func (d *Document) Section(title string) ([]Node, error) {
	start := -1
	level := 0
	for i, node := range d.Nodes {
		if node.Kind == Heading && node.Text == title {
			start, level = i, node.Level
			break
		}
	}
	if start < 0 {
		return nil, fmt.Errorf("heading %q not found", title)
	}
	end := len(d.Nodes)
	for i := start + 1; i < len(d.Nodes); i++ {
		if d.Nodes[i].Kind == Heading && d.Nodes[i].Level <= level {
			end = i
			break
		}
	}
	return d.Nodes[start:end], nil
}

func (d *Document) SectionText(title string) (string, error) {
	nodes, err := d.Section(title)
	if err != nil {
		return "", err
	}
	start := nodes[0].Span.StartLine
	end := nodes[len(nodes)-1].Span.EndLine
	lines := strings.Split(d.Source, "\n")
	if end > len(lines) {
		end = len(lines)
	}
	return strings.Join(lines[start-1:end], "\n") + "\n", nil
}
