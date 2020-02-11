package cli

import (
	"github.com/kyleu/dbui/internal/app/conn"
	"github.com/kyleu/dbui/internal/app/queries"
	"github.com/spf13/cobra"
	"strings"
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
			return conn.GetResult(info, getConnection(connNameArg), getSql(inputArg), getFormat(outputArg))
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&connNameArg, "conn", "c", "", "connection name or url")
	flags.StringVarP(&inputArg, "input", "i", "", "SQL string, named query, or file path")
	flags.StringVarP(&outputArg, "output", "o", "", "output format, one of [table, markdown, json, csv]")

	return cmd
}

func getConnection(arg string) string {
	if arg == "" {
		arg = "default"
	}
	return arg
}

func getSql(in string) string {
	sb := &strings.Builder{}
	switch in {
	case "":
		sb.WriteString("select 'use --input to specify a sql string, named query, or sql file path' as instructions")
	case "list-columns":
		queries.ListColumns(sb)
	case "list-databases":
		queries.ListDatabases(sb)
	case "list-indexes":
		queries.ListIndexes(sb)
	case "list-tables":
		queries.ListTables(sb)
	default:
		sb.WriteString(in)
	}
	return sb.String()
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
