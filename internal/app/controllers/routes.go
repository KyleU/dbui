package controllers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kyleu/dbui/internal/app/util"
	"github.com/sagikazarmark/ocmux"
)

const routesKey = "routes"
const infoKey = "info"

func BuildRouter(info util.AppInfo) (*mux.Router, error) {
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
	projects := r.PathPrefix("/projects").Subrouter()
	projects.Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(WorkspaceTest))).Name("project.list")
	r.Path("/q").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(WorkspaceTest))).Name("workspace.test")
	r.Path("/q/{p}").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(Workspace))).Name("workspace")
	r.Path("/q/{p}/t/{t}").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(WorkspaceTable))).Name("workspace.table")
	r.Path("/q/{p}/t/{t}/data").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(WorkspaceTableData))).Name("workspace.table.data")
	r.Path("/q/{p}/v/{v}").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(WorkspaceView))).Name("workspace.view")
	r.Path("/q/{p}/v/{v}/data").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(WorkspaceViewData))).Name("workspace.view.data")

	// Sandbox
	sandbox := r.Path("/sandbox").Subrouter()
	sandbox.Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(SandboxList))).Name("sandbox")
	r.Path("/sandbox/{key}").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(SandboxForm))).Name("sandbox.run")

	// Ad-hoc SQL Queries
	sql := r.Path("/sql").Subrouter()
	sql.Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(SQLForm))).Name("sql.form")
	sql.Methods(http.MethodPost).Handler(addContext(r, info, http.HandlerFunc(SQLRun))).Name("sql.run")

	// Utils
	_ = r.Path("/utils").Subrouter()
	r.Path("/about").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(About))).Name("about")
	r.Path("/health").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(Health))).Name("health")
	r.Path("/modules").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(Modules))).Name("modules")
	r.Path("/routes").Methods(http.MethodGet).Handler(addContext(r, info, http.HandlerFunc(Routes))).Name("routes")

	// Assets
	r.Path("/favicon.ico").Methods(http.MethodGet).HandlerFunc(Favicon).Name("favicon")
	r.PathPrefix("/assets").Methods(http.MethodGet).HandlerFunc(Static).Name("assets")

	r.PathPrefix("").Handler(addContext(r, info, http.HandlerFunc(NotFound)))

	return r, nil
}

func addContext(router *mux.Router, info util.AppInfo, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer InternalServerError(router, info, w, r)
		ctx := context.WithValue(r.Context(), routesKey, router)
		ctx = context.WithValue(ctx, infoKey, info)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
