package util

import (
	"github.com/sirupsen/logrus"
	logrusadapter "logur.dev/adapter/logrus"
	"logur.dev/logur"
	"os"
)

func InitLogging() logur.LoggerFacade {
	logger := logrus.New()

	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{})

	logger.SetLevel(logrus.InfoLevel)

	return logrusadapter.New(logger)
}
