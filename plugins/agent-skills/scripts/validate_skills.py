#!/usr/bin/env python3
"""Skill spec compliance gate. Mirrors .github/workflows/ci.yml:
exit 0 = clean, 1 = errors (fail), 2 = warnings-only (allow), 3 = CLI bug."""
import subprocess
import sys
from pathlib import Path

# Run from the agent-skills plugin dir so go.mod is found.
plugin_dir = Path(__file__).resolve().parent.parent
rc = subprocess.run(
    ["go", "tool", "skill-validator", "check", "skills/"],
    cwd=plugin_dir,
).returncode
sys.exit(0 if rc == 2 else rc)
