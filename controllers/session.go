package controllers

import (
	"net/http"

	"blog.ka1em.site/model"
)

func ValidateSession(w http.ResponseWriter, r *http.Request) {
	session, _ := model.SessionStore.Get(r, "app-session")

	if sid, valid := session.Values["sid"]; valid {
		currentUID, _ := model.GetSessionUID(sid.(string))
		model.UpdateSession(sid.(string), currentUID)
		model.UserSession.UserId = currentUID
	} else {
		newSID, _ := model.GenerateSessionId()
		session.Values["sid"] = newSID
		session.Save(r, w)
		model.UserSession.SessionId = newSID
		model.UpdateSession(newSID, 0)
	}
}
