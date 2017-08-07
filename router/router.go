package router

import (
	"net/http"

	"blog.ka1em.site/common"
	"github.com/gorilla/mux"
)

func InitRouters() *mux.Router {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "/tmp/www/index.html")
	})

	r = SetPageRoutes(r)
	r = SetUserRoutes(r)
	r = SetCommentRoutes(r)

	common.Suggar.Info("set route ok ")
	return r
}
