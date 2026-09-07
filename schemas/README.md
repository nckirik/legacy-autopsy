# Schemas

> Non-authoritative implementation artifacts. `protocol.md` remains the specification; derived schemas must not become a second normative source.

This directory is reserved for auditable machine projections used by parsers, validators, context assembly, fixtures, and editor tooling. No standalone schema artifact is published in M0.

The Go bootstrap has a limited structural Markdown model and protocol/routing readers. Complete typed-record boundaries, context-qualified field/table paths, enum dispatch, certification/artifact envelopes, and canonicalization remain M1/M4 work. Their implementation must route back to the applicable protocol headings—especially §§4.1.1–4.1.2 and §5.1—rather than infer rules from similarly named labels.

Future checked-in schema projections must identify protocol version and section provenance, generator/version and source inputs where generated, deterministic ordering, and boundary semantics. Generated files are regenerated, not hand-edited. Their presence never establishes protocol support or conformance.
