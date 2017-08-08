package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"blog.ka1em.site/common"
	"blog.ka1em.site/model"
	"github.com/gorilla/mux"
)

func RedirIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", 301)
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()
	if err := r.ParseForm(); err != nil {

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
		data.ResponseJson(w, common.PAGE_GUIDNOTFOUND, http.StatusNotFound)
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
	pageGuid := vars["guid"]
	common.Suggar.Debugf("page guid : %s", pageGuid)

	p := &model.Page{}
	if err := p.GetByPageGUID(pageGuid); err != nil {
		common.Suggar.Error(err.Error())
		data.ResponseJson(w, common.PAGE_GUIDNOTFOUND, http.StatusNotFound)
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
