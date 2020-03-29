package schema

import (
	"emperror.dev/errors"
	"github.com/kyleu/dbui/internal/app/config"
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

	var tables = map[string]Table{}
	if id == "_root" {
		tables, err = loadSqlite(ai.Logger, id, connection)
	} else {
		tables, err = loadPostgres(ai.Logger, id, connection)
	}
	if err != nil {
		return s, errors.WithStack(errors.Wrap(err, "Error loading columns from ["+id+"]"))
	}

	var ts []Table
	for _, table := range tables {
		ts = append(ts, table)
	}
	s.Tables.Add(ts...)

	return s, nil
}
