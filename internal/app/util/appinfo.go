package util

import (
	"emperror.dev/emperror"
	"github.com/kyleu/dbui/internal/app/config"
	"logur.dev/logur"
)

type AppInfo struct {
	AppName       string
	Debug         bool
	Version       string
	CommitHash    string
	CachePath     string
	Logger        logur.LoggerFacade
	ErrorHandler  emperror.ErrorHandlerFacade
	ConfigService *config.Service
}

func (a *AppInfo) Valid() bool {
	return a.ConfigService != nil
}
