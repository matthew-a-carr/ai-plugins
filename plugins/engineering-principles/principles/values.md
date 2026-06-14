# Values

Six statements describing **why** the principles exist. The principles answer *what* and *how*; this file answers *why*. When a principle and a value disagree, the value wins — but if that happens, the principle is wrong and should be updated.

Each value names the principles it informs. Anchors are stable; consult the linked principle for the concrete rule.

## V1 — Mechanical over documentary

A rule that lives only in documentation will be broken within six months. A rule encoded in a test, type checker, lint rule, or CI gate is broken at the moment of breaking and can be fixed before merge. The harness — tests, types, CI, architecture rules — is what actually shapes behaviour. Documentation describes the harness; the harness is the source of truth.

Informs: P14 (mechanical enforcement), P15 (Conventional Commits), P21 (supply chain), C9 (auto-rollback).

## V2 — Ship small slices

Vertical changes — database column to API field to UI control in one PR — surface integration risk early and roll back cleanly. Big-bang changes hide divergence and turn rollbacks into incidents. Branches stay short-lived; flags hide unfinished work; commits stay small.

Informs: P3 (small slices), P9 (feature flags), P11 (non-breaking changes), P17 (small commits), P31 (trunk-based).

## V3 — Boring over novel

Operational complexity is paid every day — in incidents, onboarding, upgrades, and on-call sleep. The default tool is the well-understood one. Novelty requires a written justification in an ADR naming the alternative considered, the gap it closes, and the lock-in cost. "This is interesting" is not a justification.

Informs: P28 (simplicity), C10 (CNCF default), T1 (default backend stack).

## V4 — Loud over silent failure

A test that's silently skipped, a retry that hides a deeper bug, an error swallowed because "it usually works" — these cost days when they go wrong. Loud failure costs minutes. The repo defaults to surfacing uncertainty: red tests, structured errors, explicit timeouts, alerts tied to SLOs, postmortems for every customer-impacting incident.

Informs: `behavioural-rules.md` Rule 9, P5 (observability), P22 (timeouts), P26 (SLOs), P29 (postmortems).

## V5 — Root cause over symptom

Every fix asks "why" until the answer is structural. Mitigations are explicit, dated with a `Remove-By` line, and removed when the underlying fix lands. A repeated symptom is a structural problem the chapter hasn't named yet.

Informs: P4 (root cause first), P29 (blameless postmortems).

## V6 — Reversible over heroic

A change that can be rolled back in one step is safer than a change that can't, even if it ships slower. Feature flags, expand–migrate–contract, auto-rollback on deployment, and trunk-based development are all expressions of the same value: prefer reversibility to recovery.

Informs: P9 (feature flags), P11 (non-breaking changes), C9 (auto-rollback), P31 (trunk-based development).
