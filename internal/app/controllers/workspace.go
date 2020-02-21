package controllers

import (
	"net/http"

	"github.com/kyleu/dbui/internal/app/schema"
	template "github.com/kyleu/dbui/internal/gen/templates"

	"github.com/gorilla/mux"
	"github.com/kyleu/dbui/internal/app/util"
)

func WorkspaceTest(w http.ResponseWriter, r *http.Request) {
	redir(w, r, func(ctx util.RequestContext) (string, error) {
		return ctx.Route("workspace", "p", "test"), nil
	})
}

func Workspace(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)["p"]
	act(w, r, func(ctx util.RequestContext) (int, error) {
		s := schema.GetSchema(ctx.AppInfo, p, false)
		ctx.Breadcrumbs = util.BreadcrumbsSimple(ctx.Route("workspace", "p", s.ID), s.ID)
		return template.WorkspaceOverview(s, "overview", ctx, w)
	})
}

func WorkspaceTable(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)["p"]
	t := mux.Vars(r)["t"]
	act(w, r, func(ctx util.RequestContext) (int, error) {
		s := schema.GetSchema(ctx.AppInfo, p, false)
		bc := util.Breadcrumb{Path: ctx.Route("workspace.table", "p", s.ID, "t", t), Title: t}
		ctx.Breadcrumbs = append(util.BreadcrumbsSimple(ctx.Route("workspace", "p", s.ID), s.ID), bc)
		return template.WorkspaceTable(s, t, ctx, w)
	})
}

func WorkspaceTableData(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)["p"]
	t := mux.Vars(r)["t"]
	act(w, r, func(ctx util.RequestContext) (int, error) {
		s := schema.GetSchema(ctx.AppInfo, p, false)
		dc := util.Breadcrumb{Path: ctx.Route("workspace.table.data", "p", s.ID, "t", t), Title: "data"}
		bc := util.Breadcrumb{Path: ctx.Route("workspace.table", "p", s.ID, "t", t), Title: t}
		ctx.Breadcrumbs = append(util.BreadcrumbsSimple(ctx.Route("workspace", "p", s.ID), s.ID), bc, dc)
		return template.WorkspaceTable(s, t, ctx, w)
	})
}

func WorkspaceView(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)["p"]
	v := mux.Vars(r)["v"]
	act(w, r, func(ctx util.RequestContext) (int, error) {
		s := schema.GetSchema(ctx.AppInfo, p, false)
		bc := util.Breadcrumb{Path: ctx.Route("workspace.view", "p", s.ID, "v", v), Title: v}
		ctx.Breadcrumbs = append(util.BreadcrumbsSimple(ctx.Route("workspace", "p", s.ID), s.ID), bc)
		return template.WorkspaceView(s, v, ctx, w)
	})
}

func WorkspaceViewData(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)["p"]
	v := mux.Vars(r)["v"]
	act(w, r, func(ctx util.RequestContext) (int, error) {
		s := schema.GetSchema(ctx.AppInfo, p, false)
		dc := util.Breadcrumb{Path: ctx.Route("workspace.view.data", "p", s.ID, "v", v), Title: "data"}
		bc := util.Breadcrumb{Path: ctx.Route("workspace.view", "p", s.ID, "v", v), Title: v}
		ctx.Breadcrumbs = append(util.BreadcrumbsSimple(ctx.Route("workspace", "p", s.ID), s.ID), bc, dc)
		return template.WorkspaceView(s, v, ctx, w)
	})
}
