package wxcorp

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/aosfather/bingo"
	"github.com/aosfather/bingo/mvc"
)

const (
	WXAPI_HEADER      = "https://qyapi.weixin.qq.com/cgi-bin/service/"
	CATALOG_PERMANENT = "PermanentCode"
)

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

//部门
type WxDepart struct {
	Id     int64
	Name   string
	Parent int64
	Order  int64
}

func (this *WxDepart) Init(xmlbyte []byte) {
	depart := CorpChangeDepart{}
	xml.Unmarshal(xmlbyte, &depart)
	this.Id = depart.DepartId
	this.Name = depart.DepartName
	this.Parent = depart.DepartParent
	this.Order = depart.DepartOrder
}

func convertToWxDepart(xmlbyte []byte) WxDepart {
	depart := WxDepart{}
	depart.Init(xmlbyte)
	return depart
}

//标签
type WxTag struct {
	Id            string   //标签Id
	AddUserItems  []string //	标签中新增的成员userid列表，用逗号分隔
	DelUserItems  []string //	标签中删除的成员userid列表，用逗号分隔
	AddPartyItems []string //	标签中新增的部门id列表，用逗号分隔
	DelPartyItems []string //	标签中删除的部门id列表，用逗号分隔
}

func (this *WxTag) Init(xmlbyte []byte) {
	tag := CorpChangeTag{}
	xml.Unmarshal(xmlbyte, &tag)
	this.Id = tag.TagId
	this.AddUserItems = tag.AddUserItems
	this.DelUserItems = tag.DelUserItems
	this.AddPartyItems = tag.AddPartyItems
	this.DelPartyItems = tag.DelPartyItems

}

//用户
type WxUser struct {
	Id          string
	NewId       string            //新的UserID，变更时推送（userid由系统生成时可更改一次）
	Name        string            //成员名称
	Avatar      string            //头像
	Departments []int             //所属部门
	Mobile      string            //手机
	Position    string            //职位
	Gender      int               //性别，变更时推送。1表示男性，2表示女性
	Email       string            //邮箱，变更时推送 ，仅通讯录套件可获取
	Status      int               //激活状态：1=已激活 2=已禁用
	EnglishName string            //英文名
	IsLeader    int               //是否主管 标识是否为上级。0表示普通成员，1表示上级
	Telephone   string            //座机，仅通讯录套件可获取
	ExtAttr     map[string]string //扩展属性，变更时推送，仅通讯录套件可获取
}

func (this *WxUser) Init(xmlbyte []byte) {
	usr := CorpChangeContact{}
	xml.Unmarshal(xmlbyte, &usr)
	this.Id = usr.UserID
	this.NewId = usr.NewUserID
	this.Avatar = usr.Avatar
	this.Departments = usr.Department
	this.Mobile = usr.Mobile
	this.Position = usr.Position
	this.Gender = usr.Gender
	this.Email = usr.Email
	this.Status = usr.Status
	this.EnglishName = usr.EnglishName
	this.IsLeader = usr.IsLeader
	this.Telephone = usr.Telephone
	if len(usr.ExtAttr) > 0 {
		this.ExtAttr = make(map[string]string)
		for _, attr := range usr.ExtAttr {
			this.ExtAttr[attr.Name] = attr.Value

		}

	}

}

func convertToWxUser(xmlbyte []byte) WxUser {
	wxusr := WxUser{}
	wxusr.Init(xmlbyte)
	return wxusr
}

//组织通讯录表换授权
type CorpOrgChangeHandle interface {
	//新增部门
	AddDepart(depart WxDepart)
	//更新部门
	UpdateDepart(depart WxDepart)
	//删除部门
	DeleteDepart(depart int64)
	//更新标签
	UpdateTag(tag WxTag)
	//新增员工
	AddEmployee(user WxUser)
	//更新员工
	UpdateEmployee(user WxUser)
	//删除员工
	DeleteEmployee(userId string)
}

type WxValidateRequest struct {
	Timestamp string `Field:"timestamp"`
	Nonce     string `Field:"nonce"`
	Signature string `Field:"msg_signature"`
	Echostr   string `Field:"echostr"`
}

type WxSuitInputMsg struct {
	Timestamp string `Field:"timestamp"`
	Nonce     string `Field:"nonce"`
	Signature string `Field:"msg_signature"`
	data      CorpInputputMessage
}

func (this *WxSuitInputMsg) GetInput() CorpInputputMessage {
	return this.data
}
func (this *WxSuitInputMsg) GetData() interface{} {
	return &this.data
}
func (this *WxSuitInputMsg) GetDataType() string {
	return "xml"
}

/**
  访问获取token时候返回的消息

*/
type WxCorpSuite struct {
	mvc.SimpleController
	encryted         CorpEncrypt
	token            string //初始token
	corpId           string //企业id
	corpSecret       string //应用秘钥,用于消息解密
	suiteId          string //套件ID
	suiteSecret      string //套件秘钥
	suiteAccessToken SuiteAccessToken
	suiteTicket      string //套件的ticket
	suitePreAuthCode SuitePreAuthCode
	appHandle        CorpApplicationHandle //授权handle
	orgHandle        CorpOrgChangeHandle   //组织变化handle
	dataStage        CorpDataStage         //企业及永授权code存储
	hub              wxCorpApplicationHub
	contexts         map[string]*wxAuthCorpcontext
	//	contexts         map[string]*wxCorpAppContext
}

func (this *WxCorpSuite) Init(app bingo.Application, stage CorpDataStage) {
	prefix := app.Name
	this.dataStage = stage
	this.corpId = app.GetPropertyFromConfig(prefix + ".wx.corpid")
	this.corpSecret = app.GetPropertyFromConfig(prefix + ".wx.secret")
	this.token = app.GetPropertyFromConfig(prefix + ".wx.token")
	this.suiteId = app.GetPropertyFromConfig(prefix + ".wx.suite")
	this.suiteSecret = app.GetPropertyFromConfig(prefix + ".wx.suitesecret")
	this.encryted = CorpEncrypt{}
	this.encryted.Init(this.token, this.corpId, this.suiteId, this.corpSecret)
	this.contexts = make(map[string]*wxAuthCorpcontext)
	this.hub.Init(this)

	if this.dataStage != nil {
		this.suiteTicket = this.dataStage.GetCode("suite", this.suiteId)
	}

}

func (this *WxCorpSuite) SetHandle(app CorpApplicationHandle, contact CorpOrgChangeHandle) {
	this.appHandle = app
	this.orgHandle = contact
}

func (this *WxCorpSuite) GetParameType(method string) interface{} {
	fmt.Println(method)
	if method == "GET" {
		return &WxValidateRequest{}
	} else {
		return &WxSuitInputMsg{}
	}

}

func (this *WxCorpSuite) Get(c mvc.Context, p interface{}) (interface{}, mvc.BingoError) {
	if q, ok := p.(*WxValidateRequest); ok {
		fmt.Println(q)
		ret, result := this.encryted.VerifyURL(q.Signature, q.Timestamp, q.Nonce, q.Echostr)
		if ret == 0 {
			return result, nil
		}

	}

	return "hello", nil

}

//正常的访问消息处理
func (this *WxCorpSuite) Post(c mvc.Context, p interface{}) (interface{}, mvc.BingoError) {
	if msg, ok := p.(*WxSuitInputMsg); ok {
		fmt.Printf("/ndo:%s", msg)
		ret, result := this.encryted.DecryptInputMsg(msg.Signature, msg.Timestamp, msg.Nonce, msg.GetInput())
		fmt.Printf("/n result:%i,%s", ret, result)
		if ret == 0 {
			data := []byte(result)
			inputmsg := baseCorpMsg{}
			xml.Unmarshal(data, &inputmsg)
			fmt.Printf("/nevent:%s", inputmsg)
			switch inputmsg.Type {
			case "suite_ticket": //推送suite_ticket
				this.processSuiteTicketMsg(data)
			case "create_auth": //授权成功通知
				if this.processCorpAuthMsg(data) {

					return "success", nil
				}

				return "get permanent error", nil

			case "change_auth": //变更授权通知
				this.processCorpAuthChangeMsg(data, false)
			case "cancel_auth": //取消授权通知
				this.processCorpAuthChangeMsg(data, true)
			case "change_contact": //通讯录变更通知
				this.processCorpContactMsg(data)
			}

		}

	}

	return "hi", nil

}

//解析推送的ticket消息，并更新suiteTicket
func (this *WxCorpSuite) processSuiteTicketMsg(msg []byte) {
	xmlmsg := CorpTickMsg{}
	xml.Unmarshal(msg, &xmlmsg)
	fmt.Printf("/nticket:%s", xmlmsg)
	this.suiteTicket = xmlmsg.Ticket
	if this.dataStage != nil {
		this.dataStage.SaveCode("suite", this.suiteId, this.suiteTicket)
	}
	fmt.Println(this.suiteTicket)

}

func (this *WxCorpSuite) processCorpAuthMsg(msg []byte) bool {
	xmlmsg := CorpAuthMsg{}
	xml.Unmarshal(msg, &xmlmsg)
	fmt.Printf("/ncorpauth:%s", xmlmsg)
	//通过得到的authcode获取永久授权
	this.GetPermanentCode(xmlmsg.AuthCode)

	//通知授权了。

	return true

}

func (this *WxCorpSuite) processCorpAuthChangeMsg(msg []byte, del bool) {
	xmlmsg := CorpChangeMsg{}
	xml.Unmarshal(msg, &xmlmsg)
	corpId := xmlmsg.CorpId
	if del { //删除授权
		if this.appHandle != nil {
			this.appHandle.DeleteAuth(corpId)
		}

	} else { //授权变更
		if this.appHandle != nil {
			//
			//通过corpId获取变更信息
			info := this.getCorpAuthInfo(corpId)
			this.appHandle.UpdateAuth(info.Corpinfo, info.Agents.Agents)
		}

	}

}

//处理通讯录变更消息
func (this *WxCorpSuite) processCorpContactMsg(msg []byte) {
	if this.orgHandle == nil {
		return
	}
	xmlmsg := CorpChangeMsg{}
	xml.Unmarshal(msg, &xmlmsg)
	changeType := xmlmsg.ChangeType
	switch changeType {
	//处理员工信息变更
	case "create_user":
		this.orgHandle.AddEmployee(convertToWxUser(msg))

	case "update_user":
		this.orgHandle.UpdateEmployee(convertToWxUser(msg))

	case "delete_user":
		this.orgHandle.DeleteEmployee(convertToWxUser(msg).Id)

	//处理department变更
	case "create_party":
		this.orgHandle.AddDepart(convertToWxDepart(msg))

	case "update_party":
		this.orgHandle.UpdateDepart(convertToWxDepart(msg))
	case "delete_party":
		this.orgHandle.DeleteDepart(convertToWxDepart(msg).Id)

	case "update_tag": //处理tag变更
		tag := WxTag{}
		tag.Init(msg)
		this.orgHandle.UpdateTag(tag)

	}

}

func (this *WxCorpSuite) getCorpAuthInfo(corpId string) CorpAuthInfoResult {
	query := CorpAccessTokenQuery{}
	query.CorpId = corpId
	query.Id = this.suiteId
	query.PermanentCode = this.dataStage.GetCode(CATALOG_PERMANENT, corpId)
	result := CorpAuthInfoResult{}
	this.callApi("get_auth_info", query, &result)
	return result
}

func (this *WxCorpSuite) RefreshSuiteToken() {
	nao := time.Duration(this.suiteAccessToken.Expires) * time.Second
	if this.suiteAccessToken.UpdateTime.Add(nao).Before(time.Now()) {
		query := SuiteTokenQuery{}
		query.Id = this.suiteId
		query.Secret = this.suiteSecret
		query.Ticket = this.suiteTicket

		theUrl := WXAPI_HEADER + "get_suite_token"
		accessToken := SuiteAccessToken{}
		err := PostToWx(theUrl, query, &accessToken)
		if err == nil {
			if accessToken.ErrCode == 0 {
				this.suiteAccessToken.AccessToken = accessToken.AccessToken
				this.suiteAccessToken.Expires = accessToken.Expires
				this.suiteAccessToken.UpdateTime = time.Now()
			}
		}
	}

}

func (this *WxCorpSuite) RefreshPreAuthCode() {
	//如果preauthcode过期，则执行刷新
	nao := time.Duration(this.suitePreAuthCode.Expires) * time.Second
	if this.suitePreAuthCode.UpdateTime.Add(nao).Before(time.Now()) {
		query := baseQuery{}
		query.Id = this.suiteId
		//访问api
		preauthCode := SuitePreAuthCode{}
		if this.callApi("get_pre_auth_code", query, &preauthCode) {
			if preauthCode.ErrCode == 0 {
				this.suitePreAuthCode.UpdateTime = time.Now()
				this.suitePreAuthCode.Expires = preauthCode.Expires
				this.suitePreAuthCode.AuthCode = preauthCode.AuthCode
			}
		}

	}

}

func (this *WxCorpSuite) SetSessionInfo(apps []string, authType int) bool {
	//1、先刷新preauthcode
	this.RefreshPreAuthCode()
	//2、构建设置参数
	info := SuiteSessionInfo{}
	info.AuthCode = this.suitePreAuthCode.AuthCode
	info.Info.Id = apps
	info.Info.Type = authType

	//3、访问api
	msg := baseMessage{}
	if this.callApi("set_session_info", info, &msg) {
		return msg.ErrCode == 0
	}
	return false

}
func (this *WxCorpSuite) GetPermanentCode(authCode string) bool {
	this.RefreshSuiteToken()
	query := CorpPermanentCodeQuery{}
	query.Id = this.suiteId
	query.AuthCode = authCode
	permanent := CorpPermanentAuth{}
	if this.callApi("get_permanent_code", query, &permanent) {
		fmt.Printf("get permanentCode:%v", permanent)

		//存储permanent code
		if this.dataStage != nil {
			this.dataStage.SaveCode(CATALOG_PERMANENT, permanent.Corpinfo.Id, permanent.PermanentCode)
			//
		}
		//通知handle，授权事件
		if this.appHandle != nil {
			this.appHandle.NewCorp(permanent.Corpinfo, permanent.Agents.Agents, permanent.Admin)
		}
		// 通知 通讯录变更处理，同步通讯录

		this.initCorpContact(permanent.Corpinfo.Id)

		return true
	}
	return false

}

//同步通讯录
func (this *WxCorpSuite) initCorpContact(corpId string) {
	if this.orgHandle == nil {
		return

	}
	context := this.GetAuthCorpContext(corpId)
	if context != nil {
		contactApi := WxContactApi{context}
		deptList := contactApi.GetDepartmentList()
		//循环创建新增部门
		for _, dept := range deptList.List {
			this.orgHandle.AddDepart(dept.ToWxDepart())
			//遍历部门内的员工进行同步
			usrList := contactApi.GetDepartmentUserList(dept.Id)
			for _, usr := range usrList.List {
				this.orgHandle.AddEmployee(usr.ToWxUser())
			}

		}

	}

}

//获取授权企业的访问token
func (this *WxCorpSuite) GetCorpAccessToken(corpId, permanentCode string) CorpAccessToken {
	this.RefreshSuiteToken()
	query := CorpAccessTokenQuery{}
	query.Id = this.suiteId
	query.CorpId = corpId
	query.PermanentCode = permanentCode

	accessToken := CorpAccessToken{}
	if this.callApi("get_corp_token", query, &accessToken) {
		return accessToken

	}
	return CorpAccessToken{}
}

func (this *WxCorpSuite) GetAuthCorpApi() WxAuthCorpApi {
	return &this.hub
}

func (this *WxCorpSuite) GetAuthCorpContext(corpid string) *wxAuthCorpcontext {
	context := this.contexts[corpid]
	if context == nil && this.dataStage != nil {
		//通过永久授权码获取访问accesstoken
		pcode := this.dataStage.GetCode(CATALOG_PERMANENT, corpid)
		token := this.GetCorpAccessToken(corpid, pcode)
		context = &wxAuthCorpcontext{corpid, token.Token}
		this.contexts[corpid] = context

	}

	return context

}

func (this *WxCorpSuite) callApi(apiname string, data interface{}, response interface{}) bool {
	theUrl := WXAPI_HEADER + apiname + "?suite_access_token=" + this.suiteAccessToken.AccessToken
	fmt.Println(theUrl)
	fmt.Println(data)
	err := PostToWx(theUrl, data, response)
	if err == nil {

		return true
	} else {
		fmt.Println(err.Error())
	}
	return false
}
