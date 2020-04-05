package controllers

import (
	"net/http"

	"github.com/kyleu/dbui/internal/app/web"

	"emperror.dev/errors"
	"github.com/kyleu/dbui/internal/app/conn/output"

	"github.com/kyleu/dbui/internal/app/conn"
	"github.com/kyleu/dbui/internal/gen/templates"
)

func SQLForm(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx web.RequestContext) (int, error) {
		ctx.Title = "Ad-hoc Query"
		ctx.Breadcrumbs = web.BreadcrumbsSimple(ctx.Route("sql.form"), "ad-hoc")
		return templates.SqlForm("", ctx, w)
	})
}

func SQLRun(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx web.RequestContext) (int, error) {
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
			return 0, errors.WithStack(errors.Wrap(err, "error opening connection"))
		}
		rs, err := conn.RunQueryNoTx(ctx.Logger, connection, ms, conn.Adhoc(sqlArg))
		if err != nil {
			return 0, errors.WithStack(errors.Wrap(err, "error running query"))
		}

		switch fmtArg {
		case "html":
			ctx.Title = "Ad-hoc Results"
			ctx.Breadcrumbs = web.BreadcrumbsSimple(ctx.Route("sql.form"), "ad-hoc")
			return templates.SqlResults(rs, conn.PostgreSQL, err, ctx, w)
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
