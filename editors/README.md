# Editor tooling

Non-authoritative editor support for the CDL language. The extensions and
formatter never change semantics: `cdl fmt` is deterministic and preserves
opaque normative text, and generated artifacts are regenerated, not
reformatted.

- [`vscode/`](vscode/): TextMate grammar and language configuration for `.cdl`.
- Formatter: `go run ./cmd/cdl fmt [-w | --check] <files>` (see
  [development](../docs/development.md)).
