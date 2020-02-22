package controllers

import (
	"net/http"

	"emperror.dev/errors"
	"github.com/kyleu/dbui/internal/app/conn/output"

	"github.com/kyleu/dbui/internal/app/conn"
	"github.com/kyleu/dbui/internal/app/util"
	"github.com/kyleu/dbui/internal/gen/templates"
)

func SQLForm(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx util.RequestContext) (int, error) {
		ctx.Breadcrumbs = util.BreadcrumbsSimple(ctx.Route("sql.form"), "ad-hoc")
		return templates.SqlForm("", ctx, w)
	})
}

func SQLRun(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx util.RequestContext) (int, error) {
		_ = r.ParseForm()
		sqlArg := r.Form.Get("sql")
		connArg := r.Form.Get("conn")
		fmtArg := r.Form.Get("fmt")
		if fmtArg == "" {
			fmtArg = "html"
		}
		if connArg == "_url" {
			connArg = r.Form.Get("url")
		}
		connection, ms, err := ctx.AppInfo.ConfigService.GetConnection(connArg)
		if err != nil {
			return 0, errors.WithStack(errors.Wrap(err, "Error opening connection"))
		}
		rs, err := conn.GetResult(ctx.Logger, connection, ms, sqlArg)
		switch fmtArg {
		case "html":
			ctx.Breadcrumbs = util.BreadcrumbsSimple(ctx.Route("sql.form"), "ad-hoc")
			return templates.SqlResults(rs, err, ctx, w)
		case "csv":
			content, err := output.AsString(rs)
			if err != nil {
				return 0, err
			}
			w.Header().Set("Content-Type", "text/csv; charset=utf-8")
			w.Header().Set("Content-Disposition", "attachment; filename=\"export.csv\"")
			return w.Write([]byte(content))
		case "json":
			content, err := output.AsJson(rs)
			if err != nil {
				return 0, err
			}
			w.Header().Set("Content-Type", "text/json; charset=utf-8")
			return w.Write([]byte(content))
		default:
			return 0, errors.New("Invalid output format [" + fmtArg + "]")
		}
	})
}
