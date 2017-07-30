package controllers

import (
	"crypto/sha1"
	"errors"
	"io"
	"net/http"
	"net/url"

	"blog.ka1em.site/common"
	"blog.ka1em.site/model"
	"github.com/gorilla/schema"
)

func RegisterPost(w http.ResponseWriter, r *http.Request) {

	//if err := r.ParseForm(); err != nil {
	//	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	//	common.Suggar.Error(err.Error())
	//	return
	//}

	p := &RegistParams{}

	p.UserName = r.FormValue("user_name")
	p.UserEmail = r.FormValue("user_email")
	p.UserPasswd = r.FormValue("user_passwd")

	//gure := regexp.MustCompile("[^A-Za-z0-9]+")
	//guid := gure.ReplaceAllString(params.name, "")
	passwd := weakPasswordHash(p.UserPasswd)

	common.Suggar.Debug("%s", p.UserName)
	u := &model.User{}
	u.UserName = p.UserName
	u.UserEmail = p.UserEmail
	u.UserPasswd = string(passwd)
	//u.UserGuid = guid
	common.Suggar.Debug("%+v", u)

	if err := u.CreateUser(); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		common.Suggar.Error(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

type RegistParams struct {
	UserName   string //`schema:"user_name"`
	UserEmail  string //`schema:"user_email"`
	UserPasswd string //`schema:"user_passwd"`
	//pageGUID string //`schema:""`
}

//func (rp *ResgistParams) get(u url.Values) error {
//	or := &ResgistParams{}
//	err := model.SchemaDecoder.Decode(or, u)
//	if err != nil {
//		common.Suggar.Error(err.Error())
//		return err
//	}
//	rp.UserName = or.UserName
//	rp.UserPasswd = or.UserPasswd
//	rp.UserEmail = or.UserEmail
//
//	common.Suggar.Debug(rp.UserName)
//	return nil
//}
//func (rp *ResgistParams) valid() error {
//	if rp.UserName == "" || rp.UserEmail == "" || rp.UserPasswd == "" {
//		return errors.New("name, email or passwd is null")
//	}
//	return nil
//}

func weakPasswordHash(p string) []byte {
	hash := sha1.New()
	io.WriteString(hash, p)
	return hash.Sum(nil)
}

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
