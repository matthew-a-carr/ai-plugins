// Command check-naming enforces directory and file naming conventions.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/matthew-a-carr/ai-plugins/tools/namecheck"
)

func main() {
	_, thisFile, _, _ := runtime.Caller(0)
	root := filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))

	violations := namecheck.Check(root)
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
