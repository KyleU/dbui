package controllers

import (
	"github.com/kyleu/dbui/internal/app/util"
	template "github.com/kyleu/dbui/internal/gen/templates"
	"net/http"
)

func Health(res http.ResponseWriter, req *http.Request) {
	act(res, req, "Routes", func(ctx util.RequestContext) (int, error) {
		res.Write([]byte("OK"))
		return 0, nil
	})
}

func Routes(res http.ResponseWriter, req *http.Request) {
	act(res, req, "Routes", func(ctx util.RequestContext) (int, error) {
		return template.RoutesList(ctx, res)
	})
}

func Modules(res http.ResponseWriter, req *http.Request) {
	act(res, req, "Routes", func(ctx util.RequestContext) (int, error) {
		return template.ModulesList(ctx, res)
	})
}
