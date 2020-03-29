package schema

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"logur.dev/logur"
	"strings"

	"emperror.dev/errors"
	"github.com/kyleu/dbui/internal/app/conn"
	"github.com/kyleu/dbui/internal/app/conn/results"
)

func loadSqlite(logger logur.LoggerFacade, id string, connection *sqlx.DB) (map[string]Table, error) {
	tx, rows, err := conn.GetRowsNoTx(logger, connection, "named:list-columns-sqlite")
	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "Error retrieving columns from ["+id+"]"))
	}
	defer func() {
		_ = rows.Close()
	}()

	var tables = map[string]Table{}
	for rows.Next() {
		var res SqliteColumnResult
		err := rows.StructScan(&res)
		if err != nil {
			return nil, errors.WithStack(errors.Wrap(err, "Error scanning column results from ["+id+"]"))
		}

		table, ok := tables[res.Table]
		if !ok {
			table = Table{Name: res.Table, ReadOnly: false}
		}
		dt := res.DataType
		args := ""
		if strings.Contains(dt, "(") {
			args = dt[strings.Index(dt, "("):]
			dt = dt[0:strings.Index(dt, "(")]
			args = strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(args, ")"), "("))
		}
		t := results.FieldTypeForName(logger, res.Name, dt)
		d := ""
		if res.Default.Valid {
			d = res.Default.String
		}
		p, s, l := results.ParseArgs(t, args)
		table.AddColumn(results.Column{
			T:        t,
			Name:     res.Name,
			Nullable: res.IsNullable(),
			Default:  d,
			Precision: p,
			Scale: s,
			Length: l,
		})
		tables[table.Name] = table
	}

	err = tx.Commit()
	if err != nil {
		logger.Warn(fmt.Sprintf("Error comitting config database transaction: %+v", err))
	}

	return tables, nil
}

type SqliteColumnResult struct {
	Table    string         `db:"table_name"`
	Name     string         `db:"column_name"`
	Ordinal  int32          `db:"ordinal_position"`
	Default  sql.NullString `db:"column_default"`
	NotNull  string         `db:"is_nullable"`
	DataType string         `db:"data_type"`
}

func (cr *SqliteColumnResult) IsNullable() bool {
	return cr.NotNull == "0"
}
