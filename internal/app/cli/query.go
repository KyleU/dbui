package cli

import (
	"fmt"
	"strings"

	"emperror.dev/errors"
	"github.com/kyleu/dbui/internal/app/conn"
	"github.com/kyleu/dbui/internal/app/util"
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
			info := InitApp(appName, version, commitHash)
			rs, err := conn.GetResult(info.Logger, util.GetConnection(connNameArg), util.GetSQL(inputArg))
			if err != nil {
				return errors.WithStack(errors.Wrap(err, "Error retrieving result"))
			}
			out, err := conn.OutputFor(rs, getFormat(outputArg))
			if err != nil {
				return errors.WithStack(errors.Wrap(err, "Error formatting query output"))
			}
			fmt.Println(out)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&connNameArg, "conn", "c", "", "connection name or url")
	flags.StringVarP(&inputArg, "input", "i", "", "SQL string or \"file://path/filename.sql\"")
	flags.StringVarP(&outputArg, "output", "o", "", "output format, one of [table, markdown, csv, json, xml]")

	return cmd
}

func getFormat(o string) string {
	switch strings.ToLower(o) {
	case "csv":
		return "csv"
	case "json":
		return "json"
	case "markdown":
		return "markdown"
	default:
		return "table"
	}
}
