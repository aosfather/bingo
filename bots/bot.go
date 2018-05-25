package bots

import (
	"os/exec"
	"github.com/aosfather/bingo/openapi"
	"strings"
)

var(
	_tulingSdk *openapi.TulingSDK
)

func init(){
	_tulingSdk=&openapi.TulingSDK{"808811ad0fd34abaa6fe800b44a9556a"}
}
/**
  bot 定义
 */
//简单bot执行任务
type BotDoAction func(paramters...string) (string,error)


//-----------------简单的系统bots实现----------------//
/*
执行系统命令
  第一个参数，执行的目录
  第二个参数，命令
  第三个及以后的参数，命令相关参数
 */
func RunCmdBot(paramters...string) (string,error){
	cmd := exec.Command(paramters[1], paramters[2:len(paramters)]...)
	cmd.Dir=paramters[0]
	out, err := cmd.Output()
	if err!=nil {
		return err.Error(),err
	}

	return string(out),nil
}

//有道翻译
func RunYoudaoQueryBot(paramters...string) (string,error){
	return openapi.QueryFromYoudaoAsString(strings.Join(paramters," ")),nil
}

//茉莉聊天bot
func RunMoliTalkBot(paramters...string) (string,error){
	//第一个参数是对象，后续的是用户交谈的内容
	return openapi.QueryByMoli(strings.Join(paramters[1:len(paramters)]," ")),nil
}

//图灵聊天bot
func RunTulingTalkBot(paramters...string) (string,error){
	return _tulingSdk.QueryAsString(paramters[0],strings.Join(paramters[1:len(paramters)]," ")),nil
}