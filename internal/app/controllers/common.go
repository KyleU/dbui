package controllers

import (
	"github.com/kyleu/dbui/internal/app/util"
	"net/http"
)

func prepHtml(res http.ResponseWriter, req *http.Request, title string) util.RequestContext {
	res.Header().Set("Content-Type", "text/html")
	res.WriteHeader(http.StatusOK)
	return util.ExtractContext(req, title)
}
