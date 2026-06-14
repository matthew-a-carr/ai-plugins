// Command check-naming enforces directory and file naming conventions.
package main

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
	".git":          true,
	"node_modules":  true,
	".venv":         true,
	"venv":          true,
	"__pycache__":   true,
	".pytest_cache": true,
	".ruff_cache":   true,
	"dist":          true,
	"build":         true,
	".idea":         true,
	".vscode":       true,
	"evals":         true,
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

// check walks root and returns a list of naming convention violations.
func check(root string) []string {
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

		if d.IsDir() && skipDirs[name] {
			return filepath.SkipDir
		}

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
			if strings.HasPrefix(name, ".") {
				if allowedHiddenDirs[name] {
					return nil
				}
				return filepath.SkipDir
			}
			if !lowerDirRE.MatchString(name) {
				violations = append(violations, fmt.Sprintf("directory not lowercase: %s", rel))
			}
			return nil
		}

		if filepath.Ext(name) != ".md" {
			return nil
		}

		fileDir := filepath.Dir(path)
		if fileDir == root {
			return nil
		}

		if ecosystemFiles[name] {
			return nil
		}

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

func main() {
	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	violations := check(root)
	if len(violations) > 0 {
		fmt.Println("Naming convention violations:")
		for _, v := range violations {
			fmt.Printf("  - %s\n", v)
		}
		fmt.Println()
		fmt.Println("Fix by renaming. Directories must be lowercase; markdown files")
		fmt.Println("in subdirectories must be lowercase kebab-case (e.g. cloud-native.md).")
		fmt.Println("Ecosystem files keep their casing.")
		os.Exit(1)
	}
}
