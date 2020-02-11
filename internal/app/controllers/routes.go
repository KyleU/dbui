package controllers

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/kyleu/dbui/internal/app/util"
	"github.com/sagikazarmark/ocmux"
	"net/http"
)

func BuildRouter(info util.AppInfo) (*mux.Router, error) {
	r := mux.NewRouter()
	r.Use(ocmux.Middleware())

	// Home
	r.Methods("get").Path("/").Handler(addContext(r, info, http.HandlerFunc(Home))).Name("home")

	profile := r.Path("/profile").Subrouter()
	profile.Methods("get").Handler(addContext(r, info, http.HandlerFunc(Profile))).Name("profile")
	profile.Methods("post").Handler(addContext(r, info, http.HandlerFunc(Profile))).Name("profile.save")

	settings := r.Path("/settings").Subrouter()
	settings.Methods("get").Handler(addContext(r, info, http.HandlerFunc(Settings))).Name("settings")

	// Sandbox
	sandbox := r.Path("/sandbox").Subrouter()
	sandbox.Methods("get").Handler(addContext(r, info, http.HandlerFunc(SandboxList))).Name("sandbox.list")
	r.Path("/sandbox/{key}").Methods("get").Handler(addContext(r, info, http.HandlerFunc(SandboxForm))).Name("sandbox.form")
	r.Path("/sandbox/{key}").Methods("post").Handler(addContext(r, info, http.HandlerFunc(SandboxRun))).Name("sandbox.run")

	// Utils
	_ = r.Path("/utils").Subrouter()
	r.Path("/about").Methods("get").Handler(addContext(r, info, http.HandlerFunc(About))).Name("about")

	// Assets
	r.Path("/favicon.ico").Methods("get").HandlerFunc(Favicon).Name("favicon")
	r.PathPrefix("/assets").Methods("get").HandlerFunc(Static).Name("assets")
	return r, nil
}

func addContext(r *mux.Router, info util.AppInfo, next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), "routes", r)
		ctx = context.WithValue(ctx, "info", info)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
