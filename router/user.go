package router

import (
	"blog/controllers"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func SetUserRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/user/register", controllers.RegisterPost).Methods("POST")
	r.HandleFunc("/user/login", controllers.LoginPost).Methods("POST")

	newRouer := mux.NewRouter().StrictSlash(false)
	newRouer.HandleFunc("/user/logout", controllers.LogoutGET).Methods("GET")

	r.PathPrefix("/user").Handler(negroni.New(
		negroni.HandlerFunc(controllers.ValidateSession),
		negroni.Wrap(newRouer),
	))
	return r
}
