package controllers

import (
	"github.com/kyleu/dbui/internal/app/conn"
	"github.com/kyleu/dbui/internal/app/util"
	template "github.com/kyleu/dbui/internal/gen/templates"
	"net/http"
)

func SqlForm(w http.ResponseWriter, r *http.Request) {
	act(w, r, "Ad-hoc SQL Query", func(ctx util.RequestContext) (int, error) {
		return template.SqlForm("", ctx, w)
	})
}

func SqlRun(w http.ResponseWriter, r *http.Request) {
	act(w, r, "SQL Results", func(ctx util.RequestContext) (int, error) {
		_ = r.ParseForm()
		sql := r.Form.Get("sql")
		rs, err := conn.GetResult(ctx.Logger, "", sql)
		return template.SqlResults(sql, rs, err, ctx, w)
	})
}
