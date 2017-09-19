package controllers

import (
	"net/http"

	"blog/common/zlog"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	zlog.ZapLog.Info("%+v", r.Form)
	return
}
