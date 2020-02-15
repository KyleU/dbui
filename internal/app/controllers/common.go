package controllers

import (
	"fmt"
	"github.com/kyleu/dbui/internal/app/util"
	"net/http"
	"time"
)

func act(res http.ResponseWriter, req *http.Request, title string, f func(util.RequestContext) (int, error)) {
	startNanos := time.Now().UnixNano()
	ctx := util.ExtractContext(req, title)

	if len(ctx.Flashes) > 0 {
		saveSession(res, req, ctx)
	}

	_, err := f(ctx)
	if err != nil {
		ctx.Logger.Warn("Error running action")
	}
	logComplete(startNanos, ctx, req)
}

func redir(res http.ResponseWriter, req *http.Request, f func(util.RequestContext) (string, error)) {
	startNanos := time.Now().UnixNano()
	ctx := util.ExtractContext(req, "redirect")
	url, err := f(ctx)
	if err != nil {
		ctx.Logger.Warn("Error running action")
	}
	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusFound)
	logComplete(startNanos, ctx, req)
}

func logComplete(startNanos int64, ctx util.RequestContext, req *http.Request) {
	delta := (time.Now().UnixNano() - startNanos) / int64(time.Microsecond)
	ms := util.MicrosToMillis(delta)
	ctx.Logger.Debug(fmt.Sprintf("[%v %v] processed in [%vms]", req.Method, req.URL.Path, ms), map[string]interface{} { "elapsed": delta })
}

func saveSession(res http.ResponseWriter, req *http.Request, ctx util.RequestContext) {
	err := ctx.Session.Save(req, res)
	if err != nil {
		ctx.Logger.Warn("Unable to save session to response")
	}
}
