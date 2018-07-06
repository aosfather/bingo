package wxcorp

import (
	"encoding/xml"
	"time"
)

/**

第三方suit授权体系
*/
type baseMessage struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type baseToken struct {
	baseMessage

	Expires    int `json:"expires_in"`
	UpdateTime time.Time
}

type SuiteAccessToken struct {
	baseToken
	AccessToken string `json:"suite_access_token"`
}

type SuitePreAuthCode struct {
	baseToken
	AuthCode string `json:"pre_auth_code"`
}
type sessionInfo struct {
	Id   []string `json:"appid"`
	Type int      `json:"auth_type"` //授权类型：0 正式授权， 1 测试授权， 默认值为0
}
type SuiteSessionInfo struct {
	AuthCode string `json:"pre_auth_code"`
	Info     sessionInfo
}

type baseQuery struct {
	Id string `json:"suite_id"` //套件id
}

//永久授权码查询
type CorpPermanentCodeQuery struct {
	baseQuery
	AuthCode string `json:"auth_code"`
}

type SuiteTokenQuery struct {
	baseQuery
	Secret string `json:"suite_secret"` //套件秘钥
	Ticket string `json:"suite_ticket"` //访问用的ticket
}

type CorpAccessTokenQuery struct {
	baseQuery
	CorpId        string `json:"auth_corpid"`
	PermanentCode string `json:"permanent_code"`
}

type CorpAccessToken struct {
	baseMessage
	Token string `json:"access_token"`
	//有效期，单位秒
	Expire int `json:"expires_in"`
}

//企业永久授权码
type CorpPermanentAuth struct {
	AccessToken   string        `json:"access_token"`   //授权方（企业）access_token,最长为512字节
	Expires       int           `json:"expires_in"`     //授权方（企业）access_token超时时间
	PermanentCode string        `json:"permanent_code"` //企业微信永久授权码,最长为512字节
	Corpinfo      CorpAuthInfo  `json:"auth_corp_info"` //授权方企业信息
	Agents        AuthInfo      `json:"auth_info"`      //授权的应用信息
	Admin         CorpAdminInfo `json:"auth_user_info"` //授权管理员的信息
}

type AuthInfo struct {
	Agents []CorpAgentInfo `json:"agent"` //授权的应用信息
}

type CorpAuthInfoResult struct {
	Corpinfo CorpAuthInfo `json:"auth_corp_info"` //授权方企业信息
	Agents   AuthInfo     `json:"auth_info"`      //授权的应用信息
}

//授权企业信息
type CorpAuthInfo struct {
	Id              string        `json:"corpid"`               //授权方企业微信id
	Name            string        `json:"corp_name"`            //	授权方企业微信名称
	Type            string        `json:"corp_type"`            //	授权方企业微信类型，认证号：verified, 注册号：unverified
	Logo            string        `json:"corp_square_logo_url"` //授权方企业微信方形头像
	UserMax         int           `json:"corp_user_max"`        //授权方企业微信用户规模
	AgentMax        int           `json:"corp_agent_max"`
	FullName        string        `json:"corp_full_name"`    //所绑定的企业微信主体名称
	SubjectType     int           `json:"subject_type"`      //企业类型，1. 企业; 2. 政府以及事业单位; 3. 其他组织, 4.团队号
	Wxqrcode        string        `json:"corp_wxqrcode"`     //授权方企业微信二维码
	VerifiedEndTime time.Duration `json:"verified_end_time"` //认证到期时间
}

type CorpAgentInfo struct {
	Agentid    string   `json:"agentid"`         //授权方应用id
	Name       string   `json:"name"`            //授权方应用名字
	SquareLogo string   `json:"square_logo_url"` //授权方应用方形头像
	roundLogo  string   `json:"round_logo_url"`  //授权方应用圆形头像
	Appid      string   `json:"appid"`           //服务商套件中的对应应用id
	Privilege  string   `json:"privilege"`       //应用对应的权限
	AllowParty []string `json:"allow_party"`     //应用可见范围（部门）
	AllowTag   []string `json:"allow_tag"`       //应用可见范围（标签）
	AllowUser  []string `json:"allow_user"`      //应用可见范围（成员）
	ExtraParty []string `json:"extra_party"`     //额外通讯录（部门）
	ExtraUser  []string `json:"extra_user"`      //额外通讯录（成员）
	ExtraTag   []string `json:"extra_tag"`       //额外通讯录（标签）
	Level      int      `json:"email"`           //权限等级。1:通讯录基本信息只读2:通讯录全部信息只读3:通讯录全部信息读写4:单个基本信息只读5:通讯录全部信息只写
}

//授权管理员的信息
type CorpAdminInfo struct {
	Email  string `json:"email"`  //授权管理员的邮箱，可能为空（外部管理员一定有，不可更改）
	Mobile string `json:"mobile"` //授权管理员的手机号，可能为空（内部管理员一定有，可更改）
	Userid string `json:"userid"` //授权管理员的userid，可能为空（内部管理员一定有，不可更改）
	Name   string `json:"name"`   //授权管理员的name，可能为空（内部管理员一定有，不可更改）
	Avatar string `json:"avatar"` //授权管理员的头像url
}

/**第三方回调协议
推送suite_ticket
授权成功通知
变更授权通知
取消授权通知
通讯录变更事件通知
 新增成员事件
 更新成员事件
 删除成员事件
 新增部门事件
 更新部门事件
 删除部门事件
 标签变更事件

*/

type baseCorpMsg struct {
	XMLName   xml.Name `xml:"xml"`
	Id        string   `xml:"SuiteId"`
	Type      string   `xml:"InfoType"`
	TimeStamp string   `xml:"TimeStamp"`
}

//推送suite_ticket
type CorpTickMsg struct {
	XMLName xml.Name `xml:"xml"`
	baseCorpMsg
	Ticket string `xml:"SuiteTicket"`
}

//授权成功通知
type CorpAuthMsg struct {
	XMLName xml.Name `xml:"xml"`
	baseCorpMsg
	AuthCode string `xml:"AuthCode"`
}

//变更授权通知:包括取消授权type：cancel_auth取消授权change_auth授权变更
type CorpChangeMsg struct {
	XMLName xml.Name `xml:"xml"`
	baseCorpMsg
	ChangeType string
	CorpId     string `xml:"AuthCorpId"` //	授权方的corpid
}

//变更部门信息
type CorpChangeDepart struct {
	XMLName xml.Name `xml:"xml"`
	CorpChangeMsg
	DepartId     int64  `xml:"Id"`
	DepartName   string `xml:"Name"`
	DepartParent int64  `xml:"ParentId"`
	DepartOrder  int64  `xml:"Order"`
}

//标签变更事件
type CorpChangeTag struct {
	XMLName xml.Name `xml:"xml"`
	CorpChangeMsg
	TagId         string   //标签Id
	AddUserItems  []string //	标签中新增的成员userid列表，用逗号分隔
	DelUserItems  []string //	标签中删除的成员userid列表，用逗号分隔
	AddPartyItems []string //	标签中新增的部门id列表，用逗号分隔
	DelPartyItems []string //	标签中删除的部门id列表，用逗号分隔
}

type CorpContactItem struct {
	Name  string
	Value string
}

//通讯录员工变更事件
type CorpChangeContact struct {
	XMLName xml.Name `xml:"xml"`
	CorpChangeMsg
	UserID      string
	NewUserID   string
	Name        string
	Department  string
	Mobile      string
	Position    string
	Gender      int
	Email       string
	Status      int
	Avatar      string
	EnglishName string
	IsLeader    int
	Telephone   string
	ExtAttr     []CorpContactItem `xml:"ExtAttr>Item"`
}
