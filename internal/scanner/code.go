package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var patterns = []pattern{
	{expression: regexp.MustCompile(`os\.Getenv\(\s*"([A-Z0-9_]+)"\s*\)`), group: 1},
	{expression: regexp.MustCompile(`os\.LookupEnv\(\s*"([A-Z0-9_]+)"\s*\)`), group: 1},
	{expression: regexp.MustCompile(`process\.env\.([A-Z0-9_]+)`), group: 1},
	{expression: regexp.MustCompile(`process\.env\[\s*["']([A-Z0-9_]+)["']\s*\]`), group: 1},
	{expression: regexp.MustCompile(`(?:os\.environ(?:\.get)?|environ\.get)\(\s*["']([A-Z0-9_]+)["']`), group: 1},
	{expression: regexp.MustCompile(`os\.getenv\(\s*["']([A-Z0-9_]+)["']\s*\)`), group: 1},
	{expression: regexp.MustCompile(`getenv\(\s*["']([A-Z0-9_]+)["']\s*\)`), group: 1},
}

func Scan(root string) (Result, error) {
	result := Result{
		Referenced: map[string][]string{},
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		name := d.Name()
		if d.IsDir() {
			if excludedDirectories[name] {
				return filepath.SkipDir
			}
			return nil
		}

		if !supportedExtensions[strings.ToLower(filepath.Ext(name))] {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		result.Scanned = append(result.Scanned, path)
		text := string(content)
		seenInFile := map[string]bool{}

		for _, p := range patterns {
			matches := p.expression.FindAllStringSubmatch(text, -1)
			for _, match := range matches {
				if len(match) <= p.group {
					continue
				}
				key := match[p.group]
				if seenInFile[key] {
					continue
				}
				result.Referenced[key] = append(result.Referenced[key], path)
				seenInFile[key] = true
			}
		}

		return nil
	})
	if err != nil {
		return result, err
	}

	sort.Strings(result.Scanned)
	for key := range result.Referenced {
		sort.Strings(result.Referenced[key])
	}

	return result, nil
}
