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
	"github.com/gorilla/schema"
)

func RegisterPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		common.Suggar.Error(err.Error())
		return
	}

	p := &registParams{}

	if err := p.get(r.PostForm); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		common.Suggar.Error(err.Error())
		return
	}

	if err := p.valid(); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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

	data := model.GetBaseData()

	//创建用户
	if err := u.CreateUser(); err != nil {
		if err.Error() == "exists" {
			common.Suggar.Error(err.Error())
			data.ResponseJson(w, common.ERR_USER_EXIST, err.Error(), http.StatusBadRequest)
			return
		}

		common.Suggar.Error(err.Error())
		data.ResponseJson(w, -2, err.Error(), http.StatusInternalServerError)
		return
	}

	data.ResponseJson(w, 0, "SUCCESS", http.StatusOK)
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
	ValidateSession(w, r)
	params := &loginParams{}
	if err := params.get(r.PostForm); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		common.Suggar.Error("%s", err.Error())
		return
	}

	if err := params.valid(); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		common.Suggar.Error(err.Error())
		return
	}

	params.passwd = string(weakPasswordHash(params.passwd))

	u := &model.User{}
	if ok := u.Login(params.name, params.passwd); !ok {
		common.Suggar.Error("log err")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	model.UpdateSession(model.UserSession.SessionId, u.ID)

	common.Suggar.Debug(u.ID)
	http.Redirect(w, r, "/page/hello", 301)
	return
}

type loginParams struct {
	name   string `schema:"user_name"`
	passwd string `schema:"user_passwd"`
}

func (l *loginParams) get(v url.Values) error {
	err := schema.NewDecoder().Decode(l, v)
	if err != nil {
		return err
	}
	return nil
}

func (l *loginParams) valid() error {
	if l.name == "" || l.passwd == "" {
		return errors.New("name or passwd is nil")
	}
	return nil
}
