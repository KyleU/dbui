package controllers

import (
	"github.com/kyleu/dbui/internal/app/web"
	"net/http"

	"github.com/kyleu/dbui/internal/gen/templates"
)

func Health(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx web.RequestContext) (int, error) {
		_, _ = w.Write([]byte("OK"))
		return 0, nil
	})
}

func Modules(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx web.RequestContext) (int, error) {
		ctx.Title = "Modules"
		ctx.Breadcrumbs = append(aboutBC(ctx), web.Breadcrumb{Path: ctx.Route("modules"), Title: "modules"})
		return templates.ModulesList(ctx, w)
	})
}

func Routes(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx web.RequestContext) (int, error) {
		ctx.Title = "Routes"
		ctx.Breadcrumbs = append(aboutBC(ctx), web.Breadcrumb{Path: ctx.Route("routes"), Title: "routes"})
		return templates.RoutesList(ctx, w)
	})
}

func aboutBC(ctx web.RequestContext) web.Breadcrumbs {
	return web.BreadcrumbsSimple(ctx.Route("about"), "about")
}
