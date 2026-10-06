package work

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/store"
)

// invocationID mints the next protocol invocation identity
// ITERATION-PERSONA-POV-SEQUENCE. The sequence continues from the committed
// invocation log, so identities are never reused across sessions or runs.
func (s *Session) invocationID(unit Unit) (string, error) {
	iteration := s.config.Iteration
	if iteration == "" {
		iteration = "ALFA"
	}
	persona := strings.ToUpper(strings.TrimSpace(unit.Persona))
	if persona == "" {
		persona = "NONE"
	}
	pov := strings.ToUpper(strings.TrimSpace(unit.POV))
	if pov == "" {
		pov = "NONE"
	}
	prefix := iteration + "-" + persona + "-" + pov + "-"
	ids, err := store.List(s.db, s.doc, "INVOCATION-LOG")
	if err != nil {
		return "", err
	}
	highest := 0
	for _, id := range ids {
		if !strings.HasPrefix(id, prefix) {
			continue
		}
		sequence, err := strconv.Atoi(strings.TrimPrefix(id, prefix))
		if err == nil && sequence > highest {
			highest = sequence
		}
	}
	return fmt.Sprintf("%s%03d", prefix, highest+1), nil
}

// appendInvocationLog records the mandatory `0G` entry for one finished
// invocation. It is runtime-authored; executors never write the log directly.
func (s *Session) appendInvocationLog(unit Unit, invocationID string, stats UnitStats) error {
	summary := "no-mutation"
	mutation := "None"
	if stats.Committed {
		summary = "committed"
		mutation = fmt.Sprintf("committed %d record(s)", len(stats.MutationIDs))
	}
	if stats.LastError != "" {
		summary = "stopped"
	}
	var b strings.Builder
	write := func(name, value string) {
		fmt.Fprintf(&b, "- **%s:** %s\n", name, value)
	}
	write("identity-header-summary", fmt.Sprintf("%s mode=%s scope=%s result=%s", invocationID, unit.Mode, unit.Scope, unit.ResultType))
	write("cold-resume-result-check-fingerprint", "None")
	write("declared-semantic-actions-runtime-authorization-result", fmt.Sprintf("%s (%d retries)", summary, stats.Retries))
	write("loaded-files-fingerprints", joinOrNone(stats.LoadedFingerprints))
	write("mutation-summary", mutation)
	write("tickets-created-updated", "None")
	write("frontiers-created-closed", "None")
	write("unmapped-discoveries-gap-records-buffered-merged-or-disposed", "None")
	write("exports-acquired-requested", "None")
	write("src-cmp-clm-ids-updated", joinOrNone(stats.MutationIDs))
	write("evidence-changes", "None")
	write("promotion-depromotion", "None")
	write("coverage-reconciliation-changes", "None")
	write("stale-propagation", "None")
	write("next-required-loads", "None")
	return store.Ingest(s.db, s.doc, "INVOCATION-LOG", invocationID, 1, []byte(b.String()))
}

func joinOrNone(values []string) string {
	if len(values) == 0 {
		return "None"
	}
	return strings.Join(values, ", ")
}
