# Tier 3 — Chapter Tech Stack

Concrete tool choices for systems built or extended in this chapter. Less stable than Tier 1 — revisit when the chapter adopts a new default. Fork this file and fill it in per chapter; the entries below are a worked example.

Every entry here is a **standard** (a chapter default), not a universal principle. When this plugin is installed in a project whose constitution or ADRs record a different choice — a different language, framework, error shape, or identity model — the project's record is authoritative. Cite the divergence ("project uses X per ADR NNN; T1 default not applied") instead of applying the chapter default verbatim. A divergence with no recorded decision behind it is a conflict to surface: recommend an ADR in the project repo. See `AGENTS.md` § Read order.

## T1 — Default backend stack

**Standard:** New chapter services use Spring Boot / Java / Postgres unless the project records otherwise.

- **Spring Boot / Java** is the default for new BE services.
- **Kotlin** allowed where the team is already fluent; not the default — readability and onboarding cost outweigh ergonomics for the chapter at current scale.
- **Postgres** for relational state in new services. Other RDBMS only when reading from an existing legacy system.
- A single central OIDC identity provider for the chapter; see T8.
- Async messaging and event sourcing: no chapter-wide default. When picking for a new service: pick per use case, and the project **must** record the choice and the rejected alternatives in `docs/adr/` before merging the first message handler.
- Anti-pattern: introducing a second backend language because one team prefers it. New languages need a chapter-level decision, not a project-level one.

## T2 — Strangler-fig for legacy systems

**Standard:** New bounded contexts are extracted from the monolith as separate services; the monolith is not extended.

- Where a legacy monolith exists, decomposition follows strangler-fig (see `patterns/strangler-fig-extraction.md`).
- Specific extractions are decided per-context, in the consuming project's own ADR log — not pre-committed here.
- Anti-pattern: extending the monolith with a new HTTP surface "as a stepping stone". That's reinforcement, not extraction.

## T3 — API design and specs

**Standard:** HTTP APIs use RFC 9457 errors, cursor pagination, Idempotency-Key headers, and URI versioning.

- **OpenAPI** for every sync HTTP API; **AsyncAPI** for every event surface. Spec is the contract; consumers consume the spec.
- Spring: `springdoc-openapi` generates OpenAPI from controllers. AsyncAPI specs are hand-curated or generated using one of the [AsyncAPI tools](https://www.asyncapi.com/tools) (the Spring Cloud Stream and code-generator templates are the common ones); commit the spec alongside the code.
- **Two-level URL scoping** for tenanted APIs: both customer-scoped (`/customers/{cid}/<resource>`) and tenant-scoped (`/customers/{cid}/tenants/{tid}/<resource>`) where the tenancy boundary matters. Same response shape. Tenant id is first-class — always a field on the row even when not in the URL.
- Path segments are reserved for **genuinely-invariant** primitives (customer, tenant). Churny dimensions (period, status, filters) go in the query string.
- **Canonical + adapter, not consumer-shaped.** APIs return the chapter's canonical shape; downstream consumers own transformation.
- **Snapshot, not delta, unless the consumer can guarantee replay.** Default to per-window snapshots; offer delta endpoints only when justified.
- Events use [CloudEvents](https://cloudevents.io/) (CNCF graduated) as the envelope.
- **Error responses follow [RFC 9457 Problem Details](https://www.rfc-editor.org/rfc/rfc9457.html)** — `type`, `title`, `status`, `detail`, `instance`, plus extension fields for the domain. One error shape across all services; consumers parse once.
- **Pagination is cursor-based** for collections that can grow unbounded (use `?cursor=<opaque>&limit=<n>`; the response includes `next_cursor`). Offset pagination is allowed only for small, naturally-bounded collections.
- **Idempotency-Key header** for unsafe (non-GET) operations. Clients supply a key; the server short-circuits on retry per Constitution P6.
- **URI versioning** (`/v2/...`) for sync APIs. Major bumps reserved for genuinely incompatible shape changes (P11).
- Anti-pattern: `200 OK` with `{ "error": "..." }` in the body. Use the right status code; consumers parse one shape (the Problem Details envelope).

## T4 — Integration layer

**Standard:** External consumers connect through the integration layer; no direct service-to-consumer calls.

- A single routing/translation layer sits between external consumers and downstream services. New services register with it; consumers never call them directly.
- Anti-pattern: a new service called directly by an external consumer, bypassing the integration layer. Routing now lives in N consumers instead of one place.

## T5 — Workflow orchestration

**Standard:** The first batch job in a service requires an orchestrator choice recorded in an ADR.

- Durable orchestration for batch processes is tenant-scoped, window-parameterised (see Cloud Native C5).
- The chapter has not standardised on a single orchestrator. The project **must** pick one before the first batch job ships and record the choice + rejected alternatives in `docs/adr/`. Subsequent batch jobs in the same service use the same orchestrator unless a follow-up ADR overrides.
- Anti-pattern: cron + ad-hoc scripts for multi-step batch work. State lives in the worker, retries are manual, and recovery is a pager rotation problem.

## T6 — CNCF default for cloud-native tooling

**Standard:** Same rule as C10: CNCF graduated > incubating > sandbox; non-CNCF needs an ADR.

Same rule as Cloud Native C10 — graduated first, then incubating, then sandbox; non-CNCF choices need a written rationale in the project ADR log. The current chapter stack: Kubernetes, Argo CD, Helm, Prometheus, OpenTelemetry, Envoy (via service mesh where applicable), CloudEvents.

## T7 — AI tooling

**Standard:** Every repo where chapter AI agents run installs this plugin; no drifting copies.

- This `engineering-principles` plugin is the chapter reference for AI work. Install it in any repo where chapter AI tooling runs so the rules apply consistently.
- Anti-pattern: each repo carrying its own copy of the principles, drifting independently. Install the plugin; update via writeback.

## T8 — Identity

**Standard:** Authentication uses a central OIDC provider; tenant boundaries are enforced in the IdP.

- One central OIDC identity provider for both customer-facing auth and service-to-service.
- Multi-tenant boundaries are enforced in the identity provider's own tenant model, not re-implemented per service.
- Anti-pattern: a service rolling its own JWT verification with hardcoded signing secrets. Federate through the central identity provider.

## T9 — Kubernetes as the chapter substrate

**Standard:** Kubernetes deployments declare probes, resource requests/limits, and PodDisruptionBudgets.

The chapter's realisation of Cloud Native C3 (declared health and resource needs):

- Chapter services run in Kubernetes. Helm chart or Argo CD application manifest in the repo.
- Health probes: `liveness` (process up — restart on fail), `readiness` (can serve traffic — pull from rotation on fail), `startup` (defers liveness for slow-starting apps).
- Resource requests and limits set. Never deploy without them.
- PodDisruptionBudget for anything serving customer traffic.
- Anti-pattern: a Deployment with no resource limits. One leaky pod evicts the rest of the node.

## T10 — Argo CD for delivery

**Standard:** Argo CD reconciles cluster state to Git; no imperative kubectl apply in production.

The chapter's realisation of Cloud Native C4 (declarative, Git-driven delivery):

- Argo CD reconciles cluster state to Git. Sync waves and health checks govern rollout order.
- No `kubectl apply` from a laptop in production.
- Anti-pattern: `kubectl apply` to fix a "quick thing". Cluster state and Git now disagree; the next sync clobbers the fix.
