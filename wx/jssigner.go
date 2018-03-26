package wx

import (
	"encoding/json"
	"github.com/aosfather/bingo/mvc"
	"fmt"

)

type signerRequest struct {
	Url      string `Field:"url"`
	Callback string `Field:"callback"`
}

// jsConfig 返回给用户jssdk配置信息
type jsConfig struct {
	AppID     string `json:"appId"`
	Timestamp int64  `json:"timestamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
}

type JsSigner struct {
	mvc.SimpleController
	App        *WxApplication
	DefaultStr string
}

//GetConfig 获取jssdk需要的配置参数
//uri 为当前网页地址
func (this *JsSigner) GetConfig(uri string) (config *jsConfig, err error) {
	config = new(jsConfig)
	var ticketStr string
	ticketStr, err = this.App.GetTicket()
	if err != nil {
		return
	}

	nonceStr :=this.DefaultStr
	if this.DefaultStr ==""{
		nonceStr=RandomStr(16)
	}
	timestamp := GetCurrTs()
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticketStr, nonceStr, timestamp, uri)
	fmt.Println(str)
	sigStr := Signature(str)

	config.AppID = this.App.AppID
	config.NonceStr = nonceStr
	config.Timestamp = timestamp
	config.Signature = sigStr
	return
}

func (this *JsSigner) GetParameType(method string) interface{} {
	return &signerRequest{}

}

func (this *JsSigner) Get(c mvc.Context, p interface{}) (interface{}, mvc.BingoError) {
	if value, ok := p.(*signerRequest); ok {
		if value.Url == "" {
			return "input parameter is wrong!", nil
		}
		fmt.Println(value.Url)
		config, err := this.GetConfig(value.Url)
		if err == nil {
			jsonData, _ := json.Marshal(config)
			result := value.Callback + "(" + string(jsonData) + ")"
			return result, nil
		}

		return err.Error(), nil
	}

	return "input parameter is wrong!", nil

}

