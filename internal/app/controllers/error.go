package controllers

import (
	"context"
	"fmt"
	"github.com/kyleu/dbui/internal/app/util"
	"io"
	"net/http"

	"github.com/kyleu/dbui/internal/app/config"
	"github.com/kyleu/dbui/internal/app/web"

	"emperror.dev/errors"

	"github.com/gorilla/mux"
	"github.com/kyleu/dbui/internal/gen/templates"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	ctx := web.ExtractContext(r)
	ctx.Title = "Not Found"
	ctx.Breadcrumbs = web.BreadcrumbsSimple(r.URL.Path, "not found")
	args := map[string]interface{}{"status": 500}
	ctx.Logger.Info(fmt.Sprintf("[%v %v] returned [%d]", r.Method, r.URL.Path, http.StatusNotFound), args)
	_, _ = templates.NotFound(r, ctx, w)
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func InternalServerError(router *mux.Router, info *config.AppInfo, w http.ResponseWriter, r *http.Request) {
	defer lastChanceError(w)

	if err := recover(); err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		rc := context.WithValue(r.Context(), routesKey, router)
		rc = context.WithValue(rc, infoKey, info)
		ctx := web.ExtractContext(r.WithContext(rc))
		ctx.Title = "Server Error"
		ctx.Breadcrumbs = web.BreadcrumbsSimple(r.URL.Path, "error")
			_, _ = templates.InternalServerError(util.GetErrorDetail(err.(error)), r, ctx, w)
		args := map[string]interface{}{"status": 500}
		st := http.StatusInternalServerError
		ctx.Logger.Warn(fmt.Sprintf("[%v %v] returned [%d]: %+v", r.Method, r.URL.Path, st, err), args)
	}
}

func lastChanceError(w io.Writer) {
	if err := recover(); err != nil {
		println(fmt.Sprintf("error while processing error handler: %+v", err))
		_, _ = w.Write([]byte("Internal Server Error"))
	}
}
