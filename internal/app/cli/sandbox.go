package cli

import (
	"github.com/kyleu/dbui/internal/app/conn"
	"github.com/spf13/cobra"
)

func NewSandboxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sandbox",
		Aliases: []string{"x"},
		Short:   "Runs an internal test",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return conn.Sandbox("", "", "")
		},
	}

	return cmd
}
