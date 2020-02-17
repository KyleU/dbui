package output

import (
	"strings"

	"github.com/kyleu/dbui/internal/app/conn/results"
	"github.com/olekukonko/tablewriter"
)

func AsTable(rs *results.ResultSet) (string, error) {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	columnNames := make([]string, len(rs.Columns))
	for i, c := range rs.Columns {
		columnNames[i] = c.Name
	}
	table.SetHeader(columnNames)
	for _, row := range rs.Data {
		table.Append(row)
	}
	table.Render()
	return tableString.String(), nil
}
