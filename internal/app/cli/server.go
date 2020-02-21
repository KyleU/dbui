package cli

import (
	"fmt"
	"net/http"

	"emperror.dev/errors"
	"github.com/gorilla/handlers"
	"github.com/kyleu/dbui/internal/app/controllers"
	"github.com/kyleu/dbui/internal/app/util"
	"github.com/spf13/cobra"
)

func NewServerCommand(appName string, version string, commitHash string) *cobra.Command {
	var port uint16
	var addr string

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Starts the http server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			info := InitApp(appName, version, commitHash)
			return makeServer(info, addr, port)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&addr, "address", "a", "127.0.0.1", "interface address to listen on")
	flags.Uint16VarP(&port, "port", "p", 4200, "port for http server to listen on")

	return cmd
}

func makeServer(info util.AppInfo, address string, port uint16) error {
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
