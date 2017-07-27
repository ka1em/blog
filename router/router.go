package router

import "github.com/gorilla/mux"

func InitRouters() *mux.Router {
	r := mux.NewRouter().StrictSlash(false)
	r = SetPageRoutes(r)
	r = SetUserRoutes(r)
	return r
}
