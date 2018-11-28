package openapi

import (
	"github.com/aosfather/bingo/utils"
	"fmt"
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

//type News struct {
//	Article   string `json:"article"`
//	Source    string `json:"source"`
//	Icon      string `json:"icon"`
//	DetailURL string `json:"detailurl"`
//}
//
//type Menu struct {
//	Name      string `json:"name"`
//	Icon      string `json:"icon"`
//	Info      string `json:"info"`
//	DetailURL string `json:"detailurl"`
//}

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

func (this *TulingSDK)QueryAsString(user,text string) string {
	reply:=this.Query(user,text)

	switch reply.Code {
	case 100000:
		return reply.Text
	case 200000:
		return reply.Text + " " + reply.URL
	case 302000:
		var res string
		news := reply.List
		for _, n := range news {
			res += fmt.Sprintf("%s\n%s\n", n.Article, n.Url)
		}

		return res
	case 308000:
		var res string
		menu := reply.List
		for _, m := range menu {
			res += fmt.Sprintf("%s\n%s\n%s\n", m.Name, m.Info, m.Url)
		}
		return res
	default:
		return "不知道你在说啥～"
	}

}
