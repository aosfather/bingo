package bot

import (
	"github.com/aosfather/bingo/openapi"
	"fmt"
	"strings"
	"os/exec"
)

var (
	botsMap=[]TalkBot{}
	bots=make(map[string]Bot)
)

func init(){
	fmt.Println("init bots")
	botsMap=append(botsMap,&TulingBot{&openapi.TulingSDK{"808811ad0fd34abaa6fe800b44a9556a"}})
	botsMap=append(botsMap,&MoliBot{})
	bots["翻译"]=&YoudaoBot{}
	bots["执行"]=&CmdBot{}
}
type Bot interface {
    DoAction(paramters...string) string
}

type TalkBot interface {
	Reply(user,msg string) string
}


type TulingBot struct {
	sdk *openapi.TulingSDK
}

func (this *TulingBot) Reply(user,msg string) string{
	if this.sdk!=nil {
		return this.sdk.QueryAsString(user,msg)
	}

	return "[自动回复] 暂时不在线！"

}

//茉莉机器人
type MoliBot struct {

}

func (this *MoliBot) Reply(user,msg string) string{
	return openapi.QueryByMoli(msg)
}

//有道翻译机器人
type YoudaoBot struct {

}
func (this *YoudaoBot)DoAction(paramters...string) string{
	return openapi.QueryFromYoudaoAsString(strings.Join(paramters," "))
}

//命令执行bot
type CmdBot struct {

}

func (this *CmdBot)DoAction(paramters...string) string{
	cmd := exec.Command(paramters[0], paramters[1:len(paramters)]...)
	out, err := cmd.Output()
	if err!=nil {
		return err.Error()
	}

	return string(out)
}