package cmd

import (
	"blog/common/setting"
	"blog/common/zlog"
	"blog/controllers"
	"blog/model"
	"blog/router"
	"log"
	"net/http"
	"time"

	"github.com/urfave/cli"
	"github.com/urfave/negroni"
)

const (
	// DefaultPort 默认端口
	DefaultPort = "8443"
	// DefaultConfigFile 默认配置文件
	DefaultConfigFile = "conf/dev.ini"
)

// Web blog后端启动命令
var Web = cli.Command{
	Name:  "web",
	Usage: "Start web server",
	Description: `blog server is the only thing you need to run,
and it takes care of all the other things for you`,
	Action: runWeb,
	Flags: []cli.Flag{
		stringFlag("port, p", DefaultPort, "Port number, eg: 8443"),
		stringFlag("config, c", DefaultConfigFile, "Configuration file path"),
	},
}

func runWeb(c *cli.Context) {
	port := DefaultPort
	confFile := DefaultConfigFile
	if c.IsSet("port") {
		port = c.String("port")
	}
	if c.IsSet("config") {
		confFile = c.String("config")
	}

	setting.NewContext(confFile)
	zlog.ZapLogInit()
	model.DBInit()

	r := router.InitRouters()

	n := negroni.New()
	n.Use(negroni.NewStatic(http.Dir("static")))
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())

	n.UseHandler(r)
	n.UseFunc(controllers.ValidateSession)

	s := &http.Server{
		Addr:           ":" + port,
		Handler:        n,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	zlog.ZapLog.Info("blog listening ...")
	log.Fatal(s.ListenAndServe())
}
