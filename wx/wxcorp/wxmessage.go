package wxcorp

const (
	SEND_MSG_API = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"
	MSG_TO_ALL   = 0
	MSG_TO_USER  = 1
	MSG_TO_PARTY = 2
	MSG_TO_TAG   = 3
)

//微信目的地址
type ToAddress struct {
	ToAll   bool     //是否给所有人发送
	ToUser  []string //接收人
	ToParty []string //接收部分
	ToTag   []string //接收的tag
}

//微信返回的发送消息结果
type wxMsgResult struct {
	Errcode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	Invaliduser  string `json:"invaliduser"`
	Invalidparty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
}

type wxBaseMsg struct {
	ToUser  string `json:"touser"`
	ToPart  string `json:"toparty"`
	ToTag   string `json:"totag"`
	MsgType string `json:"msgtype"`
	AgentId int    `json:"agentid"`
	Safe    int    `json:"safe"`
}

//文本消息
type wxTextMsg struct {
	wxBaseMsg
	Text wxTextContent `json:"text"`
}

type wxTextContent struct {
	Content string `json:"content"` //content字段可以支持换行、以及A标签，即可打开自定义的网页（可参考以上示例代码）(注意：换行符请用转义过的\n
}

type wxMediaContent struct {
	MediaId string `json:"media_id"`
}
type wxVideoContent struct {
	wxMediaContent
	Title       string `json:"title"`
	Description string `json:"description"`
}

type wxTextCardContent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	ButtonLabel string `json:"btntxt"`
}

type wxNewsContent struct {
	wxTextCardContent
	PicUrl string `json:"picurl"`
}
type wxMpNewsContent struct {
	Title       string `json:"title"`
	MediaId     string `json:"thumb_media_id"` //图文消息缩略图的media_id, 可以在上传多媒体文件接口中获得.
	Author      string `json:"author"`
	ContentUrl  string `json:"content_source_url"`
	Content     string `json:"content"` //图文消息的内容，支持html标签，不超过666 K个字节
	Description string `json:"digest"`  //图文消息的描述，不超过512个字节，超过会自动截断
}

type wxArticleContent struct {
	Articles []wxNewsContent `json:"articles"`
}

type wxMpArticleContent struct {
	Articles []wxMpNewsContent `json:"articles"`
}

//图片消息
type wxImageMsg struct {
	wxBaseMsg
	Image wxMediaContent `json:"image"`
}

//声音消息
type wxVoiceMsg struct {
	wxBaseMsg
	Voice wxMediaContent `json:"voice"`
}

//视频消息
type wxVideoMsg struct {
	wxBaseMsg
	Video wxVideoContent `json:"video"`
}

//文件消息
type wxFileMsg struct {
	wxBaseMsg
	File wxMediaContent `json:"file"`
}

//文本卡片消息
type wxTextCardMsg struct {
	wxBaseMsg
	TextCard wxTextCardContent `json:"textcard"`
}

//图文消息
type wxNewsMsg struct {
	wxBaseMsg
	News wxArticleContent `json:"news"`
}

//微信存储的图文消息
type wxMpNewsMsg struct {
	wxBaseMsg
	News wxArticleContent `json:"mpnews"`
}
