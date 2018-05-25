package bot

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"github.com/aosfather/bingo/utils"
)

type Logger struct {
}

func (this *Logger) Panic(f string, v ...interface{}) {
	fmt.Printf(f, v...)
}

// Fatal without custom Fields.
func (this *Logger) Fatal(f string, v ...interface{}) {
	fmt.Printf(f, v...)
}

// Error without custom Fields.
func (this *Logger) Error(f string, v ...interface{}) {
	fmt.Printf(f, v...)
}

// Warn without custom Fields.
func (this *Logger) Warn(f string, v ...interface{}) {
	fmt.Printf(f, v...)
}

// Info without custom Fields.
func (this *Logger) Info(f string, v ...interface{}) {
	fmt.Printf(f, v...)
}

// Debug without custom Fields.
func Debug(f string, v ...interface{}) {
	fmt.Printf(f, v...)
}

type Wecat struct {
	uuid        string
	baseURI     string
	redirectURI string
	loginRes    LoginResult
	deviceID    string
	syncKey     SyncKey
	user        User
	baseRequest map[string]interface{}
	syncHost    string
	client      *http.Client
	auto        bool
	showRebot   bool
	contacts    map[string]Contact
	robotName   string
	currentBot  TalkBot
	currentBotIndex int
}

const (
	LoginBaseURL = "https://login.weixin.qq.com"
	WxReferer    = "https://wx.qq.com/"
	WxUserAgent  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36"
)

var (
	log   = Logger{}
	Hosts = []string{
		"webpush.wx.qq.com",
		"webpush2.wx.qq.com",
		"webpush.wechat.com",
		"webpush1.wechat.com",
		"webpush2.wechat.com",
		"webpush1.wechatapp.com",
	}
)

func NewWecat() (*Wecat, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Error("get cookiejar fail", err)
		return nil, err
	}

	client := &http.Client{
		CheckRedirect: nil,
		Jar:           jar,
	}

	rand.Seed(time.Now().Unix())
	randID := strconv.Itoa(rand.Int())
	return &Wecat{
		client:      client,
		deviceID:    "e" + randID[2:17],
		baseRequest: make(map[string]interface{}),
		contacts:    make(map[string]Contact),
		auto:        true,
		robotName:   "图小二",
	}, nil
}

func (w *Wecat) GetUUID() error {
	if w.uuid != "" {
		return nil
	}

	uri := LoginBaseURL + "/jslogin?appid=wx782c26e4c19acffb&fun=new&lang=zh_CN&_=" + w.timestamp()
	//result: window.QRLogin.code = 200; window.QRLogin.uuid = "xxx"; //wx782c26e4c19acffb  wxeb7ec651dd0aefa9
	data, err := w.get(uri)
	if err != nil {
		log.Error("get uuid fail", err)
		return err
	}

	res := make(map[string]string)
	datas := strings.Split(string(data), ";")
	for _, d := range datas {
		kvs := strings.Split(d, " = ")
		if len(kvs) == 2 {
			res[strings.Trim(kvs[0], " ")] = strings.Trim(strings.Trim(kvs[1], " "), "\"")
		}
	}
	if res["window.QRLogin.code"] == "200" {
		if uuid, ok := res["window.QRLogin.uuid"]; ok {
			w.uuid = uuid
			return nil
		}
	}

	return fmt.Errorf(string(data))
}

func (w *Wecat) GenQrcode() error {
	if w.uuid == "" {
		err := errors.New("haven't get uuid")
		log.Error("gen qrcode fail", err)
		return err
	}

	uri := LoginBaseURL + "/qrcode/" + w.uuid + "?t=webwx&_=" + w.timestamp()
	fmt.Println(uri)
	//resp, err := w.get(uri)
	//qr := qrcode.NewQRCode("", false)
	//
	//img, err := jpeg.Decode(bytes.NewReader([]byte(resp)))
	//if err != nil {
	//	return err
	//}
	//
	//if err := qr.SetImage(img); err != nil {
	//	return err
	//}
	//qr.Output()

	return nil
}

func (w *Wecat) Login() error {
	tip := 1
	for {
		uri := fmt.Sprintf("%s/cgi-bin/mmwebwx-bin/login?tip=%d&uuid=%s&_=%s", LoginBaseURL, tip, w.uuid, w.timestamp())
		data, err := w.get(uri)
		if err != nil {
			return err
		}

		re := regexp.MustCompile(`window.code=(\d+);`)
		codes := re.FindStringSubmatch(string(data))
		if len(codes) > 1 {
			code := codes[1]
			switch code {
			case "201":
				log.Info("scan code success")
				tip = 0
			case "200":
				log.Info("login success, wait to redirect")
				re := regexp.MustCompile(`window.redirect_uri="(\S+?)";`)
				redirctURIs := re.FindStringSubmatch(string(data))

				if len(redirctURIs) > 1 {
					redirctURI := redirctURIs[1] + "&fun=new"
					w.redirectURI = redirctURI
					re = regexp.MustCompile(`/`)
					baseURIs := re.FindAllStringIndex(redirctURI, -1)
					w.baseURI = redirctURI[:baseURIs[len(baseURIs)-1][0]]
					if err := w.redirect(); err != nil {
						log.Error("%s", err)
						return err
					}
					return nil
				}

				log.Error("get redirct URL fail")

			case "408":
				log.Error("login timeout")
			default:
				log.Error("login fail")
			}
		} else {
			return errors.New("get code fail")
		}

		time.Sleep(time.Second * time.Duration(2))
	}
}

func (w *Wecat) redirect() error {
	data, err := w.get(w.redirectURI)
	if err != nil {
		log.Error("redirct fail", err)
		return err
	}

	var lr LoginResult
	if err = xml.Unmarshal(data, &lr); err != nil {
		log.Error("unmarshal fail", err)
		return err
	}

	w.loginRes = lr
	w.baseRequest["Uin"] = utils.Int64(lr.Wxuin)
	w.baseRequest["Sid"] = lr.Wxsid
	w.baseRequest["Skey"] = lr.Skey
	w.baseRequest["DeviceID"] = w.deviceID
	return nil
}

func (w *Wecat) Init() error {
	uri := fmt.Sprintf("%s/webwxinit?pass_ticket=%s&skey=%s&r=%s", w.baseURI, w.loginRes.PassTicket, w.loginRes.Skey, w.timestamp())
	params := make(map[string]interface{})
	params["BaseRequest"] = w.baseRequest
	data, err := w.post(uri, params)
	if err != nil {
		log.Error("init post fail", err)
		return err
	}

	var res InitResult
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	w.user = res.User
	w.syncKey = res.SyncKey

	if res.BaseResponse.Ret != 0 {
		fmt.Println(string(data))
		log.Error("init fail ret <> 0")
	}

	return nil
}

func (w *Wecat) strSyncKey() string {
	kvs := []string{}
	for _, list := range w.syncKey.List {
		kvs = append(kvs, utils.String(list.Key)+"_"+utils.String(list.Val))
	}

	return strings.Join(kvs, "|")
}

func (w *Wecat) SyncCheck() (retcode, selector int) {
	for _, host := range Hosts {
		uri := fmt.Sprintf("https://%s/cgi-bin/mmwebwx-bin/synccheck", host)
		v := url.Values{}
		v.Add("r", w.timestamp())
		v.Add("sid", w.loginRes.Wxsid)
		v.Add("uin", w.loginRes.Wxuin)
		v.Add("skey", w.loginRes.Skey)
		v.Add("deviceid", w.deviceID)
		v.Add("synckey", w.strSyncKey())
		v.Add("_", w.timestamp())
		uri = uri + "?" + v.Encode()

		data, err := w.get(uri)
		if err != nil {
			//log.Error("sync check fail", err)
			continue
		}

		re := regexp.MustCompile(`window.synccheck={retcode:"(\d+)",selector:"(\d+)"}`)
		codes := re.FindStringSubmatch(string(data))
		if len(codes) > 2 {
			return utils.Int(codes[1]), utils.Int(codes[2])
		}
	}

	return 9999, 0
}

func (w *Wecat) StatusNotify() error {
	uri := fmt.Sprintf("%s/webwxstatusnotify?lang=zh_CN&pass_ticket=%s", w.baseURI, w.loginRes.PassTicket)
	params := make(map[string]interface{})
	params["BaseRequest"] = w.baseRequest
	params["Code"] = 3
	params["FromUserName"] = w.user.UserName
	params["ToUserName"] = w.user.UserName
	params["ClientMsgId"] = int(time.Now().Unix())
	data, err := w.post(uri, params)
	if err != nil {
		return err
	}

	var res StatusNotifyResult

	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
    fmt.Println(string(data))
	if res.BaseResponse.Ret != 0 {
		return fmt.Errorf("%s", res.BaseResponse.ErrMsg)
	}
	return nil
}

func (w *Wecat) GetContact() error {
	uri := fmt.Sprintf("%s/webwxgetcontact?sid=%s&skey=%s&pass_ticket=%s", w.baseURI, w.loginRes.Wxsid, w.loginRes.Skey, w.loginRes.PassTicket)
	params := make(map[string]interface{})
	params["BaseRequest"] = w.baseRequest

	data, err := w.post(uri, params)
	if err != nil {
		return err
	}

	var contacts Contacts
	if err := json.Unmarshal(data, &contacts); err != nil {
		return err
	}

	for _, contact := range contacts.MemberList {
		if contact.NickName == "" {
			contact.NickName = contact.UserName
		}
		w.contacts[contact.UserName] = contact
	}

	return nil
}

func (w *Wecat) WxSync() (*Message, error) {
	uri := fmt.Sprintf("%s/webwxsync?sid=%s&skey=%s&pass_ticket=%s", w.baseURI, w.loginRes.Wxsid, w.loginRes.Skey, w.loginRes.PassTicket)
	params := make(map[string]interface{})
	params["BaseRequest"] = w.baseRequest
	params["SyncKey"] = w.syncKey
	params["rr"] = ^int(time.Now().Unix())

	data, err := w.post(uri, params)
	if err != nil {
		return nil, err
	}

	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	if msg.BaseResponse.Ret == 0 {
		w.syncKey = msg.SyncKey
	}
	//TODO
	return &msg, nil
}

func (w *Wecat) run(desc string, f func() error) {
	start := time.Now()
	log.Info(desc+"\n")
	if err := f(); err != nil {
		log.Error("FAIL, exit now", err)
		os.Exit(1)
	}

	log.Info("--> SUCCESS, use time ", strconv.FormatInt(time.Now().Sub(start).Nanoseconds(),10),"\n")
}



func (w *Wecat) getNickName(userName string) string {
	if v, ok := w.contacts[userName]; ok {
		return v.NickName
	}

	return userName
}

func (w *Wecat) handle(msg *Message) error {
	for _, contact := range msg.ModContactList {
		if _, ok := w.contacts[contact.UserName]; !ok {
			if contact.NickName == "" {
				contact.NickName = contact.UserName
			}
			w.contacts[contact.UserName] = contact
		}
	}

	for _, m := range msg.AddMsgList {
		fmt.Println(m)
		fmt.Println(w.user.NickName, w.user.RemarkName)
		m.Content = strings.Replace(m.Content, "&lt;", "<", -1)
		m.Content = strings.Replace(m.Content, "&gt;", ">", -1)
		switch m.MsgType {
		case 1:
			if m.FromUserName[:2] == "@@" { //群消息
				w.doGroupReply(m.Content, m.FromUserName)
			} else {
				if m.FromUserName != w.user.UserName {
					w.doAutoReply(m.Content, m.FromUserName)
				} else {
					w.doCommand(m.Content, m.ToUserName)
				}
			}
		case 51:
			log.Info("sync ok->")
		}
	}

	return nil
}

func (w *Wecat) Dail() error {
	for {
		retcode, selector := w.SyncCheck()
		switch retcode {
		case 1100:
			log.Info("logout with phone, bye")
			return nil
		case 1101:
			log.Info("login web wecat at other palce, bye")
			return nil
		case 0:
			switch selector {
			case 2:
				msg, err := w.WxSync()
				if err != nil {
					log.Error("%s", err)
				}

				if err := w.handle(msg); err != nil {
					log.Error("%s", err)
				}
			case 0:
				time.Sleep(time.Second)
			case 6, 4:
				w.WxSync()
				time.Sleep(time.Second)
			}
		default:
			log.Warn("unknow code", retcode)
		}
	}
}

func (w *Wecat) Start() {
	w.run("[*] get uuid ...", w.GetUUID)
	w.run("[*] generate qrcode ...", w.GenQrcode)
	w.run("[*] login ...", w.Login)
	w.run("[*] init wecat ...", w.Init)
	w.run("[*] open status notify ...", w.StatusNotify)
	w.run("[*] get contact ...", w.GetContact)
	w.run("[*] dail sync message ...", w.Dail)
}

func (w *Wecat) timestamp() string {
	return utils.String(time.Now().Unix())
}

func (w *Wecat) get(uri string) ([]byte, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Referer", WxReferer)
	req.Header.Add("User-agent", WxUserAgent)

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (w *Wecat) post(uri string, params map[string]interface{}) ([]byte, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Referer", WxReferer)
	req.Header.Add("User-agent", WxUserAgent)

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
