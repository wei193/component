package component

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//JSONError  微信错误
type JSONError struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func requsetJosn(req *http.Request) ([]byte, error) {
	client := &http.Client{Timeout: 60 * time.Second}
	// b, _ := httputil.DumpRequest(req, true)
	// log.Println(string(b))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	var errcode JSONError
	err = json.Unmarshal(res, &errcode)
	if err == nil && errcode.Errcode != 0 {
		return res, errors.New(string(res))
	}
	return res, err
}

//createRequset 生成请求requset参数
func createRequset(surl, method string, p map[string]string, d interface{}) (req *http.Request, err error) {
	buf, err := json.Marshal(d)
	if d != nil && err != nil {
		return nil, err
	}

	val := make(url.Values)
	for k, v := range p {
		val.Add(k, v)
	}
	ContentType := ""
	switch method {
	case "GET":
		if strings.Index(surl, "?") != -1 {
			surl += "&" + val.Encode()
		} else {
			surl += "?" + val.Encode()
		}
	case "POST":
		if d == nil {
			buf = []byte(val.Encode())
			ContentType = "application/x-www-form-urlencoded"
		} else {
			ContentType = "application/json"
			if strings.Index(surl, "?") != -1 {
				surl += "&" + val.Encode()
			} else {
				surl += "?" + val.Encode()
			}
		}
	default:
		if strings.Index(surl, "?") != -1 {
			surl += "&" + val.Encode()
		} else {
			surl += "?" + val.Encode()
		}
	}

	req, err = http.NewRequest(method, surl, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	if ContentType != "" {
		req.Header.Add("Content-Type", ContentType)
	}
	return req, err
}
