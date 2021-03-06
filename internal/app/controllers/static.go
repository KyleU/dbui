package controllers

import (
	"net/http"
	"path/filepath"
	"strings"

	"emperror.dev/emperror"
	"emperror.dev/errors"
	"github.com/kyleu/dbui/internal/app/controllers/assets"
)

func Favicon(w http.ResponseWriter, r *http.Request) {
	data, hash, contentType, err := assets.Asset("web/assets", "/favicon.ico")
	zipResponse(w, r, data, hash, contentType, errors.WithStack(err))
}

func Static(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(strings.TrimPrefix(r.URL.Path, "/assets"))
	if err == nil {
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		data, hash, contentType, err := assets.Asset("web/assets", path)
		zipResponse(w, r, data, hash, contentType, errors.WithStack(err))
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func zipResponse(w http.ResponseWriter, r *http.Request, data []byte, hash string, contentType string, err error) {
	if err == nil {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", contentType)
		w.Header().Add("Cache-Control", "public, max-age=31536000")
		w.Header().Add("ETag", hash)
		if r.Header.Get("If-None-Match") == hash {
			w.WriteHeader(http.StatusNotModified)
		} else {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write(data)
			emperror.Panic(errors.WithStack(errors.Wrap(err, "unable to write to response")))
		}
	} else {
		NotFound(w, r)
	}
}
