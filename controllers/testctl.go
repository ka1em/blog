package controllers

import (
	"net/http"

	"blog/common/zlog"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	zlog.ZapLog.Info("%+v %s %s", r.Form, r.RemoteAddr, r.RequestURI)
	return
}
