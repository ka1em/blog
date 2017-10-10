package controllers

import (
	"blog/common/zlog"
	"blog/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// APICommentPOST 创建评论
func APICommentPOST(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()

	uid, err := model.ValidSessionUID(r)
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.NO_USER_ID, http.StatusUnauthorized)
		return
	}

	zlog.ZapLog.Debug("api comment post user_id = %d", uid)

	if err = r.ParseForm(); err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.PARAMS_ERR, http.StatusBadRequest)
		return
	}

	param := new(commentPostParam)
	if err := model.SchemaDecoder.Decode(param, r.PostForm); err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.PARAMS_ERR, http.StatusBadRequest)
		return
	}

	cm := &model.Comment{
		//CommentName:  param.Name,
		//CommentEmail: param.Email,
		Text:   param.Comment,
		PageId: param.PageID,
		UserId: uid,
	}

	if err := cm.Add(); err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.DATABASE_ERR, http.StatusInternalServerError)
		return
	}

	data.ResponseJson(w, model.SUCCESS, http.StatusOK)
	return
}

type commentPostParam struct {
	Name    string `schema:"name"`
	Email   string `schema:"email"`
	Comment string `schema:"comment,required"`
	PageID  int64  `schema:"page_id,required"`
}

// APICommentGET 获取评论
func APICommentGET(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()
	// TODO

	data.ResponseJson(w, model.SUCCESS, http.StatusOK)
	return
}

func APICommentPUT(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()

	uid, err := model.ValidSessionUID(r)
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.NO_USER_ID, http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.PARAMS_ERR, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	idn, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.PARAMS_ERR, http.StatusBadRequest)
		return
	}

	param := new(commentPutParam)
	if err := model.SchemaDecoder.Decode(param, r.PostForm); err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.PARAMS_ERR, http.StatusBadRequest)
		return
	}

	if param.UserID != idn {
		zlog.ZapLog.Error(errors.New("not self comment"))
		data.ResponseJson(w, model.PARAMS_ERR, http.StatusBadRequest)
		return
	}

	c := &model.Comment{
		ID:     idn,
		Text:   param.Comment,
		UserId: uid,
	}

	if err := c.Update(); err != nil {
		zlog.ZapLog.Error(errors.New(fmt.Sprintf("update comment err : %s", err.Error())))
		data.ResponseJson(w, model.DATABASE_ERR, http.StatusInternalServerError)
		return
	}

	data.ResponseJson(w, model.SUCCESS, http.StatusOK)
	return
}

type commentPutParam struct {
	UserID  int64  `schema:"user_id,required"`
	Comment string `schema:"comment,required"`
}
