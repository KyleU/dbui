package cli

import (
	"emperror.dev/emperror"
	"github.com/kyleu/dbui/internal/app/util"
	"github.com/spf13/cobra"
	"logur.dev/logur"
)

var verbose bool

// Configure configures a root command.
func Configure(version string, commitHash string) cobra.Command {
	rootCmd := cobra.Command{
		Use:   "dbui",
		Short: "Command line interface for dbui",
		Long:  "A work in progress...",
	}

	flags := rootCmd.PersistentFlags()
	flags.BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(
		NewServerCommand(rootCmd.Use, version, commitHash),
		NewQueryCommand(rootCmd.Use, version, commitHash),
		NewSandboxCommand(rootCmd.Use, version, commitHash),
		NewVersionCommand(),
	)

	return rootCmd
}

func InitApp(appName string, version string, commitHash string) util.AppInfo {
	logger := util.InitLogging(verbose)
	logger = logur.WithFields(logger, map[string]interface{}{"debug": verbose, "version": version, "commit": commitHash})

	handler := emperror.WithDetails(util.AppErrorHandler{ Logger: logger }, "key", "value")

	return util.AppInfo{
		AppName:    appName,
		Debug:      verbose,
		Version:    version,
		CommitHash: commitHash,
		Logger:     logger,
		ErrorHandler: handler,
	}
}
