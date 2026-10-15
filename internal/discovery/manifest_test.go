package discovery

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

func TestPOVManifestsProjectSharedReferences(t *testing.T) {
	doc, db := testDocAndDB(t)
	root := hitsTree(t)
	files, err := Enumerate(root, nil, 100)
	if err != nil {
		t.Fatal(err)
	}
	hits, err := EnumerateHits(root, files, 100)
	if err != nil {
		t.Fatal(err)
	}
	header := []byte("- **system-namespace:** demo\n- **target-repository-evidence-roots:** app, db, jobs\n- **primary-technology-stack:** Python\n- **protocol-version:** 4.2.1\n- **canonical-storage-binding:** sqlite/1\n- **current-iteration:** ALFA\n- **default-max-traversal-depth:** 3\n- **included-environments-snapshots:** dev\n- **explicit-scope-policy:** repo\n")
	if err := store.Ingest(db, doc, "PREFLIGHT-HEADER", "0A", 1, header); err != nil {
		t.Fatal(err)
	}
	persona := []byte("- **persona:** Web Kullanıcısı\n- **prefix:** WEB\n- **persona-slug:** web-kullanici\n- **persona-directory:** WEB-web-kullanici\n- **actor-class:** human\n- **canonical-entry-ids:** None\n- **scope-paths:** app/\n- **registry-state:** Active\n")
	if err := store.Ingest(db, doc, "PERSONA-REGISTRY", "WEB", 1, persona); err != nil {
		t.Fatal(err)
	}
	srcIDs := map[string]string{}
	for _, hit := range hits {
		id, err := capabilities.SourceID("demo", string(hit.Kind), hit.Coordinate)
		if err != nil {
			t.Fatal(err)
		}
		srcIDs[hit.Coordinate] = id
		if err := store.Ingest(db, doc, "SOURCE-INVENTORY-RECORD", id, 1, []byte(inventoryPayload(string(hit.Kind), hit.Coordinate, hit.Fingerprint))); err != nil {
			t.Fatal(err)
		}
	}
	clusters := []Cluster{{ID: "WEB", Title: "Web", Owner: "WEB", PrimaryPOV: "POV-1", Track: "Normal", Entries: []string{"app/routes.py"}}}
	directives := Directives{Namespace: "demo", Snapshot: "dev", BatchSize: 10, PersonaDirectories: map[string]string{"WEB": "WEB-test"}}
	units, err := ComponentPlan(root, hits, clusters, directives)
	if err != nil {
		t.Fatal(err)
	}
	executor := fakeExecutor(func(inv work.Invocation) (work.Proposal, error) {
		var proposal work.Proposal
		for _, hit := range hits {
			if hitPOV(hit.Kind) != inv.Unit.POV {
				continue
			}
			proposal.Records = append(proposal.Records, work.Record{Payload: []byte(componentPayload(hit, srcIDs[hit.Coordinate], inv.Unit.POV, "WEB", "Normal"))})
		}
		proposal.FieldID = "ATOMIC-COMPONENT"
		return proposal, nil
	})
	session, err := work.NewSession(db, doc, work.DefaultConfig(), executor, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := session.Run(units); err != nil {
		t.Fatal(err)
	}
	components, err := LoadComponents(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	creating, err := LoadCreatingInvocations(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	for i := range components {
		if invocation, ok := creating[components[i].ID]; ok {
			components[i].CreatingInvocation = invocation
		}
	}
	refSession, err := work.NewSession(db, doc, work.DefaultConfig(), MechanicalSharedReferenceExecutor{Components: components}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := refSession.Run(SharedReferencePlan(components)); err != nil {
		t.Fatal(err)
	}
	manifests, err := POVManifests(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	if len(manifests) != 4 {
		t.Fatalf("manifests = %d, want 4", len(manifests))
	}
	paths := map[string]int{}
	for _, manifest := range manifests {
		paths[manifest.Path] = len(manifest.Markdown)
	}
	for _, want := range []string{"personas/WEB-web-kullanici/01-FRONTEND-POV.md", "personas/WEB-web-kullanici/03-BACKGROUND-POV.md", "personas/WEB-web-kullanici/04-DATA-POV.md", "personas/WEB-web-kullanici/05-WIRING-POV.md"} {
		if paths[want] == 0 {
			t.Fatalf("missing manifest %s in %v", want, paths)
		}
	}
	text := string(manifests[3].Markdown)
	for _, marker := range []string{"- **Persona Directory:** WEB-web-kullanici", "- **POV:** POV-5 (System Wiring Workflows)", "canonical location: personas/_shared/05-WIRING-POV.md -> CMP-"} {
		if !strings.Contains(text, marker) {
			t.Fatalf("manifest missing %q:\n%s", marker, text)
		}
	}
	again, err := POVManifests(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	for i := range manifests {
		if string(again[i].Markdown) != string(manifests[i].Markdown) {
			t.Fatalf("manifest projection is not deterministic at %d", i)
		}
	}
}

func TestComponentIDFromLocation(t *testing.T) {
	cases := map[string]string{
		"personas/_shared/05-WIRING-POV.md -> CMP-AAAAAAAAAAAA": "CMP-AAAAAAAAAAAA",
		"personas/_shared/05-WIRING-POV.md":                     "",
		"personas/_shared/05-WIRING-POV.md -> SRC-FILE-X":       "",
	}
	for input, want := range cases {
		if got := componentIDFromLocation(input); got != want {
			t.Fatalf("componentIDFromLocation(%q) = %q, want %q", input, got, want)
		}
	}
}
