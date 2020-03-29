package config

import (
	"database/sql"
	"emperror.dev/errors"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kyleu/dbui/internal/app/conn"
	"logur.dev/logur"
	"sort"
)

type Project struct {
	Key          string         `db:"key"`
	Title        string         `db:"title"`
	Description  string         `db:"description"`
	Owner        uuid.UUID      `db:"owner"`
	EngineString string         `db:"engine"`
	Url          string         `db:"url"`
	Username     sql.NullString `db:"username"`
	Password     sql.NullString `db:"password"`
}

func (p *Project) Engine() conn.Engine {
	return conn.EngineFromString(p.EngineString)
}

type ProjectRegistry struct {
	logger   logur.LoggerFacade
	db       *sqlx.DB
	names    []string
	projects map[string]Project
}

func NewRegistry(logger logur.LoggerFacade, db *sqlx.DB) *ProjectRegistry {
	return &ProjectRegistry{
		logger:   logger,
		db:       db,
		names:    make([]string, 0),
		projects: make(map[string]Project),
	}
}

func (s *ProjectRegistry) Refresh() error {
	tx, rows, err := conn.GetRowsNoTx(s.logger, s.db, "select * from projects")
	if err != nil {
		return errors.WithStack(errors.Wrap(err, "error selecting projects from config database"))
	}

	s.projects = make(map[string]Project)
	s.names = make([]string, 0)
	for rows.Next() {
		var res Project
		err := rows.StructScan(&res)
		if err != nil {
			return errors.WithStack(errors.Wrap(err, "error scanning project from config database"))
		}
		s.Add(res)
	}

	err = tx.Commit()
	if err != nil {
		return errors.WithStack(errors.Wrap(err, "error committing project transaction from config database"))
	}
	return nil
}

func (s *ProjectRegistry) Names() []string {
	return s.names
}

func (s *ProjectRegistry) Get(key string) Project {
	return s.projects[key]
}

func (s *ProjectRegistry) Size() int {
	return len(s.names)
}

func (s *ProjectRegistry) Add(t ...Project) {
	for _, x := range t {
		s.projects[x.Key] = x
	}
	var acc []string
	for _, x := range s.projects {
		acc = append(acc, x.Key)
	}
	sort.Strings(acc)
	s.names = acc
}
