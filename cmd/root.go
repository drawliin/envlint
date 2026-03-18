package cmd

import (
	"fmt"
	"os"

	"env-doctor/internal/audit"
	"env-doctor/internal/report"

	"github.com/spf13/cobra"
)

type options struct {
	path   string
	fix    bool
	json   bool
	strict bool
}

type ExitError struct {
	Code    int
	Message string
}

func (e ExitError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("exit status %d", e.Code)
	}
	return e.Message
}

var opts options

var rootCmd = &cobra.Command{
	Use:   "env-doctor",
	Short: "Audit .env files against your codebase",
	Long:  "env-doctor audits environment variable usage, documentation drift, duplicate keys, and basic secret hygiene.",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := audit.Run(opts.path)
		if err != nil {
			return err
		}

		if opts.fix {
			if err := audit.ApplyFixes(result); err != nil {
				return err
			}
		}

		renderOpts := report.Options{
			JSON:   opts.json,
			Strict: opts.strict,
		}

		exitCode, outputErr := report.Write(os.Stdout, result, renderOpts)
		if outputErr != nil {
			return outputErr
		}

		if opts.strict && exitCode != 0 {
			return ExitError{Code: exitCode}
		}

		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	rootCmd.Flags().StringVar(&opts.path, "path", cwd, "path to audit")
	rootCmd.Flags().BoolVar(&opts.fix, "fix", false, "auto-add undocumented .env keys to .env.example with empty values")
	rootCmd.Flags().BoolVar(&opts.json, "json", false, "render output as JSON")
	rootCmd.Flags().BoolVar(&opts.strict, "strict", false, "exit with code 1 if any issues are found")
}

func Execute() error {
	return rootCmd.Execute()
}
