# Tier 1 — Constitution

Universal architectural principles. No language. No framework. These hold for any system built in this style.

## P1 — Domain-Driven Design

**Protects:** Each domain concept has exactly one authoritative model; models do not leak across context boundaries.

A design practice with two halves; both apply. The bullets under each are the checkable constraints the practice produces — they are what a review enforces.

**Strategic design** — deciding where the model boundaries are:

- Bounded contexts are first-class; never share a model across contexts without an explicit translation.
- Ubiquitous language: code names match domain language. If product calls it an "invoice line", the type is `InvoiceLine` — not `Row`, `Record`, or `Item`.

**Tactical design** — what the model inside a context is built from:

- Aggregates own consistency boundaries; transactions stop at the aggregate edge.
- Value objects, pure domain functions, and aggregate-owned invariants — the concrete modelling rules are P13.

Anti-pattern: a single `Customer` class shared across billing, support, and onboarding contexts. Different aggregates, different invariants, different languages — share the name only inside one context.

## P2 — Clean architecture, no skipping layers

**Property:** Dependencies point inward; domain code has zero framework imports.

Following the Made Tech variant of clean architecture (https://github.com/madetech/clean-architecture) — more prescriptive than Uncle Bob's. These rules matter; the chapter applies them strictly:

- **Domain** at the centre. Pure types and business rules. Zero framework imports. ("Domain", not "Entity" — avoids framework overloading.)
- **Use Case** layer orchestrates a single verb (`SubmitClaim`, not `ClaimsService`). ("Use Case", not "Interactor" — clearer to junior developers.)
- **Gateway** is the boundary for every external dependency (DB, HTTP, broker, file system). Ports and adapters: the domain depends on the interface, the adapter implements it.
- **Presenter** formats responses. Substitutable; flexibility is a feature.
- **Delivery Mechanism** (HTTP controller, CLI, message handler) wires the use case to the outside world. Thin.
- Dependencies point inward. Domain knows nothing of HTTP or SQL. Adapters depend on domain interfaces, not the other way around.
- Constructors are for collaborators only — never side effects, never logic.
- Don't leak your internals. A gateway's signature speaks the domain's language, not the adapter's.
- Anti-pattern: a `Service` class that loads from the DB, applies business rules, and serialises to JSON. Three layers smushed into one — every change ripples across all three.

## P3 — Small slices

**Protects:** Every production change is independently reviewable, deployable, and reversible.

- Vertical slice end-to-end before widening. Database column to API field to UI control in the same PR.
- Behind a feature flag, default off. Dark-launch path verified before flip.
- Rollback ready: every change can be reverted in one commit without data corruption.
- Anti-pattern: a "foundation PR" that adds infrastructure with no caller. Either it's used end-to-end in the same PR or it's dead code waiting to be wrong.

## P4 — Root cause first

**Property:** Every fix addresses a structural cause; mitigations are explicit, time-bounded, and tracked to removal.

- Every fix asks "why" until the answer is structural, not symptomatic.
- Mitigations carry an explicit `Remove-By: <YYYY-MM-DD>` line in the PR description and in a code comment next to the workaround. Default ≤ 30 days; longer requires an ADR linking the structural fix.
- An expired `Remove-By` date is a bug — either land the structural fix or extend the date with rationale.
- Anti-pattern: "TODO: fix later" with no date or owner. Either fix now or add the `Remove-By` line.

## P5 — Observability first

**Property:** Every deployed behaviour is observable via structured logs, metrics, and traces shipped in the same change.

- Structured logs with correlation ids, business ids (customer, legal entity, invoice), and a stable schema.
- Metrics: at least latency, rate, errors per use case. Custom metrics for the domain-critical paths.
- Traces: full span from HTTP in to external call out.
- Ship the dashboard before the feature — observability lands in the same PR as the behaviour it observes, not as follow-up work.
- Anti-pattern: "we'll add metrics once it's in production". By then nobody knows what should be on the dashboard.

SLO definition is P26 — this principle covers the telemetry that lets you measure the SLO.

## P6 — Idempotency at every boundary

**Property:** Every external write accepts an idempotency key; replays return the original result without re-executing the side effect.

- Every external write accepts a client-supplied or server-derived idempotency key.
- Replays return the original result without re-doing the side effect.
- Pattern: persist `(key, result_hash)` and short-circuit on retry.
- Anti-pattern: relying on the caller to "not retry". Networks retry without asking — design for it.

## P7 — Test pyramid

**Protects:** Feedback on correctness is fast, focused, and proportional to risk.

- Unit: fast, pure, no I/O. The bulk of the suite.
- Integration: Testcontainers (real Postgres, real broker — whatever the service uses). Cover the seams.
- Contract: Pact (or equivalent) at every external boundary. Pinned in CI.
- E2E: thin layer, golden-path only.
- Anti-pattern: an inverted pyramid where E2E is the only honest signal because units mock too much. Reach into the gateways with integration tests instead.

<!-- P8 is intentionally unassigned. The numbering gap is preserved for anchor stability — renumbering existing anchors would break cross-references in consuming repos. -->

## P9 — Rollback as NFR, flags proportional to blast radius

**Protects:** Every production change is reversible within minutes; the mechanism scales with blast radius.

- Every behavioural change ships with a one-step rollback path. A feature flag with default off is the standard mechanism.
- Proportionality: the mechanism scales with blast radius. A low-risk change may rely on a clean one-commit revert (P3, P17). A behavioural change on a critical path — payments, auth, data migration, anything where a bad release is an incident — needs a flag, default off, with the rollback procedure documented in the PR.
- A consuming project may record its own threshold for "needs a flag" in an ADR. The concern a review raises is a critical-path change with *no* rollback story — not the absence of a flag on a copy tweak.
- Flags expire — orphaned flags are a smell. Quarterly cleanup pass.
- See `patterns/feature-flag-rollback.md`.
- Anti-pattern: a flag that's been on at 100% for six months. Either remove it and the old code path, or document in the flag system why it stays.

## P10 — Event-driven where the domain warrants it

**Property:** Inter-context communication uses asynchronous events only when the producer does not need a synchronous answer; all other communication is synchronous.

- Use events when the producer doesn't need to know the consumers, and the consumers don't need a synchronous answer. Otherwise use HTTP.
- Producer + consumer contracts are explicit: every event has a schema, published as an AsyncAPI document at the producer's boundary.
- Use [CloudEvents](https://cloudevents.io/) (CNCF graduated) as the envelope.
- Cross-aggregate consistency is eventual. The application layer is responsible for converging.
- Persist + emit is atomic via the outbox pattern (`patterns/event-driven-outbox.md`).
- Anti-pattern: introducing events to architect for hypothetical future consumers. Build for the consumers you have today.

## P11 — Non-breaking changes by default

**Property:** No consumer breaks during a deployment window; contract changes follow expand → migrate → contract.

- Every contract change (API, event schema, DB column) ships expand → migrate → contract. Never remove + add in the same release.
- Consumers tolerate unknown fields; producers add fields, don't rename or remove in minor versions.
- The OpenAPI / AsyncAPI document is the contract. CI breaks on a breaking diff.
- See `patterns/non-breaking-changes.md`.
- Anti-pattern: rename + remove in the same release. The deploy window is exactly when consumers break.

## P12 — Machine-readable contracts at every boundary

**Property:** Every API and event surface has a machine-readable spec (OpenAPI / AsyncAPI) that is the single source of truth, enforced by CI.

- **OpenAPI** for every synchronous HTTP API.
- **AsyncAPI** for every event surface.
- Specs are generated from source (or source-of-truth-generated from spec — pick one and enforce in CI).
- Consumers consume the spec, not screenshots of the API.
- Anti-pattern: a hand-maintained API doc that drifts from the running service. Generate one from the other; CI fails on drift.

## P13 — Domain modelling rules

**Property:** Domain logic is deterministic, composable, and independent of infrastructure.

- **No exceptions from domain logic.** Return a `Result<T, E>` (or equivalent: `Either`, sealed return types). Exceptions are for genuinely exceptional failures, not control flow.
- **Value objects over primitives.** `Money`, `Currency`, `DateRange`, `EmailAddress`, `LegalEntityId` — not raw numbers and strings. Equality is structural over all fields.
- **Units are encoded in the value object.** Money is stored as integer minor units (pence, cents) — never floats. Convert at the UI boundary only.
- **Domain functions are pure.** No I/O. No `async`. No side effects. A domain function called twice with the same inputs returns the same output.
- **Aggregates own their invariants.** Validation lives on the aggregate, not in a service that "checks" the aggregate. `Trip.allocate(amount)` enforces the budget rule itself.
- **Composition root.** Runtime dependency wiring lives in one place per service. Adapter construction is forbidden outside it; tests enforce this mechanically (see P14).
- Anti-pattern: passing `BigDecimal` around for money. Currency, precision, and rounding now live in fifty callers instead of one value object.

## P14 — Mechanical enforcement, not documentation

**Property:** Every rule the team intends to keep is enforced by the build; if the build does not check it, the rule will be broken.

If a rule matters, the build enforces it. If the build doesn't enforce it, the rule will be broken.

- **Layer import boundaries** — an architecture test asserts `domain` imports nothing from `infrastructure`, `application` doesn't import `delivery`, etc. Breaking the rule breaks the build.
- **Composition-root boundary** — only one file may construct repository adapters. Other locations attempting to are caught by a test.
- **Type safety** — strict mode on. `any` requires a written justification. Non-null assertions require a comment.
- **Spec drift** — OpenAPI / AsyncAPI specs are CI-validated against the running service.
- **Schema migrations** — backward-compatibility checks in CI for both event schemas and DB migrations.

"It's in the docs" is not enforcement. The docs describe the build; the build is the source of truth.

Anti-pattern: a `CONVENTIONS.md` that says "always do X" with no test. Within six months, half the codebase doesn't do X — and nobody notices until an incident.

## P15 — Conventional Commits

**Standard:** Every commit message follows Conventional Commits v1.0.0 format. Pass/fail.

Every commit follows [Conventional Commits v1.0.0](https://www.conventionalcommits.org/).

```text
<type>[scope]: <description>
```

Types: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`, `ci`, `perf`. Breaking changes use `!` and a `BREAKING CHANGE:` footer.

Subject line: lowercase, no trailing period, ≤ 72 characters, imperative ("add", not "added").

Why: machine-readable history enables changelog generation, semver-driven releases, and AI-assisted release notes. Without convention, agents can't reason about what changed.

Anti-pattern: a commit message of `fix stuff` or `WIP`. Forfeits changelog generation, semver inference, and `git bisect`.

## P16 — Keep a Changelog for user-facing change

**Standard:** User-facing changes have a CHANGELOG.md entry in the same commit. Pass/fail.

A `CHANGELOG.md` in [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) format. New entries under `## [Unreleased]`. Updated **in the same commit** as the behavioural change — never as a follow-up.

Rule of thumb: if a real user would notice a difference, the changelog needs an entry. `chore`, `ci`, internal `refactor`, and test-only changes don't.

Write from the user's perspective. "Improved trip budget validation" — not "Refactored TripValidator".

Anti-pattern: updating `CHANGELOG.md` in a follow-up PR after the feature ships. The entry drifts from what actually shipped.

## P17 — Small commits

**Protects:** Every commit is individually reviewable, bisectable, and revertable.

One logical change per commit. Every commit leaves the codebase in a working state — tests green, lint clean, build passing.

A commit touching ≥5 files or ≥200 lines of production code is a signal to split. Commit completed work before starting the next sub-task; don't batch unrelated change.

Why: small commits make review, `git bisect`, cherry-picking, and revert cheap. Large commits obscure intent, increase blast radius on rollback, and make CI feedback slower.

Anti-pattern: a single 50-file "refactor + feature + lint cleanup" commit. Bisect can't help; a revert takes down unrelated work with it.

## P18 — Accessibility is a non-functional requirement

**Property:** Every user-facing UI meets WCAG 2.1 Level AA; compliance is verified by automated audit in CI.

Every user-facing UI meets **WCAG 2.1 Level AA** at minimum. This is non-negotiable for any product that touches a user, not an optional polish.

- Every interactive element has an accessible name (visible label, `aria-label`, or `aria-labelledby`).
- Colour is never the sole conveyor of information.
- Minimum touch target: 44×44px.
- Colour contrast: 4.5:1 normal text, 3:1 large text and UI components.
- Mobile-first: layouts work at 375px before scaling up. Two-column forms collapse to one on mobile.
- Automated: `axe-core` (or equivalent) runs in CI at representative viewports (375 / 768 / 1280px). A failing audit blocks merge same as a failing unit test.
- Anti-pattern: "we'll add ARIA labels in a follow-up." Once a feature ships inaccessible, it stays that way.

## P19 — Zero Trust

**Property:** Every request at every service boundary is authenticated and authorised; no implicit trust based on network location.

Authenticate every request at every service boundary. No implicit trust based on network location.

- Service-to-service: mTLS or signed short-lived tokens (JWT, OAuth2 client credentials). Workload identity via [SPIFFE/SPIRE](https://spiffe.io/) where the platform supports it.
- Every endpoint enforces authorisation; there is no "internal" endpoint that skips it.
- Least privilege: every component has the minimum permissions it needs to do its job. Wide IAM grants are debt.
- Multi-tenant isolation is enforced in the auth layer, not by convention or naming.
- See Thoughtworks Tech Radar v34 — Zero Trust Architecture (Adopt).
- Anti-pattern: an internal endpoint with no auth because "it's behind the VPC". Network position is not identity.

## P20 — Secrets management

**Protects:** No credential can be extracted from source, build artefact, or log.

- No secrets in code. No `.env` in repos. No secrets in CI logs or build artefacts.
- Use a secret manager: cloud-native (Google Secret Manager, AWS Secrets Manager) or self-hosted (HashiCorp Vault, 1Password Connect).
- Prefer short-lived dynamic credentials over long-lived static secrets where the platform supports it.
- Rotation is automated and exercised — a rotation path that's never run is broken.
- Pre-commit hooks scan for accidentally committed secrets (gitleaks, trufflehog). CI re-runs the scan on every PR.
- Anti-pattern: a service that reads its DB password from an env var baked into the container image at build time. Use a secret manager and short-lived credentials.

## P21 — Supply chain security

**Protects:** Every dependency is auditable and every released artefact is traceable to source.

- Pin dependencies; commit lockfiles. Renew via Renovate / Dependabot with automated tests gating merges.
- Scan dependencies for known vulnerabilities on every PR; block on high/critical.
- Generate an SBOM (CycloneDX or SPDX) per release and publish it alongside the artefact.
- Sign released artefacts: [Sigstore](https://www.sigstore.dev/) / cosign for containers; signed tags for source releases.
- Restrict CI to trusted publishers; require code review on workflow changes; pin Actions to a SHA, not a moving tag.
- Anti-pattern: pinning a GitHub Action to `@v3`. Tags are mutable; an attacker who controls the tag controls your CI.

## P22 — Timeouts on every external call

**Property:** Every call that crosses a process boundary has an explicit timeout; no default-infinite anywhere.

Every call that crosses a process boundary has a timeout. No default-infinite anywhere.

- HTTP clients, DB drivers, message brokers, internal RPCs: all configured with explicit connect + read timeouts.
- Caller-side timeout is shorter than the upstream's worst-case response — otherwise the caller waits past the point of usefulness.
- "Slow is the new down": a hung dependency should fail fast and surface as an error, not consume the calling thread / connection pool.
- See `patterns/timeouts-and-retries.md`.
- Anti-pattern: a default-infinite HTTP client. A slow upstream silently consumes the thread pool until the service falls over.

## P23 — Retries with exponential backoff + jitter

**Protects:** Transient failures recover automatically without amplifying load on the downstream.

- Retry only **idempotent** operations (P6). Retrying a non-idempotent write multiplies bugs.
- Exponential backoff: 100ms, 200ms, 400ms, …, capped (e.g. 5s).
- Add jitter to avoid thundering herds across retrying clients.
- Cap the total retry duration. Beyond that, surface the failure.
- Retries are a tactic, not a fix — a service that needs more than 3 retries to succeed has an upstream problem.
- Anti-pattern: a retry loop with no cap on attempts or total duration. The service spends an hour failing slowly instead of failing fast.

## P24 — Circuit breakers and bulkheads

**Property:** A failing or slow dependency is isolated so it cannot exhaust resources shared with healthy dependencies.

- Wrap every dependency call in a circuit breaker. Open the circuit when downstream error rate or latency breaches a threshold.
- Half-open the circuit periodically to probe recovery. Don't flood a struggling dependency.
- Bulkhead independent dependencies: a slow auth service should not exhaust the same thread pool used to talk to the database. Separate pools / clients / queues per downstream.
- Failed-open vs failed-closed is a per-dependency decision recorded in the project ADR log.
- Anti-pattern: a shared thread pool for all downstream calls. One slow dependency starves traffic to all the healthy ones.

## P25 — Graceful degradation

**Property:** Critical paths remain useful when a non-critical dependency is unavailable; degraded modes are documented and tested.

- Critical paths have a fallback: cached result, stale-but-known-good, default response. The service stays useful when a non-critical dependency is down.
- Non-critical features fail closed — better to hide the feature than serve a broken one.
- Document degraded modes in the service repo at `docs/degraded-modes.md` — one row per dependency: "When X is down, the service does Y, customer-visible impact is Z." This is part of the service's contract, not internal trivia.
- Anti-pattern: discovering the degraded behaviour during an incident. If it isn't documented, it isn't tested, and if it isn't tested it doesn't work.

## P26 — SLOs and error budgets

**Protects:** Reliability investment is proportional to user expectation and business impact.

- Every customer-facing service has an SLO defined before launch, not after.
- SLI examples: success rate, p95 latency, freshness for batch outputs.
- Default target: 99.5% success, p95 < 5s — adjust per service and document.
- The gap between 100% and the SLO is the **error budget**. Burn it deliberately on changes; freeze changes when it's exhausted.
- Dashboards and alerts tie back to the SLO. If an alert doesn't map to an SLO breach (or imminent breach), it shouldn't page.
- See Google SRE chapter 4.
- Anti-pattern: an alert that pages without an SLO behind it. It fires forever, on-call learns to ignore it, and the real ones get lost.

## P27 — Eliminate toil

**Process:** Track time spent on manual repetitive work; automate when it exceeds ~50% of on-call time.

- Toil is manual, repetitive, automatable work that scales with service size. It is not engineering progress.
- Track time spent on toil. If it exceeds ~50% of an on-call shift, that's a signal to invest in automation.
- Convert recurring toil items into tickets, scripts, or platform features in the next sprint.
- See Google SRE chapter 5.
- Anti-pattern: a runbook step that says "run this script and paste the output into Slack" firing twice a week. Automate or remove.

## P28 — Simplicity is a design value

**Property:** Every component, dependency, and abstraction in the system is justified; operational complexity is minimised.

- Operational complexity is a tax paid on every incident, every onboarding, every change.
- Justify every added moving part: new service, new dependency, new library, new platform.
- "Boring" technology is a strength, not a weakness. Use the well-understood tool unless there's a documented reason not to.
- A system has as many failure modes as it has components — minimise components.
- See Google SRE chapter 9.
- Anti-pattern: introducing a second database "because it scales better" before the existing one is the bottleneck. Operational complexity is paid every day.

## P29 — Blameless postmortems

**Process:** Every customer-impacting incident gets a blameless postmortem with tracked action items.

- Every customer-impacting incident gets a postmortem. No exceptions for "small" ones — the lessons compound.
- Blameless means: investigate the system that allowed the human error, not the human. People act reasonably given what they know in the moment.
- Structure: timeline, root cause(s), impact, what went well, what went poorly, action items with owners and dates.
- Action items are tracked to completion in the project tracker, not buried in the postmortem doc.
- See Google SRE chapter 15.
- Anti-pattern: "human error" as the root cause. Investigate the system that allowed the error.

## P30 — Runbooks for every actionable alert

**Protects:** Every production alert is actionable by the person who receives it.

- An alert that pages a human must come with a runbook. No runbook = no page (turn it off or escalate to a ticket queue).
- Runbook structure: what fires the alert → what to check first → likely causes → mitigations → who to escalate to.
- Runbooks live with the service code (`docs/runbooks/` in the project repo), not in a separate wiki that drifts.
- New alerts ship with their runbook in the same PR.
- See Google SRE chapters 10, 11.
- Anti-pattern: a runbook step that says "check the dashboard" with no link. The on-call engineer is paged at 3am; they don't have time to guess which dashboard.

## P31 — Trunk-based development

**Protects:** The main branch is always deployable; divergence cost is paid daily, not at merge time.

- `main` is always deployable. Every commit on `main` is a candidate release.
- Branches are short-lived (hours, not days). Long-lived branches accumulate merge conflict cost and hide divergence.
- Unfinished work hides behind feature flags (P9), not behind branches.
- Branch protection on `main`: required review, required CI green, no force-pushes by humans, signed commits where the platform supports it.
- Anti-pattern: a long-lived branch that diverged so far from `main` it now needs a "big bang" merge. Pay the cost daily, not all at once.

## P32 — DORA metrics

**Process:** Measure delivery via four metrics (lead time, deploy frequency, MTTR, change failure rate). Not assessable on individual diffs.

Measure delivery the same way across services:

- **Change lead time** — time from commit to production.
- **Deployment frequency** — how often `main` lands in prod.
- **Mean time to restore (MTTR)** — outage duration from detection to recovery.
- **Change failure rate** — share of deploys that cause an incident or rollback.

The targets shift over time; the measurements stay. AI-assisted coding throughput (PR count, lines of code) is not a substitute — see Thoughtworks Tech Radar v34 caution on "Coding Throughput as Productivity Metric".

Anti-pattern: optimising lines-of-code-per-engineer with AI assistance. The vanity metric goes up; lead time and change failure rate don't move.

## P33 — Code review standard

**Process:** Approve on net improvement; block only on regressions, missing tests, security, or correctness issues.

**Approve when** the change improves the system's overall health, even if it's not perfect.
**Block when** there is a code-health regression, design flaw, missing test for new behaviour, security or correctness issue, or contract violation.

There is only "better code", not "perfect code" — adopted from Google's primary review standard.

- Don't block on: personal taste, style choices that match the codebase's conventions, optional polish.
- Prefix optional polish with `Nit:` so the author knows it's not a blocker.
- Ground arguments in technical facts and existing principles; "I prefer X" is not a review argument.
- Mentorship comments are welcome but shouldn't block approval when purely instructional.
- Anti-pattern: rubber-stamp approval. Reviewing means reading the diff and asking whether the system is healthier after it lands.
- Reference: https://google.github.io/eng-practices/review/reviewer/standard.html

## How this evolves

Agents update this file when a pattern has hardened into a principle (recurring across multiple project repos, validated by code review). Updates land directly on `main` per the writeback policy in `AGENTS.md`. The "why" goes in the commit message; the "what" goes in the file.
