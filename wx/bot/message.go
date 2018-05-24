package bot
/**
消息
1、撤回消息--revoker
2、发送文本
3、发送图片
4、发送语音
5、发送视频
 */
import (
	"fmt"
	"strconv"
	"math/rand"
)




func (w *Wecat) SendMessage(message string, to string) error {
	uri := fmt.Sprintf("%s/webwxsendmsg?pass_ticket=%s", w.baseURI, w.loginRes.PassTicket)
	clientMsgID := w.timestamp() + "0" + strconv.Itoa(rand.Int())[3:6]
	params := make(map[string]interface{})
	params["BaseRequest"] = w.baseRequest
	msg := make(map[string]interface{})
	msg["Type"] = 1
	msg["Content"] = message
	msg["FromUserName"] = w.user.UserName
	msg["ToUserName"] = to
	msg["LocalID"] = clientMsgID
	msg["ClientMsgId"] = clientMsgID
	params["Msg"] = msg
	_, err := w.post(uri, params)
	if err != nil {
		return err
	}

	return nil
}