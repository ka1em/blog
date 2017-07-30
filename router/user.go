package router

import (
	"blog.ka1em.site/controllers"

	"github.com/gorilla/mux"
)

func SetUserRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/register", controllers.RegisterPost).Methods("POST")
	r.HandleFunc("/login", controllers.LoginPost).Methods("POST")
	return r
}
