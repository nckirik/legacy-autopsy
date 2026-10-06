package work

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/certify"
)

// Invocation is the bounded executor call for one unit: the unit contract, its
// materialized dependency closure, and the capability surface whose read set
// will stale-check the commit.
type Invocation struct {
	Unit           Unit
	Context        Context
	Caps           *Caps
	Fields         []string
	RequiredFields []string
	FieldOptions   map[string][]string
	Hints          []string
}

// Executor performs one semantic work unit and returns a proposal. It never
// writes storage directly.
type Executor interface {
	Execute(Invocation) (Proposal, error)
}

// UnitStats is the per-work-item instrumentation record.
type UnitStats struct {
	UnitID       string
	ContextBytes int
	Dependencies int
	Fetches      int
	Retries      int
	Review       ReviewOutcome
	Escalated    bool
	Blocker      BlockerClass
	Committed    bool
}

// Report aggregates instrumentation for a session.
type Report struct {
	Units []UnitStats
}

// Session runs bounded semantic work units one at a time against the canonical
// store, with read-set provenance, stale-safe commits, optional independent
// review, blocker routing, and per-unit instrumentation.
type Session struct {
	db         string
	doc        *cdl.EIRDoc
	config     Config
	executor   Executor
	reviewer   Reviewer
	maxRetries int
}

// NewSession builds a session. reviewer may be nil when the policy never
// requires review.
func NewSession(db string, doc *cdl.EIRDoc, config Config, executor Executor, reviewer Reviewer) (*Session, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	if executor == nil {
		return nil, fmt.Errorf("work: session requires an executor")
	}
	maxRetries := config.MaxRetries
	if maxRetries == 0 {
		maxRetries = 3
	}
	return &Session{db: db, doc: doc, config: config, executor: executor, reviewer: reviewer, maxRetries: maxRetries}, nil
}

// Run processes units sequentially, one work item at a time, and returns the
// instrumentation report. A human-authority blocker or a protocol/harness
// failure stops the run with an error after recording the unit.
func (s *Session) Run(units []Unit) (Report, error) {
	var report Report
	for _, unit := range units {
		stats, err := s.runUnit(unit)
		report.Units = append(report.Units, stats)
		if err != nil {
			return report, err
		}
	}
	return report, nil
}

func (s *Session) runUnit(unit Unit) (UnitStats, error) {
	stats := UnitStats{UnitID: unit.ID}
	if s.config.Observer != nil && s.config.Observer.UnitStarted != nil {
		s.config.Observer.UnitStarted(unit)
	}
	defer func() {
		if s.config.Observer != nil && s.config.Observer.UnitDone != nil {
			s.config.Observer.UnitDone(stats)
		}
	}()
	if err := unit.Validate(s.doc); err != nil {
		return stats, err
	}
	feedback := ""
	for attempt := 0; attempt <= s.maxRetries; attempt++ {
		caps := NewCaps(s.db, s.doc)
		context, err := BuildContext(caps, unit)
		if err != nil {
			return stats, err
		}
		if feedback != "" {
			context.Bytes = append(context.Bytes, []byte("\n## Previous attempt rejected\n\n"+feedback+"\n")...)
		}
		stats.ContextBytes = len(context.Bytes)
		stats.Dependencies = len(unit.Reads)
		fields, required, options, err := FieldContract(s.doc, unit.ResultType)
		if err != nil {
			return stats, err
		}
		proposal, err := s.executor.Execute(Invocation{Unit: unit, Context: context, Caps: caps, Fields: fields, RequiredFields: required, FieldOptions: options, Hints: certify.ValidationHints(unit.ResultType)})
		if err != nil {
			var blocker *BlockerError
			if !errors.As(err, &blocker) {
				stats.Retries++
				feedback = "executor error: " + err.Error()
				s.notifyRetry(unit, feedback)
				continue
			}
			class := ClassifyBlocker(blocker.Blocker)
			stats.Blocker = class
			if class == BlockerUncertainty {
				stats.Retries++
				s.notifyRetry(unit, "blocker: "+blocker.Blocker.Message)
				continue
			}
			stats.Escalated = ReachesOperator(class)
			return stats, err
		}
		if s.config.NeedsReview(unit, false) && s.reviewer != nil {
			if s.config.Observer != nil && s.config.Observer.Reviewing != nil {
				s.config.Observer.Reviewing(unit)
			}
			result, err := s.review(unit, proposal)
			if err != nil {
				return stats, err
			}
			stats.Review = result.Outcome
			if result.Outcome != ReviewApproved {
				stats.Retries++
				feedback = "review " + string(result.Outcome)
				if len(result.Findings) > 0 {
					feedback += ": " + strings.Join(result.Findings, "; ")
				}
				s.notifyRetry(unit, feedback)
				continue
			}
		}
		if err := Commit(caps, unit, proposal); err != nil {
			var stale *StaleError
			if errors.As(err, &stale) {
				stats.Retries++
				feedback = ""
				s.notifyRetry(unit, "stale dependency "+stale.Ref.String())
				continue
			}
			stats.Retries++
			feedback = err.Error()
			s.notifyRetry(unit, feedback)
			continue
		}
		stats.Fetches = caps.FetchCount()
		stats.Committed = true
		return stats, nil
	}
	return stats, fmt.Errorf("work: unit %s was not committed within %d attempts", unit.ID, s.maxRetries+1)
}

func (s *Session) notifyRetry(unit Unit, reason string) {
	if s.config.Observer != nil && s.config.Observer.Retrying != nil {
		s.config.Observer.Retrying(unit, reason)
	}
}

// review constructs the reviewer's context independently of the worker
// context and asks for a verdict.
func (s *Session) review(unit Unit, proposal Proposal) (ReviewResult, error) {
	resolved, err := ResolveRecords(s.doc, unit, proposal)
	if err != nil {
		return ReviewResult{}, err
	}
	reviewCaps := NewCaps(s.db, s.doc)
	context, err := BuildReviewContext(reviewCaps, unit, resolved)
	if err != nil {
		return ReviewResult{}, err
	}
	request := ReviewRequest{Unit: unit, Context: context}
	for _, record := range resolved {
		request.Targets = append(request.Targets, Ref{FieldID: proposal.FieldID, RecordID: record.RecordID})
		request.Payloads = append(request.Payloads, record.Payload)
	}
	return s.reviewer.Review(request)
}
