# Export acquisition helper

> Non-authoritative implementation boundary derived from Protocol Part 7. This directory contains no helper implementation and grants no agent production capability.

**Current state:** M0 boundary documentation only. Implementation is planned for [M8](../../docs/roadmap.md#m8--export-acquisition-helper); the detailed [architecture proposal](architecture.md) is not an implemented CLI contract.

The planned helper is an independently runnable, operator-controlled local tool for acquiring and sanitizing Appsmith and n8n exports. It prepares approved evidence projections for Legacy Autopsy; it is not invoked by an agent with production credentials and does not make raw material repository-visible.

## Trust split

Operator-only state:

- raw Appsmith exports;
- n8n credentials and API access;
- raw workflow data;
- private matching hints;
- private entity-correlation indexes.

Agent-visible, only after operator approval:

- sanitized projections;
- raw/projection fingerprints and lineage reports that disclose no secret;
- normalized virtual maps;
- acquisition manifests and reconciliation candidates;
- operator-action queues.

Correlations and normalized maps are navigation aids, not behavioral truth. Bridge fragments are proposed patches requiring authorized review before merge.

## Future safety contract

The helper must:

- default to dry-run;
- accept credentials only through environment variables, standard input, or operating-system credential storage—never command-line arguments;
- keep raw exports, temporary sensitive data, and private matching state outside the repository with restrictive permissions and best-effort cleanup;
- perform field-aware sanitization and secret scanning rather than blind text replacement;
- fail closed on unsupported schemas, unknown sensitive fields, unresolved high-confidence secret findings, missing lineage, mixed snapshots, or rejected operator review;
- use approved read-only interfaces and make no Appsmith, n8n, or production mutation;
- transmit no exports, secrets, private hints, repository content, or production data to third parties;
- preserve stable platform IDs, fingerprints, snapshot boundaries, virtual coordinates, and sanitization lineage;
- publish only operator-approved sanitized output using atomic staging and rename;
- keep operator actions explicit when references cannot be safely acquired or resolved.

## Anticipated boundary

A future local CLI should separate `inspect`, `sanitize`, `review`, and `publish` phases. Dry-run inspection produces private operator output only. Publication requires explicit approval, writes a new immutable projection and manifest, and never overwrites raw evidence. Exact flags and schemas will be designed with the implementation; this README is not a CLI specification.

No production-connected acquisition, sample raw export, credential example, sanitizer, or publisher exists today. Protocol Part 7 remains deferred.
