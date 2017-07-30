package main

import (
	"log"
	"net/http"
	"os"
	"runtime"

	"blog.ka1em.site/common"
	"blog.ka1em.site/router"
	"github.com/go-errors/errors"
)

const PORT = ":8443"

func main() {
	r := router.InitRouters()

	//http.ListenAndServeTLS(PORT, "keys/blog.ka1em.site/214098123750645.pem",
	//	"keys/blog.ka1em.site/214098123750645.key", r)
	log.Println("Listening...")
	http.ListenAndServe(PORT, r)
}

func init() {
	var err error

	switch runtime.GOOS {
	case "darwin":
		err = os.Chdir("/Users/ka1em/go/src/blog.ka1em.site")
	case "linux":
		err = os.Chdir("/root/go/src/blog.ka1em.site")
	default:
		err = errors.New("Not darwin or linux")
	}

	if err != nil {
		common.Suggar.Error(err.Error())
		return
	}
	return
}
