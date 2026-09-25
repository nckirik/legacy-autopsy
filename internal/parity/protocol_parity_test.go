package parity

import (
	"strings"
	"testing"

	"github.com/nckirik/legacy-autopsy/cdl"
	"github.com/nckirik/legacy-autopsy/internal/fixtures"
	"github.com/nckirik/legacy-autopsy/internal/runtime"
)

// TestProtocolOracleUnchanged is test B's oracle half: the independent 24-case
// runner must remain unchanged and green.
func TestProtocolOracleUnchanged(t *testing.T) {
	summary, err := fixtures.Run(root(t))
	if err != nil {
		t.Fatalf("fixture oracle: %v", err)
	}
	if summary.Failed != 0 {
		t.Fatalf("oracle failures: %d", summary.Failed)
	}
	if summary.Passed != 24 {
		t.Fatalf("oracle fixture count changed: %d passed", summary.Passed)
	}
}

func executeFixture(t *testing.T, res *cdl.Result, name string) runtime.Trace {
	t.Helper()
	fx, err := runtime.LoadFixture(fixturePath(t, name))
	if err != nil {
		t.Fatal(err)
	}
	trace, err := runtime.Execute(res.EIR, fx, capabilityMap(fx.Capabilities))
	if err != nil {
		t.Fatal(err)
	}
	return trace
}

func record(t *testing.T, trace runtime.Trace, step string) runtime.StepRecord {
	t.Helper()
	for _, rec := range trace.Steps {
		if rec.Step == step {
			return rec
		}
	}
	t.Fatalf("no trace record for %s", step)
	return runtime.StepRecord{}
}

func blockedGates(rec runtime.StepRecord) []string {
	var out []string
	for _, e := range rec.Effects {
		if e.Operation == "block" {
			out = append(out, e.Target)
		}
	}
	return out
}

func setValue(rec runtime.StepRecord, target string) (string, bool) {
	for _, e := range rec.Effects {
		if e.Operation == "set" && e.Target == target {
			return e.Value, true
		}
	}
	return "", false
}

// TestBehaviorParity is test B's behavior half: EIR execution must reflect the
// corrected §7.7 semantics recorded in protocol.md@4.1.2.
func TestBehaviorParity(t *testing.T) {
	res := compileExample(t)

	classify := findStep(t, res.EIR, "reconcile.classify")
	if !strings.Contains(classify.Require, "Container mappings do not reconcile their children") {
		t.Fatalf("agent requirement clause lost: %q", classify.Require)
	}
	guard := findStep(t, res.EIR, "reconcile.guard")
	if len(guard.Ops) == 0 {
		t.Fatal("guard step has no operations")
	}

	t.Run("reconciled completes", func(t *testing.T) {
		trace := executeFixture(t, res, "reconciled.json")
		guardRec := record(t, trace, "reconcile.guard")
		if guardRec.Outcome != "committed" {
			t.Fatalf("guard outcome %s", guardRec.Outcome)
		}
		if len(blockedGates(guardRec)) != 0 {
			t.Fatalf("unexpected blocks: %v", blockedGates(guardRec))
		}
		fingerprintRec := record(t, trace, "reconcile.fingerprint")
		if v, ok := setValue(fingerprintRec, "fingerprint-status"); !ok || v != "computed" {
			t.Fatalf("fingerprint-status = %q (present=%v)", v, ok)
		}
	})

	t.Run("gap blocks completion", func(t *testing.T) {
		trace := executeFixture(t, res, "incomplete.json")
		guardRec := record(t, trace, "reconcile.guard")
		if guardRec.Outcome != "blocked" {
			t.Fatalf("guard outcome %s", guardRec.Outcome)
		}
		want := []string{"A-STRUCTURALLY-COMPLETE", "R-SWEPT", "Exit-A"}
		got := blockedGates(guardRec)
		if strings.Join(got, ",") != strings.Join(want, ",") {
			t.Fatalf("blocked gates %v, want %v (record %+v)", got, want, guardRec)
		}
		found := false
		for _, d := range guardRec.Diagnostics {
			if d.Code == "GUARD_BLOCKED" {
				found = true
			}
		}
		if !found {
			t.Fatalf("missing GUARD_BLOCKED diagnostic: %v", guardRec.Diagnostics)
		}
	})

	t.Run("container mapping does not reconcile", func(t *testing.T) {
		trace := executeFixture(t, res, "container-only.json")
		guardRec := record(t, trace, "reconcile.guard")
		if guardRec.Outcome != "blocked" {
			t.Fatalf("container-only guard outcome %s; a container mapping must not reconcile", guardRec.Outcome)
		}
	})

	t.Run("unavailable capability blocks only dependents", func(t *testing.T) {
		trace := executeFixture(t, res, "hash-unavailable.json")
		for _, step := range []string{"reconcile.populate", "reconcile.classify", "reconcile.count", "reconcile.guard"} {
			if rec := record(t, trace, step); rec.Outcome != "committed" {
				t.Fatalf("%s outcome %s; independent work must continue", step, rec.Outcome)
			}
		}
		fingerprintRec := record(t, trace, "reconcile.fingerprint")
		if fingerprintRec.Outcome != "committed" {
			t.Fatalf("fingerprint outcome %s", fingerprintRec.Outcome)
		}
		if v, ok := setValue(fingerprintRec, "fingerprint-status"); !ok || v != "unsupported" {
			t.Fatalf("fingerprint-status = %q (present=%v)", v, ok)
		}
		want := map[string]bool{"artifact-fingerprint": true, "artifact-completeness": true}
		got := blockedGates(fingerprintRec)
		if len(got) != len(want) {
			t.Fatalf("blocked gates %v, want %v", got, want)
		}
		for _, g := range got {
			if !want[g] {
				t.Fatalf("unexpected blocked gate %s", g)
			}
		}
	})

	t.Run("authority follows channel semantics", func(t *testing.T) {
		trace := executeFixture(t, res, "reconciled.json")
		for _, rec := range trace.Steps {
			switch rec.Owner {
			case "AGENT":
				if rec.Authority != "draft" {
					t.Fatalf("AGENT %s authority %s", rec.Step, rec.Authority)
				}
			case "MACHINE":
				if rec.Authority != "authoritative" {
					t.Fatalf("MACHINE %s authority %s", rec.Step, rec.Authority)
				}
			}
		}
	})
}

func findStep(t *testing.T, doc *cdl.EIRDoc, id string) cdl.EIRStep {
	t.Helper()
	for _, sec := range doc.Sections {
		for _, step := range sec.Steps {
			if step.ID == id {
				return step
			}
		}
	}
	t.Fatalf("no EIR step %s", id)
	return cdl.EIRStep{}
}
