package controllers

import (
	template "github.com/kyleu/dbui/internal/app/templates"
	"net/http"
)

func Settings(res http.ResponseWriter, req *http.Request) {
	template.Settings(prepHtml(res, req, ""), res)
}
