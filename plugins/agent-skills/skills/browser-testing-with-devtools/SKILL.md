---
name: browser-testing-with-devtools
description: >
  Verify a web change in a real browser — render, interaction, console, network,
  responsive layout, and visual regressions — using Chrome DevTools (via MCP
  where available) and Playwright. Use after a UI or page change, when asked to
  "check it actually works in the browser", to reproduce a visual/runtime bug,
  or to screenshot desktop + mobile before pushing. Complements Playwright e2e;
  it does not replace the e2e suite.
---

# Browser Testing with DevTools

Closes the gap between "tests pass" and "it works when a human loads it". Pairs
with `frontend-ui-engineering` and `debugging-and-error-recovery`.

## When to reach for this

- A page/component changed and you need to see it render.
- A bug only shows in the browser (layout, hydration, console error, network).
- Pre-push visual check: desktop + mobile, no overflow, no console errors.

Do **not** use this in place of the Playwright e2e suite — that's the CI gate.
Use this for interactive inspection and to author/repro before writing the
durable Playwright test.

## Driving the browser

Prefer the **Chrome DevTools MCP** when available (runtime DOM, console,
network, performance traces). Otherwise drive Playwright directly:

```bash
# about: skip the browser download, use installed Chrome
PW_BROWSER_CHANNEL=chrome pnpm exec playwright test --project=chromium --project=mobile-chrome
# travel-planner: full web e2e (Testcontainers manages the DB)
pnpm test:e2e:web
# interactive UI mode
pnpm test:e2e:web:ui
```

Ad-hoc screenshot / inspection script (no test runner):

```ts
import { chromium, devices } from "@playwright/test";
const browser = await chromium.launch({ channel: "chrome" });
for (const [name, vp] of [["desktop", { width: 1440, height: 900 }],
                          ["mobile", devices["Pixel 7"].viewport]]) {
  const page = await browser.newPage({ viewport: vp });
  const errors: string[] = [];
  page.on("console", m => m.type() === "error" && errors.push(m.text()));
  await page.goto(url, { waitUntil: "networkidle" });
  await page.screenshot({ path: `/tmp/${name}.png`, fullPage: true });
  console.log(name, "console errors:", errors);
}
await browser.close();
```

## Inspection checklist

- [ ] Renders at 1440px **and** mobile (Pixel 7 / 390px); no horizontal overflow
- [ ] **Zero console errors/warnings** (React keys, hydration mismatch, a11y)
- [ ] Network: no 4xx/5xx, no waterfalls, no oversized payloads
- [ ] Interactions work by keyboard (Tab/Enter/Esc) and pointer
- [ ] No layout shift as async content loads (watch CLS)
- [ ] `prefers-reduced-motion` respected; content visible with JS disabled
- [ ] Visual diff vs. previous state (screenshot before + after)

## Reporting

State what you observed, not just "looks fine". Attach the screenshots. Quote
any console error verbatim. If you reproduced a bug, capture the minimal repro
steps and hand them to `debugging-and-error-recovery`. Once confirmed, fold the
check into a durable Playwright assertion so CI catches the regression next time.
