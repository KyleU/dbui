package conn

import (
	"io/ioutil"
	"strings"

	"github.com/kyleu/dbui/internal/gen/queries"
)

func getConnection(arg string) string {
	if arg == "" {
		arg = "default"
	}
	return arg
}

func getSQL(in string) string {
	if len(in) == 0 {
		return "select 'specify a sql string or file:path/filename.sql' as instructions"
	}
	switch {
	case strings.HasPrefix(in, "named:"):
		sb := &strings.Builder{}
		qName := strings.TrimPrefix(in, "named:")

		if strings.HasPrefix(qName, "list-") {
			switch strings.TrimPrefix(qName, "list-") {
			case "columns":
				queries.ListColumns(sb)
			case "databases":
				queries.ListDatabases(sb)
			case "indexes":
				queries.ListIndexes(sb)
			case "tables":
				queries.ListTables(sb)
			}
		}
		if strings.HasPrefix(qName, "example-") {
			switch strings.TrimPrefix(qName, "example-") {
			case "simple":
				queries.ExampleSimple(sb)
			case "complex":
				queries.ExampleComplex(sb)
			}
		}
		if sb.Len() == 0 {
			return "select 'Cannot load named query [" + qName + "]' as error"
		}
		return sb.String()
	case strings.HasPrefix(in, "file:"):
		path := strings.TrimPrefix(in, "file:")
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return "select 'cannot read file [" + path + "] (" + err.Error() + ")' as error"
		}
		return string(bytes)
	default:
		return in
	}
}
