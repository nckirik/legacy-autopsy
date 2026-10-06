// Package synthesis runs the Synthesis pipeline: model-authored catalog records
// with runtime-derived semantic identities (domain path + record title) and
// runtime-injected block IDs.
package synthesis

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// Directives are runtime-owned synthesis values.
type Directives struct {
	Namespace      string
	Snapshot       string
	StorageBinding string
	MaxClaims      int
}

// Catalog is one synthesis record family.
type Catalog struct {
	FieldID string
	IDType  string
	Section string
}

// Catalogs maps catalog IDs to their declared record family and section.
var Catalogs = map[string]Catalog{
	"MOD":   {FieldID: "ARCHITECTURE-MODULE", IDType: "MOD", Section: "architecture-blueprint"},
	"ER":    {FieldID: "ENTITY-RECORD", IDType: "ER", Section: "entity-catalog"},
	"REL":   {FieldID: "RELATIONSHIP-RECORD", IDType: "REL", Section: "relationship-catalog"},
	"SM":    {FieldID: "STATE-MACHINE", IDType: "SM", Section: "state-machine-catalog"},
	"DR":    {FieldID: "DATABASE-ROUTINE", IDType: "DR", Section: "database-routine-catalog"},
	"BR":    {FieldID: "BUSINESS-RULE-RECORD", IDType: "BR", Section: "business-rule-catalog"},
	"UC":    {FieldID: "USE-CASE", IDType: "UC", Section: "use-case-catalog"},
	"IF":    {FieldID: "INTERFACE-RECORD", IDType: "IF", Section: "interface-catalog"},
	"DEP":   {FieldID: "DEPLOYMENT-ELEMENT", IDType: "DEP", Section: "deployment-catalog"},
	"CFG":   {FieldID: "CONFIGURATION-RECORD", IDType: "CFG", Section: "deployment-catalog"},
	"SCHED": {FieldID: "SCHEDULE-RECORD", IDType: "SCHED", Section: "deployment-catalog"},
	"NFR":   {FieldID: "NFR-RECORD", IDType: "NFR", Section: "nfr-security-catalog"},
	"SEC":   {FieldID: "SECURITY-FINDING", IDType: "SEC", Section: "nfr-security-catalog"},
	"FLT":   {FieldID: "FAULT-RECORD", IDType: "FLT", Section: "nfr-security-catalog"},
}

// ClaimRef is a bounded committed claim.
type ClaimRef struct {
	ID      string
	Subject string
	Text    string
}

// LoadClaims reads every committed claim in bounded form, sorted by ID.
func LoadClaims(db string, doc *cdl.EIRDoc) ([]ClaimRef, error) {
	names, _, _, err := work.FieldContract(doc, "CLAIM")
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, "CLAIM")
	if err != nil {
		return nil, err
	}
	claims := make([]ClaimRef, 0, len(ids))
	for _, id := range ids {
		payload, err := store.Get(db, doc, "CLAIM", id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		claims = append(claims, ClaimRef{ID: id, Subject: values["subject-id-record-version"], Text: values["claim"]})
	}
	sort.Slice(claims, func(i, j int) bool { return claims[i].ID < claims[j].ID })
	return claims, nil
}

// Plan builds the bounded synthesis unit for one catalog.
func Plan(catalog Catalog, claims []ClaimRef, directives Directives) work.Unit {
	return work.Unit{
		ID:                   "U-SYNTHESIS-" + catalog.IDType,
		Mode:                 "Partial-Synthesis",
		Scope:                directives.Namespace,
		IDType:               catalog.IDType,
		Namespace:            directives.Namespace,
		IDOwnerField:         "domain-path",
		IDDiscriminatorField: "record-title",
		RuntimeFields:        []string{"block-id"},
		RequiresReview:       true,
		ResultType:           catalog.FieldID,
		Inputs:               []work.Input{{Name: "synthesis-context.md", Bytes: contextMarkdown(catalog, claims, directives)}},
		Sections:             []string{"synthesis-scope", catalog.Section},
		ReasoningEffort:      "off",
	}
}

func contextMarkdown(catalog Catalog, claims []ClaimRef, directives Directives) []byte {
	var b strings.Builder
	b.WriteString("# Synthesis context (runtime-generated, bounded)\n\n")
	b.WriteString("Group the evidence below into " + catalog.FieldID + " records. Author one record per coherent responsibility boundary that the claims support.\n\n")
	fmt.Fprintf(&b, "## Runtime directives (must be used exactly)\n\n")
	fmt.Fprintf(&b, "- **System Namespace:** %s\n", directives.Namespace)
	fmt.Fprintf(&b, "- **Environment/Snapshot:** %s\n", directives.Snapshot)
	b.WriteString("- `block-id` is runtime-assigned: emit the field with placeholder value `PENDING`.\n")
	b.WriteString("- `record-title` is a short human title; `domain-path` is the stable normalized path or domain coordinate the record describes.\n")
	b.WriteString("- Where the evidence does not support a field, use `Unknown` (fail-inclusive) or `None`; never invent facts.\n\n")
	b.WriteString("## Committed claims (evidence)\n\n")
	limit := directives.MaxClaims
	if limit <= 0 || limit > len(claims) {
		limit = len(claims)
	}
	for _, claim := range claims[:limit] {
		fmt.Fprintf(&b, "- %s [%s]: %s\n", claim.ID, strings.TrimSpace(claim.Subject), claim.Text)
	}
	return []byte(b.String())
}
