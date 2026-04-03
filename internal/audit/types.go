package audit

import "env-doctor/internal/parser"

type Result struct {
	Root                  string              `json:"root"`
	EnvFile               parser.File         `json:"env_file"`
	ExampleEnvFile        parser.File         `json:"example_env_file"`
	ScannedFiles          []string            `json:"scanned_files"`
	Referenced            map[string][]string `json:"referenced"`
	MissingVars           []string            `json:"missing_vars"`
	UnusedVars            []string            `json:"unused_vars"`
	ExampleEnvMissingFromEnv []string         `json:"example_env_missing_from_env"`
	UndocumentedInExampleEnv []string         `json:"undocumented_in_example_env"`
	DuplicateKeys         map[string][]string `json:"duplicate_keys"`
	GitignoreWarnings     []string            `json:"gitignore_warnings"`
	FixesApplied          []string            `json:"fixes_applied"`
	IssueCount            int                 `json:"issue_count"`
	BlockingIssueCount    int                 `json:"blocking_issue_count"`
	NonBlockingIssueCount int                 `json:"non_blocking_issue_count"`
}

type Options struct {
	EnvFile        string
	ExampleEnvFile string
}
