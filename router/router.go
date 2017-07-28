package router

import (
	"blog.ka1em.site/common"
	"github.com/gorilla/mux"
)

func InitRouters() *mux.Router {
	r := mux.NewRouter().StrictSlash(false)
	r = SetPageRoutes(r)
	r = SetUserRoutes(r)
	common.Suggar.Info("set route ok ")
	return r
}
