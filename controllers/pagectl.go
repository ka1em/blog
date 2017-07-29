package controllers

import (
	"net/http"

	"encoding/json"

	"blog.ka1em.site/common"
	"blog.ka1em.site/model"
	"github.com/gorilla/mux"
)

func ServePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	pageGuid := vars["guid"]

	if pageGuid == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		common.Suggar.Error("%s", "page guid is nil")
		return
	}

	page := &model.Page{}

	err := page.GetByPageGUID(pageGuid)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		common.Suggar.Error(err.Error())
		return
	}

	data := model.GetBaseData()
	data.Code = 0
	data.Msg = "ok"
	data.Data["page"] = *page

	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusNotFound)
		common.Suggar.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	return
}
