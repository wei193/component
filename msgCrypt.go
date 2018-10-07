package component

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

//MsgDecrypt 消息解密
func (c *Component) MsgDecrypt(data string) (result string, err error) {
	buf, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	d, err := AesCBCDecrypt(buf, c.AESKey)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

//MsgEncrypt 消息加密
func (c *Component) MsgEncrypt(data string) (result string, err error) {
	buf, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	d, err := AesCBCEncrypt(buf, c.AESKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(d), nil
}

//AesCBCDecrypt 解密
func AesCBCDecrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockModel := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	plantText := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plantText, ciphertext)
	plantText = PKCS7UnPadding(plantText, block.BlockSize())
	return plantText, nil
}

//PKCS7UnPadding PKCS7删除
func PKCS7UnPadding(plantText []byte, blockSize int) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

//AesCBCEncrypt 加密
func AesCBCEncrypt(plantText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plantText = PKCS7Padding(plantText, block.BlockSize())

	blockModel := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])

	ciphertext := make([]byte, len(plantText))

	blockModel.CryptBlocks(ciphertext, plantText)
	return ciphertext, nil
}

//PKCS7Padding PKCS7填充
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
