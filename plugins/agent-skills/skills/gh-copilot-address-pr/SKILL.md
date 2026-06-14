---
name: gh-copilot-address-pr
description: Address GitHub Copilot (and other bot reviewer) PR review comments for a specific PR. Use when asked to implement or evaluate Copilot/bot PR review comments, decide which to accept, apply changes, then commit and (on confirmation) push. This is the canonical detect → triage → apply → reply → resolve-thread loop that `babysit-pr` delegates to for review feedback.
---

# Bot PR review comment workflow

Owns the comment-triage discipline for PR review feedback (Copilot, Claude, Wiz,
other `[bot]` reviewers). `babysit-pr` delegates here rather than re-deriving it.

## Tool conventions

GitHub I/O has two supported paths — use whichever the environment provides:
`gh` CLI (interactive / local, shown below) or `mcp__github__*` MCP tools
(scheduled / routine contexts, where `gh` has auth issues —
anthropics/claude-code#42743). Pick one and stay consistent within a run;
translate the `gh` commands below to their MCP equivalents as needed.

## Inputs

- PR reference: URL or `owner/repo#123` or `123` with repo.
- Copilot reviewer login string if known (default: contains `copilot`).
- If the PR is already open on the current branch, prefer using the current branch directly (no PR number needed).

## Steps

1) Resolve PR ref (no URLs in commands)
- If user gives URL, extract `owner/repo` + PR number.
- If the PR is already open for the current branch, use:
  - `gh pr view --json comments,files,number,title,url,headRefName` (current branch)
- Otherwise use:
  - `gh pr view <num> -R owner/repo --json comments,files,number,title,url,headRefName`
- Use `gh pr diff <num> -R owner/repo` for full context when needed.

2) Isolate bot comments
- Scan comments for the bot login (Copilot is often `Copilot` / contains `copilot`; others: `claude`, `wiz`, `[bot]`).
- If unsure which reviewer to act on, ask user.
- Optional filter (default Copilot; widen the regex for other bots): `gh api /repos/OWNER/REPO/pulls/NUM/comments --jq '.[] | select(.user.login|test("copilot|claude|wiz|\\[bot\\]";"i")) | {id,path,position,body,user}'` (skip if jq missing).
- Skip already-replied comments when possible (e.g., those with `in_reply_to_id`).
- Ignore comments that are already resolved. Use GraphQL to fetch `reviewThreads` and filter out threads where `isResolved` is true; only act on unresolved bot threads.

3) Triage each comment
- Read code context in file + surrounding lines.
- Decide: accept / reject / unsure.
- Accept only if: correctness improved, style consistent, no new risk.
- If unsure, ask user before changing.

4) Implement accepted changes
- Make small, reviewable edits.
- Prefer repo tooling; avoid new deps unless needed.
- Add or update tests if behavior changes and fits scope.
- Keep edits to PR scope unless fix requires wider change.

5) Verify
- Run relevant tests or checks if fast and configured.
- If blocked, state what’s missing.

6) Close the loop on the bot comments
- For each accepted comment, add a reply noting it was addressed + reviewed.
- If a review thread is open, resolve it after applying the fix.
- Use gh API if needed:
-  - Reply to a PR review comment (note `-F` for numeric type): `gh api -X POST /repos/OWNER/REPO/pulls/NUM/comments -F body="Addressed and reviewed." -F in_reply_to=COMMENT_ID`
- Resolve a review thread (GraphQL): fetch thread IDs via `gh api graphql` query for `reviewThreads` (not in `gh pr view`), check `isResolved`, then call `resolveReviewThread` only when false.
- If a thread cannot be resolved via CLI, post the reply and note the limitation in the summary.

7) Confirm before commit + push
- Ask for confirmation before commit + push.

8) Commit + push (only after explicit confirmation)
- Single Conventional Commit covering all accepted changes.
- Push to current PR branch.

## Guardrails

- Never blindly accept all comments.
- If any comment unclear, stop and ask user.
- Use `gh pr view/diff` for PRs (no URL browsing).
