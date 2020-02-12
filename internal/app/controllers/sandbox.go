package controllers

import (
	"github.com/gorilla/mux"
	template "github.com/kyleu/dbui/internal/app/templates"
	"github.com/kyleu/dbui/internal/app/util"
	"net/http"
)

var _sandboxes = []string{"routes", "gallery", "testbed"}

func SandboxList(res http.ResponseWriter, req *http.Request) {
	act(res, req, "Sandbox List", func(ctx util.RequestContext) (int, error) {
		return template.SandboxList(_sandboxes, ctx, res)
	})
}

func SandboxForm(res http.ResponseWriter, req *http.Request) {
	key := mux.Vars(req)["key"]
	act(res, req, "Sandbox [" + key + "]", func(ctx util.RequestContext) (int, error) {
		return template.SandboxForm(key, ctx, res)
	})
}
