package component

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

//Component 第三平台信息
type Component struct {
	ComponentAppid        string
	ComponentAppsecret    string
	ComponentToken        string
	ComponentAesKey       string
	ComponentVerifyTicket string
	ComponentAccessToken  string
	AccessTokenExpires    int64
	AESKey                []byte
}

// XEncryptMsg 消息
type XEncryptMsg struct {
	AppID   string `xml:"AppId"`
	Encrypt string `xml:"Encrypt"`
}

//XCEvent ticket协议和推送授权相关通知
type XCEvent struct {
	AppID                        string `xml:"AppId"`
	CreateTime                   string `xml:"CreateTime"`
	InfoType                     string `xml:"InfoType"`
	ComponentVerifyTicket        string `xml:"ComponentVerifyTicket"`
	AuthorizerAppid              string `xml:"AuthorizerAppid"`
	AuthorizationCode            string `xml:"AuthorizationCode"`
	AuthorizationCodeExpiredTime int    `xml:"AuthorizationCodeExpiredTime"`
	PreAuthCode                  string `xml:"PreAuthCode"`
}

//JPreAuthCode 预授权码
type JPreAuthCode struct {
	PreAuthCode string `json:"pre_auth_code"`
	ExpiresIn   int    `json:"expires_in"`
}

//JComponenAccessToken JComponenAccessToken
type JComponenAccessToken struct {
	ComponentAccessToken string `json:"component_access_token"`
	ExpiresIn            int    `json:"expires_in"`
}

//NewComponent 新建第三方平台
func NewComponent(appid, appsecret, token, aeskey, verifyticket, accesstoken string, expires int64) (component *Component, err error) {

	AESKey, err := base64.StdEncoding.DecodeString(aeskey + "=")
	if err != nil {
		beego.Error("Component Init Error")
		return nil, err
	}

	component = &Component{
		ComponentAppid:        appid,
		ComponentAppsecret:    appsecret,
		ComponentToken:        token,
		ComponentAesKey:       aeskey,
		ComponentVerifyTicket: verifyticket,
		ComponentAccessToken:  accesstoken,
		AccessTokenExpires:    expires,
		AESKey:                AESKey,
	}
	return component, nil
}

//GetPreAuthCode 获取预授权码
func (c *Component) GetPreAuthCode() (authcode JPreAuthCode, err error) {
	type st struct {
		ComponentAppid string `json:"component_appid"`
	}
	d := st{
		ComponentAppid: c.ComponentAppid,
	}
	req, err := createRequset("https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?component_access_token="+c.ComponentAccessToken,
		"POST", nil, d)
	if err != nil {
		return authcode, err
	}
	res, err := requsetJosn(req)
	if err != nil {
		return authcode, err
	}
	log.Println(string(res))

	err = json.Unmarshal(res, &authcode)
	if err != nil {
		return authcode, err
	}
	return
}

//AuthEventTicket 授权事件接收处理
func (c *Component) AuthEventTicket(msg, signature, timestamp, nonce string) (event *XCEvent, err error) {
	isok := c.CheckSignature(signature, msg, timestamp, nonce)
	if !isok {
		return nil, errors.New("Check Signature Error")
	}
	if len(msg) == 0 {
		return nil, errors.New("Msg Error")
	}
	result, err := c.MsgDecrypt(msg)
	if err != nil {
		return nil, errors.New("Msg Decrypt Error")
	}

	s := strings.Index(result, "<")
	e := strings.LastIndex(result, ">")
	if s == -1 || e == -1 {
		return nil, errors.New("Msg Font Error")
	}

	event = new(XCEvent)
	err = xml.Unmarshal([]byte(result[s:e+1]), event)
	if err != nil {
		beego.Error(err)
		return nil, errors.New("xml Unmarshal Error")
	}
	switch event.InfoType {
	case "component_verify_ticket":
		c.ComponentVerifyTicket = event.ComponentVerifyTicket
	}
	return event, nil
}

//GetComponentAccessToken 获取第三方AccessToken
func (c *Component) GetComponentAccessToken() (token *JComponenAccessToken, err error) {
	type st struct {
		ComponentAppid        string `json:"component_appid"`
		ComponentAppsecret    string `json:"component_appsecret"`
		ComponentVerifyTicket string `json:"component_verify_ticket"`
	}
	d := st{
		ComponentAppid:        c.ComponentAppid,
		ComponentAppsecret:    c.ComponentAppsecret,
		ComponentVerifyTicket: c.ComponentVerifyTicket,
	}

	req, err := createRequset("https://api.weixin.qq.com/cgi-bin/component/api_component_token",
		"POST", nil, d)
	if err != nil {
		return nil, err
	}
	res, err := requsetJosn(req)
	if err != nil {
		return nil, err
	}

	log.Println(string(res))

	token = new(JComponenAccessToken)
	err = json.Unmarshal(res, token)
	if err != nil {
		return nil, err
	}

	c.ComponentAccessToken = token.ComponentAccessToken
	c.AccessTokenExpires = time.Now().Unix() + int64(token.ExpiresIn)
	return token, nil
}

//QueryAuth 使用授权码换取公众号或小程序的接口调用凭据和授权信息
func (c *Component) QueryAuth(code string) (authorizer *Authorizer, err error) {
	type st struct {
		ComponentAppid    string `json:"component_appid"`
		AuthorizationCode string `json:"authorization_code"`
	}
	d := st{
		ComponentAppid:    c.ComponentAppid,
		AuthorizationCode: code,
	}
	req, err := createRequset("https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token="+c.ComponentAccessToken,
		"POST", nil, d)
	if err != nil {
		return nil, err
	}
	res, err := requsetJosn(req)
	if err != nil {
		return nil, err
	}

	log.Println(string(res))
	var auth JAuthorizer
	err = json.Unmarshal(res, &auth)
	if err != nil {
		return nil, err
	}

	authorizer = &Authorizer{
		Component:              c,
		AuthorizerAppid:        auth.AuthorizationInfo.AuthorizerAppid,
		AuthorizerAccessToken:  auth.AuthorizationInfo.AuthorizerAccessToken,
		AccessTokenExpires:     time.Now().Unix() + int64(auth.AuthorizationInfo.ExpiresIn),
		AuthorizerRefreshToken: auth.AuthorizationInfo.AuthorizerRefreshToken,
	}
	return authorizer, nil
}
