package wechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/wei193/component/common"
)

//base url
const (
	URLUSERACCESSTOKEN        = "https://api.weixin.qq.com/sns/oauth2/access_token"
	URLUSERUSERINFO           = "https://api.weixin.qq.com/sns/userinfo"
	URLUSERINFO               = "https://api.weixin.qq.com/cgi-bin/user/info"
	URLUSERBATCHGET           = "https://api.weixin.qq.com/cgi-bin/user/info/batchget"
	URLUSERGET                = "https://api.weixin.qq.com/cgi-bin/user/get"
	URLUSERGROUPCREATE        = "https://api.weixin.qq.com/cgi-bin/groups/create"
	URLUSERGROUPGET           = "https://api.weixin.qq.com/cgi-bin/groups/get"
	URLUSERGROUPGETID         = "https://api.weixin.qq.com/cgi-bin/groups/getid"
	URLUSERGROUPUPDATE        = "https://api.weixin.qq.com/cgi-bin/groups/update"
	URLUSERGROUPMEMBERSUPDATE = "https://api.weixin.qq.com/cgi-bin/groups/members/update"
	URLUSERTAGCREATE          = "https://api.weixin.qq.com/cgi-bin/tags/create"
	URLUSERTAGGET             = "https://api.weixin.qq.com/cgi-bin/tags/get"
	URLUSERTAGUPDATE          = "https://api.weixin.qq.com/cgi-bin/tags/update"
	URLUSERTAGDELETE          = "https://api.weixin.qq.com/cgi-bin/tags/delete"
	URLUSERTAGGETIDLIST       = "https://api.weixin.qq.com/cgi-bin/tags/getidlist"
	URLUSERTAGBATCHTAGGING    = "https://api.weixin.qq.com/cgi-bin/tags/members/batchtagging"
	URLUSERTAGBATCHUNTAGGING  = "https://api.weixin.qq.com/cgi-bin/tags/members/batchuntagging"
	URLUSERUPDATEREMARK       = "https://api.weixin.qq.com/cgi-bin/user/info/updateremark"
	URLJSCODE2SESSION         = "https://api.weixin.qq.com/sns/jscode2session"
	// GET https://api.weixin.qq.com/sns/jscode2session?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code

)

//STOpenids STOpenids
type STOpenids struct {
	Openid []string `json:"openid"`
}

//STOpenid STOpenid
type STOpenid struct {
	Openid string `json:"openid"`
}

//STGroups STGroups
type STGroups struct {
	Groupsid   int    `json:"id"`
	Groupsname string `json:"name"`
}

//STTag STTag
type STTag struct {
	Tagid   int    `json:"id,omitempty"`
	Tagname string `json:"name,omitempty"`
}

//STUserInfo STUserInfo
type STUserInfo struct {
	Subscribe     int    `json:"subscribe"`
	Openid        string `json:"openid"`
	Nickname      string `json:"nickname"`
	Sex           int    `json:"sex"`
	Language      string `json:"language"`
	City          string `json:"city"`
	Province      string `json:"province"`
	Country       string `json:"country"`
	Headimgurl    string `json:"headimgurl"`
	SubscribeTime int64  `json:"subscribe_time"`
	Remark        string `json:"remark"`
	Groupid       int    `json:"groupid"`
	TagidList     []int  `json:"tagid_list"`
	Unionid       string `json:"unionid"`
}

//STUserList STUserList
type STUserList struct {
	Total      int `json:"total"`
	Count      int `json:"count"`
	Data       STOpenids
	NextOpenid string `json:"next_openid"`
	Errcode    int    `json:"errcode"`
}

//GetUserInfoToken GetUserInfo 获取用户信息by user_token
func GetUserInfoToken(accessToken, openid string) (userInfo STUserInfo, err error) {
	param := make(map[string]string)
	param["access_token"] = accessToken
	param["openid"] = openid

	req, err := http.NewRequest("GET", common.Param("https://api.weixin.qq.com/sns/userinfo?lang=zh_CN", param), nil)
	resBody, err := common.RequsetJSON(req, TOKENIGNORE)
	if err != nil {
		log.Println(err)
		return userInfo, err
	}
	err = json.Unmarshal(resBody, &userInfo)
	if err != nil {
		log.Println(err)
		return userInfo, err
	}
	return userInfo, nil
}

//GetUserToken 获取用户user_token
func (wx *Wechat) GetUserToken(Code string) (uToken ResUserToken, err error) {

	param := make(map[string]string)
	param["appid"] = wx.Appid
	param["secret"] = wx.Appsecret
	param["code"] = Code
	param["grant_type"] = "authorization_code"
	req, err := http.NewRequest("GET", common.Param("https://api.weixin.qq.com/sns/oauth2/access_token", param), nil)
	if err != nil {
		return uToken, err
	}
	resBody, err := common.RequsetJSON(req, 1)
	if err != nil {
		log.Println(err)
		return uToken, err
	}
	err = json.Unmarshal(resBody, &uToken)
	if err != nil {
		log.Println(err)
		return uToken, err
	}
	return uToken, nil
}

//GetUserInfoToken GetUserInfo 获取用户信息by user_token
func (wx *Wechat) GetUserInfoToken(accessToken, openid string) (userInfo STUserInfo, err error) {
	param := make(map[string]string)
	param["access_token"] = accessToken
	param["openid"] = openid

	req, err := http.NewRequest("GET", common.Param("https://api.weixin.qq.com/sns/userinfo?lang=zh_CN", param), nil)
	resBody, err := common.RequsetJSON(req, TOKENIGNORE)
	if err != nil {
		log.Println(err)
		return userInfo, err
	}
	err = json.Unmarshal(resBody, &userInfo)
	if err != nil {
		log.Println(err)
		return userInfo, err
	}
	return userInfo, nil
}

//GetUserInfo 获取用户信息by ac_token
func (wx *Wechat) GetUserInfo(openid string) (userInfo STUserInfo, err error) {

	param := make(map[string]string)
	param["access_token"] = wx.AccessToken
	param["openid"] = openid
	req, err := http.NewRequest("GET", common.Param("https://api.weixin.qq.com/cgi-bin/user/info?lang=zh_CN", param), nil)
	resBody, err := common.RequsetJSON(req, 0)
	if err != nil {
		return userInfo, err
	}
	err = json.Unmarshal(resBody, &userInfo)
	if err != nil {
		return userInfo, err
	}

	return userInfo, nil
}

//GetUsers 批量获取用户信息
func (wx *Wechat) GetUsers(openidList []string) (userInfo []STUserInfo, err error) {
	var openids []STOpenid
	for index, openid := range openidList {
		if (index+1)%99 == 0 {
			user, err := wx.GetUsers100(openids)
			if err != nil {
				return user, err
			}
			userInfo = append(userInfo, user...)
			openids = openids[0:0]
		}
		openids = append(openids, STOpenid{openid})
	}
	user, err := wx.GetUsers100(openids)
	if err != nil {
		return user, err
	}
	userInfo = append(userInfo, user...)
	return
}

//GetUsers100 GetUsers100
func (wx *Wechat) GetUsers100(openids []STOpenid) (userInfo []STUserInfo, err error) {
	if len(openids) > 100 {
		log.Println("数据太多")
		return nil, errors.New("数据太多")
	}
	type stList struct {
		UserList []STOpenid `json:"user_list"`
	}
	t := stList{openids}
	d, err := json.Marshal(t)

	req, err := http.NewRequest("POST", "https://api.weixin.qq.com/cgi-bin/user/info/batchget?access_token="+wx.AccessToken, bytes.NewReader(d))
	resBody, err := common.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return userInfo, err
	}

	type users struct {
		UserInfoList []STUserInfo `json:"user_info_list"`
	}
	var userList users
	tempData := strings.Replace(string(resBody), "\x1c", "", -1)
	tempData = strings.Replace(tempData, "\x0e", "", -1)
	err = json.Unmarshal([]byte(tempData), &userList)
	if err != nil {
		log.Println(err, string(resBody))
		//		return userInfo, err
	}
	return userList.UserInfoList, nil
}

//GetUserList 获取用户openid列表
func (wx *Wechat) GetUserList() []string {
	var openidList []string
	nextOpenid := ""
loop:
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken
	param["next_openid"] = nextOpenid
	req, err := http.NewRequest("GET", common.Param("https://api.weixin.qq.com/cgi-bin/user/get", param), nil)
	resBody, err := common.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return openidList
	}
	var userList STUserList
	err = json.Unmarshal(resBody, &userList)
	if err != nil || userList.Errcode != 0 {
		log.Println(string(resBody))
		return openidList
	}
	openidList = append(openidList, userList.Data.Openid...)
	if userList.Total > len(openidList) {
		nextOpenid = userList.NextOpenid
		goto loop
	}
	return openidList
}

//GetAllUserInfo 获取所有用户的所有信息
func (wx *Wechat) GetAllUserInfo() (userInfo []STUserInfo, err error) {
	return wx.GetUsers(wx.GetUserList())
}

//CreateTag 创建用户标签
func (wx *Wechat) CreateTag(name string) (tagid int, err error) {

	type tag struct {
		Tag STTag `json:"tag"`
	}
	t := tag{STTag{0, name}}
	d, _ := json.Marshal(t)
	log.Println(string(d))
	req, err := http.NewRequest("POST", "https://api.weixin.qq.com/cgi-bin/tags/create?access_token="+wx.AccessToken, bytes.NewReader(d))
	resBody, err := common.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	err = json.Unmarshal(resBody, &t)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return t.Tag.Tagid, nil
}

//GetTag 获取所有用户标签
func (wx *Wechat) GetTag() (data []STTag, err error) {
	type tags struct {
		Tags []STTag `json:"tags"`
	}
	req, err := http.NewRequest("GET", "https://api.weixin.qq.com/cgi-bin/tags/get?access_token="+wx.AccessToken, nil)
	resBody, err := common.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return data, err
	}
	var tmpData tags
	err = json.Unmarshal(resBody, &tmpData)
	if err != nil {
		log.Println(err)
		return tmpData.Tags, err
	}
	return tmpData.Tags, nil
}

//UpdateTag 更新用户标签名
func (wx *Wechat) UpdateTag(tagid int, name string) (err error) {
	type tags struct {
		Tag STTag `json:"tag"`
	}
	t := tags{STTag{tagid, name}}
	d, _ := json.Marshal(t)
	req, err := http.NewRequest("POST", "https://api.weixin.qq.com/cgi-bin/tags/update?access_token="+wx.AccessToken, bytes.NewReader(d))
	_, err = common.RequsetJSON(req, 0)
	return err
}

//DelTags 删除用户标签
func (wx *Wechat) DelTags(tagid int) (err error) {
	type stTags struct {
		Tag STTag `json:"tag"`
	}
	t := stTags{STTag{tagid, ""}}
	d, _ := json.Marshal(t)
	log.Println(string(d))
	req, err := http.NewRequest("POST", "https://api.weixin.qq.com/cgi-bin/tags/delete?access_token="+wx.AccessToken, bytes.NewReader(d))
	_, err = common.RequsetJSON(req, 0)
	return err
}

//GetUserTags 获取用户所在标签
func (wx *Wechat) GetUserTags(openid string) []int {

	t := STOpenid{openid}
	d, _ := json.Marshal(t)

	req, err := http.NewRequest("POST", "https://api.weixin.qq.com/cgi-bin/tags/getidlist?access_token="+wx.AccessToken, bytes.NewReader(d))
	resBody, err := common.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return nil
	}
	type stTagid struct {
		TagidList []int `json:"tagid_list"`
	}
	var tmpTag stTagid
	err = json.Unmarshal(resBody, &tmpTag)
	if err != nil {
		log.Println(err)
		return nil
	}
	return tmpTag.TagidList
}

//BatchTags 批量为用户打标签
func (wx *Wechat) BatchTags(openid []string, tagid int) int {
	type stOpenid struct {
		Openid []string `json:"openidList"`
		Tagid  int      `json:"tagid"`
	}
	t := stOpenid{openid, tagid}
	d, _ := json.Marshal(t)
	log.Println(string(d))
	req, err := http.NewRequest("POST", "https://api.weixin.qq.com/cgi-bin/tags/members/batchtagging?access_token="+wx.AccessToken, bytes.NewReader(d))
	_, err = common.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return 0
	}

	return 1
}

//UnBatchTags 批量取消用户标签
func (wx *Wechat) UnBatchTags(openid []string, tagid int) int {
	type stOpenid struct {
		Openid []string `json:"openidList"`
		Tagid  int      `json:"tagid"`
	}
	t := stOpenid{openid, tagid}
	d, _ := json.Marshal(t)
	req, err := http.NewRequest("POST", "https://api.weixin.qq.com/cgi-bin/tags/members/batchuntagging?access_token="+wx.AccessToken, bytes.NewReader(d))
	_, err = common.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return 0
	}
	return 1
}

//UpdateRemark 设置备注名
func (wx *Wechat) UpdateRemark(openid, remark string) int {
	type stRemark struct {
		Openid string `json:"openid"`
		Remark string `json:"remark"`
	}
	t := stRemark{openid, remark}
	d, _ := json.Marshal(t)
	req, err := http.NewRequest("POST", "https://api.weixin.qq.com/cgi-bin/user/info/updateremark?access_token="+wx.AccessToken, bytes.NewReader(d))
	_, err = common.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return 0
	}

	return 1
}
