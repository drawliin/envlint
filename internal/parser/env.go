package parser

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

var envLinePattern = regexp.MustCompile(`^\s*(?:export\s+)?([A-Za-z_][A-Za-z0-9_]*)\s*=`)

// ParseEnvFile reads an env file and keeps both the parsed values and some extra metadata like duplicate keys.
func ParseEnvFile(path string) (File, error) {
	file := File{
		Path:   path,
		Values: map[string]string{},
	}

	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return file, nil
		}
		return file, fmt.Errorf("stat %s: %w", path, err)
	}
	if info.IsDir() {
		return file, fmt.Errorf("%s is a directory", path)
	}

	file.Exists = true

	values, err := godotenv.Read(path)
	if err != nil {
		return file, fmt.Errorf("parse %s: %w", path, err)
	}
	file.Values = values

	raw, err := os.Open(path)
	if err != nil {
		return file, fmt.Errorf("open %s: %w", path, err)
	}
	defer raw.Close()

	seen := map[string]bool{}
	duplicateSet := map[string]bool{}
	scanner := bufio.NewScanner(raw)
	// check duplicates
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		match := envLinePattern.FindStringSubmatch(line)
		if len(match) != 2 {
			continue
		}

		key := match[1]
		if seen[key] {
			duplicateSet[key] = true
			continue
		}
		seen[key] = true
	}

	if err := scanner.Err(); err != nil {
		return file, fmt.Errorf("scan %s: %w", path, err)
	}

	for key := range duplicateSet {
		file.Duplicates = append(file.Duplicates, key)
	}

	return file, nil
}

// EnvPaths builds the full paths for the main env file and the example env file from the chosen root.
func EnvPaths(root, envFile, exampleEnvFile string) (string, string) {
	return filepath.Join(root, envFile), filepath.Join(root, exampleEnvFile)
}
