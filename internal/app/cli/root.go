package cli

import (
	"github.com/spf13/cobra"
)

var verbose bool

// Configure configures a root command.
func Configure(rootCmd *cobra.Command, version string, commitHash string) {

	flags := rootCmd.PersistentFlags()
	flags.BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(
		NewServerCommand(rootCmd.Use, version, commitHash),
		NewQueryCommand(),
		NewSandboxCommand(),
		NewVersionCommand(),
	)
}
