package controllers

import (
	"context"
	"fmt"
	"net/http"

	"blog.ka1em.site/common"
	"blog.ka1em.site/model"
)

func ValidateSession(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	data := model.GetBaseData()
	sessionStore := model.GetSessionStore()

	session, err := sessionStore.Get(r, "app-session")
	if err != nil {
		common.Suggar.Error(err.Error())
		data.ResponseJson(w, model.MIDDLEWAREERR, http.StatusBadRequest)
		return
	}

	if sid, ok := session.Values["sid"]; ok {
		if uid, err := model.GetUserID(sid.(string)); err != nil {
			if err.Error() == "record not found" {
				common.Suggar.Error(err.Error())
				data.ResponseJson(w, model.NEEDLOGIN, http.StatusOK)
				return
			}
			common.Suggar.Error(err.Error())
			data.ResponseJson(w, model.MIDDLEWAREERR, http.StatusOK)
			return
		} else {
			ctx := context.WithValue(r.Context(), "user_id", fmt.Sprintf("%d", uid))
			next(w, r.WithContext(ctx))
		}

	} else {
		common.Suggar.Error("middleware need login")

		//next(w, r)
	}
	return
}
