package controllers

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"

	"emperror.dev/errors"
	"github.com/gofrs/uuid"
	"github.com/kyleu/dbui/internal/app/config"
	"github.com/kyleu/dbui/internal/app/web"

	"github.com/kyleu/dbui/internal/gen/templates"
)

func WorkspaceAddForm(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx web.RequestContext) (int, error) {
		ctx.Title = "New Workspace"
		bc := web.BreadcrumbsSimple(ctx.Route("workspace.add.form"), "new")
		ctx.Breadcrumbs = bc
		p := &config.Project{
			EngineString: "pgx",
		}
		return templates.WorkspaceForm(p, ctx, w)
	})
}

func WorkspaceAdd(w http.ResponseWriter, r *http.Request) {
	redir(w, r, func(ctx web.RequestContext) (string, error) {
		_ = r.ParseForm()
		key := r.Form.Get("key")
		if key == "" {
			return ctx.Route("workspace.add.form"), nil
		}
		o := r.Form.Get("owner")
		owner, err := uuid.FromString(o)
		if err != nil {
			return ctx.Route("workspace.add.form"), errors.WithStack(errors.Wrap(err, "error parsing owner uuid from ["+o+"]"))
		}
		username := sql.NullString{
			String: r.Form.Get("username"),
			Valid:  true,
		}
		password := sql.NullString{
			String: r.Form.Get("password"),
			Valid:  true,
		}
		p := config.Project{
			Key:          key,
			Title:        r.Form.Get("title"),
			Description:  r.Form.Get("description"),
			Owner:        owner,
			EngineString: r.Form.Get("engine"),
			URL:          r.Form.Get("url"),
			Username:     username,
			Password:     password,
		}
		err = ctx.AppInfo.ConfigService.ProjectRegistry.Add(true, &p)
		if err != nil {
			return ctx.Route("workspace.add.form"), errors.WithStack(errors.Wrap(err, "error adding ["+key+"] to project registry"))
		}
		return ctx.Route("workspace", "p", key), nil
	})
}

func WorkspaceEditForm(w http.ResponseWriter, r *http.Request) {
	act(w, r, func(ctx web.RequestContext) (int, error) {
		p := mux.Vars(r)["p"]
		proj, err := ctx.AppInfo.ConfigService.ProjectRegistry.Get(p)
		if err != nil {
			return 0, errors.WithStack(err)
		}
		ctx.Title = "Edit " + proj.Title
		bc := web.BreadcrumbsSimple(ctx.Route("workspace", "p", proj.Key), proj.Key)
		ec := web.BreadcrumbsSimple(ctx.Route("workspace.edit.form"), "edit")
		ctx.Breadcrumbs = append(bc, ec...)
		return templates.WorkspaceForm(proj, ctx, w)
	})
}

func WorkspaceEdit(w http.ResponseWriter, r *http.Request) {
	redir(w, r, func(ctx web.RequestContext) (string, error) {
		p := mux.Vars(r)["p"]
		_ = r.ParseForm()
		key := r.Form.Get("key")
		if key == "" {
			return ctx.Route("workspace.edit.form"), nil
		}
		o := r.Form.Get("owner")
		owner, err := uuid.FromString(o)
		if err != nil {
			return ctx.Route("workspace.edit.form"), errors.WithStack(errors.Wrap(err, "error parsing owner uuid from ["+o+"]"))
		}
		username := sql.NullString{
			String: r.Form.Get("username"),
			Valid:  true,
		}
		password := sql.NullString{
			String: r.Form.Get("password"),
			Valid:  true,
		}
		proj := &config.Project{
			Key:          key,
			Title:        r.Form.Get("title"),
			Description:  r.Form.Get("description"),
			Owner:        owner,
			EngineString: r.Form.Get("engine"),
			URL:          r.Form.Get("url"),
			Username:     username,
			Password:     password,
		}
		err = ctx.AppInfo.ConfigService.ProjectRegistry.Add(true, proj)
		if err != nil {
			return ctx.Route("workspace.edit.form"), errors.WithStack(errors.Wrap(err, "error saving ["+p+"] to project registry"))
		}
		return ctx.Route("workspace", "p", p), nil
	})
}
