package controllers

import (
	"context"
	"fmt"
	"net/http"

	"blog/common/zlog"
	"blog/model"
	"errors"
)

func ValidateSession(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	data := model.GetBaseData()

	session, err := model.SessionStore.Get(r, model.COOKIE_NAME)
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.MIDDLEWARE_ERR, http.StatusBadRequest)
		return
	}

	if sid, ok := session.Values["sid"]; ok {
		s := &model.Session{}
		if uid, ok := s.GetUserID(sid.(string), 1); !ok {
			zlog.ZapLog.Error(errors.New("record not found"))
			data.ResponseJson(w, model.NEED_LOGIN, http.StatusBadRequest)
			return
		} else {
			ctx := context.WithValue(r.Context(), "user_id", fmt.Sprintf("%d", uid))
			next(w, r.WithContext(ctx))
		}
	} else {
		zlog.ZapLog.Error("middleware need login")
		next(w, r)
	}
	return
}
