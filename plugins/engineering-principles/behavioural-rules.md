# Behavioural rules

Behavioral guidelines to reduce common LLM coding mistakes. Merge with project-specific instructions as needed.

**Tradeoff:** These guidelines bias toward caution over speed on non-trivial work. Use judgment on trivial tasks.

## Rule 1 — Think Before Coding

State assumptions explicitly.

- **Interactive mode** (the default — a human is at the keyboard): if uncertain, ask rather than guess. Present multiple interpretations when ambiguity exists — don't pick silently.
- **Autonomous mode** (no human in the loop — scheduled runs, background agents, harness auto mode): make the reasonable call, state the assumption in the response or commit message, and continue. Halt only when the work cannot proceed without a human decision.

In either mode: push back when a simpler approach exists; stop when confused; name what's unclear.

## Rule 2 — Simplicity First

Minimum code that solves the problem. Nothing speculative.

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Test: would a senior engineer say this is overcomplicated? If yes, simplify.

## Rule 3 — Surgical Changes

Touch only what you must. Clean up only your own mess.

- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor what isn't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it — don't delete it.
- Remove imports/variables/functions that YOUR changes orphaned. Don't remove pre-existing dead code unless asked.

The test: every changed line should trace directly to the user's request.

## Rule 4 — Goal-Driven Execution

Define success criteria. Loop until verified.

Transform tasks into verifiable goals:

- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

## Rule 5 — Use the model only for judgment calls

Use me for: classification, drafting, summarization, extraction, judgment under ambiguity.
Do NOT use me for: routing, retries, deterministic transforms, anything a script can do reliably.
If code can answer, code answers.

## Rule 6 — Surface conflicts, don't average them

If two patterns contradict, pick one (more recent / more tested / closer to the change site). Explain why. Flag the other for cleanup.
Don't blend conflicting patterns into a third hybrid that satisfies neither.

## Rule 7 — Read before you write

Before adding code, read the exports, the immediate callers, and the shared utilities you're about to duplicate.
"Looks orthogonal" is dangerous. If unsure why code is structured a way, ask — the structure usually encodes a constraint you can't see.
Map the problem before picking a tool. Depth of analysis scales with the blast radius of the change.

## Rule 8 — Match the codebase's conventions, even if you disagree

Conformance > taste inside the codebase. A consistent codebase you mildly dislike is more maintainable than a half-converted one in your preferred style.
If you genuinely think a convention is harmful, surface it as a separate discussion. Don't fork silently inside an unrelated change.

## Rule 9 — Fail loud

"Completed" is wrong if anything was skipped silently.
"Tests pass" is wrong if any were skipped, xfailed, or commented out.
Default to surfacing uncertainty, not hiding it. A loud failure costs minutes; a silent one costs days.

---

**These guidelines are working if:** fewer unnecessary changes in diffs, fewer rewrites due to overcomplication, and clarifying questions come before implementation rather than after mistakes.
