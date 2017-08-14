package controllers

import (
	"net/http"
	"strconv"

	"errors"

	"blog.ka1em.site/common"
	"blog.ka1em.site/model"
	"github.com/gorilla/mux"
)

func APICommentPOST(w http.ResponseWriter, r *http.Request) {
	var err error
	data := model.GetBaseData()

	var uid uint64
	var userIds interface{}

	if userIds = r.Context().Value("user_id"); userIds == nil {
		common.Suggar.Error("need login ")
		data.ResponseJson(w, common.NEEDLOGIN, http.StatusUnauthorized)
		return
	}

	uid, _ = strconv.ParseUint(userIds.(string), 10, 64)

	common.Suggar.Debug("api comment post user_id = %d", uid)

	if err = r.ParseForm(); err != nil {
		common.Suggar.Error(err.Error())
		data.ResponseJson(w, common.PARAMSERR, http.StatusBadRequest)
		return
	}

	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	comment := r.PostFormValue("comment")
	pageId := r.PostFormValue("page_id")

	if name == "" || email == "" || comment == "" || pageId == "" {
		common.Suggar.Error(err.Error())
		data.ResponseJson(w, common.PARAMSERR, http.StatusBadRequest)
		return
	}

	pageIdn, err := strconv.ParseUint(pageId, 10, 64)
	if err != nil {
		common.Suggar.Error(err.Error())
		data.ResponseJson(w, common.PARAMSERR, http.StatusBadRequest)
		return
	}

	cm := &model.Comment{
		CommentName:  name,
		CommentEmail: email,
		CommentText:  comment,
		PageId:       pageIdn,
	}

	if err := cm.AddComment(); err != nil {
		common.Suggar.Error(err.Error())
		data.ResponseJson(w, common.DATABASEERR, http.StatusInternalServerError)
		return
	}

	data.ResponseJson(w, common.SUCCESS, http.StatusOK)
	return
}

func APICommentGET(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()

	//var uid uint64
	//if userId := context.Get(r, "user_id"); userId == nil {
	//	data.ResponseJson(w, common.NEED_LOGIN, http.StatusUnauthorized)
	//	return
	//} else {
	//	uid = userId.(uint64)
	//}
	//
	//common.Suggar.Debug("api comment post user_id = %d", uid)
	data.ResponseJson(w, common.SUCCESS, http.StatusOK)
	return
}

func APICommentPUT(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()

	var uid uint64
	var userIds interface{}

	if userIds = r.Context().Value("user_id"); userIds == nil {
		common.Suggar.Error("need login ")
		data.ResponseJson(w, common.NEEDLOGIN, http.StatusUnauthorized)
		return
	}

	uid, _ = strconv.ParseUint(userIds.(string), 10, 64)

	if err := r.ParseForm(); err != nil {
		common.Suggar.Error(err.Error())
		data.ResponseJson(w, common.PARAMSERR, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	idn, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		common.Suggar.Error(err.Error())
		data.ResponseJson(w, common.PARAMSERR, http.StatusBadRequest)
		return
	}

	userId := r.PostFormValue("user_id")
	userIdn, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {

	}

	if userIdn != uid {
		common.Suggar.Error(errors.New("not self comment"))
		data.ResponseJson(w, common.PARAMSERR, http.StatusBadRequest)
		return
	}

	comment := r.PostFormValue("comment")

	c := &model.Comment{
		Id:          idn,
		CommentText: comment,
	}

	if err := c.UpdateComment(); err != nil {

	}
}
