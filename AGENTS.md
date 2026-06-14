# AGENTS.md

Owner: Matt.
Work style: telegraph; noun-phrases ok; drop grammar; min tokens.

## What this repo is

- Monorepo for AI agent plugins: marketplace catalog + two plugins.
- Plugin scope: `plugins/agent-skills/skills/` and `plugins/engineering-principles/`.
- Manifests: `.claude-plugin/` at root (marketplace) and in each plugin subtree.
- See `README.md` for install/distribution details.

## Plugins

### agent-skills (`plugins/agent-skills/`)

- Every skill follows the [agentskills.io](https://agentskills.io) open spec.
- Layout: `skills/<name>/SKILL.md` + optional `references/`, `scripts/`, `assets/` subdirs.
- Root of a skill dir contains only `SKILL.md`. Supplementary docs go in `references/`.
- Frontmatter: `name` (≤64 chars, kebab-case, must match dir name) + `description` (≤1024 chars; write to self-activate, e.g. "Use when user wants to …").
- Internal links use relative paths; cross-skill refs use `../<other-skill>/…`.
- After edits, run a quick validation pass (frontmatter present, name matches dir, links resolve).

### engineering-principles (`plugins/engineering-principles/`)

- Cross-repo engineering principles, patterns, behavioural rules.
- Read order (precedence): `behavioural-rules.md` → `principles/values.md` → `principles/index.md` → `principles/constitution.md` → `principles/cloud-native.md` → `principles/tech-stack.md` → `patterns/*` → `skills/*`.
- Skills: `apply-principles`, `architecture-review`, `enforce-principles`.
- Writeback policy (for maintainer's fork only): agents may push new patterns, new principles, clarifications, and reference/link fixes directly to `main`. Removals, inversions, tier restructures, and behavioural-rule changes need a PR.

## Versioning (P0)

- [release-please](https://github.com/googleapis/release-please) handles bumps. **Do not** edit `version` fields directly.
- Multi-package mode: each plugin is independently versioned with tags like `agent-skills-v3.1.0` and `engineering-principles-v1.2.0`.
- Conventional Commits drive the bump:
  - `fix:` → **patch**
  - `feat:` → **minor**
  - `feat!:` / `BREAKING CHANGE:` → **major**
  - `docs:`, `chore:`, `refactor:`, `test:`, `style:` → no release
- Scope commits to the plugin when changes are scoped to one (e.g. `feat(tdd): …`, `docs(principles): …`).

## Naming (P0)

- Skill dir + frontmatter `name`: kebab-case, lowercase, no underscores.
- File names inside a skill: kebab-case for prose docs, SCREAMING-KEBAB for short canonical reference docs.
- Conventional Commits everywhere.

## Style (P0)

- Replies: telegraph; noun phrases OK; drop filler; minimal tokens. No "AI slop".
- Markdown: match existing style of the touched file.
- Skills: imperative voice, concise sections, examples over abstractions.

## PRs (P0)

- Open as draft; mark ready only when title + description match the work.
- Title: Conventional Commits.
- Description: prioritised bullets, one item per line, no blank lines.

## Git (P0)

- Safe by default: `git status/diff/log`. Push only when asked.
- Destructive ops require explicit ask. No amend unless asked.

## Local setup

```sh
go test ./tools/...                                       # unit tests
go run ./tools/validate-plugins                           # validate all plugin.json
go run ./tools/check-naming                               # naming conventions
go run ./tools/check-markdown                             # markdown hygiene
lefthook install                                          # git hooks (optional)
```

## Critical thinking

- Fix root cause, not band-aid.
- Unsure: read more code; if still stuck, ask with short options.
- Conflicts: call out; pick safer path.
