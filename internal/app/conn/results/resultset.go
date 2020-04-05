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

func (rs *ResultSet) PrimaryKeys() []string {
	ret := make([]string, 0)
	for _, x := range rs.Columns {
		if x.PrimaryKey {
			ret = append(ret, x.Name)
		}
	}
	return ret
}

func (rs *ResultSet) Col(name string) *Column {
	for _, x := range rs.Columns {
		if x.Name == name {
			return &x
		}
	}
	return nil
}
