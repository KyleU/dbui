package controllers

import (
	"net/http"

	"github.com/kyleu/dbui/internal/app/util"
	"github.com/kyleu/dbui/internal/gen/templates"
)

func Health(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx util.RequestContext) (int, error) {
		_, _ = w.Write([]byte("OK"))
		return 0, nil
	})
}

func Modules(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx util.RequestContext) (int, error) {
		ctx.Title = "Modules"
		ctx.Breadcrumbs = append(aboutBC(ctx), util.Breadcrumb{Path: ctx.Route("modules"), Title: "modules"})
		return templates.ModulesList(ctx, w)
	})
}

func Routes(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx util.RequestContext) (int, error) {
		ctx.Title = "Routes"
		ctx.Breadcrumbs = append(aboutBC(ctx), util.Breadcrumb{Path: ctx.Route("routes"), Title: "routes"})
		return templates.RoutesList(ctx, w)
	})
}

func aboutBC(ctx util.RequestContext) util.Breadcrumbs {
	return util.BreadcrumbsSimple(ctx.Route("about"), "about")
}
