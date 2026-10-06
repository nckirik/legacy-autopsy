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
	default:
		return nil
	}
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
