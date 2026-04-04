package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/drawliin/envlint/internal/audit"
)

const (
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
	colorBlue   = "\033[36m"
	colorReset  = "\033[0m"
)

func Write(w io.Writer, result *audit.Result, opts Options) (int, error) {
	if opts.JSON {
		return jsonExitCode(result, opts.Strict), writeJSON(w, result)
	}
	return terminalExitCode(result, opts.Strict), writeTerminal(w, result, opts)
}

func writeJSON(w io.Writer, result *audit.Result) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func writeTerminal(w io.Writer, result *audit.Result, opts Options) error {
	if result.IssueCount == 0 {
		_, err := fmt.Fprintf(w, "%s✅  No issues found%s\n", colorGreen, colorReset)
		return err
	}

	if !result.EnvFile.Exists {
		fmt.Printf("%s⚠️  .env file missing\n", colorYellow)
	}
	if !result.ExampleEnvFile.Exists {
		fmt.Printf("%s⚠️  env example file missing\n", colorYellow)
	}
	writeCategory(w, colorRed, "❌  Missing vars", result.MissingVars)
	writeCategory(w, colorYellow, "⚠️  Unused vars", result.UnusedVars)
	writeCategory(w, colorRed, fmt.Sprintf("❌  %s missing from %s", opts.ExampleEnvFile, opts.EnvFile), result.ExampleEnvMissingFromEnv)
	writeCategory(w, colorYellow, fmt.Sprintf("⚠️  %s missing from %s", opts.EnvFile, opts.ExampleEnvFile), result.UndocumentedInExampleEnv)
	writeCategory(w, colorYellow, "⚠️  Secret leak detection", result.GitignoreWarnings)

	if len(result.DuplicateKeys) > 0 {
		fmt.Fprintf(w, "%s❌  Duplicate keys%s\n", colorRed, colorReset)
		paths := make([]string, 0, len(result.DuplicateKeys))
		for path := range result.DuplicateKeys {
			paths = append(paths, path)
		}
		sort.Strings(paths)
		for _, path := range paths {
			fmt.Fprintf(w, "  - %s: %s\n", path, strings.Join(result.DuplicateKeys[path], ", "))
		}
		fmt.Fprintln(w)
	}

	if len(result.FixesApplied) > 0 {
		fmt.Fprintf(w, "%s%s%s\n", colorBlue, "Auto-fixes applied", colorReset)
		for _, key := range result.FixesApplied {
			fmt.Fprintf(w, "  - added %s= to %s\n", key, result.ExampleEnvFile.Path)
		}
		fmt.Fprintln(w)
	}

	_, err := fmt.Fprintf(
		w,
		"Summary: %d issues (%d blocking, %d non-blocking)\n",
		result.IssueCount,
		result.BlockingIssueCount,
		result.NonBlockingIssueCount,
	)
	return err
}

func writeCategory(w io.Writer, color, title string, items []string) {
	if len(items) == 0 {
		return
	}

	fmt.Fprintf(w, "%s%s%s\n", color, title, colorReset)
	for _, item := range items {
		fmt.Fprintf(w, "  - %s\n", item)
	}
	fmt.Fprintln(w)
}

func terminalExitCode(result *audit.Result, strict bool) int {
	if strict && result.IssueCount > 0 {
		return 1
	}
	return 0
}

func jsonExitCode(result *audit.Result, strict bool) int {
	return terminalExitCode(result, strict)
}
