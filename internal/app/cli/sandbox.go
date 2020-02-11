package cli

import (
	"github.com/kyleu/dbui/internal/app/conn"
	"github.com/spf13/cobra"
)

func NewSandboxCommand(appName string, version string, commitHash string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sandbox",
		Aliases: []string{"x"},
		Short:   "Runs an internal test",
		RunE: func(cmd *cobra.Command, _ []string) error {
			info := InitApp(appName, version, commitHash)
			return conn.GetResult(info, "", "", "")
		},
	}

	return cmd
}
