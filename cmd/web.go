package cmd

import (
	"log"
	"net/http"
	"time"

	"blog/common/setting"
	"blog/common/zlog"
	"blog/model"
	"blog/router"

	"crypto/tls"

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
	model.DBInit()
	zlog.ZapLogInit()

	r := router.InitRouters()

	n := negroni.New()
	n.Use(negroni.NewStatic(http.Dir("static")))
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())

	n.UseHandler(r)
	s := &http.Server{
		Addr:           ":" + port,
		Handler:        n,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	switch setting.SSLMode {
	case false:
		log.Fatal(s.ListenAndServe())
	case true:
		var tlsMinVersion uint16
		switch setting.TLSMinVersion {
		case "SSL30":
			tlsMinVersion = tls.VersionSSL30
		case "TLS12":
			tlsMinVersion = tls.VersionTLS12
		case "TLS11":
			tlsMinVersion = tls.VersionTLS11
		case "TLS10":
			fallthrough
		default:
			tlsMinVersion = tls.VersionTLS10
		}

		s.TLSConfig = &tls.Config{
			MinVersion:               tlsMinVersion,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, // Required for HTTP/2 support.
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}

		log.Fatal(s.ListenAndServeTLS(setting.CertFile, setting.KeyFile))
	}
}
