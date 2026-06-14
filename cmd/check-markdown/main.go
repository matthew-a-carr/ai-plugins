// check-markdown scans all .md files in the repository for hygiene issues.
//
// Rules: no trailing whitespace, no hard tabs, no consecutive blank lines.
// Skips: CHANGELOG.md, evals/ directories, fenced code blocks.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/matthew-a-carr/ai-plugins/tools/mdcheck"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var allViolations []string

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip hidden dirs, node_modules, .git.
		if info.IsDir() {
			base := filepath.Base(path)
			if strings.HasPrefix(base, ".") || base == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		violations := mdcheck.CheckFile(path)
		allViolations = append(allViolations, violations...)
		return nil
	})

	if len(allViolations) > 0 {
		for _, v := range allViolations {
			fmt.Println(v)
		}
		fmt.Fprintf(os.Stderr, "\n%d markdown violation(s) found\n", len(allViolations))
		os.Exit(1)
	}
}
