package openapi

import (
	"github.com/aosfather/bingo/utils"
	"fmt"
	"strings"
)

/**
魔力 茉莉 机器人 api
 */

 const _MOLI_API_URL="http://i.itpk.cn/api.php?question=%s"

 func QueryByMoli(msg string) string {
 	url:=fmt.Sprintf(_MOLI_API_URL,msg)
 	data,err:=utils.HTTPGet(url)
 	if err!=nil {
 		return "休息，休息，休息一下！"
	}

	str:=string(data)

	str=strings.Replace(str,"[name]","",1)
	str=strings.Replace(str,"[cqname]","",1)
	return str
 }