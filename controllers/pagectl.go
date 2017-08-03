package controllers

import (
	"net/http"

	"blog.ka1em.site/common"
	"blog.ka1em.site/model"
	"github.com/gorilla/mux"
)

func ServePage(w http.ResponseWriter, r *http.Request) {
	data := model.GetBaseData()
	vars := mux.Vars(r)

	pageGuid := vars["guid"]

	if pageGuid == "" {
		common.Suggar.Error("%s", "page guid is nil")
		data.ResponseJson(w, common.PAGE_NOPAGEGUID, http.StatusBadRequest)
		return
	}

	page := &model.Page{}

	if err := page.GetByPageGUID(pageGuid); err != nil {
		common.Suggar.Error(err.Error())
		data.ResponseJson(w, common.PAGE_GUIDNOTFOUND, http.StatusNotFound)
		return
	}

	data.Data["page"] = *page

	data.ResponseJson(w, common.SUCCESS, http.StatusOK)
	return
}
