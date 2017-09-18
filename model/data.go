package model

import (
	"encoding/json"
	"net/http"

	"blog/common"
)

const (
	SUCCESS = 0

	USERNAMEEXIST = (iota * -1) - 10000 // -10000
	PARAMSERR
	PASSWDERROR
	NOUSERNAME

	DATABASEERR = (iota * -1) - 20000 // -20000

	MIDDLEWAREERR = (iota * -1) - 30000 // -30000

	NEEDLOGIN
)

var errMap = map[int]string{
	SUCCESS: "success",

	USERNAMEEXIST: "user name was exist",
	PARAMSERR:     "params error",
	DATABASEERR:   "database create user error",
	MIDDLEWAREERR: "middler ware error",
	NEEDLOGIN:     "not login",
	PASSWDERROR:   "password error",
	NOUSERNAME:    "no user name",
}

type Data struct {
	Code int                    `json:"code,string"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

func GetBaseData() *Data {
	d := &Data{}

	d.Code = 0
	d.Msg = "success"
	d.Data = map[string]interface{}{}

	return d
}

func (d *Data) ResponseJson(w http.ResponseWriter, datacode, httpStateCode int) {
	d.Code = datacode
	d.Msg = errMap[datacode]

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStateCode)
	err := json.NewEncoder(w).Encode(d)
	if err != nil {
		common.Suggar.Error(err.Error())
		panic(err)
	}

	return
}
