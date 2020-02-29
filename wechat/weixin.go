package wechat

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
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

//JSONError  微信错误
type JSONError struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

//XMLError  微信错误
type XMLError struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
}

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
func New(Appid, Appsecret, Token, Encodingaeskey string) *Wechat {
	wx := &Wechat{
		Appid:          Appid,
		Appsecret:      Appsecret,
		Token:          Token,
		Encodingaeskey: Encodingaeskey,
	}
	return wx
}

//SetMch 设置商户
func (wx *Wechat) SetMch(mchid, paykey, certpath, keypath, capath string) {
	if wx.Mch == nil {
		wx.Mch = &MchInfo{
			MchID:  mchid,
			PayKey: paykey,
		}
	} else {
		wx.Mch.MchID = mchid
		wx.Mch.PayKey = paykey
	}
	wx.Mch._tlsConfig, _ = getTLSConfig(certpath, keypath, capath)
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

	req, err := http.NewRequest("GET", Param(URLTOKEN, param), nil)

	resBody, err := requsetJSON(req, -1)
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
	_, err = requsetJSON(req, 0)
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
	req, err := http.NewRequest("GET", Param(URLGETTICKET, param), nil)
	if err != nil {
		return err
	}

	resBody, err := requsetJSON(req, 0)
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
	return SignSha1(data)
}

//发送微信请求
func requsetJSON(req *http.Request, tflag int) ([]byte, error) {

	resBody, err := requset(req)
	if err != nil {
		return nil, err
	}
	var errcode JSONError
	err = json.Unmarshal(resBody, &errcode)
	if err == nil && errcode.Errcode != 0 {
		return resBody, errors.New(string(resBody))
	}
	return resBody, nil
}

func requsetXML(req *http.Request, tflag int, isXML ...bool) ([]byte, error) {
	resBody, err := requset(req)
	if err != nil {
		return nil, err
	}
	if len(isXML) == 1 && !isXML[0] {
		return resBody, nil
	}
	var errcode XMLError
	err = xml.Unmarshal(resBody, &errcode)
	if err != nil ||
		errcode.ReturnCode != "SUCCESS" ||
		errcode.ResultCode != "SUCCESS" ||
		errcode.ErrCode != "" {
		return resBody, errors.New(string(resBody))
	}

	return resBody, nil
}

func requset(req *http.Request) ([]byte, error) {
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (wx *Wechat) httpsRequsetXML(req *http.Request, tflag int, isXML ...bool) ([]byte, error) {
	resBody, err := wx.httpsRequset(req)
	if err != nil {
		return nil, err
	}
	if len(isXML) == 1 && !isXML[0] {
		return resBody, nil
	}
	var errcode XMLError
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

//Param 生成请求url参数
func Param(urlBase string, P map[string]string) string {
	for k, v := range P {
		if strings.Index(urlBase, "?") != -1 {
			urlBase += "&"
		} else {
			urlBase += "?"
		}
		urlBase += k
		urlBase += "="
		urlBase += url.QueryEscape(v)
	}
	return urlBase
}

//SignSha1 hash1签名
func SignSha1(data map[string]interface{}) string {
	var keys []string
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	str1 := ""
	for i := range keys {
		val := data[keys[i]]
		var str string
		switch val.(type) {
		case string:
			str = val.(string)
		case bool:
			str = strconv.FormatBool(val.(bool))
		case int:
			str = strconv.Itoa(val.(int))
		case int64:
			str = strconv.FormatInt(val.(int64), 10)
		case []byte:
			str = string(val.([]byte))
		default:
			continue
		}
		if len(str) == 0 {
			continue
		}
		if len(str1) != 0 {
			str1 += "&"
		}
		str1 += keys[i] + "=" + str
	}
	t := sha1.New()
	io.WriteString(t, str1)
	return fmt.Sprintf("%x", t.Sum(nil))
}

//XMLSignMd5 MD5签名
func XMLSignMd5(data interface{}, key string) string {
	k := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	var keys []string
	m := make(map[string]interface{})
	for i := 0; i < k.NumField(); i++ {
		chKey := k.Field(i).Tag.Get("xml")
		tmpStr := strings.Split(chKey, ",")
		keys = append(keys, tmpStr[0])
		if len(tmpStr) > 1 && tmpStr[1] == "omitempty" && IsEmptyValue(v.Field(i)) {
			continue
		}
		m[tmpStr[0]] = v.Field(i).Interface()
	}
	sort.Strings(keys)
	str1 := ""
	for i := range keys {
		val := m[keys[i]]
		var str string
		switch val.(type) {
		case string:
			str = val.(string)
		case bool:
			str = strconv.FormatBool(val.(bool))
		case int:
			str = strconv.Itoa(val.(int))
		case int64:
			str = strconv.FormatInt(val.(int64), 10)
		case []byte:
			str = string(val.([]byte))
		default:
			continue
		}
		if len(str) == 0 {
			continue
		}
		if len(str1) != 0 {
			str1 += "&"
		}
		str1 += keys[i] + "=" + str
	}
	str1 += "&key=" + key
	t := md5.New()
	io.WriteString(t, str1)
	return fmt.Sprintf("%X", t.Sum(nil))
}

//CheckSignMd5 检查数据的MD5是否正确
func CheckSignMd5(data interface{}, signName, key string) (string, bool) {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Map {
		return "", false
	}
	var keys []string
	var sign string
	var signStr string

	m := make(map[string]string)
	for _, t := range val.MapKeys() {
		k := fmt.Sprint(t.Interface())
		if k == signName {
			sign = fmt.Sprint(val.MapIndex(t).Interface())
		}
		keys = append(keys, k)
		m[k] = fmt.Sprint(val.MapIndex(t).Interface())
	}
	sort.Strings(keys)
	for i, k := range keys {
		if i != 0 {
			signStr += "&"
		}
		signStr += k + "=" + m[k]
	}
	signStr += "&key=" + key
	tMd5 := md5.New()
	io.WriteString(tMd5, signStr)
	return fmt.Sprintf("%X", tMd5.Sum(nil)), fmt.Sprintf("%X", tMd5.Sum(nil)) == sign
}

// IsEmptyValue 判断值是否为空
func IsEmptyValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
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

//RandomStr 随机字符串
/*Random = 0  // 纯数字
Random = 1  // 小写字母
Random = 2  // 大写字母
Random   = 3  // 数字、大小写字母*/
func RandomStr(size int, Random int) string {
	iRandom, Randoms, result := Random, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	iAll := Random > 2 || Random < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if iAll { // random ikind
			iRandom = rand.Intn(3)
		}
		scope, base := Randoms[iRandom][0], Randoms[iRandom][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}

func getTLSConfig(certpath, keypath, capath string) (*tls.Config, error) {

	cert, err := tls.LoadX509KeyPair(certpath, keypath)
	if err != nil {
		return nil, err
	}

	caData, err := ioutil.ReadFile(capath)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}, nil
}
