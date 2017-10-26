package controllers

import (
	"context"
	"fmt"
	"net/http"

	"blog/common/zlog"
	"blog/model"
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
			zlog.ZapLog.Error(model.ErrMap[model.SessionNoUserID])
			data.ResponseJson(w, model.SessionNoUserID, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", fmt.Sprintf("%d", uid))
		next(w, r.WithContext(ctx))

	} else {
		zlog.ZapLog.Error(model.ErrMap[model.NeedLogin])
		next(w, r)
	}
}
