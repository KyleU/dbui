package controllers

import (
	"github.com/kyleu/dbui/internal/app/conn"
	"github.com/kyleu/dbui/internal/app/web"
	"net/http"

	"github.com/kyleu/dbui/internal/app/schema"
	"github.com/kyleu/dbui/internal/gen/templates"

	"github.com/gorilla/mux"
)

func WorkspaceTest(w http.ResponseWriter, r *http.Request) {
	redir(w, r, func(ctx web.RequestContext) (string, error) {
		return ctx.Route("workspace", "p", "test"), nil
	})
}

func Workspace(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx web.RequestContext) (int, error) {
		p := mux.Vars(r)["p"]
		s, bc, err := load(ctx, p, false)
		if err != nil {
			return 0, err
		}
		ctx.Title = s.Name
		ctx.Breadcrumbs = bc
		return templates.WorkspaceOverview(s, "overview", ctx, w)
	})
}

func WorkspaceTable(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx web.RequestContext) (int, error) {
		p := mux.Vars(r)["p"]
		t := mux.Vars(r)["t"]
		s, bc, err := load(ctx, p, false)
		if err != nil {
			return 0, err
		}
		ctx.Title = "Table [" + t + "]"
		ctx.Breadcrumbs = append(bc, tableBC(ctx, s.ID, t))
		return templates.WorkspaceTable(s, t, ctx, w)
	})
}

func WorkspaceData(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx web.RequestContext) (int, error) {
		p := mux.Vars(r)["p"]
		name := mux.Vars(r)["t"]
		opts := web.FromQueryString(ctx.Profile, true, r.URL.Query())
		s, bc, err := load(ctx, p, false)
		if err != nil {
			return 0, err
		}
		db, connectMS, err := ctx.AppInfo.ConfigService.GetConnection(s.ID)
		if err != nil {
			return 0, err
		}
		rs, err := conn.GetResult(ctx.AppInfo.Logger, db, connectMS, opts.ToSQL(s.Engine, name))
		if err != nil {
			return 0, err
		}
		dc := web.Breadcrumb{Path: ctx.Route("workspace.data", "p", s.ID, "t", name), Title: "data"}
		ctx.Title = "[" + name + "] Data"
		var tc= tableBC(ctx, s.ID, name)
		ctx.Breadcrumbs = append(bc, tc, dc)
		return templates.WorkspaceData(s, name, rs, opts, ctx, w)
	})
}

func tableBC(ctx web.RequestContext, id string, name string) web.Breadcrumb {
	return web.Breadcrumb{Path: ctx.Route("workspace.table", "p", id, "t", name), Title: name}
}

func viewBC(ctx web.RequestContext, id string, name string) web.Breadcrumb {
	return web.Breadcrumb{Path: ctx.Route("workspace.view", "p", id, "v", name), Title: name}
}

func load(ctx web.RequestContext, p string, forceReload bool) (*schema.Schema, web.Breadcrumbs, error) {
	s, err := schema.GetSchema(ctx.AppInfo, p, forceReload)
	if err != nil {
		return nil, nil, err
	}
	bc := web.BreadcrumbsSimple(ctx.Route("workspace", "p", s.ID), s.ID)
	return s, bc, nil
}
