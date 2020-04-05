package conn

import (
	"github.com/kyleu/dbui/internal/gen/queries"
	"io/ioutil"
	"strings"
)

type Query struct {
	Key string
	SQL string
	Values []interface{}
}

func Adhoc(sql string, values ...interface{}) Query {
	return Query{
		Key:    "adhoc",
		SQL:    sql,
		Values: values,
	}
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
			case "columns-postgres":
				queries.ListColumnsPostgres(sb)
			case "columns-mysql":
				queries.ListColumnsMySQL(sb)
			case "columns-sqlite":
				queries.ListColumnsSQLite(sb)
			case "indexes-postgres":
				queries.ListIndexesPostgres(sb)
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
