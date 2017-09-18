package cmd

import (
	"flag"
	"log"
	"net/http"
	"time"

	"blog/router"
	"github.com/urfave/cli"
	"github.com/urfave/negroni"
)

const DEFALUT_PORT = ":8443"
const DEFAULT_CONFIG_FILEPATH = "conf/dev.ini"

var Web = cli.Command{
	Name:  "web",
	Usage: "Start web server",
	Description: `blog server is the only thing you need to run,
and it takes care of all the other things for you`,
	Action: runWeb,
	Flags: []cli.Flag{
		stringFlag("port, p", DEFALUT_PORT, "eg: :8443 Temporary port number to prevent conflict"),
		stringFlag("config, c", DEFAULT_CONFIG_FILEPATH, "Configuration file path"),
	},
}

func runWeb(c *cli.Context) error {
	set := flag.NewFlagSet("contrive", 0)
	nc := cli.NewContext(c.App, set, c)

	port := nc.String("port")
	r := router.InitRouters()

	n := negroni.Classic() // 导入一些预设的中间件
	n.UseHandler(r)

	s := &http.Server{
		Addr:           port,
		Handler:        n,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())

	//http.ListenAndServeTLS(PORT, "keys/blog/214098123750645.pem",
	//	"keys/blog/214098123750645.key", r)

	//common.Suggar.Debug("Listening...")
	//
	//http.ListenAndServe(PORT, r)
}
