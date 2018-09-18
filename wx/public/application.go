package public

import (
	"encoding/xml"
	"fmt"
	"github.com/aosfather/bingo/mvc"
	"github.com/aosfather/bingo/utils"
	"github.com/aosfather/bingo/wx"
	"time"
)

//公众号

//微信的验证请求
type wxValidateRequest struct {
	Signature string `Field:"signature"`
	Timestamp string `Field:"timestamp"`
	Nonce     string `Field:"nonce"`
	Echostr   string `Field:"echostr"`
}

//微信发送过来的消息
type wxInputMsg struct {
	Signature   string `Field:"msg_signature"`
	Timestamp   string `Field:"timestamp"`
	Nonce       string `Field:"nonce"`
	EncryptType string `Field:"encrypt_type"`
	data        wx.WxEncryptInputMessage
}

func (this *wxInputMsg) GetInput() wx.WxEncryptInputMessage {
	return this.data
}
func (this *wxInputMsg) GetData() interface{} {
	return &this.data
}
func (this *wxInputMsg) GetDataType() string {
	return "xml"
}

//消息处理接口，用于实现应用自身的逻辑
type WXProcessor interface {
	OnEvent(event XMLEvent) IReplyMsg                    //事件响应
	OnMessage(msgtype string, msg interface{}) IReplyMsg //消息响应
}

type WXPublicApplication struct {
	Token  string `Value:"wx.public.token"`
	AppId  string `Value:"wx.public.appid"`
	AESKey string `Value:"wx.public.aeskey"`
	mvc.SimpleController
	logger    utils.Log
	encryted  *wx.BaseEncrypt //加解密
	Processor WXProcessor     `Inject:"processor"` //处理器
}

func (this *WXPublicApplication) Init() {
	this.logger = this.GetBeanFactory().GetLog("wx_public")
	this.encryted = &wx.BaseEncrypt{}
	this.encryted.Init(this.Token, this.AppId, this.AESKey)
}

func (this *WXPublicApplication) GetUrl() string {
	return "/wx/public_msg"
}

func (this *WXPublicApplication) GetParameType(method string) interface{} {
	if method == "GET" {
		return &wxValidateRequest{}
	} else {
		return &wxInputMsg{}
	}
}

func (this *WXPublicApplication) Get(c mvc.Context, p interface{}) (interface{}, mvc.BingoError) {
	if q, ok := p.(*wxValidateRequest); ok {
		this.logger.Info("wx validate %v", q)
		ret := this.encryted.VerifyURL(q.Signature, q.Timestamp, q.Nonce, q.Echostr)
		if ret {
			return q.Echostr, nil
		} else {
			this.logger.Info("wx validate failed！")

		}

	}

	return "hello", nil

}

//正常的访问消息处理
func (this *WXPublicApplication) Post(c mvc.Context, p interface{}) (interface{}, mvc.BingoError) {
	if msg, ok := p.(*wxInputMsg); ok {
		this.logger.Debug("msg:%s", msg)
		ret, result := this.encryted.DecryptInputMsg(msg.Signature, msg.Timestamp, msg.Nonce, msg.GetInput())
		this.logger.Debug("msg result:%i,%s", ret, result)
		if ret == 0 {
			var replymsg interface{}
			rmsg := xmlBaseMessage{}
			//消息处理
			if this.Processor != nil {
				//解析result，压入对象

				msgdata := []byte(result)
				xml.Unmarshal(msgdata, &rmsg)
				//根据msgtype类型构造对应的消息结构
				var realmsg interface{}
				switch rmsg.MsgType {
				case wx.MSGTYPE_TEXT:
					realmsg = &WXxmlTextMessage{}
				case wx.MSGTYPE_IMAGE:
					realmsg = &WXxmlImageMessage{}
				case wx.MSGTYPE_VOICE:
					realmsg = &WXxmlVoiceMessage{}
				case wx.MSGTYPE_VIDEO, wx.MSGTYPE_SHORT_VIDEO:
					realmsg = &WXxmlVideoMessage{}
				case wx.MSGTYPE_LOCATION:
					realmsg = &WXxmlLocationMessage{}
				case wx.MSGTYPE_LINK:
					realmsg = &WXxmlLocationMessage{}
				case wx.MSGTYPE_EVENT: //事件处理

					event := XMLEvent{}
					xml.Unmarshal([]byte(result), &event)
					replymsg = this.Processor.OnEvent(event)
				}

				if realmsg != nil {
					//重新解析消息
					xml.Unmarshal(msgdata, realmsg)
					this.logger.Debug("msg struct %s", realmsg)
					replymsg = this.Processor.OnMessage(rmsg.MsgType, realmsg)
				}

			}

			//如果给的返回消息不为空回复微信
			if replymsg != nil {

				//输出xml格式，加密返回
				if replyBaseMsg, ok := replymsg.(IReplyMsg); ok {
					this.logger.Debug("输出：%s", replymsg)
					replyBaseMsg.GetBaseInfo().CreateTime = time.Now().Unix()
					enmsg, _ := xml.Marshal(replymsg)
					this.logger.Debug("reply:%s", enmsg)
					_, result := this.encryted.EncryptMsg(string(enmsg), "wxcorpxingyun", fmt.Sprintf("%d", replyBaseMsg.GetBaseInfo().CreateTime))
					return result, nil
				}
				return "", nil
			} else {
				return "success", nil
			}
		}

	}

	return "hi", nil

}
