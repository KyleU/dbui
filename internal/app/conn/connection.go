package conn

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/kyleu/dbui/internal/app/conn/output"
	"github.com/kyleu/dbui/internal/app/conn/results"
	"logur.dev/logur"
	"time"
)

func GetResult(logger logur.LoggerFacade, conn string, in string) (*results.ResultSet, error) {
	return runQuery(logger, urlForConn(conn), in)
}

func OutputFor(result *results.ResultSet, out string) (string, error) {
	switch out {
	case "string":
		return output.AsString(result)
	default:
		return output.AsTable(result)
	}
}

func urlForConn(conn string) string {
	return "postgres://127.0.0.1:5432/dbui"
}

func runQuery(logger logur.LoggerFacade, url string, sql string) (*results.ResultSet, error) {
	startNanos := time.Now().UnixNano()
	logger.Debug("Running sql query", map[string]interface{}{"sql": sql, "url": url})
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "Unable to connect to database"))
	}
	defer conn.Close(ctx)

	connected := (time.Now().UnixNano() - startNanos) / int64(time.Microsecond)

	stmt, err := conn.Prepare(ctx, "test", sql)
	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "Unable to prepare query"))
	}

	rs := results.ResultSet{
		Sql:     sql,
		Timing: results.ResultSetTiming {
			Connected: connected,
		},
	}

	fields := make([]results.Column, len(stmt.Fields))
	for i, v := range stmt.Fields {
		n := string(v.Name)
		if n == "?column?" {
			n = fmt.Sprintf("col%d", i+1)
		}
		fields[i] = results.Column{Name: n, T: typeFor(v, conn.ConnInfo())}
	}
	rs.Columns = fields

	startNanos = time.Now().UnixNano()
	rows, err := conn.Query(ctx, sql)
	rs.Timing.Elapsed = (time.Now().UnixNano() - startNanos) / int64(time.Microsecond)
	if err != nil {
		return &rs, errors.WithStack(errors.Wrap(err, "Unable to execute query"))
	}

	data := make([][]string, 0)

	for rows.Next() {
		values, err := rows.Values()
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

func typeFor(f pgproto3.FieldDescription, info *pgtype.ConnInfo) results.FieldType {
	r := results.TypeInvalid
	t, ok := info.DataTypeForOID(f.DataTypeOID)
	if ok {
		switch t.Name {
		case "bool":
			r = results.TypeBool
		case "timestamp", "timestamptz":
			r = results.TypeTime
		case "json":
			r = results.TypeJSON
		case "uuid":
			r = results.TypeUUID
		case "BYTES_TODO":
			r = results.TypeBytes
		case "ENUM_TODO":
			r = results.TypeEnum
		case "text", "varchar", "bpchar":
			r = results.TypeString
		case "int2":
			r = results.TypeInt16
		case "int4":
			r = results.TypeInt32
		case "int8":
			r = results.TypeInt64
		case "float4":
			r = results.TypeFloat32
		case "float8":
			r = results.TypeFloat64
		default:
			println(string(f.Name) + ": " + t.Name)
		}
	}

	return r
}
