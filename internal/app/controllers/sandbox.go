package controllers

import (
	"net/http"

	"emperror.dev/emperror"
	"emperror.dev/errors"

	"github.com/gorilla/mux"
	"github.com/kyleu/dbui/internal/app/util"
	template "github.com/kyleu/dbui/internal/gen/templates"
)

var _sandboxes = []string{"gallery", "testbed"}

func SandboxList(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx util.RequestContext) (int, error) {
		ctx.Breadcrumbs = util.BreadcrumbsSimple(ctx.Route("sandbox"), "sandbox")
		return template.SandboxList(_sandboxes, ctx, w)
	})
}

func SandboxForm(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	if key == "testbed" {
		emperror.Panic(errors.WithStack(errors.New("error!")))
	}
	act(w, r, func(ctx util.RequestContext) (int, error) {
		bc := util.Breadcrumb{Path: ctx.Route("sandbox.run", "key", key), Title: key}
		ctx.Breadcrumbs = append(util.BreadcrumbsSimple(ctx.Route("sandbox"), "sandbox"), bc)
		return template.SandboxForm(key, ctx, w)
	})
}
