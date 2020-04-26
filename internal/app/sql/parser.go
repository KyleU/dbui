package sql

import (
	"emperror.dev/errors"
	"fmt"
	"github.com/cube2222/octosql/parser/sqlparser"
)

func Parse(sql string) ([]string, error) {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "Error parsing sql"))
	}

	out := fmt.Sprint("%v", stmt)
	return append(make([]string, 0), out), nil
}
