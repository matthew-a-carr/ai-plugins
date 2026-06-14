---
name: apply-principles
description: Ground a code change in the engineering principles BEFORE writing it. Loads the relevant tier files from this plugin, identifies which principles apply to the change at hand, and cites them by anchor (P5, C8, T3). Activate when the user is about to make a non-trivial code change, asks "what's the principle for X", asks how to approach something, or starts work where conformance to the chapter's principles matters. Not for reviewing existing diffs — that's architecture-review.
---

# Apply Principles

Load the principles, work out which apply to the change in flight, and tell the user how they shape the approach. The goal is to surface the relevant rules *before* code is written, not after.

## Step 0 — Load the tier files

Resolve the principles in this order and read everything that applies. Stop at the first source that has them.

1. `${CLAUDE_PLUGIN_ROOT}/principles/{index,values,constitution,cloud-native,tech-stack}.md` — the canonical install path when this skill runs under the engineering-principles plugin.
2. `<repo-root>/principles/{index,values,constitution,cloud-native,tech-stack}.md` — when running from a clone of `engineering-principles` directly.
3. Otherwise, search the workspace for a directory containing `principles/constitution.md` and load the set from there.

`index.md` is the routing table — one row per anchor with its kind (principle / practice / standard / process) and when it applies. Read it first; read a tier file in full only when the index points into it.

Also load `behavioural-rules.md` from the same root — Rules 1–9 always apply. If the consuming repo has its own constitution or ADR log (`CONSTITUTION.md`, `docs/decisions/`, `docs/adr/`), check it for recorded decisions that override Tier 3 defaults — the project's record is the most specific layer and wins.

If none of these resolve, **halt and report**. Do not fall back to memory — anchors and wording change.

## Step 1 — Frame the change

In one or two sentences, state what the user is trying to do. Name the surface area: API change, event change, new dep, UI change, infra change, refactor, bug fix, etc. The surface determines which principles are relevant.

## Step 2 — Pick the relevant anchors

Walk the tiers in order and pull out only the anchors that touch this change. Don't dump the whole tier — relevance is the filter. As a rough guide:

- **Always relevant** (regardless of surface): P3 (small slice), P4 (root cause), P7 (tests), P14 (mechanical enforcement), P15 (Conventional Commits), P17 (small commits), P28 (simplicity), behavioural-rules Rules 1–9.
- **API / contract change**: P11 (non-breaking), P12 (machine-readable contracts), T3 (API conventions), `patterns/non-breaking-changes.md`, `patterns/error-responses.md`, `patterns/pagination.md`.
- **Event / message change**: P10 (event-driven), P12, `patterns/event-driven-outbox.md`.
- **External call / integration**: P6 (idempotency), P22 (timeouts), P23 (retries), P24 (circuit breakers), `patterns/idempotency.md`, `patterns/timeouts-and-retries.md`.
- **Behavioural change in user-visible flow**: P9 (feature flag + rollback), P25 (graceful degradation), P26 (SLO), `patterns/feature-flag-rollback.md`.
- **New service / dep / library**: P28 (simplicity), C10 + T6 (CNCF-graduated where available), T1/T2 (right tool for the chapter), supply chain P21.
- **Domain / model code**: P2 (layering), P13 (domain modelling).
- **Infra / deployment**: cloud-native tier (C1–C…), especially C2 (twelve-factor), C8 (multi-region), C9 (auto-rollback).
- **Security surface**: P19 (Zero Trust), P20 (secrets), P21 (supply chain).
- **Strangler / legacy extraction**: T1/T2, `patterns/strangler-fig-extraction.md`.

Where two tiers say different things, Tier 3 wins (T > C > P) but flag the conflict — see behavioural-rules Rule 6.

## Step 2.5 — Resolve trade-offs

When the relevant anchors pull in different directions (e.g. P5 "ship observability in the same PR" vs P3 "keep the slice small", or P9's flag machinery vs P28 "no unjustified moving parts"), don't average them (Rule 6). Resolve in this order:

1. **The project's own record** — a constitution rule or ADR that already settles it wins. Cite it.
2. **Kind** — a standard or practice is settled; a principle is trade-off material. Standards aren't re-argued per change.
3. **Values** — argue the remaining open trade-off from `values.md` (V1–V6): which option is more mechanical (V1), smaller (V2), more boring (V3), louder on failure (V4), closer to root cause (V5), more reversible (V6)? State the call and the value it rests on.
4. **Proportionality** — practices scale with blast radius. A flag-plus-rollback plan for a copy change is overcomplication (Rule 2); the same omission on a payment path is a P9 concern.

If the resolution is a decision the project will hit again, recommend recording it as an ADR.

## Step 3 — Tell the user how it shapes the work

Output should be short and actionable. Not a lecture — a checklist they can act on.

```markdown
## Principles for this change

**Change framed as:** <one-line>

**Applicable anchors:**
- <anchor> — <one-line on what it requires here>
- <anchor> — <one-line>

**Concrete asks before code:**
- [ ] <thing the user needs to decide or design before writing>
- [ ] <thing>

**Patterns to follow:**
- `patterns/<name>.md` — <why it applies>

**Watch for:**
- <known foot-gun for this kind of change, cited by anchor>
```

If a novel decision is implied (no existing principle or ADR covers it), say so and recommend an ADR in the consuming project's repo (`docs/decisions/NNNN-<slug>.md`) — not in `engineering-principles`.

## Depth scaling

Scale the analysis to the change size (Rule 2 — simplicity first):

- **Small change** (bug fix, config tweak, single-file edit) — Frame the change, list 2–3 applicable anchors, skip the full output template. One paragraph is enough.
- **Medium change** (new endpoint, new use case, schema migration) — Full output template. Cite patterns where they apply.
- **Large change** (new service, new bounded context, infrastructure migration) — Full output template plus: recommend an ADR for novel decisions, check all three tiers, and flag any T3 defaults the project hasn't recorded a position on.

## What this skill is NOT for

- Reviewing an existing diff or PR — use `architecture-review`.
- Quoting the principles verbatim — cite by anchor, link to the file. The model and the user can read the source.
- Trivial changes (typo fix, comment tweak, dependency bump). The principles aren't free; loading four tier files for a one-line change is overcomplication (Rule 2).
