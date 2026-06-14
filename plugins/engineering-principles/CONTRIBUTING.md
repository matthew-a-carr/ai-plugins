# Contributing

Thanks for considering a contribution. This repo is the source of engineering principles distributed as a Claude Code plugin; the bar for changes is "does this principle hold across multiple repos and multiple PRs, with surviving review."

## How to propose a change

- **Open an issue first** for anything that adds, removes, or inverts a principle. Cite at least two project repos or PRs where the pattern recurred. Drive-by edits to the tier files will be closed.
- **Open a PR directly** for typo fixes, link repairs, kebab-case violations, broken cross-references.
- **Patterns** (`patterns/<name>.md`) are recipes, not essays. Each fits in one screen.

## Commits

[Conventional Commits v1.0.0](https://www.conventionalcommits.org/). See Constitution P15.

## Local setup

```sh
go run ./tools/check-naming
go run ./tools/check-markdown
lefthook install
```

CI runs the same hooks.

## Style

Plain language; see `AGENTS.md` § Style rules for banned phrases. Cite principles by anchor (P5, C8, T3) — anchors are stable, wording isn't.

## Scope

Project-specific decisions belong in *that project's* repo, in `docs/decisions/NNNN-<slug>.md`. This repo is principles only.

## Maintainer writeback

The repo maintainer's agents may push directly to `main` for the categories listed in `AGENTS.md` § Writeback policy. External contributors always open PRs.
