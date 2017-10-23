package controllers

import (
	"context"
	"fmt"
	"net/http"

	"blog/common/zlog"
	"blog/model"
	"errors"
)

// ValidateSession 验证session
func ValidateSession(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	data := model.GetBaseData()

	session, err := model.SessionStore.Get(r, model.CookieName)
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.MiddlewareErr, http.StatusBadRequest)
		return
	}

	if sid, ok := session.Values["sid"]; ok {
		s := &model.Session{}
		var uid int64
		var ok bool
		if uid, ok = s.GetUserID(sid.(string), 1); !ok {
			zlog.ZapLog.Error(errors.New("record not found"))
			data.ResponseJson(w, model.NeedLogin, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", fmt.Sprintf("%d", uid))
		next(w, r.WithContext(ctx))

	} else {
		zlog.ZapLog.Error("middleware need login")
		next(w, r)
	}
}
