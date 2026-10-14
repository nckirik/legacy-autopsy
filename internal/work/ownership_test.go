package work

import (
	"strings"
	"testing"
)

// TestInvocationOwnershipIsolation enforces §2.2 canonical entry ownership and
// §8.2 persona isolation at commit time: a persona-bound invocation may not
// create a record that claims another persona as its canonical owner, while
// SHARED/None owners and unbound invocations stay valid.
func TestInvocationOwnershipIsolation(t *testing.T) {
	db, doc := testDoc(t)
	cases := []struct {
		name    string
		persona string
		owner   string
		wantErr string
	}{
		{name: "bound owner matches", persona: "WEB", owner: "WEB", wantErr: ""},
		{name: "shared owner accepted", persona: "WEB", owner: "SHARED", wantErr: ""},
		{name: "none owner accepted", persona: "WEB", owner: "None", wantErr: ""},
		{name: "unbound invocation accepted", persona: "", owner: "AD", wantErr: ""},
		{name: "foreign owner rejected", persona: "WEB", owner: "AD", wantErr: "may not write another persona's record"},
		{name: "foreign prefixed owner rejected", persona: "WEB", owner: "AD-admin", wantErr: "may not write another persona's record"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			unit := Unit{
				ID:               "U-OWNERSHIP-" + strings.ReplaceAll(tc.name, " ", "-"),
				Mode:             "Discovery",
				Persona:          tc.persona,
				IDFromCoordinate: "concrete-coordinate",
				IDType:           "SRC",
				Namespace:        "demo",
				ResultType:       "SOURCE-INVENTORY-RECORD",
			}
			payload := "- **kind:** FILE\n" +
				"- **concrete-coordinate:** app/" + strings.ReplaceAll(tc.name, " ", "-") + ".py\n" +
				"- **parent-container-src:** None\n" +
				"- **source-fingerprint-snapshot-src-version:** sha256:aaa version 1\n" +
				"- **discovery-method:** filesystem\n" +
				"- **scope-summary:** demo\n" +
				"- **environment-runtime-state:** dev\n" +
				"- **traversal-tracks:** Normal\n" +
				"- **canonical-owner:** " + tc.owner + "\n" +
				"- **disposition-summary:** Unknown\n" +
				"- **mappings:** None\n" +
				"- **frontier-references:** None\n" +
				"- **coverage-rows:** None\n" +
				"- **coverage-summary:** C-GAP\n" +
				"- **claims:** None\n"
			_, err := CommitWithIDs(NewCaps(db, doc), unit, Proposal{FieldID: "SOURCE-INVENTORY-RECORD", Records: []Record{{Payload: []byte(payload)}}}, false)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("error = %v, want %q", err, tc.wantErr)
			}
		})
	}
}
