package cmd

import "github.com/urfave/cli"

const DEFAULT_ADMIN_PORT = "9000"

var WebAdmin = cli.Command{
	Name:  "webadmin",
	Usage: "Start web admin server",
	Description: `blog server is the only thing you need to run,
and it takes care of all the other things for you`,
	Action: runWebAdmin,
	Flags: []cli.Flag{
		stringFlag("port, p", DEFAULT_ADMIN_PORT, "Port number, eg: 8443"),
		stringFlag("config, c", DEFAULT_CONFIG_FILEPATH, "Configuration file path"),
	},
}

func runWebAdmin(c *cli.Context) {

}
