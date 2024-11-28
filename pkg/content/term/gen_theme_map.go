//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func normalizeThemeName(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)
	// Replace spaces and underscores with hyphens
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")
	// Remove any non-alphanumeric characters except hyphens
	name = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' {
			return r
		}
		return -1
	}, name)
	// Replace multiple consecutive hyphens with a single hyphen
	for strings.Contains(name, "--") {
		name = strings.ReplaceAll(name, "--", "-")
	}
	// Trim hyphens from start and end
	return strings.Trim(name, "-")
}

func main() {
	// Read all .yml files in the themes directory
	entries, err := os.ReadDir("themes")
	if err != nil {
		panic(err)
	}

	// Create a map of normalized name to actual filename
	themeMap := make(map[string]string)
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yml") {
			name := strings.TrimSuffix(entry.Name(), ".yml")
			normalizedName := normalizeThemeName(name)
			themeMap[normalizedName] = entry.Name()
		}
	}

	// Get sorted keys for consistent output
	var keys []string
	for k := range themeMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Generate the Go code
	output := "package term\n\n"
	output += "// Code generated by gen_theme_map.go; DO NOT EDIT.\n\n"
	output += "var themeFileMap = map[string]string{\n"
	for _, k := range keys {
		output += fmt.Sprintf("\t%q: %q,\n", k, themeMap[k])
	}
	output += "}\n"

	// Write to themes_gen.go
	err = os.WriteFile("themes_gen.go", []byte(output), 0644)
	if err != nil {
		panic(err)
	}
}
