package work

import (
	"errors"
	"fmt"

	"github.com/nckirik/legacy-autopsy/cdl"
)

// Invocation is the bounded executor call for one unit: the unit contract, its
// materialized dependency closure, and the capability surface whose read set
// will stale-check the commit.
type Invocation struct {
	Unit    Unit
	Context Context
	Caps    *Caps
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
	return &Session{db: db, doc: doc, config: config, executor: executor, reviewer: reviewer, maxRetries: 2}, nil
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
	if err := unit.Validate(s.doc); err != nil {
		return stats, err
	}
	for attempt := 0; attempt <= s.maxRetries; attempt++ {
		caps := NewCaps(s.db, s.doc)
		context, err := BuildContext(caps, unit)
		if err != nil {
			return stats, err
		}
		stats.ContextBytes = len(context.Bytes)
		stats.Dependencies = len(unit.Reads)
		proposal, err := s.executor.Execute(Invocation{Unit: unit, Context: context, Caps: caps})
		if err != nil {
			var blocker *BlockerError
			if !errors.As(err, &blocker) {
				return stats, err
			}
			class := ClassifyBlocker(blocker.Blocker)
			stats.Blocker = class
			if class == BlockerUncertainty {
				stats.Retries++
				continue
			}
			stats.Escalated = ReachesOperator(class)
			return stats, err
		}
		if s.config.NeedsReview(unit, false) && s.reviewer != nil {
			result, err := s.review(unit, proposal)
			if err != nil {
				return stats, err
			}
			stats.Review = result.Outcome
			if result.Outcome != ReviewApproved {
				stats.Retries++
				continue
			}
		}
		if err := Commit(caps, unit, proposal); err != nil {
			var stale *StaleError
			if errors.As(err, &stale) {
				stats.Retries++
				continue
			}
			return stats, err
		}
		stats.Fetches = caps.FetchCount()
		return stats, nil
	}
	return stats, fmt.Errorf("work: unit %s was not committed within %d attempts", unit.ID, s.maxRetries+1)
}

// review constructs the reviewer's context independently of the worker
// context and asks for a verdict.
func (s *Session) review(unit Unit, proposal Proposal) (ReviewResult, error) {
	target := Ref{FieldID: proposal.FieldID, RecordID: proposal.RecordID}
	reviewCaps := NewCaps(s.db, s.doc)
	context, err := BuildReviewContext(reviewCaps, unit, target, proposal.Payload)
	if err != nil {
		return ReviewResult{}, err
	}
	return s.reviewer.Review(ReviewRequest{Unit: unit, Target: target, Payload: proposal.Payload, Context: context})
}
