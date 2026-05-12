package scanner

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestScanSkipsExcludedDirectories(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "app.go", `package main

import "os"

func main() {
	_ = os.Getenv("ROOT_KEY")
}
`)

	excludedDirs := map[string]string{
		".venv":        "VENV_KEY",
		"node_modules": "NODE_MODULES_KEY",
		"vendor":       "VENDOR_KEY",
		"__pycache__":  "PYCACHE_KEY",
	}
	for dir, key := range excludedDirs {
		writeFile(t, root, filepath.Join(dir, "ignored.py"), "import os\nos.getenv(\""+key+"\")\n")
	}

	result, err := Scan(root)
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	if _, ok := result.Referenced["ROOT_KEY"]; !ok {
		t.Fatalf("expected ROOT_KEY to be referenced, got %#v", result.Referenced)
	}

	for dir, key := range excludedDirs {
		if _, ok := result.Referenced[key]; ok {
			t.Fatalf("expected %s from %s to be ignored, got %#v", key, dir, result.Referenced)
		}
		if slices.ContainsFunc(result.Scanned, func(path string) bool {
			return containsPathElement(path, dir)
		}) {
			t.Fatalf("expected scanner to skip directory %s, scanned %#v", dir, result.Scanned)
		}
	}
}

func writeFile(t *testing.T, root string, name string, content string) {
	t.Helper()

	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create parent directory for %s: %v", name, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func containsPathElement(path string, element string) bool {
	for {
		dir, file := filepath.Split(path)
		if file == element {
			return true
		}
		cleaned := filepath.Clean(dir)
		if cleaned == "." || cleaned == string(filepath.Separator) || cleaned == dir {
			return false
		}
		path = filepath.Clean(dir)
	}
}
