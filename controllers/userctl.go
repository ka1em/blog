package controllers

import (
	"net/http"
	"regexp"

	"crypto/sha256"
	"encoding/hex"

	"blog.ka1em.site/common"
	"blog.ka1em.site/model"
)

func RegisterPost(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()
	if err := r.ParseForm(); err != nil {
		data.ResponseJson(w, common.USER_PARSEFORM, http.StatusBadRequest)
		common.Suggar.Error(err.Error())
		return
	}

	name := r.PostFormValue("user_name")
	email := r.PostFormValue("user_email")
	passwd := r.PostFormValue("user_passwd")

	if name == "" || email == "" || passwd == "" {
		data.ResponseJson(w, common.USER_PARAMVALID, http.StatusBadRequest)
		common.Suggar.Error("register params is null")
		return
	}

	common.Suggar.Debugf("%+v", passwd)

	gure := regexp.MustCompile("[^A-Za-z0-9]+")
	guid := gure.ReplaceAllString(name, "")

	u := &model.User{}
	u.UserName = name
	u.UserEmail = email
	u.UserGuid = guid

	u.UserSalt = guid + "QWERdsfawer2314=="
	pass := passwordHash(passwd, u.UserSalt)

	u.UserPasswd = pass

	common.Suggar.Debug("%+v", u)

	//创建用户
	if err := u.CreateUser(); err != nil {
		if err.Error() == "exists" {
			common.Suggar.Error(err.Error())
			data.ResponseJson(w, common.USER_WASEXIST, http.StatusBadRequest)
			return
		}

		common.Suggar.Error(err.Error())
		data.ResponseJson(w, common.DATA_CREATEUSER, http.StatusInternalServerError)
		return
	}

	data.ResponseJson(w, common.SUCCESS, http.StatusOK)
	return
}

func passwordHash(p, solt string) string {
	hash := sha256.New()
	s := p + solt
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

//
func LoginPost(w http.ResponseWriter, r *http.Request) {
	u := &model.User{}
	data := model.GetBaseData()

	s := &model.Session{}
	s.CreateSeesion(w, r)

	if err := r.ParseForm(); err != nil {
		data.ResponseJson(w, common.USER_PARSEFORM, http.StatusBadRequest)
		common.Suggar.Error(err.Error())
		return
	}

	name := r.PostFormValue("user_name")
	passwd := r.PostFormValue("user_passwd")

	gure := regexp.MustCompile("[^A-Za-z0-9]+")
	guid := gure.ReplaceAllString(name, "")

	salt := guid + "QWERdsfawer2314=="

	u.UserName = name
	u.UserPasswd = passwordHash(passwd, salt)

	if notfound := u.Login(); notfound {
		common.Suggar.Error("log err")
		data.ResponseJson(w, common.USER_PARSEFORM, http.StatusInternalServerError)
		return
	}

	s.UserId = u.ID
	if err := s.UpdateSession(); err != nil {
		common.Suggar.Error("log err")
		data.ResponseJson(w, common.USER_PARSEFORM, http.StatusInternalServerError)
		return
	}

	common.Suggar.Debugf("login user id = %d", u.ID)

	data.ResponseJson(w, common.SUCCESS, http.StatusOK)
	return
}
