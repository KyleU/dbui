package results

import (
	"fmt"
	"strings"

	"logur.dev/logur"
)

type ResultSet struct {
	SQL     string
	Columns []Column
	Data    [][]string
	Timing  ResultSetTiming
}

type ResultSetTiming struct {
	Connected int64
	Prepared  int64
	Elapsed   int64
}

func (rs *ResultSet) Debug(l logur.LoggerFacade) {
	l.Debug(fmt.Sprintf("Returned ResultSet with [%d] columns and [%d] rows", len(rs.Columns), len(rs.Data)))
	var sb strings.Builder
	for idx, c := range rs.Columns {
		_, _ = sb.WriteString(c.Name + ": " + c.T.String())
		if idx < len(rs.Columns)-1 {
			_, _ = sb.WriteString(", ")
		}
	}
	l.Debug("Columns: " + sb.String())
}
