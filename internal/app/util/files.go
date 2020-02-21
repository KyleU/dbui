package util

import (
	"os"
	"path/filepath"

	"github.com/shibukawa/configdir"
)

var cfg = configdir.New("dbui", "dbui")
var cachePath = cfg.QueryCacheFolder().Path

func FilePath(fn string) string {
	return filepath.Join(cachePath, fn)
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
