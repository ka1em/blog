package wechat

import (
	"blog/common/setting"
	"blog/common/zlog"
	"fmt"
	"net/http"
	"net/url"

	"github.com/parnurzeal/gorequest"
)

// WXUserInfo 微信用户信息
type WXUserInfo struct {
	ErrCode    int      `json:"errcode"`    //
	ErrMsg     string   `json:"errmsg"`     //
	OpenID     string   `json:"openid"`     //普通用户的标识，对当前开发者帐号唯一
	Nickname   string   `json:"nickname"`   //普通用户昵称
	Sex        int      `json:"sex"`        //普通用户性别，1为男性，2为女性
	Province   string   `json:"province"`   //普通用户个人资料填写的省份
	City       string   `json:"city"`       //普通用户个人资料填写的城市
	Country    string   `json:"country"`    //国家，如中国为CN
	HeadImgURL string   `json:"headimgurl"` //用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空
	Privilege  []string `json:"privilege"`  //用户特权信息，json数组，如微信沃卡用户为（chinaunicom）
	UnionID    string   `json:"unionid"`    //用户统一标识。针对一个微信开放平台帐号下的应用，同一用户的unionid是唯一的。
}

const WXAccessTokenURL = "https://api.weixin.qq.com/sns/oauth2/access_token?"
const WXUserInfoURL = "https://api.weixin.qq.com/sns/userinfo?"

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
}

// GetWeChatAccessToken 获取access_token
func GetWeChatAccessToken(code string) (AccessToken, error) {
	at := AccessToken{}
	v := url.Values{}

	v.Set("code", code)
	v.Set("appid", setting.WxAppID)
	v.Set("secret", setting.WxAppSecret)
	v.Set("grant_type", "authorization_code")

	resp, body, errs := gorequest.New().Get(WXAccessTokenURL + v.Encode()).EndStruct(&at)
	if len(errs) > 0 {
		return at, fmt.Errorf("GetWeChatAccessToken errs:%+v, body : %+v", errs, string(body))
	}

	if resp.StatusCode != http.StatusOK {
		return at, fmt.Errorf("GetWeChatAccessToken body : %s", string(body))
	}

	zlog.ZapLog.Debugf("wxat: %+v", at)
	return at, nil
}

// GetWeChatUnionID 获取微信unionid
func GetWeChatUnionID(code string) (string, error) {
	at, err := GetWeChatAccessToken(code)
	return at.UnionID, err
}

// GetWeChatUserInfo 获取微信用户信息
func GetWeChatUserInfo(code string) (WXUserInfo, error) {
	w := WXUserInfo{}

	at, err := GetWeChatAccessToken(code)
	if err != nil {
		return w, err
	}

	v := url.Values{}

	v.Set("access_token", at.AccessToken)
	v.Set("openid", at.OpenID)

	resp, body, errs := gorequest.New().Get(WXUserInfoURL + v.Encode()).EndStruct(&w)
	if len(errs) > 0 {
		return w, fmt.Errorf("GetWeChatUserInfo errs %+v", errs)
	}
	if resp.StatusCode != http.StatusOK {
		return w, fmt.Errorf("GetWeChatUserInfo body : %s", string(body))
	}
	return w, nil
}
