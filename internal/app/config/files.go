package config

import (
	"path/filepath"

	"github.com/kirsle/configdir"
)

var cfgDir = ""

func ConfigPath(filename string) string {
	if cfgDir == "" {
		cfgDir = configdir.LocalConfig("dbui")
		_ = configdir.MakePath(cfgDir)
	}
	return filepath.Join(cfgDir, filename)
}
