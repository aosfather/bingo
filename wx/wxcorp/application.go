package wxcorp

import (
	"fmt"
	"strings"

	"encoding/json"

	"github.com/aosfather/bingo/mvc"
)

type WxAuthCorpApi interface {
	WxOauthGetUser(app string, corp string, agent int, p WxRedirectParamter) WxUserDetail
	WxInitContact(corp string)
}

type wxCorpAppContext struct {
	id  int
	app string

	data map[string]*wxAppcontext //corp   agent context

}

func (this *wxCorpAppContext) Init(id int, title string) {
	this.id = id
	this.app = title
	this.data = make(map[string]*wxAppcontext)

}

type wxCorpApplicationHub struct {
	mvc.SimpleController
	suit     *WxCorpSuite
	contexts map[string]*wxCorpAppContext
	notifys  map[string]*WxApprovalApi //申请同步api
}

func (this *wxCorpApplicationHub) Init(suit *WxCorpSuite) {
	this.suit = suit
	this.contexts = make(map[string]*wxCorpAppContext)
	this.notifys = make(map[string]*WxApprovalApi)
}

//回复消息
func (this *wxCorpApplicationHub) ReplyMessage() {

}

//发送消息
func (this *wxCorpApplicationHub) SendTextMessage(app, corp string, group int, whos []string, content string) {
	msg := wxTextMsg{}
	msg.MsgType = "text"
	msg.Text.Content = content
	context := this.getAppcontext(app, corp, 1)
	msg.AgentId = context.AuthAgentId
	var target string = ""
	if whos != nil {
		target = strings.Join(whos, "|")

	}
	switch group {
	case MSG_TO_ALL:
		msg.ToUser = "@all"
	case MSG_TO_USER:
		msg.ToUser = target
	case MSG_TO_PARTY:
		msg.ToPart = target
	case MSG_TO_TAG:
		msg.ToTag = target
	}

	this.sendWxMessage(context, msg)

}

//和微信发送消息
func (this *wxCorpApplicationHub) sendWxMessage(context *wxAppcontext, msg interface{}) wxMsgResult {
	targetUrl := fmt.Sprintf(SEND_MSG_API, context.auth.AccessToken)
	result := wxMsgResult{}
	err := PostToWx(targetUrl, msg, &result)
	if err != nil {
		result.Errcode = 99999 //自定义的错误码
		result.ErrMsg = err.Error()
	}
	return result
}

func (this *wxCorpApplicationHub) getAppcontext(app, corp string, agent int) *wxAppcontext {
	appcontext := this.contexts[app]
	if appcontext == nil { //初始化app的上下文
		appcontext = &wxCorpAppContext{1, app, make(map[string]*wxAppcontext)}
		this.contexts[app] = appcontext
	}

	context := appcontext.data[corp]
	if context == nil { //初始化corp的上下文
		auth := this.suit.GetAuthCorpContext(corp)
		if auth == nil { //如果系统未配置完全，则返回空对象
			return nil
		}
		context = &wxAppcontext{auth, agent}
		appcontext.data[corp] = context

	}

	return context
}

//处理微信的跳转返回
func (this *wxCorpApplicationHub) WxOauthGetUser(app string, corp string, agent int, p WxRedirectParamter) WxUserDetail {
	//	appcontext := this.contexts[app]
	//	if appcontext == nil { //初始化app的上下文
	//		appcontext = &wxCorpAppContext{1, app, make(map[string]*wxAppcontext)}
	//		this.contexts[app] = appcontext
	//	}

	//	context := appcontext.data[corp]
	//	if context == nil { //初始化corp的上下文
	//		auth := this.suit.GetAuthCorpContext(corp)
	//		if auth == nil { //如果系统未配置完全，则返回空对象
	//			return WxUserDetail{}
	//		}
	//		context = &wxAppcontext{auth, agent}
	//		appcontext.data[corp] = context

	//	}
	context := this.getAppcontext(app, corp, agent)
	if context == nil {
		return WxUserDetail{}
	}

	return this.getUserDetailByCode(context, p.Code)

}

func (this *wxCorpApplicationHub) WxInitContact(corp string) {
	this.suit.initCorpContact(corp)
}

func (this *wxCorpApplicationHub) getUserDetailByCode(context *wxAppcontext, code string) WxUserDetail {
	//通过code取获取userinfo
	info := this.getUserInfo(context, code)
	fmt.Println(info)
	//通过ticket获取用户detail信息
	detail := this.getUserDetail(context, info.Ticket)
	return detail
}

//获取用户带ticket信息
func (this *wxCorpApplicationHub) getUserInfo(context *wxAppcontext, code string) wxUserInfo {
	url := fmt.Sprintf(USER_INFO_API, context.auth.AccessToken, code)
	fmt.Println("get userinfo:" + url)
	content, err := HTTPGet(url)
	info := wxUserInfo{}
	if err == nil {
		fmt.Println(string(content))
		json.Unmarshal(content, &info)
	}
	return info
}

//获取用户的详细信息
func (this *wxCorpApplicationHub) getUserDetail(context *wxAppcontext, ticket string) WxUserDetail {
	query := wxUserDetailQuery{ticket}
	url := fmt.Sprintf(USER_DETAIL_API, context.auth.AccessToken)
	fmt.Println("get userdetail:" + url)
	fmt.Println(query)
	result := WxUserDetail{}
	PostToWx(url, query, &result)
	return result
}
