package controllers

import (
	"net/http"

	"emperror.dev/errors"

	"github.com/gorilla/mux"
	"github.com/kyleu/dbui/internal/app/util"
	"github.com/kyleu/dbui/internal/gen/templates"
)

var _sandboxes = []string{"gallery", "testbed"}

func SandboxList(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx util.RequestContext) (int, error) {
		ctx.Breadcrumbs = util.BreadcrumbsSimple(ctx.Route("sandbox"), "sandbox")
		return templates.SandboxList(_sandboxes, ctx, w)
	})
}

func SandboxForm(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	act(w, r, func(ctx util.RequestContext) (int, error) {
		if key == "testbed" {
			return 0, errors.WithStack(errors.New("error!"))
		}
		bc := util.Breadcrumb{Path: ctx.Route("sandbox.run", "key", key), Title: key}
		ctx.Breadcrumbs = append(util.BreadcrumbsSimple(ctx.Route("sandbox"), "sandbox"), bc)
		return templates.SandboxForm(key, ctx, w)
	})
}
