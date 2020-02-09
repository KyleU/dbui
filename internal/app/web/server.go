package web

import (
	"fmt"
	"github.com/kyleu/dbui/internal/app/controllers"
	"github.com/kyleu/dbui/internal/app/util"
	"net/http"
)

func MakeServer(info util.AppInfo, address string, port uint16) error {
	routes := controllers.BuildRouter(info)
	return http.ListenAndServe(fmt.Sprintf("%v:%v", address, port), routes)
}
