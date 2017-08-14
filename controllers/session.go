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
	session, _ := model.SessionStore.Get(r, "app-session")
	s := &model.Session{}
	if sid, valid := session.Values["sid"]; valid {
		s.SessionId = sid.(string)
		common.Suggar.Debugf("validate session session id  = %s", s.SessionId)
		if err := s.GetSessionUID(); err != nil {
			if err.Error() == "record not found" {
				common.Suggar.Error(err.Error())
				data.ResponseJson(w, common.NEEDLOGIN, http.StatusOK)
				return
			}
			common.Suggar.Error(err.Error())
			data.ResponseJson(w, common.MIDDLEWAREERR, http.StatusOK)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", fmt.Sprintf("%d", s.UserId))

		common.Suggar.Debugf("validate session user_id = %d", s.UserId)
		next(w, r.WithContext(ctx))

	} else {
		next(w, r)
	}
	return
}
