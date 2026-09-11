package contextpacket

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nckirik/legacy-autopsy/internal/protocol"
	"github.com/nckirik/legacy-autopsy/internal/routing"
	"github.com/nckirik/legacy-autopsy/internal/workspace"
)

func TestBuildPreflightPacket(t *testing.T) {
	root, model, manifest, options := fixture(t)
	packet, err := Build(root, model, manifest, options)
	if err != nil {
		t.Fatal(err)
	}
	if len(packet.Files) == 0 || packet.Files[0].Content == "" {
		t.Fatal("workspace content was not loaded")
	}
	text, err := packet.Markdown()
	if err != nil {
		t.Fatal(err)
	}
	if text == "" {
		t.Fatal("empty packet")
	}
}

func TestBuildRejectsMissingMandatoryInput(t *testing.T) {
	root, model, manifest, options := fixture(t)
	if err := os.Remove(filepath.Join(options.Workspace, "0A-PREFLIGHT.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := Build(root, model, manifest, options); err == nil {
		t.Fatal("expected missing-input rejection")
	}
}

func TestBuildRejectsWorkspaceEscape(t *testing.T) {
	root, model, manifest, options := fixture(t)
	options.AdditionalReadTargets = []string{"../../../etc/passwd"}
	if _, err := Build(root, model, manifest, options); err == nil {
		t.Fatal("expected traversal rejection")
	}
}

func TestBuildRejectsSymlinkEscape(t *testing.T) {
	root, model, manifest, options := fixture(t)
	if err := os.Symlink("/etc/passwd", filepath.Join(options.Workspace, "escape")); err != nil {
		t.Fatal(err)
	}
	options.AdditionalReadTargets = []string{"escape"}
	if _, err := Build(root, model, manifest, options); err == nil {
		t.Fatal("expected symlink escape rejection")
	}
}

func fixture(t *testing.T) (string, *protocol.Model, *routing.Manifest, Options) {
	t.Helper()
	root := filepath.Join("..", "..")
	model, err := protocol.Load(filepath.Join(root, "protocol.md"))
	if err != nil {
		t.Fatal(err)
	}
	manifest, err := routing.Load(filepath.Join(root, "skill", "modes.json"))
	if err != nil {
		t.Fatal(err)
	}
	workspaceRoot := filepath.Join(t.TempDir(), ".extracted")
	if err := workspace.Init(workspaceRoot); err != nil {
		t.Fatal(err)
	}
	options := Options{Mode: "Preflight", SystemNamespace: "example", Iteration: "ALFA", InvocationID: "INV-EXAMPLE", InvocationScope: "bootstrap", EnvironmentSnapshot: "dev/snapshot-1", Workspace: workspaceRoot, AllowedWriteTargets: []string{"0A-PREFLIGHT.md", "0G-DECONSTRUCTION-STATE.md"}, ForbiddenWriteTargets: []string{"all other workspace paths"}}
	return root, model, manifest, options
}

func TestBuildRejectsPersonaDirectoryPrefixMismatch(t *testing.T) {
	root, model, manifest, options := fixture(t)
	options.Mode = "Discovery"
	options.Persona = "Candidate-Web"
	options.PersonaPrefix = "PAY"
	options.PersonaDirectory = "CND-candidate-web"
	options.Cluster = "CLU-EXAMPLE"
	options.Track = "Normal"
	options.POV = "POV-2"
	options.POVFile = "personas/CND-candidate-web/02-BACKEND-POV.md"
	options.MaxTraversalDepth = 3
	options.AdditionalReadTargets = []string{"20-TRACEABILITY.md"}
	if _, err := Build(root, model, manifest, options); err == nil {
		t.Fatal("expected persona directory prefix mismatch")
	}
}
