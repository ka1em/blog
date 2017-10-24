package controllers

import (
	"blog/common/zlog"
	"blog/model"
	"net/http"

	"github.com/pkg/errors"
)

// RegisterPost 注册
func RegisterPost(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()
	if err := r.ParseForm(); err != nil {
		data.ResponseJson(w, model.ParamsErr, http.StatusBadRequest)
		zlog.ZapLog.Error(err.Error())
		return
	}

	param := new(userRegistParam)
	if err := model.SchemaDecoder.Decode(param, r.PostForm); err != nil {
		zlog.ZapLog.Errorf("%s", err)
		data.ResponseJson(w, model.ParamsErr, http.StatusBadRequest)
		return
	}

	if err := param.valid(); err != nil {
		zlog.ZapLog.Errorf("%s", err)
		data.ResponseJson(w, model.ParamsErr, http.StatusBadRequest)
		return
	}

	u := &model.User{
		Name:   param.Name,
		Email:  param.Email,
		Passwd: param.Passwd,
	}

	//创建用户
	if err := u.CreateUser(); err != nil {
		if err.Error() == model.ErrMap[model.UserNameExist] {
			zlog.ZapLog.Error(err.Error())
			data.ResponseJson(w, model.UserNameExist, http.StatusBadRequest)
			return
		}
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.DataBaseErr, http.StatusInternalServerError)
		return
	}

	data.ResponseJson(w, model.SUCCESS, http.StatusOK)
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

// LoginPost 登录
func LoginPost(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()
	//创建session_id
	sid, err := model.PreCreateSession(w, r)
	if err != nil {
		zlog.ZapLog.Errorf("%s", err.Error())
		data.ResponseJson(w, model.DataBaseErr, http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		data.ResponseJson(w, model.ParamsErr, http.StatusBadRequest)
		zlog.ZapLog.Error("user zlog in ", err.Error())
		return
	}

	param := new(loginParams)
	if err := model.SchemaDecoder.Decode(param, r.PostForm); err != nil {
		zlog.ZapLog.Errorf("%s", err)
		data.ResponseJson(w, model.ParamsErr, http.StatusBadRequest)
		return
	}

	if err := param.valid(); err != nil {
		zlog.ZapLog.Errorf("%+v", err)
		data.ResponseJson(w, model.ParamsErr, http.StatusBadRequest)
		return
	}

	u, ok, err := model.CheckPassWord(param.Name, param.Passwd)
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.NoUserName, http.StatusBadRequest)
		return
	}

	if !ok {
		zlog.ZapLog.Error("passwd error")
		data.ResponseJson(w, model.PasswordErr, http.StatusBadRequest)
		return
	}

	//登录成功，更新session，关联userid和sessionid
	if err := model.UpdateSession(u.ID, sid); err != nil {
		zlog.ZapLog.Error("zlog err %s", err.Error())
		data.ResponseJson(w, model.DataBaseErr, http.StatusInternalServerError)
		return
	}

	data.Data["redirct_url"] = "/index"
	data.ResponseJson(w, model.SUCCESS, http.StatusOK)
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

// LogoutGET 登出
func LogoutGET(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()

	uid, err := model.ValidSessionUID(r)
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.NoUserID, http.StatusBadRequest)
		return
	}

	s := &model.Session{}

	if err := s.Close(uid); err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.ParamsErr, http.StatusInternalServerError)
		return
	}

	data.ResponseJson(w, model.SUCCESS, http.StatusOK)
}
