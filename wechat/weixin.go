package wechat

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/wei193/component/common"
)

//基础定义
const (
	Text     = "text"
	Location = "location"
	Image    = "image"
	Link     = "link"
	Event    = "event"
	Music    = "music"
	News     = "news"

	TOKENIGNORE   = -1
	TOKENRETURN   = 0
	TOKENCONTINUE = 1

	URLGETCALLBACKIP = "https://api.weixin.qq.com/cgi-bin/getcallbackip"
	URLTOKEN         = "https://api.weixin.qq.com/cgi-bin/token"
	URLGETTICKET     = "https://api.weixin.qq.com/cgi-bin/ticket/getticket"
)

//ResAccessToken ResAccessToken
type ResAccessToken struct {
	AccessToken string `json:"access_token"`
	Expiresin   int    `json:"expires_in"`
	Errcode     string `json:"errcode"`
}

//ResUserToken 用户Token
type ResUserToken struct {
	AccessToken string `json:"access_token"`
	Expiresin   int    `json:"expires_in"`
	Openid      string `json:"openid"`
	Scope       string `json:"scope"`
	Errmsg      string `json:"Errmsg"`
}

type resJsTicket struct {
	Errcode   int    `json:"errcode"`
	Ticket    string `json:"ticket"`
	Errmsg    string `json:"errmsg"`
	Expiresin int    `json:"expires_in"`
}

//Wechat 微信接口
type Wechat struct {
	Appid              string
	Appsecret          string
	Token              string
	Encodingaeskey     string
	AccessToken        string
	AccessTokenExpires int64
	JsapiTicket        string
	JsapiTokenTime     int64
	JsapiTokenExpires  int64
	AuthorizedDomain   string
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

//New 创建wechat
func New(Appid, Appsecret, Token, Encodingaeskey, AuthorizedDomain string) *Wechat {
	wx := &Wechat{
		Appid:            Appid,
		Appsecret:        Appsecret,
		Token:            Token,
		Encodingaeskey:   Encodingaeskey,
		AuthorizedDomain: AuthorizedDomain,
	}
	return wx
}

//SetMch 设置商户
func (wx *Wechat) SetMch(mchid, paykey, certpath, keypath, capath string) (err error) {
	if wx.Mch == nil {
		wx.Mch = &MchInfo{
			MchID:  mchid,
			PayKey: paykey,
		}
	} else {
		wx.Mch.MchID = mchid
		wx.Mch.PayKey = paykey
	}
	wx.Mch._tlsConfig, err = common.GetTLSConfig(certpath, keypath, capath)
	return err
}

//GetAccessToken 获取 access_token
func (wx *Wechat) GetAccessToken() (err error) {
	if wx.Appsecret == "" {
		return errors.New("no secret")
	}
	param := make(map[string]string)
	param["grant_type"] = "client_credential"
	param["appid"] = wx.Appid
	param["secret"] = wx.Appsecret

	req, err := http.NewRequest("GET", common.Param(URLTOKEN, param), nil)
	if err == nil {
		return err
	}
	resBody, err := common.RequsetJSON(req, -1)
	if err != nil {
		log.Println(err)
		return err
	}
	var accToken ResAccessToken
	err = json.Unmarshal(resBody, &accToken)
	if err != nil {
		log.Println(err)
		return err
	}
	if accToken.Errcode != "" {
		return errors.New("获取Token  失败")
	}

	wx.AccessTokenExpires = time.Now().Unix() + int64(accToken.Expiresin)
	wx.AccessToken = accToken.AccessToken

	return nil
}

//CheckAccessToken 检查微信access_token有效性
func (wx *Wechat) CheckAccessToken() (err error) {
	req, err := http.NewRequest("GET", URLGETCALLBACKIP+"?access_token="+
		wx.AccessToken, nil)
	if err == nil {
		return err
	}
	_, err = common.RequsetJSON(req, 0)
	if err != nil {
		return err
	}
	return nil
}

//GetJsapiTicket 获取js的jsapi_ticket
func (wx *Wechat) GetJsapiTicket() (err error) {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken
	param["type"] = "jsapi"
	req, err := http.NewRequest("GET", common.Param(URLGETTICKET, param), nil)
	if err != nil {
		return err
	}

	resBody, err := common.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return err
	}
	var tmpTick resJsTicket
	err = json.Unmarshal(resBody, &tmpTick)
	if err != nil {
		log.Println(err)
		return err
	} else if tmpTick.Errcode == 0 {
		wx.JsapiTokenTime = time.Now().Unix()
		wx.JsapiTicket = tmpTick.Ticket
		wx.JsapiTokenExpires = time.Now().Unix() + int64(tmpTick.Expiresin)
		return nil
	} else {
		return errors.New(tmpTick.Errmsg)
	}
}

//CreateJsSignature 创建jsapi_ticket签名
func (wx *Wechat) CreateJsSignature(url, noncestr string, timestamp int64, data map[string]interface{}) string {
	if data == nil {
		data = make(map[string]interface{})
	}
	data["url"] = url
	data["noncestr"] = noncestr
	data["jsapi_ticket"] = wx.JsapiTicket
	data["timestamp"] = strconv.FormatInt(timestamp, 10)
	return common.SignSha1(data)
}

func (wx *Wechat) httpsRequsetXML(req *http.Request, tflag int, isXML ...bool) ([]byte, error) {
	resBody, err := wx.httpsRequset(req)
	if err != nil {
		return nil, err
	}
	if len(isXML) == 1 && !isXML[0] {
		return resBody, nil
	}
	var errcode common.XMLError
	err = xml.Unmarshal(resBody, &errcode)
	if err != nil ||
		errcode.ReturnCode != "SUCCESS" ||
		errcode.ResultCode != "SUCCESS" ||
		errcode.ErrCode != "" {
		return resBody, errors.New(string(resBody))
	}

	return resBody, nil
}

func (wx *Wechat) httpsRequset(req *http.Request) ([]byte, error) {
	tlsConfig := wx.Mch._tlsConfig
	if tlsConfig == nil {
		return nil, errors.New("init tls Config Error")
	}
	tr := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

//httpsPost  HttpsPost请求
func (wx *Wechat) httpsPost(url string, xmlContent []byte, ContentType string) (*http.Response, error) {
	tlsConfig := wx.Mch._tlsConfig
	if tlsConfig == nil {
		return nil, errors.New("init tls Config Error")
	}
	tr := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: tr}
	return client.Post(url,
		ContentType,
		bytes.NewBuffer(xmlContent))
}

//GetRedirectUri 获取登录连接
func (wx *Wechat) GetDefaultRedirectUri(path, scope, state string) string {
	return wx.GetRedirectUri(wx.AuthorizedDomain+path, scope, state)
}

//GetRedirectUri 获取登录连接
func (wx *Wechat) GetRedirectUri(uri, scope, state string) string {
	var tpl, _ = url.Parse("https://open.weixin.qq.com/connect/oauth2/authorize")
	params := url.Values{}
	params.Add("appid", wx.Appid)
	params.Add("redirect_uri", uri)
	params.Add("response_type", "code")
	params.Add("scope", scope)
	params.Add("state", state)
	tpl.RawQuery = params.Encode()
	return tpl.String()
}

//DecodeRequest 解析请求
func DecodeRequest(data []byte) (req *STMsgRequest, err error) {
	req = &STMsgRequest{}
	if err = xml.Unmarshal(data, req); err != nil {
		return
	}
	req.CreateTime *= time.Second
	return
}
