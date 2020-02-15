package results

type ResultSet struct {
	Sql string
	Columns []Column
	Data [][]string
	Timing ResultSetTiming
}

type ResultSetTiming struct {
	Connected int64
	Elapsed int64
}
