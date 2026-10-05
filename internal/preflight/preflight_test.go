package preflight

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/spec"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

func testTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"go.mod":                "module demo\n",
		"main.go":               "package main\n",
		"internal/service/a.go": "package service\n",
		"web/index.js":          "console.log(1)\n",
		"Dockerfile":            "FROM scratch\n",
		".git/config":           "[core]\n",
		"node_modules/x.js":     "x\n",
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

func TestScanRepoDeterministic(t *testing.T) {
	root := testTree(t)
	first, err := ScanRepo(root, 100, 10, 6)
	if err != nil {
		t.Fatal(err)
	}
	if first.Files != 5 {
		t.Fatalf("files = %d, want 5 (skips .git and node_modules)", first.Files)
	}
	if strings.Join(first.Languages, ",") != "Go,JavaScript" {
		t.Fatalf("languages = %v", first.Languages)
	}
	for _, want := range []string{"Dockerfile", "go.mod", "main.go"} {
		if !contains(first.EntryPoints, want) {
			t.Fatalf("entry points %v missing %q", first.EntryPoints, want)
		}
	}
	second, err := ScanRepo(root, 100, 10, 6)
	if err != nil {
		t.Fatal(err)
	}
	if first.Fingerprint != second.Fingerprint || string(first.Markdown()) != string(second.Markdown()) {
		t.Fatal("scan is not deterministic")
	}
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func testDoc(t *testing.T) (*cdl.EIRDoc, string) {
	t.Helper()
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	res, err := spec.Compile("../..")
	if err != nil {
		t.Fatal(err)
	}
	db := filepath.Join(t.TempDir(), "store.db")
	ddl, err := store.Schema(res.EIR, store.Provenance{Generator: "preflight/test", SourceFingerprint: "sha256:test"})
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

func TestPreflightPlanRunsToStore(t *testing.T) {
	doc, db := testDoc(t)
	scan, err := ScanRepo(testTree(t), 100, 10, 6)
	if err != nil {
		t.Fatal(err)
	}
	units := Plan(scan, DefaultDirectives())
	if len(units) != 2 {
		t.Fatalf("units = %d", len(units))
	}
	executor := fakeExecutor(func(inv work.Invocation) (work.Proposal, error) {
		context := string(inv.Context.Bytes)
		if !strings.Contains(context, "Repository scan") || !strings.Contains(context, "Canonical Storage Binding") {
			t.Errorf("unit %s context lacks the bounded scan input", inv.Unit.ID)
		}
		if !strings.Contains(context, "Normative protocol sections") {
			t.Errorf("unit %s context lacks normative protocol sections", inv.Unit.ID)
		}
		switch inv.Unit.ID {
		case "U-PREFLIGHT-HEADER":
			return work.Proposal{FieldID: "PREFLIGHT-HEADER", Records: []work.Record{{Payload: []byte(
				"- **System Namespace:** demo\n" +
					"- **Target Repository / Evidence Roots:** /tmp/demo\n" +
					"- **Primary Technology Stack:** Go, JavaScript\n" +
					"- **Protocol Version:** 4.2.1\n" +
					"- **Canonical Storage Binding:** sqlite/1\n" +
					"- **Current Iteration:** ALFA\n" +
					"- **Default Max Traversal Depth:** 3\n" +
					"- **Included Environments / Snapshots:** dev\n" +
					"- **Explicit Scope Policy:** full tree\n")}}}, nil
		case "U-PREFLIGHT-PERSONAS":
			return work.Proposal{FieldID: "PERSONA-REGISTRY", Records: []work.Record{
				{Payload: []byte("- **Persona:** Backend\n- **Prefix:** BE\n- **Persona Slug:** backend\n- **Persona Directory:** BE-backend\n- **Actor Class:** human\n- **Canonical Entry IDs:** cmd/server\n- **Scope Paths:** internal/\n- **Registry State:** Registered\n")},
				{Payload: []byte("- **Persona:** Web\n- **Prefix:** WEB\n- **Persona Slug:** web\n- **Persona Directory:** WEB-web\n- **Actor Class:** human\n- **Canonical Entry IDs:** web/index.js\n- **Scope Paths:** web/\n- **Registry State:** Registered\n")},
			}}, nil
		default:
			t.Fatalf("unexpected unit %s", inv.Unit.ID)
			return work.Proposal{}, nil
		}
	})
	session, err := work.NewSession(db, doc, work.DefaultConfig(), executor, nil)
	if err != nil {
		t.Fatal(err)
	}
	report, err := session.Run(units)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Units) != 2 || report.Units[1].Retries != 0 {
		t.Fatalf("report = %+v", report)
	}
	personas, err := store.List(db, doc, "PERSONA-REGISTRY")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(personas, ",") != "BE,WEB" {
		t.Fatalf("personas = %v", personas)
	}
	header, err := store.Get(db, doc, "PREFLIGHT-HEADER", "0A")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(header), "Canonical Storage Binding:** sqlite/1") {
		t.Fatalf("header = %s", header)
	}
}
