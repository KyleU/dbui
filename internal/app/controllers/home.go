package controllers

import (
	template "github.com/kyleu/dbui/internal/app/templates"
	"net/http"
)

func Home(res http.ResponseWriter, req *http.Request) {
	template.Index(prepHtml(res, req), res)
}
