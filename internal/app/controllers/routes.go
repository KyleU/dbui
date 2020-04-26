package controllers

import (
	"context"
	"net/http"

	"github.com/kyleu/dbui/internal/app/config"

	"github.com/gorilla/mux"
	"github.com/sagikazarmark/ocmux"
)

const routesKey = "routes"
const infoKey = "info"

func BuildRouter(info *config.AppInfo) (*mux.Router, error) {
	r := mux.NewRouter()
	r.Use(ocmux.Middleware())

	// Home
	r.Methods(http.MethodGet).Path("/").Handler(addContext(r, info, http.HandlerFunc(Home))).Name("home")
	r.Methods(http.MethodGet).Path("/s").Handler(addContext(r, info, http.HandlerFunc(Socket))).Name("websocket")

	profile := r.Path("/profile").Subrouter()
	profile.Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(Profile))).Name("profile")
	profile.Methods(http.MethodPost).Handler(addContext(r, info, http.HandlerFunc(ProfileSave))).Name("profile.save")

	settings := r.Path("/settings").Subrouter()
	settings.Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(Settings))).Name("settings")
	settings.Methods(http.MethodPost).Handler(addContext(r, info, http.HandlerFunc(SettingsSave))).Name("settings.save")

	// Project
	projects := r.PathPrefix("/workspace").Subrouter()
	projects.Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(WorkspaceTest))).Name("workspace.list")
	r.Path("/w").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(WorkspaceTest))).Name("workspace.test")
	r.Path("/w/new").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(WorkspaceAddForm))).Name("workspace.add.form")
	r.Path("/w/new").Methods(http.MethodPost).Handler(addContext(r, info, http.HandlerFunc(WorkspaceAdd))).Name("workspace.add")
	r.Path("/w/{p}").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(Workspace))).Name("workspace")
	r.Path("/w/{p}/edit").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(WorkspaceEditForm))).Name("workspace.edit.form")
	r.Path("/w/{p}/edit").Methods(http.MethodPost).Handler(addContext(r, info, http.HandlerFunc(WorkspaceEdit))).Name("workspace.edit")
	r.Path("/w/{p}/adhoc").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(WorkspaceAdhocForm))).Name("workspace.adhoc.form")
	r.Path("/w/{p}/adhoc").Methods(http.MethodPost).Handler(addContext(r, info, http.HandlerFunc(WorkspaceAdhoc))).Name("workspace.adhoc")
	r.Path("/w/{p}/t/{t}").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(WorkspaceTable))).Name("workspace.table")
	r.Path("/w/{p}/t/{t}/data").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(WorkspaceData))).Name("workspace.data")

	// Sandbox
	sandbox := r.Path("/sandbox").Subrouter()
	sandbox.Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(SandboxList))).Name("sandbox")
	r.Path("/sandbox/{key}").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(SandboxForm))).Name("sandbox.run")

	// Utils
	_ = r.Path("/utils").Subrouter()
	r.Path("/about").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(About))).Name("about")
	r.Path("/health").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(Health))).Name("health")
	r.Path("/modules").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(Modules))).Name("modules")
	r.Path("/routes").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(Routes))).Name("routes")

	// Assets
	r.Path("/favicon.ico").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(Favicon))).Name("favicon")
	r.PathPrefix("/assets").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(Static))).Name("assets")

	r.PathPrefix("").Handler(addContext(r, info, http.HandlerFunc(NotFound)))

	return r, nil
}

func addContext(router *mux.Router, info *config.AppInfo, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer InternalServerError(router, info, w, r)
		ctx := context.WithValue(r.Context(), routesKey, router)
		ctx = context.WithValue(ctx, infoKey, info)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
