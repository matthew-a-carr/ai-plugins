---
name: grill-me
description: Interview the user relentlessly about a plan or design until reaching shared understanding, resolving each branch of the decision tree. If the repo has domain docs (CONTEXT.md / ADRs), also challenge the plan against the existing domain language and update those docs inline as decisions crystallise. Use when the user wants to stress-test a plan, get grilled on a design, sharpen terminology, or mentions "grill me".
---

# Grill Me

Interview me relentlessly about every aspect of this plan until we reach a shared understanding. Walk down each branch of the design tree, resolving dependencies between decisions one-by-one. For each question, provide your recommended answer.

Ask the questions one at a time, waiting for feedback on each before continuing.

If a question can be answered by exploring the codebase, explore the codebase instead.

## If the repo has domain docs

During codebase exploration, look for existing documentation. If none exists, run the plain grilling above — the rest of this section is only for repos that carry a domain model.

### File structure

Most repos have a single context:

```
/
├── CONTEXT.md
├── docs/
│   └── decisions/
│       ├── 001-event-sourced-orders.md
│       └── 002-postgres-for-write-model.md
└── src/
```

If a `CONTEXT-MAP.md` exists at the root, the repo has multiple contexts. The map points to where each one lives:

```
/
├── CONTEXT-MAP.md
├── docs/
│   └── decisions/                    ← system-wide decisions
├── src/
│   ├── ordering/
│   │   ├── CONTEXT.md
│   │   └── docs/decisions/           ← context-specific decisions
│   └── billing/
│       ├── CONTEXT.md
│       └── docs/decisions/
```

Create files lazily — only when you have something to write. If no `CONTEXT.md` exists, create one when the first term is resolved. ADRs go in `docs/decisions/` (check `docs/agents/domain.md` for this repo's ADR home, and defer to the `write-adr` skill if the repo has one).

### During the session

- **Challenge against the glossary** — when the user uses a term that conflicts with the existing language in `CONTEXT.md`, call it out immediately. "Your glossary defines 'cancellation' as X, but you seem to mean Y — which is it?"
- **Sharpen fuzzy language** — when the user uses vague or overloaded terms, propose a precise canonical term. "You're saying 'account' — do you mean the Customer or the User? Those are different things."
- **Discuss concrete scenarios** — stress-test domain relationships with specific scenarios that probe edge cases and force precision about the boundaries between concepts.
- **Cross-reference with code** — when the user states how something works, check whether the code agrees. If you find a contradiction, surface it: "Your code cancels entire Orders, but you just said partial cancellation is possible — which is right?"
- **Update CONTEXT.md inline** — when a term is resolved, update `CONTEXT.md` right there. Don't batch these up. Use the format in [references/CONTEXT-FORMAT.md](references/CONTEXT-FORMAT.md). `CONTEXT.md` is a glossary and nothing else — totally devoid of implementation details, not a spec or scratch pad.

### Offer ADRs sparingly

Only offer to create an ADR when all three are true:

1. **Hard to reverse** — the cost of changing your mind later is meaningful.
2. **Surprising without context** — a future reader will wonder "why did they do it this way?"
3. **The result of a real trade-off** — there were genuine alternatives and you picked one for specific reasons.

If any of the three is missing, skip the ADR. Use the format in [references/ADR-FORMAT.md](references/ADR-FORMAT.md).
