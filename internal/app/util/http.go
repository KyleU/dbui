package util

import (
	"emperror.dev/emperror"
	"emperror.dev/errors"
	"github.com/gorilla/mux"
	"logur.dev/logur"
	"net/http"
)

func NotFound(res http.ResponseWriter, _ *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	res.WriteHeader(http.StatusNotFound)
	_, err := res.Write([]byte("404: Page Not Found"))
	emperror.Panic(errors.Wrap(err, "Unable to write to response"))
}

type AppInfo struct {
	AppName      string
	Debug        bool
	Version      string
	CommitHash   string
	Logger       logur.LoggerFacade
	ErrorHandler emperror.ErrorHandlerFacade
}

type RequestContext struct {
	AppInfo AppInfo
	Profile UserProfile
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
	prof := SystemProfile

	if title == "" {
		title = ai.AppName
	}

	return RequestContext{
		AppInfo: ai,
		Profile: prof,
		Routes:  r,
		Title:   title,
	}
}
