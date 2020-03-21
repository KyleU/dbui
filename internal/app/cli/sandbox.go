package cli

import (
	"fmt"

	"github.com/kyleu/dbui/internal/app/conn/output"

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
			info, err := InitApp(appName, version, commitHash)
			if err != nil {
				return errors.WithStack(errors.Wrap(err, "Error initializing application"))
			}
			connection, ms, err := info.ConfigService.GetConnection("")
			if err != nil {
				return errors.WithStack(errors.Wrap(err, "Error opening connection"))
			}
			rs, err := conn.GetResultNoTx(info.Logger, connection, ms, "")
			if err != nil {
				return errors.WithStack(errors.Wrap(err, "Error retrieving result"))
			}
			out, err := output.OutputFor(rs, "table")
			if err != nil {
				return errors.WithStack(errors.Wrap(err, "Error formatting output"))
			}
			fmt.Println(out)
			return nil
		},
	}

	return cmd
}
