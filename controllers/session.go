package controllers

import (
	"context"
	"fmt"
	"net/http"

	zlog "blog/common/log"
	"blog/model"
)

func ValidateSession(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	data := model.GetBaseData()
	sessionStore := model.GetSessionStore()

	session, err := sessionStore.Get(r, "app-session")
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.MIDDLEWAREERR, http.StatusBadRequest)
		return
	}

	if sid, ok := session.Values["sid"]; ok {
		if uid, err := model.GetUserID(sid.(string), 1); err != nil {
			if err.Error() == "record not found" {
				zlog.ZapLog.Error(err.Error())
				data.ResponseJson(w, model.NEEDLOGIN, http.StatusOK)
				return
			}
			zlog.ZapLog.Error(err.Error())
			data.ResponseJson(w, model.MIDDLEWAREERR, http.StatusOK)
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
