package router

import (
	"blog.ka1em.site/controllers"

	"blog.ka1em.site/common"
	"github.com/gorilla/mux"
)

func SetPageRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/page/{guid:[0-9a-zA\\-]+}", controllers.ServePage).Methods("GET")
	r.HandleFunc("/", controllers.RedirIndex).Methods("GET")
	r.HandleFunc("/home", controllers.ServeIndex).Methods("GET")
	r.HandleFunc("/api/pages", controllers.ServeIndex).Methods("GET")
	r.HandleFunc("/api/pages/{guid:[0-9a-zA\\-]+}", controllers.APIPage).Methods("GET")

	common.Suggar.Info("set page routes ok")
	return r
}
