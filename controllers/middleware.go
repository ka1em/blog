package controllers

import (
	"context"
	"fmt"
	"net/http"

	"blog/common/zlog"
	"blog/model"
)

// ValidateSession 验证session
//func ValidateSession(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
//	data := model.GetBaseData()
//
//	session, err := model.SessionStore.Get(r, model.CookieName)
//	if err != nil {
//		zlog.ZapLog.Error(err.Error())
//		data.ResponseJson(w, model.MiddlewareErr, http.StatusBadRequest)
//		return
//	}
//
//	if sid, ok := session.Values["sid"]; ok {
//		s := &model.Session{}
//		var uid uint64
//		var ok bool
//		if uid, ok = s.GetUserID(sid.(string), 1); !ok {
//			zlog.ZapLog.Error(model.ErrMap[model.SessionNoUserID])
//			data.ResponseJson(w, model.SessionNoUserID, http.StatusBadRequest)
//			return
//		}
//
//		ctx := context.WithValue(r.Context(), "USER_ID", fmt.Sprintf("%d", uid))
//		zlog.ZapLog.Debug("middler ware: logined")
//		next(w, r.WithContext(ctx))
//
//	} else {
//		zlog.ZapLog.Debug("middler ware: not logined")
//		next(w, r)
//	}
//}

func ValidateSession(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	data := model.GetBaseData()
	sst, err := model.SessionStore.Get(r, model.CookieName)
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.MiddlewareErr, http.StatusBadRequest)
		return
	}
	if sid, ok := sst.Values["sid"]; ok {
		s := model.Session{
			SID: sid.(string),
		}
		se, err := s.Get()
		if err != nil {
			// TODO err
			return
		}
		ctx := context.WithValue(r.Context(), "USER_ID", fmt.Sprintf("%d", se.UserID))
		//ctx := context.WithValue(r.Context(), "SID", se.SID)
		next(w, r.WithContext(ctx))
		zlog.ZapLog.Debug("middleware session has id ")
	} else {
		sid, err := model.GenerateSessionID()
		if err != nil {
			// todo
			return
		}
		sst.Values["sid"] = sid
		sst.Save(r, w)
		next(w, r)
		zlog.ZapLog.Debug("middleware session no id ")
	}
}
