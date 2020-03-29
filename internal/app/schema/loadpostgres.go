package schema

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"logur.dev/logur"

	"emperror.dev/errors"
	"github.com/kyleu/dbui/internal/app/conn"
	"github.com/kyleu/dbui/internal/app/conn/results"
)

func loadPostgres(logger logur.LoggerFacade, id string, connection *sqlx.DB) (map[string]Table, error) {
	tx, rows, err := conn.GetRowsNoTx(logger, connection, "named:list-columns")
	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "Error retrieving columns from ["+id+"]"))
	}
	defer func() {
		_ = rows.Close()
	}()

	var tables = map[string]Table{}
	for rows.Next() {
		var res ColumnResult
		err := rows.StructScan(&res)
		if err != nil {
			return nil, errors.WithStack(errors.Wrap(err, "Error scanning column results from ["+id+"]"))
		}

		table, ok := tables[res.Table]
		if !ok {
			table = Table{Name: res.Table, ReadOnly: res.Updatable == "No"}
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
		table.AddColumn(results.Column{
			T:         t,
			Name:      res.Name,
			Nullable:  res.IsNullable(),
			Default:   d,
			Precision: int64(res.NumericPrecision.Int32),
			Scale:     int64(res.NumericScale.Int32),
			Length:    int64(res.CharLength.Int32),
		})
		tables[table.Name] = table
	}

	err = tx.Commit()
	if err != nil {
		logger.Warn(fmt.Sprintf("Error comitting config database transaction: %+v", err))
	}

	return tables, nil
}

type ColumnResult struct {
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

func (cr *ColumnResult) IsNullable() bool {
	return cr.Nullable == "YES"
}
