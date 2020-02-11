package conn

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/kyleu/dbui/internal/app/conn/output"
	"github.com/kyleu/dbui/internal/app/conn/results"
	"github.com/kyleu/dbui/internal/app/util"
)

func GetResult(info util.AppInfo, conn string, in string, out string) error {
	url := "postgres://127.0.0.1:5432/dbui"
	result, err := runQuery(url, in)
	if err != nil {
		return err
	}
	str, err := output.AsTable(result)
	if err != nil {
		return errors.Wrap(err, "Unable to format output")
	}
	fmt.Println(str)
	return nil
}

func runQuery(url string, sql string) (*results.ResultSet, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to connect to database")
	}
	defer conn.Close(ctx)

	stmt, err := conn.Prepare(ctx, "test", sql)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to prepare query")
	}

	fields := make([]results.Column, len(stmt.Fields))
	for i, v := range stmt.Fields {
		fields[i] = results.Column{
			Name: fmt.Sprintf("%s", v.Name),
			T: results.TypeString,
		}
	}

	rows, err := conn.Query(ctx, sql)
	if err != nil {
		return &results.ResultSet{Columns: fields}, errors.Wrap(err, "Unable to execute query")
	}

	result := [][]string{}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return &results.ResultSet{Columns: fields, Data: result}, errors.Wrap(err, "Unable to extract values from rows")
		}

		row := make([]string, len(values))
		for i, v := range values {
			row[i] = fmt.Sprintf("%s", v)
		}
		result = append(result, row)
	}

	return &results.ResultSet{Columns: fields, Data: result}, nil
}
