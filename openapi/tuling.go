package openapi

import (
	"github.com/aosfather/comb"
	"github.com/aosfather/bingo/utils"
)

const (
	API_URL = "http://www.tuling123.com/openapi/api"
)

type TulingRequest struct {
	Key  string `json:"key"`
	Text string `json:"info"`
	Loc  string `json:"loc"`
	User string `json:"userid"`
}

/**

100000 	文本类
200000 	链接类 有url
302000 	新闻类 有list
308000 	菜谱类


*/
type TulingRespone struct {
	Code int64        `json:"code"`
	Text string       `json:"text"`
	URL  string       `json:"url"`
	List []TulingItem `json:"list"`
}

type TulingItem struct {
	Name    string `json:"name"`      //菜名
	Info    string `json:"info"`      //菜谱信息
	Icon    string `json:"icon"`      //信息图标
	Url     string `json:"detailurl"` //详情链接
	Article string `json:"article"`   //文章标题
	Source  string `json:"source"`    //来源
}

type TulingSDK struct {
	Key string
}

func (this *TulingSDK) Query(user, text string) TulingRespone {
	request := TulingRequest{}
	request.Key = this.Key
	request.User = user
	request.Text = text
	result := TulingRespone{}
	utils.Post(API_URL, request, &result)
	return result

}
