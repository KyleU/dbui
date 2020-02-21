package controllers

import (
	"net/http"

	"github.com/kyleu/dbui/internal/app/conn"
	"github.com/kyleu/dbui/internal/app/util"
	template "github.com/kyleu/dbui/internal/gen/templates"
)

func SQLForm(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx util.RequestContext) (int, error) {
		ctx.Breadcrumbs = util.BreadcrumbsSimple(ctx.Route("sql.form"), "ad-hoc")
		return template.SqlForm("", ctx, w)
	})
}

func SQLRun(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx util.RequestContext) (int, error) {
		_ = r.ParseForm()
		sql := r.Form.Get("sql")
		rs, err := conn.GetResult(ctx.Logger, "", sql)
		ctx.Breadcrumbs = util.BreadcrumbsSimple(ctx.Route("sql.form"), "ad-hoc")
		return template.SqlResults(sql, rs, err, ctx, w)
	})
}
