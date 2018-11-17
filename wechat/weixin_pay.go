package wechat

import (
	"bytes"
	"encoding/xml"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/wei193/wechat"
)

//baseurl
const (
	URLPAYUNIFIEDORDER = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	URLPAYORDERQUERY   = "https://api.mch.weixin.qq.com/pay/orderquery"
	URLPAYCLOSEORDER   = "https://api.mch.weixin.qq.com/pay/closeorder"
	URLDOWNLOADBILL    = "https://api.mch.weixin.qq.com/pay/downloadbill"
)

//ReqHongbao 红包发送结构体
type ReqHongbao struct {
	XMLName     xml.Name `xml:"xml"`
	NonceStr    string   `xml:"nonce_str"`
	Sign        string   `xml:"sign"`
	MchBillno   string   `xml:"mch_billno"`
	Mchid       string   `xml:"mch_id"`
	Wxappid     string   `xml:"wxappid"`
	SendName    string   `xml:"send_name"`
	ReOpenid    string   `xml:"re_openid"`
	TotalAmount int      `xml:"total_amount"`
	TotalNum    int      `xml:"total_num"`
	Wishing     string   `xml:"wishing"`
	ClientIP    string   `xml:"client_ip"`
	ActName     string   `xml:"act_name"`
	Remark      string   `xml:"remark"`
}

//ResHongbao ResHongbao
type ResHongbao struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
	ResultCode string   `xml:"result_code"`
	ErrCode    string   `xml:"err_code"`
	ErrCodeSes string   `xml:"err_code_des"`
	MchBillno  string   `xml:"mch_billno"`
	Mchid      string   `xml:"mch_id"`
	Wxappid    string   `xml:"wxappid"`
	ReOpenid   string   `xml:"re_openid"`
	SendListid string   `xml:"send_listid"`
}

//ReqQueryOrder 订单查询结构体
type ReqQueryOrder struct {
	XMLName       xml.Name `xml:"xml"`
	Appid         string   `xml:"appid"`
	Mchid         string   `xml:"mch_id"`
	Transactionid string   `xml:"transaction_id,omitempty"`
	OutTradeNo    string   `xml:"out_trade_no,omitempty"`
	Noncestr      string   `xml:"nonce_str"`
	Sign          string   `xml:"sign"`
}

//ResQueryOrder 订单查询
type ResQueryOrder struct {
	ReturnCode         string `xml:"return_code,omitempty"`
	ReturnMsg          string `xml:"return_msg,omitempty"`
	Appid              string `xml:"appid,omitempty"`
	Mchid              string `xml:"mch_id,omitempty"`
	DeviceInfo         string `xml:"device_info,omitempty"`
	NonceStr           string `xml:"nonce_str,omitempty"`
	Sign               string `xml:"sign,omitempty"`
	SignType           string `xml:"sign_type,omitempty"`
	ResultCode         string `xml:"result_code,omitempty"`
	ErrCode            string `xml:"err_code,omitempty"`
	ErrCodeDes         string `xml:"err_code_des,omitempty"`
	Openid             string `xml:"openid,omitempty"`
	IsSubscribe        string `xml:"is_subscribe,omitempty"`
	TradeType          string `xml:"trade_type,omitempty"`
	BankType           string `xml:"bank_type,omitempty"`
	TotalFee           int    `xml:"total_fee,omitempty"`
	SettlementTotalFee int    `xml:"settlement_total_fee,omitempty"`
	FeeType            string `xml:"fee_type,omitempty"`
	CashFee            int    `xml:"cash_fee,omitempty"`
	CashFeeType        string `xml:"cash_fee_type,omitempty"`
	CouponFee          int    `xml:"coupon_fee,omitempty"`
	CouponCount        int    `xml:"coupon_count,omitempty"`
	Transactionid      string `xml:"transaction_id,omitempty"`
	OutTradeNo         string `xml:"out_trade_no,omitempty"`
	Attach             string `xml:"attach,omitempty"`
	TimeEnd            string `xml:"time_end,omitempty"`
	TradeState         string `xml:"trade_state,omitempty"`
	TradeStateDesc     string `xml:"trade_state_desc,omitempty"`
}

//ReqUnifiedOrder 下单
type ReqUnifiedOrder struct {
	XMLName        xml.Name `xml:"xml"`
	Appid          string   `xml:"appid"`
	Mchid          string   `xml:"mch_id"`
	DeviceInfo     string   `xml:"device_info,omitempty"`
	NonceStr       string   `xml:"nonce_str"`
	Sign           string   `xml:"sign"`
	Body           string   `xml:"body"`
	Detail         string   `xml:"detail,omitempty"`
	Attach         string   `xml:"attach,omitempty"`
	OutTradeNo     string   `xml:"out_trade_no"`
	FeeType        string   `xml:"fee_type,omitempty"`
	TotalFee       int      `xml:"total_fee"`
	SpbillCreateIP string   `xml:"spbill_create_ip"`
	TimeStart      string   `xml:"time_start,omitempty"`
	TimeExpire     string   `xml:"time_expire,omitempty"`
	GoodsTag       string   `xml:"goods_tag,omitempty"`
	NotifyURL      string   `xml:"notify_url"`
	TradeType      string   `xml:"trade_type"`
	Productid      string   `xml:"product_id,omitempty"`
	LimitPay       string   `xml:"limit_pay,omitempty"`
	Openid         string   `xml:"openid,omitempty"`
	SceneInfo      string   `xml:"scene_info,omitempty"`
}

//ResUnifiedOrder 下单返回
type ResUnifiedOrder struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg,omitempty"`
	Appid      string   `xml:"appid"`
	Mchid      string   `xml:"mch_id"`
	DeviceInfo string   `xml:"device_info,omitempty"`
	NonceStr   string   `xml:"nonce_str"`
	Sign       string   `xml:"sign"`
	ResultCode string   `xml:"result_code"`
	ErrCode    string   `xml:"err_code,omitempty"`
	ErrCodeDes string   `xml:"err_code_des,omitempty"`
	TradeType  string   `xml:"trade_type"`
	Prepayid   string   `xml:"prepay_id"`
	CodeURL    string   `xml:"code_url,omitempty"`
	MwebURL    string   `xml:"mweb_url,omitempty"`
}

//ResCloseOrder 关闭订单
type ResCloseOrder struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	Appid      string `xml:"appid"`
	Mchid      string `xml:"mch_id"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
}

//ReqDownloadBill 请求下载对账单
type ReqDownloadBill struct {
	XMLName  xml.Name `xml:"xml"`
	Appid    string   `xml:"appid"`
	Mchid    string   `xml:"mch_id"`
	BiilDate string   `xml:"bill_date"`
	BillType string   `xml:"bill_type"`
	Noncestr string   `xml:"nonce_str"`
	Sign     string   `xml:"sign"`
}

//TPaySign TPaySign
type TPaySign struct {
	AppID     string `xml:"appId" json:"appId"`
	Timestamp int64  `xml:"timeStamp" json:"timeStamp"`
	NonceStr  string `xml:"nonceStr" json:"nonceStr"`
	Package   string `xml:"package" json:"package"`
	SignType  string `xml:"signType" json:"signType"`
	PaySign   string `xml:"paySign" json:"paySign"`
}

//TAPPPaySign TAPPPaySign
type TAPPPaySign struct {
	AppID     string `xml:"appid" json:"appid"`
	Partnerid string `xml:"partnerid" json:"partnerid"`
	Prepayid  string `xml:"prepayid" json:"prepayid"`
	Package   string `xml:"package" json:"package"`
	NonceStr  string `xml:"noncestr" json:"noncestr"`
	Timestamp int64  `xml:"timestamp" json:"timestamp"`
	PaySign   string `xml:"sign" json:"sign"`
}

//UnifiedOrder 支付下单https://api.mch.weixin.qq.com/pay/unifiedorder
func (wx *Wechat) UnifiedOrder(order ReqUnifiedOrder) (data ResUnifiedOrder, err error) {
	order.Appid = wx.Appid
	order.Mchid = wx.mch.MchID
	order.Sign = XMLSignMd5(order, wx.mch.PayKey)
	d, _ := xml.MarshalIndent(order, "", "\t")
	// base.PAYLOG.Info("unifiedorder send ", string(d))
	req, err := http.NewRequest("POST", URLPAYUNIFIEDORDER, bytes.NewReader(d))
	resBody, err := requsetXML(req, -1)
	// base.PAYLOG.Info("unifiedorder recv ", string(resBody))
	if err != nil {
		return data, err
	}
	err = xml.Unmarshal(resBody, &data)
	if err != nil {
		return data, err
	}
	Sign := data.Sign
	data.Sign = ""
	if XMLSignMd5(data, wx.mch.PayKey) != Sign {
		return data, errors.New("签名错误")
	}
	return data, nil
}

//QueryOrder 查询订单https://api.mch.weixin.qq.com/pay/orderquery
func (wx *Wechat) QueryOrder(transactionid, outTradeNo string) (data ResQueryOrder, err error) {

	queryOrder := ReqQueryOrder{
		Appid:         wx.Appid,
		Mchid:         wx.mch.MchID,
		Transactionid: transactionid,
		OutTradeNo:    outTradeNo,
		Noncestr:      RandomStr(20, 3)}

	queryOrder.Sign = XMLSignMd5(queryOrder, wx.mch.PayKey)
	d, _ := xml.MarshalIndent(queryOrder, "", "\t")
	req, err := http.NewRequest("POST", URLPAYORDERQUERY, bytes.NewReader(d))
	resBody, err := requsetXML(req, -1)
	if err != nil {
		log.Println(err)
		return data, err
	}
	err = xml.Unmarshal(resBody, &data)
	if err != nil {
		return data, err
	}
	Sign := data.Sign
	data.Sign = ""
	if XMLSignMd5(data, wx.mch.PayKey) != Sign {
		return data, errors.New("签名错误")
	}
	return data, nil
}

//CloseOrder 关闭订单https://api.mch.weixin.qq.com/pay/closeorder
func (wx *Wechat) CloseOrder(outTradeNo string) (data ResCloseOrder, err error) {
	queryOrder := ReqQueryOrder{
		Appid:      wx.Appid,
		Mchid:      wx.mch.MchID,
		OutTradeNo: outTradeNo,
		Noncestr:   RandomStr(20, 3)}
	queryOrder.Sign = XMLSignMd5(queryOrder, wx.mch.PayKey)
	d, _ := xml.MarshalIndent(queryOrder, "", "\t")
	req, err := http.NewRequest("POST", URLPAYCLOSEORDER, bytes.NewReader(d))
	resBody, err := requsetXML(req, -1)
	if err != nil {
		log.Println(err)
		return data, err
	}
	err = xml.Unmarshal(resBody, &data)
	if err != nil {
		return data, err
	}
	Sign := data.Sign
	data.Sign = ""
	if XMLSignMd5(data, wx.mch.PayKey) != Sign {
		return data, errors.New("签名错误")
	}
	return data, nil
}

//Refund 申请退款https://api.mch.weixin.qq.com/secapi/pay/refund
func (wx *Wechat) Refund() {

}

//Downloadbill 下载对账单https://api.mch.weixin.qq.com/pay/downloadbill
func (wx *Wechat) Downloadbill(billDate string, billType string) (data string, err error) {
	queryBill := ReqDownloadBill{
		Appid:    wx.Appid,
		Mchid:    wx.mch.MchID,
		BiilDate: billDate,
		BillType: billType,
		Noncestr: RandomStr(20, 3)}
	queryBill.Sign = XMLSignMd5(queryBill, wx.mch.PayKey)
	d, _ := xml.MarshalIndent(queryBill, "", "\t")
	req, err := http.NewRequest("POST", URLDOWNLOADBILL, bytes.NewReader(d))
	resBody, err := requsetXML(req, -1, false)
	if err != nil {
		log.Println(err)
		return "", err
	}
	// log.Println(string(resBody))
	return string(resBody), nil
}

//SendHongbao 发送红包
func (wx *Wechat) SendHongbao(hb ReqHongbao) (resp ResHongbao, err error) {
	hb.Sign = XMLSignMd5(hb, wx.mch.PayKey)
	data, err := xml.MarshalIndent(&hb, "", " ")
	if err != nil {
		return resp, err
	}
	res, err := wx.httpsPost("https://api.mch.weixin.qq.com/mmpaymkttransfers/sendredpack", data, "text/xml")
	if err != nil {
		log.Println(res, err)
		return resp, err
	}
	resBody := make([]byte, 1024)
	res.Body.Read(resBody)
	err = xml.Unmarshal(resBody, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

//CreatePaySign 创建PaySign
func (wx *Wechat) CreatePaySign(prepayid string) (data TPaySign) {
	data = TPaySign{
		AppID:     wx.Appid,
		Timestamp: time.Now().Unix(),
		NonceStr:  wechat.RandomStr(20, 3),
		Package:   "prepay_id=" + prepayid,
		SignType:  "MD5",
	}
	data.PaySign = XMLSignMd5(data, wx.mch.PayKey)
	return
}

//CreateAPPPaySign 创建APPPaySign
func (wx *Wechat) CreateAPPPaySign(prepayid string) (data TAPPPaySign) {
	data = TAPPPaySign{
		AppID:     wx.Appid,
		Partnerid: wx.mch.MchID,
		Prepayid:  prepayid,
		Package:   "Sign=WXPay",
		NonceStr:  wechat.RandomStr(20, 3),
		Timestamp: time.Now().Unix(),
	}
	data.PaySign = XMLSignMd5(data, wx.mch.PayKey)
	return
}
