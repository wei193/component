package wechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/wei193/component/common"
)

//base url
const (
	URLMediaUpload           = "https://api.weixin.qq.com/cgi-bin/media/upload"
	URLMediaGet              = "https://api.weixin.qq.com/cgi-bin/media/get"
	URLMediaAddNews          = "https://api.weixin.qq.com/cgi-bin/material/add_news"
	URLMediaUpdateNews       = "https://api.weixin.qq.com/cgi-bin/material/update_news"
	URLMediaUploadImg        = "https://api.weixin.qq.com/cgi-bin/media/uploadimg"
	URLMediaAddMaterial      = "https://api.weixin.qq.com/cgi-bin/material/add_material"
	URLMediaDelMaterial      = "https://api.weixin.qq.com/cgi-bin/material/del_material"
	URLMediaGetMaterial      = "https://api.weixin.qq.com/cgi-bin/material/get_material"
	URLMediaBatchgetMaterial = "https://api.weixin.qq.com/cgi-bin/material/batchget_material"
)

//TNews 图文消息
type TNews struct {
	Title            string `json:"title"`
	ThumbMediaid     string `json:"thumb_media_id"`
	Author           string `json:"author"`
	Digest           string `json:"digest"`
	ShowCoverPic     int    `json:"show_cover_pic"`
	Content          string `json:"content"`
	URL              string `json:"url"`
	ContentSourceURL string `json:"content_source_url"`
}

//STNews 图文信息返回
type STNews struct {
	TNews
	URL      string `json:"url"`
	ThumbURL string `json:"thumb_url"`
}

//ReqMedia ReqMedia
type ReqMedia struct {
	Mediaid string `josn:"media_id"`
	URL     string `json:"url"`
}

//ResMedia 媒体返回
type ResMedia struct {
	Type      string `json:"type"`
	Mediaid   string `json:"media_id"`
	CreatedAt int64  `json:"created_at"`
	URL       string `json:"url"`
}

//TMaterialList TMaterialList
type TMaterialList struct {
	TotalCount int         `json:"total_count"`
	ItemCount  int         `json:"item_count"`
	Item       []TMaterial `json:"item"`
	Content    []TNews     `json:"content"`
}

//TMaterial TMaterial
type TMaterial struct {
	Mediaid    string `json:"media_id"`
	Name       string `json:"name,omitempty"`
	UpdateTime int64  `json:"update_time,omitempty"`
	URL        string `json:"url,omitempty"`
}

// TNewsItem TNewsItem
type TNewsItem struct {
	Mediaid    string       `json:"media_id"`
	Content    TNewsContent `json:"content"`
	UpdateTime int64        `json:"update_time"`
}

//TNewsContent TNewsContent
type TNewsContent struct {
	NewsItem []TNews `json:"news_item"`
}

//ResNews ResNews
type ResNews struct {
	TotalCount int         `json:"total_count"`
	ItemCount  int         `json:"item_count"`
	Item       []TNewsItem `json:"item"`
}

//AddTempMaterial 增加临时资源
func (wx *Wechat) AddTempMaterial(mediaType, filepath string) (data ReqMedia, err error) {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken
	param["type"] = mediaType
	req, err := newfileUploadRequest(common.Param(URLMediaUpload, param), nil,
		"media", filepath)
	if err != nil {
		return
	}
	resBody, err := common.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return data, err
	}
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		log.Println(err)
		return data, err
	}
	return data, nil
}

//GetTempMaterial 获取临时的媒体资源
func (wx *Wechat) GetTempMaterial(basepath, mediaid string) (string, error) {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken
	param["media_id"] = mediaid

	req, err := http.NewRequest("GET", common.Param(URLMediaGet, param), nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	fileName := ""
	Disposition := resp.Header.Get("Content-Disposition")
	if Disposition != "" {
		if strings.Index(Disposition, `attachment; filename="`) == 0 {
			fileName = Disposition[len(`attachment; filename="`) : len(Disposition)-1]
			log.Println(fileName)
		}
	}
	if fileName == "" {
		return "", errors.New("not a file")
	}
	fileName = path.Join(basepath, fileName)
	file, err := os.Create(fileName)
	if err != nil {
		log.Println(err)
		return fileName, err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Println(err)
		return fileName, err
	}
	return fileName, err
}

//AddNews 增加图文消息
func (wx *Wechat) AddNews(news []TNews) (data ReqMedia, err error) {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken

	type articles struct {
		Articles []TNews `json:"articles"`
	}
	t := articles{news}
	d, _ := json.Marshal(t)
	req, err := http.NewRequest("POST", common.Param(URLMediaAddNews, param),
		bytes.NewReader(d))
	if err != nil {
		log.Println(err)
		return data, err
	}
	resBody, err := common.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return data, err
	}
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		log.Println(err)
		return data, err
	}
	return data, nil
	//	return resBody, err
}

//UpdateNews 修改图文消息
func (wx *Wechat) UpdateNews(Mediaid string, Index int, news TNews) int {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken

	type articles struct {
		Mediaid  string `json:"media_id"`
		Index    int    `json:"index"`
		Articles TNews  `json:"articles"`
	}
	t := articles{Mediaid, Index, news}
	d, _ := json.Marshal(t)
	req, err := http.NewRequest("POST", common.Param(URLMediaUpdateNews, param),
		bytes.NewReader(d))
	if err != nil {
		return 0
	}
	_, err = common.RequsetJSON(req, 0)
	if err != nil {
		return 0
	}
	return 1
}

//UploadImg 增加图文消息图片
func (wx *Wechat) UploadImg(filepath string) (data ReqMedia, err error) {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken

	req, err := newfileUploadRequest(common.Param(URLMediaUploadImg, param),
		nil, "media", filepath)
	if err != nil {
		return
	}
	resBody, err := common.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return data, err
	}
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		log.Println(err)
		return data, err
	}
	return data, nil
}

//AddMaterial 增加资源
func (wx *Wechat) AddMaterial(mediaType, filepath string) (data ReqMedia, err error) {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken
	param["type"] = mediaType

	req, err := newfileUploadRequest(common.Param(URLMediaAddMaterial, param),
		nil, "media", filepath)
	if err != nil {
		return
	}
	resBody, err := common.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return data, err
	}
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		log.Println(err)
		return data, err
	}
	return data, nil
}

//DelMaterial 删除资源
func (wx *Wechat) DelMaterial(mediaid string) int {
	type stTmp struct {
		Mediaid string `json:"media_id"`
	}
	t := stTmp{mediaid}
	d, _ := json.Marshal(t)
	req, err := http.NewRequest("POST", URLMediaDelMaterial+"?access_token="+
		wx.AccessToken, bytes.NewReader(d))
	if err == nil {
		return 1
	}
	_, err = common.RequsetJSON(req, 0)
	if err == nil {
		return 1
	}
	return 0
}

//GetMaterial 获取除文章和视频类型外的媒体资源
func (wx *Wechat) GetMaterial(basepath, mediaid string) (string, error) {
	type stTmp struct {
		Mediaid string `json:"media_id"`
	}
	t := stTmp{mediaid}
	d, _ := json.Marshal(t)
	req, err := http.NewRequest("POST", URLMediaGetMaterial+"?access_token="+
		wx.AccessToken, bytes.NewReader(d))
	if err == nil {
		return "", err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	fileName := ""
	Disposition := resp.Header.Get("Content-Disposition")
	if Disposition != "" {
		if strings.Index(Disposition, `attachment; filename="`) == 0 {
			fileName = Disposition[len(`attachment; filename="`) : len(Disposition)-1]
			log.Println(fileName)
		}
	}
	if fileName == "" {
		return "", errors.New("not a file")
	}
	fileName = path.Join(basepath, fileName)
	file, err := os.Create(fileName)
	if err != nil {
		log.Println(err)
		return fileName, err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Println(err)
		return fileName, err
	}
	return fileName, err
}

//GetMaterialsNews 获取图文资源
func (wx *Wechat) GetMaterialsNews(Type string, offset, count int) (data ResNews, err error) {
	type stTmp struct {
		Type   string `json:"type"`
		Offset int    `json:"offset"`
		Count  int    `json:"count"`
	}
	t := stTmp{Type, offset, count}
	d, _ := json.Marshal(t)
	req, err := http.NewRequest("POST", URLMediaBatchgetMaterial+"?access_token="+
		wx.AccessToken, bytes.NewReader(d))
	if err == nil {
		return data, err
	}
	resBody, err := common.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return data, err
	}
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		log.Println(err)
		return data, err
	}
	return data, err
}

//GetMaterials 获取资源
func (wx *Wechat) GetMaterials(Type string, offset, count int) (data TMaterialList, err error) {
	type stTmp struct {
		Type   string `json:"type"`
		Offset int    `json:"offset"`
		Count  int    `json:"count"`
	}
	t := stTmp{Type, offset, count}
	d, _ := json.Marshal(t)
	req, err := http.NewRequest("POST", URLMediaBatchgetMaterial+"?access_token="+
		wx.AccessToken, bytes.NewReader(d))
	if err == nil {
		return data, err
	}
	resBody, err := common.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return data, err
	}
	//	log.Println(string(resBody))
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		log.Println(err)
		return data, err
	}
	return data, err
}

//GetAllMaterials 获取除文章类型外的资源列表
func (wx *Wechat) GetAllMaterials(Type string) (materials []TMaterial, err error) {
	offset := 0
	for {
		data, err := wx.GetMaterials(Type, offset, 20)
		if err != nil {
			return materials, err
		}
		offset += data.ItemCount
		materials = append(materials, data.Item...)
		if offset >= data.TotalCount {
			break
		}
	}
	return materials, nil
}

//GetAllMaterialsNews 获取文章类型资源列表
func (wx *Wechat) GetAllMaterialsNews() (materials []TNewsItem, err error) {
	offset := 0
	for {
		data, err := wx.GetMaterialsNews("news", offset, offset+20)
		if err != nil {
			return materials, err
		}
		offset += data.ItemCount
		materials = append(materials, data.Item...)
		if offset >= data.TotalCount {
			break
		}
	}
	return materials, nil
}

//微信文件上传
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req, err
}
