// Copyright 2020 wei_193 Author. All Rights Reserved.
//
// 小程序内容安全接口

package miniprogram

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/wei193/component/base"
)

//CheckResult 检查结果
type CheckResult struct {
	TraceID string `json:"trace_id"`
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

//ImgSecCheck 图片安全检查
func (mini *MiniProgram) ImgSecCheck(media string) (res CheckResult, err error) {
	param := make(map[string]string)
	param["access_token"] = mini.AccessToken

	file, err := os.Open(media)
	if err != nil {
		return res, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	req, err := http.NewRequest("POST",
		base.Param("https://api.weixin.qq.com/wxa/img_sec_check", param),
		body)

	resBody, err := base.Requset(req)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		return res, err
	}
	if res.Errcode != 0 {
		return res, errors.New(res.Errmsg)
	}
	return res, nil
}

//MediaCheckAsync 异步检查媒体是否存在违规信息
func (mini *MiniProgram) MediaCheckAsync(url string, typ int) (res CheckResult, err error) {
	param := make(map[string]string)
	param["access_token"] = mini.AccessToken

	tmp := make(map[string]interface{})
	tmp["media_url"] = url
	tmp["media_type"] = typ

	rdata, _ := json.Marshal(tmp)

	req, err := http.NewRequest("POST",
		base.Param("https://api.weixin.qq.com/wxa/media_check_async", param),
		bytes.NewReader(rdata))
	resBody, err := base.Requset(req)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		return res, err
	}
	if res.Errcode != 0 {
		return res, errors.New(res.Errmsg)
	}
	return res, nil
}

//MsgSecCheck 文本检查
func (mini *MiniProgram) MsgSecCheck(content string) (res CheckResult, err error) {
	param := make(map[string]string)
	param["access_token"] = mini.AccessToken

	tmp := make(map[string]interface{})
	tmp["content"] = content
	rdata, _ := json.Marshal(tmp)

	req, err := http.NewRequest("POST",
		base.Param("https://api.weixin.qq.com/wxa/msg_sec_check", param),
		bytes.NewReader(rdata))

	resBody, err := base.Requset(req)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		return res, err
	}
	if res.Errcode != 0 {
		return res, errors.New(res.Errmsg)
	}
	return res, nil
}
