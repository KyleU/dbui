package output

import (
	"fmt"
	"strings"

	"github.com/kyleu/dbui/internal/app/conn/results"
)

func AsString(rs *results.ResultSet) (string, error) {
	var sb strings.Builder
	for cIdx, col := range rs.Columns {
		if cIdx > 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(col.Name)
	}
	for _, row := range rs.Data {
		sb.WriteRune('\n')
		for cIdx, cell := range row {
			if cIdx > 0 {
				sb.WriteRune(',')
			}
			sb.WriteString(fmt.Sprintf("%v", cell))
		}
	}
	return sb.String(), nil
}
