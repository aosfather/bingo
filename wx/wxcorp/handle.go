package wxcorp

//企业token等信息存储
type CorpDataStage interface {
	SaveCode(catalog, corpId, code string) //保存code
	GetCode(catalog, corpId string) string //获取保存的Code
}

//应用授权的handle
type CorpApplicationHandle interface {
	//新授权应用
	NewCorp(corp CorpAuthInfo, agents []CorpAgentInfo, admin CorpAdminInfo)
	//更新授权
	UpdateAuth(corp CorpAuthInfo, agents []CorpAgentInfo)
	//删除授权
	DeleteAuth(corpid string)
}

//组织通讯录表换授权
type CorpOrgChangeHandle interface {
	//新增部门
	AddDepart(corpid string, depart WxDepart)
	//更新部门
	UpdateDepart(corpid string, depart WxDepart)
	//删除部门
	DeleteDepart(corpid string, depart int64)
	//更新标签
	UpdateTag(corpid string, tag WxTag)
	//新增员工
	AddEmployee(corpid string, user WxUser)
	//更新员工
	UpdateEmployee(corpid string, user WxUser)
	//删除员工
	DeleteEmployee(corpid string, userId string)
}

type WXMsgType byte

const (
	WMT_UNKNOWN  WXMsgType = 0
	WMT_Text     WXMsgType = 11
	WMT_Image    WXMsgType = 12
	WMT_Voice    WXMsgType = 13
	WMT_Video    WXMsgType = 14
	WMT_Location WXMsgType = 15
	WMT_Link     WXMsgType = 16
)

type ReplyMessage struct {
	Type        WXMsgType
	Title       string //text 为内容，其它为Media id
	Description string
	Url         string
}

//被动消息处理的handle
type MessageReplyHandle interface {
	ReplyMsg(t WXMsgType, content string, url string, ext string) *ReplyMessage
}
