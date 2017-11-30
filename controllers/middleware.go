package controllers

import (
	"context"
	"net/http"

	"blog/common/zlog"
	"blog/model"
)

func ValidateSession(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	data := model.GetBaseData()
	sst, _ := model.SessionStore.Get(r, model.CookieName)

	// 获取session， 如果是新生成的则给客户端一个新的session
	if sid, ok := sst.Values["sid"]; ok {
		zlog.ZapLog.Debug("middlerware ok sid: ", sid)
		s := model.Session{
			SID: sid.(string),
		}

		// 忽略错误，有可能没有user_id
		se, _ := s.GetCache()

		s.UserID = se.UserID

		ctx := context.WithValue(r.Context(), model.UserSession, s)
		next(w, r.WithContext(ctx))
	} else {
		zlog.ZapLog.Debug("middlerware false sid: ", sid)
		sid, err := model.GenerateSessionID()
		if err != nil {
			zlog.ZapLog.Error(err.Error())
			data.ResponseJson(w, model.MiddlewareErr, http.StatusInternalServerError)
			return
		}
		sst.Values["sid"] = sid
		se := &model.Session{
			SID:    sid,
			UserID: 0,
		}
		if err := sst.Save(r, w); err != nil {
			zlog.ZapLog.Error(err.Error())
			data.ResponseJson(w, model.MiddlewareErr, http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), model.UserSession, se)
		next(w, r.WithContext(ctx))
	}
}
