package llm

import (
	"context"

	"github.com/nckirik/legacy-autopsy/internal/work"
)

// Reviewer adapts a chat client to the independent-review interface.
type Reviewer struct {
	client *Client
}

// NewReviewer wraps a client as a work reviewer.
func NewReviewer(client *Client) *Reviewer { return &Reviewer{client: client} }

// NewReviewerFromEnv builds a reviewer from environment configuration.
func NewReviewerFromEnv() (*Reviewer, error) {
	config, err := ConfigFromEnv()
	if err != nil {
		return nil, err
	}
	client, err := NewClient(config)
	if err != nil {
		return nil, err
	}
	return NewReviewer(client), nil
}

// Review asks the model to judge the proposed record from its independently
// constructed context.
func (r *Reviewer) Review(request work.ReviewRequest) (work.ReviewResult, error) {
	text, err := r.client.Complete(context.Background(), reviewerSystemPrompt, string(request.Context.Bytes))
	if err != nil {
		return work.ReviewResult{}, err
	}
	return ParseReview(text)
}
