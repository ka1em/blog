package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"blog.ka1em.site/common"
	"blog.ka1em.site/model"
	"github.com/gorilla/mux"
	"errors"
)

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()
	if err := r.ParseForm(); err != nil {
		common.Suggar.Error(err.Error())
		data.ResponseJson(w, common.PARAMSERR, http.StatusBadRequest)
		return
	}

	var (
		pIndex int = 1
		pSize  int = 10
		err    error
	)
	if pi := r.Form["page_index"]; pi != nil {
		pIndex, err = strconv.Atoi(pi[0])
		if err != nil {

		}
	}
	common.Suggar.Debugf("page_index = %d, page_size = %d", pIndex, pSize)
	p := &model.Page{}

	pages, err := p.GetAllPage(pIndex, pSize)
	if err != nil {
		common.Suggar.Error(err.Error())
		data.ResponseJson(w, common.PARAMSERR, http.StatusBadRequest)
		return
	}

	//截取正文 前150字符
	for _, v := range pages {
		v.Content = v.TruncatedText()
	}

	data.Data["pages_list"] = pages

	data.Data["page_index"] = fmt.Sprintf("%d", pIndex+1)
	data.Data["page_size"] = fmt.Sprintf("%d", pSize)
	if len(pages) < pSize {
		data.Data["page_index"] = fmt.Sprintf("%d", -1)
	}

	data.ResponseJson(w, common.SUCCESS, http.StatusOK)
	return
}

func APIPage(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()
	vars := mux.Vars(r)
	pageId := vars["id"]
	common.Suggar.Debugf("page guid : %s", pageId)

	pageIdn, err := strconv.ParseUint(pageId, 10, 64)
	if err != nil {
		common.Suggar.Error(err.Error())
		data.ResponseJson(w, common.PARAMSERR, http.StatusBadRequest)
		return
	}

	p := &model.Page{
		Id: pageIdn,
	}

	if err := p.GetByID(); err != nil {
		common.Suggar.Error(err.Error())
		data.ResponseJson(w, common.PARAMSERR, http.StatusBadRequest)
		return
	}

	c := &model.Comment{
		PageId: p.Id,
	}

	c.GetComment(1, 10)

	data.Data["page"] = *p
	data.Data["comments"] = *c

	data.ResponseJson(w, common.SUCCESS, http.StatusOK)
	return
}

func PageAdd(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()

	if err := r.ParseForm(); err != nil {
		common.Suggar.Error(err.Error())
		data.ResponseJson(w, common.PARAMSERR, http.StatusBadRequest)
		return
	}

	p := &model.Page{}
	p.Content = r.PostFormValue("content")
	p.Title = r.PostFormValue("title")

	if p.Title == "" || p.Content == "" {
		common.Suggar.Error(errors.New("title or content is nill"))
		data.ResponseJson(w, common.PARAMSERR, http.StatusBadRequest)
		return
	}

	//todo

}
