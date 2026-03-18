package audit

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"env-doctor/internal/parser"
	"env-doctor/internal/scanner"
)

type Result struct {
	Root                  string              `json:"root"`
	Env                   parser.File         `json:"env"`
	Example               parser.File         `json:"env_example"`
	ScannedFiles          []string            `json:"scanned_files"`
	Referenced            map[string][]string `json:"referenced"`
	MissingVars           []string            `json:"missing_vars"`
	UnusedVars            []string            `json:"unused_vars"`
	ExampleMissingFromEnv []string            `json:"example_missing_from_env"`
	UndocumentedInExample []string            `json:"undocumented_in_example"`
	DuplicateKeys         map[string][]string `json:"duplicate_keys"`
	GitignoreWarnings     []string            `json:"gitignore_warnings"`
	FixesApplied          []string            `json:"fixes_applied"`
	IssueCount            int                 `json:"issue_count"`
	BlockingIssueCount    int                 `json:"blocking_issue_count"`
	NonBlockingIssueCount int                 `json:"non_blocking_issue_count"`
}

func Run(root string) (*Result, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}

	envPath, examplePath := parser.EnvPaths(absRoot)
	envFile, err := parser.ParseEnvFile(envPath)
	if err != nil {
		return nil, err
	}

	exampleFile, err := parser.ParseEnvFile(examplePath)
	if err != nil {
		return nil, err
	}

	scanResult, err := scanner.Scan(absRoot)
	if err != nil {
		return nil, fmt.Errorf("scan source files: %w", err)
	}

	result := &Result{
		Root:          absRoot,
		Env:           envFile,
		Example:       exampleFile,
		ScannedFiles:  scanResult.Scanned,
		Referenced:    scanResult.Referenced,
		DuplicateKeys: map[string][]string{},
	}

	result.MissingVars = diffKeys(keysOf(scanResult.Referenced), envFile.Values)
	result.UnusedVars = diffKeys(keysOf(envFile.Values), scanResult.Referenced)
	result.ExampleMissingFromEnv = diffKeys(keysOf(exampleFile.Values), envFile.Values)
	result.UndocumentedInExample = diffKeys(keysOf(envFile.Values), exampleFile.Values)

	if len(envFile.Duplicates) > 0 {
		result.DuplicateKeys[envFile.Path] = append([]string(nil), envFile.Duplicates...)
	}
	if len(exampleFile.Duplicates) > 0 {
		result.DuplicateKeys[exampleFile.Path] = append([]string(nil), exampleFile.Duplicates...)
	}

	if warning := gitignoreWarning(absRoot, envFile.Exists); warning != "" {
		result.GitignoreWarnings = append(result.GitignoreWarnings, warning)
	}

	result.BlockingIssueCount = len(result.MissingVars) + len(result.ExampleMissingFromEnv) + duplicateCount(result.DuplicateKeys)
	result.NonBlockingIssueCount = len(result.UnusedVars) + len(result.UndocumentedInExample) + len(result.GitignoreWarnings)
	result.IssueCount = result.BlockingIssueCount + result.NonBlockingIssueCount

	return result, nil
}

func ApplyFixes(result *Result) error {
	if len(result.UndocumentedInExample) == 0 {
		return nil
	}

	if err := os.MkdirAll(result.Root, 0o755); err != nil {
		return fmt.Errorf("ensure root exists: %w", err)
	}

	f, err := os.OpenFile(result.Example.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open %s for fix: %w", result.Example.Path, err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("stat %s: %w", result.Example.Path, err)
	}

	if info.Size() > 0 {
		if _, err := f.WriteString("\n"); err != nil {
			return fmt.Errorf("prepare %s: %w", result.Example.Path, err)
		}
	}

	for _, key := range result.UndocumentedInExample {
		line := fmt.Sprintf("%s=\n", key)
		if _, err := f.WriteString(line); err != nil {
			return fmt.Errorf("write fix for %s: %w", key, err)
		}
		result.FixesApplied = append(result.FixesApplied, key)
	}

	updated, err := parser.ParseEnvFile(result.Example.Path)
	if err != nil {
		return err
	}
	result.Example = updated
	result.UndocumentedInExample = diffKeys(keysOf(result.Env.Values), updated.Values)
	result.BlockingIssueCount = len(result.MissingVars) + len(result.ExampleMissingFromEnv) + duplicateCount(result.DuplicateKeys)
	result.NonBlockingIssueCount = len(result.UnusedVars) + len(result.UndocumentedInExample) + len(result.GitignoreWarnings)
	result.IssueCount = result.BlockingIssueCount + result.NonBlockingIssueCount

	return nil
}

func diffKeys[T any](keys []string, compare map[string]T) []string {
	var missing []string
	for _, key := range keys {
		if _, ok := compare[key]; !ok {
			missing = append(missing, key)
		}
	}
	sort.Strings(missing)
	return missing
}

func keysOf[K ~string, V any](values map[K]V) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, string(key))
	}
	sort.Strings(keys)
	return keys
}

func duplicateCount(duplicates map[string][]string) int {
	total := 0
	for _, keys := range duplicates {
		total += len(keys)
	}
	return total
}

func gitignoreWarning(root string, envExists bool) string {
	if !envExists {
		return ""
	}

	gitignorePath := filepath.Join(root, ".gitignore")
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			return ".env exists but no .gitignore file was found"
		}
		return fmt.Sprintf("unable to read %s: %v", gitignorePath, err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if trimmed == ".env" || trimmed == "/.env" || trimmed == "*.env" || trimmed == ".env*" {
			return ""
		}
	}

	return ".env is not ignored by .gitignore"
}
