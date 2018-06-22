package wxcorp

import (
	"fmt"
)

/**
微信支付
需要在userid和openid之间转换
*/

const (
	OPEN_USER_API = "https://qyapi.weixin.qq.com/cgi-bin/user/convert_to_userid?access_token=%s"
	USER_OPEN_API = " https://qyapi.weixin.qq.com/cgi-bin/user/convert_to_openid?access_token=%s"
)

type wxOpenIdQuery struct {
	UserId string `json:"userid"`
}
type wxOpenIdExtQuery struct {
	wxOpenIdQuery
	AgentId int `json:"agentid"`
}

type wxOpenIdResult struct {
	baseMessage
	OpenId string `json:"openid"`

	AppId string `json:"appid"`
}

type wxUserIdQuery struct {
	OpenId string `json:"openid"`
}

type wxUserIdResult struct {
	baseMessage
	UserId string `json:"userid"`
}

type WxPayApi struct {
	context *wxAppcontext
}

//将userID 转为openid
func (this *WxPayApi) ConvertUseIdToOpenId(userid string, useAgentid bool) (string, string) {
	var query interface{}
	if useAgentid {
		query = wxOpenIdExtQuery{wxOpenIdQuery{userid}, this.context.AuthAgentId}
	} else {
		query = wxOpenIdQuery{userid}
	}

	urlstr := fmt.Sprintf(USER_OPEN_API, this.context.auth.AccessToken)

	result := wxOpenIdResult{}
	err := PostToWx(urlstr, query, &result)
	if err == nil {
		return result.OpenId, result.AppId
	}

	return "", ""

}

//openid 转为userid
func (this *WxPayApi) ConvertOpenIdToUserId(openId string) string {
	query := wxUserIdQuery{openId}
	urlstr := fmt.Sprintf(OPEN_USER_API, this.context.auth.AccessToken)
	result := wxUserIdResult{}
	err := PostToWx(urlstr, query, &result)
	if err == nil {
		return result.UserId
	}

	return ""
}
