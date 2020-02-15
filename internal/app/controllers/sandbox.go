package controllers

import (
	"github.com/gorilla/mux"
	"github.com/kyleu/dbui/internal/app/util"
	template "github.com/kyleu/dbui/internal/gen/templates"
	"net/http"
)

var _sandboxes = []string{"gallery", "testbed"}

func SandboxList(res http.ResponseWriter, req *http.Request) {
	act(res, req, "Sandbox List", func(ctx util.RequestContext) (int, error) {
		return template.SandboxList(_sandboxes, ctx, res)
	})
}

func SandboxForm(res http.ResponseWriter, req *http.Request) {
	key := mux.Vars(req)["key"]
	if key == "testbed" {
		// x := 0
		// _ = 10 / x
		panic("!!!!")
	}
	act(res, req, "Sandbox [" + key + "]", func(ctx util.RequestContext) (int, error) {
		return template.SandboxForm(key, ctx, res)
	})
}
