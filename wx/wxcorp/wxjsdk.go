package wxcorp

import (
	"encoding/json"
	"fmt"
	"time"
)

const getTicketURL = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
const AccessTokenURL = "https://api.weixin.qq.com/cgi-bin/token"

// Js struct
type Js struct {
	AppID      string
	Ticket     *resTicket
	AppSecret  string
	UpdateTime time.Time
}

// Config 返回给用户jssdk配置信息
type Config struct {
	AppID     string `json:"appId"`
	Timestamp int64  `json:"timestamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
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

//GetConfig 获取jssdk需要的配置参数
//uri 为当前网页地址
func (js *Js) GetConfig(uri string) (config *Config, err error) {
	config = new(Config)
	var ticketStr string
	ticketStr, err = js.GetTicket()
	if err != nil {
		return
	}

	nonceStr := "123456789qaw" //RandomStr(16)
	timestamp := GetCurrTs()
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticketStr, nonceStr, timestamp, uri)
	fmt.Println(str)
	sigStr := Signature(str)

	config.AppID = js.AppID
	config.NonceStr = nonceStr
	config.Timestamp = timestamp
	config.Signature = sigStr
	return
}

//GetTicket 获取jsapi_tocket
func (js *Js) GetTicket() (ticketStr string, err error) {
	if js.Ticket == nil {
		now := time.Now()
		var ticket resTicket
		ticket, err = js.getTicketFromServer()

		if err != nil {
			return
		}
		js.Ticket = &ticket
		js.UpdateTime = now
		ticketStr = ticket.Ticket
	} else {
		nao := time.Duration(7200) * time.Second
		if js.UpdateTime.Add(nao).Before(time.Now()) {
			js.Ticket = nil
			return js.GetTicket()
		}
		ticketStr = js.Ticket.Ticket
	}

	return
}

//getTicketFromServer 强制从服务器中获取ticket
func (js *Js) getTicketFromServer() (ticket resTicket, err error) {
	var accessToken string
	accessToken, err = js.GetAccessToken()
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

func (js *Js) GetAccessToken() (accessToken string, err error) {

	//从微信服务器获取
	var resAccessToken ResAccessToken
	resAccessToken, err = js.GetAccessTokenFromServer()
	if err != nil {
		return
	}

	accessToken = resAccessToken.AccessToken
	return
}

//GetAccessTokenFromServer 强制从微信服务器获取token
func (js *Js) GetAccessTokenFromServer() (resAccessToken ResAccessToken, err error) {
	url := fmt.Sprintf("%s?grant_type=client_credential&appid=%s&secret=%s", AccessTokenURL, js.AppID, js.AppSecret)
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
