package miniprogram

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/wei193/component/base"
)

//MiniProgram 小程序接口
type MiniProgram struct {
	Appid              string
	Appsecret          string
	Token              string
	Encodingaeskey     string
	AccessToken        string
	AccessTokenExpires int64
	Mch                *MchInfo
}

//MchInfo 微信商户信息
type MchInfo struct {
	MchID      string
	PayKey     string
	CertPath   string
	KeyPath    string
	CaPath     string
	_tlsConfig *tls.Config
}

//AccessToken ResAccessToken
type AccessToken struct {
	AccessToken string `json:"access_token"`
	Expiresin   int    `json:"expires_in"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}

//GetAccessToken 获取 access_token
func (mini *MiniProgram) GetAccessToken() (err error) {
	if mini.Appsecret == "" {
		return errors.New("no secret")
	}
	param := make(map[string]string)
	param["grant_type"] = "client_credential"
	param["appid"] = mini.Appid
	param["secret"] = mini.Appsecret

	req, err := http.NewRequest("GET", base.Param("https://api.weixin.qq.com/cgi-bin/token", param), nil)

	resBody, err := base.Requset(req)
	if err != nil {
		return err
	}
	var acc AccessToken
	err = json.Unmarshal(resBody, &acc)
	if err != nil {
		return err
	}
	if acc.Errcode != 0 {
		return errors.New(acc.Errmsg)
	}
	mini.AccessTokenExpires = time.Now().Unix() + int64(acc.Expiresin)
	mini.AccessToken = acc.AccessToken

	return nil
}

//Getpaidunionid 微信用户支付以后获取用户UnionId
func (mini *MiniProgram) Getpaidunionid(openid, transactionid, outtradeno string) (id string, err error) {
	param := make(map[string]string)
	param["access_token"] = mini.AccessToken
	param["openid"] = openid
	if transactionid == "" {
		param["transaction_id"] = transactionid
	}
	if outtradeno == "" {
		param["mch_id"] = mini.Mch.MchID
		param["out_trade_no"] = outtradeno
	}

	req, err := http.NewRequest("GET", base.Param("https://api.weixin.qq.com/cgi-bin/token", param), nil)

	resBody, err := base.Requset(req)
	if err != nil {
		return "", err
	}
	type st struct {
		Unionid string `json:"unionid"`
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}
	var data st
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		return "", err
	}
	if data.Errcode != 0 {
		return "", errors.New(data.Errmsg)
	}

	return data.Unionid, nil
}

//MiniSession 小程序会话
type MiniSession struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
}

//Code2Session 通过Code获取session_key
func (mini *MiniProgram) Code2Session(code string) (s MiniSession, err error) {
	if mini.Appsecret == "" {
		return s, errors.New("no secret")
	}
	param := make(map[string]string)
	param["grant_type"] = "authorization_code"
	param["appid"] = mini.Appid
	param["secret"] = mini.Appsecret
	param["js_code"] = code

	req, err := http.NewRequest("GET", base.Param("https://api.weixin.qq.com/cgi-bin/token", param), nil)

	resBody, err := base.Requset(req)
	if err != nil {
		return s, err
	}
	err = json.Unmarshal(resBody, &s)
	if err != nil {
		return s, err
	}
	if s.Errcode != 0 {
		return s, errors.New(s.Errmsg)
	}
	return s, nil
}
