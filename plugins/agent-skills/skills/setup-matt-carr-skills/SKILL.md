---
name: setup-matt-carr-skills
description: Sets up an `## Agent skills` block in AGENTS.md/CLAUDE.md and `docs/agents/` so the engineering skills know this repo's issue tracker (GitHub or local markdown), workflow label vocabulary, domain docs, and verification commands. Run before first use of `tdd`, `improve-codebase-architecture`, `grill-with-docs`, or the lifecycle skills (`draft-spec`, `implement-spec`) — or if those skills appear to be missing context about the issue tracker, labels, domain docs, or how to verify a change.
disable-model-invocation: true
---

# Setup Matt Carr's Skills

Scaffold the per-repo configuration that the engineering skills assume:

- **Issue tracker** — where issues live (GitHub by default; local markdown is also supported out of the box)
- **Workflow labels** — the label strings that drive this repo's lifecycle (e.g. the `ai:*` labels travel-planner uses)
- **Domain docs** — where the repo's domain language and ADRs live, and the consumer rules for reading them
- **Stack & verification** — this repo's toolchain and the commands skills run to verify a change

This is a prompt-driven skill, not a deterministic script. Explore, present what you found, confirm with the user, then write.

## Process

### 1. Explore

Look at the current repo to understand its starting state. Read whatever exists; don't assume:

- `git remote -v` and `.git/config` — is this a GitHub repo? Which one?
- `AGENTS.md` and `CLAUDE.md` at the repo root — does either exist? Is there already an `## Agent skills` section in either?
- The repo's domain-language home — `CONSTITUTION.md`, root/per-layer `AGENTS.md`, or a `CONTEXT.md`/`CONTEXT-MAP.md` if one exists
- `docs/decisions/` (preferred ADR home) or a legacy `docs/adr/` directory
- `docs/agents/` — does this skill's prior output already exist?
- `.scratch/` — sign that a local-markdown issue tracker convention is already in use

### 2. Present findings and ask

Summarise what's present and what's missing. Then walk the user through the three decisions **one at a time** — present a section, get the user's answer, then move to the next. Don't dump all three at once.

Assume the user does not know what these terms mean. Each section starts with a short explainer (what it is, why these skills need it, what changes if they pick differently). Then show the choices and the default.

**Section A — Issue tracker.**

> Explainer: The "issue tracker" is where issues live for this repo. The lifecycle skills (`draft-spec`, `implement-spec`) and other engineering skills read from and write to it — they need to know whether to call `gh issue create`, write a markdown file under `.scratch/`, or follow some other workflow you describe. Pick the place you actually track work for this repo.

Default posture: these skills were designed for GitHub. If a `git remote` points at GitHub, propose that. If a `git remote` points at GitLab (`gitlab.com` or a self-hosted host), propose GitLab. Otherwise (or if the user prefers), offer:

- **GitHub** — issues live in the repo's GitHub Issues (uses the `gh` CLI)
- **GitLab** — issues live in the repo's GitLab Issues (uses the [`glab`](https://gitlab.com/gitlab-org/cli) CLI)
- **Local markdown** — issues live as files under `.scratch/<feature>/` in this repo (good for solo projects or repos without a remote)
- **Other** (Jira, Linear, etc.) — ask the user to describe the workflow in one paragraph; the skill will record it as freeform prose

**Section B — Workflow labels.**

> Explainer: The autonomous lifecycle is label-driven — opening or merging an issue/PR with a given label fires the matching routine. The skills need to apply the label strings *this repo has actually configured*. Map them here so skills apply the right ones instead of creating duplicates.

The default vocabulary is travel-planner's `ai:*` lifecycle (ADR 057):

- `ai:plan` — issue → draft a SPEC (fires `draft-spec`)
- `ai:plan-epic` — issue → draft an EPIC (fires `draft-epic`)
- `ai:revise-now` — spec/epic PR → rewrite from review feedback (fires `revise-spec`)
- `ai:implement` — spec PR merged with this label → implement (fires `implement-spec`)
- `ai:done` — implementation PR is ready for review
- `ai:blocked` — a routine hit a wall and needs a human
- `ai:planned` — issue already drafted; routine won't redo it

Default: these strings as-is. Ask the user to override any that differ in this repo, or to describe a different lifecycle if the repo doesn't use `ai:*` labels.

**Section C — Domain docs.**

> Explainer: Some skills (`improve-codebase-architecture`, `grill-with-docs`, `tdd`) need to learn the project's domain language and read past architectural decisions. They need to know *where those live in this repo*. Don't assume a `CONTEXT.md` exists — name the repo's actual home.

Record two things:

- **Domain language** — where the project's vocabulary lives. Default to `CONSTITUTION.md` + root/per-layer `AGENTS.md` (travel-planner's model). A `CONTEXT.md` glossary is optional: use it if present, but don't assert it.
- **ADRs** — `docs/decisions/NNN-title.md` (the standard, with a README index and the `write-adr` skill). A legacy `docs/adr/` is tolerated if a repo already uses it, but new repos use `docs/decisions/`.

**Section D — Stack & verification commands.**

> Explainer: The universal-method skills (`tdd`, `debugging-and-error-recovery`,
> `security-and-hardening`, `code-review`, `implement-spec`, …) are
> language-agnostic — they describe the *method*, not the toolchain. They learn
> *this* repo's actual commands (lint, type-check, unit/integration tests,
> build) from one file so the same skill works against a TS, Go, Rust, or Java
> repo. This is what lets the skill set extend across languages without forking.

Detect the toolchain and propose the commands; don't assume one ecosystem:

- **Package manager / build**: presence of `package.json` (npm/pnpm/yarn),
  `go.mod` (Go), `Cargo.toml` (Rust), `pom.xml` / `build.gradle` (Java).
- **Lint/format, type-check, unit, integration/e2e, build**: read the repo's
  scripts (`package.json` `scripts`, `Makefile`, `justfile`, CI workflow) and
  propose the real commands. If the repo already documents a verification
  table in `AGENTS.md`/`CONSTITUTION.md`, point at it rather than duplicating.

Confirm the detected commands with the user before writing.

### 3. Confirm and edit

Show the user a draft of:

- The `## Agent skills` block to add to whichever of `CLAUDE.md` / `AGENTS.md` is being edited (see step 4 for selection rules)
- The contents of `docs/agents/issue-tracker.md`, `docs/agents/workflow-labels.md`, `docs/agents/domain.md`, `docs/agents/verification.md`

Let them edit before writing.

### 4. Write

**Pick the file to edit:**

- If `CLAUDE.md` exists, edit it.
- Else if `AGENTS.md` exists, edit it.
- If neither exists, ask the user which one to create — don't pick for them.

Never create `AGENTS.md` when `CLAUDE.md` already exists (or vice versa) — always edit the one that's already there.

If an `## Agent skills` block already exists in the chosen file, update its contents in-place rather than appending a duplicate. Don't overwrite user edits to the surrounding sections.

The block:

```markdown
## Agent skills

### Issue tracker

[one-line summary of where issues are tracked]. See `docs/agents/issue-tracker.md`.

### Workflow labels

[one-line summary of the lifecycle label vocabulary]. See `docs/agents/workflow-labels.md`.

### Domain docs

[one-line summary of layout — "single-context" or "multi-context"]. See `docs/agents/domain.md`.

### Stack & verification

[one-line summary of stack — e.g. "Next.js / pnpm / Biome / Vitest / Playwright"]. The commands universal skills run live in `docs/agents/verification.md`.
```

Then write the three docs files using the seed templates in this skill folder as a starting point:

- [issue-tracker-github.md](./references/issue-tracker-github.md) — GitHub issue tracker
- [issue-tracker-gitlab.md](./references/issue-tracker-gitlab.md) — GitLab issue tracker
- [issue-tracker-local.md](./references/issue-tracker-local.md) — local-markdown issue tracker
- [workflow-labels.md](./references/workflow-labels.md) — lifecycle label mapping
- [domain.md](./references/domain.md) — domain doc consumer rules + layout
- [verification.md](./references/verification.md) — stack + verification command table (the injection point for universal skills)

For "other" issue trackers, write `docs/agents/issue-tracker.md` from scratch using the user's description.

### 5. Done

Tell the user the setup is complete and which engineering skills will now read from these files. Mention they can edit `docs/agents/*.md` directly later — re-running this skill is only necessary if they want to switch issue trackers or restart from scratch.
