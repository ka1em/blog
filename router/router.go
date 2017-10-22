package router

import (
	"github.com/gorilla/mux"
)

// InitRouters 初始化路由
func InitRouters() *mux.Router {
	r := mux.NewRouter().StrictSlash(false)
	r = SetPageRoutes(r)
	r = SetUserRoutes(r)
	r = SetCommentRoutes(r)
	r = SetTestRoutes(r)
	return r
}
