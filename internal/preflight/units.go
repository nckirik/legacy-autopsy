package preflight

import (
	"fmt"

	"github.com/nckirik/legacy-autopsy/internal/work"
)

// Directives are runtime-owned values the header must declare.
type Directives struct {
	ProtocolVersion string
	StorageBinding  string
	Iteration       string
	MaxDepth        int
}

// DefaultDirectives are the pilot preflight directives.
func DefaultDirectives() Directives {
	return Directives{
		ProtocolVersion: "4.2.1",
		StorageBinding:  "sqlite/1",
		Iteration:       "ALFA",
		MaxDepth:        3,
	}
}

// Input materializes the bounded preflight context: the deterministic scan plus
// the runtime directives the model must honor.
func (d Directives) Input(scan Scan) work.Input {
	content := string(scan.Markdown())
	content += "\n## Runtime directives (must be declared exactly)\n\n"
	content += fmt.Sprintf("- **Protocol Version:** %s\n", d.ProtocolVersion)
	content += fmt.Sprintf("- **Canonical Storage Binding:** %s\n", d.StorageBinding)
	content += fmt.Sprintf("- **Current Iteration:** %s\n", d.Iteration)
	content += fmt.Sprintf("- **Default Max Traversal Depth:** %d\n", d.MaxDepth)
	content += "\n## Persona discovery guidance\n\n"
	content += "- Identify distinct actors by entry path and visibility, not by business title: for example web UI, HTTP/API routes, background jobs, database/admin surfaces, operators, external systems.\n"
	content += "- One row per materially distinct actor; never a single catch-all row when the scan shows multiple entry surfaces.\n"
	content += "- Prefix must match [A-Z][A-Z0-9]{1,7}; slug is lowercase words joined by hyphens; directory must equal prefix-slug; Actor Class is human, service, or external; Registry State is one of Active, Conditional, Disabled, Historical, Retired.\n"
	content += "\n## Entry-cluster guidance\n\n"
	content += "- One cluster per distinct concrete entry surface (web routes, API routes, auth/SAML, DB/admin, scheduled jobs); cluster-id is a short uppercase token (for example WEB, API, AUTH, DB); cluster-title is human-readable.\n"
	content += "- Concrete Entry Coordinates must be exact paths or coordinates present in the scan, never wildcards.\n"
	content += "- Enum fields must use one of the listed options.\n"
	return work.Input{Name: "preflight-context.md", Bytes: []byte(content)}
}

// Plan builds the Preflight work units in dependency order: the `0A` header
// singleton, then the persona registry batch with IDs derived from the Prefix
// field.
func Plan(scan Scan, directives Directives) []work.Unit {
	input := directives.Input(scan)
	return []work.Unit{
		{
			ID: "U-PREFLIGHT-HEADER", Mode: "Preflight", RecordID: "0A",
			ResultType: "PREFLIGHT-HEADER", Inputs: []work.Input{input},
			Sections: []string{"preflight-registry"}, ReasoningEffort: "off",
		},
		{
			ID: "U-PREFLIGHT-PERSONAS", Mode: "Preflight",
			IDFromField: "prefix", ResultType: "PERSONA-REGISTRY",
			Reads:    []work.Ref{{FieldID: "PREFLIGHT-HEADER", RecordID: "0A"}},
			Inputs:   []work.Input{input},
			Sections: []string{"persona-definition"}, ReasoningEffort: "off",
		},
		{
			ID: "U-PREFLIGHT-CLUSTERS", Mode: "Preflight",
			IDFromField: "cluster-id", ResultType: "ENTRY-CLUSTER",
			Reads:    []work.Ref{{FieldID: "PREFLIGHT-HEADER", RecordID: "0A"}},
			Inputs:   []work.Input{input},
			Sections: []string{"preflight-registry"}, ReasoningEffort: "off",
		},
	}
}
