package controllers

import (
	"fmt"
	"github.com/kyleu/dbui/internal/app/web"
	"net/http"
	"time"

	"github.com/kyleu/dbui/internal/gen/templates"
	"golang.org/x/text/language"

	"github.com/kyleu/dbui/internal/app/util"
)

func act(w http.ResponseWriter, r *http.Request, f func(web.RequestContext) (int, error)) {
	startNanos := time.Now().UnixNano()
	ctx := web.ExtractContext(r)

	if len(ctx.Flashes) > 0 {
		saveSession(w, r, ctx)
	}

	_, err := f(ctx)
	if err != nil {
		ctx.Logger.Warn(fmt.Sprintf("error running action: %+v", err))
		if ctx.Title == "" {
			ctx.Title = "Error"
		}
		_, _ = templates.InternalServerError(err.(error).Error(), err.(stackTracer).StackTrace(), r, ctx, w)
	}
	logComplete(startNanos, ctx, http.StatusOK, r)
}

func redir(w http.ResponseWriter, r *http.Request, f func(web.RequestContext) (string, error)) {
	startNanos := time.Now().UnixNano()
	ctx := web.ExtractContext(r)
	url, err := f(ctx)
	if err == nil {
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		logComplete(startNanos, ctx, http.StatusFound, r)
	} else {
		ctx.Logger.Warn(fmt.Sprintf("error running redirect: %+v", err))
		if ctx.Title == "" {
			ctx.Title = "Error"
		}
		_, _ = templates.InternalServerError(err.(error).Error(), err.(stackTracer).StackTrace(), r, ctx, w)
	}
}

func logComplete(startNanos int64, ctx web.RequestContext, status int, r *http.Request) {
	delta := (time.Now().UnixNano() - startNanos) / int64(time.Microsecond)
	ms := util.MicrosToMillis(language.AmericanEnglish, int(delta))
	args := map[string]interface{}{"elapsed": delta, "status": status}
	ctx.Logger.Debug(fmt.Sprintf("[%v %v] returned [%v] in [%v]", r.Method, r.URL.Path, status, ms), args)
}

func saveSession(w http.ResponseWriter, r *http.Request, ctx web.RequestContext) {
	err := ctx.Session.Save(r, w)
	if err != nil {
		ctx.Logger.Warn("Unable to save session to response")
	}
}
