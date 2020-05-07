package wechat

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/wei193/component/base"
)

//base url
const (
	URLMASSSENDALL  = "https://api.weixin.qq.com/cgi-bin/message/mass/sendall"
	URLMASSSENDLIST = "https://api.weixin.qq.com/cgi-bin/message/mass/send"
	URLMASSDELETE   = "https://api.weixin.qq.com/cgi-bin/message/mass/delete"
	URLMASSPREVIEW  = "https://api.weixin.qq.com/cgi-bin/message/mass/preview"
	URLMASSSEND     = "https://api.weixin.qq.com/cgi-bin/message/custom/send"
)

//STAutoReply 自动回复表
type STAutoReply struct {
	AutoReplyid int        //自动回复ID
	ReqType     string     //请求类型
	ReqContent  string     //请求内容
	ReqEvent    string     //请求事件
	ReqEventkey string     //请求事件key
	ResType     string     //自动回复类型
	ResTitle    string     //自动回复标题
	ResContent  string     //自动回复内容
	ResMediaid  string     //自动回复媒体ID
	ResArticles []TArticle //自动回复文章
}

//STMsgRequest 请求参数
type STMsgRequest struct {
	XMLName              xml.Name `xml:"xml"`
	ToUserName           string
	FromUserName         string
	CreateTime           time.Duration
	MsgType              string
	Content              string
	Mediaid              string
	LocationX, LocationY float32
	Scale                int
	Label                string
	PicURL               string
	Msgid                int
	Event                string
	EventKey             string
	Ticket               string
	Title                string
	Description          string
}

//STMediaid 媒体ID
type STMediaid struct {
	Mediaid string `json:"media_id,omitempty"`
}

//STText 文字消息
type STText struct {
	Content string `json:"content,omitempty"`
}

//STVideo 视频类
type STVideo struct {
	Title       string `json:"title,omitempty"`
	Mediaid     string `json:"media_id,omitempty"`
	Description string `json:"description,omitempty"`
}

//STMusic 音乐类
type STMusic struct {
	Title        string //音乐标题
	Description  string //音乐描述
	MusicURL     string //音乐链接
	HQMusicURL   string //高质量音乐链接，WIFI环境优先使用该链接播放音乐
	ThumbMediaid string //缩略图的媒体id，通过素材管理接口上传多媒体文件，得到的id
}

//STArticles 文章类
type STArticles struct {
	Articles []TArticle `xml:"Articles,omitempty"`
	Mediaid  string     `xml:"MediaId,omitempty" json:"media_id,omitempty"`
}

//TArticle 文章类
type TArticle struct {
	XMLName     xml.Name `xml:"item"`
	Title       string
	Description string
	PicURL      string `xml:"PicUrl" json:"PicUrl"`
	URL         string `xml:"Url" json:"Url"`
}

//STMsgResponse 返回参数
type STMsgResponse struct {
	XMLName      xml.Name      `xml:"xml"`
	ToUserName   string        `xml:",omitempty"`
	FromUserName string        `xml:",omitempty"`
	CreateTime   time.Duration `xml:",omitempty"`
	MsgType      string        `xml:",omitempty"`
	Content      string        `xml:",omitempty"`
	Image        *STMediaid    `xml:",omitempty"`
	Voice        *STMediaid    `xml:",omitempty"`
	Music        *STMusic      `xml:",omitempty"`
	Video        *STVideo      `xml:",omitempty"`
	ArticleCount int           `xml:",omitempty"`
	Articles     *STArticles   `xml:",omitempty"`
}

//STFilter STFilter
type STFilter struct {
	IsToAll bool `json:"is_to_all"`
	Tagid   int  `json:"tag_id,omitempty"`
}

//ResMsg ResMsg
//type	媒体文件类型，分别有图片（image）、语音（voice）、视频（video）和缩略图（thumb），news，即图文消息
//errcode	错误码
//errmsg	错误信息
//msgid	消息发送任务的ID
//msg_data_id
type ResMsg struct {
	Type      string `json:"type"`
	Msgid     int64  `json:"msgid"`
	MsgDataid int64  `json:"msg_data_id"`
}

//CreateTextRes 创建自动回复文字类消息
func CreateTextRes(req *STMsgRequest, Content string) (res interface{}, err error) {
	resp := STMsgResponse{}
	resp.CreateTime = time.Duration(time.Now().Unix())
	resp.ToUserName = req.FromUserName
	resp.FromUserName = req.ToUserName
	resp.MsgType = Text
	resp.Content = Content
	return resp, nil
}

//CreateImageRes 创建自动回复图片类消息
func CreateImageRes(req *STMsgRequest, Mediaid string) (res interface{}, err error) {
	resp := STMsgResponse{}
	resp.CreateTime = time.Duration(time.Now().Unix())
	resp.ToUserName = req.FromUserName
	resp.FromUserName = req.ToUserName
	resp.MsgType = Image
	resp.Image.Mediaid = Mediaid
	return resp, nil
}

//CreateVoiceRes 创建自动回复语音类消息
func CreateVoiceRes(req *STMsgRequest, Mediaid string) (res interface{}, err error) {
	resp := STMsgResponse{}
	resp.CreateTime = time.Duration(time.Now().Unix())
	resp.ToUserName = req.FromUserName
	resp.FromUserName = req.ToUserName
	resp.MsgType = "voice"
	resp.Voice.Mediaid = Mediaid
	return resp, nil
}

//CreateMusicRes 创建自动回复音乐类消息
func CreateMusicRes(req *STMsgRequest, Title, Description, MusicURL, HQMusicURL, ThumbMediaid string) (res interface{}, err error) {
	resp := STMsgResponse{}
	resp.CreateTime = time.Duration(time.Now().Unix())
	resp.ToUserName = req.FromUserName
	resp.FromUserName = req.ToUserName
	resp.MsgType = Music
	resp.Music.Title = Title
	resp.Music.Description = Description
	resp.Music.MusicURL = MusicURL
	resp.Music.HQMusicURL = HQMusicURL
	resp.Music.ThumbMediaid = ThumbMediaid
	return resp, nil
}

//CreateVideoRes 创建自动回复视频类消息
func CreateVideoRes(req *STMsgRequest, Title, Description, Mediaid string) (res interface{}, err error) {
	resp := STMsgResponse{}
	resp.CreateTime = time.Duration(time.Now().Unix())
	resp.ToUserName = req.FromUserName
	resp.FromUserName = req.ToUserName
	resp.MsgType = "video"
	resp.Video.Mediaid = Mediaid
	resp.Video.Title = Title
	resp.Video.Description = Description
	return resp, nil
}

//CreateArticlesRes 创建自动回复文章类消息
func CreateArticlesRes(req *STMsgRequest, data []TArticle) (res interface{}, err error) {
	resp := STMsgResponse{}
	resp.CreateTime = time.Duration(time.Now().Unix())
	resp.ToUserName = req.FromUserName
	resp.FromUserName = req.ToUserName
	resp.MsgType = News
	resp.ArticleCount = len(data)
	for _, a := range data {
		resp.Articles.Articles = append(resp.Articles.Articles, a)
	}
	return resp, nil
}

//createResponse 创建被动回复消息
func createResponse(req *STMsgRequest, auto STAutoReply) (resp interface{}, err error) {
	switch auto.ResType {
	case Text:
		return CreateTextRes(req, auto.ResContent)
	case Music:
		str := strings.Split(auto.ResContent, "@@")
		if len(str) != 3 {
			return nil, nil
		}
		return CreateMusicRes(req, auto.ResTitle, str[0], str[1], str[2], auto.ResMediaid)
	case "voice":
		return CreateVoiceRes(req, auto.ResMediaid)
	case Image:
		return CreateImageRes(req, auto.ResMediaid)
	case News:
		return CreateArticlesRes(req, auto.ResArticles)
	case "video":
		return CreateVideoRes(req, auto.ResTitle, auto.ResContent, auto.ResMediaid)
	}
	return nil, nil
}

//CheckSignature 检查微信消息签名
func (wx *Wechat) CheckSignature(signature, timestamp, nonce string) bool {

	tmps := []string{wx.Token, timestamp, nonce}
	sort.Strings(tmps)
	tmpStr := strings.Join(tmps, "")
	t := sha1.New()
	io.WriteString(t, tmpStr)
	tmp := fmt.Sprintf("%x", t.Sum(nil))

	if tmp == signature {
		return true
	}
	return false
}

//SendAll https://api.weixin.qq.com/cgi-bin/message/mass/sendall?access_token=ACCESS_TOKEN
func (wx *Wechat) SendAll(Tagid int, Msgtype, Content string) (sendData ResMsg, err error) {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken

	var Filter STFilter
	if Tagid == 0 {
		Filter.IsToAll = true
	} else {
		Filter.IsToAll = false
		Filter.Tagid = Tagid
	}
	var data interface{}
	switch Msgtype {
	case "mpnews":
		type stTmp struct {
			Filter  STFilter  `json:"filter"`
			Mpnews  STMediaid `json:"mpnews"`
			Msgtype string    `json:"msgtype"`
		}
		t := stTmp{Filter, STMediaid{Content}, Msgtype}
		data = t
	case "text":
		type stTmp struct {
			Filter  STFilter `json:"filter"`
			Text    STText   `json:"text"`
			Msgtype string   `json:"msgtype"`
		}
		t := stTmp{Filter, STText{Content}, Msgtype}
		data = t
	case "image":
		type stTmp struct {
			Filter  STFilter  `json:"filter"`
			Image   STMediaid `json:"image"`
			Msgtype string    `json:"msgtype"`
		}
		t := stTmp{Filter, STMediaid{Content}, Msgtype}
		data = t
	case "voice":
		type stTmp struct {
			Filter  STFilter  `json:"filter"`
			Voice   STMediaid `json:"voice"`
			Msgtype string    `json:"msgtype"`
		}
		t := stTmp{Filter, STMediaid{Content}, Msgtype}
		data = t
	case "mpvideo":
		type stTmp struct {
			Filter  STFilter  `json:"filter"`
			Mpvideo STMediaid `json:"mpvideo"`
			Msgtype string    `json:"msgtype"`
		}
		t := stTmp{Filter, STMediaid{Content}, Msgtype}
		data = t
	default:
		return sendData, errors.New("暂不支持该类型")
	}

	d, _ := json.Marshal(data)
	log.Println(string(d))
	req, err := http.NewRequest("POST", base.Param("https://api.weixin.qq.com/cgi-bin/message/mass/sendall", param),
		bytes.NewReader(d))
	if err != nil {
		log.Println(err)
		return sendData, err
	}
	resp, err := base.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err, string(resp))
		return sendData, err
	}
	err = json.Unmarshal(resp, &sendData)
	if err != nil {
		log.Println(err)
		return sendData, err
	}
	return sendData, nil
}

//SendList https://api.weixin.qq.com/cgi-bin/message/mass/send?access_token=ACCESS_TOKEN
func (wx *Wechat) SendList(userList []string, Msgtype, Content string) (sendData ResMsg, err error) {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken

	var data interface{}
	switch Msgtype {
	case "mpnews":
		type stTmp struct {
			Touser  []string  `json:"touser"`
			Mpnews  STMediaid `json:"mpnews"`
			Msgtype string    `json:"msgtype"`
		}
		t := stTmp{userList, STMediaid{Content}, Msgtype}
		data = t
	case "text":
		type stTmp struct {
			Touser  []string `json:"touser"`
			Text    STText   `json:"text"`
			Msgtype string   `json:"msgtype"`
		}
		t := stTmp{userList, STText{Content}, Msgtype}
		data = t
	case "image":
		type stTmp struct {
			Touser  []string  `json:"touser"`
			Image   STMediaid `json:"image"`
			Msgtype string    `json:"msgtype"`
		}
		t := stTmp{userList, STMediaid{Content}, Msgtype}
		data = t
	case "voice":
		type stTmp struct {
			Touser  []string  `json:"touser"`
			Voice   STMediaid `json:"voice"`
			Msgtype string    `json:"msgtype"`
		}
		t := stTmp{userList, STMediaid{Content}, Msgtype}
		data = t
	case "mpvideo":
		type stTmp struct {
			Touser  []string  `json:"touser"`
			Mpvideo STMediaid `json:"mpvideo"`
			Msgtype string    `json:"msgtype"`
		}
		t := stTmp{userList, STMediaid{Content}, Msgtype}
		data = t
	default:
		return sendData, errors.New("暂不支持该类型")
	}

	d, _ := json.Marshal(data)
	log.Println(string(d))
	req, err := http.NewRequest("POST", base.Param("https://api.weixin.qq.com/cgi-bin/message/mass/send", param),
		bytes.NewReader(d))
	if err != nil {
		log.Println(err)
		return sendData, err
	}
	resp, err := base.RequsetJSON(req, 0)
	if err != nil {
		log.Println(err)
		return sendData, err
	}
	err = json.Unmarshal(resp, &sendData)
	if err != nil {
		log.Println(err)
		return sendData, err
	}
	return sendData, nil
}

//DeleteMsg https://api.weixin.qq.com/cgi-bin/message/mass/delete?access_token=ACCESS_TOKEN
func (wx *Wechat) DeleteMsg(msgid string) int {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken
	t := STMediaid{msgid}
	d, _ := json.Marshal(t)
	req, err := http.NewRequest("POST", base.Param("https://api.weixin.qq.com/cgi-bin/message/mass/delete", param),
		bytes.NewReader(d))
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

//PreviewMsg https://api.weixin.qq.com/cgi-bin/message/mass/preview?access_token=ACCESS_TOKEN
func (wx *Wechat) PreviewMsg(openid, wxname, Msgtype, Content string) int {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken

	var data interface{}
	switch Msgtype {
	case "mpnews":
		type stTmp struct {
			Touser   string    `json:"touser,omitempty"`
			Towxname string    `json:"towxname,omitempty"`
			Mpnews   STMediaid `json:"mpnews"`
			Msgtype  string    `json:"msgtype"`
		}
		t := stTmp{openid, wxname, STMediaid{Content}, Msgtype}
		data = t
	case "text":
		type stTmp struct {
			Touser   string `json:"touser,omitempty"`
			Towxname string `json:"towxname,omitempty"`
			Text     STText `json:"text"`
			Msgtype  string `json:"msgtype"`
		}
		t := stTmp{openid, wxname, STText{Content}, Msgtype}
		data = t
	case "image":
		type stTmp struct {
			Touser   string    `json:"touser,omitempty"`
			Towxname string    `json:"towxname,omitempty"`
			Image    STMediaid `json:"image"`
			Msgtype  string    `json:"msgtype"`
		}
		t := stTmp{openid, wxname, STMediaid{Content}, Msgtype}
		data = t
	case "voice":
		type stTmp struct {
			Touser   string    `json:"touser,omitempty"`
			Towxname string    `json:"towxname,omitempty"`
			Voice    STMediaid `json:"voice"`
			Msgtype  string    `json:"msgtype"`
		}
		t := stTmp{openid, wxname, STMediaid{Content}, Msgtype}
		data = t
	case "mpvideo":
		type stTmp struct {
			Touser   string    `json:"touser,omitempty"`
			Towxname string    `json:"towxname,omitempty"`
			Mpvideo  STMediaid `json:"mpvideo"`
			Msgtype  string    `json:"msgtype"`
		}
		t := stTmp{openid, wxname, STMediaid{Content}, Msgtype}
		data = t
	default:
		return 0
	}

	d, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", base.Param("https://api.weixin.qq.com/cgi-bin/message/mass/preview", param),
		bytes.NewReader(d))
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

//SendMsg https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=ACCESS_TOKEN
func (wx *Wechat) SendMsg(data interface{}) int {
	param := make(map[string]string)
	param["access_token"] = wx.AccessToken

	d, _ := json.Marshal(data)
	req, err := http.NewRequest("POST",
		base.Param("https://api.weixin.qq.com/cgi-bin/message/custom/send", param),
		bytes.NewReader(d))
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
