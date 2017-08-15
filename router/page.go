package router

import (
	"blog.ka1em.site/controllers"

	"blog.ka1em.site/common"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func SetPageRoutes(r *mux.Router) *mux.Router {
	//r.HandleFunc("/page/{guid:[0-9a-zA\\-]+}", controllers.ServePage).Methods("GET")
	//r.HandleFunc("/", controllers.RedirIndex).Methods("GET")
	//r.HandleFunc("/home", controllers.ServeIndex).Methods("GET")
	r.HandleFunc("/api/pages", controllers.PageIndexGET).Methods("GET")
	r.HandleFunc("/api/pages/{id:[0-9a-zA\\-]+}", controllers.APIPageGET).Methods("GET")

	newRouter := mux.NewRouter().StrictSlash(false)
	newRouter.HandleFunc("/api/pages", controllers.APIPagePOST).Methods("POST")

	r.PathPrefix("/api/pages").Handler(negroni.New(
		negroni.HandlerFunc(controllers.ValidateSession),
		negroni.Wrap(newRouter),
	))

	common.Suggar.Info("set page routes ok")
	return r
}
