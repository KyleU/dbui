package conn

import (
	"fmt"
	"time"

	"emperror.dev/errors"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"logur.dev/logur"

	"github.com/kyleu/dbui/internal/app/conn/results"
	"github.com/kyleu/dbui/internal/app/util"
)

func GetRows(conn string, input string) (*sqlx.DB, *sqlx.Rows, error) {
	url := urlForConn(util.GetConnection(conn))
	sqlText := util.GetSQL(input)

	connection, _, err := connect(url)
	if err != nil {
		return nil, nil, err
	}
	stmt, _, err := prepare(*connection, sqlText)
	if err != nil {
		return nil, nil, err
	}
	rows, _, err := run(*stmt)
	return connection, rows, err
}

func GetResult(logger logur.LoggerFacade, conn string, input string) (*results.ResultSet, error) {
	url := urlForConn(util.GetConnection(conn))
	sqlText := util.GetSQL(input)

	connection, connected, err := connect(url)
	defer func() {
		if connection != nil {
			_ = connection.Close()
		}
	}()

	if err != nil {
		return nil, err
	}
	stmt, prepared, err := prepare(*connection, sqlText)
	if err != nil {
		return nil, err
	}
	rows, elapsed, err := run(*stmt)
	if err != nil {
		return nil, err
	}
	return resultset(logger, sqlText, rows, connected, prepared, elapsed)
}

func urlForConn(conn string) string {
	return "postgres://127.0.0.1:5432/dbui?sslmode=disable"
}

func connect(url string) (*sqlx.DB, int64, error) {
	startNanos := time.Now().UnixNano()
	conn, err := sqlx.Connect("pgx", url)
	connected := (time.Now().UnixNano() - startNanos) / int64(time.Microsecond)
	if err != nil {
		return nil, connected, errors.WithStack(errors.Wrap(err, "Unable to connect to database"))
	}
	return conn, connected, nil
}

func prepare(conn sqlx.DB, sqlText string) (*sqlx.Stmt, int64, error) {
	startNanos := time.Now().UnixNano()
	stmt, err := conn.Preparex(sqlText)
	prepared := (time.Now().UnixNano() - startNanos) / int64(time.Microsecond)
	if err != nil {
		return nil, prepared, errors.WithStack(errors.Wrap(err, "Unable to prepare query"))
	}
	return stmt, prepared, nil
}

func run(stmt sqlx.Stmt) (*sqlx.Rows, int64, error) {
	startNanos := time.Now().UnixNano()
	rows, err := stmt.Queryx()
	elapsed := (time.Now().UnixNano() - startNanos) / int64(time.Microsecond)
	if err != nil {
		return rows, elapsed, errors.WithStack(errors.Wrap(err, "Unable to execute query"))
	}
	return rows, elapsed, nil
}

func resultset(
	logger logur.LoggerFacade, sqlText string, rows *sqlx.Rows,
	connected int64, prepared int64, elapsed int64) (*results.ResultSet, error) {
	rs := results.ResultSet{
		SQL: sqlText,
		Timing: results.ResultSetTiming{
			Connected: connected,
			Prepared:  prepared,
			Elapsed:   elapsed,
		},
	}

	fields := make([]results.Column, 0)
	data := make([][]string, 0)

	for rows.Next() {
		if len(fields) == 0 {
			types, err := rows.ColumnTypes()
			if err != nil {
				return &rs, errors.WithStack(errors.Wrap(err, "Unable to extract column types from rows"))
			}
			for i, col := range types {
				n := col.Name()
				if n == "?column?" {
					n = fmt.Sprintf("col%d", i+1)
				}
				t := results.FieldTypeForName(logger, col.Name(), col.DatabaseTypeName())
				nullable, _ := col.Nullable()
				fields = append(fields, results.Column{Name: n, T: t, Nullable: nullable})
			}
			rs.Columns = fields
		}

		values, err := rows.SliceScan()
		if err != nil {
			rs.Data = data
			return &rs, errors.WithStack(errors.Wrap(err, "Unable to extract values from rows"))
		}

		row := make([]string, len(values))
		for i, v := range values {
			row[i] = fmt.Sprintf("%v", v)
		}
		data = append(data, row)
	}

	rs.Data = data
	return &rs, nil
}
