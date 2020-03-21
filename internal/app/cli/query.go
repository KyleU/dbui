package cli

import (
	"fmt"

	"emperror.dev/errors"
	"github.com/kyleu/dbui/internal/app/conn"
	"github.com/kyleu/dbui/internal/app/conn/output"
	"github.com/spf13/cobra"
)

func NewQueryCommand(appName string, version string, commitHash string) *cobra.Command {
	var connNameArg string
	var inputArg string
	var outputArg string

	cmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Runs the provided sql, displaying or saving the result",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			info, err := InitApp(appName, version, commitHash)
			if err != nil {
				return errors.WithStack(errors.Wrap(err, "Error initializing application"))
			}

			connection, ms, err := info.ConfigService.GetConnection(connNameArg)
			if err != nil {
				return errors.WithStack(errors.Wrap(err, "Error opening connection"))
			}
			rs, err := conn.GetResultNoTx(info.Logger, connection, ms, inputArg)
			if err != nil {
				return errors.WithStack(errors.Wrap(err, "Error retrieving result"))
			}
			out, err := output.OutputFor(rs, outputArg)
			if err != nil {
				return errors.WithStack(errors.Wrap(err, "Error formatting query output"))
			}
			fmt.Println(out)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&connNameArg, "conn", "c", "", "connection name or url")
	flags.StringVarP(&inputArg, "input", "i", "", "SQL string or \"file:path/filename.sql\"")
	flags.StringVarP(&outputArg, "output", "o", "", "output format, one of [table, markdown, csv, json, xml]")

	return cmd
}
