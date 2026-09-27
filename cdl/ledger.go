package cdl

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// IDRecord is one stable identity entry.
type IDRecord struct {
	ID     string `json:"id"`
	Kind   string `json:"kind"`
	Status string `json:"status"`
}

// LedgerSource is the identity contribution of one source file.
type LedgerSource struct {
	Path              string     `json:"path"`
	SourceFingerprint string     `json:"source-fingerprint"`
	IDs               []IDRecord `json:"ids"`
}

// Ledger is the committed identity ledger. Section-owned identities are globally
// unique across every source: §4.1 requires global uniqueness, not per-file
// uniqueness. Document-global declarations are linked at assembly.
type Ledger struct {
	Format      int            `json:"ledger-format"`
	Generator   string         `json:"generator"`
	GeneratedAt string         `json:"generated-at"`
	Sources     []LedgerSource `json:"sources"`
}

// LedgerProvenance identifies how a ledger was generated.
type LedgerProvenance struct {
	Generator   string
	GeneratedAt string
}

// LoadLedger reads and validates a committed identity ledger.
func LoadLedger(path string) (map[string]IDRecord, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var l Ledger
	if err := json.Unmarshal(b, &l); err != nil {
		return nil, fmt.Errorf("identity ledger: %w", err)
	}
	out := map[string]IDRecord{}
	for _, src := range l.Sources {
		for _, r := range src.IDs {
			if existing, dup := out[r.ID]; dup {
				return nil, fmt.Errorf("identity ledger: id %s is duplicated (%s and %s)", r.ID, existing.Kind, r.Kind)
			}
			out[r.ID] = r
		}
	}
	return out, nil
}

// WriteLedger writes a deterministic, provenance-bearing identity ledger with
// globally unique identities across sources and sources sorted by path.
func WriteLedger(path string, sources []LedgerSource, prov LedgerProvenance) error {
	sorted := append([]LedgerSource(nil), sources...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Path < sorted[j].Path })
	seen := map[string]string{}
	for si := range sorted {
		ids := append([]IDRecord(nil), sorted[si].IDs...)
		sort.Slice(ids, func(i, j int) bool { return ids[i].ID < ids[j].ID })
		for _, id := range ids {
			if other, dup := seen[id.ID]; dup {
				return fmt.Errorf("identity ledger: id %s is duplicated between %s and %s", id.ID, other, sorted[si].Path)
			}
			seen[id.ID] = sorted[si].Path
		}
		sorted[si].IDs = ids
	}
	ledger := Ledger{
		Format:      2,
		Generator:   prov.Generator,
		GeneratedAt: prov.GeneratedAt,
		Sources:     sorted,
	}
	b, err := json.MarshalIndent(ledger, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}
