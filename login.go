package main

import (
	"github.com/aosfather/bingo_mvc"
	"github.com/aosfather/bingo_utils/redis"
	"io"
	"strings"
)

/**
登录处理
1、鉴权白名单，对于特定的页面不做鉴权检查
2、权限检查

*/
type Login interface {
	DoLogin(l *UserLogin) error
}

const _login_html = `<head>
<!-- 以下方式定时转到其他页面 -->
<meta http-equiv="refresh" content="1;url=/login"> 
</head>`

type UserLogin struct {
	SessionId string //
	UserName  string `json:"username"`
	PassWord  string `json:"password"`
	Name      string
}

//登录鉴权
type LoginAccess struct {
	RedisAddr  string                    `Value:"cm.session.addr"`
	RedisDB    int                       `Value:"cm.session.db"`
	RedisPwd   string                    `Value:"cm.session.pwd"`
	Expire     int64                     `Value:"cm.session.expire"`
	Url        string                    `Value:"cm.ssourl" mapper:"name(login);url(/login);method(GET);style(HTML)"`
	SessionMan *bingo_mvc.SessionManager `mapper:"name(loginAction);url(/dologin);method(POST);style(JSON)"`
	LoginFace  Login
}

func (this *LoginAccess) Init() {
	this.SessionMan = new(bingo_mvc.SessionManager)
	this.SessionMan.CookieName = "_login"
	redissession := redis.RedisSessionStore{}
	addrs := strings.Split(this.RedisAddr, ",")
	redissession.InitByCluster(addrs, this.RedisPwd, this.Expire)
	this.SessionMan.SetStore(&redissession)
	this.SessionMan.Init()
}

func (this *LoginAccess) GetHandles() bingo_mvc.HandleMap {
	result := bingo_mvc.NewHandleMap()
	result.Add("login", this.Index, this)
	result.Add("loginAction", this.Login, &UserLogin{})
	result.Add("logoutAction", this.Logout, &UserLogin{})
	return result
}
func (this *LoginAccess) Index(a interface{}) interface{} {
	return bingo_mvc.ModelView{"login", nil}
}

func (this *LoginAccess) Login(a interface{}) interface{} {
	u := a.(*UserLogin)
	debug("login ", u.Name)
	r := Result{}
	if this.LoginFace != nil {
		err := this.LoginFace.DoLogin(u)
		if err != nil {
			errs("error", err.Error())
			r.Code = 400
			r.Msg = err.Error()
			r.Data = nil
		} else {
			r.Code = 200
			r.Msg = "成功"
			this.SetValue(u.SessionId, "name", u.Name)
			this.SetValue(u.SessionId, "user", u.UserName)
			r.Data = []interface{}{u}
		}
	} else {
		this.SetValue(u.SessionId, "user", u.UserName)
	}

	info(r)
	return &r
}

func (this *LoginAccess) Logout(a interface{}) interface{} {
	u := a.(*UserLogin)
	debug("logout ", u)
	this.SetValue(u.SessionId, "user", nil)
	this.SetValue(u.SessionId, "name", nil)
	return _login_html
}

func (this *LoginAccess) PreHandle(writer io.Writer, context bingo_mvc.HttpContext) bool {

	session := this.SessionMan.GetSession(context)
	url := context.GetRequestURI()
	debug("the url:", url)
	if url == "/login" || url == "/dologin" || url == "/health" {
		return true
	}
	if session.IsNew() || session.GetValue("user") == nil {
		writer.Write([]byte(_login_html))
		return false
	}

	return true
}
func (this *LoginAccess) InputProcess(context bingo_mvc.HttpContext, input interface{}) error {
	debug("inputprocess....")
	if u, ok := input.(*UserLogin); ok {
		session := this.SessionMan.GetSession(context)
		u.SessionId = session.ID()
		if u.UserName == "" {
			u.UserName = session.GetValue("user").(string)
			u.Name = session.GetValue("name").(string)
		}

	}

	if m, ok := input.(map[string]interface{}); ok {
		session := this.SessionMan.GetSession(context)
		m["_session_"] = session.ID()
		m["_user_"] = session.GetValue("user")
		m["_username_"] = session.GetValue("name")
	}
	debug(input)
	return nil
}

func (this *LoginAccess) SetValue(sessionId string, key string, value interface{}) {
	session := this.SessionMan.GetSessionById(sessionId)
	if session != nil {
		session.SetValue(key, value)
	}
}
func (this *LoginAccess) PostHandle(writer io.Writer, context bingo_mvc.HttpContext, mv *bingo_mvc.ModelView) bingo_mvc.BingoError {
	return nil
}
func (this *LoginAccess) AfterCompletion(writer io.Writer, context bingo_mvc.HttpContext, err bingo_mvc.BingoError) bingo_mvc.BingoError {
	return nil
}
