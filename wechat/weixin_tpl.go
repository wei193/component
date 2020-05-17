// Copyright 2019 wei_193 Author. All Rights Reserved.
//
// 微信模板消息管理

package wechat

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/wei193/component/common"
)

//设置所属行业https://api.weixin.qq.com/cgi-bin/template/api_set_industry?access_token=ACCESS_TOKEN
//获取设置的行业信息https://api.weixin.qq.com/cgi-bin/template/get_industry?access_token=ACCESS_TOKEN
//获得模板IDhttps://api.weixin.qq.com/cgi-bin/template/api_add_template?access_token=ACCESS_TOKEN
//获取模板列表https://api.weixin.qq.com/cgi-bin/template/get_all_private_template?access_token=ACCESS_TOKEN
//删除模板https://api,weixin.qq.com/cgi-bin/template/del_private_template?access_token=ACCESS_TOKEN
//

//访问地址
const (
	TEMPLATESENDURL = "https://api.weixin.qq.com/cgi-bin/message/template/send"
)

//TemplateData TemplateData
type TemplateData struct {
	Value string `json:"value,omitempty"`
	Color string `json:"color,omitempty"`
}

//Template Template
type Template struct {
	Touser     string                  `json:"touser"`
	Templateid string                  `json:"template_id"`
	URL        string                  `json:"url"`
	Data       map[string]TemplateData `json:"data,json"`
}

//SendTemplate 发送模板消息https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=ACCESS_TOKEN
func (wx *Wechat) SendTemplate(data string) (string, error) {
	buf, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", TEMPLATESENDURL+"?access_token="+
		wx.AccessToken, bytes.NewReader(buf))
	res, err := common.RequsetJSON(req, 0)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

//SendTemplateToUser 发送模板消息到用户
func (wx *Wechat) SendTemplateToUser(touser, templateid, url string,
	data map[string]TemplateData) (string, error) {
	tpl := Template{
		Touser:     touser,
		Templateid: templateid,
		URL:        url,
		Data:       data,
	}
	rdata, _ := json.Marshal(tpl)
	req, err := http.NewRequest("POST", TEMPLATESENDURL+"?access_token="+
		wx.AccessToken, bytes.NewReader(rdata))
	res, err := common.RequsetJSON(req, 0)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
