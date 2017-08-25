package main

import (
	"log"
	"net/http"

	"time"

	"blog.ka1em.site/router"
	"github.com/urfave/negroni"
)

const PORT = ":8443"

func main() {
	r := router.InitRouters()

	n := negroni.Classic() // 导入一些预设的中间件
	n.UseHandler(r)

	s := &http.Server{
		Addr:           PORT,
		Handler:        n,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())

	//http.ListenAndServeTLS(PORT, "keys/blog.ka1em.site/214098123750645.pem",
	//	"keys/blog.ka1em.site/214098123750645.key", r)

	//common.Suggar.Debug("Listening...")
	//
	//http.ListenAndServe(PORT, r)
}
