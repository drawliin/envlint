package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/drawliin/envlint/internal/audit"
)

func TestWriteUsesProvidedWriterForMissingFileWarnings(t *testing.T) {
	result := &audit.Result{
		IssueCount:            1,
		NonBlockingIssueCount: 1,
		UnusedVars:            []string{"UNUSED"},
	}
	var out bytes.Buffer
	code, err := Write(&out, result, Options{EnvFile: ".env", ExampleEnvFile: ".env.example"})
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Fatalf("non-strict run unexpectedly failed with %d", code)
	}
	for _, phrase := range []string{".env file missing", "env example file missing", "UNUSED", "Summary: 1 issues"} {
		if !strings.Contains(out.String(), phrase) {
			t.Fatalf("expected %q in supplied writer, got: %s", phrase, out.String())
		}
	}
}
