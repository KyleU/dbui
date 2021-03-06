package web

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/kyleu/dbui/internal/app/config"
	"github.com/kyleu/dbui/internal/app/util"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"logur.dev/logur"
)

type Breadcrumb struct {
	Path  string
	Title string
}

type Breadcrumbs []Breadcrumb

func (bc Breadcrumbs) Title(ai *config.AppInfo) string {
	if len(bc) == 0 {
		return util.AppName
	}
	return bc[len(bc)-1].Title
}

func BreadcrumbsSimple(path string, title string) Breadcrumbs {
	return []Breadcrumb{
		{path, title},
	}
}

type RequestContext struct {
	AppInfo     *config.AppInfo
	Logger      logur.LoggerFacade
	Profile     util.UserProfile
	Routes      *mux.Router
	Title       string
	Breadcrumbs Breadcrumbs
	Flashes     []string
	Session     sessions.Session
}

func (r *RequestContext) Route(act string, pairs ...string) string {
	route := r.Routes.Get(act)
	if route == nil {
		r.AppInfo.Logger.Warn("Cannot find route at path [" + act + "]")
		return "/routenotfound"
	}
	url, err := route.URL(pairs...)
	if err != nil {
		r.AppInfo.Logger.Warn("Cannot bind route at path [" + act + "]")
		return "/routeerror"
	}
	return url.Path
}

var sessionKey = func() string {
	x := os.Getenv("SESSION_KEY")
	if x == "" {
		x = "random_secret_key"
	}
	return x
}()

var store = sessions.NewCookieStore([]byte(sessionKey))

const sessionName = util.AppName + "-session"

func ExtractContext(r *http.Request) RequestContext {
	ai := r.Context().Value("info").(*config.AppInfo)
	routes := r.Context().Value("routes").(*mux.Router)
	prof := util.SystemProfile
	session, err := store.Get(r, sessionName)
	if err != nil {
		session = sessions.NewSession(store, sessionName)
	}

	flashes := make([]string, 0)
	for _, f := range session.Flashes() {
		flashes = append(flashes, fmt.Sprintf("%v", f))
	}

	logger := logur.WithFields(ai.Logger, map[string]interface{}{"path": r.URL.Path, "method": r.Method})

	return RequestContext{
		AppInfo:     ai,
		Logger:      logger,
		Profile:     prof,
		Routes:      routes,
		Title:       util.AppName,
		Breadcrumbs: nil,
		Flashes:     flashes,
		Session:     *session,
	}
}

func ParseFlash(s string) (string, string) {
	split := strings.SplitN(s, ":", 2)
	severity := split[0]
	content := split[1]
	switch severity {
	case "status":
		return "uk-alert-primary", content
	case "success":
		return "uk-alert-success", content
	case "warning":
		return "uk-alert-warning", content
	case "error":
		return "uk-alert-danger", content
	default:
		return "", content
	}
}
