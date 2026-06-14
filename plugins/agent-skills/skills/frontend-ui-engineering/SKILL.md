---
name: frontend-ui-engineering
description: >
  Build and review React/Next.js UI: component architecture, design-system
  tokens, state boundaries, server vs client components, and WCAG 2.1 AA
  accessibility. Use when adding or changing a user-facing component or page,
  designing a design system / token set, fixing an accessibility failure, or
  reviewing frontend work for structure and a11y. Stack-aware for Next.js App
  Router, React Server Components, TypeScript, and Playwright/axe gates.
---

# Frontend UI Engineering

Pairs with `../tdd` (write the test first), `../../engineering-principles`
(cite anchors), and `browser-testing-with-devtools` (verify in a real browser).

## Component architecture

- Default to **Server Components**. Add `"use client"` only when you need
  state, effects, refs, or browser APIs — and push it to the leaf, not the page.
- One responsibility per component. If props exceed ~6 or the name needs "and",
  split it.
- **Presentational vs container**: data-fetching/orchestration in server
  components or hooks; rendering in dumb components that take props. Keeps the
  rendering layer trivially testable.
- Co-locate: `component.tsx`, `component.test.tsx`, sibling styles. No
  `components/shared/misc` dumping ground.
- Lift state to the lowest common ancestor that needs it — no higher. URL/search
  params and server state beat client state; reach for `useState` last.

## Design system / tokens

- Never hard-code colour, spacing, radius, font size inline. Reference tokens
  (Tailwind theme / CSS custom properties). A raw `#hex` or `13px` in a
  component is a review finding.
- One spacing scale, one type scale, one colour ramp. New value → add a token,
  don't one-off it.
- Variants belong in the component API (`variant="primary"`), not in callers
  composing class strings.

## Accessibility (WCAG 2.1 AA — non-negotiable)

- Semantic HTML first: `<button>` for actions, `<a>` for navigation, one `<h1>`,
  headings in order. A `<div onClick>` is a bug.
- Every interactive element: keyboard reachable, visible focus ring, accessible
  name. Test with Tab/Shift-Tab/Enter/Esc — no mouse.
- Images: meaningful `alt`, or `alt=""` if decorative. Icon-only buttons need
  `aria-label`.
- Colour contrast ≥ 4.5:1 text / 3:1 large text & UI. Never signal state by
  colour alone.
- Respect `prefers-reduced-motion` for any animation. **Content must be present
  in server HTML** — animation is progressive enhancement, never the only path
  to visible content.
- Forms: `<label>` tied to every input; errors announced (`aria-live`,
  `aria-describedby`), not colour-only.

## Verify before pushing

- [ ] `pnpm lint && pnpm type-check`
- [ ] Component test (RTL/Vitest) asserts behaviour + role queries, not markup
      shape — query by `getByRole`/`getByLabelText`, not test ids
- [ ] axe/Playwright a11y check passes (both repos gate on this)
- [ ] Keyboard-only walkthrough of the new interaction
- [ ] Screenshot desktop (1440px) + mobile, check no horizontal overflow
      (see `browser-testing-with-devtools`)
- [ ] No layout shift: reserve space for async content

## Review heuristics (red flags)

- `"use client"` at page level when a leaf would do
- Business logic inside a component instead of a hook / use-case
- Raw design values instead of tokens
- `useEffect` for data that should be fetched on the server
- Click handlers on non-interactive elements
- New component with no test and no a11y assertion
