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
	Debug bool
	Version string
	CommitHash string
	ConfigDir string
}

type RequestContext struct {
	AppInfo AppInfo
	Routes  *mux.Router
}

func (r *RequestContext) Route(act string, pairs ...string) string {
	url, err := r.Routes.Get(act).URL(pairs...)
	if err != nil {
		return "/noroute"
	}
	return url.Path
}

func ExtractContext(req *http.Request) RequestContext {
	r := req.Context().Value("routes")
	ai := req.Context().Value("info")
	return RequestContext {
		AppInfo: ai.(AppInfo),
		Routes: r.(*mux.Router),
	}
}
