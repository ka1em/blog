package router

import (
	"blog.ka1em.site/controllers"

	"github.com/gorilla/mux"
)

func SetUserRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/register", controllers.RegisterPOST).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	return r
}
