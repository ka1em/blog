package controllers

import (
	"blog/common/zlog"
	"blog/model"
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"sort"
)

const Token = "WXKXTOKEN001"

func WeChatValidGET(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		zlog.ZapLog.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	param := new(weChatValidParam)
	if err := model.SchemaDecoder.Decode(param, r.Form); err != nil {
		zlog.ZapLog.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !validWxSign(param.Timestamp, param.Nonce, param.Signature) {
		zlog.ZapLog.Error("wx sign error")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(param.EchoStr))
}

type weChatValidParam struct {
	Signature string `form:"signature"`
	Timestamp string `form:"timestamp"`
	Nonce     string `form:"nonce"`
	EchoStr   string `form:"echostr"`
}

func wxSign(timeStamp, nonce string) string {
	var sl sort.StringSlice
	sl = append(sl, Token)
	sl = append(sl, timeStamp)
	sl = append(sl, nonce)

	sl.Sort()

	zlog.ZapLog.Debugf("wxSign sort slice: %+v", sl)

	var str string
	for _, v := range sl {
		str += v
	}

	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func validWxSign(timeStamp, nonce, signature string) bool {
	return signature == wxSign(timeStamp, nonce)
}
