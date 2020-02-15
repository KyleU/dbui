package controllers

import (
	"github.com/kyleu/dbui/internal/app/util"
	template "github.com/kyleu/dbui/internal/gen/templates"
	"net/http"
)

func Settings(res http.ResponseWriter, req *http.Request) {
	act(res, req, "Sandbox List", func(ctx util.RequestContext) (int, error) {
		return template.Settings(ctx, res)
	})
}

func SettingsSave(res http.ResponseWriter, req *http.Request) {
	redir(res, req, func(ctx util.RequestContext) (string, error) {
		ctx.Session.AddFlash("success:Settings saved")
		saveSession(res, req, ctx)
		return ctx.Route("settings"), nil
	})
}
