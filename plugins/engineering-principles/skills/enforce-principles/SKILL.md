---
name: enforce-principles
description: Set up or extend mechanical enforcement (Constitution P14) of the engineering principles in a consuming repo. Translates the anchors that apply to the project into architecture tests, lint rules, and CI gates, then writes the repo's enforcement map so future reviews can defer to the gates. Activate when bootstrapping a new repo or service, when the user asks "make this repo enforce the principles", or when architecture-review flags a relevant rule that has no gate behind it.
---

# Enforce Principles

A rule that lives only in documentation will be broken (V1). This skill turns the principles that apply to a project into checks the build runs — so agents and humans are constrained by the harness, not by memory.

## Step 0 — Load the principles

Resolve and read `principles/{index,values,constitution,cloud-native,tech-stack}.md` exactly as `apply-principles` does (plugin root first, then repo root, then workspace search). Halt and report if they can't be found.

## Step 1 — Inventory the repo

Establish what exists before adding anything (Rule 7):

- Stack: language, framework, package manager, deploy target.
- Existing gates: test suites, architecture tests, type-checker strictness, lint config, commit hooks, CI workflows, secret scanners.
- Existing record: project constitution, ADR log, enforcement map if one exists.

## Step 2 — Pick the enforceable anchors

From `principles/index.md`, list the anchors that apply to this project's surfaces **and** can be checked mechanically. The usual set, in priority order:

1. **P2 / P13 / P14** — layer-import architecture tests, composition-root boundary, no framework imports in domain, purity rules (`patterns/architecture-tests.md`).
2. **P15** — commit message lint (commitlint or equivalent) in a commit-msg hook and CI.
3. **P20 / P21** — secret scanning (gitleaks/trufflehog), lockfile + pinned dependencies, vulnerability scan blocking on high/critical, Actions pinned to SHA.
4. **P12 / P11** — spec generated from source with a CI drift check; breaking-change check on contracts where tooling exists.
5. **P7** — test suites split so unit (fast, no I/O) and integration (real backing services) run as separate CI jobs.
6. **P16** — changelog-entry check for user-facing change, where the repo distinguishes it.
7. **P18** — axe (or equivalent) audit in CI at representative viewports, for repos with a UI.
8. **P22 / P6** — lint or test rules against default-infinite clients and unkeyed external writes, where the ecosystem has tooling.

Skip anchors the project has overridden by ADR — don't build a gate the project decided against.

## Step 3 — Implement the gaps

For each missing gate, smallest honest version first (Rule 2, V2):

- The check runs in the standard test/CI path — a red architecture test, not a bespoke tool someone must remember to run.
- Failures are loud and self-explaining (V4): the assertion message names the anchor and the fix.
- One gate per commit where practical (P17), Conventional Commit messages (P15).

## Step 4 — Write the enforcement map

Add or update a table in the project's constitution (or `docs/` if none): one row per enforced anchor — anchor, what's checked, where the check lives, how it runs. This is what `architecture-review` Step 0.5 reads to mark items `enforced` instead of re-reviewing them.

## Step 5 — Report what is NOT enforced

End with the honest remainder: relevant anchors that stay review-only (e.g. P4 root-cause, P28 simplicity — judgment calls no gate can make) and any that need tooling the ecosystem lacks. Don't claim full coverage that doesn't exist (Rule 9).

## Output template

Produce a summary with this structure:

```markdown
## Enforcement report

### Gates added
| Anchor | Check | Location | Runs in |
|--------|-------|----------|---------|
| P2 | Layer-import architecture test | `src/test/java/.../ArchitectureTest.java` | CI (test job) |
| P15 | Commit message lint | `.husky/commit-msg` + CI | Pre-commit hook + CI |
| P20 | Secret scanning | `.pre-commit-config.yaml` (gitleaks) | Pre-commit hook + CI |

### Enforcement map

Add this table to the project's constitution or `docs/enforcement-map.md`:

| Anchor | What is checked | Gate location | How it runs | Added |
|--------|----------------|---------------|-------------|-------|
| P2 | Domain imports nothing from infrastructure | `src/test/.../ArchitectureTest.java` | `./gradlew test` in CI | 2026-06-14 |
| P15 | Commit message format | commitlint config | commit-msg hook + CI job | 2026-06-14 |

### Not enforced (review-only)
| Anchor | Why |
|--------|-----|
| P4 (Root cause) | Judgment call — no gate can verify structural vs symptomatic |
| P28 (Simplicity) | Subjective — requires human review |
| P13 (Domain modelling) | Partial — architecture test covers imports but not Result types or value objects |
```
