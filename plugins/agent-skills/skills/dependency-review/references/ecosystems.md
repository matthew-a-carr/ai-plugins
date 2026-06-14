# Ecosystem Commands

Per-ecosystem crib for detect / list outdated / update / audit. Use native tooling; respect the repo's existing package manager (lockfile tells you which).

## Node (package.json)

Lockfile decides manager: `package-lock.json` → npm, `pnpm-lock.yaml` → pnpm, `yarn.lock` → yarn.

| Task | npm | pnpm | yarn |
|------|-----|------|------|
| Outdated | `npm outdated` | `pnpm outdated` | `yarn outdated` |
| In-range update | `npm update` | `pnpm update` | `yarn upgrade` |
| Specific/major | `npm install pkg@latest` | `pnpm add pkg@latest` | `yarn add pkg@latest` |
| Audit | `npm audit` | `pnpm audit` | `yarn npm audit` |

Watch: `peerDependencies` conflicts, `engines` bumps, changed exports/types on majors.

## Go (go.mod)

```
go list -u -m all                  # outdated (current [available])
go get -u ./... && go mod tidy     # minor/patch within major
go get pkg/v2@latest               # major = new import path /vN — update imports too
govulncheck ./...                  # audit (install: golang.org/x/vuln/cmd/govulncheck)
```

Majors in Go change the module path (`/v2`, `/v3`): update every import.

## Python (pyproject.toml / requirements*.txt)

Tool from lockfile/config: `uv.lock` → uv, `poetry.lock` → poetry, else pip.

| Task | pip | uv | poetry |
|------|-----|----|--------|
| Outdated | `pip list --outdated` | `uv pip list --outdated` | `poetry show --outdated` |
| Update | edit requirements + `pip install -r` | `uv lock --upgrade && uv sync` | `poetry update` |
| Specific | `pip install pkg==X.Y.Z` | `uv add pkg==X.Y.Z` | `poetry add pkg@^X.Y.Z` |
| Audit | `pip-audit` | `uv tool run pip-audit` | `pip-audit` |

## Rust (Cargo.toml)

```
cargo update --dry-run                       # preview in-range updates
cargo update                                 # apply (lockfile-only, semver-compatible)
cargo upgrade -p pkg                         # bump manifest ranges (cargo-edit; majors)
cargo audit                                  # audit (cargo install cargo-audit)
```

After update: `cargo build && cargo test && cargo clippy`.

## GitHub Actions (.github/workflows/*.yml)

```
grep -rn 'uses:' .github/workflows/
```

- Compare pinned versions against latest: `gh release view -R actions/checkout`.
- Prefer major-tag pins (`@v4`) or full SHA pins; update SHAs to the latest release commit.
- Dependabot covers these with `package-ecosystem: "github-actions"`.

## Docker (Dockerfile)

```
grep -n 'FROM' Dockerfile*
```

- Check base image tags against upstream (Docker Hub / registry release notes).
- Prefer specific tags over `latest`; mind OS/runtime major bumps (e.g. `node:20` → `node:22`).
- Dependabot covers with `package-ecosystem: "docker"`.

## Dependabot starter config (.github/dependabot.yml)

Include one `updates:` entry per detected ecosystem:

```yaml
version: 2
updates:
  - package-ecosystem: "npm"          # or gomod / pip / cargo / github-actions / docker
    directory: "/"
    schedule:
      interval: "weekly"
    groups:
      minor-and-patch:
        update-types: ["minor", "patch"]
```
