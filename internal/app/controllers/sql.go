package controllers

import (
	template "github.com/kyleu/dbui/internal/app/templates"
	"github.com/kyleu/dbui/internal/app/util"
	"net/http"
)

func SqlForm(res http.ResponseWriter, req *http.Request) {
	act(res, req, "Home", func(ctx util.RequestContext) (int, error) {
		return template.SqlView("", ctx, res)
	})
}

func SqlRun(res http.ResponseWriter, req *http.Request) {
	act(res, req, "Home", func(ctx util.RequestContext) (int, error) {
		return template.SqlView("", ctx, res)
	})
}
