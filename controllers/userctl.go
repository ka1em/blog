package controllers

import (
	"net/http"
	"regexp"

	"crypto/sha256"
	"encoding/hex"

	"strconv"

	"blog.ka1em.site/common"
	"blog.ka1em.site/model"
)

func RegisterPost(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()
	if err := r.ParseForm(); err != nil {
		data.ResponseJson(w, model.PARAMSERR, http.StatusBadRequest)
		common.Suggar.Error(err.Error())
		return
	}

	name := r.PostFormValue("user_name")
	email := r.PostFormValue("user_email")
	passwd := r.PostFormValue("user_passwd")

	if name == "" || email == "" || passwd == "" {
		data.ResponseJson(w, model.PARAMSERR, http.StatusBadRequest)
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
			data.ResponseJson(w, model.USERNAMEEXIST, http.StatusBadRequest)
			return
		}

		common.Suggar.Error(err.Error())
		data.ResponseJson(w, model.DATABASEERR, http.StatusInternalServerError)
		return
	}

	data.ResponseJson(w, model.SUCCESS, http.StatusOK)
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
		data.ResponseJson(w, model.PARAMSERR, http.StatusBadRequest)
		common.Suggar.Error("user log in ", err.Error())
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
		common.Suggar.Error("login err")
		data.ResponseJson(w, model.DATABASEERR, http.StatusInternalServerError)
		return
	}

	s.UserId = u.ID
	if err := s.UpdateSession(); err != nil {
		common.Suggar.Error("log err %s", err.Error())
		data.ResponseJson(w, model.DATABASEERR, http.StatusInternalServerError)
		return
	}

	common.Suggar.Debugf("login user id = %d", u.ID)

	data.ResponseJson(w, model.SUCCESS, http.StatusOK)
	return
}

func LogoutGET(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()

	var uid uint64
	var userIds interface{}

	if userIds = r.Context().Value("user_id"); userIds == nil {
		common.Suggar.Error("need login ")
		data.ResponseJson(w, model.NEEDLOGIN, http.StatusUnauthorized)
		return
	}
	common.Suggar.Debugf("user_id in logout get : %s", userIds)

	uid, _ = strconv.ParseUint(userIds.(string), 10, 64)

	s := &model.Session{}
	s.UserId = uid

	if err := s.CloseSession(); err != nil {
		common.Suggar.Error(err.Error())
		data.ResponseJson(w, model.PARAMSERR, http.StatusInternalServerError)
		return
	}

	data.ResponseJson(w, model.SUCCESS, http.StatusOK)
	return
}
