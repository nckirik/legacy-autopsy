package handbook

import "fmt"

// Binding is one exact chapter-to-payload binding at candidate time.
type Binding struct {
	HBKID       string
	Fingerprint string
}

// Review is one readability CNF with its shared schema and chapter bindings.
type Review struct {
	ReadabilityReview
	Bindings []Binding
}

// VerifyReadabilityCoverage validates the PI-12 readability workflow
// deterministically: every review must pass the CNF schema, and the combined
// Pass/Fail reviews must bind every candidate chapter exactly once with the
// matching candidate fingerprint. Unknown or duplicate bindings, fingerprint
// drift, and uncovered chapters fail closed. Not-Applicable reviews contribute
// no coverage evidence.
func VerifyReadabilityCoverage(subject string, chapters []Binding, reviews []Review) []Finding {
	var findings []Finding
	byID := map[string]string{}
	for _, chapter := range chapters {
		if chapter.HBKID == "" {
			findings = append(findings, Finding{Code: "HBK-COVERAGE", Subject: subject, Detail: "candidate chapter has an empty HBK ID"})
			continue
		}
		if _, ok := byID[chapter.HBKID]; ok {
			findings = append(findings, Finding{Code: "HBK-COVERAGE", Subject: subject, Detail: "duplicate candidate chapter " + chapter.HBKID})
		}
		byID[chapter.HBKID] = chapter.Fingerprint
	}
	covered := map[string]bool{}
	for i, review := range reviews {
		label := fmt.Sprintf("review %d", i+1)
		schema := review.ReadabilityReview
		schema.HBKBindings = nil
		for _, binding := range review.Bindings {
			schema.HBKBindings = append(schema.HBKBindings, binding.HBKID)
		}
		if reviewFindings := ValidateReadabilityReview(label, schema); len(reviewFindings) > 0 {
			findings = append(findings, reviewFindings...)
			continue
		}
		if review.Value == "Not-Applicable" {
			if len(review.Bindings) > 0 {
				findings = append(findings, Finding{Code: "HBK-COVERAGE", Subject: label, Detail: "Not-Applicable review must not bind HBK payloads"})
			}
			continue
		}
		for _, binding := range review.Bindings {
			want, ok := byID[binding.HBKID]
			if !ok {
				findings = append(findings, Finding{Code: "HBK-COVERAGE", Subject: label, Detail: "binds unknown or non-candidate chapter " + binding.HBKID})
				continue
			}
			if binding.Fingerprint != want {
				findings = append(findings, Finding{Code: "HBK-COVERAGE", Subject: label, Detail: "fingerprint drift for chapter " + binding.HBKID})
				continue
			}
			if covered[binding.HBKID] {
				findings = append(findings, Finding{Code: "HBK-COVERAGE", Subject: label, Detail: "chapter " + binding.HBKID + " is bound more than once"})
				continue
			}
			covered[binding.HBKID] = true
		}
	}
	for id := range byID {
		if !covered[id] {
			findings = append(findings, Finding{Code: "HBK-COVERAGE", Subject: subject, Detail: "chapter " + id + " is not covered by a readability review"})
		}
	}
	return findings
}
