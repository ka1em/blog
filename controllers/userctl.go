package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"

	"blog.ka1em.site/common"
	"blog.ka1em.site/model"
	"github.com/satori/go.uuid"
)

//注册
func RegisterPost(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()
	if err := r.ParseForm(); err != nil {
		data.ResponseJson(w, model.PARAMSERR, http.StatusBadRequest)
		common.Suggar.Error(err.Error())
		return
	}

	param := new(userRegistParam)
	if err := model.SchemaDecoder().Decode(param, r.PostForm); err != nil {
		common.Suggar.Errorf("%s", err)
		data.ResponseJson(w, model.PARAMSERR, http.StatusBadRequest)
		return
	}

	if err := param.valid(); err != nil {
		common.Suggar.Errorf("%+v", err)
		data.ResponseJson(w, model.PARAMSERR, http.StatusBadRequest)
		return
	}

	salt := uuid.NewV4().String()

	u := &model.User{
		UserName:   param.Name,
		UserEmail:  param.Email,
		UserSalt:   salt,
		UserPasswd: passwordHash(param.Passwd, salt),
	}

	//创建用户
	if err := u.CreateUser(); err != nil {
		common.Suggar.Error(err.Error())
		if err.Error() == "exists" {
			data.ResponseJson(w, model.USERNAMEEXIST, http.StatusBadRequest)
			return
		}
		data.ResponseJson(w, model.DATABASEERR, http.StatusInternalServerError)
		return
	}

	data.ResponseJson(w, model.SUCCESS, http.StatusOK)
	return
}

type userRegistParam struct {
	Name   string `schema:"name"`
	Email  string `schema:"email"`
	Passwd string `schema:"passwd"`
}

func (p *userRegistParam) valid() []error {
	var err []error
	if p.Name == "" {
		err = append(err, errors.New("name is nil"))
	}

	if p.Email == "" {
		err = append(err, errors.New("mail is nil"))
	}

	if p.Passwd == "" {
		err = append(err, errors.New("passwd is nil"))
	}

	return err
}

func passwordHash(p, salt string) string {
	hash := sha256.New()
	s := p + salt
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

//登录
func LoginPost(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()
	//TODO 重复登录？

	//创建session_id
	sid, err := model.CreateSession(w, r)
	if err != nil {
		common.Suggar.Error("log err %s", err.Error())
		data.ResponseJson(w, model.DATABASEERR, http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		data.ResponseJson(w, model.PARAMSERR, http.StatusBadRequest)
		common.Suggar.Error("user log in ", err.Error())
		return
	}

	param := new(loginParams)
	if err := model.SchemaDecoder().Decode(param, r.PostForm); err != nil {
		common.Suggar.Errorf("%s", err)
		data.ResponseJson(w, model.PARAMSERR, http.StatusBadRequest)
		return
	}

	if err := param.valid(); err != nil {
		common.Suggar.Errorf("%+v", err)
		data.ResponseJson(w, model.PARAMSERR, http.StatusBadRequest)
		return
	}

	u := &model.User{}
	var ok bool

	//用户密码salt
	if u, ok = model.GetValidInfo(param.Name); !ok {
		common.Suggar.Error("No user")
		data.ResponseJson(w, model.NOUSERNAME, http.StatusBadRequest)
		return
	}

	if u.UserPasswd != passwordHash(param.Passwd, u.UserSalt) {
		common.Suggar.Error("password wrong")
		data.ResponseJson(w, model.PASSWDERROR, http.StatusBadRequest)
		return
	}

	//登录成功，更新session，关联userid和sessionid
	if err := model.UpdateSession(u.ID, sid); err != nil {
		common.Suggar.Error("log err %s", err.Error())
		data.ResponseJson(w, model.DATABASEERR, http.StatusInternalServerError)
		return
	}

	data.ResponseJson(w, model.SUCCESS, http.StatusOK)
	return
}

type loginParams struct {
	Name   string `schema:"name"`
	Passwd string `schema:"passwd"`
}

func (p *loginParams) valid() []error {
	var err []error
	if p.Name == "" {
		err = append(err, errors.New("name is nil "))
	}
	if p.Passwd == "" {
		err = append(err, errors.New("passwd is nil "))
	}
	return err
}

//登出
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

	s := &model.Session{UserId: uid}

	if err := s.CloseSession(); err != nil {
		common.Suggar.Error(err.Error())
		data.ResponseJson(w, model.PARAMSERR, http.StatusInternalServerError)
		return
	}

	data.ResponseJson(w, model.SUCCESS, http.StatusOK)
	return
}
