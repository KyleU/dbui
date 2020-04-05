package controllers

import (
	"emperror.dev/errors"
	"github.com/kyleu/dbui/internal/app/conn"
	"github.com/kyleu/dbui/internal/app/conn/output"
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
			return 0, errors.WithStack(errors.Wrap(err, "error loading workspace [" + p + "]"))
		}
		ctx.Title = s.Name
		ctx.Breadcrumbs = bc
		return templates.WorkspaceOverview(s, "overview", ctx, w)
	})
}

func WorkspaceAdhocForm(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx web.RequestContext) (int, error) {
		p := mux.Vars(r)["p"]
		sqlSet, _ := r.URL.Query()["sql"]
		sql := ""
		if len(sqlSet) > 0 {
			sql = sqlSet[0]
		}
		s, bc, err := load(ctx, p, false)
		if err != nil {
			return 0, errors.WithStack(errors.Wrap(err, "error loading workspace [" + p + "]"))
		}
		ctx.Title = "New Query"
		ctx.Breadcrumbs = append(bc, web.Breadcrumb{Path: ctx.Route("workspace.adhoc.form", "p", p), Title: "query"})
		return templates.WorkspaceAdhoc(s, sql, nil, ctx, w)
	})
}

func WorkspaceAdhoc(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx web.RequestContext) (int, error) {
		p := mux.Vars(r)["p"]
		s, bc, err := load(ctx, p, false)
		if err != nil {
			return 0, errors.WithStack(errors.Wrap(err, "error loading workspace [" + p + "]"))
		}

		_ = r.ParseForm()
		sqlArg := r.Form.Get("sql")
		fmtArg := r.Form.Get("fmt")
		if fmtArg == "" {
			fmtArg = "html"
		}
		connection, ms, err := ctx.AppInfo.ConfigService.GetConnection(s.ID)
		if err != nil {
			return 0, errors.WithStack(errors.Wrap(err, "error opening connection"))
		}

		rs, err := conn.RunQueryNoTx(ctx.Logger, connection, ms, conn.Adhoc(sqlArg))
		if err != nil {
			return 0, errors.WithStack(errors.Wrap(err, "error running query"))
		}

		switch fmtArg {
		case "html":
			ctx.Title = "Query Results"
			ctx.Breadcrumbs = append(bc, web.Breadcrumb{Path: ctx.Route("workspace.adhoc.form", "p", p), Title: "query"})
			return templates.WorkspaceAdhoc(s, sqlArg, rs, ctx, w)
			// return templates.SqlResults(rs, err, ctx, w)
		case "csv":
			content, err := output.AsString(rs)
			if err != nil {
				return 0, errors.WithStack(errors.Wrap(err, "error formatting csv output"))
			}
			w.Header().Set("Content-Type", "text/csv; charset=utf-8")
			w.Header().Set("Content-Disposition", "attachment; filename=\"export.csv\"")
			return w.Write([]byte(content))
		case "json":
			content, err := output.AsJson(rs)
			if err != nil {
				return 0, errors.WithStack(errors.Wrap(err, "error formatting json output"))
			}
			w.Header().Set("Content-Type", "text/json; charset=utf-8")
			return w.Write([]byte(content))
		default:
			return 0, errors.New("Invalid output format [" + fmtArg + "]")
		}
	})
}

func WorkspaceTable(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx web.RequestContext) (int, error) {
		p := mux.Vars(r)["p"]
		t := mux.Vars(r)["t"]
		s, bc, err := load(ctx, p, false)
		if err != nil {
			return 0, errors.WithStack(errors.Wrap(err, "error loading workspace [" + p + "]"))
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
			return 0, errors.WithStack(errors.Wrap(err, "error loading workspace [" + p + "]"))
		}
		db, connectMS, err := ctx.AppInfo.ConfigService.GetConnection(s.ID)
		if err != nil {
			return 0, errors.WithStack(errors.Wrap(err, "error opening connection to [" + s.ID + "]"))
		}
		rs, err := conn.RunQueryNoTx(ctx.AppInfo.Logger, db, connectMS, conn.Adhoc(opts.ToSQL(name)))
		if err != nil {
			return 0, errors.WithStack(errors.Wrap(err, "error running query against project [" + p + "]"))
		}

		table := s.Tables.Get(name)
		if table != nil {
			rs.Columns = table.Columns
		}

		dc := web.Breadcrumb{Path: ctx.Route("workspace.data", "p", s.ID, "t", name), Title: "data"}
		ctx.Title = "[" + name + "] Data"
		var tc = tableBC(ctx, s.ID, name)
		ctx.Breadcrumbs = append(bc, tc, dc)
		return templates.WorkspaceData(s, name, rs, opts, ctx, w)
	})
}

func tableBC(ctx web.RequestContext, id string, name string) web.Breadcrumb {
	return web.Breadcrumb{Path: ctx.Route("workspace.table", "p", id, "t", name), Title: name}
}

func load(ctx web.RequestContext, p string, forceReload bool) (*schema.Schema, web.Breadcrumbs, error) {
	s, err := schema.GetSchema(ctx.AppInfo, p, forceReload)
	if err != nil {
		return nil, nil, errors.WithStack(errors.Wrap(err, "error loading workspace [" + p + "]"))
	}
	key := s.ID
	if key == "_root" {
		key = "system"
	}
	bc := web.BreadcrumbsSimple(ctx.Route("workspace", "p", s.ID), key)
	return s, bc, nil
}
