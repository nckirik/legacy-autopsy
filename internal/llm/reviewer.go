package llm

import (
	"context"
	"fmt"
	"io"

	"github.com/nckirik/legacy-autopsy/internal/work"
)

// Reviewer adapts a chat client to the independent-review interface. When
// Debug is set, every raw model verdict is written there for inspection.
type Reviewer struct {
	client *Client
	Debug  io.Writer
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
	if r.Debug != nil {
		fmt.Fprintf(r.Debug, "\n===== review response for %s =====\n%s\n===== end review %s =====\n", request.Unit.ID, text, request.Unit.ID)
	}
	if err != nil {
		return work.ReviewResult{}, err
	}
	return ParseReview(text)
}
