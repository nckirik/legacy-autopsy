package discovery

import (
	"strings"
	"testing"
)

func TestSharedParticipantsMatchingAndOrder(t *testing.T) {
	entries := []SharedEntry{
		{
			ConcreteEntrySrc:  "app/routes.py",
			CoordinationOwner: "WEB",
			Participants:      []string{"AD", "OE"},
		},
		{
			ConcreteEntrySrc: "other/entry.py",
			Participants:     []string{"SA"},
		},
	}
	cluster := Cluster{ID: "WEB", Entries: []string{"app/routes.py"}}
	got := SharedParticipants(entries, cluster)
	if strings.Join(got, ",") != "WEB,AD,OE" {
		t.Fatalf("participants = %v, want coordination owner first then sorted", got)
	}
	if none := SharedParticipants(entries, Cluster{ID: "X", Entries: []string{"nothing.py"}}); len(none) != 0 {
		t.Fatalf("unmatched cluster got participants: %v", none)
	}
}

func TestPersonaTokensParseDeclarationForms(t *testing.T) {
	raw := "AD: local condition; OE: own requests only, WEB: shared entry"
	got := personaTokens(raw)
	if strings.Join(got, ",") != "AD,OE,WEB" {
		t.Fatalf("tokens = %v", got)
	}
	pointers := personaTokens("SD-admin/01-FRONTEND-POV.md, personas/MK-teknik/03-BACKGROUND-POV.md")
	if strings.Join(pointers, ",") != "SD,MK" {
		t.Fatalf("pointer tokens = %v", pointers)
	}
	if reserved := personaTokens("SHARED: coordination only"); len(reserved) != 0 {
		t.Fatalf("reserved token parsed as persona: %v", reserved)
	}
}

func TestHitPlanResolvesSharedEntryPersona(t *testing.T) {
	hits := []Hit{{Coordinate: "app/routes.py :: route :: GET /x", Kind: HitRoute, File: "app/routes.py"}}
	clusters := []Cluster{{ID: "WEB", Owner: "SHARED", Track: "Normal", Entries: []string{"app/routes.py"}}}
	directives := Directives{
		Namespace:          "demo",
		Snapshot:           "dev",
		BatchSize:          10,
		PersonaDirectories: map[string]string{"AD": "AD-admin"},
		SharedEntries: []SharedEntry{{
			ConcreteEntrySrc:  "app/routes.py",
			CoordinationOwner: "MK",
			Participants:      []string{"AD"},
		}},
	}
	units := HitPlan(hits, clusters, directives)
	if len(units) != 1 || units[0].Persona != "AD" {
		t.Fatalf("units = %+v", units)
	}
	context := string(units[0].Inputs[0].Bytes)
	for _, want := range []string{
		"- **Persona:** AD",
		"- **Persona Directory:** AD-admin",
		"- **Active POV File:** personas/AD-admin/05-WIRING-POV.md",
	} {
		if !strings.Contains(context, want) {
			t.Fatalf("context missing %q:\n%s", want, context)
		}
	}

	// An unregistered coordination owner falls back to a registered participant.
	directives.SharedEntries[0].CoordinationOwner = "MK"
	units = HitPlan(hits, clusters, directives)
	if len(units) != 1 || units[0].Persona != "AD" {
		t.Fatalf("fallback units = %+v", units)
	}
}

func TestPOVFilePathShared(t *testing.T) {
	if got := povFilePath("", "POV-4"); got != "personas/_shared/04-DATA-POV.md" {
		t.Fatalf("shared pov path = %q", got)
	}
	if got := povFilePath("AD-admin", "POV-4"); got != "personas/AD-admin/04-DATA-POV.md" {
		t.Fatalf("local pov path = %q", got)
	}
	if got := povFilePath("AD-admin", "POV-9"); got != "None" {
		t.Fatalf("unknown pov path = %q", got)
	}
}
