package conn

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"emperror.dev/errors"
	"github.com/jmoiron/sqlx"
	"logur.dev/logur"

	"github.com/kyleu/dbui/internal/app/conn/results"
)

func Connect(engine string, url string) (*sqlx.DB, int, error) {
	startNanos := time.Now().UnixNano()
	conn, err := sqlx.Open(engine, url)
	connected := (time.Now().UnixNano() - startNanos) / int64(time.Microsecond)
	if err != nil {
		return nil, int(connected), errors.WithStack(errors.Wrap(err, "Unable to connect to database"))
	}
	return conn, int(connected), nil
}

func Execute(connection *sqlx.DB, connectionMs int, input string) (*results.StatementResult, error) {
	sqlText := getSQL(input)

	stmt, prepared, err := prepare(*connection, sqlText)
	if err != nil {
		return nil, err
	}
	result, elapsed, err := exec(*stmt)
	if err != nil {
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		affected = 0
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		lastId = 0
	}
	return &results.StatementResult{
		SQL:          sqlText,
		RowsAffected: affected,
		ReturnedId:   lastId,
		Timing: results.ResultSetTiming{
			Connected: connectionMs,
			Prepared:  prepared,
			Elapsed:   elapsed,
		},
	}, nil
}

func GetRows(connection *sqlx.DB, input string) (*sqlx.Rows, error) {
	sqlText := getSQL(input)

	stmt, _, err := prepare(*connection, sqlText)
	if err != nil {
		return nil, err
	}
	rows, _, err := query(*stmt)
	return rows, err
}

func GetResult(logger logur.LoggerFacade, connection *sqlx.DB, connectionMs int, input string) (*results.ResultSet, error) {
	sqlText := getSQL(input)

	stmt, prepared, err := prepare(*connection, sqlText)
	if err != nil {
		return nil, err
	}
	rows, elapsed, err := query(*stmt)
	if err != nil {
		return nil, err
	}
	return resultset(logger, sqlText, rows, connectionMs, prepared, elapsed)
}

func prepare(conn sqlx.DB, sqlText string) (*sqlx.Stmt, int, error) {
	startNanos := time.Now().UnixNano()
	stmt, err := conn.Preparex(sqlText)
	prepared := (time.Now().UnixNano() - startNanos) / int64(time.Microsecond)
	if err != nil {
		return nil, int(prepared), errors.WithStack(errors.Wrap(err, "Unable to prepare query"))
	}
	return stmt, int(prepared), nil
}

func query(stmt sqlx.Stmt) (*sqlx.Rows, int, error) {
	startNanos := time.Now().UnixNano()
	rows, err := stmt.Queryx()
	elapsed := (time.Now().UnixNano() - startNanos) / int64(time.Microsecond)
	if err != nil {
		return rows, int(elapsed), errors.WithStack(errors.Wrap(err, "Unable to execute query"))
	}
	return rows, int(elapsed), nil
}

func exec(stmt sqlx.Stmt) (sql.Result, int, error) {
	startNanos := time.Now().UnixNano()
	result, err := stmt.Exec()
	elapsed := (time.Now().UnixNano() - startNanos) / int64(time.Microsecond)
	if err != nil {
		return result, int(elapsed), errors.WithStack(errors.Wrap(err, "Unable to execute query"))
	}
	return result, int(elapsed), nil
}

func resultset(
		logger logur.LoggerFacade, sqlText string, rows *sqlx.Rows,
		connected int, prepared int, elapsed int) (*results.ResultSet, error) {
	rs := results.ResultSet{
		SQL: sqlText,
		Timing: results.ResultSetTiming{
			Connected: connected,
			Prepared:  prepared,
			Elapsed:   elapsed,
		},
	}

	fields := make([]results.Column, 0)
	data := make([][]interface{}, 0)

	for rows.Next() {
		values, err := rows.SliceScan()
		if err != nil {
			rs.Data = data
			return &rs, errors.WithStack(errors.Wrap(err, "Unable to extract values from rows"))
		}

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
				if t == results.TypeUnknown {
					if strings.HasPrefix(values[i].(string), "{") {
						t = results.TypeArrayUnknown
					} else if strings.Contains(values[i].(string), "\"=>\"") {
						t = results.TypeHStore
					}
				}
				nullable, _ := col.Nullable()
				fields = append(fields, results.Column{Name: n, T: t, Nullable: nullable})
			}
			rs.Columns = fields
		}

		data = append(data, values)
	}

	rs.Data = data
	return &rs, nil
}
