package cli

import (
	"fmt"

	"emperror.dev/errors"
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
			rs, err := conn.GetResult(info.Logger, "", "")
			if err != nil {
				return errors.WithStack(errors.Wrap(err, "Error retrieving result"))
			}
			out, err := conn.OutputFor(rs, "table")
			if err != nil {
				return errors.WithStack(errors.Wrap(err, "Error formatting output"))
			}
			fmt.Println(out)
			return nil
		},
	}

	return cmd
}
