# Stack & Verification

The single source of truth for this repo's toolchain and the commands the
universal engineering skills (`tdd`, `debugging-and-error-recovery`,
`security-and-hardening`, `code-review`, `implement-spec`, …) should run.
Skills source commands from here rather than hard-coding any ecosystem's
tooling — the same skill works against a TS, Go, Rust, or Java repo by reading
that repo's copy of this file.

> Seed template. Replace the bracketed values with this repo's real stack and
> commands (detected from `package.json` / `go.mod` / `Cargo.toml` / `pom.xml`,
> the `Makefile`/`justfile`, and the CI workflow). Delete the rows that don't
> apply. If the repo already has a verification table in `AGENTS.md`, point at
> it instead of duplicating — keep one source of truth.

## Stack

- **Runtime/framework**: [e.g. Next.js + React + TS / Spring Boot / Axum]
- **Package manager / build**: [pnpm / go / cargo / maven|gradle]
- **Lint/format**: [Biome / golangci-lint / clippy + rustfmt / spotless]
- **Unit tests**: [Vitest / go test / cargo test / JUnit]
- **Integration / e2e**: [Playwright / testcontainers / …]

## Verification — run before pushing

| You changed…        | Run                                  |
| ------------------- | ------------------------------------ |
| [any source]        | `[lint] && [type-check] && [unit]`   |
| [logic / use cases] | the above, plus `[integration]`      |
| [anything]          | `[build]` before pushing             |

```bash
# Worked examples per ecosystem — keep only this repo's:
# Node/TS:  pnpm check && pnpm exec tsc --noEmit && pnpm test && pnpm build
# Go:       golangci-lint run && go vet ./... && go test ./... && go build ./...
# Rust:     cargo clippy -- -D warnings && cargo test && cargo build --release
# Java:     mvn -q verify   (spotless:check + compile + test in one)
```
