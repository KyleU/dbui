package output

import (
	"fmt"
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
		str := make([]string, len(row))
		for _, c := range row {
			str = append(str, fmt.Sprintf("%v", c))
		}
		table.Append(str)
	}
	table.Render()
	return tableString.String(), nil
}
