package cmd

type options struct {
	path           string
	fix            bool
	json           bool
	strict         bool
	envFile        string
	exampleEnvFile string
}

type ExitError struct {
	Code    int
	Message string
}
