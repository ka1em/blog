package router

import (
	"github.com/gorilla/mux"
)

func SetWeChatRoutes(r *mux.Router) *mux.Router {
	//r.HandleFunc("/api/v1/reponse/wxtoken", controllers.WeChatValidGET).Methods("GET")
	//
	//zlog.ZapLog.Info("set router page")
	return r
}
