package store

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
)

// Import restores a canonical dump into a fresh store database and returns the
// recomputed content fingerprint. The dump header must declare the current
// schema and dump formats; a non-empty target database fails closed; when an
// expected fingerprint is supplied it must match the recomputed content
// fingerprint exactly, so a modified dump is rejected before anything is
// applied. After applying, the live schema is verified against the
// EIR-derived schema.
func Import(dbPath string, doc *cdl.EIRDoc, dump []byte, expectFingerprint string) (string, error) {
	if dbPath == "" {
		return "", fmt.Errorf("store: import requires a target database")
	}
	text := string(dump)
	begin := strings.Index(text, "BEGIN;\n")
	if begin < 0 {
		return "", fmt.Errorf("store: dump has no canonical BEGIN block")
	}
	header := text[:begin]
	schemaVersion, dumpFormat, err := parseDumpHeader(header)
	if err != nil {
		return "", err
	}
	if dumpFormat != DumpVersion {
		return "", fmt.Errorf("store: dump format %d is not supported (want %d)", dumpFormat, DumpVersion)
	}
	if schemaVersion != SchemaVersion {
		return "", fmt.Errorf("store: dump schema version %d does not match the generator version %d", schemaVersion, SchemaVersion)
	}
	end := strings.Index(text[begin:], "COMMIT;\n")
	if end < 0 {
		return "", fmt.Errorf("store: dump has no canonical COMMIT block")
	}
	content := text[begin : begin+end+len("COMMIT;\n")]
	if strings.TrimSpace(text[begin+end+len("COMMIT;\n"):]) != "" {
		return "", fmt.Errorf("store: dump has content after the canonical COMMIT block")
	}
	fingerprint := cdl.Fingerprint([]byte(content))
	if expectFingerprint != "" && fingerprint != expectFingerprint {
		return "", fmt.Errorf("store: dump content fingerprint %s does not match the expected %s", fingerprint, expectFingerprint)
	}
	empty, err := isEmptyDatabase(dbPath)
	if err != nil {
		return "", err
	}
	if !empty {
		return "", fmt.Errorf("store: target database %s is not empty; import requires a fresh database", dbPath)
	}
	cmd := exec.Command("sqlite3", dbPath)
	cmd.Stdin = strings.NewReader(content)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("store: apply dump: %v: %s", err, strings.TrimSpace(stderr.String()))
	}
	stmts, err := Statements(doc)
	if err != nil {
		return "", err
	}
	if err := verifyLiveSchema(dbPath, stmts); err != nil {
		return "", err
	}
	return fingerprint, nil
}

func parseDumpHeader(header string) (int, int, error) {
	var schemaVersion, dumpFormat = -1, -1
	for _, line := range strings.Split(header, "\n") {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "-- schema-version:"):
			value, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "-- schema-version:")))
			if err != nil {
				return 0, 0, fmt.Errorf("store: invalid schema-version header %q", line)
			}
			schemaVersion = value
		case strings.HasPrefix(line, "-- dump-format:"):
			value, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "-- dump-format:")))
			if err != nil {
				return 0, 0, fmt.Errorf("store: invalid dump-format header %q", line)
			}
			dumpFormat = value
		}
	}
	if schemaVersion < 0 || dumpFormat < 0 {
		return 0, 0, fmt.Errorf("store: dump header must declare schema-version and dump-format")
	}
	return schemaVersion, dumpFormat, nil
}

func isEmptyDatabase(dbPath string) (bool, error) {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return true, nil
	} else if err != nil {
		return false, err
	}
	cmd := exec.Command("sqlite3", dbPath, "select count(*) from sqlite_master where type in ('table','trigger','index') and name not like 'sqlite_%';")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("store: inspect target database: %v: %s", err, strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(string(out)) == "0", nil
}
