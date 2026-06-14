// Package plugincheck validates .claude-plugin/plugin.json manifests.
package plugincheck

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

// requiredKeys are the keys every plugin.json must have as non-empty strings.
var requiredKeys = []string{"description", "name", "version"}

// Validate checks a single plugin.json file and returns a list of problems.
// An empty slice means the manifest is valid.
func Validate(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{fmt.Sprintf("missing %s", path)}
	}

	// Must parse as JSON.
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return []string{fmt.Sprintf("%s is not valid JSON: %v", path, err)}
	}

	// Must be a JSON object.
	obj, ok := raw.(map[string]any)
	if !ok {
		return []string{fmt.Sprintf("%s must contain a JSON object, got %T", path, raw)}
	}

	var errs []string

	// Check for missing required keys.
	var missing []string
	for _, key := range requiredKeys {
		if _, exists := obj[key]; !exists {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		errs = append(errs, fmt.Sprintf("%s missing required keys: %v", path, missing))
	}

	// Check present required keys are non-empty strings.
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

// DiscoverPlugins finds all plugins/*/.claude-plugin/plugin.json under root.
func DiscoverPlugins(root string) []string {
	pattern := filepath.Join(root, "plugins", "*", ".claude-plugin", "plugin.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil
	}
	sort.Strings(matches)
	return matches
}

// RepoRoot returns the repository root by walking up from this source file.
func RepoRoot() string {
	_, thisFile, _, _ := runtime.Caller(0)
	// thisFile is tools/plugincheck/plugincheck.go — go up 3 levels.
	return filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))
}
