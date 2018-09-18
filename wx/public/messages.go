package public

import "encoding/xml"

//-----------微信应用---------------------//
//基本消息格式
type xmlBaseMessage struct {
	FromUserName string
	ToUserName   string
	CreateTime   int64
	MsgType      string
	MsgId        int64
}

//文本消息
type WXxmlTextMessage struct {
	XMLName xml.Name `xml:"xml"`
	xmlBaseMessage
	Content string
}

//图片消息
type WXxmlImageMessage struct {
	XMLName xml.Name `xml:"xml"`
	xmlBaseMessage
	PicUrl  string //图片路径
	MediaId string //媒体id
}

//语音消息
type WXxmlVoiceMessage struct {
	XMLName xml.Name `xml:"xml"`
	xmlBaseMessage
	MediaId string //媒体id
	Format  string //文件格式
}

//视频/小视频消息
type WXxmlVideoMessage struct {
	XMLName xml.Name `xml:"xml"`
	xmlBaseMessage
	MediaId      string //视频媒体文件id，可以调用获取媒体文件接口拉取数据，仅三天内有效
	ThumbMediaId string //视频消息缩略图的媒体id，可以调用获取媒体文件接口拉取数据，仅三天内有效
}

//地址消息
type WXxmlLocationMessage struct {
	XMLName xml.Name `xml:"xml"`
	xmlBaseMessage
	X     float64 `xml:"Location_X"` //地理位置纬度
	Y     float64 `xml:"Location_Y"` //地理位置经度
	Scale int64   //地图缩放大小
	Label string  //地理位置信息
}

//连接消息
type WXxmlLinkMessage struct {
	XMLName xml.Name `xml:"xml"`
	xmlBaseMessage
	Title       string //标题
	Description string //描述
	Url         string //封面缩略图的url
}

type IReplyMsg interface {
	GetBaseInfo() *xmlReplyBaseMessage
}

//被动回复消息格式
type xmlReplyBaseMessage struct {
	FromUserName string
	ToUserName   string
	CreateTime   int64
	MsgType      string
}

func (this *xmlReplyBaseMessage) GetBaseInfo() *xmlReplyBaseMessage {
	return this
}

//文本回复消息
type WXxmlReplyTextMessage struct {
	XMLName xml.Name `xml:"xml"`
	xmlReplyBaseMessage
	Content CDATA
}

type CDATA string

func (c CDATA) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct {
		string `xml:",cdata"`
	}{string(c)}, start)
}

type WXxmlReplyImageMessage struct {
	XMLName xml.Name `xml:"xml"`
	xmlReplyBaseMessage
	Content xmlReplyImageContent
}

//语音
type WXxmlReplyVoiceMessage struct {
	XMLName xml.Name `xml:"xml"`
	xmlReplyBaseMessage
	Content xmlReplyVoiceContent
}

//视频
type WXxmlReplyVideoMessage struct {
	XMLName xml.Name `xml:"xml"`
	xmlReplyBaseMessage
	Content xmlReplyVideoContent
}

//音乐
type WXxmlReplyMusicMessage struct {
	XMLName xml.Name `xml:"xml"`
	xmlReplyBaseMessage
	Content xmlReplyMusicContent
}

//图文
type WXxmlReplyArticleMessage struct {
	XMLName xml.Name `xml:"xml"`
	xmlReplyBaseMessage
	Content xmlReplyArticlesContent
}

type xmlReplyImageContent struct {
	XMLName xml.Name `xml:"Image"`
	MediaId string   //媒体id
}

type xmlReplyVoiceContent struct {
	XMLName xml.Name `xml:"Voice"`
	MediaId string   //媒体id
}

type xmlReplyVideoContent struct {
	XMLName     xml.Name `xml:"Video"`
	MediaId     string   //媒体id
	Title       string   //标题
	Description string   //描述
}

type xmlReplyMusicContent struct {
	XMLName      xml.Name `xml:"Music"`
	Title        string   //标题
	Description  string   //描述
	MusicURL     string   //音乐链接
	HQMusicUrl   string   //高质量音乐链接，WIFI环境优先使用该链接播放音乐
	ThumbMediaId string   //必填，缩略图的媒体id，通过素材管理中的接口上传多媒体文件，得到的id
}

type xmlReplyArticlesContent struct {
	XMLName xml.Name `xml:"Articles"`
	Items   []xmlReplyArticlesItem
}

type xmlReplyArticlesItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   //标题
	Description string   //描述
	PicUrl      string   //图片链接
	Url         string   //跳转
}

//事件
type XMLEvent struct {
	XMLName      xml.Name `xml:"xml"`
	FromUserName string
	ToUserName   string
	CreateTime   int64
	MsgType      string //事件类型，subscribe(订阅)、unsubscribe(取消订阅),SCAN(扫描二维码)
	//如果用户已经关注事件类型，SCAN，未关注时候类型为subscribe
	//	事件类型，跳转 VIEW，菜单点击 CLICK
	Event string

	EventKey string //事件KEY值，qrscene_为前缀，后面为二维码的参数值；VIEW和CLICK时候，与自定义菜单接口中KEY值对应|事件KEY值，设置的跳转URL

	Ticket string //二维码的ticket，可用来换取二维码图片
	//	事件类型，LOCATION
	Latitude  float64 //地理位置纬度
	Longitude float64 //地理位置经度
	Precision float64 //地理位置精度
}
