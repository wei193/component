package wechat

import "time"

const (
	Appid          = "wxc161f1d6315f3530"
	Appsecret      = "206acd0c5422624301be1b06e1d407e8"
	Token          = "bjwonder"
	Encodingaeskey = "g6Lyi5o09BUJFdJvsUvFKX6ZEXkSC71gwWQqEKgIRDc"
	AccessToken    = ""
)

var testWx *Wechat

func init() {
	testWx = New(Appid, Appsecret, Token, Encodingaeskey)
	if testWx.AccessToken == "" {
		err := testWx.GetAccessToken()
		if err != nil {
			panic(err)
		}
	} else {
		testWx.AccessToken = ""
		testWx.AccessTokenExpires = time.Now().Add(2 * time.Hour).Unix()
	}

}
