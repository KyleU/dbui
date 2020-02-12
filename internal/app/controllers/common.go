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
		ctx.AppInfo.Logger.Warn("Error running action")
	}
	endNanos := time.Now().UnixNano()
	delta := (endNanos - startNanos) / int64(time.Millisecond)
	ctx.AppInfo.Logger.Debug(fmt.Sprintf("[%v %v] processed in [%vms]", req.Method, req.URL.Path, delta))
}

func redir(res http.ResponseWriter, req *http.Request, f func(util.RequestContext) string) {
	startNanos := time.Now().UnixNano()
	ctx := util.ExtractContext(req, "redirect")
	url := f(ctx)
	endNanos := time.Now().UnixNano()
	delta := (endNanos - startNanos) / int64(time.Millisecond)
	ctx.AppInfo.Logger.Debug(fmt.Sprintf("[%v %v] processed in [%vms]", req.Method, req.URL.Path, delta))
	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusFound)
}

func saveSession(res http.ResponseWriter, req *http.Request, ctx util.RequestContext) {
	err := ctx.Session.Save(req, res)
	if err != nil {
		ctx.AppInfo.Logger.Warn("Unable to save session to response")
	}
}
