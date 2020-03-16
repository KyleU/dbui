package config

import (
	"emperror.dev/emperror"
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
	ConfigService *Service
}

func (a *AppInfo) Valid() bool {
	return a.ConfigService != nil
}
