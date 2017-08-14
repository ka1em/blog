package main

import (
	"net/http"

	"blog.ka1em.site/common"
	"blog.ka1em.site/router"
)

const PORT = ":8443"

func main() {
	r := router.InitRouters()

	//http.ListenAndServeTLS(PORT, "keys/blog.ka1em.site/214098123750645.pem",
	//	"keys/blog.ka1em.site/214098123750645.key", r)

	common.Suggar.Debug("Listening...")

	http.ListenAndServe(PORT, r)
}
