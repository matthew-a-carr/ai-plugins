# Tier 2 — Cloud Native

Deployment-level principles. Below the application, above the infrastructure.

## C1 — Stateless applications

**Property:** No durable state lives in the application process; sessions, caches, and queues are external.

- App pods carry no state across restarts. All durable state in datastores.
- Sessions, caches, work queues: external services with their own SLOs.
- Anti-pattern: an in-process cache treated as the source of truth. The first pod restart loses it.

## C2 — Twelve-factor (https://12factor.net)

**Property:** The application is configured from the environment, logs to stdout, and runs as an identical artefact in every stage.

Adhere to the full twelve. The ones that come up most often:

- **Config** from environment. No hardcoded URLs, no checked-in secrets.
- **Logs** to stdout; no log files on disk.
- **Backing services** as attached resources — swap a Postgres URL, don't restart the world.
- **Build / release / run** strictly separated; the artefact moved through environments is identical.
- **Disposability**: graceful shutdown on SIGTERM within 30s. PreStop hooks if needed for connection draining.
- **Dev/prod parity**: containers for local dev where feasible; the gap is measured in hours, not weeks.
- Anti-pattern: a config file checked into the repo with `dev`/`staging`/`prod` sections. Promote the same artefact through environments; let environment config in via the process environment.

## C3 — Declared health and resource needs

**Property:** Every deployable declares its liveness, readiness, resource needs, and disruption tolerance.

- Every service declares to its platform: a liveness signal (process is up), a readiness signal (can serve traffic), and its resource needs. The platform — not the process — decides restart, rotation, and placement.
- Disruption tolerance is declared, not assumed: anything serving customer traffic states how many instances may be down at once.
- On Kubernetes this means probes, resource requests/limits, and PodDisruptionBudgets — the chapter's realisation is Tech Stack T9.
- On a managed platform (serverless, PaaS), health-checking and placement are platform-provided — verify the defaults rather than replicating them.
- Anti-pattern: a service with no declared limits or health signals. The platform can't tell a hung process from a busy one, and one leaky instance degrades everything sharing its host.

## C4 — Declarative, Git-driven delivery

**Property:** Deployed state equals Git state; no imperative mutation of production.

- The deployed state of an environment is whatever the Git repo says it is. Deployment is reconciliation toward that declaration, not a sequence of imperative commands.
- No imperative mutation of production from a laptop — `kubectl apply`, console edits, or any platform equivalent. Live state and Git must not disagree.
- Rollout order is governed by declared dependencies and health checks, not by a human running steps in sequence.
- On Kubernetes the chapter realises this with Argo CD — see Tech Stack T10.
- Anti-pattern: a hand-applied "quick fix" in production. Live state and Git now disagree; the next reconciliation clobbers the fix.

## C5 — Durable batch orchestration

**Property:** Multi-step batch work is orchestrated externally with idempotent, independently retryable steps.

- Long-running, multi-step batch processes (snapshot generation, mass imports, period rollups) run on a durable orchestrator that survives node restarts. The chapter does not currently standardise on a specific orchestrator — pick per service from the available shortlist and justify in the project's ADR log.
- Each step idempotent (Constitution P6). Retries automatic with exponential backoff + jitter.
- State persisted by the orchestrator, not held in the worker process.
- Jobs parameterised by their natural partitioning keys (tenant + window).
- Anti-pattern: a cron job that calls a long-running script. State lives in the worker, retries are manual, recovery is a human at 3am.

## C6 — Horizontal scale by default

**Property:** The system scales by adding instances, not by increasing instance size; the bottleneck is measured before scaling.

- HPA on CPU + a business metric (queue depth, request rate).
- No vertical-only scaling for stateless services.
- Database scale separately: read replicas for read-heavy workloads; partitioning once a hot table crosses ~100M rows or its working set no longer fits in cache.
- Anti-pattern: scaling the app horizontally while the database is the bottleneck. Measure first.

## C7 — Graceful shutdown

**Property:** The process drains in-flight work on SIGTERM before exiting; no request is dropped during a rolling deployment.

- SIGTERM → stop accepting new requests → drain in-flight → close connections → exit.
- Drain time configured per service; defaults to 30s.
- Anti-pattern: SIGTERM exits immediately. In-flight requests 502; in-flight queue messages are lost or duplicated.

## C8 — Multi-region, data sovereignty, legal residency

**Property:** Data residency, replication topology, and key management are decided per data classification before deployment.

- Treat the **region** as a first-class deployment dimension, not an afterthought.
- Data residency follows the regulation that applies to the data subject — EU customer data stays in EU regions, UK in UK, US in US. Cross-region replication is opt-in per data class, not default.
- The service knows which region it is. Use region-aware routing at the edge (e.g. GeoDNS, GCLB) so consumers reach their local instance.
- Encryption at rest with region-local KMS keys. No cross-region key sharing for regulated data.
- Disaster recovery: RTO/RPO defined per service per region. Document which region is authoritative for which data class.
- Compliance maps (GDPR, UK DPA, US state laws, sector-specific like HIPAA / PCI) are an input to the design review, not a checklist applied at the end.
- Anti-pattern: a single global database with cross-region writes "for simplicity". Latency and regulatory cost arrive together.

## C9 — Auto-rollback on deployment

**Property:** Every deployment self-verifies and reverts automatically without human intervention on failure.

- Every deployment publishes a "deployment health" signal — derived from SLO breach, error rate, latency, dependency health.
- Argo CD / Argo Rollouts (CNCF) gates progression on the signal. A bad release auto-rollbacks before it serves significant traffic.
- Canary or progressive delivery is the default for customer-facing services; recreate is only for batch / idempotent jobs.
- Rollback is exercised in non-prod regularly — a rollback path that's never run is broken.
- Anti-pattern: a manual rollback that needs an on-call engineer awake. The bad release serves traffic for the duration of the page.

## C10 — CNCF projects are the default for cloud-native tech

**Standard:** Technology choices default to CNCF graduated projects; non-CNCF choices require an ADR. Pass/fail.

When evaluating cloud-native technology, prefer projects in the [CNCF landscape](https://landscape.cncf.io/) in this order:

- **Graduated** (Kubernetes, Argo, Prometheus, Envoy, etcd, OpenTelemetry, CloudEvents, Helm, Linkerd, …) — use without ADR.
- **Incubating** — use only when no graduated project covers the need; record the maturity risk in the project's ADR log.
- **Sandbox** — use only behind a feature flag with a written exit plan; sandbox projects churn and are sometimes archived.

Non-CNCF choices require a written rationale in the project ADR log naming the alternative considered, the gap, and the lock-in cost. Vendor lock-in, ecosystem fragmentation, and orphaned tools are the failure modes this rule exists to avoid.

Anti-pattern: adopting a sandbox or non-CNCF tool because a recent blog post recommended it. The chapter's bar is operability, multi-vendor support, and a known upgrade path — none of which a blog post supplies.
