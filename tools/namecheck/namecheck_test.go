package namecheck_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matthew-a-carr/ai-plugins/tools/namecheck"
)

// scaffold creates a directory tree from a list of relative paths.
// Paths ending in "/" are directories; everything else is a file.
func scaffold(t *testing.T, root string, paths []string) {
	t.Helper()
	for _, rel := range paths {
		full := filepath.Join(root, rel)
		if rel[len(rel)-1] == '/' {
			if err := os.MkdirAll(full, 0o755); err != nil {
				t.Fatal(err)
			}
		} else {
			dir := filepath.Dir(full)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(full, []byte("# content\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestLowercaseDirectoriesPass(t *testing.T) {
	root := t.TempDir()
	scaffold(t, root, []string{
		"plugins/",
		"plugins/agent-skills/",
		"plugins/agent-skills/skills/",
		"plugins/agent-skills/skills/tdd/",
		"plugins/agent-skills/skills/tdd/SKILL.md",
	})

	violations := namecheck.Check(root)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestUppercaseDirectoryFails(t *testing.T) {
	root := t.TempDir()
	scaffold(t, root, []string{
		"plugins/",
		"plugins/BadName/",
		"plugins/BadName/readme.md",
	})

	violations := namecheck.Check(root)
	if len(violations) == 0 {
		t.Fatal("expected violations for uppercase directory")
	}
	found := false
	for _, v := range violations {
		if containsStr(v, "BadName") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected violation mentioning 'BadName', got %v", violations)
	}
}

func TestKebabCaseMarkdownPasses(t *testing.T) {
	root := t.TempDir()
	scaffold(t, root, []string{
		"patterns/",
		"patterns/circuit-breaker-setup.md",
		"patterns/event-driven-outbox.md",
	})

	violations := namecheck.Check(root)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestNonKebabMarkdownFails(t *testing.T) {
	root := t.TempDir()
	scaffold(t, root, []string{
		"patterns/",
		"patterns/My_Pattern.md",
	})

	violations := namecheck.Check(root)
	if len(violations) == 0 {
		t.Fatal("expected violation for non-kebab markdown")
	}
	found := false
	for _, v := range violations {
		if containsStr(v, "My_Pattern.md") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected violation mentioning 'My_Pattern.md', got %v", violations)
	}
}

func TestEcosystemFilesAreAllowedAnywhere(t *testing.T) {
	root := t.TempDir()
	scaffold(t, root, []string{
		"plugins/",
		"plugins/agent-skills/",
		"plugins/agent-skills/skills/",
		"plugins/agent-skills/skills/tdd/",
		"plugins/agent-skills/skills/tdd/SKILL.md",
		"plugins/agent-skills/README.md",
		"plugins/engineering-principles/",
		"plugins/engineering-principles/CONTRIBUTING.md",
		"plugins/engineering-principles/CHANGELOG.md",
	})

	violations := namecheck.Check(root)
	if len(violations) != 0 {
		t.Fatalf("expected no violations for ecosystem files, got %v", violations)
	}
}

func TestHiddenDirsAreAllowed(t *testing.T) {
	root := t.TempDir()
	scaffold(t, root, []string{
		".github/",
		".github/workflows/",
		".claude-plugin/",
		".claude-plugin/marketplace.json",
	})

	violations := namecheck.Check(root)
	if len(violations) != 0 {
		t.Fatalf("expected no violations for hidden dirs, got %v", violations)
	}
}

func TestSkipDirsAreIgnored(t *testing.T) {
	root := t.TempDir()
	scaffold(t, root, []string{
		"node_modules/",
		"node_modules/BadPackage/",
		"node_modules/BadPackage/BAD_FILE.md",
		".git/",
		".git/objects/",
	})

	violations := namecheck.Check(root)
	if len(violations) != 0 {
		t.Fatalf("expected no violations inside skip dirs, got %v", violations)
	}
}

func TestRootMarkdownKeepsCasing(t *testing.T) {
	root := t.TempDir()
	scaffold(t, root, []string{
		"README.md",
		"AGENTS.md",
		"CLAUDE.md",
		"CONTRIBUTING.md",
	})

	violations := namecheck.Check(root)
	if len(violations) != 0 {
		t.Fatalf("expected no violations for root markdown files, got %v", violations)
	}
}

func TestNonMarkdownFilesAreIgnored(t *testing.T) {
	root := t.TempDir()
	scaffold(t, root, []string{
		"plugins/",
		"plugins/agent-skills/",
		"plugins/agent-skills/go.mod",
		"plugins/agent-skills/go.sum",
		"plugins/agent-skills/scripts/",
		"plugins/agent-skills/scripts/validate_skills.py",
	})

	violations := namecheck.Check(root)
	if len(violations) != 0 {
		t.Fatalf("expected no violations for non-md files, got %v", violations)
	}
}

func TestMixedViolations(t *testing.T) {
	root := t.TempDir()
	scaffold(t, root, []string{
		"plugins/",
		"plugins/GoodPlugin/",          // uppercase dir — violation
		"patterns/",
		"patterns/Bad_Name.md",          // not kebab — violation
		"patterns/good-name.md",         // fine
		"plugins/ok-plugin/",
		"plugins/ok-plugin/SKILL.md",    // ecosystem file — fine
	})

	violations := namecheck.Check(root)
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d: %v", len(violations), violations)
	}
}

func TestReferencesAllowScreamingKebab(t *testing.T) {
	root := t.TempDir()
	scaffold(t, root, []string{
		"plugins/",
		"plugins/my-skill/",
		"plugins/my-skill/references/",
		"plugins/my-skill/references/LANGUAGE.md",
		"plugins/my-skill/references/ADR-FORMAT.md",
		"plugins/my-skill/references/interface-design.md",
	})

	violations := namecheck.Check(root)
	if len(violations) != 0 {
		t.Fatalf("expected no violations for SCREAMING-KEBAB in references, got %v", violations)
	}
}

func TestEvalsAreSkipped(t *testing.T) {
	root := t.TempDir()
	scaffold(t, root, []string{
		"plugins/",
		"plugins/my-skill/",
		"plugins/my-skill/evals/",
		"plugins/my-skill/evals/scenario-0/",
		"plugins/my-skill/evals/scenario-0/inputs/",
		"plugins/my-skill/evals/scenario-0/inputs/SPEC-042-notifications.md",
		"plugins/my-skill/evals/scenario-0/inputs/_template.md",
	})

	violations := namecheck.Check(root)
	if len(violations) != 0 {
		t.Fatalf("expected no violations in evals, got %v", violations)
	}
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
