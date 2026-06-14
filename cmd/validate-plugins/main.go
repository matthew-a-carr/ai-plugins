// Command validate-plugins discovers and validates all plugin.json manifests.
package main

import (
	"fmt"
	"os"

	"github.com/matthew-a-carr/ai-plugins/tools/plugincheck"
)

func main() {
	root := plugincheck.RepoRoot()
	paths := os.Args[1:]
	if len(paths) == 0 {
		paths = plugincheck.DiscoverPlugins(root)
	}
	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, "no plugin.json files found")
		os.Exit(1)
	}

	failed := false
	for _, path := range paths {
		errs := plugincheck.Validate(path)
		if len(errs) > 0 {
			for _, e := range errs {
				fmt.Fprintln(os.Stderr, e)
			}
			failed = true
		} else {
			fmt.Printf("%s: ok\n", path)
		}
	}
	if failed {
		os.Exit(1)
	}
}
