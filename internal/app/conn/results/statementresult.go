package results

import (
	"fmt"
	"logur.dev/logur"
)

type StatementResult struct {
	SQL          string
	RowsAffected int64
	ReturnedId   int64
	Timing       ResultSetTiming
}

func (rs *StatementResult) Debug(l logur.LoggerFacade) {
	l.Debug(fmt.Sprintf("Returned statement result that affected [%d] rows", rs.RowsAffected))
}
