# Anchor index

One row per anchor: what kind of rule it is, when it applies, and what an agent
does with it. Use this to route — read the tier file for the full rule.

## Kinds

Not every anchor is a principle, and the kind changes how an agent uses it:

- **Principle** — a property the system should have, stated so you can tell
  from an artefact whether a design has it: "dependencies point inward",
  "config from the environment", "every external write is idempotent".
  Principles guide trade-offs — when two options are both possible, the
  principle says which to prefer and why. Twelve-factor is the baseline
  example: each factor is a property, not a procedure.
- **Practice** — a method you perform. Some are narrow working habits (test
  pyramid, small commits, trunk-based); some are whole design disciplines —
  DDD with its strategic design (where the context boundaries and language
  are) and tactical design (aggregates, value objects, invariants). A
  practice is *how* the work is done; what a review checks is the
  constraints the practice produces, and trade-offs are argued from the
  property those constraints protect, not from the method itself.
- **Standard** — a named default: a tool, format, or convention (Conventional
  Commits, RFC 9457, CloudEvents, Postgres). Pass/fail. Divergence needs an
  ADR in the project repo, not a debate in a PR thread.
- **Process** — a team or operations workflow (postmortems, toil tracking,
  DORA measurement). Usually outside a code diff; an agent applies it only in
  the specific scenario that triggers it.

When arbitrating a trade-off, argue from **principles** and the values
(`values.md`); cite **standards** as settled facts; hold work to the
constraints a **practice** produces.

## Tier 1 — Constitution

| Anchor | Kind | Applies when | Agent action |
| --- | --- | --- | --- |
| P1 DDD | Practice | Modelling a domain; naming types | Strategic design fixes context boundaries + language; tactical design fixes aggregates and invariants (with P13); review the constraints it produces |
| P2 Clean architecture | Principle | Any code in a layered service | Respect the dependency rule; check imports point inward |
| P3 Small slices | Practice | Planning or reviewing a change | Slice vertically; flag foundation-only PRs |
| P4 Root cause first | Principle | Any fix or mitigation | Ask why until structural; add `Remove-By:` to workarounds |
| P5 Observability first | Principle | New behaviour shipping | Land logs/metrics/traces in the same PR |
| P6 Idempotency | Principle | Any external write | Design the key + replay path before the write |
| P7 Test pyramid | Practice | Writing or reviewing tests | Unit-heavy, integration at seams, thin e2e |
| P9 Rollback + flags | Practice | Behavioural change to a live system | One-step rollback always; flag default-off on critical paths; proportional to blast radius |
| P10 Event-driven where warranted | Principle | Choosing sync vs async | Events only without a needed synchronous answer; outbox |
| P11 Non-breaking changes | Principle | Any contract change | Expand → migrate → contract; never rename+remove together |
| P12 Machine-readable contracts | Principle | Any API or event surface | Spec generated, CI fails on drift |
| P13 Domain modelling | Principle | Domain-layer code | `Result` not throw; value objects; pure functions; aggregate invariants |
| P14 Mechanical enforcement | Principle | Any rule worth keeping | Encode the rule as a test or CI gate, not a doc |
| P15 Conventional Commits | Standard | Every commit | Follow the format |
| P16 Keep a Changelog | Standard | User-facing change | Entry in the same commit |
| P17 Small commits | Practice | Every commit | One logical change; split at ~5 files / 200 lines |
| P18 Accessibility | Principle | User-facing UI | WCAG 2.1 AA; axe audit in CI |
| P19 Zero Trust | Principle | Any endpoint or service call | Authenticate everything; no "internal" exemption |
| P20 Secrets | Practice | Config, CI, credentials | Secret manager; scanners in CI; no secrets in code |
| P21 Supply chain | Practice | Dependencies, CI workflows | Pin, scan, sign; Actions pinned to SHA |
| P22 Timeouts | Principle | Every cross-process call | Explicit connect + read timeouts; no default-infinite |
| P23 Retries | Practice | Retrying external calls | Idempotent ops only; backoff + jitter; capped |
| P24 Circuit breakers / bulkheads | Principle | Multiple downstream dependencies | Isolate failure domains; one slow dependency must not starve the rest |
| P25 Graceful degradation | Principle | Critical paths with dependencies | Define fallback; document degraded modes |
| P26 SLOs | Practice | Launching a customer-facing service | Define SLO + error budget before launch |
| P27 Eliminate toil | Process | Operating a service | Automate recurring manual work; rarely diff-relevant |
| P28 Simplicity | Principle | Adding any moving part | Justify it; prefer boring (value V3) |
| P29 Blameless postmortems | Process | After an incident | Write one; track action items |
| P30 Runbooks | Practice | Adding an alert | Runbook ships in the same PR |
| P31 Trunk-based | Practice | Branching and merging | Short-lived branches; `main` deployable |
| P32 DORA metrics | Process | Measuring delivery | Track the four; not assessable on a single diff |
| P33 Code review standard | Process | Reviewing a change | Approve on net improvement; block only on real regressions |

## Tier 2 — Cloud Native

| Anchor | Kind | Applies when | Agent action |
| --- | --- | --- | --- |
| C1 Stateless apps | Principle | Any deployable service | Durable state in datastores, never in the process |
| C2 Twelve-factor | Principle | Any deployable service | Config from env, logs to stdout, identical artefact per env |
| C3 Declared health + resources | Principle | Any deployable service | Declare liveness, readiness, resource needs, disruption tolerance; on managed platforms verify the defaults (K8s realisation: T9) |
| C4 Declarative delivery | Principle | Deployment pipelines | Deployed state = Git state; no imperative prod mutation (Argo CD realisation: T10) |
| C5 Durable batch orchestration | Principle | Multi-step batch work | Orchestrator owns state; idempotent steps |
| C6 Horizontal scale | Principle | Capacity decisions | Scale out, not up; measure the bottleneck first |
| C7 Graceful shutdown | Principle | Long-running processes | Drain on SIGTERM; platform-provided on serverless |
| C8 Multi-region / sovereignty | Principle | Regulated or multi-region data | Residency per data class; region-local keys |
| C9 Auto-rollback | Principle | Deployment pipelines | Deploys self-verify and revert without a human |
| C10 CNCF default | Standard | Picking cloud-native tech | Graduated > incubating > sandbox; non-CNCF needs an ADR |

## Tier 3 — Chapter tech stack

Tier 3 entries are all **standards** — chapter defaults, not universal rules.
A consuming project whose constitution or ADRs record a different choice
follows its own record (see `AGENTS.md` read order); the agent cites the
divergence rather than applying the chapter default verbatim.

| Anchor | Applies when | Agent action |
| --- | --- | --- |
| T1 Default backend stack | Starting a new chapter service | Spring Boot / Java / Postgres unless the project records otherwise |
| T2 Strangler-fig | A legacy monolith exists | Extract, don't extend; per-context ADRs |
| T3 API design | Any HTTP API or event surface | RFC 9457 errors, cursor pagination, Idempotency-Key, URI versioning — or the project's recorded equivalents |
| T4 Integration layer | Chapter services with external consumers | Register with the layer; no direct consumer calls |
| T5 Workflow orchestration | First batch job in a service | Pick an orchestrator via ADR before shipping |
| T6 CNCF default | Picking tooling | Same rule as C10 |
| T7 AI tooling | Any repo where chapter AI agents run | Install this plugin; no drifting copies |
| T8 Identity | Auth design | Central OIDC provider — or the project's recorded auth model |
| T9 Kubernetes substrate | Chapter services on K8s | Probes, requests/limits, PodDisruptionBudgets (realises C3) |
| T10 Argo CD delivery | Chapter services on K8s | Argo CD reconciles cluster state to Git (realises C4) |

## Resolving conflicts

Most specific wins, and conflicts are surfaced, not silently absorbed
(behavioural Rule 6):

1. **The consuming project's constitution and ADRs** — authoritative for
   project-local decisions. A divergence from a tier anchor is fine if
   recorded; if it isn't recorded, recommend the ADR.
2. **Tier 3** over **Tier 2** over **Tier 1**.
3. Open trade-offs that no anchor settles are argued from the values
   (`values.md`, V1–V6) and recorded as a project ADR.
