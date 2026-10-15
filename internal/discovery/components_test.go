package discovery

import (
	"strconv"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

func inventoryPayload(kind, coordinate, fingerprint string) string {
	return "- **kind:** " + kind + "\n" +
		"- **concrete-coordinate:** " + coordinate + "\n" +
		"- **parent-container-src:** None\n" +
		"- **source-fingerprint-snapshot-src-version:** " + fingerprint + " version 1\n" +
		"- **discovery-method:** other\n" +
		"- **scope-summary:** WEB\n" +
		"- **environment-runtime-state:** dev\n" +
		"- **traversal-tracks:** Normal\n" +
		"- **canonical-owner:** WEB\n" +
		"- **disposition-summary:** Unknown\n" +
		"- **mappings:** None\n" +
		"- **frontier-references:** None\n" +
		"- **coverage-rows:** None\n" +
		"- **coverage-summary:** C-GAP\n" +
		"- **claims:** None\n"
}

func componentPayload(hit Hit, srcID, pov, persona, track string) string {
	return "- **record-kind:** ATOMIC\n" +
		"- **source-ids:** " + srcID + "\n" +
		"- **file-path-line-range:** " + hit.File + ":" + strconv.Itoa(hit.Line) + "\n" +
		"- **virtual-coordinate:** " + hit.Coordinate + "\n" +
		"- **pov-owner:** " + pov + "\n" +
		"- **general-intent:** Serves the declared " + string(hit.Kind) + " unit.\n" +
		"- **inputs-outputs:** Unknown\n" +
		"- **auth-binding:** None\n" +
		"- **staging-status:** [S-DISCOVERED]\n" +
		"- **coverage-status:** [C-GAP]\n" +
		"- **personas:** " + persona + "\n" +
		"- **traversal-tracks:** " + track + "\n" +
		"- **capability-state:** Unknown\n" +
		"- **runtime-state-matrix:** None\n" +
		"- **depends-on-depended-on-by:** None\n" +
		"- **synthesis-references:** None\n" +
		"- **side-effects:** Unknown\n" +
		"- **exceptional-behaviors:** Unknown\n" +
		"- **edge-and-odd-cases:** Unknown\n" +
		"- **behavioral-notes-for-reconstruction:** Unknown\n" +
		"- **material-claims:** None\n" +
		"- **block-confidence-rank:** Low\n" +
		"- **block-evidence-profile:** [E-DIRECT] " + hit.Coordinate + "\n" +
		"- **evidence-anchors:** " + hit.Fingerprint + "\n" +
		"- **record-version-supersession:** version 1\n"
}

func TestComponentPlanAnchorsIdentityInHitSource(t *testing.T) {
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
	if len(units) != 4 {
		t.Fatalf("units = %d, want 4", len(units))
	}
	executor := fakeExecutor(func(inv work.Invocation) (work.Proposal, error) {
		var proposal work.Proposal
		for _, hit := range hits {
			if hitPOV(hit.Kind) != inv.Unit.POV {
				continue
			}
			payload := componentPayload(hit, srcIDs[hit.Coordinate], inv.Unit.POV, "WEB", "Normal")
			proposal.Records = append(proposal.Records, work.Record{Payload: []byte(payload)})
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
	ids, err := store.List(db, doc, "ATOMIC-COMPONENT")
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != len(hits) {
		t.Fatalf("components = %v", ids)
	}
	for _, hit := range hits {
		want, err := capabilities.TypedID("CMP", "demo", srcIDs[hit.Coordinate], "CMP")
		if err != nil {
			t.Fatal(err)
		}
		found := false
		for _, id := range ids {
			if id == want {
				found = true
			}
		}
		if !found {
			t.Fatalf("expected anchored CMP id %s in %v", want, ids)
		}
	}
	components, err := LoadComponents(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	if len(components) != len(hits) {
		t.Fatalf("loaded components = %+v", components)
	}
	for _, component := range components {
		if component.Persona != "WEB" || component.Track != "Normal" || component.StagingStatus != "[S-DISCOVERED]" || component.VirtualCoordinate == "" {
			t.Fatalf("component = %+v", component)
		}
	}
	creating, err := LoadCreatingInvocations(db, doc)
	if err != nil {
		t.Fatal(err)
	}
	for _, component := range components {
		if !strings.HasPrefix(creating[component.ID], "ALFA-WEB-POV-") {
			t.Fatalf("creating invocation for %s = %q", component.ID, creating[component.ID])
		}
	}
	refUnits := SharedReferencePlan(components)
	if len(refUnits) != 4 {
		t.Fatalf("shared-reference units = %d, want 4 (=%+v)", len(refUnits), refUnits)
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
	if _, err := refSession.Run(refUnits); err != nil {
		t.Fatal(err)
	}
	refIDs, err := store.List(db, doc, "SHARED-REFERENCE")
	if err != nil {
		t.Fatal(err)
	}
	if len(refIDs) != len(hits) {
		t.Fatalf("shared references = %v", refIDs)
	}
	payload, err := store.Get(db, doc, "SHARED-REFERENCE", refIDs[0])
	if err != nil {
		t.Fatal(err)
	}
	text := string(payload)
	if !strings.Contains(text, "personas/_shared/") || !strings.Contains(text, " -> CMP-") || !strings.Contains(text, "- **canonical-staging-status:** [S-DISCOVERED]") {
		t.Fatalf("shared reference payload:\n%s", text)
	}
	for _, marker := range []string{"- **track:** Normal", "- **local-invocation-signature:** ALFA-WEB-POV-"} {
		if !strings.Contains(text, marker) {
			t.Fatalf("shared reference payload missing %q:\n%s", marker, text)
		}
	}
}

func TestSharedReferencePlanSkipsUnboundPersonas(t *testing.T) {
	components := []ComponentRef{
		{ID: "CMP-AAAAAAAAAAAA", Persona: "WEB", POV: "POV-5"},
		{ID: "CMP-BBBBBBBBBBBB", Persona: "", POV: "POV-5"},
		{ID: "CMP-CCCCCCCCCCCC", Persona: "WEB", POV: "POV-9"},
	}
	units := SharedReferencePlan(components)
	if len(units) != 1 || units[0].ID != "U-DISCOVERY-SHARED-REFS-WEB-POV-5" || units[0].ExpectedRecords != 1 {
		t.Fatalf("units = %+v", units)
	}
}
