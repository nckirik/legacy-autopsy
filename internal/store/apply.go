package store

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Apply executes DDL against a SQLite database using the operator-provided
// sqlite3 CLI. The DDL runs inside a single transaction, so a failure mid-script
// leaves no partially-created schema; it fails closed with no partial-silent
// recovery. The foreign-key pragma is set before the transaction begins.
func Apply(dbPath, ddl string) error {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		return fmt.Errorf("store: sqlite3 CLI not found on PATH: %w", err)
	}
	script := "PRAGMA foreign_keys = ON;\nBEGIN;\n" + ddl + "\nCOMMIT;\n"
	cmd := exec.Command("sqlite3", "--bail", dbPath)
	cmd.Stdin = strings.NewReader(script)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("store: apply schema: %v: %s", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}
