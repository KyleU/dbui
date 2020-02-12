package controllers

import (
	template "github.com/kyleu/dbui/internal/app/templates"
	"github.com/kyleu/dbui/internal/app/util"
	"net/http"
)

func Home(res http.ResponseWriter, req *http.Request) {
	act(res, req, "Home", func(ctx util.RequestContext) (int, error) {
		return template.Index(ctx, res)
	})
}

func About(res http.ResponseWriter, req *http.Request) {
	act(res, req, "About", func(ctx util.RequestContext) (int, error) {
		return template.About(ctx, res)
	})
}
