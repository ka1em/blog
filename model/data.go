package model

import (
	"blog/common/zlog"
	"net/http"

	"github.com/json-iterator/go"
)

const (
	Success       = 0
	UserNameExist = (iota * -1) - 10000 // -10000
	ParamsErr
	PasswordErr
	NoUserName
	NoUserID
	SessionNoUserID
	DataBaseErr   = (iota * -1) - 20000 // -20000
	MiddlewareErr = (iota * -1) - 30000 // -30000
	NeedLogin
)

const DefaultPageSize = 20

// ErrMap 错误map
var ErrMap = map[int64]string{
	Success:         "success",
	UserNameExist:   "error: user name was exist",
	ParamsErr:       "error: params was error",
	MiddlewareErr:   "error: middleware error",
	NeedLogin:       "error: not login",
	PasswordErr:     "error: password error",
	NoUserName:      "error: no user name",
	NoUserID:        "error: no user id",
	SessionNoUserID: "error: session no user id",
}

// Data return json data
type Data struct {
	Code int64                  `json:"code,string"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

// GetBaseData return the base data
func GetBaseData() *Data {
	return &Data{
		Code: Success,
		Msg:  ErrMap[Success],
		Data: map[string]interface{}{},
	}
}

// ResponseJson write json data to response
func (d *Data) ResponseJson(w http.ResponseWriter, code int64, httpCode int) {
	d.Code = code
	d.Msg = ErrMap[code]

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	if err := jsoniter.NewEncoder(w).Encode(d); err != nil {
		zlog.ZapLog.Error(err.Error())
		panic(err)
	}
}
