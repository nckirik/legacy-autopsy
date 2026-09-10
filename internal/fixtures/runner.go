// Package fixtures runs explicitly registered bootstrap fixture groups.
// These cases do not constitute the complete Protocol §19.2 conformance suite.
package fixtures

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/canonical"
	"github.com/nckirik/legacy-autopsy/internal/identity"
	"github.com/nckirik/legacy-autopsy/internal/protocol"
	"github.com/nckirik/legacy-autopsy/internal/workspace"
)

type Case struct {
	ID                 string            `json:"id"`
	ProtocolSection    string            `json:"protocol_section"`
	Group              string            `json:"group"`
	Operation          string            `json:"operation"`
	Input              map[string]string `json:"input"`
	Expected           string            `json:"expected"`
	ExpectError        bool              `json:"expect_error"`
	ExpectedDiagnostic string            `json:"expected_diagnostic"`
}

type Summary struct {
	Passed      int
	Failed      int
	Unsupported []string
}

type diagnosticError struct {
	Code string
	Err  error
}

func (e *diagnosticError) Error() string { return e.Code + ": " + e.Err.Error() }
func (e *diagnosticError) Unwrap() error { return e.Err }

var operationGroups = map[string]string{
	"protocol-version":    "protocol-discovery",
	"protocol-mode-count": "protocol-discovery",
	"generate-id":         "identity",
	"normalize-path":      "path-normalization",
	"basic-markdown-hash": "basic-markdown-hash",
	"workspace-skeleton":  "workspace-skeleton",
}

var unsupported = []string{
	"Part 19.2 persona workspace ownership and atomic moves",
	"Part 19.2 complete context-qualified enum registry",
	"Part 19.2 invocation ownership and cold resume",
	"Part 19.2 ticket escalation and honest unresolved coverage",
	"Part 19.2 dead-code and semantic-predicate closure",
	"Part 19.2 deterministic iteration accounting",
	"Part 19.2 complete PRF semantic hashing",
	"Part 19.2 complete HBK identity and move rules",
	"Part 19.2 Profile Synchronization transitive closure",
	"Part 19.2 record-versus-artifact hashing",
	"Part 19.2 normative canonicalization and exact exclusions",
	"Part 19.2 final envelope and package-member integrity",
	"Part 19.2 packaging schemas, gate completeness, and evidence bindings",
	"Part 19.2 snapshot-consistent acyclic Exit E and final verification",
}

func Run(repoRoot string) (Summary, error) {
	cases, err := loadCases(repoRoot)
	if err != nil {
		return Summary{}, err
	}
	summary := Summary{Unsupported: append([]string(nil), unsupported...)}
	var failures []string
	for _, fixture := range cases {
		actual, runErr := execute(repoRoot, fixture)
		if fixture.ExpectError {
			var diagnostic *diagnosticError
			switch {
			case runErr == nil:
				summary.Failed++
				failures = append(failures, fixture.ID+": expected rejection")
			case !errors.As(runErr, &diagnostic):
				summary.Failed++
				failures = append(failures, fixture.ID+": rejection lacked a stable diagnostic class")
			case diagnostic.Code != fixture.ExpectedDiagnostic:
				summary.Failed++
				failures = append(failures, fmt.Sprintf("%s: expected diagnostic %q, got %q", fixture.ID, fixture.ExpectedDiagnostic, diagnostic.Code))
			default:
				summary.Passed++
			}
			continue
		}
		if runErr != nil {
			summary.Failed++
			failures = append(failures, fixture.ID+": "+runErr.Error())
		} else if actual != fixture.Expected {
			summary.Failed++
			failures = append(failures, fmt.Sprintf("%s: expected %q, got %q", fixture.ID, fixture.Expected, actual))
		} else {
			summary.Passed++
		}
	}
	if len(failures) > 0 {
		return summary, fmt.Errorf("fixture failures:\n- %s", strings.Join(failures, "\n- "))
	}
	return summary, nil
}

func loadCases(repoRoot string) ([]Case, error) {
	model, err := protocol.Load(filepath.Join(repoRoot, "protocol.md"))
	if err != nil {
		return nil, err
	}
	var cases []Case
	seen := map[string]struct{}{}
	for _, polarity := range []string{"positive", "negative"} {
		dir := filepath.Join(repoRoot, "fixtures", polarity)
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("read %s fixtures: %w", polarity, err)
		}
		for _, entry := range entries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
				continue
			}
			data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
			if err != nil {
				return nil, err
			}
			var fileCases []Case
			if err := json.Unmarshal(data, &fileCases); err != nil {
				return nil, fmt.Errorf("parse %s: %w", entry.Name(), err)
			}
			for _, fixture := range fileCases {
				if fixture.ID == "" || fixture.ProtocolSection == "" || fixture.Group == "" || fixture.Operation == "" {
					return nil, fmt.Errorf("%s contains incomplete fixture metadata", entry.Name())
				}
				if !model.HasSection(fixture.ProtocolSection) {
					return nil, fmt.Errorf("%s references missing or duplicate protocol section %q", fixture.ID, fixture.ProtocolSection)
				}
				expectedGroup, supported := operationGroups[fixture.Operation]
				if !supported || fixture.Group != expectedGroup {
					return nil, fmt.Errorf("%s has unregistered operation/group pair %q/%q", fixture.ID, fixture.Operation, fixture.Group)
				}
				if fixture.ExpectError && fixture.ExpectedDiagnostic == "" {
					return nil, fmt.Errorf("%s lacks an expected diagnostic class", fixture.ID)
				}
				if _, duplicate := seen[fixture.ID]; duplicate {
					return nil, fmt.Errorf("duplicate fixture ID %q", fixture.ID)
				}
				seen[fixture.ID] = struct{}{}
				cases = append(cases, fixture)
			}
		}
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("no implemented fixtures found")
	}
	sort.Slice(cases, func(i, j int) bool { return cases[i].ID < cases[j].ID })
	return cases, nil
}

func execute(repoRoot string, fixture Case) (string, error) {
	switch fixture.Operation {
	case "normalize-path":
		result, err := identity.NormalizeRelativePath(fixture.Input["path"])
		return result, classified("PATH_INVALID", err)
	case "generate-id":
		result, err := identity.Generate(fixture.Input["type"], fixture.Input["kind"], fixture.Input["namespace"], fixture.Input["owner"], fixture.Input["discriminator"])
		return result.ID, classified("IDENTITY_INVALID", err)
	case "protocol-version":
		model, err := protocol.Load(filepath.Join(repoRoot, "protocol.md"))
		if err != nil {
			return "", err
		}
		if model.Version != fixture.Input["required"] {
			return "", &diagnosticError{Code: "PROTOCOL_VERSION_MISMATCH", Err: fmt.Errorf("required version %s, found %s", fixture.Input["required"], model.Version)}
		}
		return model.Version, nil
	case "protocol-mode-count":
		model, err := protocol.Load(filepath.Join(repoRoot, "protocol.md"))
		if err != nil {
			return "", err
		}
		required, err := strconv.Atoi(fixture.Input["required"])
		if err != nil {
			return "", err
		}
		if len(model.Modes) != required {
			return "", &diagnosticError{Code: "PROTOCOL_MODE_COUNT_MISMATCH", Err: fmt.Errorf("required %d modes, found %d", required, len(model.Modes))}
		}
		return strconv.Itoa(len(model.Modes)), nil
	case "basic-markdown-hash":
		result, err := canonical.BasicMarkdownFingerprint(fixture.Input["markdown"])
		return result, classified("MARKDOWN_INVALID", err)
	case "workspace-skeleton":
		parent, err := os.MkdirTemp("", "legacy-autopsy-fixture-*")
		if err != nil {
			return "", err
		}
		defer os.RemoveAll(parent)
		root := filepath.Join(parent, ".extracted")
		if err := workspace.Init(root); err != nil {
			return "", err
		}
		if remove := fixture.Input["remove"]; remove != "" {
			if err := os.Remove(filepath.Join(root, filepath.FromSlash(remove))); err != nil {
				return "", err
			}
		}
		if rewrite := fixture.Input["rewrite"]; rewrite != "" {
			if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(rewrite)), []byte(fixture.Input["content"]), 0o644); err != nil {
				return "", err
			}
		}
		if err := workspace.Check(root); err != nil {
			return "", &diagnosticError{Code: "WORKSPACE_INVALID", Err: err}
		}
		return "valid", nil
	default:
		return "", fmt.Errorf("unsupported fixture operation %q", fixture.Operation)
	}
}

func classified(code string, err error) error {
	if err == nil {
		return nil
	}
	return &diagnosticError{Code: code, Err: err}
}
