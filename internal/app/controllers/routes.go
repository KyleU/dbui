package controllers

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/kyleu/dbui/internal/app/util"
	"github.com/sagikazarmark/ocmux"
	"net/http"
)

func BuildRouter(info util.AppInfo) *mux.Router {
	r := mux.NewRouter()
	r.Use(ocmux.Middleware())

	// Home
	r.Handle("/", addContext(r, info, http.HandlerFunc(Home))).Name("home")
	r.Handle("/profile", addContext(r, info, http.HandlerFunc(Profile))).Name("profile")
	r.Handle("/settings", addContext(r, info, http.HandlerFunc(Settings))).Name("settings")

	// Sandbox
	r.Handle("/sandbox", addContext(r, info, http.HandlerFunc(SandboxList))).Name("sandbox.list")
	r.Handle("/sandbox/{key}", addContext(r, info, http.HandlerFunc(SandboxForm))).Name("sandbox.form")

	// Assets
	r.HandleFunc("/favicon.ico", Favicon).Name("favicon").Name("favicon")
	r.PathPrefix("/assets").HandlerFunc(Static).Name("assets")
	return r
}

func addContext(r *mux.Router, info util.AppInfo, next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), "routes", r)
		ctx = context.WithValue(ctx, "info", info)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
