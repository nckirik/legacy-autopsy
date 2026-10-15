package discovery

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// personaPattern is the persona prefix grammar (§2.3); reserved tokens never
// bind a pointer row.
var personaPattern = regexp.MustCompile(`^[A-Z][A-Z0-9]{1,7}$`)

var reservedPersonaTokens = map[string]bool{
	"SHARED": true, "UNKNOWN": true, "PENDING": true, "TBD": true,
}

// Deterministic excerpt bounds for atomic-component authoring.
const (
	maxExcerptLines = 41
	maxExcerptBytes = 6000
)

// povFilePath returns the active POV file of a unit: the persona-local ledger
// when a persona directory is bound, otherwise the shared POV ledger.
func povFilePath(directory, pov string) string {
	file := povFiles[pov]
	if file == "" {
		return "None"
	}
	if directory == "" {
		return "personas/_shared/" + file
	}
	return "personas/" + directory + "/" + file
}

// povFiles maps the five isolated POVs to their canonically numbered shared
// POV files (protocol/foundations.cdl forensic plane layout).
var povFiles = map[string]string{
	"POV-1": "01-FRONTEND-POV.md",
	"POV-2": "02-BACKEND-POV.md",
	"POV-3": "03-BACKGROUND-POV.md",
	"POV-4": "04-DATA-POV.md",
	"POV-5": "05-WIRING-POV.md",
}

// ComponentPlan builds bounded ATOMIC-COMPONENT authoring units from hit
// evidence. Every hit must already have its committed SRC inventory record; the
// unit's `source-ids` binding anchors CMP identity in that source unit. One
// record per hit is required (§6.1 atomic-unit specificity).
func ComponentPlan(root string, hits []Hit, clusters []Cluster, directives Directives) ([]work.Unit, error) {
	batchSize := directives.BatchSize
	if batchSize <= 0 {
		batchSize = 25
	}
	type group struct {
		cluster Cluster
		bound   bool
		pov     string
		hits    []Hit
	}
	var groups []*group
	index := map[string]*group{}
	for _, hit := range hits {
		pov := hitPOV(hit.Kind)
		if pov == "" {
			continue
		}
		cluster, bound := ClusterForHit(hit, clusters)
		key := pov + "|unbound"
		if bound {
			key = pov + "|" + cluster.ID
		}
		item, ok := index[key]
		if !ok {
			item = &group{cluster: cluster, bound: bound, pov: pov}
			index[key] = item
			groups = append(groups, item)
		}
		item.hits = append(item.hits, hit)
	}
	cache := map[string][]string{}
	var units []work.Unit
	for _, item := range groups {
		persona, directory := "", ""
		if item.bound {
			persona, directory, _ = item.cluster.PersonaBinding(directives.PersonaDirectories, directives.SharedEntries)
		}
		clusterToken := "UNBOUND"
		if item.bound {
			clusterToken = item.cluster.ID
		}
		for start := 0; start < len(item.hits); start += batchSize {
			end := start + batchSize
			if end > len(item.hits) {
				end = len(item.hits)
			}
			batch := item.hits[start:end]
			context, err := componentContext(root, cache, batch, item.cluster, item.bound, persona, directory, item.pov, directives)
			if err != nil {
				return nil, err
			}
			units = append(units, work.Unit{
				ID:              fmt.Sprintf("U-DISCOVERY-CMP-%s-%s-%03d", clusterToken, item.pov, start/batchSize+1),
				Mode:            "Discovery",
				Scope:           directives.Namespace,
				Persona:         persona,
				POV:             item.pov,
				IDFromBindings:  "source-ids",
				IDType:          "CMP",
				Namespace:       directives.Namespace,
				ExpectedRecords: len(batch),
				ResultType:      "ATOMIC-COMPONENT",
				RequiresReview:  true,
				Reads:           []work.Ref{{FieldID: "PREFLIGHT-HEADER", RecordID: "0A"}},
				Inputs:          []work.Input{{Name: "discovery-components-context.md", Bytes: context}},
				Sections:        []string{"atomic-unit-rule", "atomic-component-schema", "dual-entry-discovery", "five-isolated-povs", "traversal-tracks"},
				ReasoningEffort: "off",
			})
		}
	}
	return units, nil
}

func componentContext(root string, cache map[string][]string, batch []Hit, cluster Cluster, bound bool, persona, directory, pov string, directives Directives) ([]byte, error) {
	var b strings.Builder
	b.WriteString("# Atomic component authoring batch (runtime-generated, bounded)\n\n")
	b.WriteString("Author exactly one ATOMIC-COMPONENT record per hit below, in order. One record MUST NOT aggregate unrelated units.\n\n")
	for _, hit := range batch {
		srcID, err := capabilities.SourceID(directives.Namespace, string(hit.Kind), hit.Coordinate)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(&b, "## Hit %s\n\n", hit.Coordinate)
		fmt.Fprintf(&b, "- Kind: %s\n- File: %s\n- Line: %d\n- Source unit: %s\n- Evidence fingerprint: %s\n\n", hit.Kind, hit.File, hit.Line, srcID, hit.Fingerprint)
		excerpt, err := excerptFor(root, cache, hit)
		if err != nil {
			return nil, err
		}
		b.WriteString("Source excerpt (bounded window around the hit):\n\n```text\n")
		b.WriteString(excerpt)
		b.WriteString("\n```\n\n")
	}
	b.WriteString("## Strict single-scope bindings (§8.2)\n\n")
	fmt.Fprintf(&b, "- **Persona:** %s\n", orNone(persona))
	fmt.Fprintf(&b, "- **Persona Directory:** %s\n", orNone(directory))
	fmt.Fprintf(&b, "- **Active POV File:** %s\n", povFilePath(directory, pov))
	if bound {
		fmt.Fprintf(&b, "- **Entry Cluster:** %s (%s)\n", cluster.ID, cluster.Title)
		fmt.Fprintf(&b, "- **Traversal Track:** %s\n", cluster.Track)
	} else {
		b.WriteString("- **Entry Cluster:** None (unresolved hit binding)\n- **Traversal Track:** None\n")
	}
	fmt.Fprintf(&b, "- **Active POV:** %s\n\n", pov)
	b.WriteString("## Runtime directives (must be used exactly)\n\n")
	fmt.Fprintf(&b, "- **System Namespace:** %s\n", directives.Namespace)
	fmt.Fprintf(&b, "- **Environment/Snapshot:** %s\n", directives.Snapshot)
	b.WriteString("- `record-kind` is `ATOMIC`.\n")
	b.WriteString("- `source-ids` lists exactly the hit's Source unit token, nothing else.\n")
	b.WriteString("- `file-path-line-range` is `<file>:<line>` from the hit row.\n")
	b.WriteString("- `virtual-coordinate` is the hit's coordinate copied exactly.\n")
	fmt.Fprintf(&b, "- `pov-owner` is `%s`.\n", pov)
	b.WriteString("- `staging-status` is `[S-DISCOVERED]` (dual-entry discovery into the persona staging buffer).\n")
	b.WriteString("- `coverage-status` is `[C-GAP]` for newly discovered unswept units.\n")
	b.WriteString("- `material-claims` and `synthesis-references` are `None` (reconciliation populates them).\n")
	b.WriteString("- `block-confidence-rank` is `Low` for a single lexical hit; raise it only with direct evidence.\n")
	b.WriteString("- `block-evidence-profile` states the single [E-DIRECT] hit with its coordinate and source unit.\n")
	b.WriteString("- Where the excerpt does not support a field, use `Unknown` (fail-inclusive) or `None`; never invent facts.\n")
	if bound {
		fmt.Fprintf(&b, "- `traversal-tracks` is `%s`.\n", cluster.Track)
	} else {
		b.WriteString("- `traversal-tracks` is `Normal` only when the evidence shows it; otherwise `Unknown`.\n")
	}
	if persona != "" {
		fmt.Fprintf(&b, "- `personas` is `%s`.\n", persona)
	} else {
		b.WriteString("- `personas` is `None` (no resolvable persona binding for this unit).\n")
	}
	return []byte(b.String()), nil
}

// excerptFor returns a bounded line window around the hit, cached per file.
func excerptFor(root string, cache map[string][]string, hit Hit) (string, error) {
	lines, ok := cache[hit.File]
	if !ok {
		full := filepath.Join(root, filepath.FromSlash(hit.File))
		file, err := os.Open(full)
		if err != nil {
			return "", fmt.Errorf("discovery: component excerpt %s: %w", hit.File, err)
		}
		scanner := bufio.NewScanner(file)
		scanner.Buffer(make([]byte, 64*1024), maxHitLineBytes+1)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		file.Close()
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("discovery: component excerpt %s: %w", hit.File, err)
		}
		cache[hit.File] = lines
	}
	start := hit.Line - maxExcerptLines/2
	if start < 1 {
		start = 1
	}
	end := start + maxExcerptLines - 1
	if end > len(lines) {
		end = len(lines)
	}
	var b strings.Builder
	for i := start; i <= end; i++ {
		fmt.Fprintf(&b, "%d: %s\n", i, lines[i-1])
		if b.Len() > maxExcerptBytes {
			break
		}
	}
	return strings.TrimRight(b.String(), "\n"), nil
}

// ComponentRef is one committed ATOMIC-COMPONENT record in bounded form for the
// deterministic shared-reference materialization.
type ComponentRef struct {
	ID                 string
	VirtualCoordinate  string
	Persona            string
	POV                string
	Track              string
	CapabilityState    string
	StagingStatus      string
	CreatingInvocation string
}

// LoadComponents reads every committed atomic component with the fields the
// shared-reference rows need, sorted by ID.
func LoadComponents(db string, doc *cdl.EIRDoc) ([]ComponentRef, error) {
	names, _, _, err := work.FieldContract(doc, "ATOMIC-COMPONENT")
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, "ATOMIC-COMPONENT")
	if err != nil {
		return nil, err
	}
	sort.Strings(ids)
	var components []ComponentRef
	for _, id := range ids {
		payload, err := store.Get(db, doc, "ATOMIC-COMPONENT", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		persona := ""
		for _, token := range splitTokens(values["personas"]) {
			upper := strings.ToUpper(token)
			if personaPattern.MatchString(upper) && !reservedPersonaTokens[upper] {
				persona = upper
				break
			}
		}
		track := ""
		for _, token := range splitTokens(values["traversal-tracks"]) {
			switch token {
			case "Normal", "Shadow/Conditional", "Disabled":
				track = token
			}
			if track != "" {
				break
			}
		}
		components = append(components, ComponentRef{
			ID:                id,
			VirtualCoordinate: strings.Join(strings.Fields(values["virtual-coordinate"]), " "),
			Persona:           persona,
			POV:               strings.TrimSpace(values["pov-owner"]),
			Track:             track,
			CapabilityState:   strings.TrimSpace(values["capability-state"]),
			StagingStatus:     strings.TrimSpace(values["staging-status"]),
		})
	}
	return components, nil
}

// LoadCreatingInvocations maps each committed CMP ID to the invocation ID that
// first recorded it in `0G`, so shared-reference rows carry the discovering
// invocation. Unknown components stay empty and are rendered as None.
func LoadCreatingInvocations(db string, doc *cdl.EIRDoc) (map[string]string, error) {
	names, _, _, err := work.FieldContract(doc, "INVOCATION-LOG")
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, "INVOCATION-LOG")
	if err != nil {
		return nil, err
	}
	sort.Strings(ids)
	creating := map[string]string{}
	for _, id := range ids {
		payload, err := store.Get(db, doc, "INVOCATION-LOG", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		for _, token := range splitTokens(values["src-cmp-clm-ids-updated"]) {
			if !strings.HasPrefix(token, "CMP-") {
				continue
			}
			if _, ok := creating[token]; !ok {
				creating[token] = id
			}
		}
	}
	return creating, nil
}

// SharedReferencePlan builds one mechanical unit per (persona, POV) group. Each
// unit materializes the dual-entry discovery pointer rows
// (§6.2) for the persona's discovered components: the canonical location in the
// shared POV file plus the local trigger and staging status. Components with no
// resolvable persona binding produce no pointer; they are reported instead.
func SharedReferencePlan(components []ComponentRef) []work.Unit {
	type group struct {
		persona string
		pov     string
	}
	var order []group
	index := map[group][]ComponentRef{}
	for _, component := range components {
		if component.Persona == "" || povFiles[component.POV] == "" {
			continue
		}
		key := group{persona: component.Persona, pov: component.POV}
		if _, ok := index[key]; !ok {
			order = append(order, key)
		}
		index[key] = append(index[key], component)
	}
	var units []work.Unit
	for _, key := range order {
		units = append(units, work.Unit{
			ID:              fmt.Sprintf("U-DISCOVERY-SHARED-REFS-%s-%s", key.persona, key.pov),
			Mode:            "Discovery",
			Persona:         key.persona,
			POV:             key.pov,
			IDFromFields:    []string{"canonical-location", "local-trigger-src-coordinate", "local-invocation-signature"},
			ResultType:      "SHARED-REFERENCE",
			ExpectedRecords: len(index[key]),
			Sections:        []string{"dual-entry-discovery", "staging-and-promotion"},
		})
	}
	return units
}

// MechanicalSharedReferenceExecutor materializes the persona invocation pointer
// rows for committed atomic components. It never authors semantic content; the
// invocation's persona/POV bindings select the components of its unit.
type MechanicalSharedReferenceExecutor struct {
	Components []ComponentRef
}

// Execute implements work.Executor.
func (e MechanicalSharedReferenceExecutor) Execute(inv work.Invocation) (work.Proposal, error) {
	proposal := work.Proposal{FieldID: "SHARED-REFERENCE"}
	for _, component := range e.Components {
		if component.Persona != inv.Unit.Persona || component.POV != inv.Unit.POV {
			continue
		}
		povFile := povFiles[component.POV]
		if povFile == "" {
			return work.Proposal{}, fmt.Errorf("shared-reference: component %s has unknown POV %q", component.ID, component.POV)
		}
		track := component.Track
		if track == "" {
			return work.Proposal{}, fmt.Errorf("shared-reference: component %s has no declarable traversal track", component.ID)
		}
		capabilityState := component.CapabilityState
		if capabilityState == "" {
			capabilityState = "Unknown"
		}
		stagingStatus := component.StagingStatus
		if stagingStatus == "" {
			stagingStatus = "[S-DISCOVERED]"
		}
		signature := component.CreatingInvocation
		if signature == "" {
			signature = "None"
		}
		payload := "- **canonical-location:** personas/_shared/" + povFile + " -> " + component.ID + "\n" +
			"- **local-trigger-src-coordinate:** " + component.VirtualCoordinate + "\n" +
			"- **local-invocation-signature:** " + signature + "\n" +
			"- **track:** " + track + "\n" +
			"- **capability-state:** " + capabilityState + "\n" +
			"- **canonical-staging-status:** " + stagingStatus + "\n"
		proposal.Records = append(proposal.Records, work.Record{Payload: []byte(payload)})
	}
	return proposal, nil
}

func splitTokens(raw string) []string {
	var tokens []string
	for _, token := range strings.Split(raw, ",") {
		token = strings.TrimSpace(token)
		if token == "" || strings.EqualFold(token, "none") {
			continue
		}
		tokens = append(tokens, token)
	}
	return tokens
}
