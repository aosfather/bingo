package wx

import (
	"encoding/xml"
	"time"
)

const (
	Text     = "text"
	Location = "location"
	Image    = "image"
	Link     = "link"
	Event    = "event"
	Music    = "music"
	News     = "news"
	Voice    = "voice"
	Video    = "video"
)

type baseResponseMessage struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type WxRequest struct {
	Timestamp string `Field:"timestamp"`
	Nonce     string `Field:"nonce"`
	Signature string `Field:"signature"`
	Echostr   string `Field:"echostr"`
}

type WxMessage struct {
	WxRequest
	data WxMessageBody
}

func (this *WxMessage) GetData() interface{} {
	return &this.data
}
func (this *WxMessage) GetDataType() string {
	return "xml"
}

type msgBase struct {
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string //文本消息text、图片image、语音voice、视频为video、小视频为shortvideo、地理位置location、链接消息link

}

type WxMessageBody struct {
	XMLName xml.Name `xml:"xml"`
	msgBase
	Content  string //文本内容
	Event    string //当类型为Event时候有值，事件类型：subscribe(订阅)、unsubscribe(取消订阅)、自定义菜单事件CLICK、上报位置LOCATION、用户已关注时的事件推送SCAN、菜单跳转VIEW
	EventKey string //CLICK事件KEY值，与自定义菜单接口中KEY值对应
	//为通过二维码订阅的时候事件为subscribe事件KEY值，qrscene_为前缀，后面为二维码的参数值
	Ticket                 string  //二维码的ticket，可用来换取二维码图片
	Latitude               float32 //LOCATION地理位置纬度
	Longitude              float32 //LOCATION地理位置经度
	Precision              float32 //LOCATION地理位置精度
	Location_X, Location_Y float32 //地理位置维度,地理位置经度
	Scale                  int     //地图缩放大小
	Label                  string  //地理位置信息
	PicUrl                 string  //图片地址
	Format                 string  //语音格式，如amr，speex等
	Recognition            string  //语音识别结果，UTF8编码
	MediaId                string  //图片消息媒体id，可以调用多媒体文件下载接口拉取数据。
	ThumbMediaId           string  //视频消息缩略图的媒体id，可以调用多媒体文件下载接口拉取数据。
	Title                  string  //link消息标题
	Description            string  //link消息描述
	Url                    string  //link消息链接
	MsgId                  int     //消息id，64位整型
}

type responseMsg struct {
	XMLName xml.Name `xml:"xml"`
	msgBase
}

func (this *responseMsg) Init(wbody *WxMessageBody) {
	this.FromUserName = wbody.ToUserName
	this.ToUserName = wbody.FromUserName
	this.CreateTime = time.Duration(time.Now().Unix())
}

type WxResponse struct {
	msgBase
}

type WxTextMsg struct {
	responseMsg
	Content string
}

func (this *WxTextMsg) Init(w *WxMessageBody) {
	this.responseMsg.Init(w)
	this.MsgType = Text
}
func (this *WxTextMsg) SetBody(text string) {
	this.Content = text
}

type WxImageMsg struct {
	responseMsg
	Image string `xml:"Image>MediaId"`
}

func (this *WxImageMsg) Init(w *WxMessageBody) {
	this.responseMsg.Init(w)
	this.MsgType = Image
}
func (this *WxImageMsg) SetBody(image string) {
	this.Image = image

}

type WxVoiceMsg struct {
	responseMsg
	Voice string `xml:"Voice>MediaId"`
}

func (this *WxVoiceMsg) Init(w *WxMessageBody) {
	this.responseMsg.Init(w)
	this.MsgType = Voice
}

type WxVideoMsg struct {
	responseMsg
	Video video
}

func (this *WxVideoMsg) Init(w *WxMessageBody) {
	this.responseMsg.Init(w)
	this.MsgType = Video
}
func (this *WxVideoMsg) SetBody(video string) {
	this.Video.MediaId = video
}

type video struct {
	Title       string
	Description string
	MediaId     string
}

type WxMusicMsg struct {
	responseMsg
	Music music
}

func (this *WxMusicMsg) Init(w *WxMessageBody) {
	this.responseMsg.Init(w)
	this.MsgType = Music
}

func (this *WxMusicMsg) SetBody(title, desc, murl, hqurl, mediaid string) {
	this.Music = music{title, desc, murl, hqurl, mediaid}
}

type music struct {
	Title        string
	Description  string
	MusicUrl     string
	HQMusicUrl   string //高质量音乐链接，WIFI环境优先使用该链接播放音乐
	ThumbMediaId string //缩略图的媒体id，通过素材管理中的接口上传多媒体文件，得到的id
}

type WxArticleMsg struct {
	responseMsg
	ArticleCount int            `xml:",omitempty"`
	Articles     []*articleItem `xml:"Articles>item,omitempty"`
}

func (this *WxArticleMsg) Init(w *WxMessageBody) {
	this.responseMsg.Init(w)
	this.MsgType = News
}

func (this *WxArticleMsg) AddArticle(title, desc, picurl, url string) {
	item := articleItem{title, desc, picurl, url}
	if this.Articles == nil {

		this.Articles = []*articleItem{&item}
	} else {
		this.Articles = append(this.Articles, &item)
	}
	this.ArticleCount++
}

type articleItem struct {
	Title       string
	Description string
	PicUrl      string
	Url         string
}
