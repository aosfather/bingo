package wxcorp

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	CORP_ACCESSTOKEN_API = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	APPROVAL_API         = "https://qyapi.weixin.qq.com/cgi-bin/corp/getapprovaldata?access_token=%s"
	CHECKIN_API          = "https://qyapi.weixin.qq.com/cgi-bin/checkin/getcheckindata?access_token=%s"
)

//企业自身使用的accesstoken
type CorpPrivateAccessToken struct {
	UpdateTime time.Time
	CorpAccessToken
}

/*-----------------------------------------------------------------------------------
                        审批信息获取

-----------------------------------------------------------------------------------*/
type wxApprovalQuery struct {
	StartTime int64 `json:"starttime"`  //获取审批记录的开始时间。Unix时间戳
	EndTime   int64 `json:"endtime"`    //获取审批记录的结束时间。Unix时间戳
	NextID    int64 `json:"next_spnum"` //第一个拉取的审批单号，不填从该时间段的第一个审批单拉取
}

type wxApprovalResult struct {
	baseMessage
	Count  int          `json:"count"`
	Total  int          `json:"total"`
	NextID int64        `json:"next_spnum"`
	Data   []WxApproval `json:"data"`
}

type WxApproval struct {
	SpName        string           `json:"spname"`        //单据名称
	ApplyName     string           `json:"apply_name"`    //申请人姓名
	ApplyOrg      string           `json:"apply_org"`     //申请人部门
	ApprovalNames []string         `json:"approval_name"` //审批人姓名
	NotifyNames   []string         `json:"notify_name"`   //抄送人姓名
	SpStatus      int              `json:"sp_status"`     //审批状态：1审批中；2 已通过；3已驳回；4已取消
	SpId          int64            `json:"sp_num"`        //审批单号
	ApplyTime     int64            `json:"apply_time"`    //审批单提交时间
	ApplyUser     string           `json:"apply_user_id"` //审批单提交者的userid
	Leave         wxLeave          `json:"leave"`         //请假
	Expense       wxExpense        `json:"expense"`       //报销
	Comm          wxApprovalCustom `json:"comm"`          //审批模板信息
	Medias        []string         `json:"mediaids"`      //审批的附件media_id，可使用media/get获取附件
}

//请假
type wxLeave struct {
	Type      int    `json:"leave_type"` //请假类型：1年假；2事假；3病假；4调休假；5婚假；6产假；7陪产假；8其他
	Unit      int    `json:"timeunit"`   //请假时间单位：0半天；1小时
	StartTime int64  `json:"start_time"` //请假开始时间，unix时间
	EndTime   int64  `json:"end_time"`   //请假结束时间，unix时间
	Duration  int    `json:"duration"`   //请假时长，单位小时
	Reason    string `json:"reason"`     //请假事由

}

//报销
type wxExpense struct {
	Type   int             `json:"expense_type"` //报销类型：1差旅费；2交通费；3招待费；4其他报销
	Reason string          `json:"reason"`       // 	报销事由
	Items  []wxExpenseItem `json:"item"`         //报销明细
}

//费用明细
type wxExpenseItem struct {
	Type   int     `json:"expenseitem_type"` //费用类型：1飞机票；2火车票；3的士费；4住宿费；5餐饮费；6礼品费；7活动费；8通讯费；9补助；10其他
	Time   int64   `json:"time"`             //发生时间，unix时间
	Sum    float32 `json:"sums"`             //费用金额，单位元
	Reason string  `json:"reason"`           //明细事由
}

//自定义模板
type wxApprovalCustom struct {
	Data  string `json:"apply_data"`
	Items map[string]wxApprovalCustomItem
}

//自定义字段
type wxApprovalCustomItem struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Validate    bool   `json:"validate"`
	UnPrint     bool   `json:"un_print"`
	PlaceHolder string `json:"placeholder"`
	SetValue    string `json:"setvalue"`
	Warning     string `json:"warning"`
}

//申请单通知接口
type ApprovalNotify interface {
	NotifyApproval(approval *WxApproval) //单据通知
}
type WxApprovalApi struct {
	token   CorpPrivateAccessToken
	rate    int64          //查询频率，单位秒
	Notify  ApprovalNotify //监听通知
	running bool
}

func (this *WxApprovalApi) Run(second int64) {
	this.rate = second
	ticker := time.NewTicker(time.Second * time.Duration(this.rate))
	go func() {
		this.running = true
		for _ = range ticker.C {
			fmt.Printf("run query approval at %v\n", time.Now())
			this.doQuery()
		}
	}()
}

func (this *WxApprovalApi) doQuery() {
	now := GetCurrTs()
	this.queryNext(now-3600, now, 0)
}

func (this *WxApprovalApi) queryNext(start, end, next int64) {
	if this.token.Token == "" {
		return
	}
	url := fmt.Sprintf(APPROVAL_API, this.token.Token)
	query := wxApprovalQuery{}
	query.StartTime = start
	query.EndTime = end
	query.NextID = next
	fmt.Println(query)
	result := wxApprovalResult{}
	err := PostToWx(url, query, &result)
	if err == nil {
		if result.ErrCode == 0 {
			if this.Notify != nil {
				if result.Count == 0 {
					return
				}
				//一条条通知
				for _, sp := range result.Data {
					sp.Comm.Items = make(map[string]wxApprovalCustomItem)
					json.Unmarshal([]byte(sp.Comm.Data), &sp.Comm.Items)
					this.Notify.NotifyApproval(&sp)
				}
				//如果有下一页，继续
				if result.Total > result.Count {
					this.queryNext(start, end, result.NextID)
				}
			}
		} else {
			fmt.Println(result)
		}
	} else {
		fmt.Println(err.Error())
	}
}

/*----------------------------------------------------------------------
                            打卡记录获取



微信规定：
        获取记录时间跨度不超过三个月
        用户列表不超过100个。若用户超过100个，请分批获取
        有打卡记录即可获取打卡数据，与当前”打卡应用”是否开启无关


----------------------------------------------------------------------*/

type wxcheckinQuery struct {
	Type      int      `json:"opencheckindatatype"` //打卡类型。1：上下班打卡；2：外出打卡；3：全部打卡
	StartTime int64    `json:"starttime"`
	EndTime   int64    `json:"endtime"`
	Users     []string `json:"useridlist"` //需要获取打卡记录的用户列表
}

type wxcheckinResult struct {
	baseMessage
	List []WxcheckinItem `json:"result"`
}
type WxcheckinItem struct {
	User      string `json:"userid"`         //用户id
	Group     string `json:"groupname"`      //打卡规则名称
	Type      string `json:"checkin_type"`   //打卡类型
	Exception string `json:"exception_type"` //异常类型，如果有多个异常，以分号间隔
	Time      int64  `json:"checkin_time"`   //打卡时间。Unix时间戳

	Location string `json:"location_title"`  //打卡地点title
	Detail   string `json:"location_detail"` //打卡地点详情

	Wifi    string   `json:"wifiname"` //打卡wifi名称
	WifiMac string   `json:"wifimac"`  //打卡的MAC地址/bssid
	Notes   string   `json:"notes"`    //打卡备注
	Medias  []string `json:"mediaids"` //打卡的附件media_id，可使用media/get获取附件
}

//获取打卡记录
type WxCheckInApi struct {
	token CorpPrivateAccessToken
}

func (this *WxCheckInApi) QueryCheckInByAll(start, end int64, users []string) []WxcheckinItem {
	return this.query(3, start, end, users)

}

func (this *WxCheckInApi) QueryCheckInByWork(start, end int64, users []string) []WxcheckinItem {
	return this.query(1, start, end, users)

}

func (this *WxCheckInApi) QueryCheckInByOut(start, end int64, users []string) []WxcheckinItem {
	return this.query(2, start, end, users)

}

func (this *WxCheckInApi) query(checktype int, start, end int64, users []string) []WxcheckinItem {
	query := wxcheckinQuery{}
	query.Type = checktype
	query.StartTime = start
	query.EndTime = end
	query.Users = users

	result := wxcheckinResult{}
	url := fmt.Sprintf(CHECKIN_API, this.token.Token)
	err := PostToWx(url, query, &result)
	if err == nil {
		return result.List
	} else {
		fmt.Println(err.Error())
	}
	return nil
}

/*-------------------------------------------------------------
                 用于hold用户企业的sercret

--------------------------------------------------------------*/
type WxBaseApplicationHold struct {
	corp           string
	checkInSecret  string
	approvalSecret string
	messageSecret  string
	checkInApi     WxCheckInApi
	approvalApi    WxApprovalApi
	messageApi     WxSendMessageApi
}

func (this *WxBaseApplicationHold) Init(corp, checkInSecret, approvalSecret, messageSecret string) {
	this.corp = corp
	this.checkInSecret = checkInSecret
	this.approvalSecret = approvalSecret
	this.messageSecret = messageSecret
	this.autoRefresh()

}

func (this *WxBaseApplicationHold) autoRefresh() {
	ticker := time.NewTicker(time.Second * time.Duration(120))
	go func() {
		this.refreshAllToken()
		for _ = range ticker.C {
			fmt.Printf("autorefreshToken at %v\n", time.Now())
			//刷新token
			this.refreshAllToken()

		}
	}()
}

func (this *WxBaseApplicationHold) refreshAllToken() {
	//	this.refreshToken(this.checkInSecret, &this.checkInApi.token)
	//	this.refreshToken(this.approvalSecret, &this.approvalApi.token)
	this.refreshToken(this.messageSecret, &this.messageApi.token)
}

func (this *WxBaseApplicationHold) refreshToken(secret string, token *CorpPrivateAccessToken) {
	if secret != "" { //后续加入根据token过期时间来刷新
		nao := time.Duration(token.Expire) * time.Second
		theNow := time.Now()
		if token.UpdateTime.Add(nao).Before(theNow) {
			t := this.getAccessToken(secret)
			token.Token = t.Token
			token.Expire = t.Expire
			token.UpdateTime = theNow

		}
	}

}

func (this *WxBaseApplicationHold) getAccessToken(secret string) *CorpAccessToken {
	url := fmt.Sprintf(CORP_ACCESSTOKEN_API, this.corp, secret)
	content, err := HTTPGet(url)
	token := CorpAccessToken{}
	if err == nil {
		json.Unmarshal(content, &token)
		if token.ErrCode != 0 {
			fmt.Println(token.ErrMsg)
		}

	} else {
		fmt.Println(err.Error())
	}
	return &token

}

func (this *WxBaseApplicationHold) GetCheckInAPI() *WxCheckInApi {
	return &this.checkInApi
}

func (this *WxBaseApplicationHold) GetMessageApi() *WxSendMessageApi {
	return &this.messageApi
}

func (this *WxBaseApplicationHold) Run(second int64, notify ApprovalNotify) {
	if this.approvalApi.running {
		return
	}
	this.approvalApi.Notify = notify
	this.approvalApi.Run(second)

}

type WxSendMessageApi struct {
	token   CorpPrivateAccessToken
	agentId int
}

//发送消息
func (this *WxSendMessageApi) SendTextMessage(group int, whos []string, content string) {
	msg := wxTextMsg{}
	this.agentId = 3010040
	this.fillToWho("text", this.agentId, group, whos, &msg.wxBaseMsg)
	msg.Text.Content = content

	this.sendWxMessage(msg)

}

func (this *WxSendMessageApi) fillToWho(msgtype string, agentId, group int, whos []string, msg *wxBaseMsg) {
	msg.MsgType = msgtype
	msg.AgentId = agentId
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
}

//和微信发送消息
func (this *WxSendMessageApi) sendWxMessage(msg interface{}) wxMsgResult {
	targetUrl := fmt.Sprintf(SEND_MSG_API, this.token.CorpAccessToken.Token)
	result := wxMsgResult{}
	err := PostToWx(targetUrl, msg, &result)
	if err != nil {
		result.Errcode = 99999 //自定义的错误码
		result.ErrMsg = err.Error()
	}
	return result
}
