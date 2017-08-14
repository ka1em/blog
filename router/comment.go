package router

import (
	"blog.ka1em.site/common"
	"blog.ka1em.site/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func SetCommentRoutes(r *mux.Router) *mux.Router {
	newRouter := mux.NewRouter().StrictSlash(false)
	newRouter.HandleFunc("/api/comments", controllers.APICommentPOST).Methods("POST")
	newRouter.HandleFunc("/api/comments", controllers.APICommentGET).Methods("GET")
	newRouter.HandleFunc("/api/comments/{id[\\w\\d=]+}", controllers.APICommentPUT).Methods("PUT")

	r.PathPrefix("/api/comments").Handler(negroni.New(
		negroni.HandlerFunc(controllers.ValidateSession),
		negroni.Wrap(newRouter),
	))

	common.Suggar.Info("set comment routes ok")
	return r
}