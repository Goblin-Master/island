package qqLogin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"tgwp/global"
)

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	OpenID       string `json:"openid"`
	RefreshToken string `json:"refresh_token"` //refresh_token仅一次有效
}
type UserInfo struct {
	Ret         int    `json:"ret"`
	Msg         string `json:"msg"`
	Nickname    string `json:"nickname"`
	FigureurlQQ string `json:"figureurl_qq"` //头像
}
type QQUserInfo struct {
	OpenID   string `json:"openid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"` //头像
}

func getAccessToken(code string) (at AccessToken, err error) {
	qq := global.Config.QQ
	baseUrl, err := url.Parse("https://graph.qq.com/oauth2.0/token")
	if err != nil {
		err = errors.New("getAccessToken:url.Parse错误")
		return
	}
	p := url.Values{}
	p.Add("client_id", qq.AppID)
	p.Add("client_secret", qq.AppKey)
	p.Add("grant_type", "authorization_code")
	p.Add("redirect_uri", qq.RedirectUrl)
	p.Add("fmt", "json")
	p.Add("code", code)
	p.Add("need_openid", "1")
	baseUrl.RawQuery = p.Encode()
	res, err := http.Get(baseUrl.String())
	if err != nil {
		err = errors.New("getAccessToken:http.Get错误")
		return
	}
	byteData, err := io.ReadAll(res.Body)
	if err != nil {
		err = errors.New("getAccessToken:io.ReadAll错误")
		return
	}
	err = json.Unmarshal(byteData, &at)
	if err != nil {
		err = errors.New("getAccessToken:json.Unmarshal错误")
		return
	}
	if at.AccessToken == "" {
		//错误细分
		if strings.Contains(string(byteData), "code is reused error") {
			err = errors.New("getAccessToken:code失效")
			return
		}
		if strings.Contains(string(byteData), "client_secret is wrong") {
			err = errors.New("getAccessToken:secret错误")
			return
		}
		err = errors.New("getAccessToken:获取access_token失败")
		return
	}
	return
}
func getUserInfo(at AccessToken) (userinfo UserInfo, err error) {
	qq := global.Config.QQ
	baseUrl, err := url.Parse("https://graph.qq.com/user/get_user_info")
	if err != nil {
		err = errors.New("getUserInfo:url.Parse错误")
		return
	}
	p := url.Values{}
	p.Add("oauth_consumer_key", qq.AppID)
	p.Add("openid", at.OpenID)
	p.Add("access_token", at.AccessToken)
	baseUrl.RawQuery = p.Encode()
	res, err := http.Get(baseUrl.String())
	if err != nil {
		err = errors.New("getUserInfo:http.Get错误")
		return
	}
	byteData, err := io.ReadAll(res.Body)
	if err != nil {
		err = errors.New("getUserInfo:io.ReadAll错误")
		return
	}
	err = json.Unmarshal(byteData, &userinfo)
	if err != nil {
		err = errors.New("getUserInfo:json.Unmarshal错误")
		return
	}
	if userinfo.Ret != 0 {
		err = errors.New(fmt.Sprintf("%s:%s", "getUserInfo", userinfo.Msg))
		return
	}
	return
}
func GetUserInfo(code string) (info QQUserInfo, err error) {
	at, err := getAccessToken(code)
	if err != nil {
		return
	}
	u, err := getUserInfo(at)
	if err != nil {
		return
	}
	return QQUserInfo{
		OpenID:   at.OpenID,
		Nickname: u.Nickname,
		Avatar:   u.FigureurlQQ,
	}, nil
}
