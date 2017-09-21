package controllers

import (
	"blog/common/zlog"
	"blog/model"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

//注册
func RegisterPost(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()
	if err := r.ParseForm(); err != nil {
		data.ResponseJson(w, model.PARAMS_ERR, http.StatusBadRequest)
		zlog.ZapLog.Error(err.Error())
		return
	}

	param := new(userRegistParam)
	if err := model.SchemaDecoder().Decode(param, r.PostForm); err != nil {
		zlog.ZapLog.Errorf("%s", err)
		data.ResponseJson(w, model.PARAMS_ERR, http.StatusBadRequest)
		return
	}

	if err := param.valid(); err != nil {
		zlog.ZapLog.Errorf("%s", err)
		data.ResponseJson(w, model.PARAMS_ERR, http.StatusBadRequest)
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
		zlog.ZapLog.Error(err.Error())
		if err.Error() == "exists" {
			data.ResponseJson(w, model.USER_NAME_EXIST, http.StatusBadRequest)
			return
		}
		data.ResponseJson(w, model.DATABASE_ERR, http.StatusInternalServerError)
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

func (p *userRegistParam) valid() error {
	if p.Name == "" {
		return errors.New("regist name is nil")
	}
	if p.Email == "" {
		return errors.New("regist mail is nil")
	}
	if p.Passwd == "" {
		return errors.New("regist passwd is nil")
	}
	return nil
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
	//创建session_id
	sid, err := model.PreCreateSession(w, r)
	if err != nil {
		zlog.ZapLog.Errorf("%s", err.Error())
		data.ResponseJson(w, model.DATABASE_ERR, http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		data.ResponseJson(w, model.PARAMS_ERR, http.StatusBadRequest)
		zlog.ZapLog.Error("user zlog in ", err.Error())
		return
	}

	param := new(loginParams)
	if err := model.SchemaDecoder().Decode(param, r.PostForm); err != nil {
		zlog.ZapLog.Errorf("%s", err)
		data.ResponseJson(w, model.PARAMS_ERR, http.StatusBadRequest)
		return
	}

	if err := param.valid(); err != nil {
		zlog.ZapLog.Errorf("%+v", err)
		data.ResponseJson(w, model.PARAMS_ERR, http.StatusBadRequest)
		return
	}

	u := &model.User{}
	var ok bool

	//用户密码salt
	if u, ok = model.GetValidInfo(param.Name); !ok {
		zlog.ZapLog.Error("No user")
		data.ResponseJson(w, model.NO_USER_NAME, http.StatusBadRequest)
		return
	}

	if u.UserPasswd != passwordHash(param.Passwd, u.UserSalt) {
		zlog.ZapLog.Error("password wrong")
		data.ResponseJson(w, model.PASSWD_ERR, http.StatusBadRequest)
		return
	}

	//登录成功，更新session，关联userid和sessionid
	if err := model.UpdateSession(u.ID, sid); err != nil {
		zlog.ZapLog.Error("zlog err %s", err.Error())
		data.ResponseJson(w, model.DATABASE_ERR, http.StatusInternalServerError)
		return
	}

	data.Data["redirct_url"] = "/index"
	data.ResponseJson(w, model.SUCCESS, http.StatusOK)
	return
}

type loginParams struct {
	Name   string `schema:"name"`
	Passwd string `schema:"passwd"`
}

func (p *loginParams) valid() error {
	if p.Name == "" {
		return errors.New("login param name is nil ")
	}
	if p.Passwd == "" {
		return errors.New("login param passwd is nil ")
	}
	return nil
}

//登出
func LogoutGET(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()

	var uid uint64
	var userIds interface{}

	if userIds = r.Context().Value("user_id"); userIds == nil {
		zlog.ZapLog.Error("need login ")
		data.ResponseJson(w, model.NEED_LOGIN, http.StatusUnauthorized)
		return
	}
	zlog.ZapLog.Debugf("user_id in logout get : %s", userIds)

	uid, _ = strconv.ParseUint(userIds.(string), 10, 64)

	s := &model.Session{
		UserId: uid,
	}

	if err := s.Close(); err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.PARAMS_ERR, http.StatusInternalServerError)
		return
	}

	data.ResponseJson(w, model.SUCCESS, http.StatusOK)
	return
}
