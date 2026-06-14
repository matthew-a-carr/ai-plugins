package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestCleanFileHasNoViolations(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "good.md", "# Hello\n\nSome content.\n")

	violations := checkFile(filepath.Join(dir, "good.md"))
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestTrailingWhitespaceDetected(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "bad.md", "# Hello  \n\nTrailing spaces.   \n")

	violations := checkFile(filepath.Join(dir, "bad.md"))
	if len(violations) == 0 {
		t.Fatal("expected violations for trailing whitespace")
	}
	for _, v := range violations {
		if !contains(v, "trailing whitespace") {
			t.Fatalf("expected 'trailing whitespace' in violation, got: %s", v)
		}
	}
}

func TestHardTabsDetected(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "bad.md", "# Hello\n\n\tindented with tab\n")

	violations := checkFile(filepath.Join(dir, "bad.md"))
	if len(violations) == 0 {
		t.Fatal("expected violations for hard tabs")
	}
	found := false
	for _, v := range violations {
		if contains(v, "hard tab") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected 'hard tab' in violations, got %v", violations)
	}
}

func TestMultipleBlankLinesDetected(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "bad.md", "# Hello\n\n\n\nToo many blanks.\n")

	violations := checkFile(filepath.Join(dir, "bad.md"))
	if len(violations) == 0 {
		t.Fatal("expected violations for multiple blank lines")
	}
	found := false
	for _, v := range violations {
		if contains(v, "consecutive blank lines") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected 'consecutive blank lines' in violations, got %v", violations)
	}
}

func TestSingleBlankLineIsAllowed(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "good.md", "# Hello\n\nParagraph one.\n\nParagraph two.\n")

	violations := checkFile(filepath.Join(dir, "good.md"))
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestChangelogIsSkipped(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "CHANGELOG.md", "# Changelog\n\n\n\nBad formatting.\t\n")

	violations := checkFile(filepath.Join(dir, "CHANGELOG.md"))
	if len(violations) != 0 {
		t.Fatalf("expected CHANGELOG to be skipped, got %v", violations)
	}
}

func TestEvalsAreSkipped(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "skills/tdd/evals/scenario-0/task.md", "bad\t\ttabs  \n\n\n\n")

	violations := checkFile(filepath.Join(dir, "skills/tdd/evals/scenario-0/task.md"))
	if len(violations) != 0 {
		t.Fatalf("expected evals to be skipped, got %v", violations)
	}
}

func TestMultipleViolationsOnOneLine(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "bad.md", "\tbad line with tab and trailing space \n")

	violations := checkFile(filepath.Join(dir, "bad.md"))
	if len(violations) < 2 {
		t.Fatalf("expected at least 2 violations, got %d: %v", len(violations), violations)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
