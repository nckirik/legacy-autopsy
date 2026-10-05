package artifact

import (
	"fmt"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	md "github.com/nckirik/legacy-autopsy/internal/markdown"
)

// VerifyBundle is the deterministic §15.5 pilot check set: outer-manifest
// payload validity, package-member set and fingerprint recomputation, and direct
// snapshot equality across every snapshot-bearing member.
func VerifyBundle(artifactType, path string, outer []byte, members []Member) error {
	var failures []string
	if err := validateOuterPayload(outer); err != nil {
		failures = append(failures, err.Error())
	}
	declared, err := declaredMembers(outer)
	if err != nil {
		failures = append(failures, err.Error())
	}
	actual, err := PackageMembers(members)
	if err != nil {
		return err
	}
	if declared != nil {
		if err := compareMemberSets(declared, actual, path); err != nil {
			failures = append(failures, err.Error())
		}
	}
	if err := verifySnapshotEquality(members); err != nil {
		failures = append(failures, err.Error())
	}
	_, err = PayloadFingerprint(artifactType, path, outer)
	if err != nil {
		failures = append(failures, err.Error())
	}
	if len(failures) > 0 {
		return fmt.Errorf("artifact: bundle verification failed: %s", strings.Join(failures, "; "))
	}
	return nil
}

// compareMemberSets compares the declared table (which excludes the outer
// manifest itself) with the actual bundle member set.
func compareMemberSets(declared map[string]string, actual []ManifestEntry, outerPath string) error {
	comparable := make([]ManifestEntry, 0, len(actual))
	for _, entry := range actual {
		if entry.Path == outerPath {
			continue
		}
		comparable = append(comparable, entry)
	}
	if len(declared) != len(comparable) {
		return fmt.Errorf("member set mismatch: manifest declares %d, bundle has %d", len(declared), len(comparable))
	}
	for _, entry := range comparable {
		want, ok := declared[entry.Path]
		if !ok {
			return fmt.Errorf("bundle member %q is not declared in the outer manifest", entry.Path)
		}
		if want != entry.Fingerprint {
			return fmt.Errorf("member %q fingerprint mismatch", entry.Path)
		}
	}
	return nil
}

// declaredMembers parses the outer manifest Package Members table.
func declaredMembers(outer []byte) (map[string]string, error) {
	parsed, err := md.Parse(string(outer))
	if err != nil {
		return nil, fmt.Errorf("outer manifest: %w", err)
	}
	inTable := false
	header := true
	members := map[string]string{}
	for _, node := range parsed.Nodes {
		if node.Kind == md.Field {
			inTable = capabilities.ColumnName(node.Label) == "package-members"
			header = true
			continue
		}
		if !inTable || node.Kind != md.TableRow {
			continue
		}
		cells := splitRow(node.Raw)
		if header {
			header = false
			continue
		}
		if isSeparator(cells) {
			continue
		}
		if len(cells) != 3 {
			return nil, fmt.Errorf("package members row has %d cells, want 3", len(cells))
		}
		members[strings.TrimSpace(cells[1])] = strings.TrimSpace(cells[2])
	}
	if len(members) == 0 {
		return nil, fmt.Errorf("outer manifest has no Package Members table")
	}
	return members, nil
}

// verifySnapshotEquality requires every snapshot-bearing member to declare the
// same System / Protocol / Snapshot value.
func verifySnapshotEquality(members []Member) error {
	snapshot := ""
	for _, member := range members {
		parsed, err := md.Parse(string(member.Bytes))
		if err != nil {
			continue
		}
		for _, node := range parsed.Nodes {
			if node.Kind != md.Field || capabilities.ColumnName(node.Label) != "system-protocol-snapshot" {
				continue
			}
			if snapshot == "" {
				snapshot = node.Value
			} else if snapshot != node.Value {
				return fmt.Errorf("snapshot mismatch in member %q", member.Path)
			}
		}
	}
	return nil
}

func validateOuterPayload(outer []byte) error {
	parsed, err := md.Parse(string(outer))
	if err != nil {
		return fmt.Errorf("outer manifest payload: %w", err)
	}
	for _, node := range parsed.Nodes {
		if node.Kind == md.Field && capabilities.ColumnName(node.Label) == "exit-e-status" {
			if strings.TrimSpace(node.Value) != "Passed" {
				return fmt.Errorf("outer manifest EXIT-E-STATUS is %q, want Passed", node.Value)
			}
			return nil
		}
	}
	return fmt.Errorf("outer manifest has no EXIT-E-STATUS field")
}

func splitRow(raw string) []string {
	trimmed := strings.TrimSpace(raw)
	trimmed = strings.TrimPrefix(trimmed, "|")
	trimmed = strings.TrimSuffix(trimmed, "|")
	parts := strings.Split(trimmed, "|")
	cells := make([]string, len(parts))
	for i, part := range parts {
		cells[i] = strings.TrimSpace(part)
	}
	return cells
}

func isSeparator(cells []string) bool {
	for _, cell := range cells {
		if cell == "" {
			continue
		}
		if strings.Trim(cell, ":-") != "" {
			return false
		}
	}
	return true
}
