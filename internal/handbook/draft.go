package handbook

import (
	"fmt"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// ChapterSpec maps one handbook chapter to its synthesis catalog families.
type ChapterSpec struct {
	Number   string
	Title    string
	Catalogs []string
}

// DefaultChapters is the normative handbook layout with the implemented
// catalogs bound to their chapters.
func DefaultChapters() []ChapterSpec {
	return []ChapterSpec{
		{Number: "01", Title: "System Overview", Catalogs: []string{"ARCHITECTURE-MODULE"}},
		{Number: "02", Title: "Architecture", Catalogs: []string{"ARCHITECTURE-MODULE", "DEPLOYMENT-ELEMENT"}},
		{Number: "03", Title: "Domain and Data", Catalogs: []string{"ENTITY-RECORD", "RELATIONSHIP-RECORD", "STATE-MACHINE", "DATABASE-ROUTINE"}},
		{Number: "04", Title: "Personas, Auth, and Visibility"},
		{Number: "05", Title: "Use Cases", Catalogs: []string{"USE-CASE"}},
		{Number: "06", Title: "Business Rules", Catalogs: []string{"BUSINESS-RULE-RECORD"}},
		{Number: "07", Title: "Interfaces and Integrations", Catalogs: []string{"INTERFACE-RECORD"}},
		{Number: "08", Title: "Background Processing", Catalogs: []string{"SCHEDULE-RECORD", "DATABASE-ROUTINE"}},
		{Number: "09", Title: "Deployment and Configuration", Catalogs: []string{"DEPLOYMENT-ELEMENT", "CONFIGURATION-RECORD"}},
		{Number: "10", Title: "Security, Privacy, and Audit", Catalogs: []string{"SECURITY-FINDING"}},
		{Number: "11", Title: "Non-Functional Profile", Catalogs: []string{"NFR-RECORD"}},
		{Number: "12", Title: "Disabled and Dormant Capabilities"},
		{Number: "13", Title: "Modernization Decisions"},
		{Number: "14", Title: "Equivalence and Acceptance"},
		{Number: "15", Title: "Known Gaps and Risks", Catalogs: []string{"FAULT-RECORD"}},
		{Number: "16", Title: "Glossary and Reference"},
	}
}

// Draft renders the handbook chapters from committed synthesis records and
// returns the Markdown with per-chapter lint findings. The draft binds observed
// records and their evidence; narrative prose stays model-authored later.
func Draft(db string, doc *cdl.EIRDoc) ([]byte, []Finding, error) {
	specs := DefaultChapters()
	var out strings.Builder
	out.WriteString("# Legacy System Reconstruction Handbook\n\n")
	out.WriteString("> Generated reconstruction projection. The canonical protocol state is authoritative.\n\n")
	var findings []Finding
	for i, spec := range specs {
		chapter := Chapter{
			Number:          spec.Number,
			Title:           spec.Title,
			AudiencePurpose: fmt.Sprintf("This chapter shows the %s of the legacy system for maintainers and reviewers.", strings.ToLower(spec.Title)),
		}
		for _, catalogID := range spec.Catalogs {
			records, err := loadRecords(db, doc, catalogID)
			if err != nil {
				return nil, nil, err
			}
			chapter.Observed = append(chapter.Observed, records...)
		}
		if len(chapter.Observed) == 0 {
			chapter.TargetDecisions = "Not-Applicable: this chapter class is empty because no records exist and no evidence supports one."
		} else {
			chapter.TargetDecisions = "Observed records wait for modernization decisions. The decision log records each approved choice."
		}
		if i > 0 {
			chapter.Navigation = append(chapter.Navigation, fmt.Sprintf("previous: %s-%s", specs[i-1].Number, specs[i-1].Title))
		}
		if i < len(specs)-1 {
			chapter.Navigation = append(chapter.Navigation, fmt.Sprintf("next: %s-%s", specs[i+1].Number, specs[i+1].Title))
		}
		body := chapterBody(Assemble([]Chapter{chapter}))
		out.Write(body)
		out.WriteString("\n")
		findings = append(findings, LintChapter(spec.Number, string(body))...)
	}
	return []byte(out.String()), findings, nil
}

// chapterBody removes the assembler's document header from a single-chapter
// render.
func chapterBody(rendered []byte) []byte {
	text := string(rendered)
	if idx := strings.Index(text, "## "); idx >= 0 {
		return []byte(text[idx:])
	}
	return rendered
}

func loadRecords(db string, doc *cdl.EIRDoc, fieldID string) ([]Record, error) {
	names, _, _, err := work.FieldContract(doc, fieldID)
	if err != nil {
		return nil, err
	}
	ids, err := store.List(db, doc, fieldID)
	if err != nil {
		return nil, err
	}
	records := make([]Record, 0, len(ids))
	for _, id := range ids {
		payload, err := store.Get(db, doc, fieldID, id)
		if err != nil {
			return nil, err
		}
		values, err := work.PayloadText(names, payload)
		if err != nil {
			return nil, err
		}
		header := "- **Record:** " + id + " (" + fieldID + ")\n"
		if title := strings.TrimSpace(values["record-title"]); title != "" {
			header = "- **Record:** " + id + ": " + title + "\n"
		}
		records = append(records, Record{ID: id, Payload: []byte(header + string(payload))})
	}
	return records, nil
}
