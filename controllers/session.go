package controllers

import (
	"net/http"

	"blog.ka1em.site/common"
	"blog.ka1em.site/model"
	"github.com/gorilla/context"
)

func ValidateSession(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	data := model.GetBaseData()
	session, _ := model.SessionStore.Get(r, "app-session")
	s := &model.Session{}
	if sid, valid := session.Values["sid"]; valid {
		s.SessionId = sid.(string)

		if err := s.GetSessionUID(); err != nil {
			common.Suggar.Error(err.Error())
			data.ResponseJson(w, common.MIDDLEWAREERR, http.StatusOK)
			return
		}

		context.Set(r, "user_id", s.UserId)
		next(w, r)
		return
	}

	http.Redirect(w, r, "/login", 301)
	return
}
