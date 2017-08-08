package controllers

import (
	"net/http"

	"strconv"

	"blog.ka1em.site/common"
	"blog.ka1em.site/model"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func APICommentPOST(w http.ResponseWriter, r *http.Request) {
	var err error
	data := model.GetBaseData()

	userId := context.Get(r, "user_id")
	if userId == nil {
		data.ResponseJson(w, common.NEED_LOGIN, http.StatusUnauthorized)
		return
	}

	if err = r.ParseForm(); err != nil {
		common.Suggar.Error(err.Error())
		data.ResponseJson(w, common.USER_PARAMVALID, http.StatusBadRequest)
		return
	}

	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	comment := r.PostFormValue("comment")
	pageId := r.PostFormValue("page_id")

	if name == "" || email == "" || comment == "" || pageId == "" {
		common.Suggar.Error(err.Error())
		data.ResponseJson(w, common.USER_PARAMVALID, http.StatusBadRequest)
		return
	}

	pageIdn, err := strconv.ParseUint(pageId, 10, 64)
	if err != nil {
		common.Suggar.Error(err.Error())
		data.ResponseJson(w, common.USER_PARAMVALID, http.StatusBadRequest)
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
		data.ResponseJson(w, common.USER_PARAMVALID, http.StatusInternalServerError)
		return
	}

	data.ResponseJson(w, common.SUCCESS, http.StatusOK)
	return
}

func APICommentGET(w http.ResponseWriter, r *http.Request) {

}

func APICommentPUT(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()

	if err := r.ParseForm(); err != nil {
		common.Suggar.Error(err.Error())
		data.ResponseJson(w, common.USER_PARAMVALID, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	idn, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		common.Suggar.Error(err.Error())
		data.ResponseJson(w, common.USER_PARAMVALID, http.StatusBadRequest)
		return
	}

	comment := r.FormValue("comment")

	c := &model.Comment{
		Id:          idn,
		CommentText: comment,
	}

	if err := c.UpdateComment(); err != nil {

	}
}
