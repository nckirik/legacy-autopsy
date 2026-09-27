package capabilities

import "testing"

func TestNextIteration(t *testing.T) {
	tokens := []string{"ALFA", "BRAVO", "CHARLIE"}
	next, err := NextIteration(tokens, "ALFA")
	if err != nil || next != "BRAVO" {
		t.Fatalf("ALFA -> %s, %v", next, err)
	}
	if _, err := NextIteration(tokens, "CHARLIE"); err == nil {
		t.Fatal("final token did not fail closed")
	}
	if _, err := NextIteration(tokens, "UNKNOWN"); err == nil {
		t.Fatal("unknown token accepted")
	}
	if _, err := NextIteration(nil, "ALFA"); err == nil {
		t.Fatal("empty sequence accepted")
	}
}

func TestNextIterationFullBudget(t *testing.T) {
	tokens := []string{
		"ALFA", "BRAVO", "CHARLIE", "DELTA", "ECHO", "FOXTROT", "GOLF", "HOTEL",
		"INDIA", "JULIETT", "KILO", "LIMA", "MIKE", "NOVEMBER", "OSCAR", "PAPA",
		"QUEBEC", "ROMEO", "SIERRA", "TANGO", "UNIFORM", "VICTOR", "WHISKEY",
		"X-RAY", "YANKEE", "ZULU",
	}
	current := "ALFA"
	steps := 0
	for {
		next, err := NextIteration(tokens, current)
		if err != nil {
			break
		}
		current = next
		steps++
	}
	if current != "ZULU" || steps != 25 {
		t.Fatalf("advanced %d steps to %s, want 25 steps to ZULU", steps, current)
	}
}
