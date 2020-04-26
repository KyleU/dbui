package schema

import (
	"emperror.dev/errors"
	"github.com/kyleu/dbui/internal/app/config"
	"github.com/kyleu/dbui/internal/app/conn"
	"github.com/kyleu/dbui/internal/app/util"
	"logur.dev/logur"
	"strings"
)

var cache = map[string]Schema{}
var rootSchema Schema
var configSchema Schema

func GetSchema(ai *config.AppInfo, id string, forceReload bool) (*Schema, error) {
	if strings.HasSuffix(id, ".config") {
		if id == ".config" {
			return getRootSchema(ai.Logger)
		} else {
			return getConfigSchema(ai.Logger, strings.TrimSuffix(id, ".config"))
		}
	} else {
		c, ok := cache[id]
		if ok && !forceReload {
			return &c, nil
		}
		p, err := ai.ConfigService.ProjectRegistry.Get(id)
		if err != nil {
			return &c, errors.WithStack(errors.Wrap(err, "error loading project"))
		}
		c, err = loadSchema(ai.Logger, id, p.Engine(), p.Title, p.URL)
		if err != nil {
			return &c, errors.WithStack(errors.Wrap(err, "error loading schema"))
		}

		cache[id] = c
		return &c, nil
	}
}

func getRootSchema(logger logur.LoggerFacade) (*Schema, error) {
	if rootSchema.ID == "" {
		s, err := loadSchema(logger, ".config", conn.SQLite, "System Config", util.AppName + ".db")
		if err != nil {
			return &s, errors.WithStack(errors.Wrap(err, "error reading root schema"))
		}
		rootSchema = s
	}
	return &rootSchema, nil
}

func getConfigSchema(logger logur.LoggerFacade, id string) (*Schema, error) {
	if configSchema.ID == "" {
		s, err := loadSchema(logger, "_proto.config", conn.SQLite, "Config","_proto.config.db")
		if err != nil {
			return &s, errors.WithStack(errors.Wrap(err, "error reading config schema"))
		}
		configSchema = s
	}
	up := configSchema
	up.ID = id + ".config"
	up.Name = id + " Config"
	return &up, nil
}

func loadSchema(logger logur.LoggerFacade, id string, engine conn.Engine, title string, url string) (Schema, error) {
	s := NewSchema(id, engine, title)
	connection, _, err := conn.Connect(engine, url)
	if err != nil {
		return s, errors.WithStack(errors.Wrap(err, "error connecting to ["+id+"]"))
	}
	defer func() {
		if connection != nil {
			_ = connection.Close()
		}
	}()

	var tables = map[string]Table{}
	switch engine {
	case conn.PostgreSQL:
		tables, err = loadPostgres(logger, id, connection)
	case conn.MySQL:
		tables, err = loadMySQL(logger, id, connection)
	case conn.SQLite:
		tables, err = loadSqlite(logger, id, connection)
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
