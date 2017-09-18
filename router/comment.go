package router

import (
	"blog/common/log"
	"blog/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func SetCommentRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/api/comments", controllers.APICommentGET).Methods("GET")

	newRouter := mux.NewRouter().StrictSlash(false)
	newRouter.HandleFunc("/api/comments", controllers.APICommentPOST).Methods("POST")
	newRouter.HandleFunc("/api/comments/{id[\\w\\d=]+}", controllers.APICommentPUT).Methods("PUT")

	r.PathPrefix("/api/comments").Handler(negroni.New(
		negroni.HandlerFunc(controllers.ValidateSession),
		negroni.Wrap(newRouter),
	))

	log.Suggar.Info("set comment routes ok")
	return r
}
