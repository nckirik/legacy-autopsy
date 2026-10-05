package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/artifact"
)

// specArtifact validates executor-authored artifact payloads and verifies a
// handoff bundle directory against its outer manifest.
func specArtifact(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: legacy-autopsy spec artifact <validate|verify> [flags]")
	}
	op := args[0]
	switch op {
	case "validate", "verify":
	default:
		return fmt.Errorf("spec artifact operation %q is unsupported; supported operations: validate, verify", op)
	}
	set := flags("spec artifact "+op, stderr)
	repo := set.String("repo", ".", "repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	field := set.String("field", "", "artifact FIELD ID (validate)")
	payloadPath := set.String("payload", "", "payload file (validate)")
	artifactType := set.String("type", "", "artifact type, e.g. OUTER-BUNDLE-MANIFEST")
	artifactPath := set.String("path", "", "normalized artifact path")
	dir := set.String("dir", "", "bundle directory (verify)")
	if err := set.Parse(args[1:]); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	switch op {
	case "validate":
		if *field == "" || *payloadPath == "" {
			return fmt.Errorf("spec artifact validate requires --field and --payload")
		}
		payload, err := os.ReadFile(*payloadPath)
		if err != nil {
			return err
		}
		if err := artifact.ValidatePayload(res.EIR, *field, payload); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "valid %s payload\n", *field)
		if *artifactType != "" && *artifactPath != "" {
			fingerprint, err := artifact.PayloadFingerprint(*artifactType, *artifactPath, payload)
			if err != nil {
				return err
			}
			fmt.Fprintf(stdout, "payload fingerprint: %s\n", fingerprint)
		}
		return nil
	case "verify":
		if *artifactType == "" || *artifactPath == "" || *dir == "" {
			return fmt.Errorf("spec artifact verify requires --type, --path, and --dir")
		}
		outer, err := os.ReadFile(filepath.Join(*dir, filepath.FromSlash(*artifactPath)))
		if err != nil {
			return err
		}
		members, err := readMembers(*dir)
		if err != nil {
			return err
		}
		if err := artifact.VerifyBundle(*artifactType, *artifactPath, outer, members); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "bundle verified: %d member(s)\n", len(members))
		return nil
	}
	return nil
}

func readMembers(dir string) ([]artifact.Member, error) {
	var members []artifact.Member
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		members = append(members, artifact.Member{Path: filepath.ToSlash(relative), Bytes: data})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(members, func(i, j int) bool { return members[i].Path < members[j].Path })
	return members, nil
}
