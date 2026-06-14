# AI Plugins

Marketplace catalog and plugins for AI coding agents — [Claude Code](https://docs.anthropic.com/en/docs/claude-code), [Codex CLI](https://github.com/openai/codex), [Cursor](https://cursor.sh), [Antigravity CLI](https://github.com/google-deepmind/antigravity), and any host that speaks the [Agent Skills](https://agentskills.io) open spec.

## Plugins

| Plugin | What it does |
|--------|--------------|
| [`agent-skills`](plugins/agent-skills/) | Personal agent skills — TDD loop, handoff, design grilling, GitHub PR helpers, spec lifecycle, CLI design, etc. (26 skills) |
| [`engineering-principles`](plugins/engineering-principles/) | Cross-repo engineering principles (constitution, cloud-native, tech stack), patterns, behavioural rules for AI assistants, and skills (`apply-principles`, `architecture-review`, `enforce-principles`). |

## Install

### Claude Code

```text
/plugin marketplace add matthew-a-carr/ai-plugins
/plugin install engineering-principles@matthew-a-carr
/plugin install agent-skills@matthew-a-carr
```

Or wire it directly in `~/.claude/settings.json`:

```jsonc
{
  "extraKnownMarketplaces": {
    "matthew-a-carr": {
      "source": { "source": "github", "repo": "matthew-a-carr/ai-plugins" }
    }
  },
  "enabledPlugins": {
    "engineering-principles@matthew-a-carr": true,
    "agent-skills@matthew-a-carr": true
  }
}
```

### Codex CLI

```text
$skill-installer matthew-a-carr/ai-plugins/plugins/agent-skills
```

### Any agent (skills.sh)

```bash
npx skills@latest add matthew-a-carr/ai-plugins/plugins/agent-skills
```

## Repo structure

```
ai-plugins/
├── .claude-plugin/marketplace.json   # Plugin catalog (points to subdirs)
├── plugins/
│   ├── agent-skills/                 # Skills plugin (26 skills)
│   │   ├── .claude-plugin/plugin.json
│   │   └── skills/
│   └── engineering-principles/       # Principles plugin (docs + 3 skills)
│       ├── .claude-plugin/plugin.json
│       ├── principles/
│       ├── patterns/
│       ├── skills/
│       └── hooks/
├── .github/workflows/                # Unified CI
├── release-please-config.json        # Multi-package release-please
└── pyproject.toml                    # Python tooling (pre-commit, lint)
```

Attribution for forked skills: see [`plugins/agent-skills/skills/ATTRIBUTION.md`](plugins/agent-skills/skills/ATTRIBUTION.md).

## License

MIT — see [LICENSE](LICENSE).
