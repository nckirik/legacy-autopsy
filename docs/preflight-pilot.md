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

## Discovery inventory (same target)

`spec discovery run --target veri-yonetimi --max-files 40` enumerated candidates
deterministically (paths + streamed SHA-256), planned two inventory units, and the
model authored 40 `SOURCE-INVENTORY-RECORD` rows with runtime-computed
`SRC-FILE-<hash>` IDs (kind + normalized coordinate in the canonical key). Known gap:
frontier records and coverage rows are not yet generated; `coverage-summary` is
`C-GAP` until swept.

## Assurance claims (same target)

`spec assurance claims --batch-size 20 --max-excerpt 800` read the 40 committed
source units, materialized bounded file excerpts, and the model authored 40
`CLAIM` records with runtime IDs (`CLM|namespace|subject|claim`) and evidence-bound
anchors. Cross-record validation rejects any claim whose subject is not a committed
source unit, whose source fingerprint does not match the inventory row, or whose
anchor does not start at the unit coordinate.

## Coverage rows (same target)

`spec assurance coverage` mechanically expanded the immutable
`src-version | persona | track | environment | snapshot` tuple into 40 idempotent
`SOURCE-COVERAGE-ROW` records (fail-inclusive `Unmapped`, `In-Scope`, claims bound by
subject) and printed the arithmetic: `sources[FILE]: 40`, `coverage[Unmapped]: 40`.

## Terminal frontiers (same target)

`spec discovery frontier` mechanically recorded one terminal `FRONTIER-RECORD` per
source unit (state `Traversed`, discovery-method `other`, runtime `FRT` IDs) bound to
the real invocation `ALFA-NONE-NONE-001`; the runtime also appended its first `0G`
entry to `INVOCATION-LOG`.

## Required traversal matrix (same target)

After the v4.3 revision, `spec preflight matrix` expanded 25 keyed
`REQUIRED-TRAVERSAL-MATRIX-ROW` cells (1 cluster × 5 personas × 5 POVs,
fail-inclusive `Required`/`[R-UNSWEPT]`), and `spec assurance gate` now evaluates
real closure (`TRAVERSAL-MATRIX-CLOSURE Fail 0/25`).

## Synthesis (same target)

After v4.3, `spec synthesis run --catalog MOD` authored **16 `ARCHITECTURE-MODULE`
records** from the 40 claims: runtime IDs from `domain-path`+`record-title`, injected
`block-id`, responsibility boundaries grounded in claim text. Known blocker:
synthesis-to-source mappings are not declared (PI-18), so coverage cannot reach
`Mapped` yet.

## Tracing bindings (v4.4, same target)

`spec synthesis run --catalog MOD` re-authored 5 modules declaring `source-bindings`;
`spec synthesis bindings` materialized 41 keyed `DERIVED-FROM-EDGE` rows
(claim bindings resolved to their subject source units, invocation
`ALFA-NONE-NONE-006`), and `spec assurance advance` recomputed dispositions:
**40/40 coverage rows `Mapped`**. The gate now reports `COVERAGE-DENOMINATOR Pass`.
**Exit A is now green on the pilot**: v4.5 keyed sweep records and per-cell sweeps
closed all 25 required cells, and the synthesis-coverage condition passes because all
40 mapped source units carry `DERIVED-FROM-EDGE` synthesis provenance:

```
ACTIVE-TICKETS               Pass
COVERAGE-DENOMINATOR         Pass  40/40 mapped or excluded
COVERAGE-SOURCE-INVENTORY    Pass  40 source units, 40 coverage rows
FRONTIER-CLOSURE             Pass  40 frontiers, none open
SYNTHESIS-COVERAGE           Pass  40/40 mapped units have synthesis provenance
TRAVERSAL-MATRIX-CLOSURE     Pass  25/25 matrix cells closed
Exit A: Pass
```

## Identity and update policy (all vague entities)

Synthesis records and claims anchor identity in evidence (first bound source unit;
claim subject + source anchor). Verified live: re-running identical evidence reproduced
8/8 `MOD-*` IDs; byte-identical records are idempotent no-ops; changed bodies at the
same identity fail closed by default and are replaced at `version+1` with
`--update-existing`.

## Fresh rehearsal and re-run semantics (2026-10-06)

A from-scratch workspace (12-file window) reached **Exit A Pass** with all six
conditions green (120/120 matrix cells, 12/12 mapped and provenanced). Re-running the
model stages on the same evidence without flags reproduced every identity
(`SRC-FILE-*`, `CLM-*`, `MOD-*`) and refused changed content as
"explicit update required". With `--update-existing`, discovery, claims, and synthesis
committed versioned replacements, and the gate remained **Pass** - identities, edges,
dispositions, and cell closure all survived the updates.

## Next steps

- Remaining catalog runs; export reconciliation; symbol/route/AST kinds.
