package decision

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/capabilities"
	"github.com/nckirik/legacy-autopsy/internal/store"
)

// CommitApproved creates one canonical DECISION record per approved draft. The
// decision-content fingerprint is computed under §5.6 (carrier and approval
// envelope excluded) before the record is stored; the approval fields come from
// the recorded human review, never from the model.
func CommitApproved(db string, doc *cdl.EIRDoc, namespace string, entries []DraftEntry, review Review) ([]string, error) {
	if len(entries) == 0 {
		return nil, fmt.Errorf("decision commit: no draft entries to commit")
	}
	bySubject := map[string]DraftEntry{}
	for _, entry := range entries {
		bySubject[entry.SubjectID] = entry
	}
	var created []string
	for _, row := range review.Rows {
		if row.Verdict != "Approve" {
			continue
		}
		entry, ok := bySubject[row.SubjectID]
		if !ok {
			return nil, fmt.Errorf("decision commit: review approved %s but the draft file has no such entry", row.SubjectID)
		}
		recordID, err := capabilities.TypedID("DEC", namespace, entry.SubjectID, entry.DecisionType)
		if err != nil {
			return nil, err
		}
		payload := buildDecisionPayload(entry, review)
		fingerprint, err := capabilities.DecisionContentFingerprint([]byte(payload), "decision-content-fingerprint")
		if err != nil {
			return nil, err
		}
		payload = strings.Replace(payload, "- **decision-content-fingerprint:** PENDING", "- **decision-content-fingerprint:** "+fingerprint, 1)
		payload = strings.Replace(payload, "- **approved-content-binding:** PENDING", "- **approved-content-binding:** "+recordID+" version 1 "+fingerprint, 1)
		if err := store.Ingest(db, doc, "DECISION", recordID, 1, []byte(payload)); err != nil {
			return nil, fmt.Errorf("decision commit: %s: %w", recordID, err)
		}
		created = append(created, recordID)
	}
	sort.Strings(created)
	return created, nil
}

func buildDecisionPayload(entry DraftEntry, review Review) string {
	sections := entry.Sections
	combine := func(parts ...string) string {
		var values []string
		for _, part := range parts {
			part = strings.Join(strings.Fields(part), " ")
			if part != "" {
				values = append(values, part)
			}
		}
		return strings.Join(values, " | ")
	}
	var b strings.Builder
	b.WriteString("- **decision-content-version:** 1\n")
	b.WriteString("- **decision-content-fingerprint:** PENDING\n")
	fmt.Fprintf(&b, "- **type:** %s\n", entry.DecisionType)
	fmt.Fprintf(&b, "- **observed-inputs:** %s\n", strings.Join(entry.ObservedLines, "; "))
	fmt.Fprintf(&b, "- **options-decision-rationale-risks:** %s\n", combine(
		"Options: "+sections["Options"],
		"Decision: "+sections["Decision"],
		"Rationale: "+sections["Rationale"],
		"Risks: "+sections["Risks"],
	))
	fmt.Fprintf(&b, "- **exact-scope-and-denominator-effect:** %s\n", combine(sections["Exact Scope and Denominator Effect"]))
	fmt.Fprintf(&b, "- **affected-records-rollout-rollback:** %s\n", combine(sections["Affected Records / Rollout / Rollback"]))
	b.WriteString("- **content-supersession:** None\n\n")
	b.WriteString("<!-- DECISION-APPROVAL-ENVELOPE: excluded from decision-content version/fingerprint -->\n\n")
	b.WriteString("- **approval-status:** Approved\n")
	b.WriteString("- **approval-envelope-version:** 1\n")
	fmt.Fprintf(&b, "- **owner-approver-authority-basis-approval-timestamp:** %s\n", combine(
		review.Reviewer,
		"scope: "+review.Scope,
		"limitations: "+review.Limitations,
		review.Timestamp,
	))
	b.WriteString("- **approved-content-binding:** PENDING\n")
	fmt.Fprintf(&b, "- **approval-signature-audit:** decision review by %s at %s\n", review.Reviewer, review.Timestamp)
	b.WriteString("<!-- END-DECISION-APPROVAL-ENVELOPE -->\n")
	return b.String()
}
