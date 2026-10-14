package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nckirik/legacy-autopsy/internal/handbook"
	"github.com/nckirik/legacy-autopsy/internal/llm"
	"github.com/nckirik/legacy-autopsy/internal/store"
	"github.com/nckirik/legacy-autopsy/internal/synthesis"
)

// specHandbook renders the human handbook draft from committed synthesis records.
func specHandbook(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: legacy-autopsy spec handbook <draft|narrate> [flags]")
	}
	if args[0] == "narrate" {
		return specHandbookNarrate(args[1:], stdout, stderr)
	}
	if args[0] == "index" {
		return specHandbookIndex(args[1:], stdout, stderr)
	}
	if args[0] == "evidence" {
		return specHandbookEvidence(args[1:], stdout, stderr)
	}
	if args[0] == "mint" {
		return specHandbookMint(args[1:], stdout, stderr)
	}
	if args[0] != "draft" {
		return fmt.Errorf("spec handbook operation %q is unsupported; supported operations: draft, narrate, index, evidence, mint", args[0])
	}
	set := flags("spec handbook draft", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	out := set.String("out", "", "write the draft to this path (default stdout)")
	if err := set.Parse(args[1:]); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" {
		return fmt.Errorf("spec handbook draft requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	draft, findings, err := handbook.Draft(*db, res.EIR)
	if err != nil {
		return err
	}
	if *out != "" {
		if err := os.WriteFile(*out, draft, 0o644); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "handbook draft: %s\n", *out)
	} else {
		stdout.Write(draft)
	}
	fmt.Fprintf(stdout, "lint findings: %d\n", len(findings))
	for _, finding := range findings {
		fmt.Fprintf(stdout, "  %s\n", finding.String())
	}
	return nil
}

// specHandbookNarrate authors one chapter's narrative with STE lint feedback.
func specHandbookNarrate(args []string, stdout, stderr io.Writer) error {
	set := flags("spec handbook narrate", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	chapter := set.String("chapter", "", "chapter number (for example 01)")
	out := set.String("out", "", "write the chapter to this path (default stdout)")
	llmURL := set.String("llm-url", "", "override "+llm.EnvURL)
	llmModel := set.String("llm-model", "", "override "+llm.EnvModel)
	llmReasoning := set.String("llm-reasoning", "", "override "+llm.EnvReasoning)
	feedback := set.String("feedback", "", "reviewer feedback to address in this revision")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" || *chapter == "" {
		return fmt.Errorf("spec handbook narrate requires --db and --chapter")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	config, err := llm.ConfigFromEnv()
	if err != nil {
		return err
	}
	if *llmURL != "" {
		config.BaseURL = *llmURL
	}
	if *llmModel != "" {
		config.Model = *llmModel
	}
	effort := config.Reasoning
	if *llmReasoning != "" {
		effort = *llmReasoning
	}
	client, err := llm.NewClient(config)
	if err != nil {
		return err
	}
	narrative, err := handbook.Narrate(context.Background(), &handbookCompleter{client: client, effort: effort}, *db, res.EIR, *chapter, *feedback)
	if err != nil {
		return err
	}
	if *out != "" {
		markdown := strings.TrimRight(narrative.Markdown, "\n") + "\n"
		if err := os.WriteFile(*out, []byte(markdown), 0o644); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "chapter %s: %s\n", *chapter, *out)
		return nil
	}
	_, err = io.WriteString(stdout, narrative.Markdown)
	return err
}

type handbookCompleter struct {
	client *llm.Client
	effort string
}

func (h *handbookCompleter) Complete(ctx context.Context, system, user string) (string, error) {
	return h.client.CompleteEffort(ctx, system, user, h.effort)
}

// specHandbookIndex builds the tag-free top-level index from rendered chapters.
func specHandbookIndex(args []string, stdout, stderr io.Writer) error {
	set := flags("spec handbook index", stderr)
	dir := set.String("dir", "", "directory with chapter-NN.md files (required)")
	out := set.String("out", "", "write the index to this path (default stdout)")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *dir == "" {
		return fmt.Errorf("spec handbook index requires --dir")
	}
	entries, err := handbook.LoadIndexEntries(*dir)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return fmt.Errorf("spec handbook index: no chapter-*.md files in %s", *dir)
	}
	index := handbook.IndexMarkdown(entries)
	if *out != "" {
		if err := os.WriteFile(*out, index, 0o644); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "handbook index: %s (%d chapters)\n", *out, len(entries))
		return nil
	}
	_, err = stdout.Write(index)
	return err
}

// specHandbookEvidence writes the evidence register appendix.
func specHandbookEvidence(args []string, stdout, stderr io.Writer) error {
	set := flags("spec handbook evidence", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path (required)")
	out := set.String("out", "", "write the register to this path (default stdout)")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if *db == "" {
		return fmt.Errorf("spec handbook evidence requires --db")
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	register, err := handbook.EvidenceRegister(*db, res.EIR)
	if err != nil {
		return err
	}
	if *out != "" {
		if err := os.WriteFile(*out, register, 0o644); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "evidence register: %s\n", *out)
		return nil
	}
	_, err = stdout.Write(register)
	return err
}

// specHandbookMint wraps rendered chapters in typed HBK semantic payloads and
// certification envelopes, then optionally writes them back in place.
func specHandbookMint(args []string, stdout, stderr io.Writer) error {
	set := flags("spec handbook mint", stderr)
	repo := set.String("repo", ".", "protocol repository root")
	assembly := set.String("assembly", defaultSpecAssembly, "assembly manifest path")
	source := set.String("source", "", "single CDL source (default: full assembly)")
	ledger := set.String("ledger", defaultSpecLedger, "identity ledger path")
	db := set.String("db", "", "SQLite database path")
	dir := set.String("dir", "preview", "clean human handbook chapter directory")
	out := set.String("out", "preview/package/handbook/sections", "output directory for the minted HBK section payloads")
	restore := set.Bool("restore", false, "strip HBK wrappers from --dir instead of minting")
	if err := set.Parse(args); err != nil {
		return err
	}
	if len(set.Args()) != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(set.Args(), " "))
	}
	if !*restore && *db == "" {
		return fmt.Errorf("spec handbook mint requires --db")
	}
	if *restore {
		restored, err := handbook.RestoreDir(*dir)
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "restored %d clean chapters in %s\n", restored, *dir)
		return nil
	}
	res, _, err := compileSpec(*repo, *assembly, *source, *ledger)
	if err != nil {
		return err
	}
	header, err := preflightHeaderValues(*db, res.EIR)
	if err != nil {
		return err
	}
	namespaces := strings.TrimSpace(header["system-namespace"])
	resolver := func(id string) (handbook.UpstreamRef, bool) {
		fieldID, ok := synthesis.FieldForID(id)
		if !ok {
			return handbook.UpstreamRef{}, false
		}
		version, err := store.VersionOf(*db, res.EIR, fieldID, id)
		if err != nil {
			return handbook.UpstreamRef{}, false
		}
		fingerprint, err := store.FingerprintOf(*db, res.EIR, fieldID, id)
		if err != nil {
			return handbook.UpstreamRef{}, false
		}
		return handbook.UpstreamRef{ID: id, Version: version, Fingerprint: fingerprint}, true
	}
	sections, err := handbook.MintSections(*dir, namespaces, resolver)
	if err != nil {
		return err
	}
	if len(sections) == 0 {
		return fmt.Errorf("handbook mint: no chapter files in %s", *dir)
	}
	if err := os.MkdirAll(*out, 0o755); err != nil {
		return err
	}
	if err := handbook.WriteSections(*out, sections); err != nil {
		return err
	}
	for _, section := range sections {
		fmt.Fprintf(stdout, "%-20s %-24s v%d %s upstream=%d\n", section.ID, section.Path, section.Version, section.Fingerprint, len(section.Upstream))
	}
	fmt.Fprintf(stdout, "handbook sections: %d (source=%s, out=%s)\n", len(sections), *dir, *out)
	return nil
}
