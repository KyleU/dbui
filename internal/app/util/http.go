package util

import (
	"emperror.dev/emperror"
	"emperror.dev/errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	"os"
	"strings"
)

func NotFound(res http.ResponseWriter, _ *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	res.WriteHeader(http.StatusNotFound)
	_, err := res.Write([]byte("404: Page Not Found"))
	emperror.Panic(errors.Wrap(err, "Unable to write to response"))
}

type RequestContext struct {
	AppInfo AppInfo
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

func ExtractContext(req *http.Request, title string) RequestContext {
	ai := req.Context().Value("info").(AppInfo)
	r := req.Context().Value("routes").(*mux.Router)
	prof := SystemProfile
	session, err := store.Get(req, sessionName)
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

	return RequestContext{
		AppInfo: ai,
		Profile: prof,
		Routes:  r,
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
