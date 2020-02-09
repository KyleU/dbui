package cli

import (
	"github.com/kyleu/dbui/internal/app/util"
	"github.com/kyleu/dbui/internal/app/web"
	"github.com/spf13/cobra"
)

func NewServerCommand(info util.AppInfo) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "server",
		Short:   "Starts the http server",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return web.MakeServer(info)
		},
	}

	return cmd
}
