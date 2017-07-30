package model

import "encoding/json"

type Data struct {
	Code int                    `json:"code,string"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

func GetBaseData() *Data {
	d := &Data{}

	d.Code = 0
	d.Msg = ""
	d.Data = map[string]interface{}{}

	return d
}

func (d *Data) Marshal() ([]byte, error) {
	return json.Marshal(d)
}
