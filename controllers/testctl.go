package controllers

import (
	"net/http"

	zlog "blog/common/log"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	zlog.ZapLog.Info("%+v", r.Form)
	return
}
