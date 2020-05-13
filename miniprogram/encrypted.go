package miniprogram

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
)

//Watermark 小程序水印
type Watermark struct {
	Appid     string `json:"appid"`
	Timestamp int64  `json:"timestamp"`
}

//PhoneNumber 获取手机号
type PhoneNumber struct {
	PhoneNumber     string    `json:"phoneNumber"`
	PurePhoneNumber string    `json:"purePhoneNumber"`
	CountryCode     string    `json:"countryCode"`
	Watermark       Watermark `json:"watermark"`
}

//UserInfo 用户信息
type UserInfo struct {
	OpenID    string `json:"openId"`
	NickName  string `json:"nickName"`
	AvatarURL string `json:"avatarUrl"`
	Gender    int    `json:"gender"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Language  string `json:"language"`
	UnionID   string `json:"unionId"`
}

//PKCS7Padding PKCS7Padding
func PKCS7Padding(ciphertext []byte) []byte {
	padding := aes.BlockSize - len(ciphertext)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//PKCS7UnPadding PKCS7UnPadding
func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

//DecryptData 解密数据
func DecryptData(sessionKey, encryptedData, iv string) (data []byte, err error) {
	key, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, err
	}
	plaintext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	eiv, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCDecrypter(block, eiv)
	mode.CryptBlocks(ciphertext, plaintext)

	ciphertext = PKCS7UnPadding(ciphertext)
	return ciphertext, nil
}

//DecryptPhoneNumer 解析手机号码
func DecryptPhoneNumer(sessionKey, encryptedData, iv string) (phone PhoneNumber, err error) {
	buf, err := DecryptData(sessionKey, encryptedData, iv)
	if err != nil {
		return phone, err
	}
	err = json.Unmarshal(buf, &phone)
	if err != nil {
		return phone, err
	}
	return phone, nil
}

//DecryptUserInfo 解析用户信息
func DecryptUserInfo(sessionKey, encryptedData, iv string) (user UserInfo, err error) {
	buf, err := DecryptData(sessionKey, encryptedData, iv)
	if err != nil {
		return user, err
	}
	err = json.Unmarshal(buf, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}
