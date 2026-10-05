package work

// ReviewOutcome is the closed verdict set for an independent review.
type ReviewOutcome string

const (
	ReviewApproved     ReviewOutcome = "approved"
	ReviewRejected     ReviewOutcome = "rejected"
	ReviewInconclusive ReviewOutcome = "inconclusive"
)

// ReviewResult is a reviewer verdict with its findings.
type ReviewResult struct {
	Outcome  ReviewOutcome
	Findings []string
}

// ReviewRequest is what a reviewer receives: the unit contract, the proposed
// record, and a context constructed independently of the worker invocation.
type ReviewRequest struct {
	Unit    Unit
	Target  Ref
	Payload []byte
	Context Context
}

// Reviewer independently challenges a worker proposal. It MAY use the same
// model; independence comes from a separately constructed context and role.
type Reviewer interface {
	Review(ReviewRequest) (ReviewResult, error)
}
