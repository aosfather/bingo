package openapi

import (
	"fmt"
	"net/url"
	"github.com/aosfather/bingo/utils"
	"encoding/json"
	"strings"
)

/**
有道

 */
 const _YOUDAO_API_URL="http://fanyi.youdao.com/openapi.do?keyfrom=go-aida&key=145986666&type=data&doctype=json&version=1.1&q=%s"

 type YoudaoResponse struct {
    Translation []string `json:"translation"` //翻译内容
    Basic basicDict `json:"basic"` //基础词典
 	Query string `json:"query"` //查询的词
 	Code int `json:"errorCode"`  //错误码 0 -代表成功
 	WebItems []webItem `json:"web"` //网络词条
 	/*

 	{"translation":["男孩"],
 	"basic":{"us-phonetic":"bɔɪ","phonetic":"bɒɪ","uk-phonetic":"bɒɪ",
 	"explains":["n. 男孩；男人","n. (Boy)人名；(英、德、西、意、刚(金)、印尼、瑞典)博伊；(法)布瓦"]},
 	"query":"boy","errorCode":0,
 	"web":[{"value":["男孩","男孩","毛利男孩"],"key":"Boy"}
 	,{"value":["老男孩","老生","原罪犯"],"key":"Old boy"}
 	,{"value":["快乐男声","超级男孩","元气超人"],"key":"Super Boy"}]}
 	 */
 }

 type basicDict struct {
 	Us string  `json:"us-phonetic"`  //美式发音
 	Uk string  `json:"uk-phonetic"`  //英式发音
 	Stand string `json:"phonetic"`   //标准发音
 	Explains []string  `json:"explains"` //解释

 }

 type webItem struct {
 	Key string
 	Value []string
 }

 func QueryFromYoudao(msg string) *YoudaoResponse {
 	u:=fmt.Sprintf(_YOUDAO_API_URL,url.QueryEscape(msg))
 	data,err:=utils.HTTPGet(u)
 	if err==nil {
 		res:=YoudaoResponse{}
 		json.Unmarshal(data,&res)
 		return &res
	}
 	return nil
 }

 func QueryFromYoudaoAsString(msg string) string {
 	res:=QueryFromYoudao(msg)
 	if res!=nil {
      if res.Code!=0 {
      	return "查询出错了！未找到词条"
	  }

	  if len(res.Translation) >0 {
	  	var str string
	  	for _,item:=range res.Translation {
	  		str+=item+"\n"
		}
		if res.Basic.Stand!="" {
			str+="\n词典：\n"
			str+=fmt.Sprintf("%s：[%s] (us[%s],uk[%s])\n",res.Query,res.Basic.Stand,res.Basic.Us,res.Basic.Uk)
			for _,item:=range res.Basic.Explains {
				str+="\t"+item+"\n"
			}
		}

		if len(res.WebItems) >0{
			str+="\n网络词条:\n"
			for index,item:= range res.WebItems {
				str+=fmt.Sprintf("\t%d)、%s\t %s\n",index+1,item.Key,strings.Join(item.Value,"；"))
			}
		}

		return  str
	  }
	}

	return "没有找到！"
 }
