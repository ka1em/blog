package router

import (
	"blog/controllers"

	"blog/common/zlog"

	"github.com/gorilla/mux"
)

// SetTestRoutes 配置测试路由
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
	zlog.ZapLog.Info("set router test")
	return r
}
