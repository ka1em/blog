package controllers

import (
	"blog/common/zlog"
	"blog/model"
	"net/http"
	"time"

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

	// 生成uid
	uid, err := model.SF.NextID()
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.DataBaseErr, http.StatusInternalServerError)
		return
	}
	u := &model.User{
		ID:     uid,
		Name:   param.Name,
		Email:  param.Email,
		Passwd: param.Password,
	}

	//创建用户
	if err := u.Create(); err != nil {
		if err.Error() == model.ErrMap[model.UserNameExist] {
			zlog.ZapLog.Error(err.Error())
			data.ResponseJson(w, model.UserNameExist, http.StatusBadRequest)
			return
		}
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.DataBaseErr, http.StatusInternalServerError)
		return
	}

	data.ResponseJson(w, model.Success, http.StatusOK)
}

type userRegistParam struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (p *userRegistParam) valid() error {
	if p.Name == "" {
		return errors.New("register name is nil")
	}
	if p.Email == "" {
		return errors.New("register mail is nil")
	}
	if p.Password == "" {
		return errors.New("register passwd is nil")
	}
	return nil
}

// LoginPost 登录
func LoginPost(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()
	//创建session_id
	//sid, err := model.PreCreateSession(w, r)
	//if err != nil {
	//	zlog.ZapLog.Errorf("%s", err.Error())
	//	data.ResponseJson(w, model.DataBaseErr, http.StatusInternalServerError)
	//	return
	//}

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

	u := model.User{
		Name:   param.Name,
		Passwd: param.Password,
	}
	_, err := u.CheckPassWord()
	if err != nil {
		var errType int64
		var httpCode int
		switch err.Error() {
		case model.ErrMap[model.PasswordErr]:
			errType = model.PasswordErr
			httpCode = http.StatusBadRequest
		case model.ErrMap[model.PasswordHashErr]:
			errType = model.PasswordHashErr
			httpCode = http.StatusInternalServerError
		case model.ErrMap[model.NoUserName]:
			errType = model.NoUserName
			httpCode = http.StatusBadRequest
		default:
			errType = model.DataBaseErr
			httpCode = http.StatusBadRequest
		}
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, errType, httpCode)
		return
	}

	//登录成功，更新session，关联userid和sessionid
	s := model.Session{
		//SID:         sid, // TODO
		UserID:      u.ID,
		CreatedUnix: time.Now().Unix(),
	}
	if err := s.Save(); err != nil {
		zlog.ZapLog.Error("zlog err %s", err.Error())
		data.ResponseJson(w, model.DataBaseErr, http.StatusInternalServerError)
		return
	}

	data.ResponseJson(w, model.Success, http.StatusOK)
}

type loginParams struct {
	Name     string `schema:"name"`
	Password string `schema:"password"`
}

func (p *loginParams) valid() error {
	if p.Name == "" {
		return errors.New("login param name is nil ")
	}
	if p.Password == "" {
		return errors.New("login param passwd is nil ")
	}
	return nil
}

// LogoutGET 登出
func LogoutGET(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()
	// todo
	//
	//uid, err := model.ValidSessionUID(r)
	//if err != nil {
	//	zlog.ZapLog.Error(err.Error())
	//	data.ResponseJson(w, model.NoUserID, http.StatusBadRequest)
	//	return
	//}
	////sid, err := model.ValidSessionID(r)
	//
	////s := &model.Session{
	////	SID: sid,
	////	//UserID: uid,
	////}
	//
	//if err := s.Del(); err != nil {
	//	zlog.ZapLog.Error(err.Error())
	//	data.ResponseJson(w, model.ParamsErr, http.StatusInternalServerError)
	//	return
	//}

	data.ResponseJson(w, model.Success, http.StatusOK)
}

// LoginWeChatPOST 微信登录
func LoginWeChatPOST(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()

	data.ResponseJson(w, model.Success, http.StatusOK)
}

type loginweChatParam struct {
	Code int64
	//DeviceID string
}
