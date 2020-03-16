package lib

import "github.com/kyleu/dbui/internal/app/cli"

func Run() {
	ai, err := cli.InitApp(cli.AppName, "0.0.0", "master")
	if err != nil {
		panic(err)
	}
	err = cli.MakeServer(ai, "127.0.0.1", 4200)
	if err != nil {
		panic(err)
	}
}
