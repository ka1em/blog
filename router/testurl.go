package router

import (
	"blog.ka1em.site/controllers"
	"github.com/gorilla/mux"
)

func SetTestRoutes(r *mux.Router) *mux.Router {

	r.HandleFunc("/test", controllers.TestHandler).Methods("GET")
	//r.HandleFunc("/user/login", controllers.LoginPost).Methods("POST")

	//newRouer := mux.NewRouter().StrictSlash(false)
	//newRouer.HandleFunc("/user/logout", controllers.LogoutGET).Methods("GET")
	//
	//r.PathPrefix("/user").Handler(negroni.New(
	//	negroni.HandlerFunc(controllers.ValidateSession),
	//	negroni.Wrap(newRouer),
	//))
	return r
}
