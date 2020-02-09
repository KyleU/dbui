package controllers

import (
	"github.com/kyleu/dbui/internal/app/controllers/assets"
	"github.com/kyleu/dbui/internal/app/util"
	"net/http"
	"path/filepath"
	"strings"
)

func Static(res http.ResponseWriter, req *http.Request) {
	path, err := filepath.Abs(strings.TrimPrefix(req.URL.Path, "/assets"))
	if err == nil {
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		data, hash, contentType, err := assets.Asset("web/assets", path)
		zipResponse(res, req, data, hash, contentType, err)
	} else {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
}

func Favicon(res http.ResponseWriter, req *http.Request) {
	data, hash, contentType, err := assets.Asset("web/assets", "/favicon.ico")
	zipResponse(res, req, data, hash, contentType, err)
}

func zipResponse(res http.ResponseWriter, req *http.Request, data []byte, hash string, contentType string, err error) {
	if err == nil {
		res.Header().Set("Content-Encoding", "gzip")
		res.Header().Set("Content-Type", contentType)
		res.Header().Add("Cache-Control", "public, max-age=31536000")
		res.Header().Add("ETag", hash)
		if req.Header.Get("If-None-Match") == hash {
			res.WriteHeader(http.StatusNotModified)
		} else {
			res.WriteHeader(http.StatusOK)
			_, err := res.Write(data)
			if err != nil {
				panic(err)
			}
		}
	} else {
		util.NotFound(res, req)
	}
}
