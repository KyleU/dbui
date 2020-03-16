package output

import (
	"fmt"
	"strings"

	"github.com/kyleu/dbui/internal/app/conn/results"
)

func AsJson(rs *results.ResultSet) (string, error) {
	sb := &strings.Builder{}
	sb.WriteString("[\n")
	for _, row := range rs.Data {
		sb.WriteString("  {\n")
		for cIdx, cell := range row {
			sb.WriteString("    \"" + rs.Columns[cIdx].Name + "\": ")

			sb.WriteRune('"')
			sb.WriteString(fmt.Sprintf("%v", cell))
			sb.WriteRune('"')
			if cIdx < len(rs.Columns)-1 {
				sb.WriteString(",")
			}
			sb.WriteRune('\n')
		}
		sb.WriteString("  }\n")
	}
	sb.WriteString("]\n")
	return sb.String(), nil
}
