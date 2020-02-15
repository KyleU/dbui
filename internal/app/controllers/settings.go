package controllers

import (
	"github.com/kyleu/dbui/internal/app/util"
	template "github.com/kyleu/dbui/internal/gen/templates"
	"net/http"
)

func Settings(w http.ResponseWriter, r *http.Request) {
	act(w, r, "Settings", func(ctx util.RequestContext) (int, error) {
		return template.Settings(ctx, w)
	})
}

func SettingsSave(w http.ResponseWriter, r *http.Request) {
	redir(w, r, func(ctx util.RequestContext) (string, error) {
		ctx.Session.AddFlash("success:Settings saved")
		saveSession(w, r, ctx)
		return ctx.Route("settings"), nil
	})
}
