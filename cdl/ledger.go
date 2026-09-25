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

// Ledger is the committed identity ledger.
type Ledger struct {
	Format            int        `json:"ledger-format"`
	Generator         string     `json:"generator"`
	Source            string     `json:"source"`
	SourceFingerprint string     `json:"source-fingerprint"`
	GeneratedAt       string     `json:"generated-at"`
	IDs               []IDRecord `json:"ids"`
}

// LedgerProvenance identifies how a ledger was generated.
type LedgerProvenance struct {
	Generator         string
	Source            string
	SourceFingerprint string
	GeneratedAt       string
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
	for _, r := range l.IDs {
		if _, dup := out[r.ID]; dup {
			return nil, fmt.Errorf("identity ledger: duplicate id %s", r.ID)
		}
		out[r.ID] = r
	}
	return out, nil
}

// WriteLedger writes a deterministic, provenance-bearing identity ledger.
func WriteLedger(path string, ids []IDRecord, prov LedgerProvenance) error {
	sorted := append([]IDRecord(nil), ids...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ID < sorted[j].ID })
	ledger := Ledger{
		Format:            1,
		Generator:         prov.Generator,
		Source:            prov.Source,
		SourceFingerprint: prov.SourceFingerprint,
		GeneratedAt:       prov.GeneratedAt,
		IDs:               sorted,
	}
	b, err := json.MarshalIndent(ledger, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}
