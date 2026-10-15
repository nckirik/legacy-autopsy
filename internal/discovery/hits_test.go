package discovery

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

func hitsTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"app/routes.py": strings.Join([]string{
			"import flask",
			"app = flask.Flask(__name__)",
			"",
			"# @app.route(\"/commented-out\")",
			"@app.route(\"/requests\", methods=[\"GET\"])",
			"def list_requests(): ...",
			"",
			"@app.post(\"/requests/<int:rid>/approve\")",
			"def approve(rid): ...",
		}, "\n") + "\n",
		"app/server.js":  "router.get('/api/units', handler)\n// router.post('/api/ignored', h)\n",
		"db/schema.sql":  "CREATE TABLE malzemeler (\n  id INT\n);\nCREATE VIEW vw_talep AS SELECT 1;\nCREATE TRIGGER trg_touch BEFORE INSERT ON malzemeler BEGIN SELECT 1; END;\n-- SELECT * FROM secret\nSELECT id, ad FROM malzemeler;\nSELECT id, ad FROM malzemeler;\n",
		"jobs/tasks.py":  "import schedule\nschedule.every(10).minutes.do(scan)\n",
		"pages/units.py": "img = qrcode.make(value)\nst.download_button(\"QR PNG\", data, file_name=\"unit.png\", mime=\"image/png\")\n",
		"notes.md":       "@app.route(\"/ignored-markdown\")\n",
		"decoy/plain.py": "print('nothing to see')\n",
	}
	for path, content := range files {
		full := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestEnumerateHitsDeterministicAndBounded(t *testing.T) {
	root := hitsTree(t)
	files, err := Enumerate(root, nil, 100)
	if err != nil {
		t.Fatal(err)
	}
	hits, err := EnumerateHits(root, files, 100)
	if err != nil {
		t.Fatal(err)
	}
	type want struct {
		coordinate string
		kind       HitKind
		line       int
	}
	wants := []want{
		{"app/routes.py :: route :: /requests", HitRoute, 5},
		{"app/routes.py :: route :: POST /requests/<int:rid>/approve", HitRoute, 8},
		{"app/server.js :: route :: GET /api/units", HitRoute, 1},
		{"db/schema.sql :: table :: malzemeler", HitTable, 1},
		{"db/schema.sql :: view :: vw_talep", HitView, 4},
		{"db/schema.sql :: dr :: trg_touch", HitDR, 5},
		{"db/schema.sql :: query :: SELECT id, ad FROM malzemeler", HitQuery, 7},
		{"db/schema.sql :: query :: SELECT id, ad FROM malzemeler #2", HitQuery, 8},
		{"jobs/tasks.py :: job :: schedule.every(10).minutes.do(scan)", HitJob, 2},
		{"pages/units.py :: action :: st.download_button(\"QR PNG\", data, file_name=\"unit.png\", mime=\"image/png\")", HitAction, 2},
	}
	if len(hits) != len(wants) {
		for _, hit := range hits {
			t.Logf("hit: %s %s line %d", hit.Kind, hit.Coordinate, hit.Line)
		}
		t.Fatalf("hits = %d, want %d", len(hits), len(wants))
	}
	for i, expected := range wants {
		if hits[i].Coordinate != expected.coordinate || hits[i].Kind != expected.kind || hits[i].Line != expected.line {
			t.Fatalf("hit[%d] = %+v, want %+v", i, hits[i], expected)
		}
		if !strings.HasPrefix(hits[i].Fingerprint, "sha256:") {
			t.Fatalf("hit[%d] fingerprint = %q", i, hits[i].Fingerprint)
		}
	}
	for _, hit := range hits {
		for _, decoy := range []string{"commented-out", "secret", "ignored", "nothing to see"} {
			if strings.Contains(hit.Coordinate, decoy) || strings.Contains(hit.Evidence, decoy) {
				t.Fatalf("decoy %q enumerated: %+v", decoy, hit)
			}
		}
	}
	again, err := EnumerateHits(root, files, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(again) != len(hits) {
		t.Fatalf("second run = %d hits", len(again))
	}
	for i := range hits {
		if again[i] != hits[i] {
			t.Fatalf("enumeration is not deterministic at %d: %+v vs %+v", i, again[i], hits[i])
		}
	}
	bounded, err := EnumerateHits(root, files, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(bounded) != 2 || bounded[0] != hits[0] || bounded[1] != hits[1] {
		t.Fatalf("bounded = %+v", bounded)
	}
}

func TestClusterForHitBinding(t *testing.T) {
	clusters := []Cluster{
		{ID: "API", Entries: []string{"api/routes.py"}},
		{ID: "WEB", Entries: []string{"web/app.py"}},
	}
	cases := []struct {
		file string
		want string
	}{
		{"web/app.py", "WEB"},
		{"web/pages/dashboard.py", "WEB"},
		{"api/handlers/users.py", "API"},
		{"shared/util.py", ""},
	}
	for _, tc := range cases {
		cluster, ok := ClusterForHit(Hit{File: tc.file}, clusters)
		if tc.want == "" {
			if ok {
				t.Fatalf("hit %s bound to %s, want unbound", tc.file, cluster.ID)
			}
			continue
		}
		if !ok || cluster.ID != tc.want {
			t.Fatalf("hit %s bound = %+v/%v, want %s", tc.file, cluster, ok, tc.want)
		}
	}
	cluster, ok := ClusterForHit(Hit{File: "shared/util.py"}, []Cluster{{ID: "WEB", Entries: []string{"web/app.py"}}})
	if !ok || cluster.ID != "WEB" {
		t.Fatalf("single-cluster fallback = %+v/%v", cluster, ok)
	}
}

func TestHitPlanLeavesSharedOwnerUnbound(t *testing.T) {
	hits := []Hit{{Coordinate: "web/routes.py :: route :: GET /x", Kind: HitRoute, File: "web/routes.py"}}
	units := HitPlan(hits, []Cluster{{ID: "WEB", Owner: "SHARED", Track: "Normal", Entries: []string{"web/routes.py"}}}, Directives{Namespace: "demo", Snapshot: "dev", BatchSize: 10, PersonaDirectories: map[string]string{"WEB": "WEB-test", "AD": "AD-test"}})
	if len(units) != 1 {
		t.Fatalf("units = %+v", units)
	}
	if units[0].Persona != "" || units[0].POV != "POV-5" || units[0].ID != "U-DISCOVERY-HITS-WEB-POV-5-001" {
		t.Fatalf("unit = %+v", units[0])
	}
	if !strings.Contains(string(units[0].Inputs[0].Bytes), "shared and no committed shared-entry record") {
		t.Fatalf("context lacks the shared-ownership note:\n%s", units[0].Inputs[0].Bytes)
	}
}

func TestHitPlanCommitsInventoryWithRuntimeIDs(t *testing.T) {
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
	units := HitPlan(hits, []Cluster{{ID: "WEB", Title: "Web", Owner: "WEB", PrimaryPOV: "POV-1", Track: "Normal", Entries: []string{"app/routes.py"}}}, Directives{Namespace: "demo", Snapshot: "dev", StorageBinding: "sqlite/1", BatchSize: 10, PersonaDirectories: map[string]string{"WEB": "WEB-test"}})
	if len(units) != 4 {
		t.Fatalf("units = %+v", units)
	}
	povs := map[string]string{}
	for _, unit := range units {
		povs[unit.POV] = unit.ID
		if unit.Persona != "WEB" {
			t.Fatalf("unit %s persona = %q, want WEB", unit.ID, unit.Persona)
		}
	}
	if povs["POV-1"] != "U-DISCOVERY-HITS-WEB-POV-1-001" || povs["POV-3"] != "U-DISCOVERY-HITS-WEB-POV-3-001" || povs["POV-4"] != "U-DISCOVERY-HITS-WEB-POV-4-001" || povs["POV-5"] != "U-DISCOVERY-HITS-WEB-POV-5-001" {
		t.Fatalf("pov units = %v", povs)
	}
	executor := fakeExecutor(func(inv work.Invocation) (work.Proposal, error) {
		if !strings.Contains(string(inv.Context.Bytes), "## Candidate hits") || !strings.Contains(string(inv.Context.Bytes), inv.Unit.POV) {
			t.Error("context lacks hit coordinates or the bound POV")
		}
		var proposal work.Proposal
		for _, hit := range hits {
			if hitPOV(hit.Kind) != inv.Unit.POV {
				continue
			}
			payload := "- **kind:** " + string(hit.Kind) + "\n" +
				"- **concrete-coordinate:** " + hit.Coordinate + "\n" +
				"- **parent-container-src:** None\n" +
				"- **source-fingerprint-snapshot-src-version:** " + hit.Fingerprint + " version 1\n" +
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
			proposal.Records = append(proposal.Records, work.Record{Payload: []byte(payload)})
		}
		proposal.FieldID = "SOURCE-INVENTORY-RECORD"
		return proposal, nil
	})
	session, err := work.NewSession(db, doc, work.DefaultConfig(), executor, nil)
	if err != nil {
		t.Fatal(err)
	}
	report, err := session.Run(units)
	if err != nil {
		t.Fatal(err)
	}
	personaBound := false
	for _, stats := range report.Units {
		if strings.HasPrefix(stats.InvocationID, "ALFA-WEB-POV-") {
			personaBound = true
		}
	}
	if !personaBound {
		t.Fatalf("no persona/POV-bound invocation ID: %+v", report.Units)
	}
	ids, err := store.List(db, doc, "SOURCE-INVENTORY-RECORD")
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != len(hits) {
		t.Fatalf("ids = %v", ids)
	}
	for _, kind := range []string{"ROUTE", "QUERY", "JOB"} {
		found := false
		for _, hit := range hits {
			if string(hit.Kind) != kind {
				continue
			}
			want, err := capabilities.SourceID("demo", kind, hit.Coordinate)
			if err != nil {
				t.Fatal(err)
			}
			for _, id := range ids {
				if id == want {
					found = true
				}
			}
		}
		if !found {
			t.Fatalf("expected a runtime SRC-%s id in %v", kind, ids)
		}
	}
}
