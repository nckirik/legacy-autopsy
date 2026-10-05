package store

import (
	"os/exec"
	"strings"
	"testing"
)

const preflightFixture = `# 0A-PREFLIGHT: Demo

## 1. Metadata and Boundaries

- **System Namespace:** demo
- **Target Repository / Evidence Roots:** /tmp/demo
- **Primary Technology Stack:** Go
- **Protocol Version:** 4.2.1
- **Current Iteration:** ALFA
- **Default Max Traversal Depth:** 3
- **Included Environments / Snapshots:** dev
- **Explicit Scope Policy:** full tree

## 2. Persona Registry

| Persona | Prefix | Persona Slug | Persona Directory | Actor Class | Canonical Entry IDs | Scope Paths | Registry State |
| :-- | :-- | :-- | :-- | :-- | :-- | :-- | :-- |
| Backend | BE | backend | BE-backend | human | SRC-1 | src/ | Registered |
| Operator | OP | operator | OP-operator | external | SRC-2 | ops/ | Registered |

## 3. Entry-Point Clusters

### [CL-1] Backend entry

- **Concrete Entry Coordinates:** cmd/server
`

func TestImportPreflight(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	db, res := ingestStore(t)
	if err := ImportPreflight(db, res.EIR, []byte(preflightFixture)); err != nil {
		t.Fatal(err)
	}
	ids, err := List(db, res.EIR, "PREFLIGHT-HEADER")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(ids, ",") != PreflightRecordID {
		t.Fatalf("preflight header IDs = %v", ids)
	}
	header, err := Get(db, res.EIR, "PREFLIGHT-HEADER", PreflightRecordID)
	if err != nil {
		t.Fatal(err)
	}
	for _, marker := range []string{
		"- **System Namespace:** demo",
		"- **Canonical Storage Binding:** markdown/1",
		"- **Default Max Traversal Depth:** 3",
	} {
		if !strings.Contains(string(header), marker) {
			t.Fatalf("preflight payload missing %q:\n%s", marker, header)
		}
	}
	personas, err := List(db, res.EIR, "PERSONA-REGISTRY")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(personas, ",") != "BE,OP" {
		t.Fatalf("persona IDs = %v", personas)
	}
	backend, err := Get(db, res.EIR, "PERSONA-REGISTRY", "BE")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(backend), "- **Persona Directory:** BE-backend") {
		t.Fatalf("persona payload:\n%s", backend)
	}
	view, err := Render(db, res.EIR, "PERSONA-REGISTRY")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(view), "## BE") || !strings.Contains(string(view), "## OP") {
		t.Fatalf("human render missing records:\n%s", view)
	}
	for _, forbidden := range []string{"la_semantic_fingerprint", "DEBUG RENDER", "sha256:"} {
		if strings.Contains(string(view), forbidden) {
			t.Fatalf("human render contains mechanical content %q", forbidden)
		}
	}
}

func TestImportPreflightFailsClosed(t *testing.T) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not available")
	}
	db, res := ingestStore(t)
	if err := ImportPreflight(db, res.EIR, []byte("# nothing here\n")); err == nil {
		t.Fatal("missing metadata accepted")
	}
	noPersonas := strings.Split(preflightFixture, "## 2. Persona Registry")[0]
	if err := ImportPreflight(db, res.EIR, []byte(noPersonas)); err == nil {
		t.Fatal("missing persona registry accepted")
	}
}
