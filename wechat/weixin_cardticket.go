// Copyright 2019 wei_193 Author. All Rights Reserved.
//
// 微信公众号卡券相关

package wechat

import (
	"bytes"
	"encoding/json"
	"net/http"
)

//卡券相关接口地址
const (
	URLCardCreate          = "https://api.weixin.qq.com/card/create"              //创建卡券
	URLCardPaycell         = "https://api.weixin.qq.com/card/paycell/set"         //设置买单接口
	URLCardSelfconsumecell = "https://api.weixin.qq.com/card/selfconsumecell/set" //设置自助核销
	URLCardQrcode          = "https://api.weixin.qq.com/card/qrcode/create"       //创建二维码
	URLCardBatchget        = "https://api.weixin.qq.com/card/batchget"            //批量获取卡券

)

//TACard 卡券创建
type TACard struct {
	CardType      string          `json:"card_type"`
	Groupon       *TGroupon       `json:"groupon,omitempty"`
	Cash          *TCash          `json:"cash,omitempty"`
	Discount      *TDiscount      `json:"discount,omitempty"`
	Gift          *TGift          `json:"gift,omitempty"`
	GeneralCoupon *TGeneralCoupon `json:"generalcoupon,omitempty"`
}

//TGroupon 团购券类型
type TGroupon struct {
	BaseInfo     TBaseInfo     `json:"base_info"`
	AdvancedInfo TAdvancedInfo `json:"advanced_info"`
	DealDetail   string        `json:"deal_detail"`
}

//TCash 代金券类型
type TCash struct {
	BaseInfo     TBaseInfo     `json:"base_info"`
	AdvancedInfo TAdvancedInfo `json:"advanced_info"`
	LeastCost    int           `json:"least_cost"`
	ReduceCost   int           `json:"reduce_cost"`
}

//TDiscount 折扣券
type TDiscount struct {
	BaseInfo     TBaseInfo     `json:"base_info"`
	AdvancedInfo TAdvancedInfo `json:"advanced_info"`
	Discount     int           `json:"discount"`
}

//TGift 兑换券类型
type TGift struct {
	BaseInfo     TBaseInfo     `json:"base_info"`
	AdvancedInfo TAdvancedInfo `json:"advanced_info"`
	Gift         string        `json:"gift"`
}

//TGeneralCoupon 优惠券类型
type TGeneralCoupon struct {
	BaseInfo      TBaseInfo     `json:"base_info"`
	AdvancedInfo  TAdvancedInfo `json:"advanced_info"`
	DefaultDetail string        `json:"default_detail"`
}

//TBaseInfo 卡券基础信息字段
type TBaseInfo struct {
	Logourl                   string    `json:"logo_url"`
	BrandName                 string    `json:"brand_name"`
	CodeType                  string    `json:"code_type"`
	Title                     string    `json:"title"`
	Color                     string    `json:"color"`
	Notice                    string    `json:"notice"`
	ServicePhone              string    `json:"service_phone,omitempty"`
	Description               string    `json:"description"`
	DateInfo                  TDateInfo `json:"date_info"`
	Sku                       TSku      `json:"sku"`
	UseLimit                  int       `json:"use_limit,omitempty"`
	GetLimit                  int       `json:"get_limit,omitempty"`
	UseCustomCode             bool      `json:"use_custom_code,omitempty"`
	BindOpenid                bool      `json:"ind_openid,omitempty"`
	CanShare                  bool      `json:"can_share,omitempty"`
	CanGiveFriend             bool      `json:"can_give_friend,omitempty"`
	LocationIDList            []int     `json:"location_id_list,omitempty"`
	CenterTitle               string    `json:"center_title,omitempty"`
	CenterSubTitle            string    `json:"center_sub_title,omitempty"`
	CustomAppBrandUserName    string    `json:"custom_app_brand_user_name,omitempty"`
	CustomAppBrandPass        string    `json:"custom_app_brand_pass,omitempty"`
	CenterURL                 string    `json:"center_url,omitempty"`
	CustomURLName             string    `json:"custom_url_name,omitempty"`
	CustomURL                 string    `json:"custom_url,omitempty"`
	CustomURLSubTitle         string    `json:"custom_url_sub_title,omitempty"`
	PromotionAppBrandUserName string    `json:"promotion_app_brand_user_name,omitempty"`
	PromotionAppBrandPass     string    `json:"promotion_app_brand_pass,omitempty"`
	PromotionURLName          string    `json:"promotion_url_name,omitempty"`
	PromotionURL              string    `json:"promotion_url,omitempty"`
	PromotionURLSubTitle      string    `json:"promotion_url_sub_title,omitempty"`
	// source                        string    `json:"promotion_url"`
}

//TAdvancedInfo 卡券高级信息
type TAdvancedInfo struct {
	UseCondition    *TUseCondition   `json:"use_condition,omitempty"`
	Abstract        *TAbstract       `json:"abstract,omitempty"`
	TextImageList   []TTextImageList `json:"text_image_list,omitempty"`
	TimeLimit       []TTimeLimit     `json:"time_limit,omitempty"`
	BusinessService []string         `json:"business_service,omitempty"`
}

//TDateInfo 使用日期，有效期的信息。
type TDateInfo struct {
	Type           string `json:"type"`
	BeginTimestamp int64  `json:"begin_timestamp"`
	EndTimestamp   int64  `json:"end_timestamp,omitempty"`
	FixedTerm      int    `json:"fixed_term,omitempty"`
	FixedBeginTerm int    `json:"fixed_begin_term,omitempty"`
}

//TSku 商品信息
type TSku struct {
	Quantity int `json:"quantity"`
}

//TUseCondition //使用门槛
type TUseCondition struct {
	AcceptCategory          string `json:"accept_category,omitempty"`
	RejectCategory          string `json:"reject_category,omitempty"`
	CanUseWithOtherDiscount bool   `json:"can_use_with_other_discount,omitempty"`
	LeastCost               int    `json:"least_cost,omitempty"`
	ObjectUseFor            string `json:"object_use_for,omitempty"`
}

//TAbstract 封面摘要简介
type TAbstract struct {
	Abstract    string   `json:"abstract,omitempty"`
	IconURLList []string `json:"icon_url_list,omitempty"`
}

//TTextImageList 图文列表
type TTextImageList struct {
	ImageURL string `json:"image_url,omitempty"`
	Text     string `json:"text,omitempty"`
}

//TTimeLimit 使用时段限制，包含以下字段
type TTimeLimit struct {
	Type        string `json:"type,omitempty"`
	BeginHour   int    `json:"begin_hour,omitempty"`
	EndHour     int    `json:"end_hour,omitempty"`
	BeginMinute int    `json:"begin_minute,omitempty"`
	EndMinute   int    `json:"end_minute,omitempty"`
}

//CardCreate 创建卡券
func (wx *Wechat) CardCreate(card TACard) (cardid string, err error) {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken

	d, _ := json.Marshal(card)
	req, err := http.NewRequest("POST", Param(URLCardCreate, param),
		bytes.NewReader(d))
	type st struct {
		CardID string `json:"card_id"`
	}
	var data st
	resBody, err := requsetJSON(req, 0)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		return "", err
	}
	return data.CardID, nil
}

//CardPaycell 设置买单接口
func (wx *Wechat) CardPaycell(cardid string, isopen bool) (err error) {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken

	type st struct {
		CardID string `json:"card_id"`
		IsOpen bool   `json:"is_open"`
	}
	data := st{
		CardID: cardid,
		IsOpen: isopen,
	}
	d, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", Param(URLCardPaycell, param),
		bytes.NewReader(d))

	_, err = requsetJSON(req, 0)
	if err != nil {
		return err
	}
	return nil
}

//CardSelfconsumecell  设置自助核销接口
func (wx *Wechat) CardSelfconsumecell(cardid string, isopen bool) (err error) {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken

	type st struct {
		CardID string `json:"card_id"`
		IsOpen bool   `json:"is_open"`
	}
	data := st{
		CardID: cardid,
		IsOpen: isopen,
	}
	d, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", Param(URLCardSelfconsumecell, param),
		bytes.NewReader(d))

	_, err = requsetJSON(req, 0)
	if err != nil {
		return err
	}
	return nil
}

//CardQrcode 创建二维码接口
func (wx *Wechat) CardQrcode() {

}
