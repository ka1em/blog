package router

import (
	"blog/common/log"

	"github.com/gorilla/mux"
	"blog/common/setting"
)

func InitRouters() *mux.Router {
	setting.NewContext()

	r := mux.NewRouter().StrictSlash(false)

	r = SetPageRoutes(r)
	r = SetUserRoutes(r)
	r = SetCommentRoutes(r)
	r = SetTestRoutes(r)

	log.Suggar.Info("set route ok ")
	return r
}
