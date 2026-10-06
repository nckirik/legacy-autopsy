package render

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var spinnerFrames = []string{"|", "/", "-", "\\"}

// Activity is a single-line progress indicator for long model calls. On a
// terminal it redraws in place as lifecycle events and stream deltas arrive; on
// a non-terminal it stays quiet during the call and prints one line per unit on
// completion, so logs and pipes stay readable.
type Activity struct {
	out       io.Writer
	tty       bool
	unit      string
	reasoning int
	content   int
	frame     int
	start     time.Time
	last      time.Time
	reviewing bool
}

// NewActivity builds an activity indicator that detects whether out is a
// terminal.
func NewActivity(out io.Writer) *Activity { return newActivity(out, isTerminal(out)) }

func newActivity(out io.Writer, tty bool) *Activity { return &Activity{out: out, tty: tty} }

func isTerminal(out io.Writer) bool {
	file, ok := out.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

// Start begins a unit, resetting counters.
func (a *Activity) Start(unit string) {
	a.unit = unit
	a.reasoning = 0
	a.content = 0
	a.reviewing = false
	a.start = time.Now()
	a.last = time.Time{}
	if a.tty {
		a.draw()
	}
}

// Add accumulates streamed output sizes and repaints at a bounded rate.
func (a *Activity) Add(reasoning, content int) {
	a.reasoning += reasoning
	a.content += content
	if a.tty && time.Since(a.last) >= 80*time.Millisecond {
		a.draw()
	}
}

// Review marks that the unit is now under independent review.
func (a *Activity) Review() {
	a.reviewing = true
	if a.tty {
		a.draw()
	}
}

// Retry replaces the indicator with a one-line retry notice.
func (a *Activity) Retry(reason string) {
	if !a.tty {
		return
	}
	fmt.Fprintf(a.out, "\r%s%s retry: %s\n", clearLine(), spinner(), truncate(reason, 96))
	a.reasoning = 0
	a.content = 0
	a.start = time.Now()
}

// Done finishes the unit with a one-line summary.
func (a *Activity) Done(summary string, retries int) {
	elapsed := time.Since(a.start).Round(time.Second)
	detail := summary
	if retries > 0 {
		detail = fmt.Sprintf("%s, %d retries", summary, retries)
	}
	if a.tty {
		fmt.Fprintf(a.out, "\r%s%s %s: %s (%s)\n", clearLine(), spinner(), a.unit, detail, elapsed)
		return
	}
	fmt.Fprintf(a.out, "unit %s: %s (%s)\n", a.unit, detail, elapsed)
}

func (a *Activity) draw() {
	a.last = time.Now()
	a.frame++
	status := ""
	if a.reviewing {
		status = " reviewing"
	}
	fmt.Fprintf(a.out, "\r%s%s %s%s reasoning %s content %s %s",
		clearLine(), spinner(), a.unit, status,
		humanCount(a.reasoning), humanCount(a.content), time.Since(a.start).Round(time.Second))
}

func spinner() string {
	return spinnerFrames[int(time.Now().UnixNano()/int64(250*time.Millisecond))%len(spinnerFrames)]
}

func clearLine() string { return "\033[K" }

func humanCount(n int) string {
	switch {
	case n >= 1000000:
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	case n >= 1000:
		return fmt.Sprintf("%.1fk", float64(n)/1000)
	default:
		return fmt.Sprintf("%d", n)
	}
}

func truncate(value string, max int) string {
	value = strings.Join(strings.Fields(value), " ")
	if len(value) <= max {
		return value
	}
	return value[:max] + "..."
}
