package certify

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	clusterIDPattern   = regexp.MustCompile(`^[A-Z][A-Z0-9-]{1,15}$`)
	sourceIDToken      = regexp.MustCompile(`^(None|SRC-[A-Za-z0-9][A-Za-z0-9-]*)$`)
	personaOwnerPrefix = regexp.MustCompile(`^[A-Z][A-Z0-9]{1,7}$`)
)

// ReferenceLookup reports whether a record exists.
type ReferenceLookup func(fieldID, recordID string) (bool, error)

// ValidateReferences enforces cross-record references that cannot be expressed
// as column types. values maps declared field names to their payload text.
// Unknown FIELD IDs pass through.
func ValidateReferences(fieldID string, values map[string]string, lookup ReferenceLookup) error {
	switch fieldID {
	case "ENTRY-CLUSTER":
		return validateEntryCluster(values, lookup)
	case "FRONTIER-RECORD":
		return validateFrontier(values, lookup)
	default:
		return nil
	}
}

func validateFrontier(values map[string]string, lookup ReferenceLookup) error {
	for _, name := range []string{"from-src-cmp", "target-src"} {
		value := strings.TrimSpace(values[name])
		if value == "" || value == "None" || value == "root" {
			continue
		}
		exists, err := lookup("SOURCE-INVENTORY-RECORD", value)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("FRONTIER-RECORD.%s %q does not reference a committed source unit", name, value)
		}
	}
	return nil
}

func validateEntryCluster(values map[string]string, lookup ReferenceLookup) error {
	if id := strings.TrimSpace(values["cluster-id"]); id != "" && !clusterIDPattern.MatchString(id) {
		return fmt.Errorf("ENTRY-CLUSTER.cluster-id %q must match [A-Z][A-Z0-9-]{1,15}", id)
	}
	owner := strings.TrimSpace(values["canonical-owner-persona"])
	if owner != "" {
		prefix := strings.SplitN(owner, "/", 2)[0]
		if !personaOwnerPrefix.MatchString(prefix) && prefix != "SHARED" {
			return fmt.Errorf("ENTRY-CLUSTER.canonical-owner-persona %q must start with a persona prefix", owner)
		}
		if prefix != "SHARED" {
			exists, err := lookup("PERSONA-REGISTRY", prefix)
			if err != nil {
				return err
			}
			if !exists {
				return fmt.Errorf("ENTRY-CLUSTER.canonical-owner-persona %q does not reference a committed persona prefix", owner)
			}
		}
	}
	if raw := strings.TrimSpace(values["source-ids"]); raw != "" {
		for _, token := range strings.Split(raw, ",") {
			token = strings.TrimSpace(token)
			if token == "" {
				continue
			}
			if !sourceIDToken.MatchString(token) {
				return fmt.Errorf("ENTRY-CLUSTER.source-ids token %q must be None or an SRC-* ID", token)
			}
		}
	}
	return nil
}

// RecordReader returns the declared field values of one stored record.
type RecordReader func(fieldID, recordID string) (map[string]string, error)

// ValidateCrossReferences enforces cross-record evidence bindings: a claim must
// reference a committed source-inventory record whose fingerprint and
// coordinate match the claim's anchor and source fingerprint.
func ValidateCrossReferences(fieldID string, values map[string]string, read RecordReader) error {
	switch fieldID {
	case "CLAIM":
		return validateClaim(values, read)
	case "FRONTIER-RECORD":
		return validateFrontierCoordinates(values, read)
	default:
		return nil
	}
}

func validateFrontierCoordinates(values map[string]string, read RecordReader) error {
	target := strings.TrimSpace(values["target-src"])
	if target == "" || target == "None" {
		return nil
	}
	inventory, err := read("SOURCE-INVENTORY-RECORD", target)
	if err != nil {
		return nil // existence is enforced by the reference validator
	}
	coordinate := strings.TrimSpace(values["target-concrete-coordinate-or-resolution-rule"])
	sourceCoordinate := strings.TrimSpace(inventory["concrete-coordinate"])
	if coordinate != "" && sourceCoordinate != "" && coordinate != sourceCoordinate {
		return fmt.Errorf("FRONTIER-RECORD.target-concrete-coordinate %q does not match source unit %s coordinate %q", coordinate, target, sourceCoordinate)
	}
	return nil
}

func validateClaim(values map[string]string, read RecordReader) error {
	subject := firstToken(values["subject-id-record-version"])
	if subject == "" {
		return fmt.Errorf("CLAIM.subject-id-record-version is empty")
	}
	inventory, err := read("SOURCE-INVENTORY-RECORD", subject)
	if err != nil {
		return fmt.Errorf("CLAIM.subject-id-record-version %q does not reference a committed source unit: %w", subject, err)
	}
	fingerprint := strings.TrimSpace(values["source-fingerprint-snapshot"])
	inventoryVersion := strings.TrimSpace(inventory["source-fingerprint-snapshot-src-version"])
	inventoryFingerprint := firstTokenField(inventoryVersion)
	if fingerprint != "" && inventoryFingerprint != "" && !strings.HasPrefix(inventoryFingerprint, fingerprint) && !strings.HasPrefix(fingerprint, inventoryFingerprint) {
		return fmt.Errorf("CLAIM.source-fingerprint-snapshot %q does not match source unit %s fingerprint %q", fingerprint, subject, inventoryFingerprint)
	}
	coordinate := strings.TrimSpace(inventory["concrete-coordinate"])
	anchor := strings.TrimSpace(values["anchor"])
	if coordinate != "" && anchor != "" {
		if anchor != coordinate && !strings.HasPrefix(anchor, coordinate+":") && !strings.HasPrefix(anchor, coordinate+"#") {
			return fmt.Errorf("CLAIM.anchor %q must start at the source unit coordinate %q", anchor, coordinate)
		}
	}
	return nil
}

func firstToken(value string) string {
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return ""
	}
	return strings.TrimSuffix(fields[0], "/")
}

func firstTokenField(value string) string {
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}
