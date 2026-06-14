package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTestFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestValidManifestPasses(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "plugin.json",
		`{"name":"engineering-principles","version":"1.0.0","description":"d"}`)

	errs := validate(path)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestExtraKeysAreAllowed(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "plugin.json",
		`{"name":"x","version":"1.0.0","description":"d","author":{"name":"x"},"license":"MIT"}`)

	errs := validate(path)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestMissingFileIsReported(t *testing.T) {
	errs := validate("/nonexistent/plugin.json")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
	if !containsStr(errs[0], "missing") {
		t.Fatalf("expected 'missing' in error, got: %s", errs[0])
	}
}

func TestInvalidJSONIsReported(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "plugin.json", `{not json`)

	errs := validate(path)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
	if !containsStr(errs[0], "not valid JSON") {
		t.Fatalf("expected 'not valid JSON' in error, got: %s", errs[0])
	}
}

func TestNonObjectTopLevelIsReported(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "plugin.json", `["a","list"]`)

	errs := validate(path)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
	if !containsStr(errs[0], "JSON object") {
		t.Fatalf("expected 'JSON object' in error, got: %s", errs[0])
	}
}

func TestMissingRequiredKeysAreNamed(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "plugin.json", `{"name":"x"}`)

	errs := validate(path)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
	if !containsStr(errs[0], "description") || !containsStr(errs[0], "version") {
		t.Fatalf("expected missing keys named in error, got: %s", errs[0])
	}
}

func TestEmptyStringValueIsReported(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "plugin.json",
		`{"name":"x","version":"1.0.0","description":"   "}`)

	errs := validate(path)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
	if !containsStr(errs[0], "description") {
		t.Fatalf("expected 'description' in error, got: %s", errs[0])
	}
}

func TestNonStringValueIsReported(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "plugin.json",
		`{"name":"x","version":2,"description":"d"}`)

	errs := validate(path)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
	if !containsStr(errs[0], "version") {
		t.Fatalf("expected 'version' in error, got: %s", errs[0])
	}
}

func TestDiscoverPlugins(t *testing.T) {
	root := t.TempDir()
	for _, rel := range []string{
		"plugins/alpha/.claude-plugin",
		"plugins/beta/.claude-plugin",
	} {
		if err := os.MkdirAll(filepath.Join(root, rel), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	writeTestFile(t, filepath.Join(root, "plugins/alpha/.claude-plugin"), "plugin.json",
		`{"name":"alpha","version":"1.0.0","description":"a"}`)
	writeTestFile(t, filepath.Join(root, "plugins/beta/.claude-plugin"), "plugin.json",
		`{"name":"beta","version":"1.0.0","description":"b"}`)

	paths := discoverPlugins(root)
	if len(paths) != 2 {
		t.Fatalf("expected 2 plugins, got %d: %v", len(paths), paths)
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
