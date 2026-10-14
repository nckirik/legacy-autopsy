package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/artifact"
	"github.com/nckirik/legacy-autopsy/internal/handoff"
)

// specArtifact validates executor-authored artifact payloads and verifies a
// handoff bundle directory against its outer manifest.
func specArtifact(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: legacy-autopsy spec artifact <validate|verify> [flags]")
	}
	op := args[0]
	switch op {
	case "validate", "verify", "chain":
	default:
		return fmt.Errorf("spec artifact operation %q is unsupported; supported operations: validate, verify, chain", op)
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
	candidate := set.String("candidate", "", "EXIT-E-CANDIDATE-REPORT payload (chain)")
	manifest := set.String("manifest", "", "CANDIDATE-PAYLOAD-MANIFEST payload (chain)")
	readiness := set.String("readiness", "", "EXIT-E-CONTENT-READINESS-REPORT payload (chain)")
	certificate := set.String("certificate", "", "SCOPE-CERTIFICATE payload (chain)")
	outer := set.String("outer", "", "OUTER-BUNDLE-MANIFEST payload (chain)")
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
	case "chain":
		paths := []struct {
			flag, artifactType, path string
		}{
			{"candidate", "EXIT-E-CANDIDATE-REPORT", handoff.CandidateReportDefaultPath},
			{"manifest", "CANDIDATE-PAYLOAD-MANIFEST", handoff.CandidateManifestDefaultPath},
			{"readiness", "EXIT-E-CONTENT-READINESS-REPORT", handoff.ContentReadinessDefaultPath},
			{"certificate", "SCOPE-CERTIFICATE", handoff.ScopeCertificateDefaultPath},
			{"outer", "OUTER-BUNDLE-MANIFEST", handoff.OuterManifestDefaultPath},
		}
		values := map[string]string{"candidate": *candidate, "manifest": *manifest, "readiness": *readiness, "certificate": *certificate, "outer": *outer}
		var stages []artifact.ChainStage
		for _, entry := range paths {
			file := values[entry.flag]
			if file == "" {
				return fmt.Errorf("spec artifact chain requires --%s", entry.flag)
			}
			payload, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			stages = append(stages, artifact.ChainStage{Type: entry.artifactType, Path: entry.path, Payload: payload})
		}
		if err := artifact.VerifyChain(res.EIR, stages); err != nil {
			return err
		}
		fmt.Fprintln(stdout, "chain verified: candidate -> manifest -> readiness -> certificate -> outer")
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
