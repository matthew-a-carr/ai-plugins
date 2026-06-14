# Pattern — Justifying a new dependency

Anchors: P28, V3, C10

Before adding a new library, service, or platform dependency, evaluate and record the justification. The size of the record matches the size of the decision.

## Evaluation criteria

Answer these before adding the dependency:

1. **What gap does it close?** What does this dependency do that the current stack cannot? If the current stack can do it with reasonable effort, prefer that (P28 — simplicity, V3 — boring technology).
2. **Is there a CNCF-graduated alternative?** For cloud-native tooling, graduated projects are the default (C10). Choosing a non-CNCF or sandbox-stage tool needs justification.
3. **What is the operational cost?** On-call burden, upgrade cadence, monitoring, learning curve. A dependency is not free after the import line.
4. **What is the lock-in cost?** How hard is it to replace this dependency later? Does it spread across the codebase (framework) or stay behind an adapter (library)?
5. **What alternatives were considered?** At least one alternative, even if it is "do nothing" or "build it ourselves".

## Short format (PR description)

For small additions (a utility library, a test helper), record the justification in the PR description. No ADR needed.

```markdown
## New dependency: Awaitility (test scope)

**Gap:** Polling assertions in async integration tests. Current approach uses
`Thread.sleep()` which is flaky and slow.
**CNCF:** N/A — test utility, not infrastructure.
**Operational cost:** Test-scope only. No runtime impact. Widely used in the
Spring ecosystem.
**Lock-in:** Low — confined to test classes, easy to remove.
**Alternative:** Hand-rolled polling loop — works but duplicates code across tests.
```

## Full ADR format

For significant additions — a new framework, a new managed service, a new data store — write a full ADR in `docs/decisions/`. See `adr-template.md` for the template.

## Threshold: when does a dependency need an ADR?

| Situation | Format |
| --- | --- |
| Test-scope library | PR description note |
| Utility library behind an adapter | PR description note |
| Framework or library that shapes the codebase structure | Full ADR |
| New managed service (database, broker, cache) | Full ADR |
| Replacing an existing dependency with an alternative | Full ADR |
| Non-CNCF choice where a graduated alternative exists | Full ADR |

The dividing line: if removing the dependency later would require changes across more than one module or layer, it warrants an ADR.

## Anti-patterns

- **Adding a dependency without mentioning it.** The dependency appears in a lockfile diff with no justification anywhere. Reviewers should flag this.
- **"Everyone uses it" as justification.** Popularity is evidence, not a reason. State what gap it closes.
- **Adding a framework to solve a library-sized problem.** If you need HTTP client retry logic, you need Resilience4j (a library), not a service mesh (a framework).

## References

- Constitution P28 (simplicity — justify every moving part)
- Values V3 (prefer boring technology)
- Cloud Native C10 (CNCF graduated as default)
- `adr-template.md` — full ADR template
