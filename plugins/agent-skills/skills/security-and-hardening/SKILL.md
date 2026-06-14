---
name: security-and-hardening
description: >
  Find and fix web-app security issues across the OWASP Top 10: authn/authz,
  injection, secrets, SSRF, XSS, CSRF, insecure deserialization, access-control
  gaps, and dependency CVEs. Use when adding auth or an API endpoint, handling
  user input or secrets, reviewing a diff for security, or responding to a
  vulnerability report. Stack-aware for Next.js route handlers, server actions,
  Drizzle/Postgres, and bearer-token / session auth.
---

# Security & Hardening

Defensive review for authorized work on your own repos. Pairs with
`../dependency-review` (CVE triage) and the principles review (the
`architecture-review` skill, or the repo's implementation-review step).

## Threat checklist (OWASP Top 10, web-app slice)

**Access control (the #1 real-world bug)**
- Every endpoint / server action authenticates **and** authorizes. Auth ≠ authz:
  prove the *current user* owns the *specific resource*, not just that they're
  logged in.
- No IDOR: never trust an id from the client to scope a query — filter by
  `userId` server-side in the query, not after fetching.
- Server actions and `/api/*` enforce their own auth. Confirm the middleware /
  proxy matcher does not silently exclude a protected path.

**Injection**
- Use parameterized queries / the ORM query builder. Never string-concat SQL.
- Validate + parse all input at the boundary with a schema (zod). Reject, don't
  coerce. Untrusted input is never a type until it's parsed.

**XSS / output**
- No `dangerouslySetInnerHTML` with user content. If unavoidable, sanitize with
  a vetted allowlist library.
- Set a Content-Security-Policy. Escape by default (React does — don't defeat it).

**Secrets**
- Zero secrets in source, logs, error messages, or client bundles. Server-only
  env vars must never be imported into a client component.
- Secrets come from env / a secret store. Rotate on suspected exposure.
- `git grep` for high-entropy strings before push; run secret-scanning on the PR.

**Auth specifics**
- Tokens: short-lived, signed with a key **distinct** from the session secret;
  verify signature + expiry + audience server-side every request.
- Session cookies: `HttpOnly`, `Secure`, `SameSite`. CSRF protection on
  state-changing form posts.

**SSRF / outbound**
- Never fetch a user-supplied URL without an allowlist. Block internal ranges.

**Dependencies**
- Triage CVEs via `../dependency-review`. Distinguish runtime-reachable from
  dev-only/build-time transitive (don't panic-bump the latter).

## Process

1. Map the change's trust boundaries: where does untrusted data enter?
2. Walk each boundary against the checklist above.
3. Rank findings: **Critical** (exploitable now) / **High** / **Advisory**.
   Lead with the exploit path, not the category name.
4. Fix root cause (validate at the boundary, scope the query) — not a band-aid
   downstream.
5. Add a regression test that encodes the attack (e.g. user B cannot read user
   A's trip) so the gate is mechanical, not doc-only.

## Out of scope / refuse

This skill is for hardening your own authorized systems. Decline requests to
build attacks against third parties, evade detection for malicious ends, or
target systems you don't own.
