# Canonical Deconstruction Language (CDL) — frozen source-language contract

Not normative relative to [`protocol.md`](../protocol.md). CDL is the source language of
the Canonical Deconstruction Protocol and becomes the sole normative protocol source
only after a parity report is accepted (see BOOTSTRAP). Until then `protocol.md` remains
authoritative, and Markdown, validator, IR, and reference views are generated
projections that are never hand-edited.

`LANGUAGE-SURFACE: FROZEN-FOR-PILOT` — this revision incorporates the pre-freeze design
decisions below and Errata 1, and freezes the surface. No further language features until
the §7.7 pilot and parity fixtures have run; pilot findings may correct semantics but not
add surface. Pre-freeze design drafts are not retained.

Design decisions incorporated before freeze:

1. MACHINE and AGENT authority are two state machines, not one flat enum; the absence
   of a human-override path for proposed MACHINE results is a deliberate contract;
2. declarations are complete in every example: `TABLE` and `TYPE` are in the surface,
   all referenced scalar types are declared, and the Light omission now cites a real
   `RULE`, not a `STEP`;
3. `LET` is explicitly step-local; durable cross-step values are declared `VALUE` with
   a lifetime, which is what cold resume needs;
4. capability closure derives from operation-level `REQUIRES`; section-level
   `REQUIRES` is a summary that must equal the union of its steps;
5. states declare `VALUES`, `INITIAL`, and optional `ALLOWS`; transitions are derived
   from `SET` effects and checked against `ALLOWS`;
6. cross-section references are qualified (`section.state`) and imported with
   `USES STATE`; `BLOCK` targets resolve to declared gates;
7. stable-ID non-reuse is checked against a committed identity ledger, not the current
   AST alone;
8. opaque normative text is explicit: `REQUIRE """..."""` clauses and defined text
   fields; the AST has no "some lines are syntax, some happen to be English" category.

ERRATA 1 (pilot-blocking consistency fixes; language surface unchanged)

E1. MACHINE authority is an enumeration derived from channel, not a transition
    system; only AGENT has a transition (`draft -> confirmed`).
E2. Computed fingerprints and similar produced values are `VALUE` declarations;
    `SET` never targets a `FIELD`.
E3. A step that performs a `SET` effect cannot be omitted from a prompt-channel
    projection; its renderer states the effect. `SUPPLIED-BY RULE` supplies a
    predicate, never a state effect.
E4. Predicate-typed states are compared explicitly (`== true`); there is no
    truthiness, and `unknown` fails closed.
E5. Rule predicates are referenced as ordinary paths (`rule.predicate`); `COMPUTE`
    is removed from the effect grammar.
E6. An assignment that does not change a value is a stutter, not a transition;
    `ALLOWS` constrains changes. Redundant assignments are still removed from
    examples.

The `E` labels are stable references cited by the execution contracts. Once the pilot
exits, these pre-freeze notes may be folded away in a later revision without changing
any normative rule.

================================================================================
0. AUTHORITY AND BOOTSTRAP
================================================================================

BOOTSTRAP
  SOURCE-OF-TRUTH protocol.md@4.1.2 UNTIL parity-report accepted-by-human
  PARITY
    machine-generated parity-report OVER hashes-of
      source protocol.md
      compiled spec
      compiler + stdlib + renderer versions
      fixture set
      generated projections
    AND existing-fixtures (24)
    AND no-silent-semantic-diff
  END
  MIGRATION section-by-section WITH dual-source drift-check
  IDENTITY-LEDGER committed AND checked by every compile
  CROWN-DSL-AS-NORMATIVE ONLY-AFTER parity-report accepted-by-human
END

RULE GENERATED-PROJECTIONS-ARE-DERIVED
  REQUIRE """
    Generated projections MUST NOT introduce normative information absent from the
    source. A generated file is never edited to change the protocol.
  """
END

RULE PROJECTION-FIDELITY
  WHEN generated-projection suspected-inaccurate
    STOP projection-use
    REPORT source-rule, projection, renderer-version
    NEVER treat projection as normative
  END
END

RULE SEMANTIC-IS-NORMATIVE
  REQUIRE """
    AGENT instructions are normative requirements whose satisfaction requires semantic
    interpretation. They are not advisory and are not weaker than MACHINE rules.
  """
END

================================================================================
1. CHANNELS AND AUTHORITY
================================================================================

Execution is described by two authority state machines plus one channel axis:

    WHO determines it?  MACHINE | AGENT
    HOW is it run?      native  | prompt

    MACHINE authority:  {proposed, authoritative}     (derived from channel; no promotion)
    AGENT authority:    draft -> confirmed            (append-only attestation)

CHANNEL native
  MACHINE-RESULT authoritative
  AGENT-RESULT draft
  DRAFT MAY become confirmed ONLY THROUGH human-confirmation
END

CHANNEL prompt
  MACHINE-RESULT proposed
  AGENT-RESULT draft
  PROPOSED-RESULT MAY-NOT close-gate
  PROPOSED-RESULT MAY-NOT attest-completion
END

Authority derivation:

    MACHINE + native  -> authoritative
    MACHINE + prompt  -> proposed
    AGENT   + any     -> draft until human-confirmation, then confirmed

Deliberate contract (CDL decision): a prompt-channel MACHINE result has no
human-promotion path (`proposed -> human-verified`) because the protocol already
forbids accepting executor-supplied deterministic values as authoritative without
independent deterministic reproduction. Manual execution therefore records `unsupported`
or `proposed`, never gate-closing state. If a future language version adds
human-verified machine results, that is a normative change and requires its own rule,
fixtures, and version bump.

Consequences encoded in every projection:

- a prompt-channel run can produce drafts and proposed machine values, and may use them
  to organize work, but cannot close a gate or attest completion;
- only native execution plus required confirmations reaches authoritative gate state;
- a confirmation is an append-only attestation bound to the immutable payload hash of
  the confirmed record; it never mutates payload identity;
- the IR records provenance for every evaluated value as a tagged union:

    value: true
    owner: MACHINE
    channel: prompt
    authority: proposed                       # machine-authority: proposed|authoritative
    rule: export-reconciliation.reconcile.guard
    renderer/stdlib versions where applicable

    record: 0H-RESUME.report
    owner: AGENT
    channel: native
    authority: confirmed                      # agent-authority: draft|confirmed
    attestation: CNF-0007 payload-hash sha256:...

================================================================================
2. DECLARATIONS AND REFERENCES
================================================================================

Declarations (global or section-local):

    ARTIFACT <id>      canonical artifact and its owner
    REGISTRY <id>      a named collection/table population
    TABLE <id>         artifact table shape: rows, row type, key, columns
    TYPE <id>          named value type (alias of a primitive or composite type)
    VALUE <id>         declared computed value with TYPE and LIFETIME
    STATE <id>         runtime or artifact state: VALUES, INITIAL, optional ALLOWS
    FIELD <id>         record/field shape and required-ness
    ENUM <id>          closed value domain
    CAPABILITY <id>    supplied validated capability (hash, table serialization, ...)
    RULE <id>          normative rule identity
    SECTION <id>       normative section identity (NUMBER remains a locator)
    STEP <id>          executable step identity
    GATES <id>, ...    named gate/completion identifiers that BLOCK may reference
    WORKFLOW-TARGET <id>, ...  workflow destinations that RETURN/LOOP may reference

References:

    USES RULE <id>          reuse of normative machinery (algorithms, profiles)
    USES REGISTRY <id>      structural dependency on a declared population
    USES STATE <id>         read/effect dependency on a declared state
    REQUIRES CAPABILITY ...  supplied validated capability needed by this step/section
    AVAILABLE(<capability>) / UNAVAILABLE(<capability>)   capability predicates

Rules:

- every referenced identifier MUST resolve; unresolved references are compile errors;
- cross-section references MUST be qualified (`section-id.local-id`) even when the
  target is imported through `USES`; unqualified names resolve only within the owning
  section. Qualified references keep autopsy traces inspectable;
- `USES` is strictly stronger than `DEPENDS`: it names reused normative machinery,
  not generic ordering. Local `STEP` order describes execution; `USES` describes
  semantic dependency;
- the `USES` graph MUST be acyclic. Workflow cycles exist only through explicit
  `LOOP`/`REOPEN` transitions, never through reuse edges;
- stable IDs are mandatory and unique forever; a retired ID is never reused. `NUMBER`
  (7.7) is a human/canonical locator and may change; IDs do not. Non-reuse is checked
  against the committed identity ledger, not the current AST alone: a compile without
  the ledger (or protocol history) cannot prove non-reuse and must fail the invariant
  rather than pass silently;
- `LET` bindings are step-local and never persist (see §4); anything that must survive
  a step boundary or a resume is a declared `VALUE` with an explicit lifetime;
- opaque normative text is explicit: `REQUIRE """..."""` clauses and the text fields
  `GOAL`, `TEXT`, and `REASON`. There is no category of source line that is "sometimes
  syntax, sometimes English"; a semantic clause is always a declared text field.

================================================================================
3. CAPABILITY CONDUCT (global source rule)
================================================================================

RULE NO-SUBSTITUTE-MACHINERY
  GOAL
    Deterministic work happens only through supplied validated capabilities.
  END
  WHEN REQUIRED-CAPABILITY UNAVAILABLE
    SET dependent-result := unsupported
    BLOCK dependent-effects
    CONTINUE independent-work
    PROHIBIT executor-authored-substitute
  END
END

RULE PROPOSED-CANNOT-CLOSE
  WHEN RESULT authority == proposed
    PROHIBIT close-gate
    PROHIBIT attest-completion
  END
END

Cascade semantics: capability closure derives from **operation-level**
`REQUIRES CAPABILITY` declarations on steps, composed with step references and `USES`
edges. A section-level `REQUIRES` is a summary only and MUST equal the union of its
steps' capability requirements; the compiler rejects drift between summary and steps.
An unavailable capability marks only its dependent results `unsupported`, blocks only
their dependent effects, and never blocks independent work. A naive section-wide
implementation that over-blocks violates this rule. Projections carry
`NO-SUBSTITUTE-MACHINERY`; the compiler verifies every capability reference has a
declared unavailable-behavior path.

================================================================================
4. EXPRESSIONS AND EFFECTS
================================================================================

Expressions are pure and total-or-error. Effects mutate declared state and are the
only operations that can cause a protocol-visible change.

    EXPR    := literal | path | call | unary | binary
    literal := integer | string | boolean | enum-value | set-literal
    path    := identifier ( "." identifier )*
    unary   := "NOT" EXPR | "-" EXPR
    binary  := EXPR op EXPR
    op      := "==" | "!=" | "<" | "<=" | ">" | ">=" | "+" | "-"
             | "AND" | "OR" | "IN"
    call    := COUNT ( selector ) | HASH ( EXPR ) | NEXT ( sequence )
             | ALL ( predicate ) | ANY ( predicate )
             | AVAILABLE ( capability ) | UNAVAILABLE ( capability )

    binding := "LET" identifier ":=" EXPR
    effect  := "SET" path ":=" EXPR
             | "APPEND" record "FOR" source
             | "REMOVE" record
             | "BLOCK" identifier ( "," identifier )*
             | "REOPEN" identifier ( "," identifier )*
             | "REQUEST" "human-input"
             | "RETURN" target
             | "LOOP" target

Reserved spellings (no grammar cleverness):

    ":="  binding/assignment
    "=="  equality
    "!="  inequality

Types: integer, string, boolean, enum, set, path, id, hash, record, registry-selector.

Value lifetime:

- `LET` bindings are **step-local**, ephemeral, never persisted, and MUST NOT be read
  from another step. A cross-step read is a compile error;
- anything that crosses a step boundary or survives resume MUST be declared:

    VALUE current-fingerprint
      TYPE hash
      LIFETIME section          # step | section | run | checkpoint
    END

- `SET` may target declared `STATE` or `VALUE` identifiers only;
- `LIFETIME checkpoint` values are durable: the resumption layer records their
  producing step and must recompute or revalidate them before use after a resume;
- a `RULE` predicate is exposed as a qualified path (`rule-id.predicate-id`) and is
  used like any other expression; there is no separate `COMPUTE` effect;
- state declarations are explicit about initialization and legal transitions:

    STATE resume-state
      VALUES unknown, fresh, stale
      INITIAL unknown
      ALLOWS
        unknown -> fresh
        unknown -> stale
        fresh   -> stale
        stale   -> fresh
      END
    END

- the transition graph is derived from `SET` effects and checked against `ALLOWS`;
  transition tables are not duplicated as normative source, and unknown/uninitialized
  state is explicit, never implicitly nullable;
- an assignment whose value is unchanged is a **stutter**, not a transition; `ALLOWS`
  constrains value changes. Examples still avoid redundant assignments.

Evaluation rules:

- `LET` is a pure binding; `SET` is the only assignment form and changes declared
  `STATE`/`VALUE` identifiers (or stutters);
- all values are typed; type errors are compile errors;
- predicate expressions require boolean operands; there is no truthiness, so a
  predicate-typed state is compared explicitly (`state == true`) and an `unknown`
  value fails closed;
- runtime evaluation errors are fail-closed: record a diagnostic, treat the result as
  unavailable, and block dependent effects — never continue with a guessed value;
- `COUNT`/`ALL`/`ANY` selectors reference declared registries and fields only;
- `HASH`/`NEXT` and any future operation are stdlib operations with one versioned
  owner; the Go runtime is tested for conformance to that stdlib, not vice versa.

================================================================================
5. STEP OWNERSHIP AND AGENT CONTRACTS
================================================================================

RULE OWNERSHIP-IS-ATOMIC
  REQUIRE """
    Exactly one of MACHINE or AGENT owns each executable STEP. Ownership is never
    inferred from verbs, wording, or projection. Changing ownership is a
    protocol/language change.
  """
END

RULE ORDER-IS-LOCAL
  REQUIRE """
    STEP order defines default execution order within a SECTION. Cross-section reuse
    is declared with USES, not inferred from document order.
  """
END

AGENT steps declare a structurally checkable contract:

    STEP <id>
    AGENT:
      EVIDENCE
        <artifact|registry|record>[.<field>]
      END
      PRODUCES
        <record>.<field>
      END
      RESULT
        <declared-type-or-enum>
      END
      REQUIRE """
        <normative semantic clause, preserved verbatim>
      """
    END

The runtime may deterministically establish: every required row received a result;
each result belongs to the declared domain; evidence bindings exist; required output
fields exist. The runtime MUST NOT claim the semantic judgment was correct; that
obligation remains AGENT-owned and human-confirmable.

================================================================================
6. EXAMPLE 1 — export reconciliation (§7.7), full
================================================================================

GLOBAL DECLARATIONS
  CAPABILITY hash.sha256
  CAPABILITY table.serialize
  REGISTRY normalized-units
  ARTIFACT approved-evidence
  RULE canonical-fingerprint
  GATES A-STRUCTURALLY-COMPLETE, R-SWEPT, Exit-A, artifact-completeness
  WORKFLOW-TARGET handoff, section-pass
END

SECTION export-reconciliation
NUMBER 7.7
TITLE "Export reconciliation"
ARTIFACT 13-EXPORT-RECONCILIATION.md

GOAL
  Every normalized export coordinate has exactly one dispositioned row.
END

USES REGISTRY normalized-units
USES RULE canonical-fingerprint
REQUIRES CAPABILITY hash.sha256, table.serialize     # summary; must equal step union

TYPE EXPORT-ID              := string
TYPE COORDINATE             := path
TYPE SRC-ID                 := id
TYPE CMP-OR-POV5-EDGE-ID    := id
TYPE EXCLUSION-DECISION     := string
TYPE CLM-ID                 := id
TYPE CLM-LIST               := set<CLM-ID>

STATE reconciliation-complete
  VALUES unknown, true, false
  INITIAL unknown
  ALLOWS
    unknown -> true
    unknown -> false
    true    -> false
    false   -> true
  END
END

STATE fingerprint-status
  VALUES unknown, computed, unsupported
  INITIAL unknown
  ALLOWS
    unknown    -> computed
    unknown    -> unsupported
    computed   -> unsupported
    unsupported -> computed
  END
END

VALUE fingerprint
  TYPE hash
  LIFETIME section
END

ENUM RECONCILIATION-STATE
  Mapped-Atomic
  Mapped-Container-Only
  Approved-Excluded
  Gap
  Superseded
END

FIELD RECONCILIATION-ROW
  export-id              -> EXPORT-ID                required
  normalized-coordinate  -> COORDINATE               required
  src-id                 -> SRC-ID                   required
  atomic-edge-id         -> CMP-OR-POV5-EDGE-ID      required
  exclusion-decision     -> EXCLUSION-DECISION       required
  reconciliation-state   -> RECONCILIATION-STATE     required
  clm                    -> CLM-LIST                 optional
END

TABLE 13-EXPORT-RECONCILIATION
  ROWS normalized-units
  ROW-TYPE RECONCILIATION-ROW
  KEY normalized-coordinate
END

RULE reconciliation-completeness
  GOAL
    Define the exact completion predicate for export reconciliation.
  END
  PREDICATE complete :=
    COUNT(13-EXPORT-RECONCILIATION WHERE reconciliation-state == Mapped-Atomic)
    + COUNT(13-EXPORT-RECONCILIATION WHERE reconciliation-state == Approved-Excluded)
    == COUNT(13-EXPORT-RECONCILIATION)
END

STEP reconcile.populate
MACHINE:
  REQUIRES CAPABILITY table.serialize
  FOR EACH unit IN normalized-units
    APPEND RECONCILIATION-ROW FOR unit
  END
END

STEP reconcile.classify
AGENT:
  EVIDENCE
    approved-evidence
  END
  PRODUCES
    13-EXPORT-RECONCILIATION.reconciliation-state
    13-EXPORT-RECONCILIATION.clm
  END
  RESULT
    RECONCILIATION-STATE
  END
  REQUIRE """
    Container mappings do not reconcile their children, and entity correlations
    never satisfy reconciliation.
  """
END

STEP reconcile.count
MACHINE:
  USES RULE reconciliation-completeness
  SET reconciliation-complete := reconciliation-completeness.complete
END

STEP reconcile.guard
MACHINE:
  GUARD reconciliation-complete == true
  OTHERWISE BLOCK A-STRUCTURALLY-COMPLETE, R-SWEPT, Exit-A
END

STEP reconcile.fingerprint
MACHINE:
  REQUIRES CAPABILITY hash.sha256
  IF AVAILABLE(hash.sha256)
    SET fingerprint := HASH(13-EXPORT-RECONCILIATION)
    SET fingerprint-status := computed
  ELSE
    SET fingerprint-status := unsupported
    BLOCK fingerprint, artifact-completeness
  END
END

PROJECTION FULL
  CHANNEL prompt
END

PROJECTION LIGHT
  CHANNELS prompt, native
  COVERS SECTION export-reconciliation

  STEP reconcile.classify
    TEXT
      "Classify every reconciliation row against approved evidence. Container mappings
       are non-closing until resolved, and entity correlations do not reconcile."
    END
  END

  STEP reconcile.guard
    TEXT
      "Incomplete reconciliation blocks structural completion and Exit A."
    END
  END

  STEP reconcile.count
    TEXT
      "Reconciliation is complete only when the mapped-atomic count plus the
       approved-excluded count equals the total number of normalized units."
    END
  END

  OMIT STEP reconcile.fingerprint
    CHANNELS native
    REASON "The native harness recomputes the artifact fingerprint from canonical bytes."
    SUPPLIED-BY native-harness
  END
END

================================================================================
7. PROJECTIONS, OMISSION, AND CHANNEL COMPATIBILITY
================================================================================

A projection is one logical view materialized per channel:

    PROJECTION FULL      CHANNELS prompt        -> full.prompt.md
    PROJECTION LIGHT     CHANNELS prompt, native -> light.prompt.md, light.native.md
    PROJECTION REFERENCE CHANNEL all            -> reference.md
    HARNESS-IR           CHANNEL native         -> ir.json
    VALIDATORS           CHANNEL native         -> validators/

Omission rules:

1. `OMIT STEP <id>` is legal only when the omitted step's entire contract is supplied
   by the named source in the projection's channel.
2. `SUPPLIED-BY` values:
     native-harness   legal only in native-channel projections; supplies runtime
                      execution of the step's effects and results;
     RULE <id>        legal when that rule is present in the same projection, defines
                      the omitted step's computation, and renders its semantics into
                      the projection text. A rule supplies a predicate, never a
                      `SET` effect; a step with `SET` effects is omissible only under
                      `native-harness` in a native channel;
     CAPABILITY <id>  legal only when the channel guarantees that capability.
3. A step that owns a gate, state transition, required artifact, evidence obligation,
   or completion claim may not be omitted for a channel unless its `SUPPLIED-BY`
   resolves in that channel. There is no implicit reliance on a runtime that the
   channel does not have.
4. `OMIT ... CHANNELS native` inside a projection that supports both channels omits
   only for native; the prompt channel still renders the step (possibly with generated
   unsupported/proposed wording).

Compiler checks: an omission without `SUPPLIED-BY`; a `SUPPLIED-BY RULE` that names a
step, a rule that does not define the omitted computation, or a step carrying `SET`
effects; a `SUPPLIED-BY native-harness` in a prompt-channel projection; an omitted
guard whose behavior is not rendered elsewhere; an override that references a deleted
step identity — all are compile errors.

Renderer:

    AST form -> versioned rendering template -> deterministic text

- every AST operation has exactly one template per projection kind;
- a new or changed AST operation does not render until its template and golden
  fixtures land in the same language/stdlib version;
- wording is never generated freely; overrides are authored source material, bound to
  a step identity, and fingerprinted like any other normative-adjacent input;
- projections carry a provenance header:

    GENERATED by cdl <version>
    language: cdl/<v>  stdlib: cdl-stdlib/<v>  renderer: <v>
    source: protocol.cdl sha256:<...>
    projection-format: <n>  channel: <native|prompt>
    generated-at: <deterministic-or-declared timestamp input>

================================================================================
8. EXAMPLE 2 — cold resume (§8.5), channel-annotated
================================================================================

SECTION cold-resume
NUMBER 8.5
TITLE "Cold resume and staleness"

GOAL
  Resume only from a checkpoint whose inputs are unchanged.
END

ARTIFACT source-snapshot
ARTIFACT 0H-RESUME
REGISTRY reverse-dependents

FIELD checkpoint
  fingerprint -> hash
END

VALUE current-fingerprint
  TYPE hash
  LIFETIME section
END

STATE resume-state
  VALUES unknown, fresh, stale
  INITIAL unknown
  ALLOWS
    unknown -> fresh
    unknown -> stale
    fresh   -> stale
    stale   -> fresh
  END
END

STEP resume.hash
MACHINE:
  SET current-fingerprint := HASH(source-snapshot)
END

STEP resume.compare
MACHINE:
  IF current-fingerprint != checkpoint.fingerprint
    SET resume-state := stale
    REOPEN reverse-dependents(checkpoint.fingerprint)
  ELSE
    SET resume-state := fresh
  END
END

STEP resume.report
AGENT:
  EVIDENCE
    source-snapshot
    checkpoint.fingerprint
  END
  PRODUCES
    0H-RESUME.report
  END
  REQUIRE """
    Report which loaded inputs changed before proposing any mutation.
  """
END

PROJECTION LIGHT
  CHANNEL prompt
  COVERS SECTION cold-resume

  STEP resume.compare
    TEXT
      "If the snapshot changed since the checkpoint, treat the checkpoint as stale
       and reopen its reverse dependents."
    END
  END
END

Notes:

- In prompt channel, `current-fingerprint` and the `resume-state` transition are
  proposed machine values; they may inform the agent's report but cannot close a
  stale/fresh decision on their own.
- In native channel, `resume.compare` is authoritative; the agent report remains a
  draft until confirmation.
- `current-fingerprint` is a declared `VALUE` with `LIFETIME section` precisely
  because it crosses a step boundary; the resumption layer therefore knows it must be
  recomputed or revalidated after a cold resume.

================================================================================
9. EXAMPLE 3 — workflow control (§10.6)
================================================================================

SECTION iteration
NUMBER 10.6
TITLE "Iteration and pass control"

GOAL
  Advance the canonical iteration token until closure or exhaustion.
END

WORKFLOW-TARGET section-pass, handoff
GATES iteration-advancement

USES STATE export-reconciliation.reconciliation-complete

STATE synthesis-gaps-empty
  VALUES unknown, true, false
  INITIAL unknown
END

STATE iteration-token
  VALUES ALFA..ZULU, exhausted
  INITIAL ALFA
END

STATE iteration-state
  VALUES running, exhausted, handoff
  INITIAL running
  ALLOWS
    running   -> exhausted
    running   -> handoff
    exhausted -> handoff
  END
END

STEP iteration.check
MACHINE:
  IF export-reconciliation.reconciliation-complete == true AND synthesis-gaps-empty == true
    SET iteration-state := handoff
    RETURN handoff
  END
END

STEP iteration.advance
MACHINE:
  LET next-token := NEXT(iteration-token)
  IF next-token == exhausted
    SET iteration-state := exhausted
    BLOCK iteration-advancement
    REQUEST human-input
  ELSE
    SET iteration-token := next-token
    LOOP section-pass
  END
END

Notes:

- Cross-section state reads are qualified and imported (`USES STATE
  export-reconciliation.reconciliation-complete`), so an autopsy trace can show
  exactly which foreign state a step depends on.
- `BLOCK iteration-advancement` resolves to a declared gate; an undeclared block
  target is a compile error.
- `RETURN`, `LOOP`, `REOPEN`, `BLOCK`, `REQUEST` are workflow semantics owned by
  MACHINE. In prompt channel they render as explicit instructions, but their results
  are proposed and cannot attest completion.
- Workflow primitives are planned to live in a versioned standard module rather than
  the frozen core grammar; the section-level surface does not change when they move.

================================================================================
10. GENERATED VIEWS (from Example 1, steps reconcile.count/guard)
================================================================================

full.prompt.md, template-rendered from `RULE reconciliation-completeness` plus the
guard step, prompt channel:

    Reconciliation is complete exactly when the count of `Mapped-Atomic` rows plus
    the count of `Approved-Excluded` rows equals the total number of normalized
    units. Less than 100% reconciliation blocks `[A-STRUCTURALLY-COMPLETE]`,
    `[R-SWEPT]`, and `Exit-A`.

validators (native channel), generated from the same predicate rule:

    mapped   := count(rows, "Mapped-Atomic")
    excluded := count(rows, "Approved-Excluded")
    complete := mapped+excluded == count(rows)
    if !complete {
        block("A-STRUCTURALLY-COMPLETE", "R-SWEPT", "Exit-A")
    }

light.prompt.md, prompt channel, overrides + omission applied:

    Reconcile every normalized unit. Classify every reconciliation row against
    approved evidence. Container mappings are non-closing until resolved, and entity
    correlations do not reconcile. Reconciliation is complete only when the
    mapped-atomic count plus the approved-excluded count equals the total number of
    normalized units. Incomplete reconciliation blocks structural completion and
    Exit A. Fingerprints and other deterministic values are proposed; unavailable
    capabilities block dependent gates.

harness IR, native channel (excerpt):

    {section: "export-reconciliation",
     step: "reconcile.guard",
     owner: "MACHINE", channel: "native",
     authority: {machine: "authoritative"},
     guard: "reconciliation-complete",
     block: ["A-STRUCTURALLY-COMPLETE", "R-SWEPT", "Exit-A"]}

================================================================================
11. CONFORMANCE INVARIANTS (compiler/runtime)
================================================================================

1. Every executable STEP has exactly one owner (`MACHINE`/`AGENT`).
2. Every referenced section, step, rule, artifact, table, type, value, state, field,
   enum, registry, gate, workflow target, and capability resolves; unresolved
   references are compile errors. Cross-section references are qualified.
3. MACHINE expressions use only declared constructs and stdlib operations; all values
   are typed; predicate operands are boolean (no truthiness, explicit `== true`);
   runtime evaluation errors are fail-closed.
4. `LET` bindings never cross a step boundary; cross-step and resumed values are
   declared `VALUE` with a lifetime; `SET` targets `STATE`/`VALUE` only.
5. States declare `VALUES`, `INITIAL`, and `ALLOWS`; the transition graph derived from
   `SET` effects conforms to `ALLOWS` (identity assignments are stutters, not
   transitions), and uninitialized state is explicit.
6. Capability closure derives from operation-level `REQUIRES`; a section-level summary
   equals the union of its steps, and unavailability never over-blocks independent work.
7. Generated projections cannot add normative requirements absent from the source.
8. Projection overrides and omissions point to existing identities and satisfy the
   channel-compatibility rules; a `SUPPLIED-BY RULE` must name a rule that actually
   defines the omitted computation and never a `SET` effect; gate/state/artifact/
   evidence/completion behavior is never silently dropped.
9. Stable IDs are unique and never reused; the check runs against the committed
   identity ledger, and a compile without ledger/history fails the invariant rather
   than passing silently. `NUMBER` is a locator, not an identity.
10. Harness execution and generated validators are two implementations of the same
    stdlib and language semantics; neither redefines them.
11. Every AST operation has a reviewed template and golden fixture before rendering.
12. Workflow control is machine semantics, never hidden in prose.
13. Unavailable capabilities produce the declared unsupported/blocked behavior; no
    executor-authored substitute is ever accepted as authoritative.
14. MACHINE and AGENT authority are distinct state machines; prompt-channel MACHINE
    results are proposed and cannot close a gate or attest completion.
15. AGENT contracts are shape-checkable (evidence, produced fields, result domain)
    without any attempt to validate semantic truth.
16. Opaque normative text appears only in declared text fields (`GOAL`, `TEXT`,
    `REASON`, `REQUIRE """..."""`); the AST has no mixed syntax/prose line category.

================================================================================
12. SUBSTRATE COVERAGE PROFILE (next step after the CDL freeze)
================================================================================

Classify `protocol.md`@4.1.2 with the CDL vocabulary and report:

        rule count
        token/line share
        gate ownership
        state-transition ownership
        artifact lifecycle ownership
        cross-rule reuse (USES density)
        capability-dependent rules
        AGENT rules with structurally checkable contracts
        pure free-semantic prose

This is a profile, not a ratio gate. A textually small MACHINE layer that owns gates,
transitions, identities, and completion conditions is the skeleton of the protocol and
justifies an executable substrate. The profile decides the shape of the toolchain, not
whether the substrate exists.

================================================================================
13. RESOLVED DECISIONS (language surface frozen)
================================================================================

1. Transitions are derived from `SET` effects and checked against optional `ALLOWS`;
   transition tables are not duplicated as normative source.
2. One logical projection per view, materialized per channel
   (`light.prompt.md`, `light.native.md`).
3. Workflow primitives move to a versioned standard module after stabilization; the
   frozen core grammar keeps only their syntax.
4. `CAPABILITY` is first-class; explicit profile/version metadata is deferred but the
   declaration syntax already permits it.
5. Confirmations are append-only attestations bound to immutable payload hashes; they
   never mutate payload identity.
6. The parity report is machine-generated and human-accepted over hashes of source,
   compiler, stdlib, renderer, fixtures, and generated outputs.

================================================================================
14. PILOT SCOPE — §7.7 END-TO-END (next concrete work)
================================================================================

Deliverables:

1. `RULE reconciliation-completeness` and `SECTION export-reconciliation` as valid
   frozen-grammar source;
2. compiler checks for references, capability closure, state `ALLOWS`, omission
   validity, and `LET` lifetime;
3. renderers for the AST forms used by this section, with golden fixtures;
4. one prompt projection and one native validator projection;
5. parity fixtures: the corrected §7.7 behavior (positive exact-line, negative old
   inversion) plus preservation of the existing 24 fixtures;
6. dual-source drift check between `protocol.md@4.1.2` and the compiled section.

Golden negative fixtures already identified by the errata (the first compiler must
reject each):

1. `OMIT STEP reconcile.count` in a prompt channel — omitted step carries a `SET`
   effect and `SUPPLIED-BY RULE` cannot supply it;
2. `SET fingerprint := ...` where `fingerprint` is a `FIELD` — assignment targets
   must be `STATE`/`VALUE`;
3. `GUARD reconciliation-complete` without `== true` — no truthiness, `unknown`
   must fail closed;
4. `SET reconciliation-complete := COMPUTE reconciliation-completeness` — `COMPUTE`
   is not in the grammar; the rule predicate is a qualified path;
5. `BLOCK undeclared-target` — block targets must resolve to declared gates;
6. `SET iteration-state := running` from `exhausted` — a value change not listed in
   `ALLOWS` (the example avoids the stutter case, so a true violation is testable).

Exit criteria: all compiler checks pass; generated outputs reproduce byte-for-byte; the
parity report shows no silent semantic difference; `protocol.md` remains authoritative.

Stop conditions: any reference that cannot resolve without inventing surface; any
parity mismatch attributable to language semantics rather than a protocol defect; any
temptation to expand the grammar before the pilot exits. Record such findings as
protocol or language issues; do not fix them by growing the frozen surface.
