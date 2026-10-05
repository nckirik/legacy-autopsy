package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// queryRows runs a SQL query through the operator-provided sqlite3 CLI and
// returns the rows as JSON-decoded maps. Numbers stay json.Number so callers
// control integer formatting. No Go SQLite dependency is introduced.
func queryRows(dbPath, query string) ([]map[string]any, error) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		return nil, fmt.Errorf("store: sqlite3 CLI not found on PATH: %w", err)
	}
	cmd := exec.Command("sqlite3", "--json", dbPath, query)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("store: query: %v: %s", err, strings.TrimSpace(stderr.String()))
	}
	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" || trimmed == "[]" {
		return nil, nil
	}
	var rows []map[string]any
	dec := json.NewDecoder(strings.NewReader(trimmed))
	dec.UseNumber()
	if err := dec.Decode(&rows); err != nil {
		return nil, fmt.Errorf("store: parse query rows: %w", err)
	}
	return rows, nil
}
