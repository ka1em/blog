package model

import (
	"encoding/json"
	"net/http"

	"blog.ka1em.site/common"
)

type Data struct {
	Code int                    `json:"code,string"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

func GetBaseData() *Data {
	d := &Data{}

	d.Code = 0
	d.Msg = "sucess"
	d.Data = map[string]interface{}{}

	return d
}

func (d *Data) Marshal() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Data) ResponseJson(w http.ResponseWriter, datacode int, httpStateCode int) {

	var (
		err error
		da  []byte
	)

	d.Code = datacode
	d.Msg = common.ERRMAP[datacode]

	if da, err = d.Marshal(); err != nil {
		common.Suggar.Error(err.Error())
		panic(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStateCode)
	w.Write(da)
	return
}
