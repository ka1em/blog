package controllers

import (
	"net/http"

	"blog/common/log"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	log.Suggar.Info("%+v", r.Form)
	return
}
