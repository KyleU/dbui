package config

import (
	"fmt"

	"emperror.dev/errors"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/kyleu/dbui/internal/app/conn"
	_ "github.com/mattn/go-sqlite3"
	"logur.dev/logur"
)

type Service struct {
	Path   string
	db     *sqlx.DB
	logger logur.LoggerFacade
}

func (s *Service) GetConnection(connArg string) (*sqlx.DB, int, error) {
	engine := ""
	url := ""
	switch connArg {
	case "_root":
		engine = "sqlite3"
		url = ConfigPath(s.logger, "dbui.db")
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
	svc := Service{Path: path, db: db, logger: logger}

	err = initIfNeeded(svc)
	if err != nil {
		logger.Warn(fmt.Sprintf("Error initializing config database: %+v", err))
		return nil, err
	}

	logger.Debug("Config service started at [" + path + "]")
	return &svc, nil
}

func initIfNeeded(svc Service) error {
	rows, err := conn.GetRows(svc.db, "select 1")
	if err != nil {
		return err
	}
	for rows.Next() {
		svc.logger.Error("ROW!")
	}
	return nil
}
