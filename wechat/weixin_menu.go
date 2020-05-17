package wechat

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/wei193/component/common"
)

//base url
const (
	URLMENUGET            = "https://api.weixin.qq.com/cgi-bin/menu/get"
	URLMENUCREATE         = "https://api.weixin.qq.com/cgi-bin/menu/create"
	URLMENUADDCONDITIONAL = "https://api.weixin.qq.com/cgi-bin/menu/addconditional"
	URLMENUDELCONDITIONAL = "https://api.weixin.qq.com/cgi-bin/menu/delconditional"
	URLMENUDELETE         = "https://api.weixin.qq.com/cgi-bin/menu/delete"
)

//STButton 按钮元素
type STButton struct {
	Name      string     `json:"name,omitempty"`
	Type      string     `json:"type,omitempty"`
	Key       string     `json:"key,omitempty"`
	URL       string     `json:"url,omitempty"`
	Mediaid   string     `json:"media_id,omitempty"`
	SubButton []STButton `json:"sub_button,omitempty"`
}

//STMenu 普通菜单
type STMenu struct {
	Button []STButton `json:"button"`
}

//STCondMenu 自定义菜单
type STCondMenu struct {
	Button    []STButton  `json:"button"`
	Matchrule STMatchrule `json:"matchrule"`
	Menuid    string      `json:"menuid,omitempty"`
}

//STMatchrule 匹配规则
type STMatchrule struct {
	Groupid            string `json:"group_id,omitempty"`
	Sex                string `json:"sex,omitempty"`
	ClientPlatformType string `json:"client_platform_type,omitempty"`
	Country            string `json:"country,omitempty"`
	Province           string `json:"province,omitempty"`
	City               string `json:"city,omitempty"`
	Language           string `json:"language,omitempty"`
}

//STMenus 获取菜单
type STMenus struct {
	Menu     STMenu       `json:"menu"`
	CondMenu []STCondMenu `json:"conditionalmenu"`
}

//GetMenu 获取菜单
func (wx *Wechat) GetMenu() (data STMenus, err error) {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken

	req, err := http.NewRequest("GET", common.Param(URLMENUGET, param), nil)
	if err != nil {
		log.Println(err)
		return data, err
	}
	resp, err := common.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return data, err
	}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		log.Println(err)
		return data, err
	}
	return data, err
}

//CreatMenu 创建普通菜单
func (wx *Wechat) CreatMenu(menu STMenu) int {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken

	d, _ := json.Marshal(menu)
	d = bytes.Replace(d, []byte("\\u0026"), []byte("&"), -1)
	d = bytes.Replace(d, []byte("\\u003c"), []byte("<"), -1)
	d = bytes.Replace(d, []byte("\\u003e"), []byte(">"), -1)
	d = bytes.Replace(d, []byte("\\u003d"), []byte("="), -1)

	req, err := http.NewRequest("POST", common.Param("https://api.weixin.qq.com/cgi-bin/menu/create", param),
		bytes.NewReader(d))
	if err != nil {
		log.Println(err)
		return 0
	}
	_, err = common.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return 0
	}
	return 1
}

//CreatConditionalMenu 创建自定义菜单
func (wx *Wechat) CreatConditionalMenu(menu STCondMenu) (string, error) {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken

	d, _ := json.Marshal(menu)
	req, err := http.NewRequest("POST", common.Param("https://api.weixin.qq.com/cgi-bin/menu/addconditional", param),
		bytes.NewReader(d))
	if err != nil {
		log.Println(err)
		return "", err
	}
	resBody, err := common.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return "", err
	}
	type stTmp struct {
		Menuid int `json:"menuid"`
	}
	var temp stTmp
	err = json.Unmarshal(resBody, &temp)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return strconv.Itoa(temp.Menuid), nil
}

//DeleteConditionalMenu 删除自定义菜单
func (wx *Wechat) DeleteConditionalMenu(menuid string) int {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken

	type stTmp struct {
		Menuid string `json:"menuid"`
	}
	temp := stTmp{menuid}
	d, _ := json.Marshal(temp)
	req, err := http.NewRequest("GET", common.Param("https://api.weixin.qq.com/cgi-bin/menu/delconditional", param), bytes.NewReader(d))
	if err != nil {
		log.Println(err)
		return 0
	}
	_, err = common.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return 0
	}
	return 1
}

//DeleteAllMenu 删除所有菜单
func (wx *Wechat) DeleteAllMenu() int {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken

	req, err := http.NewRequest("GET", common.Param("https://api.weixin.qq.com/cgi-bin/menu/delete", param), nil)
	if err != nil {
		log.Println(err)
		return 0
	}
	_, err = common.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return 0
	}
	//	log.Println(string(resp))
	return 1
}

//https://api.weixin.qq.com/cgi-bin/get_current_selfmenu_info?access_token=ACCESS_TOKEN
