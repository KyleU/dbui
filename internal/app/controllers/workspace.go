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
		ctx.Breadcrumbs = append(bc, tableBC(ctx, s.ID, t))
		return templates.WorkspaceTable(s, t, ctx, w)
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
		ctx.Breadcrumbs = append(bc, viewBC(ctx, s.ID, v))
		return templates.WorkspaceView(s, v, ctx, w)
	})
}

func WorkspaceData(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)["p"]
	model := mux.Vars(r)["t"]
	name := mux.Vars(r)["n"]
	act(w, r, func(ctx util.RequestContext) (int, error) {
		s, bc, err := load(ctx, p, false)
		if err != nil {
			return 0, err
		}
		db, connectMS, err := ctx.AppInfo.ConfigService.GetConnection(s.ID)
		if err != nil {
			return 0, err
		}
		rs, err := conn.GetResult(ctx.AppInfo.Logger, db, connectMS, "select * from \"" + name + "\"")
		if err != nil {
			return 0, err
		}
		dc := util.Breadcrumb{Path: ctx.Route("workspace.data", "p", s.ID, "t", model, "n", name), Title: "data"}
		var mc = util.Breadcrumb{}
		switch model {
		case "t":
			mc = tableBC(ctx, s.ID, name)
		case "v":
			mc = viewBC(ctx, s.ID, name)
		}
		ctx.Breadcrumbs = append(bc, mc, dc)
		return templates.WorkspaceData(s, model, name, rs, ctx, w)
	})
}

func tableBC(ctx util.RequestContext, id string, name string) util.Breadcrumb {
	return util.Breadcrumb{Path: ctx.Route("workspace.table", "p", id, "t", name), Title: name}
}

func viewBC(ctx util.RequestContext, id string, name string) util.Breadcrumb {
	return util.Breadcrumb{Path: ctx.Route("workspace.view", "p", id, "v", name), Title: name}
}

func load(ctx util.RequestContext, p string, forceReload bool) (*schema.Schema, util.Breadcrumbs, error) {
	s, err := schema.GetSchema(ctx.AppInfo, p, forceReload)
	if err != nil {
		return nil, nil, err
	}
	bc := util.BreadcrumbsSimple(ctx.Route("workspace", "p", s.ID), s.ID)
	return s, bc, nil
}
