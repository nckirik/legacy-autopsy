package store

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Apply executes DDL against a SQLite database using the operator-provided
// sqlite3 CLI. It fails closed and performs no partial-silent recovery.
func Apply(dbPath, ddl string) error {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		return fmt.Errorf("store: sqlite3 CLI not found on PATH: %w", err)
	}
	cmd := exec.Command("sqlite3", "--bail", dbPath)
	cmd.Stdin = strings.NewReader(ddl)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("store: apply schema: %v: %s", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}
