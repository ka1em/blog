package router

import (
	"blog/common/zlog"
	"blog/controllers"

	"github.com/gorilla/mux"
)

// SetPageRoutes 设置page路由
func SetPageRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/api/pages", controllers.PageIndexGET).Methods("GET")
	r.HandleFunc("/api/pages/{id:[0-9a-zA\\-]+}", controllers.APIPageGET).Methods("GET")
	r.HandleFunc("/api/pages", controllers.APIPagePOST).Methods("POST")
	zlog.ZapLog.Info("set router page")
	return r
}
