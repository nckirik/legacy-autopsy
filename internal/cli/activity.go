package cli

import (
	"io"

	"github.com/nckirik/legacy-autopsy/internal/llm"
	"github.com/nckirik/legacy-autopsy/internal/render"
	"github.com/nckirik/legacy-autopsy/internal/work"
)

// wireActivity attaches a stdout progress indicator driven by work lifecycle
// events and model stream deltas, for when raw model output is off.
func wireActivity(config *work.Config, executor *llm.Executor, out io.Writer) {
	activity := render.NewActivity(out)
	config.Observer = &work.Observer{
		UnitStarted: func(unit work.Unit) { activity.Start(unit.ID) },
		Reviewing:   func(work.Unit) { activity.Review() },
		Retrying:    func(_ work.Unit, reason string) { activity.Retry(reason) },
		UnitDone: func(stats work.UnitStats) {
			summary := "stopped"
			if stats.Committed {
				summary = "committed"
			}
			if stats.Escalated {
				summary = "escalated to operator"
			}
			activity.Done(summary, stats.Retries)
		},
	}
	executor.OnDelta = func(_ string, delta llm.Delta) {
		activity.Add(len(delta.Reasoning), len(delta.Content))
	}
}
