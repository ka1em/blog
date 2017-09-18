package controllers

import (
	"net/http"

	"blog/common"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	common.Suggar.Info("%+v", r.Form)
	return
}
