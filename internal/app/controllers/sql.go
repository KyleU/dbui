package controllers

import (
	"github.com/kyleu/dbui/internal/app/conn"
	"github.com/kyleu/dbui/internal/app/util"
	template "github.com/kyleu/dbui/internal/gen/templates"
	"net/http"
)

func SqlForm(res http.ResponseWriter, req *http.Request) {
	act(res, req, "Ad-hoc SQL Query", func(ctx util.RequestContext) (int, error) {
		return template.SqlForm("", ctx, res)
	})
}

func SqlRun(res http.ResponseWriter, req *http.Request) {
	act(res, req, "SQL Results", func(ctx util.RequestContext) (int, error) {
		_ = req.ParseForm()
		sql := req.Form.Get("sql")
		rs, err := conn.GetResult(ctx.Logger, "", sql)
		return template.SqlResults(sql, rs, err, ctx, res)
	})
}
