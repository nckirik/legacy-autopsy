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
3. **Work units** (continued): `U-PREFLIGHT-CLUSTERS` writes entry clusters with
   IDs derived from the payload `Cluster ID` field; enum-typed fields carry their
   closed option domains into the prompt.
4. **Bounded execution** (`work.Session`): one work item at a time, read-set
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

- The apparent "empty response" flakiness was our own request model: a 120s whole-
  request HTTP timeout killed long `xhigh` reasoning runs (non-streaming) or aborted
  mid-stream. Opencode, by contrast, streams with no such total cap.
- The client now streams by default (SSE) and splits timeouts: a 120s
  time-to-first-byte guard plus a 45-minute whole-stream deadline
  (`LEGACY_AUTOPSY_LLM_STREAM_TIMEOUT`), so long reasoning completes. Verified with
  `--llm-reasoning xhigh` forced on all three units: committed in ~12m24s.
- Empty-content responses are still retried twice and reported with `finish_reason`.
- `LEGACY_AUTOPSY_LLM_MAX_TOKENS` bounds completion tokens (default 16384; 32768 for
  `xhigh` runs). Mechanical units may declare `ReasoningEffort: "off"`;
  `--llm-reasoning` forces an override. `LEGACY_AUTOPSY_LLM_STREAM=0` falls back to
  non-streaming.

## Diagnostics added for this pilot

- `--show-model off|stderr|stdout`: streams raw model output (default off; diagnostic
  and never an artifact); stderr keeps stdout clean for data piping.
- With model output off, a stdout activity indicator follows the stream: on a
  terminal it redraws one line with a spinner, unit ID, reasoning/content counters,
  review state, and elapsed time, prints retry notices, and closes with a per-unit
  summary; on a non-terminal it stays quiet during the call and prints one completion
  line per unit. The same `work.Observer` events feed a future TUI.
- Ingest validation failures report the parsed declared fields, and retries receive
  the rejection reason inline (`## Previous attempt rejected`).

## Next steps

- Acquisition registry (the `0A` registry currently holds header, personas, and
  clusters; acquisition rows and `SRC` binding come with the Discovery pipeline).
- Discovery pipeline over source inventory/frontier.
- Cluster owner references and source IDs are free text at preflight and must be
  reconciled once personas and the inventory exist; a reconciliation unit is the next
  quality step.
- Per-field enum registries for registry-typed columns (`actor-class`,
  `registry-state`) so invalid tokens fail at validation rather than semantics.
- Persona deduplication/reconciliation unit with independent review.
