package work

import (
	"fmt"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	md "github.com/nckirik/legacy-autopsy/internal/markdown"
)

// ResolvedRecord is one proposed record with its runtime-assigned ID.
type ResolvedRecord struct {
	RecordID string
	Payload  []byte
}

// ResolveRecords assigns record IDs per the unit contract: a single declared
// target, or the value of the declared ID field extracted from each payload.
// Duplicate IDs fail closed.
func ResolveRecords(doc *cdl.EIRDoc, unit Unit, p Proposal) ([]ResolvedRecord, error) {
	if len(p.Records) == 0 {
		return nil, fmt.Errorf("work: unit %s proposal has no records", unit.ID)
	}
	if unit.IDFromField == "" {
		if len(p.Records) != 1 {
			return nil, fmt.Errorf("work: unit %s targets a single record but proposed %d", unit.ID, len(p.Records))
		}
		return []ResolvedRecord{{RecordID: unit.RecordID, Payload: p.Records[0].Payload}}, nil
	}
	seen := map[string]bool{}
	resolved := make([]ResolvedRecord, 0, len(p.Records))
	for i, record := range p.Records {
		id, err := extractFieldValue(unit.IDFromField, record.Payload)
		if err != nil {
			return nil, fmt.Errorf("work: unit %s record %d: %w", unit.ID, i+1, err)
		}
		if seen[id] {
			return nil, fmt.Errorf("work: unit %s proposed duplicate record ID %q", unit.ID, id)
		}
		seen[id] = true
		resolved = append(resolved, ResolvedRecord{RecordID: id, Payload: record.Payload})
	}
	return resolved, nil
}

// PayloadText returns the declared field values of one payload keyed by the
// declared spec names. Duplicate labels fail closed.
func PayloadText(names []string, payload []byte) (map[string]string, error) {
	parsed, err := md.Parse(string(payload))
	if err != nil {
		return nil, err
	}
	byColumn := map[string]string{}
	for _, node := range parsed.Nodes {
		if node.Kind != md.Field {
			continue
		}
		column := capabilities.ColumnName(node.Label)
		if _, duplicate := byColumn[column]; duplicate {
			return nil, fmt.Errorf("payload repeats field %q", node.Label)
		}
		byColumn[column] = strings.TrimSpace(node.Value)
	}
	values := map[string]string{}
	for _, name := range names {
		if value, ok := byColumn[capabilities.ColumnName(name)]; ok {
			values[name] = value
		}
	}
	return values, nil
}

func extractFieldValue(fieldName string, payload []byte) (string, error) {
	parsed, err := md.Parse(string(payload))
	if err != nil {
		return "", err
	}
	want := capabilities.ColumnName(fieldName)
	for _, node := range parsed.Nodes {
		if node.Kind != md.Field {
			continue
		}
		if capabilities.ColumnName(node.Label) != want {
			continue
		}
		value := strings.TrimSpace(node.Value)
		if value == "" {
			return "", fmt.Errorf("ID field %q is empty", fieldName)
		}
		return value, nil
	}
	return "", fmt.Errorf("ID field %q is missing", fieldName)
}
