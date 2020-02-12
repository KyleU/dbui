package util

import (
	"emperror.dev/emperror"
	"logur.dev/logur"
)

type AppInfo struct {
	AppName      string
	Debug        bool
	Version      string
	CommitHash   string
	Logger       logur.LoggerFacade
	ErrorHandler emperror.ErrorHandlerFacade
}
