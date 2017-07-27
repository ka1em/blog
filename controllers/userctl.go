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

func RegisterPOST(w http.ResponseWriter, r *http.Request) {
	params := &resgistParams{}
	if err := params.get(r.PostForm); err != nil {
		common.Suggar.Error(err.Error())
		return
	}

	if err := params.valid(); err != nil {
		common.Suggar.Error(err.Error())
		return
	}

	gure := regexp.MustCompile("[^A-Za-z0-9]+")
	guid := gure.ReplaceAllString(params.name, "")
	passwd := weakPasswordHash(params.pass)

	u := &model.User{}
	u.UserName = params.name
	u.UserEmail = params.email
	u.UserPasswd = string(passwd)
	u.UserGuid = guid

}

type resgistParams struct {
	name  string `schema:"user_name"`
	email string `schema:"user_email"`
	pass  string `schema:"user_passwd"`
	//pageGUID string `schema:""`
}

func (rp *resgistParams) get(u url.Values) error {
	decoder := schema.NewDecoder()
	err := decoder.Decode(rp, u)
	if err != nil {
		common.Suggar.Error(err.Error())
		return err
	}
	return nil
}
func (rp *resgistParams) valid() error {
	if rp.name == "" || rp.email == "" || rp.pass == "" {
		return errors.New("name, email or passwd is null")
	}
	return nil
}

func weakPasswordHash(p string) []byte {
	hash := sha1.New()
	io.WriteString(hash, p)
	return hash.Sum(nil)
}
