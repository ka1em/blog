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
	d.Msg = "SUCCESS"
	d.Data = map[string]interface{}{}

	return d
}

func (d *Data) Marshal() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Data) ResponseJson(w http.ResponseWriter, datacode int, datamsg string, httpStateCode int) {
	d.Code = datacode
	//d.Msg =  common.ERRMAP[datacode]
	d.Msg = datamsg
	da, err := d.Marshal()
	if err != nil {
		common.Suggar.Error(err.Error())
		panic(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStateCode)
	w.Write(da)
	return
}
