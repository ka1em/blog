package controllers

import (
	"crypto/sha1"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"

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

	p := &registParams{}

	if err := p.get(r.PostForm); err != nil {
		data.ResponseJson(w, common.USER_PARAMAGET, http.StatusBadRequest)
		common.Suggar.Error(err.Error())
		return
	}

	if err := p.valid(); err != nil {
		data.ResponseJson(w, common.USER_PARAMVALID, http.StatusBadRequest)
		common.Suggar.Error(err.Error())
		return
	}

	common.Suggar.Debugf("%+v", p)

	gure := regexp.MustCompile("[^A-Za-z0-9]+")
	guid := gure.ReplaceAllString(p.UserName, "")
	passwd := weakPasswordHash(p.UserPasswd)

	common.Suggar.Debugf("%s", p.UserName)

	u := &model.User{}
	u.UserName = p.UserName
	u.UserEmail = p.UserEmail
	u.UserPasswd = string(passwd)
	u.UserGuid = guid

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

type registParams struct {
	UserName   string `schema:"user_name"`
	UserEmail  string `schema:"user_email"`
	UserPasswd string `schema:"user_passwd"`
}

func (rp *registParams) get(u url.Values) error {
	err := model.SchemaDecoder.Decode(rp, u)
	if err != nil {
		common.Suggar.Error(err.Error())
		return err
	}

	common.Suggar.Debug(rp.UserName)
	return nil
}
func (rp *registParams) valid() error {
	if rp.UserName == "" || rp.UserEmail == "" || rp.UserPasswd == "" {
		return errors.New("name, email or passwd is null")
	}
	return nil
}

func weakPasswordHash(p string) []byte {
	hash := sha1.New()
	io.WriteString(hash, p)
	return hash.Sum(nil)
}

//
func LoginPost(w http.ResponseWriter, r *http.Request) {
	model.ValidateSeesion(w, r)
	data := model.GetBaseData()
	if err := r.ParseForm(); err != nil {
		data.ResponseJson(w, common.USER_PARSEFORM, http.StatusBadRequest)
		common.Suggar.Error(err.Error())
		return
	}

	name := r.PostFormValue("user_name")
	passwd := r.PostFormValue("user_passwd")

	passwd = string(weakPasswordHash(passwd))

	u := &model.User{}
	if ok := u.Login(name, passwd); !ok {
		common.Suggar.Error("log err")
		data.ResponseJson(w, common.USER_PARSEFORM, http.StatusInternalServerError)
		return
	}

	model.UpdateSession(string(model.UserSession.UserId), u.ID)

	common.Suggar.Debugf("login user id = %d", u.ID)

	http.Redirect(w, r, "/page/hello", 301)
	return
}

//
//type loginParams struct {
//	name   string `schema:"user_name"`
//	passwd string `schema:"user_passwd"`
//}
//
//func (l *loginParams) get(v url.Values) error {
//	err := schema.NewDecoder().Decode(l, v)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (l *loginParams) valid() error {
//	if l.name == "" || l.passwd == "" {
//		return errors.New("name or passwd is nil")
//	}
//	return nil
//}
