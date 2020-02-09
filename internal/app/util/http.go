package util

import (
	"github.com/gorilla/mux"
	"net/http"
)

func NotFound(res http.ResponseWriter, _ *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	res.WriteHeader(http.StatusNotFound)
	_, err := res.Write([]byte("404: Page Not Found"))
	if err != nil {
		panic(err)
	}
}

type AppInfo struct {
	AppName    string
	Debug      bool
	Version    string
	CommitHash string
}

type RequestContext struct {
	AppInfo AppInfo
	Routes  *mux.Router
	Title   string
}

func (r *RequestContext) Route(act string, pairs ...string) string {
	url, err := r.Routes.Get(act).URL(pairs...)
	if err != nil {
		return "/noroute"
	}
	return url.Path
}

func ExtractContext(req *http.Request, title string) RequestContext {
	ai := req.Context().Value("info").(AppInfo)
	r := req.Context().Value("routes").(*mux.Router)
	if title == "" {
		title = ai.AppName
	}
	return RequestContext{
		AppInfo: ai,
		Routes:  r,
		Title:   title,
	}
}
