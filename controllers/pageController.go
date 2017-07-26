package controllers

import (
	"net/http"
	"blog.ka1em.site/model"
	"log"
	//"blog.ka1em.site/common"
	//"blog.ka1em.site/data"

	"github.com/gorilla/mux"
)

func ServePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageGUID := vars["guid"]

	thisPage := model.Page{}

	resp := &model.Page{}
	//resp.DB = common.GetDB()

	err := resp.GetByPageID(pageGUID, &thisPage)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		log.Println(err.Error())
		return
	}

	html := `<html>
	        	<head>
			<title>` + thisPage.Title + `</title>
		  	</head>
		<body>
			<h1>` + thisPage.Content + `</h1>
			<div>` + thisPage.Date + `</div>
		</body>
		</html>`

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
	return
}
