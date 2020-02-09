package web

import (
	"github.com/kyleu/dbui/internal/app/controllers"
	"github.com/kyleu/dbui/internal/app/util"
	"net/http"
)

func MakeServer(info util.AppInfo) error {
	routes := controllers.BuildRouter(info)
	return http.ListenAndServe(":4200", routes)
}
