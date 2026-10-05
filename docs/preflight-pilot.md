# Preflight pipeline pilot (2026-10-05)

> Non-authoritative run record. The CDL sources under `protocol/` remain normative;
> this documents the first end-to-end semantic mode built on the W3 work substrate.

## Pipeline

`internal/preflight` plus `spec preflight scan|run`:

1. **Deterministic scan** (`preflight.ScanRepo`): bounded depth/file count, skipped
   directories (`.git`, `node_modules`, `vendor`, ...), language histogram, entry-point
   heuristics, and bounded evidence excerpts from `README`, `package.json`,
   `docker-compose`, `Dockerfile`, `index.js`, `routes/*`, `passport/*`. The scan is
   fingerprinted and is the only target repository access.
2. **Work units** (`preflight.Plan`): `U-PREFLIGHT-HEADER` writes the `0A` singleton;
   `U-PREFLIGHT-PERSONAS` writes the persona registry with record IDs derived from the
   payload `Prefix` field. Each unit carries routed normative sections
   (`preflight-registry`, `persona-definition`) and the bounded scan as a runtime input.
3. **Bounded execution** (`work.Session`): one work item at a time, read-set
   provenance, stale-safe commits, retry-with-feedback on rejected proposals,
   per-unit instrumentation.

## Pilot result (target: `/home/can/dev/veri-yonetimi`)

- Model: OpenAI-compatible endpoint, `qwen3.8-flash-next`.
- Committed: `PREFLIGHT-HEADER/0A` with technology stack and environments derived from
  the evidence excerpts (Node.js 20/Express 4, PostgreSQL/pg-promise, Twig/twing,
  Passport SAML, Docker/Compose), and a seven-row persona registry (web UI, API, DB
  admin, operator, SAML IdP, and related surfaces) with valid prefixes/directories.
- A previous run without normative sections produced a single catch-all persona; the
  routed sections plus evidence fixed that.
- The run is repeatable; `scripts/check.sh` covers the deterministic parts offline.

## Provider behavior found (important)

- With `reasoning_effort` at `high`/`xhigh`, long prompts intermittently return HTTP
  200 with empty `content` (reasoning consumed the output budget, `finish_reason`
  often `length`). The client now retries empty responses twice and reports
  `finish_reason`; empty-after-retry fails closed.
- Mechanical extraction units therefore declare `ReasoningEffort: "off"`
  (`--llm-reasoning` still forces an override). This made the pilot deterministic in
  practice while leaving reasoning available for synthesis-class units.
- `LEGACY_AUTOPSY_LLM_MAX_TOKENS` bounds completion tokens (default 16384).

## Diagnostics added for this pilot

- `--show-model`: prints raw model responses and review verdicts to stderr; never an
  artifact.
- Ingest validation failures report the parsed declared fields, and retries receive
  the rejection reason inline (`## Previous attempt rejected`).

## Next steps

- Entry-cluster acquisition (needs a declared cluster identifier field), acquisition
  registry, then the Discovery pipeline over source inventory/frontier.
- Per-field enum registries for registry-typed columns (`actor-class`,
  `registry-state`) so invalid tokens fail at validation rather than semantics.
- Persona deduplication/reconciliation unit with independent review.
