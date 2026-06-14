#!/usr/bin/env python3
"""Enforce the repo's naming conventions:

- All directories (except .git, node_modules, .venv, etc.) are lowercase.
- All markdown files under subdirectories are lowercase, kebab-case.
- Ecosystem-convention files keep their original casing in any location:
  README.md, AGENTS.md, CLAUDE.md, SKILL.md, LICENSE, LICENSE.md, CHANGELOG.md,
  CONTRIBUTING.md, CONSTITUTION.md.

Exits non-zero on any violation, printing the offending paths.
"""

import re
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent

# Files that are allowed to keep ecosystem-convention casing anywhere.
ECOSYSTEM_FILES = {
    "README.md",
    "AGENTS.md",
    "CLAUDE.md",
    "SKILL.md",
    "LICENSE",
    "LICENSE.md",
    "CHANGELOG.md",
    "CONTRIBUTING.md",
    "CONSTITUTION.md",
}

# Directories not to scan into.
SKIP_DIRS = {
    ".git",
    "node_modules",
    ".venv",
    "venv",
    "__pycache__",
    ".pytest_cache",
    ".ruff_cache",
    "dist",
    "build",
    ".idea",
    ".vscode",
}

# Hidden top-level dirs we tolerate as-is (lowercase or convention).
ALLOWED_HIDDEN_DIRS = {".github", ".claude-plugin"}

KEBAB_RE = re.compile(r"^[a-z0-9]+(-[a-z0-9]+)*\.md$")
LOWER_DIR_RE = re.compile(r"^[a-z0-9][a-z0-9._-]*$")


def main() -> int:
    violations: list[str] = []

    for path in REPO_ROOT.rglob("*"):
        # Skip excluded trees.
        parts = path.relative_to(REPO_ROOT).parts
        if any(part in SKIP_DIRS for part in parts):
            continue

        rel = path.relative_to(REPO_ROOT)
        name = path.name

        if path.is_dir():
            # Allow .github, .claude-plugin, etc.
            if name.startswith(".") and name in ALLOWED_HIDDEN_DIRS:
                continue
            if name.startswith("."):
                # Other dotfiles dirs (e.g. .ruff_cache) shouldn't normally exist; ignore.
                continue
            if not LOWER_DIR_RE.match(name):
                violations.append(f"directory not lowercase: {rel}")
            continue

        # File checks: only enforce naming on .md files outside the repo root.
        if path.suffix != ".md":
            continue
        if path.parent == REPO_ROOT:
            # Root markdown files (README, AGENTS, CLAUDE) keep their casing.
            continue
        if name in ECOSYSTEM_FILES:
            continue
        if not KEBAB_RE.match(name):
            violations.append(f"markdown file not lowercase kebab-case: {rel}")

    if violations:
        print("Naming convention violations:")
        for v in violations:
            print(f"  - {v}")
        print(
            "\nFix by renaming. Directories must be lowercase; "
            "markdown files in subdirectories must be lowercase kebab-case "
            "(e.g. cloud-native.md). Ecosystem files keep their casing."
        )
        return 1

    return 0


if __name__ == "__main__":
    sys.exit(main())
