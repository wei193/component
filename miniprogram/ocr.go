package miniprogram

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/wei193/component/common"
)

//BankcardInfo 银行卡信息
type BankcardInfo struct {
	Number  string `json:"number"`
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

//Bankcard 银行卡 OCR 识别
func (mini *MiniProgram) Bankcard(media string) (info BankcardInfo, err error) {
	param := make(map[string]string)
	param["access_token"] = mini.AccessToken

	file, err := os.Open(media)
	if err != nil {
		return info, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("img", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	req, err := http.NewRequest("POST",
		common.Param("https://api.weixin.qq.com/cv/ocr/bankcard", param),
		body)

	resBody, err := common.Requset(req)
	if err != nil {
		return info, err
	}
	err = json.Unmarshal(resBody, &info)
	if err != nil {
		return info, err
	}
	if info.Errcode != 0 {
		return info, errors.New(info.Errmsg)
	}
	return info, nil
}

//BankcardByURL 银行卡 OCR 识别
func (mini *MiniProgram) BankcardByURL(url string) (info BankcardInfo, err error) {
	param := make(map[string]string)
	param["access_token"] = mini.AccessToken
	param["img_url"] = url

	req, err := http.NewRequest("POST",
		common.Param("https://api.weixin.qq.com/cv/ocr/bankcard", param),
		nil)

	resBody, err := common.Requset(req)
	if err != nil {
		return info, err
	}
	err = json.Unmarshal(resBody, &info)
	if err != nil {
		return info, err
	}
	if info.Errcode != 0 {
		return info, errors.New(info.Errmsg)
	}
	return info, nil
}
