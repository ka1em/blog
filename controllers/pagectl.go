package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"blog/common/zlog"
	"blog/model"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// PageIndexGET
func PageIndexGET(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()
	if err := r.ParseForm(); err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.ParamsErr, http.StatusBadRequest)
		return
	}

	param := new(pageIndexParam)

	if err := model.SchemaDecoder.Decode(param, r.Form); err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.ParamsErr, http.StatusBadRequest)
		return
	}

	if err := param.valid(); err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.ParamsErr, http.StatusBadRequest)
		return
	}

	pages, err := model.GetAllPage(param.PageIndex, param.PageSize)
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.ParamsErr, http.StatusBadRequest)
		return
	}

	//截取正文 前150字符
	for _, v := range pages {
		v.Content = model.TruncatedText(v.Content)
	}

	if len(pages) < param.PageSize {
		data.Data["page_index"] = fmt.Sprintf("%d", -1)
	} else {
		data.Data["page_index"] = fmt.Sprintf("%d", param.PageIndex+1)
	}

	data.Data["page_size"] = fmt.Sprintf("%d", param.PageSize)
	data.Data["pages_list"] = pages

	data.ResponseJson(w, model.SUCCESS, http.StatusOK)
}

type pageIndexParam struct {
	PageIndex int `schema:"page_index"`
	PageSize  int `schema:"page_size"`
}

func (p *pageIndexParam) valid() error {
	var err error
	if p.PageIndex < 0 {
		return errors.New("page_index < 0")
	}
	if p.PageSize < 0 {
		return errors.New("page_size < 0")
	}
	if p.PageIndex == 0 {
		p.PageIndex = 1
	}
	if p.PageSize == 0 {
		p.PageSize = model.DefaultPageSize
	}
	return err
}

func APIPageGET(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()

	vars := mux.Vars(r)

	pageID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.ParamsErr, http.StatusBadRequest)
		return
	}

	page, err := model.GetByID(pageID)
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.ParamsErr, http.StatusBadRequest)
		return
	}

	data.Data["page"] = page
	data.ResponseJson(w, model.SUCCESS, http.StatusOK)
}

func APIPagePOST(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()

	if err := r.ParseForm(); err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.ParamsErr, http.StatusBadRequest)
		return
	}

	param := new(pagePostParam)
	if err := model.SchemaDecoder.Decode(param, r.PostForm); err != nil {
		zlog.ZapLog.Error(err.Error())
		data.ResponseJson(w, model.ParamsErr, http.StatusBadRequest)
		return
	}

	p := &model.Page{
		Content: param.Content,
		Title:   param.Title,
	}

	if err := p.Add(); err != nil {
		zlog.ZapLog.Error("%s", err.Error())
		data.ResponseJson(w, model.ParamsErr, http.StatusInternalServerError)
		return
	}

	data.ResponseJson(w, model.SUCCESS, http.StatusOK)
}

type pagePostParam struct {
	Content string `schema:"content,required"`
	Title   string `schema:"title,required"`
}
