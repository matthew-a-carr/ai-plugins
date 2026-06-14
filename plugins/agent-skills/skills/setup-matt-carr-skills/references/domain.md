# Domain Docs

How the engineering skills should learn this repo's domain language and read its
architectural decisions. Replace the bracketed values with this repo's reality.

## Domain language lives in

[Where the project's vocabulary is defined. Default: `CONSTITUTION.md` + the
root and per-layer `AGENTS.md` files — travel-planner's model. A `CONTEXT.md`
glossary is optional: read it **if present**, but do not assume it exists or
create it upfront.]

When output names a domain concept (issue title, refactor proposal, test name),
use the term as the repo defines it. Don't drift to synonyms.

## Architectural decisions live in

[ADR home. Default: `docs/decisions/NNN-title.md`, with a `docs/decisions/README.md`
index. Write new ADRs with the `write-adr` skill (CONSTITUTION §7 template), which
also updates the index and any superseded ADR's status. A legacy `docs/adr/` is
only used if this repo already has one.]

Read the ADRs that touch the area you're about to work in. If your output
contradicts an existing ADR, surface it explicitly rather than silently
overriding it.

## Layout

[**Single-context** — one domain, decisions at the repo root (most repos). Or
**multi-context** — a monorepo with per-area domain docs and, optionally,
area-scoped decisions. Name which this repo is and where each area lives.]
