package schema

import (
	"emperror.dev/errors"
	"github.com/kyleu/dbui/internal/app/config"
	"github.com/kyleu/dbui/internal/app/conn"
)

var cache = map[string]Schema{}

func GetSchema(ai *config.AppInfo, id string, forceReload bool) (*Schema, error) {
	c, ok := cache[id]
	if ok && !forceReload {
		return &c, nil
	}
	c, err := LoadSchema(ai, id)
	if err != nil {
		return &c, errors.WithStack(errors.Wrap(err, "error loading schema"))
	}

	cache[id] = c
	return &c, nil
}

func LoadSchema(ai *config.AppInfo, id string) (Schema, error) {
	p := ai.ConfigService.ProjectRegistry.Get(id)
	s := NewSchema(id, p.Engine(), p.Title)
	connection, _, err := conn.Connect(p.Engine(), p.URL)
	if err != nil {
		return s, errors.WithStack(errors.Wrap(err, "error connecting to ["+id+"]"))
	}
	defer func() {
		if connection != nil {
			_ = connection.Close()
		}
	}()

	var tables = map[string]Table{}
	switch p.Engine() {
	case conn.PostgreSQL:
		tables, err = loadPostgres(ai.Logger, id, connection)
	case conn.MySQL:
		tables, err = loadMySQL(ai.Logger, id, connection)
	case conn.SQLite:
		tables, err = loadSqlite(ai.Logger, id, connection)
	}
	if err != nil {
		return s, errors.WithStack(errors.Wrap(err, "error loading columns from ["+id+"]"))
	}

	var ts []Table
	for _, table := range tables {
		ts = append(ts, table)
	}
	s.Tables.Add(ts...)

	return s, nil
}
