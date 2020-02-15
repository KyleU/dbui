package controllers

import (
	"github.com/kyleu/dbui/internal/app/util"
	template "github.com/kyleu/dbui/internal/gen/templates"
	"net/http"
)

func Workspace(res http.ResponseWriter, req *http.Request) {
	schema := "Test"
	act(res, req, schema, func(ctx util.RequestContext) (int, error) {
		return template.Workspace(schema, ctx, res)
	})
}
