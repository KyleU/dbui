package schema

import (
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
	"logur.dev/logur"

	"emperror.dev/errors"
	"github.com/kyleu/dbui/internal/app/conn"
	"github.com/kyleu/dbui/internal/app/conn/results"
)

func loadPostgres(logger logur.LoggerFacade, id string, connection *sqlx.DB) (map[string]Table, error) {
	var tables = map[string]Table{}

	getTable := func(key string) Table {
		table, ok := tables[key]
		if !ok {
			table = Table{Name: key}
		}
		return table
	}

	_, err := conn.GetRowsNoTx(logger, connection, conn.Adhoc("named:list-columns-postgres"), func(rows *sqlx.Rows) error {
		var res PostgresColumnResult
		err := rows.StructScan(&res)
		if err != nil {
			return errors.WithStack(errors.Wrap(err, "error scanning column results from ["+id+"]"))
		}

		tn := res.UDTName
		if tn == "ARRAY" {
			tn = res.ArrayType.String
		}
		t := results.FieldTypeForName(logger, res.Name, tn)
		d := ""
		if res.Default.Valid {
			d = res.Default.String
		}
		table := getTable(res.Table)
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

	_, err = conn.GetRowsNoTx(logger, connection, conn.Adhoc("named:list-indexes-postgres"), func(rows *sqlx.Rows) error {
		var res PostgresIndexResult
		err := rows.StructScan(&res)
		if err != nil {
			return errors.WithStack(errors.Wrap(err, "error scanning index results from ["+id+"]"))
		}

		table := getTable(res.Table)
		table.AddIndex(results.Index{
			Table:      res.Table,
			Index:      res.Index,
			PrimaryKey: res.PrimaryKey,
			Unique:     res.Unique,
			Columns:    strings.Split(res.ColumnNames, ","),
		})
		tables[table.Name] = table
		return nil
	})
	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "error retrieving indexes from ["+id+"]"))
	}

	return tables, nil
}

type PostgresColumnResult struct {
	Schema                string         `db:"table_schema"`
	Table                 string         `db:"table_name"`
	Name                  string         `db:"column_name"`
	Ordinal               int32          `db:"ordinal_position"`
	Default               sql.NullString `db:"column_default"`
	Nullable              string         `db:"is_nullable"`
	DataType              string         `db:"data_type"`
	ArrayType             sql.NullString `db:"array_type"`
	CharLength            sql.NullInt32  `db:"character_maximum_length"`
	OctetLength           sql.NullInt32  `db:"character_octet_length"`
	NumericPrecision      sql.NullInt32  `db:"numeric_precision"`
	NumericPrecisionRadix sql.NullInt32  `db:"numeric_precision_radix"`
	NumericScale          sql.NullInt32  `db:"numeric_scale"`
	DatetimePrecision     sql.NullInt32  `db:"datetime_precision"`
	IntervalType          sql.NullInt32  `db:"interval_type"`
	DomainSchema          sql.NullString `db:"domain_schema"`
	DomainName            sql.NullString `db:"domain_name"`
	UDTSchema             string         `db:"udt_schema"`
	UDTName               string         `db:"udt_name"`
	DTDIdentifier         string         `db:"dtd_identifier"`
	Updatable             string         `db:"is_updatable"`
}

func (cr *PostgresColumnResult) IsNullable() bool {
	return cr.Nullable == "YES"
}

type PostgresIndexResult struct {
	Table       string `db:"table_name"`
	Index       string `db:"index_name"`
	PrimaryKey  bool   `db:"pk"`
	Unique      bool   `db:"u"`
	ColumnNames string `db:"column_names"`
}
