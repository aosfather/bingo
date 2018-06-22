package wxcorp

import (
	"fmt"
	"net/url"
)

/**
企业微信获取用户信息
流程：
  用户访问应用，应用跳转到企业微信进行验证
  企业微信验证后，通过跳转压入code参数
  应用通过code获取用户信息
  用于根据用户信息返回的user_ticket来获取用户的企业通讯录信息



*/

const (
	OAUTH_SERVICE_API = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&agentid=%d&state=%s#wechat_redirect"
	USER_DETAIL_API   = "https://qyapi.weixin.qq.com/cgi-bin/user/getuserdetail?access_token=%s"
	USER_INFO_API     = "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=%s&code=%s"
	/*
			 	应用授权作用域。
		snsapi_base：静默授权，可获取成员的基础信息；
		snsapi_userinfo：静默授权，可获取成员的详细信息，但不包含手机、邮箱；
		snsapi_privateinfo：手动授权，可获取成员的详细信息，包含手机、邮箱。
	*/
	OAUTH_LEVEL_BASE   = "snsapi_base"
	OAUTH_LEVEL_INFO   = "snsapi_userinfo"
	OAUTH_LEVEL_DETAIL = "snsapi_privateinfo"
)

//微信回调参数
type WxRedirectParamter struct {
	Code  string `Field:"code"`
	State string `Field:"state"`
}

type wxUserInfo struct {
	baseMessage
	UserId   string `json:"UserId"`      //成员UserID
	OpenId   string `json:"OpenId"`      //非企业成员的标识，对当前企业唯一
	DeviceId string `json:"DeviceId"`    //手机设备号(由企业微信在安装时随机生成，删除重装会改变，升级不受影响)
	Ticket   string `json:"user_ticket"` //成员票据，最大为512字节。scope为snsapi_userinfo或snsapi_privateinfo，且用户在应用可见范围之内时返回此参数。后续利用该参数可以获取用户信息或敏感信息。
	Expire   int64  `json:"expires_in"`  //user_token的有效时间（秒），随user_ticket一起返回
}

type wxUserDetailQuery struct {
	Ticket string `json:"user_ticket"`
}

type WxUserDetail struct {
	Id          string `json:"userid"`     //成员UserID
	Name        string `json:"name"`       //成员姓名
	Departments []int  `json:"department"` //成员所属部门
	Position    string `json:"position"`   //职位信息
	Mobile      string `json:"mobile"`     //成员手机号，仅在用户同意snsapi_privateinfo授权时返回
	Gender      int    `json:"gender"`     //性别。0表示未定义，1表示男性，2表示女性
	Email       string `json:"email"`      //成员邮箱，仅在用户同意snsapi_privateinfo授权时返回
	Avatar      string `json:"avatar"`     //头像url。例如："http://shp.qpic.cn/bizmp/xxxxxxxxxxx/0" 注：如果要获取小图将url最后的”/0”改成”/100”即可
}

func BuildRedirectUrl(corpid string, agentid int, redirectURI, level, state string) string {
	//重定向后会带上state参数，企业可以填写a-zA-Z0-9的参数值，长度不可超过128个字节
	if !checkStateContent(state) {
		return ""
	}

	urlStr := url.QueryEscape(redirectURI)
	if corpid == "" {
		return fmt.Sprintf(OAUTH_SERVICE_API, "$CORPID$", urlStr, level, "$AGENTID$", state)
	} else {

		return fmt.Sprintf(OAUTH_SERVICE_API, corpid, urlStr, level, agentid, state)
	}

}

func checkStateContent(state string) bool {
	if len(state) > 128 {
		return false
	}

	return true
}
