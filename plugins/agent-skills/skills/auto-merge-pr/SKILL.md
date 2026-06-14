---
name: auto-merge-pr
description: >
  Babysit a PR to merge: monitor CI pipelines and automated review feedback
  (Wiz security scans, GitHub Copilot, Claude, or other AI reviewers), fix
  failures and address findings, re-push until green, then merge. Use when
  user wants to auto-merge a PR, get a PR mergeable, watch a PR until it
  lands, or address CI/security/AI review feedback and merge — even if they
  just say "babysit", "land", or "ship" the PR.
compatibility: Requires authenticated gh CLI and network access
---

# Auto Merge PR

Drive one PR to merge. Loop: check gates → fix what's red → push → re-check. Merge only when everything passes.

## Merge gates

All must hold before merging:

- All required CI checks `SUCCESS`.
- No unresolved automated review threads (Copilot, Claude, Wiz, bots).
- No open Wiz / security scan findings on the PR.
- `mergeable` is `MERGEABLE`; `reviewDecision` is `APPROVED` if reviews required.
- Not draft.

Human review requirements are never bypassed — if a human approval is required and missing, stop and report.

## Loop

### 1) Snapshot the PR

```
gh pr view <num> --json number,title,isDraft,mergeable,reviewDecision,statusCheckRollup,headRefName
gh pr checks <num>
```

Check out the PR branch locally: `gh pr checkout <num>`.

### 2) Fix CI failures

For each failing check:

```
gh run list --branch <headRefName> --json databaseId,name,conclusion
gh run view <run-id> --log-failed
```

- Classify first: flaky (timeout, runner/infra error, test unrelated to the diff) vs real.
  Suspected flake → `gh run rerun <run-id> --failed` once; fails again → treat as real.
- Read the failing step's log; fix the root cause (test, lint, build, type error), not the symptom.
- Reproduce locally where possible (run the same test/lint command) before pushing.
- Commit (Conventional Commits) and push to the PR branch; checks re-run.
- Merge conflicts: `gh pr update-branch <num>` or rebase locally, resolve, push.

### 3) Address automated review feedback

Collect unresolved threads from bot/AI reviewers — match `user.login` against `copilot`, `claude`, `wiz`, `[bot]`:

```
gh api /repos/{owner}/{repo}/pulls/<num>/comments --jq '.[] | select(.user.login|test("copilot|claude|wiz|\\[bot\\]";"i")) | {id,path,body,user:.user.login}'
```

Filter out resolved threads via GraphQL `reviewThreads` (`isResolved: false` only).

Per finding, triage: accept / reject / unsure — full reply/resolve workflow in the `gh-copilot-address-pr` skill (use it if available). In short:

- **Accept**: fix it, reply "addressed" with what changed, resolve the thread.
- **Reject**: reply with the reason (false positive, out of scope), resolve the thread.
- **Unsure**: stop, ask user.

**Wiz / security findings** (PR comments or failed scan checks: vulnerabilities, secrets, IaC misconfigs):

- Vulnerable dependency → bump to the first patched version.
- Secret in diff → remove it, load from env/secret store; flag for rotation in the report.
- IaC/config finding → apply the recommended remediation from the finding body.
- Believed false positive → do NOT silently suppress; reply with justification and flag for human sign-off.

Security findings outrank style feedback — address them first.

### 4) Re-check

Push fixes, then watch checks:

```
gh pr checks <num> --watch
```

New feedback may arrive on new commits — go back to step 3. Repeat until all gates pass.

**Circuit breaker**: same check failing after 3 fix attempts, or fixes ballooning beyond PR scope → stop, summarise the diagnosis, hand to user.

### 5) Merge

Use the repo default method:

```
gh repo view --json viewerDefaultMergeMethod
gh pr merge <num> --squash --delete-branch
```

All gates passed except checks still running → enable auto-merge instead of polling:

```
gh pr merge <num> --auto --squash --delete-branch
```

Only after feedback threads are resolved — auto-merge fires on green checks, it won't wait for you.

### 6) Report

- **Merged**: PR number + title, merge method.
- **CI fixes**: check name → root cause → fix commit.
- **Feedback addressed**: finding → action (fixed / rejected + reason).
- **Flagged for human**: unresolved items, secrets needing rotation, false-positive sign-offs.

## Guardrails

- Never merge with red required checks. Never `--admin` / force-merge.
- Never dismiss or suppress a finding without a written justification in the thread.
- Never push fixes that change behaviour outside the PR's scope — flag instead.
- Unclear feedback or judgement calls: ask user, don't guess.
