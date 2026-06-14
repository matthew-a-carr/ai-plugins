# Changelog

## [3.0.0](https://github.com/matthew-a-carr/dev-skills/compare/v2.0.0...v3.0.0) (2026-06-14)


### ⚠ BREAKING CHANGES

* plugin name changes from `dev-skills` to `agent-skills`. Consumers pinning `dev-skills@matthew-a-carr` must update to `agent-skills@matthew-a-carr`.

### Features

* rename dev-skills → agent-skills ([#12](https://github.com/matthew-a-carr/dev-skills/issues/12)) ([e2c4d17](https://github.com/matthew-a-carr/dev-skills/commit/e2c4d176493c268d753f4860257fa5f14e1b56f0))
* web-app skills + verification injection point + convergence cleanup ([#9](https://github.com/matthew-a-carr/dev-skills/issues/9)) ([3b51b19](https://github.com/matthew-a-carr/dev-skills/commit/3b51b19c2418d747ee8a34bea03d34639625db8c))


### Bug Fixes

* rename lifecycle labels claude:* -&gt; ai:* (vendor-neutral) ([#11](https://github.com/matthew-a-carr/dev-skills/issues/11)) ([cb6f086](https://github.com/matthew-a-carr/dev-skills/commit/cb6f0863eba7227267160fa45eb2da64b8041e6b))

## [2.0.0](https://github.com/matthew-a-carr/dev-skills/compare/v1.0.0...v2.0.0) (2026-06-13)


### ⚠ BREAKING CHANGES

* complete agent-scripts to dev-skills namespace rename ([#7](https://github.com/matthew-a-carr/dev-skills/issues/7))
* **dependency-review:** dependabot-pr-merge skill removed; superseded by dependency-review.
* plugin identifier changed from matthew-a-carr@matthew-a-carr-skills to dev-skills@matthew-a-carr. GitHub's repo-rename redirect keeps old URLs working, but settings that reference the old marketplace key or plugin name need updating.

### Features

* complete agent-scripts to dev-skills namespace rename ([#7](https://github.com/matthew-a-carr/dev-skills/issues/7)) ([caae63a](https://github.com/matthew-a-carr/dev-skills/commit/caae63ae821524f4d74fc3aa7d7b92479afa621d))
* **dependency-review:** add multi-ecosystem dependency review skill ([#6](https://github.com/matthew-a-carr/dev-skills/issues/6)) ([1838b78](https://github.com/matthew-a-carr/dev-skills/commit/1838b7838ffa25daae040d985877a77260208723))
* rename plugin to dev-skills, ship 2.0.0 ([3bf725e](https://github.com/matthew-a-carr/dev-skills/commit/3bf725ecf7d8977b4a6c8e833dc380d3cfd7540c))


### Bug Fixes

* **create-cli:** use relative reference path ([70074f6](https://github.com/matthew-a-carr/dev-skills/commit/70074f6d2c47f7b5b64fa381d15e9be849fbe90b))

## [1.0.0](https://github.com/matthew-a-carr/dev-skills/compare/v2.0.0...v1.0.0) (2026-06-12)


### ⚠ BREAKING CHANGES

* **dependency-review:** dependabot-pr-merge skill removed; superseded by dependency-review.
* plugin identifier changed from matthew-a-carr@matthew-a-carr-skills to dev-skills@matthew-a-carr. GitHub's repo-rename redirect keeps old URLs working, but settings that reference the old marketplace key or plugin name need updating.
* bristol-bin-collection no longer ships a hardcoded default postcode and changed its activation description. Treating today's state as v1.0.0.

### Features

* baseline at 1.0.0 with release-please automation ([4d3ea87](https://github.com/matthew-a-carr/dev-skills/commit/4d3ea87ef3f71b4ad367de77fb01d64106346302))
* **dependency-review:** add multi-ecosystem dependency review skill ([#6](https://github.com/matthew-a-carr/dev-skills/issues/6)) ([1838b78](https://github.com/matthew-a-carr/dev-skills/commit/1838b7838ffa25daae040d985877a77260208723))
* rename plugin to dev-skills, ship 2.0.0 ([3bf725e](https://github.com/matthew-a-carr/dev-skills/commit/3bf725ecf7d8977b4a6c8e833dc380d3cfd7540c))


### Bug Fixes

* **create-cli:** use relative reference path ([70074f6](https://github.com/matthew-a-carr/dev-skills/commit/70074f6d2c47f7b5b64fa381d15e9be849fbe90b))


### Miscellaneous Chores

* pin baseline release at 1.0.0 ([29bb00b](https://github.com/matthew-a-carr/dev-skills/commit/29bb00bd5f16e394b1fb059b24a0707ec0690d60))

## [2.0.0] (2026-05-23)


### ⚠ BREAKING CHANGES

* Plugin renamed from `matthew-a-carr` to `dev-skills`, and the GitHub repo
  from `agent-scripts` to `dev-skills`. Distribution moves under the new
  central `matthew-a-carr` marketplace at `matthew-a-carr/claude-plugins`.
  Skill prefixes change from `matthew-a-carr:<skill>` to `dev-skills:<skill>`.
  Update `enabledPlugins` from `matthew-a-carr@matthew-a-carr-skills` to
  `dev-skills@matthew-a-carr`, and `extraKnownMarketplaces` to point at
  `matthew-a-carr/claude-plugins`. GitHub's repo-rename redirect keeps old
  URLs working for now, but updating refs is recommended.

## [1.0.1](https://github.com/matthew-a-carr/agent-scripts/compare/v1.0.0...v1.0.1) (2026-05-18)


### Bug Fixes

* **create-cli:** use relative reference path ([70074f6](https://github.com/matthew-a-carr/agent-scripts/commit/70074f6d2c47f7b5b64fa381d15e9be849fbe90b))

## [1.0.0](https://github.com/matthew-a-carr/agent-scripts/compare/v1.0.0...v1.0.0) (2026-05-18)


### ⚠ BREAKING CHANGES

* bristol-bin-collection no longer ships a hardcoded default postcode and changed its activation description. Treating today's state as v1.0.0.

### Features

* baseline at 1.0.0 with release-please automation ([4d3ea87](https://github.com/matthew-a-carr/agent-scripts/commit/4d3ea87ef3f71b4ad367de77fb01d64106346302))


### Miscellaneous Chores

* pin baseline release at 1.0.0 ([29bb00b](https://github.com/matthew-a-carr/agent-scripts/commit/29bb00bd5f16e394b1fb059b24a0707ec0690d60))
