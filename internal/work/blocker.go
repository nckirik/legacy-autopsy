package work

import "fmt"

// BlockerClass routes a blocker to the only party that can resolve it. Only
// genuine human authority or domain questions may reach the operator.
type BlockerClass string

const (
	// BlockerHumanAuthority is a decision, scope approval, confirmation, or
	// domain question reserved for an authorized human.
	BlockerHumanAuthority BlockerClass = "human-authority"
	// BlockerUncertainty is executor doubt: retry or independently review.
	BlockerUncertainty BlockerClass = "executor-uncertainty"
	// BlockerProtocol is a protocol interpretation or mechanics question: the
	// runtime/spec owns it.
	BlockerProtocol BlockerClass = "protocol-mechanics"
	// BlockerHarness is a tooling, environment, or harness defect: block.
	BlockerHarness BlockerClass = "harness-defect"
)

// ReachesOperator reports whether a class may become Human Hatch traffic.
func ReachesOperator(class BlockerClass) bool { return class == BlockerHumanAuthority }

// Blocker is an executor-reported obstruction with a closed kind.
type Blocker struct {
	Kind    string
	Message string
}

// ClassifyBlocker maps a declared blocker kind to its routing class. Unknown
// kinds fail toward protocol mechanics rather than reaching a human.
func ClassifyBlocker(blocker Blocker) BlockerClass {
	switch blocker.Kind {
	case "human-authority", "decision-approval", "scope-approval", "confirmation", "domain-question":
		return BlockerHumanAuthority
	case "uncertainty", "unresolved-evidence":
		return BlockerUncertainty
	case "harness", "tooling", "environment":
		return BlockerHarness
	default:
		return BlockerProtocol
	}
}

// BlockerError is the executor error carrying a blocker.
type BlockerError struct {
	Blocker Blocker
}

func (e *BlockerError) Error() string {
	return fmt.Sprintf("work: blocker %s: %s", e.Blocker.Kind, e.Blocker.Message)
}
