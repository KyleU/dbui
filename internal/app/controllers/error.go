package controllers

import (
	"context"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kyleu/dbui/internal/app/util"
	template "github.com/kyleu/dbui/internal/gen/templates"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	ctx := util.ExtractContext(r, "Not Found")
	_, _ = template.NotFound(r, ctx, w)
}

func InternalServerError(router *mux.Router, info util.AppInfo, w http.ResponseWriter, r *http.Request) {
	defer lastChanceError(w)

	if err := recover(); err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		rc := context.WithValue(r.Context(), "routes", router)
		rc = context.WithValue(rc, "info", info)
		ctx := util.ExtractContext(r.WithContext(rc), "Internal Server Error")
		_, _ = template.InternalServerError(r, ctx, w)
	}
}

func lastChanceError(w io.Writer) {
	if err := recover(); err != nil {
		println("PANIC AT THE SERVER!!!!")
		_, _ = w.Write([]byte("Internal Server Error"))
	}
}
