package router

import (
	"blog.ka1em.site/common"
	"blog.ka1em.site/controllers"
	"github.com/gorilla/mux"
)

func SetCommentRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/api/comments", controllers.APICommentPOST).Methods("POST")
	r.HandleFunc("/api/comments", controllers.APICommentGET).Methods("GET")
	r.HandleFunc("/api/comments/{id[\\w\\d=]+}", controllers.APICommentPUT).Methods("PUT")
	common.Suggar.Info("set comment routes ok")
	return r
}
