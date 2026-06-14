---
name: debugging-and-error-recovery
description: >
  Systematically diagnose and fix a bug, failing test, flaky CI job, or
  production error using a reproduce → localize → reduce → fix → guard loop.
  Use when something is broken and the cause is not obvious, when a test or CI
  job fails intermittently, when triaging a Sentry/production error, or when a
  first fix attempt did not work. Stops guess-and-check; forces evidence before
  edits.
---

# Debugging & Error Recovery

Pairs with `browser-testing-with-devtools` (repro in a browser) and `../tdd`
(the guard test). The rule: **no fix until you can reproduce and explain.**

## The loop

### 1. Reproduce
- Get a deterministic repro before touching code. Note exact steps, inputs,
  environment, and the **verbatim** error / stack trace.
- Flaky? Run it in a loop to measure the failure rate — that number is your
  signal that the fix worked.
- Can't reproduce locally? Reproduce the *environment* (same Node version, same
  data, Testcontainers DB, prod-like build) before assuming it's "just CI".

### 2. Localize
- Read the stack trace top-down to the first frame in *your* code.
- Bisect: `git bisect` for a regression; comment/binary-search the data path for
  logic bugs. Halve the search space each step.
- Add evidence, don't guess: a log/breakpoint that confirms or kills a hypothesis.
  State the hypothesis *before* you look, so the result actually decides it.

### 3. Reduce
- Shrink to the **minimal failing case** — strip everything that still fails
  without. The minimal repro usually names the cause.

### 4. Fix
- Fix the **root cause**, not the symptom. A `?.`, a try/catch, or a retry that
  hides the error is a band-aid — call it out if you're forced to ship one.
- One change at a time; re-run the repro after each. If a change doesn't move
  the failure rate, revert it before trying the next.

### 5. Guard
- Write a test that fails on the old code and passes on the fix. This is the
  deliverable — it stops the bug returning.
- If CI-flaky: the guard is determinism (fixed clock, seeded data, awaited async,
  no shared state), not a re-run.

## Anti-patterns

- Editing code before reproducing the failure.
- Multiple simultaneous changes — you won't know which fixed it (or broke more).
- "Fixed" without a regression test.
- Blaming flakiness/"the environment" without reproducing it.
- Swallowing the error to make the symptom disappear.

## Production errors (Sentry et al.)

- Group by fingerprint; fix the highest-frequency × highest-impact first.
- Use the breadcrumb trail + release/commit to localize before reproducing.
- After deploy, confirm the error rate actually drops — see `deploy-smoke-test`.
