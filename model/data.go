package model

import (
	"blog/common/zlog"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

const (
	SUCCESS         = 0
	USER_NAME_EXIST = (iota * -1) - 10000 // -10000
	PARAMS_ERR
	PASSWD_ERR
	NO_USER_NAME
	NO_USER_ID
	DATABASE_ERR   = (iota * -1) - 20000 // -20000
	MIDDLEWARE_ERR = (iota * -1) - 30000 // -30000
	NEED_LOGIN
)

var errMap = map[int]string{
	SUCCESS:         "success",
	USER_NAME_EXIST: "user name was exist",
	PARAMS_ERR:      "params error",
	//NO_USER_NAME:    "database create user error",
	MIDDLEWARE_ERR:  "middler ware error",
	NEED_LOGIN:      "not login",
	PASSWD_ERR:      "password error",
	NO_USER_NAME:    "no user name",
	NO_USER_ID:      "no user id",
}

type data struct {
	Code int                    `json:"code,string"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

// GetBaseData return the base data
func GetBaseData() *data {
	return &data{Code: 0, Msg: "success", Data: map[string]interface{}{}}
}

// ResponseJson write json data to response
func (d *data) ResponseJson(w http.ResponseWriter, datacode, httpStateCode int) {
	d.Code = datacode
	d.Msg = errMap[datacode]

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStateCode)
	err := jsoniter.NewEncoder(w).Encode(d)
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		panic(err)
	}

	return
}
