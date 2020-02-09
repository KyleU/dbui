package conn

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/kyleu/dbui/internal/app/conn/output"
	"github.com/kyleu/dbui/internal/app/conn/query"
	"github.com/kyleu/dbui/internal/app/conn/results"
)

func Sandbox(conn string, in string, out string) error {
	result, err := runQuery(conn, in, out)
	if err != nil {
		return err
	}
	str, err := output.AsTable(result)
	if err != nil {
		return err
	}
	fmt.Println(str)
	return nil
}

func runQuery(conn string, in string, out string) (*results.ResultSet, error) {
	url := "postgres://127.0.0.1:5432/dbui"
	sql := query.ListDatabases
	return connTest(url, sql)
}

func connTest(url string, sql string) (*results.ResultSet, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	stmt, err := conn.Prepare(ctx, "test", sql)
	if err != nil {
		return nil, err
	}

	fields := make([]results.Column, len(stmt.Fields))
	for i, v := range stmt.Fields {
		fields[i] = results.Column{
			Name: fmt.Sprintf("%s", v.Name),
			T: results.FieldTypes["string"],
		}
	}

	rows, err := conn.Query(ctx, sql)
	if err != nil {
		return &results.ResultSet{Columns: fields}, err
	}

	result := [][]string{}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return &results.ResultSet{Columns: fields, Data: result}, err
		}

		row := make([]string, len(values))
		for i, v := range values {
			row[i] = fmt.Sprintf("%s", v)
		}
		result = append(result, row)
	}

	return &results.ResultSet{Columns: fields, Data: result}, nil
}
