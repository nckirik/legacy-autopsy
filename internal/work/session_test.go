package work

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestMetricsSidecarRecordsEveryUnit(t *testing.T) {
	db, doc := testDoc(t)
	metricsPath := db + ".metrics.jsonl"
	unit := Unit{ID: "U-METRICS", Mode: "Discovery", Scope: "demo", RecordID: "INV-METRICS", ResultType: "INVARIANT"}
	executor := funcExecutor(func(inv Invocation) (Proposal, error) {
		return invariantProposal("metrics"), nil
	})
	config := DefaultConfig()
	config.MetricsPath = metricsPath
	session, err := NewSession(db, doc, config, executor, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := session.Run([]Unit{unit}); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(metricsPath)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) != 1 {
		t.Fatalf("metrics lines = %d, want 1", len(lines))
	}
	var decoded map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded["unit"] != "U-METRICS" {
		t.Fatalf("metrics unit = %v", decoded["unit"])
	}
	if decoded["result"] != "committed" {
		t.Fatalf("metrics result = %v", decoded["result"])
	}
	if _, ok := decoded["duration_ms"]; !ok {
		t.Fatalf("metrics missing duration_ms")
	}
	if _, ok := decoded["timestamp"]; !ok {
		t.Fatalf("metrics missing timestamp")
	}
}
