package web

import (
	"emperror.dev/errors"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/kyleu/dbui/internal/app/controllers"
	"github.com/kyleu/dbui/internal/app/util"
	"net/http"
)

func MakeServer(info util.AppInfo, address string, port uint16) error {
	routes, err := controllers.BuildRouter(info)
	if err != nil {
		return errors.WithMessage(err, "Unable to construct routes")
	}
	var msg = fmt.Sprintf("%v is starting", info.AppName)
	if info.Debug {
		msg = msg + " (verbose)"
	}
	info.Logger.Info(msg, map[string]interface{} { "address": address, "port": port})
	err = http.ListenAndServe(fmt.Sprintf("%v:%v", address, port),  handlers.CORS()(routes))
	return errors.Wrap(err, "Unable to run http server")
}
