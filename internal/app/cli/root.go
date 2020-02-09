package cli

import (
	"github.com/kyleu/dbui/internal/app/util"
	"github.com/spf13/cobra"
)


// Configure configures a root command.
func Configure(rootCmd *cobra.Command, version string, commitHash string) {
	var config string
	var verbose bool

	flags := rootCmd.PersistentFlags()
	flags.StringVarP(&config, "config", "c", ".", "Configuration directory, url, or file")
	flags.BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	info := util.AppInfo {
		Debug: verbose,
		Version: version,
		CommitHash: commitHash,
		ConfigDir: config,
	}

	rootCmd.AddCommand(
		NewServerCommand(info),
		NewQueryCommand(),
		NewSandboxCommand(),
		NewVersionCommand(),
	)
}
