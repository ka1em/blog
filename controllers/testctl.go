package controllers

import (
	"net/http"
	"sync"

	"blog/common/zlog"
	"blog/model"
)

type Counter struct {
	Count int
	m     *sync.RWMutex
}

var Count Counter

func init() {
	Count.m = new(sync.RWMutex)
}

// TestHandler 测试handler
func TestHandler(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()
	if err := r.ParseForm(); err != nil {
		zlog.ZapLog.Error(err.Error())
	}

	Count.m.Lock()
	Count.Count++
	Count.m.Unlock()

	Count.m.RLock()
	data.Data["count"] = Count.Count
	Count.m.RUnlock()

	zlog.ZapLog.Infof("%+v %s %s", r.Form, r.RemoteAddr, r.RequestURI)
	data.Data["remoteAddr"] = r.RemoteAddr
	data.Data["requestURI"] = r.RequestURI
	data.Data["body"] = r.Body
	data.ResponseJson(w, model.Success, http.StatusOK)
}
