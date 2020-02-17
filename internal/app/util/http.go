package util

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"logur.dev/logur"
	"net/http"
	"os"
	"strings"
)

type RequestContext struct {
	AppInfo AppInfo
	Logger  logur.LoggerFacade
	Profile UserProfile
	Routes  *mux.Router
	Title   string
	Flashes []string
	Session sessions.Session
}

func (r *RequestContext) Route(act string, pairs ...string) string {
	route := r.Routes.Get(act)
	if route == nil {
		fmt.Println("Cannot find route at path [" + act + "]")
		return "/routenotfound"
	}
	url, err := route.URL(pairs...)
	if err != nil {
		fmt.Println("Cannot bind route at path [" + act + "]")
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

const sessionName = "dbui-session"

func ExtractContext(r *http.Request, title string) RequestContext {
	ai := r.Context().Value("info").(AppInfo)
	routes := r.Context().Value("routes").(*mux.Router)
	prof := SystemProfile
	session, err := store.Get(r, sessionName)
	if err != nil {
		session = sessions.NewSession(store, sessionName)
	}

	if title == "" {
		title = ai.AppName
	}

	flashes := make([]string, 0)
	for _, f := range session.Flashes() {
		flashes = append(flashes, fmt.Sprintf("%v", f))
	}

	logger := logur.WithFields(ai.Logger, map[string]interface{}{ "path": r.URL.Path, "method": r.Method})

	return RequestContext{
		AppInfo: ai,
		Logger:  logger,
		Profile: prof,
		Routes:  routes,
		Title:   title,
		Flashes: flashes,
		Session: *session,
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
