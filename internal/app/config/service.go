package config

import (
	"fmt"
	"github.com/kyleu/dbui/internal/gen/queries"
	"strings"

	"emperror.dev/errors"
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
	p := s.ProjectRegistry.Get(connArg)
	db, elapsed, err := conn.Connect(p.EngineString, p.Url)
	return db, elapsed, errors.WithStack(errors.Wrap(err, "Error connecting to database"))
}

func NewService(logger logur.LoggerFacade) (*Service, error) {
	path := ConfigPath(logger, "dbui.db")
	db, _, err := conn.Connect("sqlite3", path)
	if err != nil {
		logger.Warn(fmt.Sprintf("Error opening config database: %+v", err))
		return nil, err
	}
	defer func() {
		_ = db.Close()
	}()

	pr := NewRegistry(logger, db)
	svc := Service{ProjectRegistry: pr, configDB: db, logger: logger}

	err = initIfNeeded(db, logger)
	if err != nil {
		logger.Warn(fmt.Sprintf("Error initializing config database: %+v", err))
		return nil, err
	}

	err = pr.Refresh()
	if err != nil {
		logger.Warn(fmt.Sprintf("Error initializing project registry: %+v", err))
		return nil, err
	}

	logger.Debug("Config service started at [" + path + "]")
	return &svc, nil
}

func initIfNeeded(db *sqlx.DB, logger logur.LoggerFacade) error {
	tx, err := db.Beginx()
	if err != nil {
		logger.Warn(fmt.Sprintf("Error opening config database transaction: %+v", err))
		return err
	}
	exec("burn-it-down", tx, logger, func(sb *strings.Builder) { queries.ResetConfigDatabase(sb) })
	exec("create-table-project", tx, logger, func(sb *strings.Builder) { queries.CreateTableProject(sb) })
	exec("insert-data-project", tx, logger, func(sb *strings.Builder) { queries.InsertDataProject(sb) })
	err = tx.Commit()
	if err != nil {
		logger.Warn(fmt.Sprintf("Error comitting config database transaction: %+v", err))
		return err
	}
	return nil
}

func exec(name string, tx *sqlx.Tx, logger logur.LoggerFacade, f func(*strings.Builder)) {
	sb := &strings.Builder{}
	f(sb)
	result, err := conn.Execute(tx, 0, sb.String())
	if err != nil {
		panic(err)
	}
	logger.Debug(fmt.Sprintf("Ran [%s] in [%vms], [%v] rows affected", name, result.Timing.Elapsed, result.RowsAffected))
}
