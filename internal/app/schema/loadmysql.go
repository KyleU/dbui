package schema

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"logur.dev/logur"

	"emperror.dev/errors"
	"github.com/kyleu/dbui/internal/app/conn"
	"github.com/kyleu/dbui/internal/app/conn/results"
)

func loadMySQL(logger logur.LoggerFacade, id string, connection *sqlx.DB) (map[string]Table, error) {
	var tables = map[string]Table{}

	_, err := conn.GetRowsNoTx(logger, connection, conn.Adhoc("named:list-columns-mysql"), func(rows *sqlx.Rows) error {
		var res MySQLColumnResult
		err := rows.StructScan(&res)
		if err != nil {
			return errors.WithStack(errors.Wrap(err, "error scanning column results from ["+id+"]"))
		}

		table, ok := tables[res.Table]
		if !ok {
			table = Table{Name: res.Table, ReadOnly: false}
		}
		t := results.FieldTypeForName(logger, res.Name, res.DataType)
		d := ""
		if res.Default.Valid {
			d = res.Default.String
		}
		table.AddColumn(results.Column{
			T:          t,
			Name:       res.Name,
			Nullable:   res.IsNullable(),
			PrimaryKey: false,
			Indexed:    false,
			Default:    d,
			Precision:  int64(res.NumericPrecision.Int32),
			Scale:      int64(res.NumericScale.Int32),
			Length:     int64(res.CharLength.Int32),
		})
		tables[table.Name] = table
		return nil
	})
	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "error retrieving columns from ["+id+"]"))
	}

	return tables, nil
}

type MySQLColumnResult struct {
	Schema            string         `db:"table_schema"`
	Table             string         `db:"table_name"`
	Name              string         `db:"column_name"`
	Ordinal           int32          `db:"ordinal_position"`
	Default           sql.NullString `db:"column_default"`
	Nullable          string         `db:"is_nullable"`
	DataType          string         `db:"data_type"`
	NumericPrecision  sql.NullInt32  `db:"numeric_precision"`
	NumericScale      sql.NullInt32  `db:"numeric_scale"`
	CharLength        sql.NullInt32  `db:"character_maximum_length"`
	DatetimePrecision sql.NullInt32  `db:"datetime_precision"`
}

func (cr *MySQLColumnResult) IsNullable() bool {
	return cr.Nullable == "YES"
}

type MySQLColumn struct {
	T sql.NullString `db:"t"`
}
