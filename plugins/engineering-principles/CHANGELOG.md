# Changelog

All notable changes to this repo are documented here.
Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/);
this repo follows [SemVer](https://semver.org/) for the plugin.

## [1.1.1](https://github.com/matthew-a-carr/engineering-principles/compare/v1.1.0...v1.1.1) (2026-06-14)


### Bug Fixes

* **deps:** sync uv.lock version with pyproject.toml ([89db527](https://github.com/matthew-a-carr/engineering-principles/commit/89db52730c6b20d7bdab79cca5baa3478887848e))

## [1.1.0](https://github.com/matthew-a-carr/engineering-principles/compare/v1.0.0...v1.1.0) (2026-06-12)


### Features

* engineering-principles Claude Code plugin ([fd53a85](https://github.com/matthew-a-carr/engineering-principles/commit/fd53a859337216d138aa26af055523f9ba88eb7d))

## [Unreleased]

## [1.0.0] — 2026-05-23

### Changed

- **BREAKING**: Plugin renamed from `principles` to `engineering-principles`.
  Distributed via the new central `matthew-a-carr` marketplace
  (`matthew-a-carr/claude-plugins`). Update `enabledPlugins` from
  `principles@engineering-principles` to
  `engineering-principles@matthew-a-carr`, and `extraKnownMarketplaces`
  to point at `matthew-a-carr/claude-plugins`.
- `architecture-review` skill now also loads `principles/values.md`
  alongside the three tier files, so reviews can ground concerns in
  values (V4 loud failure, V5 root cause, V6 reversible) when no tier
  anchor fits cleanly.
- Session manifest clarifies that relative paths in tier files resolve
  under `$CLAUDE_PLUGIN_ROOT`, and names the actual skill prefix
  (`engineering-principles:`) so other agents can address them.

## [0.1.0] — 2026-05-22

### Added

- Initial public release.
- Behavioural rules (1–9) for AI agents.
- Tier 1 Constitution (P1–P33).
- Tier 2 Cloud Native (C1–C10).
- Tier 3 Chapter Tech Stack (T1–T8, worked example).
- Patterns: idempotency, feature flags, event-driven outbox, non-breaking
  changes, strangler-fig extraction, architecture tests, timeouts and
  retries, error responses (RFC 9457), cursor pagination.
- Skills: `apply-principles`, `architecture-review`.
- SessionStart manifest hook.
