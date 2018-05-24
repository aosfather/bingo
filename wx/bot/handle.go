package bot

import (
	"fmt"
	"strings"
)

//执行指令
func (w *Wecat) doCommand(c string,to string) {
	if strings.HasPrefix(c,"#") {
	c=c[1:len(c)]
	switch c {
	case "退下":
		w.auto = false
	case "来人":
		w.auto = true
		fmt.Println("set auto",w.auto)
	case "显示":
		w.showRebot = true
		fmt.Println("set showrebot ",w.showRebot)
	case "隐身":
		w.showRebot = false
	case "换人":
		w.currentBotIndex++
		if w.currentBotIndex >=len(botsMap) {
			w.currentBotIndex=0
		}
		w.currentBot=botsMap[w.currentBotIndex]
	default:
		parames:=strings.Split(c," ")
		reply:=w.callBot(parames...)
		if reply=="" {
			fmt.Println("[unknown command] ", w.user.UserName, ": ", c)
			w.doAutoReply(c,to)
		}else {
			w.SendMessage(reply,to)
		}

	}

	}else {
		//其它时候交谈
		if w.isAtme(c) {
			content := strings.Replace(c, "@"+w.user.NickName, "", -1)
			content = strings.Replace(content, "@"+w.robotName, "", -1)
			w.doAutoReply(content,to)
		}

	}
}


func(w*Wecat) callBot(parames...string) string {
	size:=len(parames)
	if size>0 {
        bot:=bots[parames[0]]
        if bot!=nil {
        	if size>1 {
				return bot.DoAction(parames[1:len(parames)]...)
			}else {
				return bot.DoAction()
			}
		}
	}
	return ""
}

func(w *Wecat)doAutoReply(c string,fromUserName string) {
	fmt.Println("[*] ", w.getNickName(fromUserName), ": ", c)
	if w.auto {
		fmt.Println("auto reply")
		reply, err := w.getReply(c, fromUserName)
		if err != nil {
			fmt.Println(err)
			return
		}

		if w.showRebot {
			reply = w.robotName + ": " + reply
		}
		if err := w.SendMessage(reply, fromUserName); err != nil {
			fmt.Println("send error ", err)
			return
		}
		fmt.Println("[->#] ", w.user.NickName, ": ", reply)
	}

}


func (w *Wecat)doGroupReply(c string,fromUserName string){
	contents := strings.Split(c, ":<br/>")
	content := contents[1]
    fmt.Println(content)

	if w.isAtme(content) {
		content = strings.Replace(content, "@"+w.user.NickName, "", -1)
		content = strings.Replace(content, "@"+w.user.RemarkName, "", -1)
		content = strings.Replace(content, "@"+w.robotName, "", -1)
        w.doAutoReply(content,fromUserName)
	} else {
		fmt.Println(contents[0])
		fmt.Println("[**]", w.getNickName(contents[0]), ": ", contents[1])



	}
}

//判断是否at了我，现在的
func (w *Wecat)isAtme(content string) bool {
	if (strings.Contains(content, "@"+w.robotName)||strings.Contains(content,"@"+w.user.NickName))  {
			return true
	}
	return false
}


func (w *Wecat) getReply(msg string, uid string) (string, error) {
	if w.currentBot==nil {
	  w.currentBotIndex=0
      w.currentBot=botsMap[w.currentBotIndex]
	}
     return w.currentBot.Reply(uid,msg),nil
}

