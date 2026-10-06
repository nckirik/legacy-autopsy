package discovery

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/spec"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

func testTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"app/a.js":          "console.log(1)\n",
		"routes/api.js":     "module.exports = 1\n",
		"node_modules/x.js": "x\n",
		"db/schema.sql":     "select 1;\n",
		"public/index.html": "<html></html>\n",
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

func TestEnumerateDeterministicAndScoped(t *testing.T) {
	root := testTree(t)
	all, err := Enumerate(root, nil, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 4 {
		t.Fatalf("candidates = %d, want 4", len(all))
	}
	if all[0].Path != "app/a.js" {
		t.Fatalf("order = %v", all)
	}
	if !strings.HasPrefix(all[0].Fingerprint, "sha256:") {
		t.Fatalf("fingerprint = %q", all[0].Fingerprint)
	}
	scoped, err := Enumerate(root, []string{"routes"}, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(scoped) != 1 || scoped[0].Path != "routes/api.js" {
		t.Fatalf("scoped = %v", scoped)
	}
	again, _ := Enumerate(root, nil, 100)
	if again[0].Fingerprint != all[0].Fingerprint {
		t.Fatal("enumeration is not deterministic")
	}
}

func testDocAndDB(t *testing.T) (*cdl.EIRDoc, string) {
	t.Helper()
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	res, err := spec.Compile("../..")
	if err != nil {
		t.Fatal(err)
	}
	db := filepath.Join(t.TempDir(), "store.db")
	ddl, err := store.Schema(res.EIR, store.Provenance{Generator: "discovery/test", SourceFingerprint: "sha256:test"})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Apply(db, ddl); err != nil {
		t.Fatal(err)
	}
	return res.EIR, db
}

type fakeExecutor func(work.Invocation) (work.Proposal, error)

func (f fakeExecutor) Execute(inv work.Invocation) (work.Proposal, error) { return f(inv) }

func TestPlanCommitsInventoryWithRuntimeIDs(t *testing.T) {
	doc, db := testDocAndDB(t)
	root := testTree(t)
	candidates, err := Enumerate(root, nil, 100)
	if err != nil {
		t.Fatal(err)
	}
	header := []byte("- **system-namespace:** demo\n- **target-repository-evidence-roots:** app, routes, db, public\n- **primary-technology-stack:** JavaScript\n- **protocol-version:** 4.2.1\n- **canonical-storage-binding:** sqlite/1\n- **current-iteration:** ALFA\n- **default-max-traversal-depth:** 3\n- **included-environments-snapshots:** dev\n- **explicit-scope-policy:** repo\n")
	if err := store.Ingest(db, doc, "PREFLIGHT-HEADER", "0A", 1, header); err != nil {
		t.Fatal(err)
	}
	units := Plan(candidates, []Cluster{{ID: "WEB", Title: "Web", Owner: "WEB"}}, Directives{Namespace: "demo", Snapshot: "dev", StorageBinding: "sqlite/1", BatchSize: 10})
	if len(units) != 1 || units[0].ExpectedRecords != 4 {
		t.Fatalf("units = %+v", units)
	}
	executor := fakeExecutor(func(inv work.Invocation) (work.Proposal, error) {
		if !strings.Contains(string(inv.Context.Bytes), "routes/api.js") {
			t.Error("context lacks candidates")
		}
		var proposal work.Proposal
		for _, candidate := range candidates {
			payload := "- **kind:** FILE\n" +
				"- **concrete-coordinate:** " + candidate.Path + "\n" +
				"- **parent-container-src:** None\n" +
				"- **source-fingerprint-snapshot-src-version:** " + candidate.Fingerprint + " version 1\n" +
				"- **discovery-method:** filesystem\n" +
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
	if len(ids) != 4 {
		t.Fatalf("ids = %v", ids)
	}
	want, _ := capabilities.SourceID("demo", "FILE", "routes/api.js")
	found := false
	for _, id := range ids {
		if id == want {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected runtime SRC id %s in %v", want, ids)
	}
}
