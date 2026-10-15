package handbook

import (
	"strings"
	"testing"
)

func TestValidateDiagramsAcceptsSupportedHeads(t *testing.T) {
	cases := map[string]string{
		"flowchart": "```mermaid\nflowchart TD\n  A[Request] --> B{Review}\n  B -->|approve| C[Ready]\n  B -->|reject| A\n```\n",
		"state":     "```mermaid\nstateDiagram-v2\n  [*] --> Open\n  Open --> Pending\n  Pending --> [*]\n```\n",
		"er":        "```mermaid\nerDiagram\n  MALZEMELER ||--o{ BIRIMLER : has\n  MALZEMELER {\n    int id\n  }\n```\n",
		"sequence":  "```mermaid\nsequenceDiagram\n  Client->>Server: request\n  Server-->>Client: response\n```\n",
	}
	for name, markdown := range cases {
		if findings := ValidateDiagrams("01", markdown); len(findings) != 0 {
			t.Fatalf("%s diagram rejected: %v", name, findings)
		}
	}
}

func TestValidateDiagramsRejectsBrokenBlocks(t *testing.T) {
	cases := map[string]struct {
		markdown string
		code     string
	}{
		"unknown head":       {"```mermaid\nquadrantChart\n  x-axis A --> B\n```\n", "HBK-DIAGRAM-UNSUPPORTED-HEAD"},
		"empty":              {"```mermaid\n```\n", "HBK-DIAGRAM-EMPTY"},
		"head only":          {"```mermaid\nflowchart TD\n```\n", "HBK-DIAGRAM-EMPTY"},
		"no links":           {"```mermaid\nflowchart TD\n  A[Request]\n  B[Review]\n```\n", "HBK-DIAGRAM-NO-STRUCTURE"},
		"unbalanced":         {"```mermaid\nflowchart TD\n  A[Request --> B\n```\n", "HBK-DIAGRAM-UNBALANCED"},
		"unterminated quote": {"```mermaid\nflowchart TD\n  A[\"Request] --> B\n```\n", "HBK-DIAGRAM-UNBALANCED"},
	}
	for name, tc := range cases {
		findings := ValidateDiagrams("01", tc.markdown)
		found := false
		for _, finding := range findings {
			if finding.Code == tc.code {
				found = true
			}
		}
		if !found {
			t.Fatalf("%s: findings = %v, want %s", name, findings, tc.code)
		}
	}
}

func TestValidateDiagramsIgnoresOtherFences(t *testing.T) {
	markdown := "```text\nflowchart without mermaid info string\n```\n\n```mermaid\nflowchart LR\n  A --> B\n```\n"
	findings := ValidateDiagrams("01", markdown)
	if len(findings) != 0 {
		t.Fatalf("findings = %v", findings)
	}
}

func TestProseViewAndAssertionsSkipFencedBlocks(t *testing.T) {
	markdown := "The page checks the role first [^1].\n\n```mermaid\nflowchart TD\n  A[a;b] --> B[it's done]\n```\n\nThe lock serializes access [^2].\n"
	if got := AssertionCount(markdown); got != 2 {
		t.Fatalf("assertion count = %d, want 2", got)
	}
	prose := proseView(markdown)
	if strings.Contains(prose, "it's") || strings.Contains(prose, "a;b") {
		t.Fatalf("prose view still contains diagram text:\n%s", prose)
	}
	if !strings.Contains(prose, "The page checks the role first") {
		t.Fatalf("prose view lost narrative:\n%s", prose)
	}
}

func TestLintRequiredBlocks(t *testing.T) {
	flowchart := "```mermaid\nflowchart TD\n  A[Client] --> B[Server]\n```\n"
	er := "```mermaid\nerDiagram\n  A ||--o{ B : has\n```\n"
	pseudocode := "```text\n1. acquire the lock for the unit\n2. check conflicting reservations\n3. release the lock\n```\n"
	cases := []struct {
		name     string
		subject  string
		markdown string
		want     string
	}{
		{"02 without diagram", "02", "Prose only.\n", "HBK-DIAGRAM-REQUIRED"},
		{"03 without er", "03", flowchart, "HBK-ER-DIAGRAM"},
		{"03 without pseudocode", "03", flowchart + er, "HBK-PSEUDOCODE"},
		{"08 without pseudocode", "08", flowchart, "HBK-PSEUDOCODE"},
	}
	for _, tc := range cases {
		findings := LintChapter(tc.subject, tc.markdown)
		found := false
		for _, finding := range findings {
			if finding.Code == tc.want {
				found = true
			}
		}
		if !found {
			t.Fatalf("%s: findings = %v, want %s", tc.name, findings, tc.want)
		}
	}
	complete := flowchart + er + pseudocode
	for _, finding := range LintChapter("03", complete) {
		if strings.HasPrefix(finding.Code, "HBK-DIAGRAM") || finding.Code == "HBK-ER-DIAGRAM" || finding.Code == "HBK-PSEUDOCODE" {
			t.Fatalf("complete chapter flagged: %v", finding)
		}
	}
	escape := "Not-Applicable: this chapter class is empty because no records exist and no evidence supports one.\n"
	for _, finding := range LintChapter("08", escape) {
		if strings.HasPrefix(finding.Code, "HBK-DIAGRAM") || finding.Code == "HBK-PSEUDOCODE" {
			t.Fatalf("reason-bearing Not-Applicable rejected: %v", finding)
		}
	}
}
