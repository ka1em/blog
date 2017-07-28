package controllers

import (
	"net/http"
	"time"

	"blog.ka1em.site/model"

	"blog.ka1em.site/common"
	"github.com/gorilla/mux"
)

func ServePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageGUID := vars["guid"]

	page := &model.Page{}

	err := page.GetByPageGUID(pageGUID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		common.Suggar.Error(err.Error())
		return
	}

	html := `<html>
	        	<head>
			<title>` + page.Title + `</title>
		  	</head>
		<body>
			<h1>` + page.Content + `</h1>
			<div>` + page.CreatedAt.Format(time.ANSIC) + `</div>
		</body>
		</html>`

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
	return
}
