package util

import (
	"path/filepath"

	"github.com/shibukawa/configdir"
)

var cfg = configdir.New("dbui", "dbui")
var cachePath = cfg.QueryCacheFolder().Path

func ReadFile(info AppInfo, fn string) string {
	return filepath.Join(cachePath, fn)
}
