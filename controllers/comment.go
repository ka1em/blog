package controllers

import (
	"blog/common/log"
	"blog/model"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func APICommentPOST(w http.ResponseWriter, r *http.Request) {
	var err error
	data := model.GetBaseData()

	var uid uint64
	var userIds interface{}

	if userIds = r.Context().Value("user_id"); userIds == nil {
		log.Suggar.Error("need login ")
		data.ResponseJson(w, model.NEEDLOGIN, http.StatusUnauthorized)
		return
	}

	uid, _ = strconv.ParseUint(userIds.(string), 10, 64)

	log.Suggar.Debug("api comment post user_id = %d", uid)

	if err = r.ParseForm(); err != nil {
		log.Suggar.Error(err.Error())
		data.ResponseJson(w, model.PARAMSERR, http.StatusBadRequest)
		return
	}

	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	comment := r.PostFormValue("comment")
	pageId := r.PostFormValue("page_id")

	if name == "" || email == "" || comment == "" || pageId == "" {
		log.Suggar.Error(err.Error())
		data.ResponseJson(w, model.PARAMSERR, http.StatusBadRequest)
		return
	}

	pageIdn, err := strconv.ParseUint(pageId, 10, 64)
	if err != nil {
		log.Suggar.Error(err.Error())
		data.ResponseJson(w, model.PARAMSERR, http.StatusBadRequest)
		return
	}

	cm := &model.Comment{CommentName: name, CommentEmail: email, CommentText: comment, PageId: pageIdn}

	if err := cm.AddComment(); err != nil {
		log.Suggar.Error(err.Error())
		data.ResponseJson(w, model.DATABASEERR, http.StatusInternalServerError)
		return
	}

	data.ResponseJson(w, model.SUCCESS, http.StatusOK)
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
	data.ResponseJson(w, model.SUCCESS, http.StatusOK)
	return
}

func APICommentPUT(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()

	userIds := r.Context().Value("user_id")
	if userIds == nil {
		log.Suggar.Error("need login ")
		data.ResponseJson(w, model.NEEDLOGIN, http.StatusUnauthorized)
		return
	}

	uid, _ := strconv.ParseUint(userIds.(string), 10, 64)

	if err := r.ParseForm(); err != nil {
		log.Suggar.Error(err.Error())
		data.ResponseJson(w, model.PARAMSERR, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	idn, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Suggar.Error(err.Error())
		data.ResponseJson(w, model.PARAMSERR, http.StatusBadRequest)
		return
	}

	userId := r.PostFormValue("user_id")
	userIdn, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {

	}

	if userIdn != uid {
		log.Suggar.Error(errors.New("not self comment"))
		data.ResponseJson(w, model.PARAMSERR, http.StatusBadRequest)
		return
	}

	comment := r.PostFormValue("comment")

	c := &model.Comment{Id: idn, CommentText: comment}

	if err := c.UpdateComment(); err != nil {
		log.Suggar.Error(errors.New(fmt.Sprintf("update comment err : %s", err.Error())))
		data.ResponseJson(w, model.DATABASEERR, http.StatusInternalServerError)
		return
	}

	data.ResponseJson(w, model.SUCCESS, http.StatusOK)
	return
}
