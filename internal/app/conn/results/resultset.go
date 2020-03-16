package results

import (
	"fmt"
	"strings"

	"logur.dev/logur"
)

type ResultSet struct {
	SQL     string
	Columns []Column
	Data    [][]interface{}
	Timing  ResultSetTiming
}

type ResultSetTiming struct {
	Connected int
	Prepared  int
	Elapsed   int
}

func (rs *ResultSet) Debug(l logur.LoggerFacade) {
	l.Debug(fmt.Sprintf("Returned ResultSet with [%d] columns and [%d] rows", len(rs.Columns), len(rs.Data)))
	sb := &strings.Builder{}
	for idx, c := range rs.Columns {
		_, _ = sb.WriteString(c.Name + ": " + c.T.String())
		if idx < len(rs.Columns)-1 {
			_, _ = sb.WriteString(", ")
		}
	}
	l.Debug("Columns: " + sb.String())
}
