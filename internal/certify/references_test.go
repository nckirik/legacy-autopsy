package certify

import "testing"

func lookupNone(string, string) (bool, error) { return false, nil }

func TestValidateReferencesEntryCluster(t *testing.T) {
	lookup := func(fieldID, recordID string) (bool, error) {
		return fieldID == "PERSONA-REGISTRY" && recordID == "WEB", nil
	}
	valid := map[string]string{
		"cluster-id":              "WEB",
		"canonical-owner-persona": "WEB",
		"source-ids":              "None, SRC-abc-123",
	}
	if err := ValidateReferences("ENTRY-CLUSTER", valid, lookup); err != nil {
		t.Fatal(err)
	}
	for name, values := range map[string]map[string]string{
		"lowercase cluster id": {"cluster-id": "web"},
		"unknown owner":        {"canonical-owner-persona": "NOPE"},
		"bad source token":     {"source-ids": "passport/index.js"},
		"lowercase owner":      {"canonical-owner-persona": "web/ui"},
	} {
		if err := ValidateReferences("ENTRY-CLUSTER", merge(valid, values), lookup); err == nil {
			t.Fatalf("%s accepted", name)
		}
	}
	if err := ValidateReferences("INVARIANT", map[string]string{"x": "y"}, lookupNone); err != nil {
		t.Fatal(err)
	}
	shared := merge(valid, map[string]string{"canonical-owner-persona": "SHARED/entry"})
	if err := ValidateReferences("ENTRY-CLUSTER", shared, lookup); err != nil {
		t.Fatalf("shared owner rejected: %v", err)
	}
}

func merge(base, overrides map[string]string) map[string]string {
	out := map[string]string{}
	for key, value := range base {
		out[key] = value
	}
	for key, value := range overrides {
		out[key] = value
	}
	return out
}
