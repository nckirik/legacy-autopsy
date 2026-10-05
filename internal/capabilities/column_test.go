package capabilities

import "testing"

func TestColumnName(t *testing.T) {
	cases := map[string]string{
		"Verification Vectors":    "verification-vectors",
		"Concrete Coordinate":     "concrete-coordinate",
		"SRC ID":                  "src-id",
		"row-kind":                "row-kind",
		"Exit E Candidate Report": "exit-e-candidate-report",
		"  Mixed...punctuation  ": "mixed-punctuation",
		"UPPER_case":              "upper-case",
		"already-kebab":           "already-kebab",
	}
	for input, want := range cases {
		if got := ColumnName(input); got != want {
			t.Errorf("ColumnName(%q) = %q, want %q", input, got, want)
		}
	}
}
