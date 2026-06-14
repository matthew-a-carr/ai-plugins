#!/usr/bin/env python3
"""Validate the shape of .claude-plugin/plugin.json.

Checks that the file exists, parses as a JSON object, and carries the
required keys (name, version, description) as non-empty strings.

Run: uv run scripts/validate_plugin.py [path]
Exits non-zero on any problem, printing each one to stderr.
"""

import json
import sys
from pathlib import Path

DEFAULT_PATH = Path("plugins/engineering-principles/.claude-plugin/plugin.json")
REQUIRED_KEYS = frozenset({"name", "version", "description"})


def validate(path: Path) -> list[str]:
    """Return a list of problems with the plugin manifest; empty means valid."""
    if not path.exists():
        return [f"missing {path}"]
    try:
        data = json.loads(path.read_text())
    except json.JSONDecodeError as exc:
        return [f"{path} is not valid JSON: {exc}"]
    if not isinstance(data, dict):
        return [f"{path} must contain a JSON object, got {type(data).__name__}"]
    errors = []
    missing = REQUIRED_KEYS - data.keys()
    if missing:
        errors.append(f"{path} missing required keys: {sorted(missing)}")
    for key in sorted(REQUIRED_KEYS & data.keys()):
        value = data[key]
        if not isinstance(value, str) or not value.strip():
            errors.append(f"{path}: {key!r} must be a non-empty string")
    return errors


def main(argv: list[str]) -> int:
    path = Path(argv[1]) if len(argv) > 1 else DEFAULT_PATH
    errors = validate(path)
    if errors:
        for error in errors:
            print(error, file=sys.stderr)
        return 1
    print(f"{path}: ok")
    return 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv))
