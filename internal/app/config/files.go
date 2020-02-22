package config

import (
	"logur.dev/logur"
	"path/filepath"

	"github.com/kirsle/configdir"
)

var cfgDir = ""

func ConfigPath(logger logur.LoggerFacade, filename string) string {
	if cfgDir == "" {
		cfgDir = configdir.LocalConfig("dbui")
		err := configdir.MakePath(cfgDir)
		if err != nil {
			logger.Error("Unable to make path [" + cfgDir + "]")
		}
	}
	return filepath.Join(cfgDir, filename)
}
