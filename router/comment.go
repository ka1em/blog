package router

import (
	"blog/controllers"

	"github.com/gorilla/mux"
)

// SetCommentRoutes 设置评论路由
func SetCommentRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/api/comments", controllers.APICommentGET).Methods("GET")
	r.HandleFunc("/api/comments", controllers.APICommentPOST).Methods("POST")
	r.HandleFunc("/api/comments/{id[\\w\\d=]+}", controllers.APICommentPUT).Methods("PUT")
	return r
}
