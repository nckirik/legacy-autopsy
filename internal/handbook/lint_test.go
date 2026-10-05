package handbook

import (
	"fmt"
	"strings"
	"testing"
)

func chapter(assertions int) string {
	var b strings.Builder
	b.WriteString("# Overview\n\n## Audience Purpose\nFor maintainers.\n\n## Observed Legacy Behavior\n")
	for i := 0; i < assertions; i++ {
		fmt.Fprintf(&b, "- The scheduler retries jobs after failures (CLM-%04d).\n", i+1)
	}
	b.WriteString("\n## Target Decisions\nKeep retries but bound them.\n\n## Navigation\n- next: data model\n")
	return b.String()
}

func TestLintChapterSubstance(t *testing.T) {
	if findings := LintChapter("overview", chapter(8)); len(findings) != 0 {
		t.Fatalf("valid chapter rejected: %v", findings)
	}
	findings := LintChapter("overview", chapter(7))
	if len(findings) != 1 || findings[0].Code != "HBK-SUBSTANCE" {
		t.Fatalf("seven assertions not rejected: %v", findings)
	}
	missing := strings.Replace(chapter(8), "## Navigation", "## Notes", 1)
	findings = LintChapter("overview", missing)
	if len(findings) != 1 || findings[0].Code != "HBK-MISSING-SECTION" {
		t.Fatalf("missing section not detected: %v", findings)
	}
	empty := chapter(8) + "\nNot-Applicable: this chapter class is legitimately empty because no records exist.\n"
	empty = strings.Replace(empty, "CLM-0001", "claim", 1)
	for i := 2; i <= 8; i++ {
		empty = strings.Replace(empty, fmt.Sprintf("CLM-%04d", i), "claim", 1)
	}
	if findings := LintChapter("overview", empty); len(findings) != 0 {
		t.Fatalf("reason-bearing N/A rejected: %v", findings)
	}
	noReason := chapter(3) + "\nNot-Applicable\n"
	if findings := LintChapter("overview", noReason); len(findings) != 1 {
		t.Fatalf("reasonless N/A accepted: %v", findings)
	}
}

func TestAssertionCountSkipsHeadingsAndNav(t *testing.T) {
	markdown := "## CLM-0001 heading claim\n- CLM-0002 bullet assertion\n- [ ] task CLM-0003\nReal statement (CLM-0004).\nCLM-0005 and CLM-0006 on one line.\n"
	if count := AssertionCount(markdown); count != 3 {
		t.Fatalf("count = %d, want 3", count)
	}
}

func TestValidateReadabilityReview(t *testing.T) {
	valid := ReadabilityReview{Value: "Pass", Scope: "all candidate chapters", Limitations: "static-only scope", Timestamp: "2026-10-05T00:00:00Z", HBKBindings: []string{"HBK-1"}}
	if findings := ValidateReadabilityReview("hbk", valid); len(findings) != 0 {
		t.Fatalf("valid review rejected: %v", findings)
	}
	for name, review := range map[string]ReadabilityReview{
		"empty scope":      {Value: "Pass", Limitations: "x", Timestamp: "t", HBKBindings: []string{"HBK-1"}},
		"no timestamp":     {Value: "Fail", Scope: "x", Limitations: "y", HBKBindings: []string{"HBK-1"}},
		"no bindings":      {Value: "Pass", Scope: "x", Limitations: "y", Timestamp: "t"},
		"bad value":        {Value: "Maybe", Scope: "x"},
		"na without scope": {Value: "Not-Applicable"},
	} {
		if findings := ValidateReadabilityReview("hbk", review); len(findings) == 0 {
			t.Fatalf("%s: invalid review accepted", name)
		}
	}
	na := ReadabilityReview{Value: "Not-Applicable", Scope: "non-handbook CNF"}
	if findings := ValidateReadabilityReview("hbk", na); len(findings) != 0 {
		t.Fatalf("valid Not-Applicable rejected: %v", findings)
	}
}
