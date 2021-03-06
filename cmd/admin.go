package cmd

import "github.com/urfave/cli"

// DefaultAdminPort 默认后台端口
const DefaultAdminPort = "9000"

// WebAdmin 管理后台启动命令
var WebAdmin = cli.Command{
	Name:  "webadmin",
	Usage: "Start web admin server",
	Description: `blog server is the only thing you need to run,
and it takes care of all the other things for you`,
	Action: runWebAdmin,
	Flags: []cli.Flag{
		stringFlag("port, p", DefaultAdminPort, "Port number, eg: 8443"),
		stringFlag("config, c", DefaultConfigFile, "Configuration file path"),
	},
}

func runWebAdmin(c *cli.Context) {

}
