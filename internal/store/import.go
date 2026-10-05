package store

import (
	"fmt"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	md "github.com/nckirik/legacy-autopsy/internal/markdown"
)

// PreflightRecordID is the singleton record ID for the `0A` preflight header.
const PreflightRecordID = "0A"

// ImportPreflight ingests the bootstrap `0A-PREFLIGHT.md` metadata and persona
// registry into the store. `0A` is the one Markdown file that must exist before
// the store because it declares the canonical storage binding; this importer
// covers the preflight header (singleton) and every persona registry row.
// Cluster, traversal-matrix, and other discovery records are ingested by later
// modes, not by bootstrap import.
func ImportPreflight(dbPath string, doc *cdl.EIRDoc, source []byte) error {
	parsed, err := md.Parse(string(source))
	if err != nil {
		return fmt.Errorf("store: parse preflight: %w", err)
	}
	meta := sectionNodes(parsed, "1. Metadata and Boundaries")
	if len(meta) == 0 {
		return fmt.Errorf("store: preflight metadata section not found")
	}
	payload, err := preflightHeaderPayload(meta)
	if err != nil {
		return err
	}
	if err := Ingest(dbPath, doc, "PREFLIGHT-HEADER", PreflightRecordID, 1, payload); err != nil {
		return fmt.Errorf("store: import preflight header: %w", err)
	}

	registry := sectionNodes(parsed, "2. Persona Registry")
	rows := tableRows(registry)
	if len(rows) < 2 {
		return fmt.Errorf("store: preflight persona registry table not found")
	}
	labels := splitTableRow(rows[0])
	prefixAt := -1
	for i, label := range labels {
		if capabilities.ColumnName(label) == "prefix" {
			prefixAt = i
		}
	}
	if prefixAt < 0 {
		return fmt.Errorf("store: preflight persona registry has no Prefix column")
	}
	imported := 0
	for _, row := range rows[1:] {
		cells := splitTableRow(row)
		if isTableSeparator(cells) {
			continue
		}
		if len(cells) != len(labels) {
			return fmt.Errorf("store: preflight persona row has %d cells, want %d", len(cells), len(labels))
		}
		var b strings.Builder
		for i, label := range labels {
			fmt.Fprintf(&b, "- **%s:** %s\n", label, cells[i])
		}
		recordID := strings.TrimSpace(cells[prefixAt])
		if recordID == "" {
			return fmt.Errorf("store: preflight persona row has an empty Prefix")
		}
		if err := Ingest(dbPath, doc, "PERSONA-REGISTRY", recordID, 1, []byte(b.String())); err != nil {
			return fmt.Errorf("store: import persona %s: %w", recordID, err)
		}
		imported++
	}
	if imported == 0 {
		return fmt.Errorf("store: preflight persona registry has no rows")
	}
	return nil
}

// preflightHeaderPayload builds the PREFLIGHT-HEADER record payload from the
// metadata bullets. When the binding bullet is absent, the protocol default
// `markdown/1` is declared explicitly rather than inferred.
func preflightHeaderPayload(nodes []md.Node) ([]byte, error) {
	var b strings.Builder
	hasBinding := false
	for _, node := range nodes {
		if node.Kind != md.Field {
			continue
		}
		if capabilities.ColumnName(node.Label) == "canonical-storage-binding" {
			hasBinding = true
		}
		fmt.Fprintf(&b, "- **%s:** %s\n", node.Label, node.Value)
	}
	if b.Len() == 0 {
		return nil, fmt.Errorf("store: preflight metadata has no fields")
	}
	if !hasBinding {
		fmt.Fprintf(&b, "- **Canonical Storage Binding:** %s\n", "markdown/1")
	}
	return capabilities.CanonicalMarkdown([]byte(b.String()))
}

// sectionNodes returns the nodes of one `##` section, excluding nested headings.
func sectionNodes(doc *md.Document, title string) []md.Node {
	start := -1
	level := 0
	for i, node := range doc.Nodes {
		if node.Kind == md.Heading && node.Text == title {
			start = i
			level = node.Level
			break
		}
	}
	if start < 0 {
		return nil
	}
	var out []md.Node
	for _, node := range doc.Nodes[start+1:] {
		if node.Kind == md.Heading && node.Level <= level {
			break
		}
		out = append(out, node)
	}
	return out
}

func tableRows(nodes []md.Node) []string {
	var rows []string
	for _, node := range nodes {
		if node.Kind == md.TableRow {
			rows = append(rows, node.Raw)
		}
	}
	return rows
}

func splitTableRow(raw string) []string {
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

func isTableSeparator(cells []string) bool {
	if len(cells) == 0 {
		return false
	}
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
