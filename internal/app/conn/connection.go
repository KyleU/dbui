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

func Connect(engine Engine, url string) (*sqlx.DB, int, error) {
	startNanos := time.Now().UnixNano()
	conn, err := sqlx.Open(engine.Key, url)
	connected := (time.Now().UnixNano() - startNanos) / int64(time.Microsecond)
	if err != nil {
		return nil, int(connected), errors.WithStack(errors.Wrap(err, "unable to connect to database"))
	}
	return conn, int(connected), nil
}

func Execute(logger logur.LoggerFacade, tx *sqlx.Tx, connectionMs int, q Query) (*results.StatementResult, error) {
	sqlText := getSQL(q.SQL)

	prepped, prepared, err := prepare(logger, *tx, sqlText)
	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "error preparing query for execution"))
	}
	result, elapsed, err := exec(*prepped, q.Values)
	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "error executing query"))
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

func ExecuteNoTx(logger logur.LoggerFacade, db *sqlx.DB, q Query) (*results.StatementResult, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "error opening transaction"))
	}
	res, err := Execute(logger, tx, 0, q)
	if err == nil {
		_ = tx.Commit()
	} else {
		_ = tx.Rollback()
	}
	return res, errors.WithStack(errors.Wrap(err, "error executing query"))
}

func GetRows(logger logur.LoggerFacade, tx *sqlx.Tx, q Query, rowFn func(rows *sqlx.Rows) error) (int, error) {
	sqlText := getSQL(q.SQL)

	res, _, err := prepare(logger, *tx, sqlText)
	if err != nil {
		return 0, errors.WithStack(errors.Wrap(err, "error opening transaction"))
	}
	rows, elapsed, err := query(*res, q.Values)
	if err != nil {
		return elapsed, errors.WithStack(errors.Wrap(err, "error running query"))
	}
	defer func() {
		_ = rows.Close()
		_ = tx.Commit()
	}()
	for rows.Next() {
		err = rowFn(rows)
		if err != nil {
			return 0, errors.WithStack(errors.Wrap(err, "error executing row function"))
		}
	}
	return elapsed, nil
}

func GetRowsNoTx(logger logur.LoggerFacade, db *sqlx.DB, q Query, rowFn func(rows *sqlx.Rows) error) (int, error) {
	tx, err := db.Beginx()
	if err != nil {
		return 0, errors.WithStack(errors.Wrap(err, "error opening transaction"))
	}
	res, err := GetRows(logger, tx, q, rowFn)
	if err == nil {
		_ = tx.Commit()
	} else {
		_ = tx.Rollback()
	}
	return res, errors.WithStack(errors.Wrap(err, "error getting rows from query"))
}

func RunQuery(logger logur.LoggerFacade, tx *sqlx.Tx, connectionMs int, q Query) (*results.ResultSet, error) {
	sqlText := getSQL(q.SQL)

	res, prepared, err := prepare(logger, *tx, sqlText)
	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "error preparing query"))
	}
	rows, elapsed, err := query(*res, q.Values)
	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "error executing query"))
	}
	defer func() {
		_ = rows.Close()
	}()
	return resultset(logger, sqlText, rows, connectionMs, prepared, elapsed)
}

func RunQueryNoTx(logger logur.LoggerFacade, db *sqlx.DB, connectionMs int, q Query) (*results.ResultSet, error) {
	tx, err := db.Beginx()
	if err != nil {
		logger.Warn(fmt.Sprintf("error opening database transaction: %+v", err))
		return nil, errors.WithStack(errors.Wrap(err, "error opening transaction"))
	}
	rs, err := RunQuery(logger, tx, connectionMs, q)
	if err == nil {
		_ = tx.Commit()
	} else {
		_ = tx.Rollback()
	}
	return rs, errors.WithStack(errors.Wrap(err, "error executing query"))
}

func prepare(logger logur.LoggerFacade, conn sqlx.Tx, sqlText string) (*sqlx.Stmt, int, error) {
	startNanos := time.Now().UnixNano()
	stmt, err := conn.Preparex(sqlText)
	prepared := (time.Now().UnixNano() - startNanos) / int64(time.Microsecond)
	if err != nil {
		logger.Warn(fmt.Sprintf("error preparing SQL [%v]", sqlText))
		return nil, int(prepared), errors.WithStack(errors.Wrap(err, "unable to prepare query"))
	}
	return stmt, int(prepared), nil
}

func query(stmt sqlx.Stmt, args []interface{}) (*sqlx.Rows, int, error) {
	startNanos := time.Now().UnixNano()
	rows, err := stmt.Queryx(args...)
	elapsed := (time.Now().UnixNano() - startNanos) / int64(time.Microsecond)
	if err != nil {
		return rows, int(elapsed), errors.WithStack(errors.Wrap(err, "unable to execute query"))
	}
	return rows, int(elapsed), nil
}

func exec(stmt sqlx.Stmt, args []interface{}) (sql.Result, int, error) {
	startNanos := time.Now().UnixNano()
	result, err := stmt.Exec(args...)
	elapsed := (time.Now().UnixNano() - startNanos) / int64(time.Microsecond)
	if err != nil {
		return result, int(elapsed), errors.WithStack(errors.Wrap(err, "unable to execute query"))
	}
	return result, int(elapsed), nil
}

func resultset(
		logger logur.LoggerFacade, sqlText string, rows *sqlx.Rows,
		connected int, prepared int, elapsed int,
) (*results.ResultSet, error) {
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
			return &rs, errors.WithStack(errors.Wrap(err, "unable to extract values from rows"))
		}

		if len(fields) == 0 {
			types, err := rows.ColumnTypes()
			if err != nil {
				return &rs, errors.WithStack(errors.Wrap(err, "unable to extract column types from rows"))
			}
			for i, col := range types {
				n := col.Name()
				if n == "?column?" {
					n = fmt.Sprintf("col%d", i+1)
				}
				dt := col.DatabaseTypeName()
				args := ""
				if strings.Contains(dt, "(") {
					args = dt[strings.Index(dt, "("):]
					dt = dt[0:strings.Index(dt, "(")]
					args = strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(args, ")"), "("))
				}
				t := results.FieldTypeForName(logger, col.Name(), dt)
				if t == results.TypeUnknown {
					if strings.HasPrefix(values[i].(string), "{") {
						t = results.TypeArrayUnknown
					} else if strings.Contains(values[i].(string), "\"=>\"") {
						t = results.TypeHStore
					}
				}
				nullable, _ := col.Nullable()
				p, s, _ := col.DecimalSize()
				l, _ := col.Length()
				if l > 1000000 {
					l = 0
				}
				if len(args) > 0 {
					p2, s2, l2 := results.ParseArgs(t, args)
					if p == 0 {
						p = p2
					}
					if s == 0 {
						s = s2
					}
					if l == 0 {
						l = l2
					}
				}
				fields = append(fields, results.Column{
					Name:       n,
					T:          t,
					Nullable:   nullable,
					PrimaryKey: false,
					Indexed:    false,
					Precision:  p,
					Scale:      s,
					Length:     l,
				})
			}
			rs.Columns = fields
		}

		data = append(data, values)
	}

	rs.Data = data
	return &rs, nil
}
