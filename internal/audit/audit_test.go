package audit

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func writeAuditFixture(t *testing.T, root, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestRunReportsMissingDriftDuplicatesAndIgnoreWarning(t *testing.T) {
	root := t.TempDir()
	writeAuditFixture(t, root, ".env", "API_KEY=defined\nUNUSED_KEY=value\nDUP_KEY=first\nDUP_KEY=second\n")
	writeAuditFixture(t, root, ".env.example", "API_KEY=\nDOC_ONLY=\n")
	writeAuditFixture(t, root, "main.go", "package main\nimport \"os\"\nfunc main() { _ = os.Getenv(\"API_KEY\"); _ = os.Getenv(\"MISSING_KEY\") }\n")

	result, err := Run(root, Options{EnvFile: ".env", ExampleEnvFile: ".env.example"})
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(result.MissingVars, []string{"MISSING_KEY"}) {
		t.Fatalf("unexpected missing variables: %v", result.MissingVars)
	}
	if !slices.Equal(result.UnusedVars, []string{"DUP_KEY", "UNUSED_KEY"}) {
		t.Fatalf("unexpected unused variables: %v", result.UnusedVars)
	}
	if !slices.Equal(result.ExampleEnvMissingFromEnv, []string{"DOC_ONLY"}) {
		t.Fatalf("unexpected example-only variables: %v", result.ExampleEnvMissingFromEnv)
	}
	if !slices.Equal(result.UndocumentedInExampleEnv, []string{"DUP_KEY", "UNUSED_KEY"}) {
		t.Fatalf("unexpected undocumented variables: %v", result.UndocumentedInExampleEnv)
	}
	if !slices.Equal(result.DuplicateKeys[result.EnvFile.Path], []string{"DUP_KEY"}) {
		t.Fatalf("unexpected duplicates: %v", result.DuplicateKeys)
	}
	if len(result.GitignoreWarnings) != 1 {
		t.Fatalf("expected missing gitignore warning, got %v", result.GitignoreWarnings)
	}
	if result.BlockingIssueCount != 3 || result.NonBlockingIssueCount != 5 {
		t.Fatalf("unexpected issue counts: blocking=%d nonblocking=%d", result.BlockingIssueCount, result.NonBlockingIssueCount)
	}
}

func TestApplyFixesWritesEmptyValuesAndIsIdempotent(t *testing.T) {
	root := t.TempDir()
	writeAuditFixture(t, root, ".env", "DB_PASSWORD=do-not-copy-this\n")
	writeAuditFixture(t, root, ".gitignore", ".env\n")
	result, err := Run(root, Options{EnvFile: ".env", ExampleEnvFile: ".env.example"})
	if err != nil {
		t.Fatal(err)
	}

	if err := ApplyFixes(result); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(result.FixesApplied, []string{"DB_PASSWORD"}) {
		t.Fatalf("unexpected fixes: %v", result.FixesApplied)
	}
	if err := ApplyFixes(result); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(root, ".env.example"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "DB_PASSWORD=\n" || strings.Contains(string(got), "do-not-copy-this") {
		t.Fatalf("fix should write only an empty placeholder, got %q", got)
	}
}
