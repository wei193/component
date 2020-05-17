package common

import (
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
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

//JSONError  微信错误
type JSONError struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

//XMLError  微信错误
type XMLError struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code,omitempty"`
	ReturnMsg  string   `xml:"return_msg,omitempty"`
	ResultCode string   `xml:"result_code,omitempty"`
	ErrCode    string   `xml:"err_code,omitempty"`
	ErrCodeDes string   `xml:"err_code_des,omitempty"`
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

//RandomStr 随机字符串
/*Random = 0  // 纯数字
Random = 1  // 小写字母
Random = 2  // 大写字母
Random   = 3  // 数字、大小写字母*/
func RandomStr(size int, Random int) string {
	iRandom, Randoms, result := Random, [][]int{{10, 48}, {26, 97}, {26, 65}}, make([]byte, size)
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

//GetTLSConfig 获取TLS Config
func GetTLSConfig(certpath, keypath, capath string) (*tls.Config, error) {

	cert, err := tls.LoadX509KeyPair(certpath, keypath)
	if err != nil {
		return nil, err
	}

	if capath != "" {
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

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
	}, nil
}

//RequsetJSON 发送微信请求
func RequsetJSON(req *http.Request, tflag int) ([]byte, error) {

	resBody, err := Requset(req)
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

//RequsetXML 发送微信请求
func RequsetXML(req *http.Request, tflag int, isXML ...bool) ([]byte, error) {
	resBody, err := Requset(req)
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

//Requset Requset
func Requset(req *http.Request) ([]byte, error) {
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
