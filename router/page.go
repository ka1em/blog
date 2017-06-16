package router

import (
	"blog.ka1em.site/controllers"

	"github.com/gorilla/mux"
)

func SetPageRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/page/{guid:[0-9a-zA\\-]+}", controllers.ServePage).Methods("GET")
	return r
}
