// Package namecheck enforces the repo's naming conventions:
//   - All directories (except hidden/skip dirs) are lowercase.
//   - All markdown files under subdirectories are lowercase, kebab-case.
//   - Ecosystem-convention files keep their original casing anywhere.
package namecheck

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ecosystemFiles may keep their casing in any location.
var ecosystemFiles = map[string]bool{
	"README.md":       true,
	"AGENTS.md":       true,
	"CLAUDE.md":       true,
	"SKILL.md":        true,
	"LICENSE":         true,
	"LICENSE.md":      true,
	"CHANGELOG.md":    true,
	"CONTRIBUTING.md": true,
	"CONSTITUTION.md": true,
	"ATTRIBUTION.md":  true,
}

// skipDirs are never scanned.
var skipDirs = map[string]bool{
	".git":           true,
	"node_modules":   true,
	".venv":          true,
	"venv":           true,
	"__pycache__":    true,
	".pytest_cache":  true,
	".ruff_cache":    true,
	"dist":           true,
	"build":          true,
	".idea":          true,
	".vscode":        true,
	"evals":          true, // test fixture dirs — intentionally varied naming
}

// allowedHiddenDirs are hidden directories we tolerate.
var allowedHiddenDirs = map[string]bool{
	".github":       true,
	".claude-plugin": true,
	".tessl-plugin":  true,
}

var kebabRE = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*\.md$`)
var screamingKebabRE = regexp.MustCompile(`^[A-Z0-9]+(-[A-Z0-9]+)*\.md$`)
var lowerDirRE = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*$`)

// Check walks root and returns a list of naming convention violations.
func Check(root string) []string {
	var violations []string

	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		rel, _ := filepath.Rel(root, path)
		if rel == "." {
			return nil
		}
		name := d.Name()

		// Skip excluded trees.
		if d.IsDir() && skipDirs[name] {
			return filepath.SkipDir
		}

		// Check parts of the path for skip dirs (handles nested cases).
		parts := strings.Split(rel, string(filepath.Separator))
		for _, part := range parts {
			if skipDirs[part] {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		if d.IsDir() {
			// Allow known hidden dirs.
			if strings.HasPrefix(name, ".") {
				if allowedHiddenDirs[name] {
					return nil
				}
				// Other hidden dirs — skip silently (e.g. .ruff_cache).
				return filepath.SkipDir
			}
			if !lowerDirRE.MatchString(name) {
				violations = append(violations, fmt.Sprintf("directory not lowercase: %s", rel))
			}
			return nil
		}

		// File checks: only enforce naming on .md files outside the root.
		if filepath.Ext(name) != ".md" {
			return nil
		}

		// Root-level markdown files keep their casing.
		fileDir := filepath.Dir(path)
		if fileDir == root {
			return nil
		}

		// Ecosystem files are always allowed.
		if ecosystemFiles[name] {
			return nil
		}

		// references/ dirs allow SCREAMING-KEBAB (e.g. LANGUAGE.md, ADR-FORMAT.md).
		inReferences := false
		for _, part := range parts {
			if part == "references" {
				inReferences = true
				break
			}
		}
		if inReferences {
			if !kebabRE.MatchString(name) && !screamingKebabRE.MatchString(name) {
				violations = append(violations, fmt.Sprintf("markdown file not lowercase kebab-case or SCREAMING-KEBAB: %s", rel))
			}
		} else if !kebabRE.MatchString(name) {
			violations = append(violations, fmt.Sprintf("markdown file not lowercase kebab-case: %s", rel))
		}

		return nil
	})

	return violations
}
