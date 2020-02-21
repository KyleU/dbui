package schema

import (
	"database/sql"
	"fmt"
	"github.com/kyleu/dbui/internal/app/conn"
	"github.com/kyleu/dbui/internal/app/conn/results"
	"github.com/kyleu/dbui/internal/app/util"
)

var cache = map[string]Schema{}

func GetSchema(ai util.AppInfo, id string, forceReload bool) Schema {
	c, ok := cache[id]
	if ok && !forceReload {
		return c
	}
	c = LoadSchema(ai, id)
	cache[id] = c
	return c
}

func LoadSchema(ai util.AppInfo, id string) Schema {
	s := NewSchema(id, "Test Schema")

	connection, rows, err := conn.GetRows(id, "named:list-columns")
	defer func() {
		_ = rows.Close()
		if connection != nil {
			_ = connection.Close()
		}
	}()

	if err != nil {
		ai.Logger.Warn("Error loading schema")
		return s
	}

	var tables = map[string]Table{}
	var views = map[string]View{}
	for rows.Next() {
		var res ColumnResult
		err := rows.StructScan(&res)
		if err != nil {
			println(err.Error())
			ai.Logger.Warn(fmt.Sprintf("Error scanning column results: %v", err))
			return s
		}

		table, ok := tables[res.Table]
		if !ok {
			table = Table{Name: res.Table}
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

	var vs []View
	for _, view := range views {
		vs = append(vs, view)
	}
	s.Views.Add(vs...)

	return s
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
