package web

import (
	"fmt"
	"github.com/kyleu/dbui/internal/app/controllers"
	"github.com/kyleu/dbui/internal/app/util"
	"net/http"
)

func MakeServer(info util.AppInfo, address string, port uint16) error {
	routes := controllers.BuildRouter(info)
	fmt.Println(fmt.Sprintf("[%v] starting on %v:%v (verbose)", info.AppName, address, port))
	err := http.ListenAndServe(fmt.Sprintf("%v:%v", address, port), routes)
	return err
}
