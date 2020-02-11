package cli

import (
	"github.com/kyleu/dbui/internal/app/web"
	"github.com/spf13/cobra"
)


func NewServerCommand(appName string, version string, commitHash string) *cobra.Command {
	var port uint16
	var addr string

	cmd := &cobra.Command{
		Use:     "server",
		Short:   "Starts the http server",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			info := InitApp(appName, version, commitHash)
			return web.MakeServer(info, addr, port)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&addr, "address", "a", "127.0.0.1", "interface address to listen on")
	flags.Uint16VarP(&port, "port", "p", 4200, "port for http server to listen on")

	return cmd
}
