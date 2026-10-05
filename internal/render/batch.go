// Package render formats human decision and review surfaces. It is a pure
// formatter: classifications and payloads come from the runtime/store, and the
// output is a display projection, never an operational record.
package render

import (
	"fmt"
	"strings"
)

// Item is one decision or review item in a batch.
type Item struct {
	ID      string
	Summary string
	Payload []byte
}

// Batch is a human-readable decision/review batch per PI-6: a first screen that
// says what changed, what is blocked, and what is needed; a one-line-per-item
// index; then the full items.
type Batch struct {
	Title   string
	Changed []string
	Blocked []string
	Needed  []string
	Items   []Item
}

// Markdown renders the batch deterministically.
func Markdown(batch Batch) []byte {
	var b strings.Builder
	title := batch.Title
	if title == "" {
		title = "Decision batch"
	}
	fmt.Fprintf(&b, "# %s\n\n", title)
	b.WriteString("## First screen\n\n")
	writeList(&b, "Changed", batch.Changed)
	writeList(&b, "Blocked", batch.Blocked)
	writeList(&b, "Needed", batch.Needed)
	b.WriteString("## Index\n\n")
	for _, item := range batch.Items {
		summary := strings.TrimSpace(item.Summary)
		if summary == "" {
			summary = "no summary"
		}
		fmt.Fprintf(&b, "- %s - %s\n", item.ID, summary)
	}
	b.WriteString("\n## Items\n\n")
	for _, item := range batch.Items {
		fmt.Fprintf(&b, "### %s\n\n", item.ID)
		payload := strings.TrimSpace(string(item.Payload))
		if payload == "" {
			payload = "None"
		}
		fmt.Fprintf(&b, "%s\n\n", payload)
	}
	return []byte(b.String())
}

func writeList(b *strings.Builder, label string, values []string) {
	if len(values) == 0 {
		fmt.Fprintf(b, "- **%s:** None\n", label)
		return
	}
	fmt.Fprintf(b, "- **%s:**\n", label)
	for _, value := range values {
		fmt.Fprintf(b, "  - %s\n", value)
	}
}

// Option is one operator choice with its explicit consequence.
type Option struct {
	Label       string
	Consequence string
}

// Question is one plain-language operator question.
type Question struct {
	ID      string
	Prompt  string
	Options []Option
}

// QuestionsMarkdown renders plain-language questions with options and
// consequences for the operator surface.
func QuestionsMarkdown(questions []Question) []byte {
	var b strings.Builder
	b.WriteString("# Operator questions\n\n")
	if len(questions) == 0 {
		b.WriteString("None\n")
		return []byte(b.String())
	}
	for _, question := range questions {
		fmt.Fprintf(&b, "## %s\n\n%s\n\n", question.ID, strings.TrimSpace(question.Prompt))
		for _, option := range question.Options {
			fmt.Fprintf(&b, "- **%s** - %s\n", option.Label, option.Consequence)
		}
		b.WriteString("\n")
	}
	return []byte(b.String())
}
