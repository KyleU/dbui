package util

import (
	"github.com/sirupsen/logrus"
	logrusadapter "logur.dev/adapter/logrus"
	"logur.dev/logur"
	"os"
)

func InitLogging(verbose bool) logur.LoggerFacade {
	logger := logrus.New()

	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{})

	if verbose {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	return logrusadapter.New(logger)
}
