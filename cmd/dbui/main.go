package main

import (
	"github.com/kyleu/dbui/internal/app/cli"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// Provisioned by ldflags
// nolint: gochecknoglobals
var (
	version    string
	commitHash string
)

func main() {
	rootCmd := &cobra.Command{
		Use: "dbui",
		Short: "Command line interface for dbui",
		Long: "A work in progress...",
	}

	cli.Configure(rootCmd, version, commitHash)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
