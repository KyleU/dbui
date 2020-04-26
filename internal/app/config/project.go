package config

import (
	"database/sql"
	"github.com/kyleu/dbui/internal/app/util"
	"sort"
	"strings"

	"emperror.dev/errors"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kyleu/dbui/internal/app/conn"
	"github.com/kyleu/dbui/internal/app/conn/results"
	"logur.dev/logur"
)

type Project struct {
	Key          string         `db:"key"`
	Title        string         `db:"title"`
	Description  string         `db:"description"`
	Owner        uuid.UUID      `db:"owner"`
	EngineString string         `db:"engine"`
	URL          string         `db:"url"`
	Username     sql.NullString `db:"username"`
	Password     sql.NullString `db:"password"`
}

func (p *Project) Engine() conn.Engine {
	return conn.EngineFromString(p.EngineString)
}

type ProjectRegistry struct {
	logger   logur.LoggerFacade
	names    []string
	projects map[string]*Project
}

func NewRegistry(logger logur.LoggerFacade) *ProjectRegistry {
	x := &ProjectRegistry{
		logger:   logger,
		names:    make([]string, 0),
		projects: make(map[string]*Project),
	}
	return x
}

func (s *ProjectRegistry) Refresh(db *sqlx.DB) error {
	s.projects = make(map[string]*Project)
	s.names = make([]string, 0)

	_, err := conn.GetRowsNoTx(s.logger, db, conn.Adhoc("select * from projects"), func(rows *sqlx.Rows) error {
		var res Project
		err := rows.StructScan(&res)
		if err != nil {
			return errors.WithStack(errors.Wrap(err, "error scanning project from config database"))
		}
		err = s.Add(false, &res)
		return errors.WithStack(errors.Wrap(err, "error adding projects to registry"))
	})
	if err != nil {
		return errors.WithStack(errors.Wrap(err, "error selecting projects from config database"))
	}

	return nil
}

func (s *ProjectRegistry) Names() []string {
	return s.names
}

func (s *ProjectRegistry) Get(key string) (*Project, error) {
	cfg := strings.HasSuffix(key, ".config")
	if cfg {
		key = strings.TrimSuffix(key, ".config")
	}
	ret, ok := s.projects[key]
	if len(key) > 0 && !ok {
		return nil, errors.WithStack(errors.New("cannot find project [" + key + "]"))
	}
	if cfg {
		title := key + " Configuration"
		desc := "Local SQLite database for settings and caches"
		url := ConfigPath(key + ".config.db")
		if len(key) == 0 {
			title = "System Configuration"
			desc = "SQLite database for system settings and configuration"
			url = ConfigPath(util.AppName + ".db")
		}
		ret = &Project{
			Key:          key + ".config",
			Title:        title,
			Description:  desc,
			EngineString: "sqlite3",
			URL:          url,
		}
	}
	return ret, nil
}

func (s *ProjectRegistry) Size() int {
	return len(s.names)
}

func (s *ProjectRegistry) Add(addToDb bool, t ...*Project) error {
	for _, proj := range t {
		if addToDb {
			_, err := update(s, proj.Key, proj)
			if err != nil {
				return errors.WithStack(errors.Wrap(err, "error updating project database"))
			}
		}
		s.projects[proj.Key] = proj
	}
	var acc []string
	for _, x := range s.projects {
		acc = append(acc, x.Key)
	}
	sort.Strings(acc)
	s.names = acc
	return nil
}

func update(s *ProjectRegistry, key string, proj *Project) (*results.StatementResult, error) {
	_, pExists := s.projects[key]

	root, err := s.Get(".config")
	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "error opening root database"))
	}

	db, _, err := conn.Connect(root.Engine(), root.URL)
	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "error opening config database"))
	}
	defer func() {
		_ = db.Close()
	}()

	values := []interface{}{
		proj.Key,
		proj.Title,
		proj.Description,
		proj.Owner,
		proj.EngineString,
		proj.URL,
		proj.Username,
		proj.Password,
	}

	if pExists {
		delete(s.projects, key)
		q := "update projects set key = ?, title = ?, description = ?, owner = ?, engine = ?, url = ?, username = ?, password = ? where key = ?"
		values = append(values, key)
		res, err := conn.ExecuteNoTx(s.logger, db, conn.Adhoc(q, values...))
		return res, errors.WithStack(errors.Wrap(err, "error updating project in config database"))
	} else {
		q := "insert into projects (key, title, description, owner, engine, url, username, password) values (?, ?, ?, ?, ?, ?, ?, ?)"
		res, err := conn.ExecuteNoTx(s.logger, db, conn.Adhoc(q, values...))
		return res, errors.WithStack(errors.Wrap(err, "error inserting project in config database"))
	}
}
