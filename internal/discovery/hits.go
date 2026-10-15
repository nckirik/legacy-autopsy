package discovery

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/work"
)

// HitKind is one non-file source kind the deterministic hit enumerator
// recognizes. The values match the declared SRC kinds (protocol/typed-id.cdl):
// SRC-ROUTE, SRC-QUERY, SRC-JOB.
type HitKind string

const (
	HitRoute  HitKind = "ROUTE"
	HitQuery  HitKind = "QUERY"
	HitJob    HitKind = "JOB"
	HitTable  HitKind = "TABLE"
	HitView   HitKind = "VIEW"
	HitDR     HitKind = "DR"
	HitAction HitKind = "ACTION"
)

// Hit is one deterministically enumerated non-file source unit: a route
// registration, a SQL statement, or a scheduled job found inside an enumerated
// file. The evidence excerpt and its fingerprint are the bounded observed bytes;
// the model authors the inventory record, the runtime computes the SRC identity.
type Hit struct {
	Coordinate  string  `json:"coordinate"`
	Kind        HitKind `json:"kind"`
	File        string  `json:"file"`
	Line        int     `json:"line"`
	Evidence    string  `json:"evidence"`
	Fingerprint string  `json:"fingerprint"`
}

// Deterministic enumeration bounds. Values are declared so runs are reproducible
// and large or minified inputs cannot stall the scan.
const (
	maxHitScanBytes   = 2 << 20 // skip files larger than 2 MiB
	maxHitLineBytes   = 4096    // skip lines longer than this
	maxHitsPerFile    = 200     // bounded hits per file
	maxHitEvidenceLen = 240     // bounded evidence excerpt length
)

// scannableExtensions is the closed set of textual file extensions the hit
// enumerator reads. Everything else is skipped without inspection.
var scannableExtensions = map[string]bool{
	".py": true, ".js": true, ".mjs": true, ".cjs": true, ".ts": true,
	".tsx": true, ".jsx": true, ".php": true, ".rb": true, ".go": true,
	".java": true, ".cs": true, ".sh": true, ".ps1": true, ".sql": true,
	".yaml": true, ".yml": true, ".json": true, ".toml": true, ".ini": true,
	".conf": true,
}

// routeReceivers and routeMethods declare the Express/Koa-style call forms.
var routeReceivers = map[string]bool{
	"app": true, "router": true, "api": true, "server": true, "bp": true,
	"route": true, "r": true, "mux": true, "http": true, "web": true,
}

var routeMethods = map[string]bool{
	"route": true, "get": true, "post": true, "put": true, "patch": true,
	"delete": true, "head": true, "options": true, "all": true,
	"handlefunc": true, "handle": true,
}

// jobMarkers declare the scheduling surfaces the enumerator recognizes.
var jobMarkers = []string{
	"crontab", "cron.schedule(", "schedule.every(", "@scheduled",
	"backgroundscheduler", "apscheduler", "celerybeat", "@celery.task",
	"setinterval(", "cron_expression",
}

// queryKeywords declare the SQL statement openers the enumerator recognizes as
// query units. Schema-object openers (CREATE TABLE/VIEW/TRIGGER/PROCEDURE/
// FUNCTION) classify as their own kinds instead.
var queryKeywords = []string{
	"ALTER TABLE", "DROP TABLE", "SELECT", "INSERT INTO",
	"UPDATE ", "DELETE FROM",
}

// schemaOpeners declare the DDL object declarations the enumerator recognizes
// and the SRC kind each object carries.
var schemaOpeners = []struct {
	keyword string
	kind    HitKind
}{
	{"CREATE MATERIALIZED VIEW", HitView},
	{"CREATE TABLE", HitTable},
	{"CREATE VIEW", HitView},
	{"CREATE TRIGGER", HitDR},
	{"CREATE PROCEDURE", HitDR},
	{"CREATE FUNCTION", HitDR},
}

// outputMarkers declare user-visible output/export constructs (file downloads
// and attachments) that produce artifacts rather than data reads.
var outputMarkers = []string{
	"st.download_button(", "download_button(", "send_file(", "send_from_directory(",
	"fileresponse(", "res.download(", "res.sendfile(", "content-disposition",
}

// EnumerateHits scans the enumerated candidate files for route, query, and job
// source units. The scan is lexical and line-based with declared bounds; it
// never interprets semantics. Returns at most maxHits hits in deterministic
// file-then-line order.
func EnumerateHits(root string, files []Candidate, maxHits int) ([]Hit, error) {
	if maxHits <= 0 {
		maxHits = 500
	}
	ordered := make([]Candidate, len(files))
	copy(ordered, files)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Path < ordered[j].Path })
	var hits []Hit
	for _, file := range ordered {
		if len(hits) >= maxHits {
			break
		}
		extension := strings.ToLower(filepath.Ext(file.Path))
		if !scannableExtensions[extension] {
			continue
		}
		fileHits, err := enumerateFileHits(root, file.Path, maxHits-len(hits))
		if err != nil {
			return nil, err
		}
		hits = append(hits, fileHits...)
	}
	return hits, nil
}

func enumerateFileHits(root, relative string, remaining int) ([]Hit, error) {
	full := filepath.Join(root, filepath.FromSlash(relative))
	info, err := os.Stat(full)
	if err != nil {
		return nil, fmt.Errorf("discovery: hit scan %s: %w", relative, err)
	}
	if info.Size() > maxHitScanBytes {
		return nil, nil
	}
	file, err := os.Open(full)
	if err != nil {
		return nil, fmt.Errorf("discovery: hit scan %s: %w", relative, err)
	}
	defer file.Close()
	var hits []Hit
	occurrence := map[string]int{}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), maxHitLineBytes+1)
	line := 0
	for scanner.Scan() {
		line++
		text := scanner.Text()
		if len(text) > maxHitLineBytes {
			continue
		}
		trimmed := strings.TrimSpace(strings.TrimPrefix(text, "\ufeff"))
		if trimmed == "" || isCommentLine(trimmed) {
			continue
		}
		hit, ok := classifyHit(relative, line, trimmed)
		if !ok {
			continue
		}
		occurrence[hit.Coordinate]++
		if occurrence[hit.Coordinate] > 1 {
			hit.Coordinate = fmt.Sprintf("%s #%d", hit.Coordinate, occurrence[hit.Coordinate])
		}
		hits = append(hits, hit)
		if len(hits) >= maxHitsPerFile || len(hits) >= remaining {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("discovery: hit scan %s: %w", relative, err)
	}
	return hits, nil
}

// isCommentLine skips declared comment forms so commented-out declarations are
// not enumerated as live units.
func isCommentLine(trimmed string) bool {
	for _, prefix := range []string{"#", "//", "--", ";", "*", "/*"} {
		if strings.HasPrefix(trimmed, prefix) {
			return true
		}
	}
	return false
}

func classifyHit(file string, line int, trimmed string) (Hit, bool) {
	if hit, ok := routeHit(file, line, trimmed); ok {
		return hit, true
	}
	if hit, ok := schemaHit(file, line, trimmed); ok {
		return hit, true
	}
	if hit, ok := queryHit(file, line, trimmed); ok {
		return hit, true
	}
	if hit, ok := jobHit(file, line, trimmed); ok {
		return hit, true
	}
	if hit, ok := outputHit(file, line, trimmed); ok {
		return hit, true
	}
	return Hit{}, false
}

// schemaHit recognizes DDL object declarations and classifies them as their
// declared SRC kinds: tables, views, and database routines.
func schemaHit(file string, line int, trimmed string) (Hit, bool) {
	body := strings.TrimLeft(trimmed, "\"'`([ ")
	upper := strings.ToUpper(body)
	for _, opener := range schemaOpeners {
		if !strings.HasPrefix(upper, opener.keyword+" ") {
			continue
		}
		name := ddlObjectName(strings.TrimSpace(body[len(opener.keyword):]))
		if name == "" {
			return Hit{}, false
		}
		coordinate := normalizeHitText(file + " :: " + strings.ToLower(string(opener.kind)) + " :: " + name)
		return newHit(opener.kind, file, line, coordinate, trimmed), true
	}
	return Hit{}, false
}

// ddlObjectName extracts the declared object name, tolerating IF NOT EXISTS and
// identifier quoting.
func ddlObjectName(rest string) string {
	if strings.HasPrefix(strings.ToUpper(rest), "IF NOT EXISTS ") {
		rest = strings.TrimSpace(rest[len("IF NOT EXISTS "):])
	}
	rest = strings.TrimLeft(rest, "`\"[ ")
	end := strings.IndexAny(rest, " `\"([\t")
	if end < 0 {
		end = len(rest)
	}
	name := strings.TrimRight(rest[:end], ";,")
	if name == "" {
		return ""
	}
	return name
}

// outputHit recognizes user-visible file output/export constructs.
func outputHit(file string, line int, trimmed string) (Hit, bool) {
	lower := strings.ToLower(trimmed)
	for _, marker := range outputMarkers {
		if strings.Contains(lower, marker) {
			if len(trimmed) > 120 {
				trimmed = trimmed[:120]
			}
			coordinate := normalizeHitText(file + " :: action :: " + trimmed)
			return newHit(HitAction, file, line, coordinate, trimmed), true
		}
	}
	return Hit{}, false
}

// routeHit recognizes decorator and call-style route registrations with a
// declared receiver or decorator method. The discriminator is the first string
// literal on the line (the route path) when present.
func routeHit(file string, line int, trimmed string) (Hit, bool) {
	method, ok := routeRegistration(trimmed)
	if !ok {
		return Hit{}, false
	}
	path, hasPath := firstStringLiteral(trimmed)
	discriminator := strings.ToUpper(method)
	if hasPath {
		discriminator = strings.ToUpper(method) + " " + path
		if strings.EqualFold(method, "route") {
			discriminator = path
		}
	}
	coordinate := normalizeHitText(file + " :: route :: " + discriminator)
	return newHit(HitRoute, file, line, coordinate, trimmed), true
}

// routeRegistration returns the registration form name for a route line.
func routeRegistration(trimmed string) (string, bool) {
	if strings.HasPrefix(trimmed, "@") {
		body := trimmed[1:]
		if idx := strings.Index(body, "("); idx >= 0 {
			body = body[:idx]
		}
		parts := strings.Split(body, ".")
		if len(parts) >= 2 {
			method := strings.ToLower(strings.TrimSpace(parts[len(parts)-1]))
			if routeMethods[method] {
				return method, true
			}
		}
		return "", false
	}
	open := strings.Index(trimmed, "(")
	if open <= 0 {
		return "", false
	}
	call := trimmed[:open]
	dot := strings.LastIndex(call, ".")
	if dot <= 0 {
		return "", false
	}
	receiver := strings.ToLower(strings.TrimSpace(call[:dot]))
	method := strings.ToLower(strings.TrimSpace(call[dot+1:]))
	if !routeReceivers[receiver] || !routeMethods[method] {
		return "", false
	}
	if _, ok := firstStringLiteral(trimmed[open:]); !ok {
		return "", false
	}
	return method, true
}

// firstStringLiteral returns the first single- or double-quoted literal value.
func firstStringLiteral(text string) (string, bool) {
	for i := 0; i < len(text); i++ {
		quote := text[i]
		if quote != '"' && quote != '\'' {
			continue
		}
		end := strings.IndexByte(text[i+1:], quote)
		if end < 0 {
			return "", false
		}
		value := text[i+1 : i+1+end]
		if value == "" {
			return "", false
		}
		return value, true
	}
	return "", false
}

// queryHit recognizes SQL statements in scanned files; the discriminator is the
// normalized statement text.
func queryHit(file string, line int, trimmed string) (Hit, bool) {
	statement, ok := queryStatement(trimmed)
	if !ok {
		return Hit{}, false
	}
	coordinate := normalizeHitText(file + " :: query :: " + statement)
	return newHit(HitQuery, file, line, coordinate, trimmed), true
}

func queryStatement(trimmed string) (string, bool) {
	body := strings.TrimLeft(trimmed, "\"'`([ ")
	upper := strings.ToUpper(body)
	for _, keyword := range queryKeywords {
		if strings.HasPrefix(upper, keyword) {
			statement := strings.TrimRight(body, ";")
			statement = strings.Join(strings.Fields(statement), " ")
			return statement, true
		}
	}
	return "", false
}

// jobHit recognizes declared scheduling markers on the line.
func jobHit(file string, line int, trimmed string) (Hit, bool) {
	lower := strings.ToLower(trimmed)
	for _, marker := range jobMarkers {
		if strings.Contains(lower, marker) {
			coordinate := normalizeHitText(file + " :: job :: " + normalizeHitText(trimmed))
			return newHit(HitJob, file, line, coordinate, trimmed), true
		}
	}
	return Hit{}, false
}

func newHit(kind HitKind, file string, line int, coordinate, evidence string) Hit {
	excerpt := normalizeHitText(evidence)
	if len(excerpt) > maxHitEvidenceLen {
		excerpt = excerpt[:maxHitEvidenceLen]
	}
	sum := sha256.Sum256([]byte(evidence))
	return Hit{
		Coordinate:  coordinate,
		Kind:        kind,
		File:        file,
		Line:        line,
		Evidence:    excerpt,
		Fingerprint: fmt.Sprintf("sha256:%x", sum[:]),
	}
}

// normalizeHitText collapses whitespace and trims the result.
func normalizeHitText(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

// hitPOV maps a hit kind to its isolated point of view (§2 five isolated
// POVs): queries, tables, views, and database routines are data movement
// (POV-4), routes are system wiring (POV-5), jobs are background workflows
// (POV-3), and output/export actions are frontend workflows (POV-1).
func hitPOV(kind HitKind) string {
	switch kind {
	case HitQuery, HitTable, HitView, HitDR:
		return "POV-4"
	case HitRoute:
		return "POV-5"
	case HitJob:
		return "POV-3"
	case HitAction:
		return "POV-1"
	default:
		return ""
	}
}

// ClusterForHit resolves one hit file to exactly one cluster: an exact entry
// coordinate first, then a file under an entry coordinate's directory. When
// several clusters could match, the hit stays unbound (fail-inclusive); when
// exactly one cluster exists it is the fallback owner of the bounded scan.
func ClusterForHit(hit Hit, clusters []Cluster) (Cluster, bool) {
	ordered := make([]Cluster, len(clusters))
	copy(ordered, clusters)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].ID < ordered[j].ID })
	for _, cluster := range ordered {
		for _, entry := range cluster.Entries {
			if entry == hit.File {
				return cluster, true
			}
		}
	}
	matches := map[string]Cluster{}
	dir := path.Dir(hit.File)
	for _, cluster := range ordered {
		for _, entry := range cluster.Entries {
			entryDir := path.Dir(entry)
			if entryDir == "." {
				continue
			}
			if dir == entryDir || strings.HasPrefix(hit.File, entryDir+"/") {
				matches[cluster.ID] = cluster
				break
			}
		}
	}
	if len(matches) == 1 {
		for _, cluster := range matches {
			return cluster, true
		}
	}
	if len(ordered) == 1 {
		return ordered[0], true
	}
	return Cluster{}, false
}

// HitPlan builds one bounded inventory unit per bound (cluster, POV) hit batch.
// Each unit must author exactly one SOURCE-INVENTORY-RECORD per hit, in order;
// the runtime computes every SRC ID from the payload kind and coordinate. Units
// carry the strict single-scope bindings (§8.2): one persona where the cluster
// owner is a registered prefix, one cluster, one traversal track, one POV.
func HitPlan(hits []Hit, clusters []Cluster, directives Directives) []work.Unit {
	batchSize := directives.BatchSize
	if batchSize <= 0 {
		batchSize = 25
	}
	type group struct {
		cluster Cluster
		bound   bool
		pov     string
		hits    []Hit
	}
	var groups []*group
	index := map[string]*group{}
	for _, hit := range hits {
		pov := hitPOV(hit.Kind)
		if pov == "" {
			continue
		}
		cluster, bound := ClusterForHit(hit, clusters)
		key := pov + "|unbound"
		if bound {
			key = pov + "|" + cluster.ID
		}
		item, ok := index[key]
		if !ok {
			item = &group{cluster: cluster, bound: bound, pov: pov}
			index[key] = item
			groups = append(groups, item)
		}
		item.hits = append(item.hits, hit)
	}
	var units []work.Unit
	for _, item := range groups {
		for start := 0; start < len(item.hits); start += batchSize {
			end := start + batchSize
			if end > len(item.hits) {
				end = len(item.hits)
			}
			batch := item.hits[start:end]
			persona := ""
			if item.bound {
				if prefix, ok := item.cluster.RegisteredPersona(directives.PersonaPrefixes); ok {
					persona = prefix
				}
			}
			clusterToken := "UNBOUND"
			if item.bound {
				clusterToken = item.cluster.ID
			}
			unit := work.Unit{
				ID:               fmt.Sprintf("U-DISCOVERY-HITS-%s-%s-%03d", clusterToken, item.pov, start/batchSize+1),
				Mode:             "Discovery",
				Scope:            directives.Namespace,
				Persona:          persona,
				POV:              item.pov,
				IDFromCoordinate: "concrete-coordinate",
				IDType:           "SRC",
				Namespace:        directives.Namespace,
				ExpectedRecords:  len(batch),
				ResultType:       "SOURCE-INVENTORY-RECORD",
				Reads:            []work.Ref{{FieldID: "PREFLIGHT-HEADER", RecordID: "0A"}},
				Inputs:           []work.Input{{Name: "discovery-hits-context.md", Bytes: hitContextMarkdown(batch, item.cluster, item.bound, persona, item.pov, directives)}},
				Sections:         []string{"atomic-unit-rule", "source-inventory", "five-isolated-povs", "traversal-tracks"},
				ReasoningEffort:  "off",
			}
			units = append(units, unit)
		}
	}
	return units
}

func hitContextMarkdown(batch []Hit, cluster Cluster, bound bool, persona, pov string, directives Directives) []byte {
	var b strings.Builder
	b.WriteString("# Discovery hits batch (runtime-generated, bounded)\n\n")
	b.WriteString("## Candidate hits - author exactly one SOURCE-INVENTORY-RECORD per row, in this order\n\n")
	b.WriteString("| Concrete Coordinate | Kind | File | Line | Evidence Fingerprint |\n| :-- | :-- | :-- | --: | :-- |\n")
	for _, hit := range batch {
		fmt.Fprintf(&b, "| %s | %s | %s | %d | %s |\n", hit.Coordinate, hit.Kind, hit.File, hit.Line, hit.Fingerprint)
	}
	b.WriteString("\n## Evidence excerpts (exact observed bytes; do not extend)\n\n")
	for _, hit := range batch {
		fmt.Fprintf(&b, "### %s\n\n> %s\n\n", hit.Coordinate, hit.Evidence)
	}
	b.WriteString("## Strict single-scope bindings (§8.2; do not switch roots, personas, POVs, or tracks)\n\n")
	fmt.Fprintf(&b, "- **Persona:** %s\n", orNone(persona))
	if bound {
		fmt.Fprintf(&b, "- **Entry Cluster:** %s (%s)\n", cluster.ID, cluster.Title)
		fmt.Fprintf(&b, "- **Traversal Track:** %s\n", cluster.Track)
	} else {
		b.WriteString("- **Entry Cluster:** None (no declared cluster covers every hit file; the record states fail-inclusive None bindings)\n")
		b.WriteString("- **Traversal Track:** None\n")
	}
	fmt.Fprintf(&b, "- **Active POV:** %s\n\n", pov)
	if !bound {
		b.WriteString("The hit-to-cluster binding is unresolved. Emit `canonical-owner: None`, `scope-summary: Unbound`, and `traversal-tracks: Normal` only when the evidence shows it; otherwise `Unknown`.\n\n")
	} else if persona == "" {
		b.WriteString("The cluster owner is not a registered persona (shared ownership). Emit `canonical-owner: None` and leave persona attribution to the shared-entry record; do not invent a persona.\n\n")
	}
	b.WriteString("## Runtime directives (must be used exactly)\n\n")
	fmt.Fprintf(&b, "- **System Namespace:** %s\n", directives.Namespace)
	fmt.Fprintf(&b, "- **Environment/Snapshot:** %s\n", directives.Snapshot)
	fmt.Fprintf(&b, "- **Canonical Storage Binding:** %s\n", directives.StorageBinding)
	b.WriteString("- `kind` must be the row Kind exactly (`ROUTE`, `QUERY`, or `JOB`).\n")
	b.WriteString("- `concrete-coordinate` must be copied exactly from the candidate row.\n")
	b.WriteString("- `source-fingerprint-snapshot-src-version` is the row Evidence Fingerprint plus ` version 1`.\n")
	b.WriteString("- `discovery-method` is `other` for this bounded lexical scan; use a stronger declared method only with evidence.\n")
	b.WriteString("- `parent-container-src`, `mappings`, `frontier-references`, `coverage-rows`, and `claims` are `None` unless evidence-backed.\n")
	if bound {
		fmt.Fprintf(&b, "- `traversal-tracks` is `%s` (the bound cluster track).\n", cluster.Track)
	}
	if persona != "" {
		fmt.Fprintf(&b, "- `canonical-owner` is `%s` (the bound cluster's registered owner persona).\n", persona)
	}
	b.WriteString("- `disposition-summary` is `Unknown` (fail-inclusive) until swept.\n")
	b.WriteString("- `coverage-summary` is the closed enum value `C-GAP` for newly discovered unswept units.\n")
	return []byte(b.String())
}

func orNone(value string) string {
	if strings.TrimSpace(value) == "" {
		return "None"
	}
	return value
}
