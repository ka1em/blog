package router

import (
	"blog/common/zlog"

	"github.com/gorilla/mux"
)

func InitRouters() *mux.Router {
	r := mux.NewRouter().StrictSlash(false)

	r = SetPageRoutes(r)
	r = SetUserRoutes(r)
	r = SetCommentRoutes(r)
	r = SetTestRoutes(r)

	zlog.ZapLog.Info("set route ok ")
	return r
}
