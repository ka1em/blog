package controllers

import (
	"net/http"

	"blog/common/zlog"
	"blog/model"
)

var Count = 0

// TestHandler 测试handler
func TestHandler(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()
	if err := r.ParseForm(); err != nil {
		zlog.ZapLog.Error(err.Error())
	}
	Count++
	data.Data["count"] = Count
	zlog.ZapLog.Info("%+v %s %s", r.Form, r.RemoteAddr, r.RequestURI)
	data.Data["remoteAddr"] = r.RemoteAddr
	data.Data["requestURI"] = r.RequestURI
	data.Data["body"] = r.Body
	data.ResponseJson(w, model.Success, http.StatusOK)
}
