package router

import "github.com/gorilla/mux"

func InitRouters() *mux.Router {
	r := mux.NewRouter().StrictSlash(false)
	r = SetPageRoutes(r)
	return r
}
