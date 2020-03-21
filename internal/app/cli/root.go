package cli

import (
	"emperror.dev/emperror"
	"emperror.dev/errors"
	logurhandler "emperror.dev/handler/logur"
	"github.com/kyleu/dbui/internal/app/config"
	"github.com/kyleu/dbui/internal/app/util"
	"github.com/spf13/cobra"
	"logur.dev/logur"
)

var verbose bool

const AppName = "dbui"

// Configure configures a root command.
func Configure(version string, commitHash string) cobra.Command {
	rootCmd := cobra.Command{
		Use:   AppName,
		Short: "Command line interface for dbui",
		Long:  "A work in progress...",
	}

	flags := rootCmd.PersistentFlags()
	flags.BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(
		NewServerCommand(rootCmd.Use, version, commitHash),
		NewQueryCommand(rootCmd.Use, version, commitHash),
		NewSandboxCommand(rootCmd.Use, version, commitHash),
		NewVersionCommand(version),
	)

	return rootCmd
}

func InitApp(appName string, version string, commitHash string) (*config.AppInfo, error) {
	logger := util.InitLogging(verbose)
	logger = logur.WithFields(logger, map[string]interface{}{"debug": verbose, "version": version, "commit": commitHash})

	errorHandler := logurhandler.New(logger)
	defer emperror.HandleRecover(errorHandler)

	handler := emperror.WithDetails(util.AppErrorHandler{Logger: logger}, "key", "value")

	cfg, err := config.NewService(logger)
	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "Error creating config service"))
	}

	ai := config.AppInfo{
		AppName:       appName,
		Debug:         verbose,
		Version:       version,
		CommitHash:    commitHash,
		Logger:        logger,
		ErrorHandler:  handler,
		ConfigService: cfg,
	}

	return &ai, nil
}
