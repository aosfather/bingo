package bingo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Application struct {
	Name    string
	config  map[string]string
	router  defaultRouter
	factory *SessionFactory
	port    int
	logfactory *LogFactory
}

func (this *Application) Validate(obj interface{}) []BingoError {
	if this.router.validates.factory == nil {
		this.router.validates.Init(&defaultValidaterFactory{})
	}
	return this.router.validates.Validate(obj)

}

func (this *Application) AddInterceptor(h CustomHandlerInterceptor) {
	if h != nil {
		this.router.interceptor.addInterceptor(h)
	}

}

func (this *Application) GetSession() *TxSession {
	if this.factory != nil {
		return this.factory.GetSession()
	}
	return nil
}

func (this *Application)GetLog(module string)Log {
	return this.logfactory.GetLog(module)
}

//不能获取bingo自身的属性，只能获取应用自身的扩展属性
func (this *Application) GetPropertyFromConfig(key string) string {
	if strings.HasPrefix(key, "bingo.") {
		return ""
	}
	return this.getProperty(key)
}

func (this *Application) getProperty(key string) string {
	if this.config == nil {
		return ""
	}
	return this.config[key]
}

func (this *Application) AddHandler(url string, handler HttpMethodHandler) {
	var rule routerRule
	rule.url = url
	rule.methodHandler = handler
	this.router.addRouter(&rule)
}

func (this *Application) init() {
	if this.config == nil {
		this.config = make(map[string]string)
	}
	if this.config["bingo.system.usedb"] == "true" {
		this.logfactory.write(LEVEL_INFO,"bingo","init db")

		var sqlfactory SessionFactory
		sqlfactory.DBtype = this.config["bingo.db.type"]
		sqlfactory.DBname = this.config["bingo.db.name"]
		sqlfactory.DBurl = this.config["bingo.db.url"]
		sqlfactory.DBuser = this.config["bingo.db.user"]
		sqlfactory.DBpassword = this.config["bingo.db.password"]
		sqlfactory.Init()

		this.factory = &sqlfactory
	}
	this.router.Init(this.factory)
	this.port = 8990
	if this.config["bingo.system.port"] != "" {
		port, err := strconv.Atoi(this.config["bingo.system.port"])
		if err == nil {
			this.port = port
		}

	}
	this.router.setTemplateRoot(this.config["bingo.mvc.template"])
	//设置静态处理
	this.router.staticHandler = &staticController{staticDir: this.config["bingo.mvc.static"],log:this.logfactory.GetLog("bingo.static")}


}

func (this *Application) Run() {
    defer this.logfactory.Close()
	this.init()
	http.ListenAndServe(":"+strconv.Itoa(this.port), &this.router)
}

func (this *Application) Load(file string) {
	if file != "" && isFileExist(file) {
		f, err := os.Open(file)
		if err == nil {
			txt, _ := ioutil.ReadAll(f)
			json.Unmarshal(txt, &this.config)
		}

	}
	this.logfactory=&LogFactory{}
	this.logfactory.SetConfig(LogConfig{true,this.config["bingo.log.file"]})
	this.router.logger=this.logfactory.GetLog("bingo.router")
}
