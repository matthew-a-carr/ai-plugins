// Command validate-plugins discovers and validates all plugin.json manifests.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// requiredKeys are the keys every plugin.json must have as non-empty strings.
var requiredKeys = []string{"description", "name", "version"}

// validate checks a single plugin.json file and returns a list of problems.
func validate(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{fmt.Sprintf("missing %s", path)}
	}

	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return []string{fmt.Sprintf("%s is not valid JSON: %v", path, err)}
	}

	obj, ok := raw.(map[string]any)
	if !ok {
		return []string{fmt.Sprintf("%s must contain a JSON object, got %T", path, raw)}
	}

	var errs []string

	var missing []string
	for _, key := range requiredKeys {
		if _, exists := obj[key]; !exists {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		errs = append(errs, fmt.Sprintf("%s missing required keys: %v", path, missing))
	}

	for _, key := range requiredKeys {
		val, exists := obj[key]
		if !exists {
			continue
		}
		str, isStr := val.(string)
		if !isStr || strings.TrimSpace(str) == "" {
			errs = append(errs, fmt.Sprintf("%s: %q must be a non-empty string", path, key))
		}
	}

	return errs
}

// discoverPlugins finds all plugins/*/.claude-plugin/plugin.json under root.
func discoverPlugins(root string) []string {
	pattern := filepath.Join(root, "plugins", "*", ".claude-plugin", "plugin.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil
	}
	sort.Strings(matches)
	return matches
}

func main() {
	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	paths := os.Args[1:]
	if len(paths) == 0 {
		paths = discoverPlugins(root)
	}
	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, "no plugin.json files found")
		os.Exit(1)
	}

	failed := false
	for _, path := range paths {
		errs := validate(path)
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
