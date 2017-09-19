package cmd

import (
	"log"
	"net/http"
	"time"

	"blog/common/setting"
	"blog/router"

	"github.com/urfave/cli"
	"github.com/urfave/negroni"
)

const (
	DEFALUT_PORT            = "8443"
	DEFAULT_CONFIG_FILEPATH = "config/dev.ini"
)

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
	confFile := DEFAULT_CONFIG_FILEPATH
	if c.IsSet("port") {
		port = c.String("port")
	}
	if c.IsSet("config") {
		confFile = c.String("config")
	}

	setting.NewContext(confFile)

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

	if setting.SSL_ON && setting.RUN_MODE == "prod" {
		log.Fatal(s.ListenAndServe())
	} else {
		log.Fatal(s.ListenAndServeTLS(setting.CERT_FILE, setting.KEY_FILE))
	}
}
