---
name: dependency-review
description: >
  Review, update, and upgrade project dependencies across ecosystems
  (npm/pnpm/yarn, Go, Python, Rust, GitHub Actions, Docker). Use when user
  wants to review/update/upgrade dependencies, triage or merge Dependabot
  PRs, audit a project for outdated or vulnerable packages, or apply a major
  version upgrade safely — even if they just say "update everything" or
  "check what's outdated". Patch/minor bumps batch and auto-apply when CI is
  green and release notes are clean; majors get one branch each with a full
  test gate. Supersedes dependabot-pr-merge.
compatibility: Requires authenticated gh CLI, network access, and the project's native package tooling
---

# Dependency Review

Review every dependency in the repo, update what's safe, upgrade majors carefully, prioritise vulnerabilities, report the rest.

## Risk policy

| Tier | Action |
|------|--------|
| Security fix | Highest priority. Apply first, any bump size. |
| Patch / minor | Batch together. Auto-apply/merge only when **confident** (guardrails below). |
| Major | One upgrade per branch/PR. Testing is key. Never auto-merge without green tests. |
| Uncertain | Skip, flag for human with reason. |

**Confidence guardrails (patch/minor auto-apply) — all must hold:**
- CI green (all required checks).
- Lockfile-only or compatible-range manifest change.
- Release notes show no breaking changes.
- No peer-dependency, config, or engine/runtime requirement changes.

## Workflow

### 1) Detect ecosystems

Find manifests; use native tooling per ecosystem. Once you know which ecosystems are present, read [references/ecosystems.md](references/ecosystems.md) for the outdated/update/audit commands of those ecosystems only.

```
find . -path ./node_modules -prune -o -path ./vendor -prune -o \
  \( -name package.json -o -name go.mod -o -name pyproject.toml \
     -o -name 'requirements*.txt' -o -name Cargo.toml -o -name 'Dockerfile*' \) -print
ls .github/workflows/ 2>/dev/null
```

Also check GitHub Actions workflow pins (`uses:` versions) and Dockerfile base image tags.

### 2) Security audit first

Run per detected ecosystem (`npm audit`, `pip-audit`, `cargo audit`, `govulncheck` — see reference), plus repo alerts:

```
gh api repos/{owner}/{repo}/dependabot/alerts --jq '.[] | select(.state=="open") | {pkg: .dependency.package.name, severity: .security_advisory.severity, fixed: .security_vulnerability.first_patched_version.identifier}'
```

Vulnerable packages jump the queue regardless of bump size.

### 3) Triage open Dependabot PRs

```
gh pr list --search "author:app/dependabot is:pr is:open" --json number,title,headRefName
```

For each PR, verify merge gates:

```
gh pr view <number> --json number,title,isDraft,mergeable,reviewDecision,statusCheckRollup
gh pr diff <number>
```

Only merge when: not draft, `mergeable` is `MERGEABLE`, all `statusCheckRollup` entries `SUCCESS`, `reviewDecision` is `APPROVED` if reviews required — **and** the confidence guardrails pass.

Check release notes per bumped dependency:

```
gh release view vX.Y.Z -R owner/repo   # fallback: gh release view X.Y.Z -R owner/repo
```

No GitHub release → read changelog with WebFetch. Scan codebase for any APIs the notes deprecate/remove.

Merge with the repo default method:

```
gh repo view --json viewerDefaultMergeMethod
gh pr merge <number> --squash --delete-branch
```

Patch/minor PRs passing all gates: merge as a batch. Major PRs: route through step 5 — check out the branch, run the full test gate first.

If `.github/dependabot.yml` missing, offer a starter config covering each detected ecosystem (template in [references/ecosystems.md](references/ecosystems.md)).

### 4) Batch patch + minor updates

For deps with no open Dependabot PR: update in-range per ecosystem (reference file has commands), one branch for the whole batch.

- Verify each bump against guardrails (release notes, peers, config).
- Run install + build + typecheck + full test suite locally before pushing.
- Anything failing a guardrail: drop from batch, record reason.

### 5) Major upgrades — one at a time

Per major bump, on its own branch/PR:

1. Read the migration guide / changelog between current and target version.
2. Apply the upgrade (`npm install pkg@latest`, edit `go.mod`/`pyproject.toml`/`Cargo.toml`, etc.).
3. Run full test suite + build + typecheck.
4. Fix compile/test breakage from API changes; follow the migration guide.
5. Green → open PR with a summary of breaking changes handled. Not green, or judgement call (behavioural change, licence change, rewrite-scale migration) → stop, flag for human with findings.

Never stack two majors on one branch. Never auto-merge a major.

### 6) Report

End with a summary:

- **Updated/merged**: package, old → new, how (Dependabot merge / batch / major PR).
- **Skipped**: package + exact guardrail that failed.
- **Majors needing attention**: package, breaking changes found, suggested next step.
- **Security**: alerts resolved / still open.
