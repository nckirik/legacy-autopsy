# Conformance fixtures

> Non-authoritative Legacy Autopsy test-corpus guidance. Protocol §19.2 defines implementation-independent normative cases. The JSON layout and runner described here are this repository's concrete realization; neither is required to use or implement the protocol.

The implemented runner loads immutable JSON cases from:

```text
fixtures/
  positive/   # accepted bootstrap cases
  negative/   # rejected bootstrap cases
```

Each case declares a stable ID, exact protocol-heading provenance, registered capability group/operation, inputs, and a stable expected result or diagnostic class. Provenance identifies the rule being approached; it does not claim the case implements the complete cited section or a complete Part 19.2 family.

Fixtures are synthetic test data, never an evidence corpus or conformance shortcut. Do not include production exports, credentials, secrets, private correlation hints, production identifiers, personal data, or real workspace/probe content.

M0 cases in `positive/foundations.json` and `negative/foundations.json` cover protocol discovery, foundational IDs/path normalization, explicitly limited basic Markdown hashing, and workspace-skeleton structure. CI runs all registered cases and the CLI lists every complete §19.2 family as unsupported. See [docs/conformance.md](../docs/conformance.md) and [docs/roadmap.md](../docs/roadmap.md).
