package schema

import (
	"database/sql"
	"github.com/kyleu/dbui/internal/app/config"

	"emperror.dev/errors"
	"github.com/kyleu/dbui/internal/app/conn"
	"github.com/kyleu/dbui/internal/app/conn/results"
)

var cache = map[string]Schema{}

func GetSchema(ai *config.AppInfo, id string, forceReload bool) (*Schema, error) {
	c, ok := cache[id]
	if ok && !forceReload {
		return &c, nil
	}
	c, err := LoadSchema(ai, id)
	if err != nil {
		return &c, err
	}

	cache[id] = c
	return &c, nil
}

func LoadSchema(ai *config.AppInfo, id string) (Schema, error) {
	s := NewSchema(id, "Test Schema")
	connection, _, err := ai.ConfigService.GetConnection(id)
	if err != nil {
		return s, errors.WithStack(errors.Wrap(err, "Error connecting to ["+id+"]"))
	}
	defer func() {
		if connection != nil {
			_ = connection.Close()
		}
	}()

	rows, err := conn.GetRows(connection, "named:list-columns")
	if err != nil {
		return s, errors.WithStack(errors.Wrap(err, "Error retrieving columns from ["+id+"]"))
	}
	defer func() {
		_ = rows.Close()
	}()

	var tables = map[string]Table{}
	for rows.Next() {
		var res ColumnResult
		err := rows.StructScan(&res)
		if err != nil {
			return s, errors.WithStack(errors.Wrap(err, "Error scanning column results from ["+id+"]"))
		}

		table, ok := tables[res.Table]
		if !ok {
			table = Table{Name: res.Table, ReadOnly: res.Updatable == "No"}
		}
		tn := res.UDTName
		if tn == "ARRAY" {
			tn = res.ArrayType.String
		}
		t := results.FieldTypeForName(ai.Logger, res.Name, tn)
		table.AddColumn(results.Column{T: t, Name: res.Name, Nullable: res.IsNullable()})
		tables[table.Name] = table
	}

	var ts []Table
	for _, table := range tables {
		ts = append(ts, table)
	}
	s.Tables.Add(ts...)

	return s, nil
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
