package controllers

import (
	"blog/common/zlog"
	"blog/model"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// APICommentPOST 创建评论
func APICommentPOST(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()

	uid, err := model.ValidSessionUID(r)
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.NoUserID, http.StatusUnauthorized)
		return
	}

	if err = r.ParseForm(); err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.ParamsErr, http.StatusBadRequest)
		return
	}

	param := new(commentPostParam)
	if err := model.SchemaDecoder.Decode(param, r.PostForm); err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.ParamsErr, http.StatusBadRequest)
		return
	}
	cid, err := model.SF.NextID()
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.DataBaseErr, http.StatusInternalServerError)
		return
	}
	cm := &model.Comment{
		ID:     cid,
		Text:   param.Comment,
		PageID: param.PageID,
		UserID: uid,
	}

	if err := cm.Add(); err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.DataBaseErr, http.StatusInternalServerError)
		return
	}

	data.ResponseJson(w, model.Success, http.StatusOK)
}

type commentPostParam struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Comment  string `schema:"comment,required"`
	PageID   uint64 `schema:"page_id,required"`
	ToUserID uint64 `schema:"to_user_id"`
}

// APICommentGET 获取评论
func APICommentGET(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()
	// TODO

	data.ResponseJson(w, model.Success, http.StatusOK)
}

// APICommentPUT 更新评论
func APICommentPUT(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()

	uid, err := model.ValidSessionUID(r)
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.NoUserID, http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.ParamsErr, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	idn, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.ParamsErr, http.StatusBadRequest)
		return
	}

	param := new(commentPutParam)
	if err := model.SchemaDecoder.Decode(param, r.PostForm); err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.ParamsErr, http.StatusBadRequest)
		return
	}

	if param.UserID != idn {
		zlog.ZapLog.Error(errors.New("not self comment"))
		data.ResponseJson(w, model.ParamsErr, http.StatusBadRequest)
		return
	}

	c := &model.Comment{
		ID:     idn,
		Text:   param.Comment,
		UserID: uid,
	}

	if err := c.Update(); err != nil {
		zlog.ZapLog.Error(fmt.Errorf("update comment err : %s", err.Error()))
		data.ResponseJson(w, model.DataBaseErr, http.StatusInternalServerError)
		return
	}

	data.ResponseJson(w, model.Success, http.StatusOK)
}

type commentPutParam struct {
	UserID  uint64 `schema:"user_id,required"`
	Comment string `schema:"comment,required"`
}
