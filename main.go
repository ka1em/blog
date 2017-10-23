package main

import (
	"blog/cmd"
	"os"

	"github.com/urfave/cli"
)

const appVer = "0.0.1"

func main() {
	app := cli.NewApp()
	app.Name = "blog"
	app.Usage = "blog backend"
	app.Version = appVer
	app.Commands = []cli.Command{
		cmd.Web,
		cmd.WebAdmin,
	}
	//app.Flags = append(app.Flags, cmd.Web.Flags...)
	app.Run(os.Args)
}
