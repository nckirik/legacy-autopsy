package parity

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// revisionMarkers are distinctive phrases introduced by accepted protocol
// revisions. They must appear in the generated edition so a revision cannot be
// silently lost during regeneration.
var revisionMarkers = []string{
	"Fail-inclusive Unknown applies at absent authority",
	"Implementation tooling - runtimes, validators",
	"Commented-out code, commented call sites",
	"Dormant and disabled units follow the extraction-first rule",
	"Commented-out or disabled wiring is not a zero-caller proof",
	"ITERATION-PERSONA-POV-SEQUENCE",
	"Bootstrap: the first Preflight invocation",
	"Closure and materialization: the mandatory read set",
	"Blocker classification: before any Human Hatch escalation",
	"Static-only runs are eligible through both gates",
	"at least eight",
	"Handbook Readability Review",
	"Static-only eligibility: no condition in this section",
	"relay-only attestation",
}

func TestGeneratedEditionCarriesRevisionMarkers(t *testing.T) {
	data, err := os.ReadFile(filepath.Join(root(t), "protocol.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	for _, marker := range revisionMarkers {
		if !strings.Contains(text, marker) {
			t.Errorf("generated protocol.md is missing revision marker %q", marker)
		}
	}
}
