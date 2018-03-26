package wx

import (
	"encoding/json"
	"fmt"
	"time"
)

const getTicketURL = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
const AccessTokenURL = "https://api.weixin.qq.com/cgi-bin/token"



//
type WxApplication struct {
	AppID      string
	Ticket     *resTicket
	AppSecret  string
	UpdateTime time.Time
}

type baseMessage struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

//ResAccessToken struct
type ResAccessToken struct {
	baseMessage

	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// resTicket 请求jsapi_tikcet返回结果
type resTicket struct {
	baseMessage

	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"`
}



//GetTicket 获取jsapi_tocket
func (this *WxApplication) GetTicket() (ticketStr string, err error) {
	if this.Ticket == nil {
		now := time.Now()
		var ticket resTicket
		ticket, err = this.getTicketFromServer()

		if err != nil {
			return
		}
		this.Ticket = &ticket
		this.UpdateTime = now
		ticketStr = ticket.Ticket
	} else {
		nao := time.Duration(7200) * time.Second
		if this.UpdateTime.Add(nao).Before(time.Now()) {
			this.Ticket = nil
			return this.GetTicket()
		}
		ticketStr = this.Ticket.Ticket
	}

	return
}

//getTicketFromServer 强制从服务器中获取ticket
func (this *WxApplication) getTicketFromServer() (ticket resTicket, err error) {
	var accessToken string
	accessToken, err = this.GetAccessToken()
	if err != nil {
		return
	}

	var response []byte
	url := fmt.Sprintf(getTicketURL, accessToken)
	response, err = HTTPGet(url)
	fmt.Println(string(response))
	err = json.Unmarshal(response, &ticket)
	if err != nil {
		return
	}
	if ticket.ErrCode != 0 {
		err = fmt.Errorf("getTicket Error : errcode=%d , errmsg=%s", ticket.ErrCode, ticket.ErrMsg)
		return
	}

	return
}

func (this *WxApplication) GetAccessToken() (accessToken string, err error) {

	//从微信服务器获取
	var resAccessToken ResAccessToken
	resAccessToken, err = this.GetAccessTokenFromServer()
	if err != nil {
		return
	}

	accessToken = resAccessToken.AccessToken
	return
}

//GetAccessTokenFromServer 强制从微信服务器获取token
func (this *WxApplication) GetAccessTokenFromServer() (resAccessToken ResAccessToken, err error) {
	url := fmt.Sprintf("%s?grant_type=client_credential&appid=%s&secret=%s", AccessTokenURL, this.AppID, this.AppSecret)
	var body []byte
	body, err = HTTPGet(url)
	if err != nil {
		return
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &resAccessToken)
	if err != nil {
		return
	}
	if resAccessToken.ErrMsg != "" {
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v", resAccessToken.ErrCode, resAccessToken.ErrMsg)
		return
	}

	return
}


