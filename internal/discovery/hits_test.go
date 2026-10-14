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
		"db/schema.sql":  "CREATE TABLE malzemeler (\n  id INT\n);\n-- SELECT * FROM secret\nSELECT id, ad FROM malzemeler;\nSELECT id, ad FROM malzemeler;\n",
		"jobs/tasks.py":  "import schedule\nschedule.every(10).minutes.do(scan)\n",
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
		{"db/schema.sql :: query :: CREATE TABLE malzemeler (", HitQuery, 1},
		{"db/schema.sql :: query :: SELECT id, ad FROM malzemeler", HitQuery, 5},
		{"db/schema.sql :: query :: SELECT id, ad FROM malzemeler #2", HitQuery, 6},
		{"jobs/tasks.py :: job :: schedule.every(10).minutes.do(scan)", HitJob, 2},
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
	units := HitPlan(hits, []Cluster{{ID: "WEB", Title: "Web", Owner: "WEB"}}, Directives{Namespace: "demo", Snapshot: "dev", StorageBinding: "sqlite/1", BatchSize: 10})
	if len(units) != 1 || units[0].ExpectedRecords != len(hits) {
		t.Fatalf("units = %+v", units)
	}
	executor := fakeExecutor(func(inv work.Invocation) (work.Proposal, error) {
		if !strings.Contains(string(inv.Context.Bytes), "app/routes.py") {
			t.Error("context lacks hit coordinates")
		}
		var proposal work.Proposal
		for _, hit := range hits {
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
	if _, err := session.Run(units); err != nil {
		t.Fatal(err)
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
