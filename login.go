package main

import (
	"fmt"
	"github.com/aosfather/bingo_mvc"
	"github.com/aosfather/bingo_mvc/sqltemplate"
	"github.com/aosfather/bingo_utils/codes"
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
<!-- 以下方式定时转到其他页面,如果不是顶层，就让父页面跳转 -->
<meta http-equiv="refresh" content="1;url=/login"> 
<script language="JavaScript">
if (window != top)
top.location.href = location.href;
</script>
</head>`

type UserLogin struct {
	SessionId string //
	UserName  string `json:"username"`
	PassWord  string `json:"password"`
	Name      string
	Roles     string
}

//登录鉴权
type LoginAccess struct {
	RedisAddr  string                    `Value:"session.addr"`
	RedisDB    int                       `Value:"session.db"`
	RedisPwd   string                    `Value:"session.pwd"`
	Expire     int64                     `Value:"session.expire"`
	Url        string                    `Value:"session.login" mapper:"name(login);url(/login);method(GET);style(HTML)"`
	SessionMan *bingo_mvc.SessionManager `mapper:"name(loginAction);url(/dologin);method(POST);style(JSON)"`
	LoginFace  Login                     `Inject:"" mapper:"name(logoutAction);url(/dologout);method(GET);style(HTML)"`
	Exemption  string                    `Value:"session.exemption"`
	exemptions []string
	Title      string `Value:"app.title"`
}

func (this *LoginAccess) Init() {
	this.SessionMan = new(bingo_mvc.SessionManager)
	this.SessionMan.CookieName = "_login"
	if this.RedisAddr != "" {
		redissession := redis.RedisSessionStore{}
		addrs := strings.Split(this.RedisAddr, ",")
		redissession.InitByCluster(addrs, this.RedisPwd, this.Expire)
		this.SessionMan.SetStore(&redissession)
	}
	if this.Exemption != "" {
		this.exemptions = strings.Split(this.Exemption, ";")
	}
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
	data := make(map[string]string)
	if this.Title != "" {
		data["Title"] = this.Title
	} else {
		data["Title"] = "Bingo 管理平台"
	}

	return bingo_mvc.ModelView{"login", data}
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
			this.SetValue(u.SessionId, "roles", u.Roles)
			r.Data = []interface{}{u}
		}
	} else {
		this.SetValue(u.SessionId, "user", u.UserName)
		this.SetValue(u.SessionId, "name", u.UserName)
		r.Code = 200
		r.Msg = "成功"
		r.Data = []interface{}{u}
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

//访问的url是否属于登录鉴权豁免地址
func (this *LoginAccess) isExemptioned(url string) bool {
	if url == "/login" || url == "/dologin" {
		return true
	}
	for _, u := range this.exemptions {
		if url == u {
			return true
		}
	}
	return false
}

func (this *LoginAccess) PreHandle(writer io.Writer, context bingo_mvc.HttpContext) bool {

	session := this.SessionMan.GetSession(context)
	url := context.GetRequestURI()
	debug("the url:", url)
	if this.isExemptioned(url) {
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

	} else if u, ok := input.(*FormRequest); ok {
		session := this.SessionMan.GetSession(context)
		u._Roles = session.GetValue("roles").(string)
	} else if m, ok := input.(map[string]interface{}); ok {
		session := this.SessionMan.GetSession(context)
		m["_session_"] = session.ID()
		m["_user_"] = session.GetValue("user")
		m["_username_"] = session.GetValue("name")
		m["_roles_"] = session.GetValue("roles")
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

//默认的登录实现
type DefaultLogin struct {
	Dao  *sqltemplate.MapperDao `Inject:"user"`
	Salt string                 `Value:"session.salt"`
}

func (this *DefaultLogin) Init() {
	if this.Salt == "" {
		this.Salt = "$$##_default_salt_##$$"
	}
}
func (this *DefaultLogin) DoLogin(l *UserLogin) error {
	if l.UserName == "" || l.PassWord == "" {
		return fmt.Errorf("password or name error!")
	}

	u := &User{}
	u.Id = l.UserName
	debug(u)
	exist := this.Dao.Find(u, "findbyid")
	debug(exist)
	if !exist {
		return fmt.Errorf("password or name error!")
	}
	debug(u)
	pwd := codes.ToMd5Hex(fmt.Sprintf("%s##%s", l.PassWord, this.Salt))
	debug(pwd)
	if u.Pwd == pwd {
		l.Name = u.Name
		l.Roles = u.Roles
	} else {
		return fmt.Errorf("password or name error!")
	}

	return nil
}

//db 中的User对象
type User struct {
	Id     string `Table:"bingo_user"`
	Name   string
	Depart string
	Pwd    string
	Roles  string
}
