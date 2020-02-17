package web

import (
	"fmt"
	"net/http"

	"emperror.dev/errors"
	"github.com/gorilla/handlers"
	"github.com/kyleu/dbui/internal/app/controllers"
	"github.com/kyleu/dbui/internal/app/util"
)

func MakeServer(info util.AppInfo, address string, port uint16) error {
	routes, err := controllers.BuildRouter(info)
	if err != nil {
		return errors.WithStack(errors.WithMessage(err, "Unable to construct routes"))
	}
	var msg = fmt.Sprintf("%v is starting on [%v:%v]", info.AppName, address, port)
	if info.Debug {
		msg += " (verbose)"
	}
	info.Logger.Info(msg, map[string]interface{}{"address": address, "port": port})
	err = http.ListenAndServe(fmt.Sprintf("%v:%v", address, port), handlers.CORS()(routes))
	return errors.WithStack(errors.Wrap(err, "Unable to run http server"))
}
