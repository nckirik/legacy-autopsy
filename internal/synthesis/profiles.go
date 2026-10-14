package synthesis

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// PersonaRow is one committed persona-registry row in bounded form.
type PersonaRow struct {
	Prefix     string
	Persona    string
	Slug       string
	Directory  string
	ActorClass string
	Entries    string
	Scope      string
	State      string
}

// LoadPersonas reads the 0A persona-prefix registry rows sorted by prefix.
func LoadPersonas(db string, doc *cdl.EIRDoc) ([]PersonaRow, error) {
	names, _, _, err := work.FieldContract(doc, "PERSONA-REGISTRY")
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, "PERSONA-REGISTRY")
	if err != nil {
		return nil, err
	}
	sort.Strings(ids)
	rows := make([]PersonaRow, 0, len(ids))
	for _, id := range ids {
		payload, err := store.Get(db, doc, "PERSONA-REGISTRY", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		rows = append(rows, PersonaRow{
			Prefix:     values["prefix"],
			Persona:    values["persona"],
			Slug:       values["persona-slug"],
			Directory:  values["persona-directory"],
			ActorClass: values["actor-class"],
			Entries:    values["canonical-entry-ids"],
			Scope:      values["scope-paths"],
			State:      values["registry-state"],
		})
	}
	return rows, nil
}

// FieldForID maps a typed record ID to its declared synthesis FIELD ID.
func FieldForID(id string) (string, bool) {
	prefix := id
	if idx := strings.Index(id, "-"); idx > 0 {
		prefix = id[:idx]
	}
	for _, catalog := range Catalogs {
		if catalog.IDType == prefix {
			return catalog.FieldID, true
		}
	}
	switch prefix {
	case "PRF":
		return "PERSONA-PROFILE", true
	case "PRV":
		return "PRIVACY-FINDING", true
	case "CAP":
		return "CAPABILITY-RECORD", true
	}
	return "", false
}

// ProfileID computes the §4.1 canonical PRF identity: type, namespace, persona
// prefix, and the fixed `persona-profile` discriminator.
func ProfileID(namespace, prefix string) (string, error) {
	return capabilities.TypedID("PRF", namespace, prefix, "persona-profile")
}

// ProfilePlan builds one Profile Synchronization unit per persona registry row.
// Each unit authors the complete persona profile payload with a runtime-injected
// PRF ID and semantic content fingerprint.
func ProfilePlan(db string, doc *cdl.EIRDoc, namespace, snapshot string) ([]work.Unit, error) {
	personas, err := LoadPersonas(db, doc)
	if err != nil {
		return nil, err
	}
	if len(personas) == 0 {
		return nil, nil
	}
	digest, reads, err := profileDigest(db, doc)
	if err != nil {
		return nil, err
	}
	registryReads := make([]work.Ref, 0, len(personas))
	for _, row := range personas {
		registryReads = append(registryReads, work.Ref{FieldID: "PERSONA-REGISTRY", RecordID: row.Prefix})
	}
	var units []work.Unit
	for _, row := range personas {
		id, err := ProfileID(namespace, row.Prefix)
		if err != nil {
			return nil, err
		}
		var b strings.Builder
		b.WriteString("# Persona profile context (runtime-generated, bounded)\n\n")
		fmt.Fprintf(&b, "Author the complete PERSONA-PROFILE payload for persona `%s` (%s).\n\n", row.Prefix, row.Persona)
		b.WriteString("## Runtime directives (must be used exactly)\n\n")
		fmt.Fprintf(&b, "- **System Namespace:** %s\n", namespace)
		fmt.Fprintf(&b, "- **Environment/Snapshot:** %s\n", snapshot)
		fmt.Fprintf(&b, "- `prf-id` is runtime-assigned: emit the field with placeholder value `PENDING`.\n")
		fmt.Fprintf(&b, "- `prefix-actor-class` must begin with the exact prefix `%s`.\n", row.Prefix)
		b.WriteString("- `semantic-record-version` is `1`.\n")
		b.WriteString("- `semantic-content-fingerprint` is runtime-computed: emit the field with placeholder value `PENDING`.\n")
		b.WriteString("- `derived-from` lists the exact record IDs used in this profile, comma-separated, from the digest below. Do not write versions or fingerprints: the runtime annotates each ID with its committed version and fingerprint.\n")
		b.WriteString("- `material-claims` and `observed-legacy-facts` list CLM IDs from the digest; use only IDs shown, never invent.\n")
		b.WriteString("- `target-decisions` is `None` while no approved decision records exist.\n")
		b.WriteString("- Where the evidence does not support a field, use `Unknown` (fail-inclusive) or `None`.\n")
		b.WriteString("- Table-shaped values (for example the authorization matrix) serialize into the single field line as semicolon-separated rows.\n\n")
		b.WriteString("## Persona registry row\n\n")
		fmt.Fprintf(&b, "- **Persona:** %s\n- **Prefix:** %s\n- **Persona-slug:** %s\n- **Persona-directory:** %s\n- **Actor-class:** %s\n- **Canonical-entry-ids:** %s\n- **Scope-paths:** %s\n- **Registry-state:** %s\n\n",
			row.Persona, row.Prefix, row.Slug, row.Directory, row.ActorClass, row.Entries, row.Scope, row.State)
		b.WriteString(digest)
		units = append(units, work.Unit{
			ID:                       "U-PROFILE-" + row.Prefix,
			Mode:                     "Profile Synchronization",
			Scope:                    namespace,
			RecordID:                 id,
			RuntimeFields:            []string{"prf-id"},
			RuntimeFingerprintFields: []string{"semantic-content-fingerprint"},
			RuntimeAnnotatedFields:   []string{"derived-from"},
			ExpectedRecords:          1,
			RequiresReview:           true,
			ResultType:               "PERSONA-PROFILE",
			Reads:                    append(append([]work.Ref{}, registryReads...), reads...),
			Inputs:                   []work.Input{{Name: "profile-context.md", Bytes: []byte(b.String())}},
			Sections:                 []string{"persona-profile", "dependency-provenance"},
			ReasoningEffort:          "off",
		})
	}
	return units, nil
}

// profileDigest renders a bounded digest of committed semantic records and the
// read set it was built from.
func profileDigest(db string, doc *cdl.EIRDoc) (string, []work.Ref, error) {
	var b strings.Builder
	var reads []work.Ref
	b.WriteString("## Committed record digest\n\n")
	type digest struct {
		fieldID string
		label   string
		fields  []string
	}
	digests := []digest{
		{fieldID: "USE-CASE", label: "Use cases", fields: []string{"record-title", "goal", "primary-persona-secondary-actors", "capability-classification"}},
		{fieldID: "INTERFACE-RECORD", label: "Interfaces", fields: []string{"record-title", "channel-transport-direction"}},
		{fieldID: "ENTITY-RECORD", label: "Entities", fields: []string{"record-title", "domain-context-purpose"}},
		{fieldID: "ARCHITECTURE-MODULE", label: "Modules", fields: []string{"record-title", "responsibility-domain-boundary"}},
		{fieldID: "FAULT-RECORD", label: "Faults", fields: []string{"record-title", "failure-class-trigger-detection-point"}},
		{fieldID: "SECURITY-FINDING", label: "Security findings", fields: []string{"record-title", "category", "risk-level"}},
		{fieldID: "PRIVACY-FINDING", label: "Privacy findings", fields: []string{"data-threat-review-category", "synthesis-sec-mirror"}},
		{fieldID: "BUSINESS-RULE-RECORD", label: "Business rules", fields: []string{"record-title", "condition"}},
		{fieldID: "CLAIM", label: "Claims", fields: []string{"claim", "subject-id-record-version"}},
	}
	for _, d := range digests {
		names, _, _, err := work.FieldContract(doc, d.fieldID)
		if err != nil {
			return "", nil, err
		}
		ids, err := store.List(db, doc, d.fieldID)
		if err != nil {
			return "", nil, err
		}
		if len(ids) == 0 {
			continue
		}
		fmt.Fprintf(&b, "### %s\n\n", d.label)
		for _, id := range ids {
			payload, err := store.Get(db, doc, d.fieldID, id)
			if err != nil {
				return "", nil, err
			}
			values, err := work.PayloadText(names, payload)
			if err != nil {
				return "", nil, err
			}
			parts := make([]string, 0, len(d.fields))
			for _, field := range d.fields {
				value := strings.Join(strings.Fields(values[field]), " ")
				if d.fieldID == "CLAIM" && len(value) > 240 {
					value = value[:240] + "…"
				}
				if value == "" {
					continue
				}
				parts = append(parts, field+": "+value)
			}
			fmt.Fprintf(&b, "- %s: %s\n", id, strings.Join(parts, " | "))
			reads = append(reads, work.Ref{FieldID: d.fieldID, RecordID: id})
		}
		b.WriteString("\n")
	}
	sort.Slice(reads, func(i, j int) bool {
		if reads[i].FieldID != reads[j].FieldID {
			return reads[i].FieldID < reads[j].FieldID
		}
		return reads[i].RecordID < reads[j].RecordID
	})
	return b.String(), reads, nil
}

// ProfileStats reports PRF identity and closure coverage.
type ProfileStats struct {
	RegistryRows int
	Profiles     int
	Missing      []string
	Mismatched   []string
	BadVersion   []string
	BadFinger    []string
	Dangling     []string
}

var profileTokenPattern = regexp.MustCompile(`\b(CLM|BR|UC|IF|SCHED|DEP|CMP|ER|MOD|SM|FLT|SEC|NFR|REL|CFG|PRF|PRV|CAP|DR|SRC)-[A-Z0-9-]*-?[A-F0-9]{12,64}\b`)

// ProfileCoverage evaluates the candidate-stage portions of §15.1 condition 12:
// every persona registry row has its legal PRF identity, the payload is
// fingerprinted correctly, and the declared dependency closure resolves.
func ProfileCoverage(db string, doc *cdl.EIRDoc, namespace string) (ProfileStats, error) {
	var stats ProfileStats
	personas, err := LoadPersonas(db, doc)
	if err != nil {
		return stats, err
	}
	stats.RegistryRows = len(personas)
	names, _, _, err := work.FieldContract(doc, "PERSONA-PROFILE")
	if err != nil {
		return stats, err
	}
	profileIDs, err := store.List(db, doc, "PERSONA-PROFILE")
	if err != nil {
		return stats, err
	}
	stats.Profiles = len(profileIDs)
	matched := map[string]bool{}
	for _, row := range personas {
		expected, err := ProfileID(namespace, row.Prefix)
		if err != nil {
			return stats, err
		}
		payload, err := store.Get(db, doc, "PERSONA-PROFILE", expected)
		if err != nil {
			stats.Missing = append(stats.Missing, expected)
			continue
		}
		matched[expected] = true
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return stats, err
		}
		if strings.TrimSpace(values["prf-id"]) != expected {
			stats.Mismatched = append(stats.Mismatched, expected+" (prf-id)")
		}
		if firstPrefixToken(values["prefix-actor-class"]) != row.Prefix {
			stats.Mismatched = append(stats.Mismatched, expected+" (prefix-actor-class)")
		}
		version, err := strconv.Atoi(strings.TrimSpace(values["semantic-record-version"]))
		if err != nil || version < 1 {
			stats.BadVersion = append(stats.BadVersion, expected)
		}
		recomputed, err := capabilities.SemanticContentFingerprint(payload, "semantic-content-fingerprint")
		if err != nil {
			return stats, err
		}
		if strings.TrimSpace(values["semantic-content-fingerprint"]) != recomputed {
			stats.BadFinger = append(stats.BadFinger, expected)
		}
		for _, field := range []string{"derived-from", "material-claims", "observed-legacy-facts"} {
			for _, token := range profileTokenPattern.FindAllString(values[field], -1) {
				resolved, err := ProfileTokenResolves(db, doc, token)
				if err != nil {
					return stats, err
				}
				if !resolved {
					stats.Dangling = append(stats.Dangling, expected+" -> "+token)
				}
			}
		}
	}
	for _, id := range profileIDs {
		if !matched[id] {
			stats.Mismatched = append(stats.Mismatched, id+" (no registry row)")
		}
	}
	sort.Strings(stats.Missing)
	sort.Strings(stats.Mismatched)
	sort.Strings(stats.BadVersion)
	sort.Strings(stats.BadFinger)
	sort.Strings(stats.Dangling)
	return stats, nil
}

// firstPrefixToken extracts the leading prefix token from a
// `Prefix / Actor Class` value, tolerating the declared slash/space separators.
func firstPrefixToken(raw string) string {
	for _, token := range strings.FieldsFunc(raw, func(r rune) bool {
		return r == '/' || r == '|' || r == ';' || r == ',' || unicode.IsSpace(r)
	}) {
		if token != "" {
			return strings.TrimSpace(token)
		}
	}
	return ""
}

// ProfileTokenResolves resolves one typed token against its record table.
func ProfileTokenResolves(db string, doc *cdl.EIRDoc, token string) (bool, error) {
	if strings.HasPrefix(token, "SRC-") {
		return store.Exists(db, doc, "SOURCE-INVENTORY-RECORD", token)
	}
	prefix := token
	if idx := strings.Index(token, "-"); idx > 0 {
		prefix = token[:idx]
	}
	fieldID, ok := profileTokenField[prefix]
	if !ok {
		return true, nil
	}
	return store.Exists(db, doc, fieldID, token)
}

var profileTokenField = map[string]string{
	"CLM":   "CLAIM",
	"BR":    "BUSINESS-RULE-RECORD",
	"UC":    "USE-CASE",
	"IF":    "INTERFACE-RECORD",
	"SCHED": "SCHEDULE-RECORD",
	"DEP":   "DEPLOYMENT-ELEMENT",
	"CMP":   "ATOMIC-COMPONENT",
	"ER":    "ENTITY-RECORD",
	"MOD":   "ARCHITECTURE-MODULE",
	"SM":    "STATE-MACHINE",
	"FLT":   "FAULT-RECORD",
	"SEC":   "SECURITY-FINDING",
	"NFR":   "NFR-RECORD",
	"REL":   "RELATIONSHIP-RECORD",
	"CFG":   "CONFIGURATION-RECORD",
	"PRF":   "PERSONA-PROFILE",
	"PRV":   "PRIVACY-FINDING",
	"CAP":   "CAPABILITY-RECORD",
	"DR":    "DATABASE-ROUTINE",
}
