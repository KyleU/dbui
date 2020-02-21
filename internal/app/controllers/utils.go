package controllers

import (
	"net/http"

	"github.com/kyleu/dbui/internal/app/util"
	template "github.com/kyleu/dbui/internal/gen/templates"
)

func Health(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx util.RequestContext) (int, error) {
		_, _ = w.Write([]byte("OK"))
		return 0, nil
	})
}

func Modules(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx util.RequestContext) (int, error) {
		bc := util.Breadcrumb{Path: ctx.Route("modules"), Title: "modules"}
		ctx.Breadcrumbs = append(util.BreadcrumbsSimple(ctx.Route("about"), "about"), bc)
		return template.ModulesList(ctx, w)
	})
}

func Routes(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx util.RequestContext) (int, error) {
		bc := util.Breadcrumb{Path: ctx.Route("routes"), Title: "routes"}
		ctx.Breadcrumbs = append(util.BreadcrumbsSimple(ctx.Route("about"), "about"), bc)
		return template.RoutesList(ctx, w)
	})
}
