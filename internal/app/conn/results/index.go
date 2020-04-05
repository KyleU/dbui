package results

type Index struct {
	Table      string
	Index      string
	PrimaryKey bool
	Unique     bool
	Columns    []string
}
