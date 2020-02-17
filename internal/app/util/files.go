package util

import (
	"github.com/shibukawa/configdir"
	"path/filepath"
)

var cfg = configdir.New("dbui", "dbui")
var cachePath = cfg.QueryCacheFolder().Path

func ReadFile(info AppInfo, fn string) string {
	return filepath.Join(cachePath, fn)
}
