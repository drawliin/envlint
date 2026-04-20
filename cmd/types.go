package cmd

// options stores CLI flag values after Cobra parses them.
type options struct {
	path           string
	fix            bool
	json           bool
	strict         bool
	envFile        string
	exampleEnvFile string
}

// ExitError lets command code return a specific process exit code without calling os.Exit directly.
type ExitError struct {
	Code    int
	Message string
}
