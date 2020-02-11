package controllers

import (
	"github.com/gorilla/mux"
	template "github.com/kyleu/dbui/internal/app/templates"
	"net/http"
)

var _sandboxes = []string{"routes", "gallery", "testbed"}

func SandboxList(res http.ResponseWriter, req *http.Request) {
	template.SandboxList(_sandboxes, prepHtml(res, req, "Sandbox List"), res)
}

func SandboxForm(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	key := vars["key"]
	template.SandboxForm(key, prepHtml(res, req, "Sandbox [" + key + "]"), res)
}

func SandboxRun(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	key := vars["key"]
	template.SandboxForm(key, prepHtml(res, req, "Sandbox [" + key + "]"), res)
}

