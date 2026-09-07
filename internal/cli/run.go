// Package cli wires commands to provider-neutral deterministic services.
package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/contextpacket"
	"github.com/nckirik/legacy-autopsy/internal/fixtures"
	"github.com/nckirik/legacy-autopsy/internal/identity"
	"github.com/nckirik/legacy-autopsy/internal/protocol"
	"github.com/nckirik/legacy-autopsy/internal/routing"
	"github.com/nckirik/legacy-autopsy/internal/workspace"
)

func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		usage(stderr)
		return 2
	}
	var err error
	switch args[0] {
	case "doctor":
		err = doctor(args[1:], stdout, stderr)
	case "init":
		err = initWorkspace(args[1:], stdout, stderr)
	case "protocol":
		err = protocolCommand(args[1:], stdout, stderr)
	case "workspace":
		err = workspaceCommand(args[1:], stdout, stderr)
	case "id":
		err = idCommand(args[1:], stdout, stderr)
	case "context":
		err = contextCommand(args[1:], stdout, stderr)
	case "validate":
		err = validateCommand(args[1:], stdout, stderr)
	case "help", "-h", "--help":
		usage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		usage(stderr)
		return 2
	}
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	return 0
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "legacy-autopsy <doctor|init|protocol check|workspace check|id|context|validate fixtures|validate routing>")
}

func flags(name string, stderr io.Writer) *flag.FlagSet {
	set := flag.NewFlagSet(name, flag.ContinueOnError)
	set.SetOutput(stderr)
	return set
}

func doctor(args []string, stdout, stderr io.Writer) error {
	set := flags("doctor", stderr)
	repo := set.String("repo", ".", "repository root")
	if err := set.Parse(args); err != nil {
		return err
	}
	for _, path := range []string{"protocol.md", "skill/SKILL.md", "skill/modes.json"} {
		if _, err := os.Stat(filepath.Join(*repo, filepath.FromSlash(path))); err != nil {
			return fmt.Errorf("required repository file %s: %w", path, err)
		}
	}
	model, manifest, err := loadRouting(*repo)
	if err != nil {
		return err
	}
	if err := routing.Validate(*repo, manifest, model); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "usable: protocol v%s, %d invocation modes, %s\n", model.Version, len(model.Modes), runtime.Version())
	fmt.Fprintln(stdout, "supported: protocol routing, IDs, path normalization, workspace skeleton/check, bounded context packets, bootstrap fixtures")
	fmt.Fprintln(stdout, "unsupported: full Protocol Part 19 conformance, Exit A, Exit E, acquisition, synthesis, confirmation, and packaging")
	return nil
}

func initWorkspace(args []string, stdout, stderr io.Writer) error {
	set := flags("init", stderr)
	root := set.String("workspace", ".extracted", "workspace root")
	if err := set.Parse(args); err != nil {
		return err
	}
	if err := workspace.Init(*root); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "initialized non-fabricated Protocol v4 workspace at %s\n", *root)
	return nil
}

func protocolCommand(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] != "check" {
		return fmt.Errorf("usage: legacy-autopsy protocol check [flags]")
	}
	set := flags("protocol check", stderr)
	repo := set.String("repo", ".", "repository root")
	protocolPath := set.String("protocol", "protocol.md", "protocol path relative to repo")
	manifestPath := set.String("manifest", "skill/modes.json", "mode manifest relative to repo")
	if err := set.Parse(args[1:]); err != nil {
		return err
	}
	model, err := protocol.Load(filepath.Join(*repo, filepath.FromSlash(*protocolPath)))
	if err != nil {
		return err
	}
	manifest, err := routing.Load(filepath.Join(*repo, filepath.FromSlash(*manifestPath)))
	if err != nil {
		return err
	}
	if err := routing.Validate(*repo, manifest, model); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "protocol check passed: v%s, %d modes, %d projections, fingerprint %s\n", model.Version, len(model.Modes), len(manifest.Projections), model.Fingerprint)
	return nil
}

func workspaceCommand(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] != "check" {
		return fmt.Errorf("usage: legacy-autopsy workspace check [flags]")
	}
	set := flags("workspace check", stderr)
	root := set.String("workspace", ".extracted", "workspace root")
	if err := set.Parse(args[1:]); err != nil {
		return err
	}
	if err := workspace.Check(*root); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "workspace structure passed: %s (%d required Markdown files)\n", *root, len(workspace.RequiredFiles()))
	return nil
}

func idCommand(args []string, stdout, stderr io.Writer) error {
	set := flags("id", stderr)
	recordType := set.String("type", "", "protocol ID type")
	kind := set.String("kind", "", "SRC kind")
	namespace := set.String("namespace", "", "system namespace")
	owner := set.String("owner", "", "normalized owner coordinate")
	discriminator := set.String("discriminator", "", "semantic discriminator")
	pathValue := set.String("normalize-path", "", "normalize a relative path instead of generating an ID")
	if err := set.Parse(args); err != nil {
		return err
	}
	if *pathValue != "" {
		result, err := identity.NormalizeRelativePath(*pathValue)
		if err != nil {
			return err
		}
		fmt.Fprintln(stdout, result)
		return nil
	}
	result, err := identity.Generate(*recordType, *kind, *namespace, *owner, *discriminator)
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "%s\ncanonical-key: %s\n", result.ID, result.CanonicalKey)
	return nil
}

type repeated []string

func (r *repeated) String() string         { return strings.Join(*r, ",") }
func (r *repeated) Set(value string) error { *r = append(*r, value); return nil }

func contextCommand(args []string, stdout, stderr io.Writer) error {
	set := flags("context", stderr)
	repo := set.String("repo", ".", "repository root")
	options := contextpacket.Options{}
	set.StringVar(&options.Mode, "mode", "", "exact invocation mode")
	set.StringVar(&options.SystemNamespace, "system-namespace", "", "system namespace")
	set.StringVar(&options.Iteration, "iteration", "", "canonical ALFA..ZULU iteration")
	set.StringVar(&options.InvocationID, "invocation-id", "", "globally unique invocation ID")
	set.StringVar(&options.InvocationScope, "scope", "", "exact invocation scope")
	set.StringVar(&options.Persona, "persona", "", "persona")
	set.StringVar(&options.PersonaPrefix, "persona-prefix", "", "persona prefix")
	set.StringVar(&options.Cluster, "cluster", "", "entry cluster")
	set.StringVar(&options.Track, "track", "", "traversal track")
	set.StringVar(&options.POV, "pov", "", "active POV")
	set.StringVar(&options.POVFile, "pov-file", "", "active POV file relative to workspace")
	set.StringVar(&options.QuestionLedger, "question-ledger", "", "active question ledger relative to workspace")
	set.StringVar(&options.EnvironmentSnapshot, "environment-snapshot", "", "environment/snapshot binding")
	set.StringVar(&options.HumanHatchAction, "human-hatch-action", "", "Human Hatch subtype")
	set.StringVar(&options.Workspace, "workspace", ".extracted", "workspace root")
	set.IntVar(&options.MaxTraversalDepth, "max-depth", 0, "maximum traversal depth")
	var reads, writes, forbidden repeated
	set.Var(&reads, "read-target", "additional mandatory workspace read target (repeatable)")
	set.Var(&writes, "write-target", "allowed write target (repeatable)")
	set.Var(&forbidden, "forbidden-write-target", "forbidden write target (repeatable)")
	if err := set.Parse(args); err != nil {
		return err
	}
	options.AdditionalReadTargets = reads
	options.AllowedWriteTargets = writes
	options.ForbiddenWriteTargets = forbidden
	model, manifest, err := loadRouting(*repo)
	if err != nil {
		return err
	}
	if err := routing.Validate(*repo, manifest, model); err != nil {
		return err
	}
	packet, err := contextpacket.Build(*repo, model, manifest, options)
	if err != nil {
		return err
	}
	text, err := packet.Markdown()
	if err != nil {
		return err
	}
	fmt.Fprint(stdout, text)
	return nil
}

func validateCommand(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: legacy-autopsy validate <fixtures|routing>")
	}
	set := flags("validate "+args[0], stderr)
	repo := set.String("repo", ".", "repository root")
	implemented := set.Bool("implemented", false, "run implemented fixture groups")
	if err := set.Parse(args[1:]); err != nil {
		return err
	}
	switch args[0] {
	case "routing":
		model, manifest, err := loadRouting(*repo)
		if err != nil {
			return err
		}
		if err := routing.Validate(*repo, manifest, model); err != nil {
			return err
		}
		fmt.Fprintln(stdout, "routing validation passed")
		return nil
	case "fixtures":
		_ = implemented
		summary, err := fixtures.Run(*repo)
		fmt.Fprintf(stdout, "implemented bootstrap fixtures: %d passed, %d failed\n", summary.Passed, summary.Failed)
		fmt.Fprintln(stdout, "unsupported conformance groups:")
		for _, item := range summary.Unsupported {
			fmt.Fprintln(stdout, "-", item)
		}
		return err
	default:
		return fmt.Errorf("validator %q is unsupported; supported validators: fixtures, routing", args[0])
	}
}

func loadRouting(repo string) (*protocol.Model, *routing.Manifest, error) {
	model, err := protocol.Load(filepath.Join(repo, "protocol.md"))
	if err != nil {
		return nil, nil, err
	}
	manifest, err := routing.Load(filepath.Join(repo, "skill", "modes.json"))
	if err != nil {
		return nil, nil, err
	}
	return model, manifest, nil
}
