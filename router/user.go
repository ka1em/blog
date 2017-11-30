package router

import (
	"blog/controllers"

	"blog/common/zlog"

	"github.com/gorilla/mux"
)

// SetUserRoutes 设置用户路由
func SetUserRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/user/register", controllers.RegisterPost).Methods("POST")
	r.HandleFunc("/user/login", controllers.LoginPost).Methods("POST")
	r.HandleFunc("/user/logout", controllers.LogoutGET).Methods("GET")
	zlog.ZapLog.Info("set router user")
	return r
}
