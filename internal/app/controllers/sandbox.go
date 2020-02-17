package controllers

import (
	"github.com/gorilla/mux"
	"github.com/kyleu/dbui/internal/app/util"
	template "github.com/kyleu/dbui/internal/gen/templates"
	"net/http"
)

var _sandboxes = []string{"gallery", "testbed"}

func SandboxList(w http.ResponseWriter, r *http.Request) {
	act(w, r, "Sandbox List", func(ctx util.RequestContext) (int, error) {
		return template.SandboxList(_sandboxes, ctx, w)
	})
}

func SandboxForm(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	if key == "testbed" {
		panic("!!!!")
	}
	act(w, r, "Sandbox [" + key + "]", func(ctx util.RequestContext) (int, error) {
		return template.SandboxForm(key, ctx, w)
	})
}
