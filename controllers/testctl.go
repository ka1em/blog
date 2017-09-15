package controllers

import (
	"net/http"

	"blog.ka1em.site/common"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	common.Suggar.Info("%+v", r.Form)
	return
}
