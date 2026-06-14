---
name: architecture-review
description: Review a PR diff against the engineering-principles repo. Cites anchors from constitution.md / cloud-native.md / tech-stack.md. Suggests a new ADR in the consuming project's repo if a novel decision is detected. Invoked on demand via `/architecture-review` in any session with this plugin installed.
---

# Architecture Review

Apply the engineering principles to a PR diff. Produce a structured review comment.

## Step 0 — Load the principles into context

This skill expects four reference files to be in context. Resolve them in this order and read all four before continuing:

1. `${CLAUDE_PLUGIN_ROOT}/principles/{index,values,constitution,cloud-native,tech-stack}.md` — the canonical install path when this skill runs under the engineering-principles plugin.
2. `<repo-root>/principles/{index,values,constitution,cloud-native,tech-stack}.md` — when running from a clone of `engineering-principles` directly.
3. Otherwise, search the workspace for a directory containing `principles/constitution.md` and load the set from there.

`index.md` classifies each anchor (principle / practice / standard / process) and says when it applies — use it to decide which "if relevant" checks the diff actually triggers.

`values.md` is the "why" check (V4 loud-failure, V5 root-cause, V6 reversible) — read it so concerns can be grounded in a value when no tier anchor fits cleanly.

Also useful when relevant: the `patterns/*.md` files in the same directory, cited by the checklist below.

If any of the four files cannot be found via the resolution above, **halt and report** — the skill cannot do a meaningful review without the principles loaded. Do not fall back to memory.

## Step 0.5 — Find the consuming repo's enforcement map

Before reviewing by hand, find what the repo already enforces mechanically (P14): an enforcement map in the project's constitution, architecture tests, lint/type gates, CI workflows. Two consequences:

- A checklist item the repo's harness already enforces gets the disposition **`enforced (<gate>)`** — cite the gate (e.g. "P2: enforced — `src/__tests__/architecture.test.ts`, runs in CI") instead of re-deriving the check from the diff. Review effort goes to what the harness *cannot* check.
- A relevant rule with **no** gate is itself a P14 finding: note it under Concerns and suggest the `enforce-principles` skill to add the gate.

Also read the project's ADR log: a recorded decision that diverges from a Tier 3 standard is `aligned` (cite the ADR), not a concern. An unrecorded divergence is the concern — recommend the ADR.

## Doc-only diffs

If the diff touches only documentation (markdown, ADRs, specs), run must-check items P15, P17, and P33 only, mark the rest `n/a (doc-only diff)` in a single line, and skip the "if relevant" list. Don't produce a full review skeleton for a typo fix (behavioural Rule 2).

## When to invoke

- On demand: `/architecture-review <PR url or #number>`.

## Inputs

- The PR diff and description.
- The three principles files (read in Step 0).
- Any ADRs already present in the consuming project's repo (typically under `docs/adr/` or `docs/decisions/`).

## Review checklist (cite by anchor)

The checklist has two parts. **Must-check** runs on every PR — the review output must explicitly list each item with one of `aligned`, `enforced (<gate>)`, `concern`, or `n/a (reason)`. **If relevant** runs only when the diff touches the surface the item covers. Reviews that silently omit must-check items are non-compliant.

### Severity levels

Every concern in the output uses one of three severities:

- **blocker** — Security vulnerability, correctness bug, contract break, data loss risk, or missing auth. The PR should not merge until addressed.
- **concern** — Principle violation, missing test for new behaviour, architectural drift, or unrecorded divergence from a tier standard. Should be addressed; may merge with a tracked follow-up if the author justifies.
- **suggestion** — Practice improvement, optional polish, or a pattern that would improve but isn't required. Does not block merge. Prefix with `Nit:` in the review comment.

### Must-check on every PR

1. **Conventional Commits (P15)** — Commit messages follow the convention?
2. **Small commits (P17)** — One logical change per commit? Working state preserved?
3. **Small slice (P3)** — One vertical change? Or sprawling?
4. **Root cause (P4)** — Does the description name the root cause, not just a symptom? Mitigations carry an explicit `Remove-By: <date>` line?
5. **Tests (P7)** — Unit + integration + contract where the boundary warrants? No skipped / xfailed tests added?
6. **Trunk-based (P31)** — Branch short-lived? `main` deployable? Unfinished work behind a flag (not a branch)?
7. **Code review etiquette (P33)** — The review itself: blocks only on real concerns; prefixes optional polish with `Nit:`; argues from principles not preference.

### If relevant to the diff

Run each item whose surface area the diff touches. Skip entire groups whose surface the diff does not touch. When in doubt, include the item with `n/a (reason)` rather than silently omitting it.

#### API / contract surface

- **Non-breaking (P11, `patterns/non-breaking-changes.md`)** — expand → migrate → contract? No rename+remove in same release?
- **Machine-readable contracts (P12)** — OpenAPI / AsyncAPI updated where APIs / events changed?
- **API conventions (T3, `patterns/error-responses.md`, `patterns/pagination.md`)** — RFC 9457 errors? Cursor pagination? `Idempotency-Key`? URI-versioned?
- **Changelog (P16)** — User-facing change has a `CHANGELOG.md` entry in the same commit?
- **ADR adherence** — Any change contradicts an existing ADR in the project's repo?

#### Domain / model surface

- **Layering (P2)** — Domain free of framework imports? Adapters at edges? Dependencies inward?
- **Domain modelling (P13)** — `Result` over exceptions? Value objects over primitives? Pure domain functions? Aggregates own invariants?
- **Event-driven (P10, `patterns/event-driven-outbox.md`)** — If events: outbox? CloudEvents? AsyncAPI updated?

#### Security surface

- **Zero Trust (P19)** — Auth at every endpoint? Least privilege on IAM? Service-to-service authenticated?
- **Secrets (P20)** — No secrets in diff or logs? Secret-scanner clean? Secrets from a manager?
- **Supply chain (P21)** — New deps pinned? Vulnerability scan clean? Workflow changes reviewed and SHA-pinned?

#### Reliability / resilience surface

- **Idempotency (P6, `patterns/idempotency.md`)** — Every external write idempotent? Key persisted?
- **Timeouts (P22, `patterns/timeouts-and-retries.md`)** — Every external call has connect + read timeouts? No default-infinite?
- **Retries (P23)** — Backoff + jitter? Idempotent ops only? Capped attempts and duration?
- **Circuit breakers / bulkheads (P24)** — Dependency calls wrapped? Independent pools per downstream?
- **Graceful degradation (P25)** — Critical paths have fallback? `docs/degraded-modes.md` updated?
- **SLO (P26)** — New customer-facing surface has an SLO defined? Dashboard/alert tied to it?

#### Observability / operations surface

- **Observability (P5)** — New behaviour shipped with logs/metrics/traces? Dashboard updated in same PR?
- **Runbooks (P30)** — New alert? Runbook ships in same PR.
- **Rollback + flags (P9, `patterns/feature-flag-rollback.md`)** — Behavioural change has one-step rollback? Critical-path change has flag default off?
- **Toil (P27)** — Did this PR add manual repetitive work?
- **Postmortems (P29)** — If fixing an incident: postmortem written? Action items tracked?

#### Deployment / infrastructure surface

- **Cloud Native (C1–C10)** — Twelve-factor (C2)? Multi-region considered (C8)? Auto-rollback wired (C9)? Durable batch orchestration (C5)? CNCF default (C10, T6)?
- **Mechanical enforcement (P14, `patterns/architecture-tests.md`)** — New rule covered by a test, not just docs?
- **Simplicity (P28)** — New moving parts justified in an ADR? Boring tech preferred?
- **Tech stack (T1, T2)** — Right tool for the chapter?

#### UI surface

- **Accessibility (P18)** — WCAG 2.1 AA at 375/768/1280px? Touch targets ≥44px? Axe audit green?

## Review depth

Scale the review to the PR size. A 10-line config change does not need the same skeleton as a 500-line feature PR:

- **Small (≤50 lines, single surface)** — Must-check status only. "If relevant" items addressed inline where the diff triggers them. No full skeleton.
- **Medium (50–300 lines, 1–2 surfaces)** — Must-check status + relevant surface groups only. Skip entire groups the diff doesn't touch.
- **Large (300+ lines or 3+ surfaces)** — Full skeleton with all relevant groups.

## Output

Post a single PR comment with this shape. The **Must-check status** block must list all seven must-check items explicitly — `aligned`, `enforced (<gate>)`, `concern`, or `n/a (reason)`. A review without this block is non-compliant.

```markdown
## Architecture review

### Must-check status
- P15 (Conventional Commits): aligned | concern: <…> | n/a (<reason>)
- P17 (Small commits): …
- P3 (Small slice): …
- P4 (Root cause + Remove-By): …
- P7 (Tests): …
- P31 (Trunk-based): …
- P33 (Review etiquette): …

### Aligned (relevant checks that passed)
- ✓ <anchor> — <evidence>

### Concerns
- 🔴 **blocker** <anchor> — <what's missing or wrong> — <suggested fix>
- 🟡 **concern** <anchor> — <what's missing or wrong> — <suggested fix>
- 💡 **suggestion** <anchor> — <suggestion>

### Novel decisions
- A new architectural decision detected: <one-line>. Recommend an ADR in <project repo>/docs/decisions/.

### Out of scope but noted
- <observation> — not blocking.
```

## ADR suggestions

If the PR encodes a "we decided X because Y" pattern not already covered by an ADR in the project's repo:

1. Recommend the team draft an ADR in **the consuming project's repo** (`docs/decisions/NNNN-<slug>.md`), not in `engineering-principles`.
2. Do **not** open the ADR PR automatically; humans own the decision.
3. Reference the relevant principle anchors so the ADR slots into the principles framing.

## Writeback to engineering-principles

If the review detects a pattern recurring across multiple project repos that isn't yet in the principles, the agent may push a writeback commit (or open a PR for changes that need human review) into `engineering-principles`. See `AGENTS.md` § Writeback policy for what may be pushed directly to `main` vs. what needs review.
