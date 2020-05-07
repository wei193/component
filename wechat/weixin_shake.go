package wechat

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"strconv"

	"github.com/wei193/component/base"
)

//STShakePage STShakePage
type STShakePage struct {
	Pageid      int    `json:"page_id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Pageurl     string `json:"page_url,omitempty"`
	Comment     string `json:"comment,omitempty"`
	Iconurl     string `json:"icon_url,omitempty"`
}

//ReqLotteryinfo ReqLotteryinfo
type ReqLotteryinfo struct {
	Title        string `json:"title"`
	Desc         string `json:"desc"`
	Onoff        int    `json:"onoff"`
	BeginTime    int64  `json:"begin_time"`
	ExpireTime   int64  `json:"expire_time"`
	SponsorAppid string `json:"sponsor_appid"`
	Total        int64  `json:"total"`
	Jumpurl      string `json:"jump_url"`
	Key          string `json:"key"`
}

//ResLotteryinfo ResLotteryinfo
type ResLotteryinfo struct {
	Errcode   int    `json:"errcode"`
	Lotteryid string `json:"lottery_id"`
	Pageid    int    `json:"page_id"`
}

//ReqHbpreorder ReqHbpreorder
type ReqHbpreorder struct {
	XMLName     xml.Name `xml:"xml"`
	NonceStr    string   `xml:"nonce_str"`
	Sign        string   `xml:"sign"`
	MchBillno   string   `xml:"mch_billno"`
	Mchid       string   `xml:"mch_id"`
	Wxappid     string   `xml:"wxappid"`
	SendName    string   `xml:"send_name"`
	HbType      string   `xml:"hb_type"`
	TotalAmount int      `xml:"total_amount"`
	TotalNum    int      `xml:"total_num"`
	AmtType     string   `xml:"amt_type"`
	Wishing     string   `xml:"wishing"`
	ActName     string   `xml:"act_name"`
	Remark      string   `xml:"remark"`
	AuthMchid   string   `xml:"auth_mchid"`
	AuthAppid   string   `xml:"auth_appid"`
	RiskCntl    string   `xml:"risk_cntl"`
}

//ResHbpreorder ResHbpreorder
type ResHbpreorder struct {
	XMLName    xml.Name `xml:"xml"`
	SpTicket   string   `xml:"sp_ticket"`
	Detailid   string   `xml:"detail_id"`
	SendTime   string   `xml:"send_time"`
	ReturnCode string   `xml:"return_code"`
	ResultCode string   `xml:"result_code"`
}

//STticket STticket
type STticket struct {
	Ticket string `json:"ticket"`
}

//ReqPrizebucket ReqPrizebucket
type ReqPrizebucket struct {
	Lotteryid     string     `json:"lottery_id"`
	Mchid         string     `json:"mchid"`
	SponsorAppid  string     `json:"sponsor_appid"`
	PrizeInfoList []STticket `json:"prize_info_list"`
}

//ResPrizebucket ResPrizebucket
type ResPrizebucket struct {
	Errcode                  int        `json:"errcode"`
	RepeatTicketList         []STticket `json:"repeat_ticket_list"`
	ExpireTicketList         []STticket `json:"expire_ticket_list"`
	InvalidAmountTicketList  []STticket `json:"invalid_amount_ticket_list"`
	SuccessNum               int        `json:"success_num"`
	WrongAuthmchidTicketList []STticket `json:"wrong_authmchid_ticket_list"`
	InvalidTicketList        []STticket `json:"invalid_ticket_list"`
}

//ShakeDevice ShakeDevice
type ShakeDevice struct {
	Deviceid int    `json:"device_id"`
	UUID     string `json:"uuid"`
	Major    int    `json:"major"`
	Minor    int    `json:"minor"`
}

//ResShakeDevice ResShakeDevice
type ResShakeDevice struct {
	Applyid      int    `json:"apply_id"`
	ApplyTime    int64  `json:"apply_time"`
	AuditStatus  int    `json:"audit_status"`
	AuditComment string `json:"audit_comment"`
	AuditTime    int64  `json:"audit_time"`
}

//AddShakeDevice 申请设备ID
func (wx *Wechat) AddShakeDevice(quantity int, applyReason, comment string, poiid int) (res ResShakeDevice) {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken
	type stTmp struct {
		Quantity    int    `json:"quantity"`
		ApplyReason string `json:"apply_reason"`
		Comment     string `json:"comment,omitempty"`
		Poiid       int    `json:"poi_id,omitempty"`
	}
	t := stTmp{quantity, applyReason, comment, poiid}
	d, _ := json.Marshal(t)
	req, err := http.NewRequest("POST", base.Param("https://api.weixin.qq.com/shakearound/device/applyid", param),
		bytes.NewReader(d))
	resBody, err := base.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return res
	}

	type stRes struct {
		Data ResShakeDevice `json:"data"`
	}
	var data stRes
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		log.Println(err)
		return res
	}
	return data.Data
}

//ChkeckShakeDevice 查询设备ID申请审核状态
func (wx *Wechat) ChkeckShakeDevice(applyid int) (res ResShakeDevice) {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken
	type stTmp struct {
		Applyid int `json:"apply_id"`
	}
	t := stTmp{applyid}
	d, _ := json.Marshal(t)
	req, err := http.NewRequest("POST", base.Param("https://api.weixin.qq.com/shakearound/device/applystatus", param),
		bytes.NewReader(d))
	resBody, err := base.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return res
	}

	type stRes struct {
		Data ResShakeDevice `json:"data"`
	}
	var data stRes
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		log.Println(err)
		return res
	}
	return data.Data
}

//AddShakePage 添加一个页面
func (wx *Wechat) AddShakePage(title, description, pageurl, comment, iconurl string) int {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken

	t := STShakePage{0, title, description, pageurl, comment, iconurl}
	d, _ := json.Marshal(t)
	req, err := http.NewRequest("POST", base.Param("https://api.weixin.qq.com/shakearound/page/add", param),
		bytes.NewReader(d))
	resBody, err := base.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return 0
	}
	var data STShakePage
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		log.Println(err)
		return 0
	}
	return data.Pageid
}

//UpdateShakePage 编辑一个页面
func (wx *Wechat) UpdateShakePage(pageid int, title, description, pageurl, comment, iconurl string) int {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken

	t := STShakePage{pageid, title, description, pageurl, comment, iconurl}
	d, _ := json.Marshal(t)
	req, err := http.NewRequest("POST", base.Param("https://api.weixin.qq.com/shakearound/page/update", param),
		bytes.NewReader(d))
	resBody, err := base.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return 0
	}
	var data STShakePage
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		log.Println(err)
		return 0
	}
	return data.Pageid
}

//SearchShakePage 查询页面
func (wx *Wechat) SearchShakePage(typeid int, pageids []int, begin, count int) []STShakePage {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken

	type stTmp struct {
		Type    int   `json:"type"`
		Pageids []int `json:"page_ids,omitempty"`
		Begin   int   `json:"begin,omitempty"`
		Count   int   `json:"count,omitempty"`
	}

	t := stTmp{typeid, pageids, begin, count}
	d, _ := json.Marshal(t)
	req, err := http.NewRequest("POST", base.Param("https://api.weixin.qq.com/shakearound/page/search", param),
		bytes.NewReader(d))
	resBody, err := base.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return nil
	}
	type stPage struct {
		Pages []STShakePage `josn:"pages"`
	}
	type stData struct {
		Data stPage `josn:"data"`
	}
	var data stData
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		log.Println(err)
		return nil
	}
	return data.Data.Pages
}

//DeleteShakePage 删除一个页面
func (wx *Wechat) DeleteShakePage(pageid int) int {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken

	t := STShakePage{pageid, "", "", "", "", ""}
	d, _ := json.Marshal(t)
	req, err := http.NewRequest("POST", base.Param("https://api.weixin.qq.com/shakearound/page/delete", param),
		bytes.NewReader(d))
	_, err = base.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return 0
	}
	return 1
}

//UploadShakeImage 上传摇一摇图片素材
func (wx *Wechat) UploadShakeImage(mediaType, filepath string) string {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken
	param["type"] = mediaType
	req, err := newfileUploadRequest(base.Param("https://api.weixin.qq.com/shakearound/material/add", param), nil,
		"media", filepath)
	if err != nil {
		return ""
	}
	resBody, err := base.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return ""
	}
	type stPic struct {
		Picurl string `json:"pic_url"`
	}
	type stRes struct {
		Data stPic `json:"data"`
	}
	var data stRes
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		log.Println(err)
		return ""
	}
	return data.Data.Picurl
}

//SetLotterySwitch 设置红包活动抽奖开关
func (wx *Wechat) SetLotterySwitch(lotteryid string, onoff int) int {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken
	param["lottery_id"] = lotteryid
	param["onoff"] = strconv.Itoa(onoff)

	req, err := http.NewRequest("POST", base.Param("https://api.weixin.qq.com/shakearound/lottery/setlotteryswitch", param),
		nil)
	if err != nil {
		log.Println(err)
		return 0
	}
	_, err = base.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return 0
	}
	return 1

}

//applyid 申请设备ID
func (wx *Wechat) applyid(quantity int, applyreason, comment, poiid string) int {
	type stReq struct {
		Quantity    int    `json:"quantity"`
		Applyreason string `json:"apply_reason"`
		Comment     string `json:"comment"`
		Poiid       string `json:"poi_id"`
	}
	tmpStr := stReq{quantity, applyreason, comment, poiid}
	data, _ := json.Marshal(tmpStr)
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken

	req, err := http.NewRequest("POST", base.Param("https://api.weixin.qq.com/shakearound/device/applyid", param),
		bytes.NewReader(data))
	if err != nil {
		log.Println(err)
		return 0
	}
	_, err = base.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return 0
	}
	return 1

}

//查询设备ID申请审核状态
func applystatus(wx *Wechat, applyid int) int {
	//	resBody, err := post("https://api.weixin.qq.com/shakearound/device/applystatus?access_token="+wx.Access_token, []byte(`{"apply_id": `+strconv.Itoa(apply_id)+`}`), "application/json")

	//	if err != nil {
	//		log.Println(err)
	//	}
	//	type m_Errcode struct {
	//		Errcode       int    `json:"errcode"`
	//		Errmsg        string `json:"errmsg"`
	//		apply_time    int64
	//		audit_comment string
	//		audit_status  int
	//		audit_time    int64
	//	}
	//	var t m_Errcode
	//	json.Unmarshal(resBody, &t)
	//	err = json.Unmarshal(resBody, &t)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	if t.Errcode == 0 {
	//		return 1
	//	}
	return 0
}

//配置设备与页面的关联关系https://api.weixin.qq.com/shakearound/device/bindpage?access_token=ACCESS_TOKEN
func bindpage(wx *Wechat, deviceid int, pageids []int, Append int) int {
	//	type stReq struct {
	//		Device_identifier shakea_device `json:"device_identifier"`
	//		Page_ids          []int         `json:"page_ids"`
	//		Bind              int           `json:"bind"`
	//		Append            int           `json:"append"`
	//	}
	//	tmpStr := stReq{shakea_device{Device_id: device_id}, page_ids, 1, Append}
	//	d, _ := json.Marshal(tmpStr)
	//	resBody, err := post("https://api.weixin.qq.com/shakearound/device/bindpage?access_token="+wx.Access_token,
	//		d, "application/json")
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	type m_Errcode struct {
	//		Errcode int    `json:"errcode"`
	//		Errmsg  string `json:"errmsg"`
	//	}
	//	var t m_Errcode
	//	err = json.Unmarshal(resBody, &t)
	//	if err != nil {
	//		log.Println(err, ":", resBody)
	//		return 0
	//	}
	//	if t.Errcode == 0 {
	//		return 1
	//	}
	return 0
}

//查询设备列表https://api.weixin.qq.com/shakearound/device/search?access_token=ACCESS_TOKEN
func search(wx *Wechat) {
}

//编辑页面信息https://api.weixin.qq.com/shakearound/page/update?access_token=ACCESS_TOKEN
func updatePage(wx *Wechat, pageid int, title, description, pageurl, comment, iconurl string) int {
	//	type stReq struct {
	//		Page_id     int    `json:"page_id"`
	//		Title       string `json:"title"`
	//		Description string `json:"description"`
	//		Page_url    string `json:"page_url"`
	//		Comment     string `json:"comment"`
	//		Icon_url    string `json:"icon_url"`
	//	}
	//	tmpStr := stReq{page_id, title, description, page_url, comment, icon_url}
	//	data, _ := json.Marshal(tmpStr)
	//	resBody, err := post("https://api.weixin.qq.com/shakearound/page/update?access_token="+wx.Access_token,
	//		data, "application/json")
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	type m_Errcode struct {
	//		Errcode int    `json:"errcode"`
	//		Errmsg  string `json:"errmsg"`
	//	}
	//	var t m_Errcode
	//	err = json.Unmarshal(resBody, &t)
	//	if err != nil {
	//		log.Println(err, ":", resBody)
	//		return 0
	//	}
	//	if t.Errcode == 0 {
	//		return 1
	//	}
	return 0
}

//删除页面https://api.weixin.qq.com/shakearound/page/delete?access_token=ACCESS_TOKEN
func deletePage(wx *Wechat, pageid int) int {
	//	type stReq struct {
	//		Page_id int `json:"page_id"`
	//	}
	//	tmpStr := stReq{page_id}
	//	data, _ := json.Marshal(tmpStr)
	//	resBody, err := post("https://api.weixin.qq.com/shakearound/page/delete?access_token="+wx.Access_token,
	//		data, "application/json")
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	type m_Errcode struct {
	//		Errcode int    `json:"errcode"`
	//		Errmsg  string `json:"errmsg"`
	//	}
	//	var t m_Errcode
	//	err = json.Unmarshal(resBody, &t)
	//	if err != nil {
	//		log.Println(err, ":", resBody)
	//		return 0
	//	}
	//	if t.Errcode == 0 {
	//		return 1
	//	}
	return 0
}

func delTicket(list []STticket, ticket []string) []string {
	for _, v := range list {
		for i, t := range ticket {
			if t == v.Ticket {
				if i+1 != len(ticket) {
					ticket = append(ticket[:i], ticket[i+1:]...)
				} else {
					ticket = ticket[:i]
				}
			}
		}
	}
	return ticket
}
