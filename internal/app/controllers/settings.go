package controllers

import (
	"net/http"

	"github.com/kyleu/dbui/internal/app/web"

	"github.com/kyleu/dbui/internal/gen/templates"
)

func Settings(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx web.RequestContext) (int, error) {
		ctx.Title = "Settings"
		ctx.Breadcrumbs = web.BreadcrumbsSimple(ctx.Route("settings"), "settings")
		return templates.Settings(ctx, w)
	})
}

func SettingsSave(w http.ResponseWriter, r *http.Request) {
	redir(w, r, func(ctx web.RequestContext) (string, error) {
		ctx.Session.AddFlash("success:Settings saved")
		saveSession(w, r, ctx)
		return ctx.Route("settings"), nil
	})
}
