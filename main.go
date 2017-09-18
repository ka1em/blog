package main

import (
	"os"

	"blog/cmd"
	"github.com/urfave/cli"
)

const APP_VER = "0.0.1"

func main() {
	app := cli.NewApp()
	app.Name = "blog"
	app.Usage = "blog backend"
	app.Version = APP_VER

	app.Commands = []cli.Command{
		cmd.Web,
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
