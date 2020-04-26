package config

import (
	"fmt"
	"strings"

	"github.com/kyleu/dbui/internal/gen/queries"

	"emperror.dev/errors"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/kyleu/dbui/internal/app/conn"
	_ "github.com/mattn/go-sqlite3"
	"logur.dev/logur"
)

type Service struct {
	ProjectRegistry *ProjectRegistry
	configDB        *sqlx.DB
	logger          logur.LoggerFacade
}

func (s *Service) GetConnection(connArg string) (*sqlx.DB, int, error) {
	p, err := s.ProjectRegistry.Get(connArg)
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}
	db, elapsed, err := conn.Connect(p.Engine(), p.URL)
	return db, elapsed, errors.WithStack(errors.Wrap(err, "error connecting to database"))
}

func NewService(logger logur.LoggerFacade) (*Service, error) {
	pr := NewRegistry(logger)
	root, err := pr.Get(".config")
	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "error getting root project"))
	}
	db, _, err := conn.Connect(root.Engine(), root.URL)
	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "error opening config database"))
	}
	defer func() {
		_ = db.Close()
	}()

	svc := Service{ProjectRegistry: pr, configDB: db, logger: logger}

	err = initIfNeeded(db, logger)
	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "error initializing config database"))
	}

	err = pr.Refresh(db)
	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "error initializing project registry"))
	}

	logger.Debug("Config service started at [" + root.URL + "]")
	return &svc, nil
}

func initIfNeeded(db *sqlx.DB, logger logur.LoggerFacade) error {
	exec("burn-it-down", db, logger, func(sb *strings.Builder) { queries.ResetConfigDatabase(sb) })
	exec("create-table-project", db, logger, func(sb *strings.Builder) { queries.CreateTableProject(sb) })
	exec("insert-data-project", db, logger, func(sb *strings.Builder) { queries.InsertDataProject(sb) })
	return nil
}

func exec(name string, db *sqlx.DB, logger logur.LoggerFacade, f func(*strings.Builder)) {
	sb := &strings.Builder{}
	f(sb)
	result, err := conn.ExecuteNoTx(logger, db, conn.Adhoc(sb.String()))
	if err != nil {
		panic(errors.WithStack(err))
	}
	logger.Debug(fmt.Sprintf("Ran [%s] in [%vms], [%v] rows affected", name, result.Timing.Elapsed, result.RowsAffected))
}
