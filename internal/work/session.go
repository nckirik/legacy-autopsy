package work

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/certify"
)

// Invocation is the bounded executor call for one unit: the unit contract, its
// materialized dependency closure, and the capability surface whose read set
// will stale-check the commit.
type Invocation struct {
	Unit           Unit
	InvocationID   string
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
	DurationMS   int64
	ContextBytes int
	Dependencies int
	Fetches      int
	Retries      int
	Review       ReviewOutcome
	Escalated    bool
	Blocker      BlockerClass
	Committed    bool
	LastError    string
	ZeroDomain   string
	// InvocationID is the protocol invocation identity ITERATION-PERSONA-POV-SEQUENCE.
	InvocationID       string
	LoadedFingerprints []string
	MutationIDs        []string
}

// Report aggregates instrumentation for a session.
type Report struct {
	Units []UnitStats
}

// Session runs bounded semantic work units one at a time against the canonical
// store, with read-set provenance, stale-safe commits, optional independent
// review, blocker routing, and per-unit instrumentation.
type Session struct {
	db          string
	doc         *cdl.EIRDoc
	config      Config
	executor    Executor
	reviewer    Reviewer
	maxRetries  int
	metricsPath string
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
	metricsPath := config.MetricsPath
	if metricsPath == "" {
		metricsPath = db + ".metrics.jsonl"
	}
	return &Session{db: db, doc: doc, config: config, executor: executor, reviewer: reviewer, maxRetries: maxRetries, metricsPath: metricsPath}, nil
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

func (s *Session) runUnit(unit Unit) (stats UnitStats, err error) {
	started := time.Now()
	stats = UnitStats{UnitID: unit.ID}
	if err := unit.Validate(s.doc); err != nil {
		return stats, err
	}
	invocationID, err := s.invocationID(unit)
	if err != nil {
		return stats, err
	}
	stats.InvocationID = invocationID
	if s.config.Observer != nil && s.config.Observer.UnitStarted != nil {
		s.config.Observer.UnitStarted(unit)
	}
	defer func() {
		stats.DurationMS = time.Since(started).Milliseconds()
		if logErr := s.appendInvocationLog(unit, invocationID, stats); logErr != nil && err == nil {
			err = logErr
		}
		s.appendMetrics(unit, stats, err)
		if s.config.Observer != nil && s.config.Observer.UnitDone != nil {
			s.config.Observer.UnitDone(stats)
		}
	}()
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
		proposal, err := s.executor.Execute(Invocation{Unit: unit, InvocationID: invocationID, Context: context, Caps: caps, Fields: fields, RequiredFields: required, FieldOptions: options, Hints: certify.ValidationHints(unit.ResultType)})
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
		if proposal.ZeroDomain != nil {
			stats.ZeroDomain = proposal.ZeroDomain.Reason
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
		committedIDs, commitErr := CommitWithIDs(caps, unit, proposal, s.config.UpdateExisting)
		if err := commitErr; err != nil {
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
		stats.MutationIDs = committedIDs
		var loaded []string
		reads := caps.ReadSet()
		refs := make([]Ref, 0, len(reads))
		for ref := range reads {
			refs = append(refs, ref)
		}
		sort.Slice(refs, func(i, j int) bool { return refs[i].String() < refs[j].String() })
		for _, ref := range refs {
			loaded = append(loaded, ref.String()+"@"+reads[ref])
		}
		stats.LoadedFingerprints = loaded
		return stats, nil
	}
	stats.LastError = feedback
	return stats, fmt.Errorf("work: unit %s was not committed within %d attempts; last rejection: %s", unit.ID, s.maxRetries+1, feedback)
}

func (s *Session) notifyRetry(unit Unit, reason string) {
	if s.config.Observer != nil && s.config.Observer.Retrying != nil {
		s.config.Observer.Retrying(unit, reason)
	}
}

// review constructs the reviewer's context independently of the worker
// context and asks for a verdict.
func (s *Session) review(unit Unit, proposal Proposal) (ReviewResult, error) {
	if proposal.ZeroDomain != nil {
		context := Context{UnitID: unit.ID, Bytes: []byte(
			"## Zero-domain declaration\n\n- **Result contract:** " + unit.ResultType +
				"\n- **Reason:** " + proposal.ZeroDomain.Reason +
				"\n- **Evidence:** " + strings.Join(proposal.ZeroDomain.Evidence, ", ") + "\n")}
		return s.reviewer.Review(ReviewRequest{Unit: unit, ZeroDomain: proposal.ZeroDomain, Context: context})
	}
	sessionCaps := NewCaps(s.db, s.doc)
	resolved, err := ResolveRecords(s.doc, unit, proposal, storeSourceResolver(sessionCaps))
	if err != nil {
		return ReviewResult{}, err
	}
	// Show the reviewer the records as they will be committed: runtime-owned
	// fields carry their resolved values, not placeholders.
	if len(unit.RuntimeFields) > 0 {
		for i := range resolved {
			injected, err := injectRuntimeFields(resolved[i].Payload, unit.RuntimeFields, resolved[i].RecordID)
			if err == nil {
				resolved[i].Payload = injected
			}
		}
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

// unitMetrics is one non-authoritative telemetry line. It never carries record
// content and is not part of any semantic hash.
type unitMetrics struct {
	Timestamp    string `json:"timestamp"`
	Unit         string `json:"unit"`
	Mode         string `json:"mode"`
	Invocation   string `json:"invocation"`
	Result       string `json:"result"`
	DurationMS   int64  `json:"duration_ms"`
	Retries      int    `json:"retries"`
	ContextBytes int    `json:"context_bytes"`
	Dependencies int    `json:"dependencies"`
	Fetches      int    `json:"fetches"`
	Review       string `json:"review"`
	Escalated    bool   `json:"escalated"`
	Error        string `json:"error,omitempty"`
}

// appendMetrics writes one telemetry line to the non-authoritative sink. Write
// failures never change the execution outcome.
func (s *Session) appendMetrics(unit Unit, stats UnitStats, runErr error) {
	result := "no-op"
	switch {
	case runErr != nil:
		result = "failed"
	case stats.ZeroDomain != "":
		result = "zero-domain"
	case stats.Committed:
		result = "committed"
	}
	line := unitMetrics{
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		Unit:         stats.UnitID,
		Mode:         unit.Mode,
		Invocation:   stats.InvocationID,
		Result:       result,
		DurationMS:   stats.DurationMS,
		Retries:      stats.Retries,
		ContextBytes: stats.ContextBytes,
		Dependencies: stats.Dependencies,
		Fetches:      stats.Fetches,
		Review:       string(stats.Review),
		Escalated:    stats.Escalated,
	}
	if runErr != nil {
		line.Error = runErr.Error()
	}
	encoded, err := json.Marshal(line)
	if err != nil {
		return
	}
	file, err := os.OpenFile(s.metricsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer file.Close()
	_, _ = file.Write(append(encoded, '\n'))
}
