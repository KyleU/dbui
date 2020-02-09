package util

import (
	"github.com/gorilla/mux"
)

func ExtractRoutes(r *mux.Router) []string {
	ret := []string{}
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, _ := route.GetPathTemplate()
		ret = append(ret, pathTemplate)
		return nil
	})
	return ret
}
