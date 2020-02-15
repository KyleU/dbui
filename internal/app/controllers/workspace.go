package controllers

import (
	"github.com/gorilla/mux"
	"github.com/kyleu/dbui/internal/app/models"
	"github.com/kyleu/dbui/internal/app/util"
	template "github.com/kyleu/dbui/internal/gen/templates"
	"net/http"
)

func Workspace(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	view(getSchema(key), w, r)
}

func view(s models.Schema, w http.ResponseWriter, r *http.Request) {
	act(w, r, s.Name, func(ctx util.RequestContext) (int, error) {
		return template.WorkspaceOverview(s, ctx, w)
	})
}

func getSchema(key string) models.Schema {
	s := models.NewSchema("test", "Test Schema")

	if key == "" {
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
