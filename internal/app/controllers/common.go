package controllers

import (
	"fmt"
	"github.com/kyleu/dbui/internal/app/util"
	"net/http"
	"time"
)

func act(w http.ResponseWriter, r *http.Request, title string, f func(util.RequestContext) (int, error)) {
	startNanos := time.Now().UnixNano()
	ctx := util.ExtractContext(r, title)

	if len(ctx.Flashes) > 0 {
		saveSession(w, r, ctx)
	}

	_, err := f(ctx)
	if err != nil {
		ctx.Logger.Warn("Error running action")
	}
	logComplete(startNanos, ctx, r)
}

func redir(w http.ResponseWriter, r *http.Request, f func(util.RequestContext) (string, error)) {
	startNanos := time.Now().UnixNano()
	ctx := util.ExtractContext(r, "redirect")
	url, err := f(ctx)
	if err != nil {
		ctx.Logger.Warn("Error running action")
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
	logComplete(startNanos, ctx, r)
}

func logComplete(startNanos int64, ctx util.RequestContext, r *http.Request) {
	delta := (time.Now().UnixNano() - startNanos) / int64(time.Microsecond)
	ms := util.MicrosToMillis(delta)
	ctx.Logger.Debug(fmt.Sprintf("[%v %v] processed in [%vms]", r.Method, r.URL.Path, ms), map[string]interface{} { "elapsed": delta })
}

func saveSession(w http.ResponseWriter, r *http.Request, ctx util.RequestContext) {
	err := ctx.Session.Save(r, w)
	if err != nil {
		ctx.Logger.Warn("Unable to save session to response")
	}
}
