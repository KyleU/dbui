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
	Path   string
	Logger logur.LoggerFacade
}

func (s *Service) GetConnection(connArg string) (*sqlx.DB, int, error) {
	engine := ""
	url := ""
	switch connArg {
	case "_root":
		engine = "sqlite3"
		url = ConfigPath(s.Logger, "dbui.db")
	case "test":
		engine = "pgx"
		url = "postgres://127.0.0.1:5432/dbui?sslmode=disable"
	default:
		return nil, 0, errors.WithStack(errors.New("Unknown database [" + connArg + "]"))
	}
	db, elapsed, err := conn.Connect(engine, url)
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
	svc := Service{Path: path, Logger: logger}

	err = initIfNeeded(db, svc.Logger)
	if err != nil {
		logger.Warn(fmt.Sprintf("Error initializing config database: %+v", err))
		return nil, err
	}

	logger.Debug("Config service started at [" + path + "]")
	return &svc, nil
}

func initIfNeeded(db *sqlx.DB, logger logur.LoggerFacade) error {
	exec("burn-it-down", db, logger, func(sb *strings.Builder) { queries.ResetConfigDatabase(sb) })
	exec("create-table-project", db, logger, func(sb *strings.Builder) { queries.CreateTableProject(sb) })
	return nil
}

func exec(name string, db *sqlx.DB, logger logur.LoggerFacade, f func(*strings.Builder)) {
	sb := &strings.Builder{}
	f(sb)
	result, err := conn.Execute(db, 0, sb.String())
	if err != nil {
		panic(err)
	}
	logger.Debug(fmt.Sprintf("Ran [%s] in [%vms], [%v] rows affected", name, result.Timing.Elapsed, result.RowsAffected))
}
