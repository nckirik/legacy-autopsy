package render

import (
	"bytes"
	"strings"
	"testing"
)

func TestActivityInteractiveRedraws(t *testing.T) {
	var out bytes.Buffer
	activity := newActivity(&out, true)
	activity.Start("U-1")
	activity.Add(1500, 200)
	activity.Review()
	activity.Retry("review rejected: missing evidence")
	activity.Done("committed", 1)
	text := out.String()
	for _, marker := range []string{"U-1", "reasoning 1.5k", "reviewing", "retry: review rejected", "committed, 1 retries"} {
		if !strings.Contains(text, marker) {
			t.Fatalf("activity output missing %q:\n%q", marker, text)
		}
	}
	if !strings.Contains(text, clearLine()) {
		t.Fatal("interactive activity never cleared its line")
	}
	if !strings.HasSuffix(text, "\n") {
		t.Fatal("activity did not finish with a newline")
	}
}

func TestActivityNonTerminalIsQuietUntilDone(t *testing.T) {
	var out bytes.Buffer
	activity := newActivity(&out, false)
	activity.Start("U-2")
	activity.Add(5000, 100)
	activity.Review()
	activity.Retry("ignored")
	if out.Len() != 0 {
		t.Fatalf("non-terminal activity emitted progress: %q", out.String())
	}
	activity.Done("stopped", 0)
	text := out.String()
	if !strings.Contains(text, "unit U-2: stopped") {
		t.Fatalf("non-terminal completion missing: %q", text)
	}
	if strings.Contains(text, "\033") {
		t.Fatalf("non-terminal output contains ANSI escapes: %q", text)
	}
}
