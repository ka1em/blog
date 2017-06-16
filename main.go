package main

import (
	"blog.ka1em.site/router"
	"net/http"
	"os"
	"log"
	"runtime"
	"github.com/go-errors/errors"
)

const PORT  = ":8443"

func main() {
	r := router.InitRouters()

	http.ListenAndServeTLS(PORT, "keys/blog.ka1em.site/214098123750645.pem",
	"keys/blog.ka1em.site/214098123750645.key", r)
}

func init()  {
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
		log.Fatal(err.Error())
		return
	}
	return
}

