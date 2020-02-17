package controllers

import (
	"github.com/gorilla/mux"
	"github.com/kyleu/dbui/internal/app/models"
	"github.com/kyleu/dbui/internal/app/util"
	template "github.com/kyleu/dbui/internal/gen/templates"
	"net/http"
)

func WorkspaceTest(w http.ResponseWriter, r *http.Request) {
	redir(w, r, func(ctx util.RequestContext) (string, error) {
		return ctx.Route("workspace", "p", "test"), nil
	})
}
func Workspace(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)["p"]
	view(getSchema(p), w, r)
}

func WorkspaceTable(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)["p"]
	t := mux.Vars(r)["t"]
	println("Table: " + t)
	view(getSchema(p), w, r)
}

func WorkspaceView(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)["p"]
	v := mux.Vars(r)["v"]
	println("View: " + v)
	view(getSchema(p), w, r)
}

func view(s models.Schema, w http.ResponseWriter, r *http.Request) {
	act(w, r, s.Name, func(ctx util.RequestContext) (int, error) {
		return template.WorkspaceOverview(s, ctx, w)
	})
}

func getSchema(key string) models.Schema {
	s := models.NewSchema("test", "Test Schema")

	if key == "test" {
		s.Tables.Add(
			models.Table{Name: "C"},
			models.Table{Name: "B"},
			models.Table{Name: "A"},
		)
		s.Views.Add(
			models.View{Name: "C"},
			models.View{Name: "B"},
			models.View{Name: "A"},
		)
	}

	return s
}
