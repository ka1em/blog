package controllers

import (
	"net/http"

	"blog/common/zlog"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		zlog.ZapLog.Error(err.Error())
	}

	zlog.ZapLog.Info("%+v %s %s", r.Form, r.RemoteAddr, r.RequestURI)
}
