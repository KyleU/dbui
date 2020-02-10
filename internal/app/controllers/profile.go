package controllers

import (
	template "github.com/kyleu/dbui/internal/app/templates"
	"net/http"
)

func Profile(res http.ResponseWriter, req *http.Request) {
	template.Profile(prepHtml(res, req, ""), res)
}
