// Package mdcheck performs lightweight markdown hygiene checks.
//
// Rules enforced:
//   - No trailing whitespace (except inside fenced code blocks)
//   - No hard tabs (except inside fenced code blocks)
//   - No multiple consecutive blank lines
//
// Files named CHANGELOG.md and files under evals/ directories are skipped.
package mdcheck

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// skipFile returns true if the file should not be checked.
func skipFile(path string) bool {
	base := filepath.Base(path)
	if base == "CHANGELOG.md" {
		return true
	}
	for _, part := range strings.Split(filepath.ToSlash(path), "/") {
		if part == "evals" {
			return true
		}
	}
	return false
}

// CheckFile returns a list of violations for a single markdown file.
func CheckFile(path string) []string {
	if skipFile(path) {
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		return []string{fmt.Sprintf("%s: %v", path, err)}
	}
	defer f.Close()

	var violations []string
	scanner := bufio.NewScanner(f)
	lineNum := 0
	blankRun := 0
	inFence := false

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Track fenced code blocks — don't lint inside them.
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inFence = !inFence
		}

		// Blank line tracking (applies everywhere).
		if strings.TrimSpace(line) == "" {
			blankRun++
			if blankRun > 1 {
				violations = append(violations,
					fmt.Sprintf("%s:%d: consecutive blank lines", path, lineNum))
			}
			continue
		}
		blankRun = 0

		// Skip content checks inside fenced code blocks.
		if inFence {
			continue
		}

		if strings.ContainsRune(line, '\t') {
			violations = append(violations,
				fmt.Sprintf("%s:%d: hard tab", path, lineNum))
		}

		if len(line) > 0 && (line[len(line)-1] == ' ' || line[len(line)-1] == '\t') {
			violations = append(violations,
				fmt.Sprintf("%s:%d: trailing whitespace", path, lineNum))
		}
	}

	return violations
}
