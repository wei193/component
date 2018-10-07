package component

import (
	"encoding/json"
	"log"
	"time"

	"github.com/wei193/component/wechat"
)

//Authorizer Authorizer
type Authorizer struct {
	Component              *Component
	AuthorizerAppid        string
	AuthorizerAccessToken  string
	AccessTokenExpires     int64
	AuthorizerRefreshToken string
}

//JAuthorizer 授权信息
type JAuthorizer struct {
	AuthorizerInfo    JAuthorizerInfo    `json:"authorizer_info"`
	AuthorizationInfo JAuthorizationInfo `json:"authorization_info"`
}

//JAuthorizerInfo 授权信息
type JAuthorizerInfo struct {
	NickName        string        `json:"nick_name"`
	HeadImg         string        `json:"head_img"`
	ServiceTypeInfo JTypeInfo     `json:"service_type_info"`
	VerifyTypeInfo  JTypeInfo     `json:"verify_type_info"`
	UserName        string        `json:"user_name"`
	PrincipalName   string        `json:"principal_name"`
	BusinessInfo    JBusinessInfo `json:"business_info"`
	Alias           string        `json:"alias"`
	QrcodeURL       string        `json:"qrcode_url"`
}

//JBusinessInfo JBusinessInfo
type JBusinessInfo struct {
	OpenStore int `json:"open_store"`
	OpenScan  int `json:"open_scan"`
	OpenPay   int `json:"open_pay"`
	OpenCard  int `json:"open_card"`
	OpenShake int `json:"open_shake"`
}

//JAuthorizationInfo 授权信息
type JAuthorizationInfo struct {
	AuthorizerAppid        string      `json:"authorizer_appid"`
	AuthorizationAppid     string      `json:"authorization_appid"`
	AuthorizerAccessToken  string      `json:"authorizer_access_token"`
	ExpiresIn              int         `json:"expires_in"`
	AuthorizerRefreshToken string      `json:"authorizer_refresh_token"`
	FuncInfo               []JFuncInfo `json:"func_info"`
}

//JTypeInfo 类型信息
type JTypeInfo struct {
	ID int `json:"id"`
}

//JFuncInfo 功能列表
type JFuncInfo struct {
	FuncscopeCategory JFuncscopeCategory `json:"funcscope_category"`
}

// JFuncscopeCategory  JFuncscopeCategory
type JFuncscopeCategory struct {
	ID int `json:"id"`
}

//JAuthorizerAccessToken JAuthorizerAccessToken
type JAuthorizerAccessToken struct {
	AuthorizerAccessToken  string `json:"authorizer_access_token"`
	ExpiresIn              int    `json:"expires_in"`
	AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
}

//JUserAccessToken 用户AccessToken
type JUserAccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}

//NewAuthorizer 新建Authorizer
func (c *Component) NewAuthorizer(appid, accesstoken string, tokenexpires int64, refreshtoken string) (authorizer *Authorizer, err error) {
	authorizer = &Authorizer{
		Component:              c,
		AuthorizerAppid:        appid,
		AuthorizerAccessToken:  accesstoken,
		AccessTokenExpires:     tokenexpires,
		AuthorizerRefreshToken: refreshtoken,
	}
	return
}

//GetAuthorizerInfo 获取授权详细信息
func (a *Authorizer) GetAuthorizerInfo() (authorizer *JAuthorizer, err error) {
	type st struct {
		ComponentAppid  string `json:"component_appid"`
		AuthorizerAppid string `json:"authorizer_appid"`
	}
	d := st{
		ComponentAppid:  a.Component.ComponentAppid,
		AuthorizerAppid: a.AuthorizerAppid,
	}
	req, err := createRequset("https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_info?component_access_token="+a.Component.ComponentAccessToken,
		"POST", nil, d)
	if err != nil {
		return nil, err
	}
	res, err := requsetJosn(req)
	if err != nil {
		return nil, err
	}
	log.Println(string(res))

	var auth JAuthorizer
	err = json.Unmarshal(res, &auth)
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

//GetAuthorizerAccessToken 获取Authorizer AccessToken
func (a *Authorizer) GetAuthorizerAccessToken() (token *JAuthorizerAccessToken, err error) {
	type st struct {
		ComponentAppid         string `json:"component_appid"`
		AuthorizerAppid        string `json:"authorizer_appid"`
		AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
	}
	d := st{
		ComponentAppid:         a.Component.ComponentAppid,
		AuthorizerAppid:        a.AuthorizerAppid,
		AuthorizerRefreshToken: a.AuthorizerRefreshToken,
	}
	req, err := createRequset("https://api.weixin.qq.com/cgi-bin/component/api_authorizer_token?component_access_token="+a.Component.ComponentAccessToken,
		"POST", nil, d)
	if err != nil {
		return nil, err
	}
	res, err := requsetJosn(req)
	if err != nil {
		return nil, err
	}
	log.Println(string(res))
	token = new(JAuthorizerAccessToken)
	err = json.Unmarshal(res, token)
	if err != nil {
		return nil, err
	}
	a.AuthorizerAccessToken = token.AuthorizerAccessToken
	a.AuthorizerRefreshToken = token.AuthorizerRefreshToken
	a.AccessTokenExpires = time.Now().Unix() + int64(token.ExpiresIn)
	return token, nil
}

//CodeToAccessToken 通过code换取access_token
func (a *Authorizer) CodeToAccessToken(code string) (token *JUserAccessToken, err error) {
	param := make(map[string]string)
	param["appid"] = a.AuthorizerAppid
	param["code"] = code
	param["grant_type"] = "authorization_code"
	param["component_appid"] = a.Component.ComponentAppid
	param["component_access_token"] = a.Component.ComponentAccessToken

	req, err := createRequset("https://api.weixin.qq.com/sns/oauth2/component/access_token",
		"GET", param, nil)
	if err != nil {
		return nil, err
	}
	res, err := requsetJosn(req)
	if err != nil {
		return nil, err
	}
	log.Println(string(res))
	token = new(JUserAccessToken)
	err = json.Unmarshal(res, token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

//GetWecaht 获取微信方法
func (a *Authorizer) GetWecaht() (w *wechat.Wechat) {
	return &wechat.Wechat{
		Appid: a.AuthorizerAppid,
		// Appsecret       :,
		// Token           :,
		// Encodingaeskey  :,
		AccessToken:        a.AuthorizerAccessToken,
		AccessTokenExpires: a.AccessTokenExpires,
		// JsapiTicket     :,
		// JsapiTokenTime  :,
	}
}
