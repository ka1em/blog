package cmd

import (
	"log"
	"net/http"
	"time"

	"blog/router"

	"github.com/urfave/cli"
	"github.com/urfave/negroni"
)

const DEFALUT_PORT = "8443"
const DEFAULT_CONFIG_FILEPATH = "conf/dev.ini"

var Web = cli.Command{
	Name:  "web",
	Usage: "Start web server",
	Description: `blog server is the only thing you need to run,
and it takes care of all the other things for you`,
	Action: runWeb,
	Flags: []cli.Flag{
		stringFlag("port, p", DEFALUT_PORT, "Port number, eg: 8443"),
		stringFlag("config, c", DEFAULT_CONFIG_FILEPATH, "Configuration file path"),
	},
}

func runWeb(c *cli.Context) {
	port := DEFALUT_PORT
	configFile := DEFAULT_CONFIG_FILEPATH
	if c.IsSet("port") {
		port = c.String("port")
	}

	if c.IsSet("config") {
		configFile = c.String("config")
	}




	r := router.InitRouters()

	n := negroni.Classic() // 导入一些预设的中间件
	n.UseHandler(r)

	s := &http.Server{
		Addr:           ":" + port,
		Handler:        n,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Listeng port", port)
	log.Fatal(s.ListenAndServe())

	//s.ListenAndServeTLS()
	//http.ListenAndServeTLS(PORT, "keys/blog/214098123750645.pem",
	//	"keys/blog/214098123750645.key", r)

	//common.Suggar.Debug("Listening...")
	//
	//http.ListenAndServe(PORT, r)
}
