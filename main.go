package main

import (
	"fmt"
	"os"

	"github.com/drawliin/envlint/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		var exitErr cmd.ExitError
		if ok := AsExitError(err, &exitErr); ok {
			if exitErr.Message != "" {
				fmt.Fprintln(os.Stderr, exitErr.Message)
			}
			os.Exit(exitErr.Code)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func AsExitError(err error, target *cmd.ExitError) bool {
	exitErr, ok := err.(cmd.ExitError)
	if !ok {
		return false
	}
	*target = exitErr
	return true
}
