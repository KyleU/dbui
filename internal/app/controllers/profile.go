package controllers

import (
	"github.com/kyleu/dbui/internal/app/util"
	template "github.com/kyleu/dbui/internal/gen/templates"
	"net/http"
)

func Profile(res http.ResponseWriter, req *http.Request) {
	act(res, req, "", func(ctx util.RequestContext) (int, error) {
		return template.Profile(ctx, res)
	})
}

func ProfileSave(res http.ResponseWriter, req *http.Request) {
	redir(res, req, func(ctx util.RequestContext) (string, error) {
		_ = req.ParseForm()
		util.SystemProfile.Name = req.Form.Get("username")
		util.SystemProfile.Theme = util.ThemeFromString(req.Form.Get("theme"))
		util.SystemProfile.NavColor = req.Form.Get("navbar-color")
		util.SystemProfile.LinkColor = req.Form.Get("link-color")
		ctx.Session.AddFlash("success:Profile saved")
		saveSession(res, req, ctx)
		return ctx.Route("home"), nil
	})
}
