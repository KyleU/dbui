package controllers

import (
	"github.com/kyleu/dbui/internal/app/conn"
	"net/http"

	"github.com/kyleu/dbui/internal/app/schema"
	"github.com/kyleu/dbui/internal/gen/templates"

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
		s, bc, err := load(ctx, p, false)
		if err != nil {
			return 0, err
		}
		ctx.Breadcrumbs = bc
		return templates.WorkspaceOverview(s, "overview", ctx, w)
	})
}

func WorkspaceTable(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)["p"]
	t := mux.Vars(r)["t"]
	act(w, r, func(ctx util.RequestContext) (int, error) {
		s, bc, err := load(ctx, p, false)
		if err != nil {
			return 0, err
		}
		tc := util.Breadcrumb{Path: ctx.Route("workspace.table", "p", s.ID, "t", t), Title: t}
		ctx.Breadcrumbs = append(bc, tc)
		return templates.WorkspaceTable(s, t, ctx, w)
	})
}

func WorkspaceTableData(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)["p"]
	t := mux.Vars(r)["t"]
	act(w, r, func(ctx util.RequestContext) (int, error) {
		s, bc, err := load(ctx, p, false)
		if err != nil {
			return 0, err
		}
		db, connectMS, err := ctx.AppInfo.ConfigService.GetConnection(s.ID)
		if err != nil {
			return 0, err
		}
		rs, err := conn.GetResult(ctx.AppInfo.Logger, db, connectMS, "select * from \"" + t + "\"")
		if err != nil {
			return 0, err
		}
		dc := util.Breadcrumb{Path: ctx.Route("workspace.table.data", "p", s.ID, "t", t), Title: "data"}
		tc := util.Breadcrumb{Path: ctx.Route("workspace.table", "p", s.ID, "t", t), Title: t}
		ctx.Breadcrumbs = append(bc, tc, dc)
		return templates.WorkspaceData(s, "table", t, rs, ctx, w)
	})
}

func WorkspaceView(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)["p"]
	v := mux.Vars(r)["v"]
	act(w, r, func(ctx util.RequestContext) (int, error) {
		s, bc, err := load(ctx, p, false)
		if err != nil {
			return 0, err
		}
		vc := util.Breadcrumb{Path: ctx.Route("workspace.view", "p", s.ID, "v", v), Title: v}
		ctx.Breadcrumbs = append(bc, vc)
		return templates.WorkspaceView(s, v, ctx, w)
	})
}

func WorkspaceViewData(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)["p"]
	v := mux.Vars(r)["v"]
	act(w, r, func(ctx util.RequestContext) (int, error) {
		s, bc, err := load(ctx, p, false)
		if err != nil {
			return 0, err
		}
		db, connectMS, err := ctx.AppInfo.ConfigService.GetConnection(s.ID)
		if err != nil {
			return 0, err
		}
		rs, err := conn.GetResult(ctx.AppInfo.Logger, db, connectMS, "select * from \"" + v + "\"")
		if err != nil {
			return 0, err
		}
		dc := util.Breadcrumb{Path: ctx.Route("workspace.view.data", "p", s.ID, "v", v), Title: "data"}
		tc := util.Breadcrumb{Path: ctx.Route("workspace.view", "p", s.ID, "v", v), Title: v}
		ctx.Breadcrumbs = append(bc, tc, dc)
		return templates.WorkspaceData(s, "view", v, rs, ctx, w)
	})
}

func load(ctx util.RequestContext, p string, forceReload bool) (*schema.Schema, util.Breadcrumbs, error) {
	s, err := schema.GetSchema(ctx.AppInfo, p, forceReload)
	if err != nil {
		return nil, nil, err
	}
	bc := util.BreadcrumbsSimple(ctx.Route("workspace", "p", s.ID), s.ID)
	return s, bc, nil
}
