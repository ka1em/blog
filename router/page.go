package router

import (
	"blog/common/log"
	"blog/controllers"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func SetPageRoutes(r *mux.Router) *mux.Router {

	r.HandleFunc("/api/pages", controllers.PageIndexGET).Methods("GET")
	r.HandleFunc("/api/pages/{id:[0-9a-zA\\-]+}", controllers.APIPageGET).Methods("GET")

	newRouter := mux.NewRouter().StrictSlash(false)
	newRouter.HandleFunc("/api/pages", controllers.APIPagePOST).Methods("POST")

	r.PathPrefix("/api/pages").Handler(negroni.New(
		negroni.HandlerFunc(controllers.ValidateSession),
		negroni.Wrap(newRouter),
	))

	log.Suggar.Info("set page routes ok")
	return r
}
