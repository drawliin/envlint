package parser

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestParseEnvFileDetectsDuplicates(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("# comment\nAPI_KEY=first\nAPI_KEY=second\nOTHER=ok\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := ParseEnvFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Exists || !slices.Equal(got.Duplicates, []string{"API_KEY"}) {
		t.Fatalf("expected an existing file with a duplicate API_KEY, got %+v", got)
	}
	if _, ok := got.Values["OTHER"]; !ok {
		t.Fatalf("expected OTHER in parsed values: %+v", got.Values)
	}
}

func TestParseEnvFileMissingIsNotFatal(t *testing.T) {
	got, err := ParseEnvFile(filepath.Join(t.TempDir(), ".env"))
	if err != nil || got.Exists || len(got.Values) != 0 {
		t.Fatalf("unexpected missing-file result: %+v, err: %v", got, err)
	}
}
