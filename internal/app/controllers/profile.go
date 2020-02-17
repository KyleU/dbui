package controllers

import (
	"net/http"

	"github.com/kyleu/dbui/internal/app/util"
	template "github.com/kyleu/dbui/internal/gen/templates"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	act(w, r, "User Profile", func(ctx util.RequestContext) (int, error) {
		return template.Profile(ctx, w)
	})
}

func ProfileSave(w http.ResponseWriter, r *http.Request) {
	redir(w, r, func(ctx util.RequestContext) (string, error) {
		_ = r.ParseForm()
		util.SystemProfile.Name = r.Form.Get("username")
		util.SystemProfile.Theme = util.ThemeFromString(r.Form.Get("theme"))
		util.SystemProfile.NavColor = r.Form.Get("navbar-color")
		util.SystemProfile.LinkColor = r.Form.Get("link-color")
		ctx.Session.AddFlash("success:Profile saved")
		saveSession(w, r, ctx)
		return ctx.Route("home"), nil
	})
}
